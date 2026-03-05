package discovery

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/rahulkj/cp-discovery/internal/model"
)

func DiscoverKafka(config model.KafkaConfig, detailed bool) (model.KafkaReport, error) {
	report := model.KafkaReport{
		Available: false,
		Topics:    make([]model.TopicInfo, 0),
		Brokers:   make([]model.BrokerInfo, 0),
	}

	if config.BootstrapServers == "" {
		return report, fmt.Errorf("bootstrap servers not configured")
	}

	// Parse bootstrap servers
	brokers := strings.Split(config.BootstrapServers, ",")

	// Create transport with security settings
	transport, err := createTransport(config)
	if err != nil {
		return report, fmt.Errorf("creating transport: %w", err)
	}

	// Create client
	client := &kafka.Client{
		Addr:      kafka.TCP(brokers...),
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get cluster metadata
	metadata, err := client.Metadata(ctx, &kafka.MetadataRequest{})
	if err != nil {
		return report, fmt.Errorf("getting metadata: %w", err)
	}

	report.Available = true
	report.BrokerCount = len(metadata.Brokers)

	// Get broker information
	for _, broker := range metadata.Brokers {
		host, portStr, _ := net.SplitHostPort(broker.Host)
		port, _ := strconv.Atoi(portStr)
		if host == "" {
			host = broker.Host
		}

		brokerInfo := model.BrokerInfo{
			ID:   broker.ID,
			Host: host,
			Port: port,
		}
		report.Brokers = append(report.Brokers, brokerInfo)
	}

	// Detect controller type (ZooKeeper vs KRaft)
	report.ControllerType = detectControllerTypeFromBrokers(metadata.Brokers)

	// Estimate controller count based on cluster type
	report.ControllerCount = estimateControllerCount(report.ControllerType, report.BrokerCount)

	// Get topic information
	if len(metadata.Topics) > 0 {
		report.TopicCount = len(metadata.Topics)
		report.InternalTopics = 0
		report.ExternalTopics = 0
		report.TotalPartitions = 0

		// Categorize topics as internal or external
		for _, topic := range metadata.Topics {
			if isInternalTopic(topic.Name) {
				report.InternalTopics++
			} else {
				report.ExternalTopics++
			}

			partitionCount := len(topic.Partitions)
			report.TotalPartitions += partitionCount

			// Get replication factor from first partition
			replicationFactor := 0
			if len(topic.Partitions) > 0 {
				replicationFactor = len(topic.Partitions[0].Replicas)
			}

			// Always collect topic information
			topicInfo := model.TopicInfo{
				Name:              topic.Name,
				IsInternal:        isInternalTopic(topic.Name),
				Partitions:        partitionCount,
				ReplicationFactor: replicationFactor,
				AssociatedSchemas: make([]string, 0),
			}

			// Get topic configuration for retention settings
			if configs, err := getTopicConfigs(client, ctx, topic.Name); err == nil {
				topicInfo.RetentionMs = configs.RetentionMs
				topicInfo.RetentionBytes = configs.RetentionBytes
			}

			// Calculate topic storage size by querying partition offsets
			topicSize, err := calculateTopicSize(client, ctx, topic.Name, topic.Partitions)
			if err == nil {
				topicInfo.SizeBytes = topicSize
			}

			// Throughput will be populated from metrics if available
			topicInfo.ThroughputBytesInPerSec = 0
			topicInfo.ThroughputBytesOutPerSec = 0

			report.Topics = append(report.Topics, topicInfo)
		}
	}

	// Get security configuration
	report.SecurityConfig = getKafkaSecurityConfig(config)

	// Calculate cluster metrics
	report.ClusterMetrics = calculateClusterMetrics(report)

	return report, nil
}

func createTransport(config model.KafkaConfig) (*kafka.Transport, error) {
	transport := &kafka.Transport{
		DialTimeout: 10 * time.Second,
		IdleTimeout: 30 * time.Second,
	}

	// Configure SASL
	if config.SaslMechanism != "" {
		var mechanism sasl.Mechanism
		var err error

		switch config.SaslMechanism {
		case "PLAIN":
			mechanism = plain.Mechanism{
				Username: config.SaslUsername,
				Password: config.SaslPassword,
			}
		case "SCRAM-SHA-256":
			mechanism, err = scram.Mechanism(scram.SHA256, config.SaslUsername, config.SaslPassword)
			if err != nil {
				return nil, fmt.Errorf("creating SCRAM-SHA-256 mechanism: %w", err)
			}
		case "SCRAM-SHA-512":
			mechanism, err = scram.Mechanism(scram.SHA512, config.SaslUsername, config.SaslPassword)
			if err != nil {
				return nil, fmt.Errorf("creating SCRAM-SHA-512 mechanism: %w", err)
			}
		default:
			return nil, fmt.Errorf("unsupported SASL mechanism: %s", config.SaslMechanism)
		}

		transport.SASL = mechanism
	}

	// Configure SSL/TLS
	if config.SecurityProtocol != "" && strings.Contains(config.SecurityProtocol, "SSL") {
		tlsConfig := &tls.Config{}

		// Load CA certificate if provided
		if config.SslCaLocation != "" {
			caCert, err := os.ReadFile(config.SslCaLocation)
			if err != nil {
				return nil, fmt.Errorf("reading CA cert: %w", err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}

		// Load client certificate if provided
		if config.SslCertLocation != "" && config.SslKeyLocation != "" {
			cert, err := tls.LoadX509KeyPair(config.SslCertLocation, config.SslKeyLocation)
			if err != nil {
				return nil, fmt.Errorf("loading client cert: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// Set endpoint identification
		if config.SslEndpointIdentification == "" || config.SslEndpointIdentification == "https" {
			tlsConfig.InsecureSkipVerify = false
		} else {
			tlsConfig.InsecureSkipVerify = true
		}

		transport.TLS = tlsConfig
	}

	return transport, nil
}

func detectControllerTypeFromBrokers(brokers []kafka.Broker) string {
	if len(brokers) == 0 {
		return "unknown"
	}

	// Collect broker IDs
	brokerIDs := make([]int32, 0, len(brokers))
	for _, broker := range brokers {
		brokerIDs = append(brokerIDs, int32(broker.ID))
	}

	// Method 1: Check if broker IDs are sequential starting from 0 or 1
	isSequentialFrom0 := isSequentialBrokerIDs(brokerIDs, 0)
	isSequentialFrom1 := isSequentialBrokerIDs(brokerIDs, 1)

	if isSequentialFrom0 || isSequentialFrom1 {
		return "kraft"
	}

	// Method 2: Check for broker ID patterns
	minID := brokerIDs[0]
	maxID := brokerIDs[0]
	for _, id := range brokerIDs {
		if id < minID {
			minID = id
		}
		if id > maxID {
			maxID = id
		}
	}

	idRange := maxID - minID

	if (minID == 0 || minID == 1) && maxID < 100 && idRange < 100 {
		return "kraft"
	}

	if maxID < 10 {
		return "kraft"
	}

	hasLargeIDs := maxID > 100

	if !hasLargeIDs && minID <= 10 {
		return "kraft"
	}

	return "zookeeper"
}

func estimateControllerCount(controllerType string, brokerCount int) int {
	if controllerType == "kraft" {
		if brokerCount <= 3 {
			return brokerCount
		} else {
			return 3
		}
	} else if controllerType == "zookeeper" {
		return 1
	}
	return 0
}

func isSequentialBrokerIDs(ids []int32, startFrom int32) bool {
	if len(ids) == 0 {
		return false
	}

	idSet := make(map[int32]bool)
	for _, id := range ids {
		idSet[id] = true
	}

	for i := int32(0); i < int32(len(ids)); i++ {
		expectedID := startFrom + i
		if !idSet[expectedID] {
			return false
		}
	}

	return true
}

type TopicConfigs struct {
	RetentionMs    int64
	RetentionBytes int64
}

func getTopicConfigs(client *kafka.Client, ctx context.Context, topicName string) (TopicConfigs, error) {
	configs := TopicConfigs{
		RetentionMs:    -1,
		RetentionBytes: -1,
	}

	// Describe topic configs
	req := &kafka.DescribeConfigsRequest{
		Resources: []kafka.DescribeConfigRequestResource{
			{
				ResourceType: kafka.ResourceTypeTopic,
				ResourceName: topicName,
			},
		},
	}

	resp, err := client.DescribeConfigs(ctx, req)
	if err != nil {
		return configs, err
	}

	if len(resp.Resources) > 0 {
		for _, entry := range resp.Resources[0].ConfigEntries {
			switch entry.ConfigName {
			case "retention.ms":
				fmt.Sscanf(entry.ConfigValue, "%d", &configs.RetentionMs)
			case "retention.bytes":
				fmt.Sscanf(entry.ConfigValue, "%d", &configs.RetentionBytes)
			}
		}
	}

	return configs, nil
}

func calculateClusterMetrics(report model.KafkaReport) model.ClusterMetrics {
	metrics := model.ClusterMetrics{
		BytesInPerSec:             0,
		BytesOutPerSec:            0,
		MessagesInPerSec:          0,
		TotalDiskUsageBytes:       0,
		UnderReplicatedPartitions: 0,
	}

	// Calculate total cluster storage from topic sizes
	var totalTopicStorage int64 = 0
	for _, topic := range report.Topics {
		totalTopicStorage += topic.SizeBytes
	}
	metrics.TotalDiskUsageBytes = totalTopicStorage

	return metrics
}

// calculateTopicSize estimates topic storage size by querying partition high watermarks
func calculateTopicSize(client *kafka.Client, ctx context.Context, topicName string, partitions []kafka.Partition) (int64, error) {
	var totalSize int64 = 0

	// Query each partition's watermarks
	for _, partition := range partitions {
		// Get high watermark (latest offset)
		highReq := &kafka.ListOffsetsRequest{
			Topics: map[string][]kafka.OffsetRequest{
				topicName: {
					{
						Partition: partition.ID,
						Timestamp: kafka.LastOffset,
					},
				},
			},
		}

		highResp, err := client.ListOffsets(ctx, highReq)
		if err != nil {
			continue
		}

		// Get low watermark (first offset)
		lowReq := &kafka.ListOffsetsRequest{
			Topics: map[string][]kafka.OffsetRequest{
				topicName: {
					{
						Partition: partition.ID,
						Timestamp: kafka.FirstOffset,
					},
				},
			},
		}

		lowResp, err := client.ListOffsets(ctx, lowReq)
		if err != nil {
			continue
		}

		var high, low int64
		if topicOffsets, ok := highResp.Topics[topicName]; ok && len(topicOffsets) > 0 {
			high = topicOffsets[0].LastOffset
		}
		if topicOffsets, ok := lowResp.Topics[topicName]; ok && len(topicOffsets) > 0 {
			low = topicOffsets[0].LastOffset
		}

		// Calculate approximate size based on offset range
		// Average message size assumption: 1KB per message
		offsetRange := high - low
		estimatedSize := offsetRange * 1024

		totalSize += estimatedSize
	}

	return totalSize, nil
}

// Helper function to determine if a topic is internal
func isInternalTopic(topicName string) bool {
	if strings.HasPrefix(topicName, "_") {
		return true
	}

	internalPrefixes := []string{
		"connect-configs",
		"connect-offsets",
		"connect-status",
		"_confluent",
		"default_ksql_processing_log",
	}

	for _, prefix := range internalPrefixes {
		if strings.HasPrefix(topicName, prefix) {
			return true
		}
	}

	return false
}

func getKafkaSecurityConfig(config model.KafkaConfig) model.SecurityConfig {
	secConfig := model.SecurityConfig{
		SaslMechanisms:    make([]string, 0),
		SecurityProtocols: make([]string, 0),
		SslEnabled:        false,
		SaslEnabled:       false,
	}

	// Start with client connection security info
	if config.SecurityProtocol != "" && config.SecurityProtocol != "PLAINTEXT" {
		secConfig.SecurityProtocols = append(secConfig.SecurityProtocols, config.SecurityProtocol)
		if strings.Contains(config.SecurityProtocol, "SSL") {
			secConfig.SslEnabled = true
		}
		if strings.Contains(config.SecurityProtocol, "SASL") {
			secConfig.SaslEnabled = true
		}
	} else {
		secConfig.SecurityProtocols = append(secConfig.SecurityProtocols, "PLAINTEXT")
	}

	if config.SaslMechanism != "" {
		secConfig.SaslMechanisms = append(secConfig.SaslMechanisms, config.SaslMechanism)
		secConfig.SaslEnabled = true
	}

	// Determine authentication method
	secConfig.AuthenticationMethod = determineKafkaAuthMethod(secConfig)

	return secConfig
}

func determineKafkaAuthMethod(config model.SecurityConfig) string {
	if len(config.SaslMechanisms) == 0 && !config.SslEnabled {
		return "PLAINTEXT"
	}

	methods := []string{}

	if config.SslEnabled && !config.SaslEnabled {
		methods = append(methods, "SSL/TLS")
	}

	if config.SaslEnabled {
		if len(config.SaslMechanisms) > 0 {
			for _, mech := range config.SaslMechanisms {
				methods = append(methods, "SASL/"+mech)
			}
		} else {
			methods = append(methods, "SASL")
		}
	}

	if config.SslEnabled && config.SaslEnabled {
		return strings.Join(methods, " + ")
	}

	if len(methods) > 0 {
		return strings.Join(methods, ", ")
	}

	return "PLAINTEXT"
}
