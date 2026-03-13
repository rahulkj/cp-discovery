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

	// Fetch additional detailed information
	if detailed {
		additionalInfo := fetchAdditionalKafkaInfo(client, ctx, metadata, report.Topics)
		report.AdditionalInfo = &additionalInfo
	}

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

// fetchAdditionalKafkaInfo collects extended information about the Kafka cluster
func fetchAdditionalKafkaInfo(client *kafka.Client, ctx context.Context, metadata *kafka.MetadataResponse, topics []model.TopicInfo) model.KafkaAdditionalInfo {
	info := model.KafkaAdditionalInfo{
		ConsumerGroups:     make([]model.KafkaConsumerGroup, 0),
		DetailedPartitions: make([]model.DetailedPartitionInfo, 0),
		BrokerConfigs:      make([]model.BrokerConfigInfo, 0),
		ApiVersions:        make([]model.ApiVersionInfo, 0),
	}

	// Extract cluster ID from metadata
	if metadata != nil && metadata.ClusterID != "" {
		info.ClusterID = metadata.ClusterID
	}

	// Identify controller
	if metadata != nil && len(metadata.Brokers) > 0 {
		// In KRaft mode, try to find the active controller
		// In ZooKeeper mode, there's typically one controller
		info.ControllerID = metadata.Controller.ID
	}

	// Fetch consumer groups
	consumerGroups, activeCount := fetchConsumerGroups(client, ctx)
	info.ConsumerGroups = consumerGroups
	info.TotalConsumerGroups = len(consumerGroups)
	info.ActiveConsumerGroups = activeCount

	// Fetch detailed partition information
	info.DetailedPartitions = fetchDetailedPartitions(client, ctx, metadata)

	// Fetch broker configurations
	info.BrokerConfigs = fetchBrokerConfigs(client, ctx, metadata)

	// Fetch API versions
	info.ApiVersions = fetchApiVersions(client, ctx)

	return info
}

// fetchConsumerGroups retrieves all consumer groups and their details
func fetchConsumerGroups(client *kafka.Client, ctx context.Context) ([]model.KafkaConsumerGroup, int) {
	groups := make([]model.KafkaConsumerGroup, 0)
	activeCount := 0

	// List consumer groups
	listReq := &kafka.ListGroupsRequest{}
	listResp, err := client.ListGroups(ctx, listReq)
	if err != nil {
		return groups, activeCount
	}

	// For each group, get detailed information
	for _, groupInfo := range listResp.Groups {
		group := model.KafkaConsumerGroup{
			GroupID:      groupInfo.GroupID,
			ProtocolType: groupInfo.ProtocolType,
			Coordinator:  groupInfo.Coordinator,
			Members:      make([]model.ConsumerGroupMember, 0),
			Partitions:   make([]model.ConsumerGroupPartition, 0),
		}

		// Describe group to get member details and state
		descReq := &kafka.DescribeGroupsRequest{
			GroupIDs: []string{groupInfo.GroupID},
		}
		descResp, err := client.DescribeGroups(ctx, descReq)
		if err == nil && len(descResp.Groups) > 0 {
			groupDetail := descResp.Groups[0]
			group.State = groupDetail.GroupState
			group.MemberCount = len(groupDetail.Members)

			// Count active groups
			if groupDetail.GroupState == "Stable" || groupDetail.GroupState == "PreparingRebalance" || groupDetail.GroupState == "CompletingRebalance" {
				activeCount++
			}

			// Extract member information
			for _, member := range groupDetail.Members {
				groupMember := model.ConsumerGroupMember{
					MemberID:           member.MemberID,
					ClientID:           member.ClientID,
					ClientHost:         member.ClientHost,
					AssignedTopics:     make([]string, 0),
					AssignedPartitions: 0,
				}

				// Parse member assignment to get topics and partitions
				if len(member.MemberAssignments.Topics) > 0 {
					topicsMap := make(map[string]bool)
					for _, topicAssignment := range member.MemberAssignments.Topics {
						topicsMap[topicAssignment.Topic] = true
						groupMember.AssignedPartitions += len(topicAssignment.Partitions)
					}
					for topic := range topicsMap {
						groupMember.AssignedTopics = append(groupMember.AssignedTopics, topic)
					}
				}

				group.Members = append(group.Members, groupMember)
			}
		}

		// Fetch offsets for the group
		offsetFetchReq := &kafka.OffsetFetchRequest{
			GroupID: groupInfo.GroupID,
		}
		offsetFetchResp, err := client.OffsetFetch(ctx, offsetFetchReq)
		if err == nil {
			var totalLag int64 = 0

			for topic, partitions := range offsetFetchResp.Topics {
				for _, partition := range partitions {
					// Get high watermark for this partition
					highReq := &kafka.ListOffsetsRequest{
						Topics: map[string][]kafka.OffsetRequest{
							topic: {
								{
									Partition: partition.Partition,
									Timestamp: kafka.LastOffset,
								},
							},
						},
					}

					highResp, err := client.ListOffsets(ctx, highReq)
					if err != nil {
						continue
					}

					var logEndOffset int64
					if topicOffsets, ok := highResp.Topics[topic]; ok && len(topicOffsets) > 0 {
						logEndOffset = topicOffsets[0].LastOffset
					}

					lag := logEndOffset - partition.CommittedOffset
					if lag < 0 {
						lag = 0
					}
					totalLag += lag

					cgPartition := model.ConsumerGroupPartition{
						Topic:         topic,
						Partition:     partition.Partition,
						CurrentOffset: partition.CommittedOffset,
						LogEndOffset:  logEndOffset,
						Lag:           lag,
					}

					group.Partitions = append(group.Partitions, cgPartition)
				}
			}

			group.TotalLag = totalLag
		}

		groups = append(groups, group)
	}

	return groups, activeCount
}

// fetchDetailedPartitions retrieves detailed partition information
func fetchDetailedPartitions(client *kafka.Client, ctx context.Context, metadata *kafka.MetadataResponse) []model.DetailedPartitionInfo {
	partitions := make([]model.DetailedPartitionInfo, 0)

	if metadata == nil {
		return partitions
	}

	// For each topic, get partition details
	for _, topic := range metadata.Topics {
		for _, partition := range topic.Partitions {
			detailedPartition := model.DetailedPartitionInfo{
				Topic:           topic.Name,
				Partition:       partition.ID,
				Leader:          partition.Leader.ID,
				Replicas:        make([]int, 0),
				ISR:             make([]int, 0),
				OfflineReplicas: make([]int, 0),
			}

			// Collect replica IDs
			for _, replica := range partition.Replicas {
				detailedPartition.Replicas = append(detailedPartition.Replicas, replica.ID)
			}

			// Collect ISR IDs
			for _, isr := range partition.Isr {
				detailedPartition.ISR = append(detailedPartition.ISR, isr.ID)
			}

			// Identify offline replicas (replicas not in ISR)
			isrMap := make(map[int]bool)
			for _, isr := range partition.Isr {
				isrMap[isr.ID] = true
			}
			for _, replica := range partition.Replicas {
				if !isrMap[replica.ID] {
					detailedPartition.OfflineReplicas = append(detailedPartition.OfflineReplicas, replica.ID)
				}
			}

			// Get offset information
			// Get high watermark (last offset)
			highReq := &kafka.ListOffsetsRequest{
				Topics: map[string][]kafka.OffsetRequest{
					topic.Name: {
						{
							Partition: partition.ID,
							Timestamp: kafka.LastOffset,
						},
					},
				},
			}
			highResp, err := client.ListOffsets(ctx, highReq)
			if err == nil {
				if topicOffsets, ok := highResp.Topics[topic.Name]; ok && len(topicOffsets) > 0 {
					detailedPartition.LastOffset = topicOffsets[0].LastOffset
				}
			}

			// Get low watermark (first offset)
			lowReq := &kafka.ListOffsetsRequest{
				Topics: map[string][]kafka.OffsetRequest{
					topic.Name: {
						{
							Partition: partition.ID,
							Timestamp: kafka.FirstOffset,
						},
					},
				},
			}
			lowResp, err := client.ListOffsets(ctx, lowReq)
			if err == nil {
				if topicOffsets, ok := lowResp.Topics[topic.Name]; ok && len(topicOffsets) > 0 {
					detailedPartition.FirstOffset = topicOffsets[0].LastOffset
				}
			}

			// Calculate message count
			detailedPartition.MessageCount = detailedPartition.LastOffset - detailedPartition.FirstOffset

			partitions = append(partitions, detailedPartition)
		}
	}

	return partitions
}

// fetchBrokerConfigs retrieves broker configurations
func fetchBrokerConfigs(client *kafka.Client, ctx context.Context, metadata *kafka.MetadataResponse) []model.BrokerConfigInfo {
	brokerConfigs := make([]model.BrokerConfigInfo, 0)

	if metadata == nil {
		return brokerConfigs
	}

	// For each broker, fetch its configuration
	for _, broker := range metadata.Brokers {
		brokerConfig := model.BrokerConfigInfo{
			BrokerID: broker.ID,
			Configs:  make(map[string]model.BrokerConfigEntry),
		}

		// Describe broker configs
		req := &kafka.DescribeConfigsRequest{
			Resources: []kafka.DescribeConfigRequestResource{
				{
					ResourceType: kafka.ResourceTypeBroker,
					ResourceName: fmt.Sprintf("%d", broker.ID),
				},
			},
		}

		resp, err := client.DescribeConfigs(ctx, req)
		if err != nil {
			continue
		}

		if len(resp.Resources) > 0 {
			for _, entry := range resp.Resources[0].ConfigEntries {
				configEntry := model.BrokerConfigEntry{
					Name:      entry.ConfigName,
					Value:     entry.ConfigValue,
					Source:    string(entry.ConfigSource),
					Sensitive: entry.IsSensitive,
					ReadOnly:  entry.ReadOnly,
				}
				brokerConfig.Configs[entry.ConfigName] = configEntry
			}
		}

		brokerConfigs = append(brokerConfigs, brokerConfig)
	}

	return brokerConfigs
}

// fetchApiVersions retrieves supported API versions from the broker
func fetchApiVersions(client *kafka.Client, ctx context.Context) []model.ApiVersionInfo {
	apiVersions := make([]model.ApiVersionInfo, 0)

	// API versions request
	req := &kafka.ApiVersionsRequest{}
	resp, err := client.ApiVersions(ctx, req)
	if err != nil {
		return apiVersions
	}

	// Map of API keys to names
	apiNames := map[int16]string{
		0:  "Produce",
		1:  "Fetch",
		2:  "ListOffsets",
		3:  "Metadata",
		8:  "OffsetCommit",
		9:  "OffsetFetch",
		10: "FindCoordinator",
		11: "JoinGroup",
		12: "Heartbeat",
		13: "LeaveGroup",
		14: "SyncGroup",
		15: "DescribeGroups",
		16: "ListGroups",
		17: "SaslHandshake",
		18: "ApiVersions",
		19: "CreateTopics",
		20: "DeleteTopics",
		21: "DeleteRecords",
		22: "InitProducerId",
		23: "OffsetForLeaderEpoch",
		24: "AddPartitionsToTxn",
		25: "AddOffsetsToTxn",
		26: "EndTxn",
		27: "WriteTxnMarkers",
		28: "TxnOffsetCommit",
		29: "DescribeAcls",
		30: "CreateAcls",
		31: "DeleteAcls",
		32: "DescribeConfigs",
		33: "AlterConfigs",
		36: "SaslAuthenticate",
		37: "CreatePartitions",
		38: "CreateDelegationToken",
		39: "RenewDelegationToken",
		40: "ExpireDelegationToken",
		41: "DescribeDelegationToken",
		42: "DeleteGroups",
		43: "ElectLeaders",
		44: "IncrementalAlterConfigs",
		45: "AlterPartitionReassignments",
		46: "ListPartitionReassignments",
		50: "DescribeUserScramCredentials",
		51: "AlterUserScramCredentials",
	}

	for _, apiKey := range resp.ApiKeys {
		apiVersion := model.ApiVersionInfo{
			ApiKey:     int16(apiKey.ApiKey),
			MinVersion: int16(apiKey.MinVersion),
			MaxVersion: int16(apiKey.MaxVersion),
		}

		// Add API name if known
		if name, ok := apiNames[int16(apiKey.ApiKey)]; ok {
			apiVersion.ApiName = name
		}

		apiVersions = append(apiVersions, apiVersion)
	}

	return apiVersions
}
