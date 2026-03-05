package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/rahulkj/cp-discovery/internal/model"
	httpauth "github.com/rahulkj/cp-discovery/internal/http"
)

type RestProxyInfo struct {
	Version string `json:"version"`
}

type V3ClustersResponse struct {
	Data []V3ClusterInfo `json:"data"`
}

type V3ClusterInfo struct {
	ClusterID  string                 `json:"cluster_id"`
	Controller V3ControllerInfo       `json:"controller"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type V3ControllerInfo struct {
	Related string `json:"related"` // URL to the controller broker resource
}

type V3ControllerResponse struct {
	Kind        string                 `json:"kind"`
	Metadata    map[string]interface{} `json:"metadata"`
	ClusterID   string                 `json:"cluster_id"`
	BrokerID    int                    `json:"broker_id"`
	Host        string                 `json:"host"`
	Port        int                    `json:"port"`
	Rack        string                 `json:"rack"`
}

type V3BrokersResponse struct {
	Data []V3BrokerInfo `json:"data"`
}

type V3BrokerInfo struct {
	BrokerID int                    `json:"broker_id"`
	Host     string                 `json:"host"`
	Port     int                    `json:"port"`
	Rack     string                 `json:"rack"`
	Related  map[string]interface{} `json:"related"`
}

type V3BrokerDetailResponse struct {
	ClusterID string                 `json:"cluster_id"`
	BrokerID  int                    `json:"broker_id"`
	Host      string                 `json:"host"`
	Port      int                    `json:"port"`
	Rack      string                 `json:"rack"`
	Configs   map[string]interface{} `json:"configs"`
}

type V3TopicsResponse struct {
	Data []V3TopicInfo `json:"data"`
}

type V3TopicInfo struct {
	TopicName         string      `json:"topic_name"`
	PartitionsCount   int         `json:"partitions_count"`
	Partitions        interface{} `json:"partitions"`
	ReplicationFactor int         `json:"replication_factor"`
	IsInternal        bool        `json:"is_internal"`
}

type V3TopicDetail struct {
	ClusterID         string                 `json:"cluster_id"`
	TopicName         string                 `json:"topic_name"`
	IsInternal        bool                   `json:"is_internal"`
	PartitionsCount   int                    `json:"partitions_count"`
	ReplicationFactor int                    `json:"replication_factor"`
	Partitions        map[string]interface{} `json:"partitions"`
}

type V3TopicConfigsResponse struct {
	Data []V3TopicConfig `json:"data"`
}

type V3TopicConfig struct {
	Name   string      `json:"name"`
	Value  string      `json:"value"`
	Source string      `json:"source"`
}

type V3BrokerConfigsResponse struct {
	Data []V3BrokerConfig `json:"data"`
}

type V3BrokerConfig struct {
	Kind        string                 `json:"kind"`
	Metadata    map[string]interface{} `json:"metadata"`
	ClusterID   string                 `json:"cluster_id"`
	BrokerID    int                    `json:"broker_id"`
	Name        string                 `json:"name"`
	Value       string                 `json:"value"`
	IsDefault   bool                   `json:"is_default"`
	IsReadOnly  bool                   `json:"is_read_only"`
	IsSensitive bool                   `json:"is_sensitive"`
	Source      string                 `json:"source"`
	Synonyms    []interface{}          `json:"synonyms"`
}

type V3ConsumerGroupsResponse struct {
	Data []V3ConsumerGroup `json:"data"`
}

type V3ConsumerGroup struct {
	GroupID        string `json:"consumer_group_id"`
	IsSimple       bool   `json:"is_simple"`
	PartitionAssignor string `json:"partition_assignor"`
	State          string `json:"state"`
	Coordinator    V3CoordinatorInfo `json:"coordinator"`
}

type V3CoordinatorInfo struct {
	Related string `json:"related"`
}

type V3ConsumerGroupDetail struct {
	GroupID           string                  `json:"consumer_group_id"`
	IsSimple          bool                    `json:"is_simple"`
	PartitionAssignor string                  `json:"partition_assignor"`
	State             string                  `json:"state"`
	Consumers         V3ConsumersInfo         `json:"consumers"`
}

type V3ConsumersInfo struct {
	Related string `json:"related"`
}

type V3ConsumersResponse struct {
	Data []V3Consumer `json:"data"`
}

type V3Consumer struct {
	ConsumerID    string `json:"consumer_id"`
	InstanceID    string `json:"instance_id"`
	ClientID      string `json:"client_id"`
}

type V3AclsResponse struct {
	Data []V3Acl `json:"data"`
}

type V3Acl struct {
	ClusterID      string `json:"cluster_id"`
	ResourceType   string `json:"resource_type"`
	ResourceName   string `json:"resource_name"`
	PatternType    string `json:"pattern_type"`
	Principal      string `json:"principal"`
	Host           string `json:"host"`
	Operation      string `json:"operation"`
	Permission     string `json:"permission"`
}

type V3ClusterConfigsResponse struct {
	Data []V3ClusterConfig `json:"data"`
}

type V3ClusterConfig struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	IsDefault  bool   `json:"is_default"`
	IsReadOnly bool   `json:"is_read_only"`
	IsSensitive bool  `json:"is_sensitive"`
	Source     string `json:"source"`
}

func DiscoverRestProxy(config model.RestProxyConfig, detailed bool) (model.RestProxyReport, error) {
	report := model.RestProxyReport{
		Available: false,
	}

	if config.URL == "" {
		return report, fmt.Errorf("REST proxy URL not configured")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if REST Proxy is available
	url := fmt.Sprintf("%s/", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return report, fmt.Errorf("creating request: %w", err)
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return report, fmt.Errorf("connecting to REST proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("REST proxy returned status: %d", resp.StatusCode)
	}

	report.Available = true

	// Try to get version from response
	body, _ := io.ReadAll(resp.Body)

	// REST Proxy root endpoint may return version info
	var info map[string]interface{}
	if json.Unmarshal(body, &info) == nil {
		if version, ok := info["version"].(string); ok {
			report.Version = version
		}
	}

	// Try to get cluster info and broker count from v3 API
	clustersURL := fmt.Sprintf("%s/v3/clusters", config.URL)
	req, err = http.NewRequest("GET", clustersURL, nil)
	if err == nil {
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				// If version not set, indicate v3 API is available
				if report.Version == "" {
					report.Version = "v3+"
				}

				// Get cluster ID to fetch detailed information
				body, _ := io.ReadAll(resp.Body)
				var clustersResp V3ClustersResponse
				if json.Unmarshal(body, &clustersResp) == nil && len(clustersResp.Data) > 0 {
					clusterID := clustersResp.Data[0].ClusterID
					report.ClusterID = clusterID

					// Fetch detailed cluster information
					getClusterDetails(client, config, clusterID, &report, detailed)
				}
			}
		}
	}

	return report, nil
}

func getClusterDetails(client *http.Client, config model.RestProxyConfig, clusterID string, report *model.RestProxyReport, detailed bool) {
	// Get detailed cluster information including controller
	clusterURL := fmt.Sprintf("%s/v3/clusters/%s", config.URL, clusterID)
	req, err := http.NewRequest("GET", clusterURL, nil)
	if err != nil {
		return
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	activeControllerID := -1
	var clusterMetadata V3ClusterInfo
	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if json.Unmarshal(body, &clusterMetadata) == nil {
			// Try to fetch controller information directly from the API
			if clusterMetadata.Controller.Related != "" {
				controllerInfo := getControllerInfo(client, config, clusterMetadata.Controller.Related)
				if controllerInfo != nil {
					activeControllerID = controllerInfo.BrokerID
					report.ControllerID = activeControllerID
				} else {
					// Fallback: Extract controller ID from related URL
					if controllerID := extractControllerID(clusterMetadata.Controller.Related); controllerID >= 0 {
						report.ControllerID = controllerID
						activeControllerID = controllerID
					}
				}
			}
		}
	}

	// Fetch broker information with controller ID for comparison
	brokers, controllerMode, controllerCount := getBrokerDetails(client, config, clusterID, activeControllerID, clusterMetadata)
	report.BrokerCount = len(brokers)
	report.ControllerMode = controllerMode
	report.ControllerCount = controllerCount
	if detailed {
		report.Brokers = brokers
	}

	// Fetch security configuration from broker configs
	securityConfig := getSecurityConfig(client, config, clusterID, brokers)
	report.SecurityConfig = securityConfig

	// Fetch topic count and partition count
	topicCount, internalCount, externalCount, partitionCount, avgRF := getTopicStats(client, config, clusterID)
	report.TopicCount = topicCount
	report.InternalTopics = internalCount
	report.ExternalTopics = externalCount
	report.PartitionCount = partitionCount
	report.AvgReplicationFactor = avgRF

	// Fetch consumer groups information
	consumerGroups, activeCount := getConsumerGroups(client, config, clusterID, detailed)
	report.ConsumerGroups = consumerGroups
	report.ConsumerGroupCount = len(consumerGroups)
	report.ActiveConsumerGroups = activeCount

	// Fetch ACLs information
	acls := getAcls(client, config, clusterID, detailed)
	report.Acls = acls
	report.AclCount = len(acls)

	// Fetch cluster-level configurations
	if detailed {
		clusterConfig := getClusterConfig(client, config, clusterID)
		report.ClusterConfig = clusterConfig
	}
}

func getControllerInfo(client *http.Client, config model.RestProxyConfig, controllerURL string) *V3ControllerResponse {
	// The controllerURL is a relative path like /v3/clusters/{cluster_id}/brokers/{broker_id}
	// Construct the full URL
	fullURL := config.URL + controllerURL

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	var controllerResp V3ControllerResponse
	if json.Unmarshal(body, &controllerResp) == nil {
		return &controllerResp
	}

	return nil
}

func getBrokerDetails(client *http.Client, config model.RestProxyConfig, clusterID string, activeControllerID int, clusterMetadata V3ClusterInfo) ([]model.RestProxyBrokerInfo, string, int) {
	brokersURL := fmt.Sprintf("%s/v3/clusters/%s/brokers", config.URL, clusterID)
	req, err := http.NewRequest("GET", brokersURL, nil)
	if err != nil {
		return []model.RestProxyBrokerInfo{}, "unknown", 0
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.RestProxyBrokerInfo{}, "unknown", 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []model.RestProxyBrokerInfo{}, "unknown", 0
	}

	body, _ := io.ReadAll(resp.Body)
	var brokersResp V3BrokersResponse
	if json.Unmarshal(body, &brokersResp) != nil {
		return []model.RestProxyBrokerInfo{}, "unknown", 0
	}

	// Convert to model.RestProxyBrokerInfo and determine controller mode
	brokers := make([]model.RestProxyBrokerInfo, 0, len(brokersResp.Data))
	brokersWithControllerRole := 0
	brokerIDs := make([]int, 0, len(brokersResp.Data))

	for _, broker := range brokersResp.Data {
		brokerIDs = append(brokerIDs, broker.BrokerID)

		// Check if this broker is the active controller by comparing IDs
		isActiveController := (broker.BrokerID == activeControllerID)

		// Additionally check if broker has controller role (KRaft mode)
		// by checking broker configs for process.roles
		hasControllerRole := checkBrokerControllerRole(client, config, clusterID, broker.BrokerID)

		if hasControllerRole {
			brokersWithControllerRole++
		}

		brokers = append(brokers, model.RestProxyBrokerInfo{
			BrokerID:           broker.BrokerID,
			Host:               broker.Host,
			Port:               broker.Port,
			IsActiveController: isActiveController,
			HasControllerRole:  hasControllerRole,
		})
	}

	// Determine controller mode using multiple indicators
	controllerMode := determineControllerMode(brokersWithControllerRole, len(brokers), brokerIDs, activeControllerID, clusterMetadata)

	return brokers, controllerMode, brokersWithControllerRole
}

func determineControllerMode(brokersWithControllerRole int, totalBrokers int, brokerIDs []int, activeControllerID int, clusterMetadata V3ClusterInfo) string {
	// Method 1: Check if we found explicit controller roles
	if brokersWithControllerRole > 0 {
		if brokersWithControllerRole == totalBrokers {
			return "kraft-combined" // All brokers have controller role
		} else {
			return "kraft-separated" // Some nodes are dedicated controllers
		}
	}

	// Method 2: Use heuristics to detect KRaft when process.roles is not available
	// KRaft clusters typically have:
	// 1. Sequential broker IDs starting from 0 or 1
	// 2. Active controller ID that matches a broker ID
	// 3. Cluster metadata might have KRaft-specific fields

	if totalBrokers > 0 && activeControllerID >= 0 {
		// Check if broker IDs are sequential (common in KRaft)
		isSequential := isSequentialIDs(brokerIDs)

		// Check if active controller ID is in the broker list
		controllerIsBroker := false
		for _, id := range brokerIDs {
			if id == activeControllerID {
				controllerIsBroker = true
				break
			}
		}

		// If we have sequential IDs starting from 0 or 1, and controller is a broker,
		// it's very likely KRaft combined mode
		if isSequential && controllerIsBroker {
			minID := brokerIDs[0]
			for _, id := range brokerIDs {
				if id < minID {
					minID = id
				}
			}
			// KRaft typically starts from 0 or 1
			if minID <= 1 {
				return "kraft-combined"
			}
		}

		// Check cluster metadata for KRaft indicators
		if clusterMetadata.Metadata != nil {
			// Look for metadata.server.type or similar KRaft indicators
			if serverType, ok := clusterMetadata.Metadata["server.type"].(string); ok {
				if strings.Contains(strings.ToLower(serverType), "kraft") {
					return "kraft-combined"
				}
			}
		}
	}

	// If we have a controller but can't determine KRaft, default to ZooKeeper
	// Only return "zookeeper" if we have clear indicators it's NOT KRaft
	if activeControllerID >= 0 {
		// If controller is among brokers with sequential IDs, likely KRaft
		if len(brokerIDs) > 0 && isSequentialIDs(brokerIDs) {
			return "kraft-combined"
		}
		return "zookeeper"
	}

	return "unknown"
}

func isSequentialIDs(ids []int) bool {
	if len(ids) == 0 {
		return false
	}

	// Sort IDs to check sequence
	sortedIDs := make([]int, len(ids))
	copy(sortedIDs, ids)

	// Simple bubble sort
	for i := 0; i < len(sortedIDs); i++ {
		for j := i + 1; j < len(sortedIDs); j++ {
			if sortedIDs[i] > sortedIDs[j] {
				sortedIDs[i], sortedIDs[j] = sortedIDs[j], sortedIDs[i]
			}
		}
	}

	// Check if sequential
	for i := 1; i < len(sortedIDs); i++ {
		if sortedIDs[i] != sortedIDs[i-1]+1 {
			return false
		}
	}

	return true
}

func checkBrokerControllerRole(client *http.Client, config model.RestProxyConfig, clusterID string, brokerID int) bool {
	// First, try to get all broker configs to check for process.roles
	configsURL := fmt.Sprintf("%s/v3/clusters/%s/brokers/%d/configs", config.URL, clusterID, brokerID)
	req, err := http.NewRequest("GET", configsURL, nil)
	if err != nil {
		return false
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var configsResp V3BrokerConfigsResponse
		if json.Unmarshal(body, &configsResp) == nil {
			// Look for process.roles or related configs
			for _, cfg := range configsResp.Data {
				// Check multiple possible config names
				if cfg.Name == "process.roles" || cfg.Name == "server.process.roles" {
					value := strings.TrimSpace(cfg.Value)
					// process.roles can be "controller", "broker", or "broker,controller"
					if value != "" && (strings.Contains(value, "controller") || value == "broker,controller" || value == "controller,broker") {
						return true
					}
				}
				// Also check for node.id which is KRaft-specific
				if cfg.Name == "node.id" && cfg.Value != "" {
					// node.id exists in KRaft mode, not in ZooKeeper mode
					// This is an indicator but not definitive
					// We'll use it in combination with other checks
				}
			}
		}
	}

	// Fallback 1: try direct config endpoint for process.roles
	processRolesURL := fmt.Sprintf("%s/v3/clusters/%s/brokers/%d/configs/process.roles", config.URL, clusterID, brokerID)
	if hasControllerInConfig(client, config, processRolesURL) {
		return true
	}

	// Fallback 2: try server.process.roles
	serverProcessRolesURL := fmt.Sprintf("%s/v3/clusters/%s/brokers/%d/configs/server.process.roles", config.URL, clusterID, brokerID)
	if hasControllerInConfig(client, config, serverProcessRolesURL) {
		return true
	}

	return false
}

func hasControllerInConfig(client *http.Client, config model.RestProxyConfig, configURL string) bool {
	req, err := http.NewRequest("GET", configURL, nil)
	if err != nil {
		return false
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var configData V3BrokerConfig
		if json.Unmarshal(body, &configData) == nil {
			value := strings.TrimSpace(configData.Value)
			if value != "" && strings.Contains(value, "controller") {
				return true
			}
		}
	}

	return false
}

func extractControllerID(relatedURL string) int {
	// Extract broker ID from URL like: /v3/clusters/{cluster_id}/brokers/{broker_id}
	// Simple parsing - split by '/' and get last segment
	if relatedURL == "" {
		return -1
	}

	parts := make([]rune, 0)
	for i := len(relatedURL) - 1; i >= 0; i-- {
		if relatedURL[i] == '/' {
			break
		}
		parts = append([]rune{rune(relatedURL[i])}, parts...)
	}

	var controllerID int
	if _, err := fmt.Sscanf(string(parts), "%d", &controllerID); err == nil {
		return controllerID
	}

	return -1
}

func getTopicStats(client *http.Client, config model.RestProxyConfig, clusterID string) (int, int, int, int, float64) {
	topicsURL := fmt.Sprintf("%s/v3/clusters/%s/topics", config.URL, clusterID)
	req, err := http.NewRequest("GET", topicsURL, nil)
	if err != nil {
		return 0, 0, 0, 0, 0.0
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, 0.0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, 0, 0.0
	}

	body, _ := io.ReadAll(resp.Body)
	var topicsResp V3TopicsResponse
	if json.Unmarshal(body, &topicsResp) != nil {
		return 0, 0, 0, 0, 0.0
	}

	topicCount := len(topicsResp.Data)
	internalCount := 0
	externalCount := 0
	partitionCount := 0
	totalRF := 0
	rfCount := 0

	// Count total partitions across all topics and categorize
	for _, topic := range topicsResp.Data {
		// Categorize as internal or external
		if isInternalTopicRP(topic.TopicName) {
			internalCount++
		} else {
			externalCount++
		}

		// Try to get partition count and replication factor from partitions_count field
		if topic.PartitionsCount > 0 {
			partitionCount += topic.PartitionsCount
		} else if topic.Partitions != nil {
			// If partitions_count not available, count partitions array
			if partitions, ok := topic.Partitions.([]interface{}); ok {
				partitionCount += len(partitions)
			}
		}

		// If replication factor is available in the topic list response
		if topic.ReplicationFactor > 0 {
			totalRF += topic.ReplicationFactor
			rfCount++
		} else {
			// Fallback: fetch topic details to get replication factor and partition count
			topicDetail := getTopicDetail(client, config, clusterID, topic.TopicName)
			if topicDetail != nil {
				if topicDetail.ReplicationFactor > 0 {
					totalRF += topicDetail.ReplicationFactor
					rfCount++
				}
				if topic.PartitionsCount == 0 && topicDetail.PartitionsCount > 0 {
					partitionCount += topicDetail.PartitionsCount
				}
			}
		}
	}

	avgRF := 0.0
	if rfCount > 0 {
		avgRF = float64(totalRF) / float64(rfCount)
	}

	return topicCount, internalCount, externalCount, partitionCount, avgRF
}

func getTopicDetail(client *http.Client, config model.RestProxyConfig, clusterID, topicName string) *V3TopicDetail {
	topicDetailURL := fmt.Sprintf("%s/v3/clusters/%s/topics/%s", config.URL, clusterID, topicName)
	req, err := http.NewRequest("GET", topicDetailURL, nil)
	if err != nil {
		return nil
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	var topicDetail V3TopicDetail
	if json.Unmarshal(body, &topicDetail) == nil {
		return &topicDetail
	}

	return nil
}

// Helper function to determine if a topic is internal (for REST Proxy)
func isInternalTopicRP(topicName string) bool {
	// Import strings package for string operations
	// Internal topics typically start with underscore
	if len(topicName) > 0 && topicName[0] == '_' {
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
		if len(topicName) >= len(prefix) && topicName[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

func getSecurityConfig(client *http.Client, config model.RestProxyConfig, clusterID string, brokers []model.RestProxyBrokerInfo) model.SecurityConfig {
	secConfig := model.SecurityConfig{
		SaslMechanisms:    make([]string, 0),
		SecurityProtocols: make([]string, 0),
		SslEnabled:        false,
		SaslEnabled:       false,
	}

	if len(brokers) == 0 {
		return secConfig
	}

	// Query the first broker's configs to determine security settings
	// All brokers should have the same security configuration
	brokerID := brokers[0].BrokerID

	configsURL := fmt.Sprintf("%s/v3/clusters/%s/brokers/%d/configs", config.URL, clusterID, brokerID)
	req, err := http.NewRequest("GET", configsURL, nil)
	if err != nil {
		return secConfig
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return secConfig
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return secConfig
	}

	body, _ := io.ReadAll(resp.Body)
	var configsResp V3BrokerConfigsResponse
	if json.Unmarshal(body, &configsResp) != nil {
		return secConfig
	}

	// Parse security-related configs
	saslMechanismsMap := make(map[string]bool)
	securityProtocolsMap := make(map[string]bool)
	listenerSecurityMap := make(map[string]string)

	for _, cfg := range configsResp.Data {
		switch cfg.Name {
		case "sasl.enabled.mechanisms":
			// Parse comma-separated list of SASL mechanisms
			if cfg.Value != "" {
				mechanisms := strings.Split(cfg.Value, ",")
				for _, mech := range mechanisms {
					trimmed := strings.TrimSpace(mech)
					if trimmed != "" {
						saslMechanismsMap[trimmed] = true
						secConfig.SaslEnabled = true
					}
				}
			}

		case "listener.security.protocol.map":
			// Parse listener to protocol mapping
			// Format: PLAINTEXT:PLAINTEXT,SSL:SSL,SASL_PLAINTEXT:SASL_PLAINTEXT,SASL_SSL:SASL_SSL
			if cfg.Value != "" {
				mappings := strings.Split(cfg.Value, ",")
				for _, mapping := range mappings {
					parts := strings.Split(mapping, ":")
					if len(parts) == 2 {
						listener := strings.TrimSpace(parts[0])
						protocol := strings.TrimSpace(parts[1])
						listenerSecurityMap[listener] = protocol
						securityProtocolsMap[protocol] = true

						// Check if SSL or SASL is enabled
						if strings.Contains(protocol, "SSL") {
							secConfig.SslEnabled = true
						}
						if strings.Contains(protocol, "SASL") {
							secConfig.SaslEnabled = true
						}
					}
				}
			}

		case "security.inter.broker.protocol":
			// Inter-broker security protocol
			if cfg.Value != "" {
				securityProtocolsMap[cfg.Value] = true
				if strings.Contains(cfg.Value, "SSL") {
					secConfig.SslEnabled = true
				}
				if strings.Contains(cfg.Value, "SASL") {
					secConfig.SaslEnabled = true
				}
			}

		case "listeners", "advertised.listeners":
			// Parse listener configurations to detect protocols
			// Format: PLAINTEXT://host:port,SSL://host:port,SASL_SSL://host:port
			if cfg.Value != "" {
				listeners := strings.Split(cfg.Value, ",")
				for _, listener := range listeners {
					// Extract protocol from listener (before ://)
					if idx := strings.Index(listener, "://"); idx > 0 {
						protocol := listener[:idx]
						if secProto, ok := listenerSecurityMap[protocol]; ok {
							securityProtocolsMap[secProto] = true
						}
					}
				}
			}
		}
	}

	// Convert maps to slices
	for mech := range saslMechanismsMap {
		secConfig.SaslMechanisms = append(secConfig.SaslMechanisms, mech)
	}
	for proto := range securityProtocolsMap {
		secConfig.SecurityProtocols = append(secConfig.SecurityProtocols, proto)
	}

	// Determine authentication method
	secConfig.AuthenticationMethod = determineAuthMethod(secConfig)

	return secConfig
}

func determineAuthMethod(config model.SecurityConfig) string {
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

func getConsumerGroups(client *http.Client, config model.RestProxyConfig, clusterID string, detailed bool) ([]model.ConsumerGroupInfo, int) {
	cgURL := fmt.Sprintf("%s/v3/clusters/%s/consumer-groups", config.URL, clusterID)
	req, err := http.NewRequest("GET", cgURL, nil)
	if err != nil {
		return []model.ConsumerGroupInfo{}, 0
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.ConsumerGroupInfo{}, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []model.ConsumerGroupInfo{}, 0
	}

	body, _ := io.ReadAll(resp.Body)
	var cgResp V3ConsumerGroupsResponse
	if json.Unmarshal(body, &cgResp) != nil {
		return []model.ConsumerGroupInfo{}, 0
	}

	groups := make([]model.ConsumerGroupInfo, 0, len(cgResp.Data))
	activeCount := 0

	for _, cg := range cgResp.Data {
		cgInfo := model.ConsumerGroupInfo{
			GroupID:           cg.GroupID,
			State:             cg.State,
			PartitionAssignor: cg.PartitionAssignor,
		}

		// Count active consumer groups (state = "Stable" or "PreparingRebalance")
		if cg.State == "Stable" || cg.State == "PreparingRebalance" {
			activeCount++
		}

		// If detailed, fetch member count for each group
		if detailed {
			memberCount := getConsumerGroupMemberCount(client, config, clusterID, cg.GroupID)
			cgInfo.MemberCount = memberCount
		}

		groups = append(groups, cgInfo)
	}

	return groups, activeCount
}

func getConsumerGroupMemberCount(client *http.Client, config model.RestProxyConfig, clusterID, groupID string) int {
	membersURL := fmt.Sprintf("%s/v3/clusters/%s/consumer-groups/%s/consumers", config.URL, clusterID, groupID)
	req, err := http.NewRequest("GET", membersURL, nil)
	if err != nil {
		return 0
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	body, _ := io.ReadAll(resp.Body)
	var consumersResp V3ConsumersResponse
	if json.Unmarshal(body, &consumersResp) != nil {
		return 0
	}

	return len(consumersResp.Data)
}

func getAcls(client *http.Client, config model.RestProxyConfig, clusterID string, detailed bool) []model.AclInfo {
	if !detailed {
		return []model.AclInfo{}
	}

	aclsURL := fmt.Sprintf("%s/v3/clusters/%s/acls", config.URL, clusterID)
	req, err := http.NewRequest("GET", aclsURL, nil)
	if err != nil {
		return []model.AclInfo{}
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.AclInfo{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []model.AclInfo{}
	}

	body, _ := io.ReadAll(resp.Body)
	var aclsResp V3AclsResponse
	if json.Unmarshal(body, &aclsResp) != nil {
		return []model.AclInfo{}
	}

	acls := make([]model.AclInfo, 0, len(aclsResp.Data))
	for _, acl := range aclsResp.Data {
		acls = append(acls, model.AclInfo{
			ResourceType: acl.ResourceType,
			ResourceName: acl.ResourceName,
			PatternType:  acl.PatternType,
			Principal:    acl.Principal,
			Operation:    acl.Operation,
			Permission:   acl.Permission,
		})
	}

	return acls
}

func getClusterConfig(client *http.Client, config model.RestProxyConfig, clusterID string) map[string]string {
	configURL := fmt.Sprintf("%s/v3/clusters/%s/broker-configs", config.URL, clusterID)
	req, err := http.NewRequest("GET", configURL, nil)
	if err != nil {
		return make(map[string]string)
	}

	httpauth.ApplyRestProxyAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return make(map[string]string)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return make(map[string]string)
	}

	body, _ := io.ReadAll(resp.Body)
	var configResp V3ClusterConfigsResponse
	if json.Unmarshal(body, &configResp) != nil {
		return make(map[string]string)
	}

	clusterConfig := make(map[string]string)
	// Only include important cluster-wide configs, not defaults
	importantConfigs := map[string]bool{
		"compression.type":                  true,
		"log.retention.hours":               true,
		"log.retention.bytes":               true,
		"log.segment.bytes":                 true,
		"message.max.bytes":                 true,
		"replica.fetch.max.bytes":           true,
		"num.partitions":                    true,
		"default.replication.factor":        true,
		"min.insync.replicas":               true,
		"unclean.leader.election.enable":    true,
		"auto.create.topics.enable":         true,
		"delete.topic.enable":               true,
		"offsets.retention.minutes":         true,
		"transaction.state.log.replication.factor": true,
		"transaction.state.log.min.isr":     true,
	}

	for _, cfg := range configResp.Data {
		if importantConfigs[cfg.Name] && !cfg.IsDefault {
			clusterConfig[cfg.Name] = cfg.Value
		}
	}

	return clusterConfig
}
