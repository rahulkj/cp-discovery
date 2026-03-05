package discovery

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

	// Create admin client configuration
	adminConfig := kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
	}

	// Add security configuration
	if config.SecurityProtocol != "" && config.SecurityProtocol != "PLAINTEXT" {
		adminConfig["security.protocol"] = config.SecurityProtocol
	}
	if config.SaslMechanism != "" {
		adminConfig["sasl.mechanism"] = config.SaslMechanism
		adminConfig["sasl.username"] = config.SaslUsername
		adminConfig["sasl.password"] = config.SaslPassword
	}

	// Add SSL/TLS configuration
	if config.SslCaLocation != "" {
		adminConfig["ssl.ca.location"] = config.SslCaLocation
	}
	if config.SslCertLocation != "" {
		adminConfig["ssl.certificate.location"] = config.SslCertLocation
	}
	if config.SslKeyLocation != "" {
		adminConfig["ssl.key.location"] = config.SslKeyLocation
	}
	if config.SslKeyPassword != "" {
		adminConfig["ssl.key.password"] = config.SslKeyPassword
	}
	if config.SslEndpointIdentification != "" {
		adminConfig["ssl.endpoint.identification.algorithm"] = config.SslEndpointIdentification
	}

	// Create admin client
	adminClient, err := kafka.NewAdminClient(&adminConfig)
	if err != nil {
		return report, fmt.Errorf("creating admin client: %w", err)
	}
	defer adminClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get cluster metadata
	metadata, err := adminClient.GetMetadata(nil, false, int(30*time.Second/time.Millisecond))
	if err != nil {
		return report, fmt.Errorf("getting metadata: %w", err)
	}

	report.Available = true
	report.BrokerCount = len(metadata.Brokers)

	// Detect controller type (ZooKeeper vs KRaft)
	report.ControllerType = detectControllerType(metadata)

	// Estimate controller count based on cluster type
	report.ControllerCount = estimateControllerCount(report.ControllerType, report.BrokerCount)

	// Get broker information
	for _, broker := range metadata.Brokers {
		brokerInfo := model.BrokerInfo{
			ID:   int(broker.ID),
			Host: broker.Host,
			Port: broker.Port,
		}
		report.Brokers = append(report.Brokers, brokerInfo)
	}

	// Get topic information
	topics := make([]string, 0, len(metadata.Topics))
	for topicName := range metadata.Topics {
		topics = append(topics, topicName)
	}

	if len(topics) > 0 {
		report.TopicCount = len(topics)
		report.InternalTopics = 0
		report.ExternalTopics = 0

		// Categorize topics as internal or external
		for _, topicName := range topics {
			if isInternalTopic(topicName) {
				report.InternalTopics++
			} else {
				report.ExternalTopics++
			}
		}

		// Get detailed topic information
		topicMetadata, err := adminClient.GetMetadata(&topics[0], true, int(30*time.Second/time.Millisecond))
		if err == nil {
			report.TotalPartitions = 0
			for topicName, topic := range topicMetadata.Topics {
				partitionCount := len(topic.Partitions)
				report.TotalPartitions += partitionCount

				// Always collect topic information (basic details)
				topicInfo := model.TopicInfo{
					Name:              topicName,
					IsInternal:        isInternalTopic(topicName),
					Partitions:        partitionCount,
					ReplicationFactor: getReplicationFactor(topic),
					AssociatedSchemas: make([]string, 0),
				}

				// Get topic configuration for retention settings
				if configs, err := getTopicConfigs(adminClient, ctx, topicName); err == nil {
					topicInfo.RetentionMs = configs.RetentionMs
					topicInfo.RetentionBytes = configs.RetentionBytes
				}

				// Calculate topic storage size by querying partition offsets
				// This is done for all topics to show cluster storage
				topicSize, err := calculateTopicSize(config, topicName, partitionCount)
				if err == nil {
					topicInfo.SizeBytes = topicSize
				}

				// Throughput will be populated from metrics if available
				// For now, these are placeholders that can be enriched from Prometheus/JMX
				topicInfo.ThroughputBytesInPerSec = 0
				topicInfo.ThroughputBytesOutPerSec = 0

				report.Topics = append(report.Topics, topicInfo)
			}
		}
	}

	// Get security configuration
	report.SecurityConfig = getKafkaSecurityConfig(adminClient, ctx, config, metadata)

	// Calculate cluster metrics (simulated - in real scenario, use JMX or metrics API)
	report.ClusterMetrics = calculateClusterMetrics(report)

	return report, nil
}

func detectControllerType(metadata *kafka.Metadata) string {
	// In KRaft mode, controller.id is typically visible in broker configs
	// For this implementation, we'll use heuristics to detect KRaft vs ZooKeeper

	if len(metadata.Brokers) == 0 {
		return "unknown"
	}

	// Collect broker IDs
	brokerIDs := make([]int32, 0, len(metadata.Brokers))
	for _, broker := range metadata.Brokers {
		brokerIDs = append(brokerIDs, broker.ID)
	}

	// Method 1: Check if broker IDs are sequential starting from 0 or 1
	// This is common in KRaft combined mode
	isSequentialFrom0 := isSequentialBrokerIDs(brokerIDs, 0)
	isSequentialFrom1 := isSequentialBrokerIDs(brokerIDs, 1)

	if isSequentialFrom0 || isSequentialFrom1 {
		// Sequential IDs starting from 0 or 1 strongly indicate KRaft
		return "kraft"
	}

	// Method 2: Check for broker ID patterns
	// Find min and max broker IDs
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

	// KRaft clusters typically have:
	// - IDs starting from 0 or 1
	// - Relatively small IDs (< 100)
	// - Compact ID range
	idRange := maxID - minID

	// If IDs start from 0 or 1 and have small range, likely KRaft
	if (minID == 0 || minID == 1) && maxID < 100 && idRange < 100 {
		return "kraft"
	}

	// If all IDs are very small (< 10) even if not sequential, likely KRaft
	if maxID < 10 {
		return "kraft"
	}

	// Method 3: Check for large broker IDs (> 100)
	// ZooKeeper clusters often have arbitrary IDs assigned
	hasLargeIDs := maxID > 100

	if !hasLargeIDs && minID <= 10 {
		// Small IDs, more likely KRaft
		return "kraft"
	}

	// Default to ZooKeeper if we have large or arbitrary IDs
	return "zookeeper"
}

func estimateControllerCount(controllerType string, brokerCount int) int {
	if controllerType == "kraft" {
		// KRaft mode can have:
		// 1. Combined mode: All brokers are also controllers
		// 2. Dedicated mode: Separate controller nodes (typically 3)
		// 3. Mixed mode: Some brokers are controllers

		// For typical production KRaft deployments:
		// - Small clusters (1-3 brokers): Combined mode (all are controllers)
		// - Medium/Large clusters: Dedicated controllers (usually 3)

		if brokerCount <= 3 {
			// Small cluster, likely combined mode
			return brokerCount
		} else {
			// Larger cluster, likely has dedicated controllers
			// Standard is 3 controllers for HA
			return 3
		}
	} else if controllerType == "zookeeper" {
		// ZooKeeper mode: 1 active controller elected from brokers
		return 1
	}

	// Unknown type
	return 0
}

func isSequentialBrokerIDs(ids []int32, startFrom int32) bool {
	if len(ids) == 0 {
		return false
	}

	// Create a set of IDs for quick lookup
	idSet := make(map[int32]bool)
	for _, id := range ids {
		idSet[id] = true
	}

	// Check if we have sequential IDs starting from startFrom
	for i := int32(0); i < int32(len(ids)); i++ {
		expectedID := startFrom + i
		if !idSet[expectedID] {
			return false
		}
	}

	return true
}

func getReplicationFactor(topic kafka.TopicMetadata) int {
	if len(topic.Partitions) == 0 {
		return 0
	}
	// Return replication factor of first partition
	return len(topic.Partitions[0].Replicas)
}

type TopicConfigs struct {
	RetentionMs    int64
	RetentionBytes int64
}

func getTopicConfigs(adminClient *kafka.AdminClient, ctx context.Context, topicName string) (TopicConfigs, error) {
	configs := TopicConfigs{
		RetentionMs:    -1,
		RetentionBytes: -1,
	}

	// Describe topic configs
	resources := []kafka.ConfigResource{
		{
			Type: kafka.ResourceTopic,
			Name: topicName,
		},
	}

	results, err := adminClient.DescribeConfigs(ctx, resources)
	if err != nil {
		return configs, err
	}

	if len(results) > 0 {
		for _, entry := range results[0].Config {
			switch entry.Name {
			case "retention.ms":
				fmt.Sscanf(entry.Value, "%d", &configs.RetentionMs)
			case "retention.bytes":
				fmt.Sscanf(entry.Value, "%d", &configs.RetentionBytes)
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
		TotalDiskUsageBytes:            0,
		UnderReplicatedPartitions: 0,
	}

	// Calculate total cluster storage from topic sizes
	var totalTopicStorage int64 = 0
	for _, topic := range report.Topics {
		totalTopicStorage += topic.SizeBytes
	}
	metrics.TotalDiskUsageBytes = totalTopicStorage

	// In a real implementation, these would be fetched from:
	// 1. JMX metrics (kafka.server:type=BrokerTopicMetrics,name=BytesInPerSec)
	// 2. Metrics Reporter API
	// 3. Confluent Control Center API
	// 4. Prometheus metrics if enabled

	// These would be real-time metrics in production
	// Placeholder values to show structure
	metrics.BytesInPerSec = 0
	metrics.BytesOutPerSec = 0
	metrics.MessagesInPerSec = 0

	return metrics
}

// calculateTopicSize estimates topic storage size by querying partition high watermarks
func calculateTopicSize(config model.KafkaConfig, topicName string, partitionCount int) (int64, error) {
	// Create consumer configuration
	consumerConfig := kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          "cp-discovery-size-calculator",
		"auto.offset.reset": "earliest",
	}

	// Add security configuration
	if config.SecurityProtocol != "" && config.SecurityProtocol != "PLAINTEXT" {
		consumerConfig["security.protocol"] = config.SecurityProtocol
	}
	if config.SaslMechanism != "" {
		consumerConfig["sasl.mechanism"] = config.SaslMechanism
		consumerConfig["sasl.username"] = config.SaslUsername
		consumerConfig["sasl.password"] = config.SaslPassword
	}

	// Add SSL/TLS configuration
	if config.SslCaLocation != "" {
		consumerConfig["ssl.ca.location"] = config.SslCaLocation
	}
	if config.SslCertLocation != "" {
		consumerConfig["ssl.certificate.location"] = config.SslCertLocation
	}
	if config.SslKeyLocation != "" {
		consumerConfig["ssl.key.location"] = config.SslKeyLocation
	}
	if config.SslKeyPassword != "" {
		consumerConfig["ssl.key.password"] = config.SslKeyPassword
	}
	if config.SslEndpointIdentification != "" {
		consumerConfig["ssl.endpoint.identification.algorithm"] = config.SslEndpointIdentification
	}

	// Create consumer
	consumer, err := kafka.NewConsumer(&consumerConfig)
	if err != nil {
		return 0, fmt.Errorf("creating consumer: %w", err)
	}
	defer consumer.Close()

	var totalSize int64 = 0

	// Query each partition's watermarks
	for partition := 0; partition < partitionCount; partition++ {
		// Get low and high watermarks
		low, high, err := consumer.QueryWatermarkOffsets(topicName, int32(partition), 5000)
		if err != nil {
			// If we can't query, skip this partition
			continue
		}

		// Calculate approximate size based on offset range
		// This is an estimate: actual size depends on message sizes
		// Average message size assumption: 1KB (can be refined based on actual data)
		offsetRange := high - low
		estimatedSize := offsetRange * 1024 // 1KB per message estimate

		totalSize += estimatedSize
	}

	return totalSize, nil
}

// Helper function to determine if a topic is internal
func isInternalTopic(topicName string) bool {
	// Internal topics typically start with underscore
	// Common internal topics:
	// - __consumer_offsets (consumer group offsets)
	// - __transaction_state (transactional state)
	// - _schemas (Schema Registry)
	// - _confluent* (Confluent Platform internal topics)
	// - connect-* (Kafka Connect internal topics)
	// - ksql* (ksqlDB internal topics)

	if strings.HasPrefix(topicName, "_") {
		return true
	}

	// Additional patterns for internal topics
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

func getKafkaSecurityConfig(adminClient *kafka.AdminClient, ctx context.Context, config model.KafkaConfig, metadata *kafka.Metadata) model.SecurityConfig {
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

	// Try to get broker configs to see all enabled SASL mechanisms
	if len(metadata.Brokers) > 0 {
		// Query first broker's configs
		brokerID := int32(metadata.Brokers[0].ID)

		// Try to describe broker configs
		resources := []kafka.ConfigResource{
			{
				Type: kafka.ResourceBroker,
				Name: fmt.Sprintf("%d", brokerID),
			},
		}

		results, err := adminClient.DescribeConfigs(ctx, resources)
		if err == nil && len(results) > 0 {
			for _, entry := range results[0].Config {
				switch entry.Name {
				case "sasl.enabled.mechanisms":
					if entry.Value != "" {
						mechanisms := strings.Split(entry.Value, ",")
						mechanismsMap := make(map[string]bool)
						for _, mech := range mechanisms {
							trimmed := strings.TrimSpace(mech)
							if trimmed != "" {
								mechanismsMap[trimmed] = true
							}
						}
						// Add client mechanism if not already in list
						if config.SaslMechanism != "" {
							mechanismsMap[config.SaslMechanism] = true
						}
						// Convert map to slice
						secConfig.SaslMechanisms = make([]string, 0, len(mechanismsMap))
						for mech := range mechanismsMap {
							secConfig.SaslMechanisms = append(secConfig.SaslMechanisms, mech)
						}
						if len(secConfig.SaslMechanisms) > 0 {
							secConfig.SaslEnabled = true
						}
					}

				case "listener.security.protocol.map":
					if entry.Value != "" {
						protocolsMap := make(map[string]bool)
						mappings := strings.Split(entry.Value, ",")
						for _, mapping := range mappings {
							parts := strings.Split(mapping, ":")
							if len(parts) == 2 {
								protocol := strings.TrimSpace(parts[1])
								protocolsMap[protocol] = true
								if strings.Contains(protocol, "SSL") {
									secConfig.SslEnabled = true
								}
								if strings.Contains(protocol, "SASL") {
									secConfig.SaslEnabled = true
								}
							}
						}
						// Add client protocol if not already in list
						if config.SecurityProtocol != "" && config.SecurityProtocol != "PLAINTEXT" {
							protocolsMap[config.SecurityProtocol] = true
						}
						// Convert map to slice
						secConfig.SecurityProtocols = make([]string, 0, len(protocolsMap))
						for proto := range protocolsMap {
							secConfig.SecurityProtocols = append(secConfig.SecurityProtocols, proto)
						}
					}

				case "security.inter.broker.protocol":
					if entry.Value != "" {
						protocolExists := false
						for _, proto := range secConfig.SecurityProtocols {
							if proto == entry.Value {
								protocolExists = true
								break
							}
						}
						if !protocolExists {
							secConfig.SecurityProtocols = append(secConfig.SecurityProtocols, entry.Value)
						}
						if strings.Contains(entry.Value, "SSL") {
							secConfig.SslEnabled = true
						}
						if strings.Contains(entry.Value, "SASL") {
							secConfig.SaslEnabled = true
						}
					}
				}
			}
		}
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

// Helper function to check ZooKeeper connectivity (optional)
func detectZooKeeperNodes(bootstrapServers string) int {
	// This would require ZooKeeper client library
	// For now, return 0 as placeholder
	// In production, parse ZooKeeper connection string and check connectivity

	// Typical implementation would:
	// 1. Get ZooKeeper connection string from broker config
	// 2. Connect to ZooKeeper
	// 3. Count nodes in /brokers/ids

	return 0
}

// Helper to parse broker rack information
func getBrokerRack(brokerID int32, metadata *kafka.Metadata) string {
	// Rack information would come from broker metadata
	// This requires additional API calls or JMX
	return ""
}

// Helper to get disk usage per broker (requires JMX or metrics API)
func getBrokerDiskUsage(brokerID int32, host string, port int) int64 {
	// In production, this would:
	// 1. Connect to JMX on broker
	// 2. Query kafka.log:type=Log,name=Size metric
	// 3. Or use Confluent Metrics API

	return 0
}
