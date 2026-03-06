package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/rahulkj/cp-discovery/internal/model"
	httpauth "github.com/rahulkj/cp-discovery/internal/http"
)

type ControlCenterClusterInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type C3KafkaCluster struct {
	ClusterID   string `json:"clusterId"`
	ClusterName string `json:"clusterName"`
	Hosts       []string `json:"hosts"`
	BrokerCount int    `json:"brokerCount"`
	TopicCount  int    `json:"topicCount"`
	PartitionCount int `json:"partitionCount"`
}

type C3KafkaClusterDetail struct {
	Cluster struct {
		ClusterID string `json:"clusterId"`
		Name      string `json:"name"`
		Health    struct {
			Status string `json:"status"`
		} `json:"health"`
	} `json:"cluster"`
	Brokers []struct {
		ID   int    `json:"id"`
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"brokers"`
	Topics []struct {
		TopicName      string `json:"topicName"`
		PartitionCount int    `json:"partitionCount"`
	} `json:"topics"`
}

type C3ConnectCluster struct {
	ClusterID   string `json:"clusterId"`
	ClusterName string `json:"clusterName"`
	KafkaClusterID string `json:"kafkaClusterId"`
}

type C3ConnectClusterDetail struct {
	Cluster struct {
		ClusterID string `json:"clusterId"`
	} `json:"cluster"`
	Connectors []struct {
		Name  string `json:"name"`
		State string `json:"state"`
		Type  string `json:"type"`
	} `json:"connectors"`
	Workers []struct {
		WorkerID string `json:"workerId"`
	} `json:"workers"`
}

type C3SchemaRegistryCluster struct {
	ClusterID   string `json:"clusterId"`
	ClusterName string `json:"clusterName"`
	KafkaClusterID string `json:"kafkaClusterId"`
}

type C3SchemaRegistryDetail struct {
	Cluster struct {
		ClusterID string `json:"clusterId"`
		Version   string `json:"version"`
	} `json:"cluster"`
	SubjectCount int `json:"subjectCount"`
}

type C3KsqlCluster struct {
	ClusterID      string `json:"clusterId"`
	ClusterName    string `json:"clusterName"`
	KafkaClusterID string `json:"kafkaClusterId"`
}

type C3KsqlClusterDetail struct {
	Cluster struct {
		ClusterID string `json:"clusterId"`
	} `json:"cluster"`
	Queries []struct {
		QueryID string `json:"queryId"`
	} `json:"queries"`
	Streams []struct {
		Name string `json:"name"`
	} `json:"streams"`
	Tables []struct {
		Name string `json:"name"`
	} `json:"tables"`
}

type C3ConsumerGroupLag struct {
	GroupID string `json:"groupId"`
	Lag     int64  `json:"lag"`
}

func DiscoverControlCenter(config model.ControlCenterConfig, detailed bool) (model.ControlCenterReport, error) {
	report := model.ControlCenterReport{
		Available: false,
	}

	if config.URL == "" {
		return report, fmt.Errorf("control center URL not configured")
	}

	report.URL = config.URL

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if Control Center is available
	healthURL := fmt.Sprintf("%s/health", config.URL)
	req, err := http.NewRequest("GET", healthURL, nil)
	if err != nil {
		return report, fmt.Errorf("creating request: %w", err)
	}

	httpauth.ApplyControlCenterAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return report, fmt.Errorf("connecting to control center: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("control center returned status: %d", resp.StatusCode)
	}

	report.Available = true
	report.NodeCount = 1 // Control Center doesn't expose cluster info; assume at least 1 node

	// Try to get version and cluster information
	clustersURL := fmt.Sprintf("%s/2.0/clusters/kafka", config.URL)
	req, err = http.NewRequest("GET", clustersURL, nil)
	if err == nil {
		httpauth.ApplyControlCenterAuth(req, config)

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var clusters []ControlCenterClusterInfo
				if json.Unmarshal(body, &clusters) == nil {
					report.MonitoredClusters = len(clusters)
					// Get version from first cluster if available
					if len(clusters) > 0 && clusters[0].Version != "" {
						report.Version = clusters[0].Version
					}
				}
			}
		}
	}

	// Try alternative version endpoint
	if report.Version == "" {
		versionURL := fmt.Sprintf("%s/api/version", config.URL)
		req, err := http.NewRequest("GET", versionURL, nil)
		if err == nil {
			httpauth.ApplyControlCenterAuth(req, config)

			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					body, _ := io.ReadAll(resp.Body)
					var versionInfo map[string]interface{}
					if json.Unmarshal(body, &versionInfo) == nil {
						if version, ok := versionInfo["version"].(string); ok {
							report.Version = version
						}
					}
				}
			}
		}
	}

	// Fetch detailed information if requested
	if detailed {
		getC3DetailedInfo(client, config, &report)
	}

	return report, nil
}

func getC3DetailedInfo(client *http.Client, config model.ControlCenterConfig, report *model.ControlCenterReport) {
	// Fetch Kafka clusters with details
	kafkaClusters := getC3KafkaClusters(client, config)
	report.Clusters = kafkaClusters

	// Fetch Connect clusters
	connectClusters := getC3ConnectClusters(client, config)
	report.ConnectClusters = connectClusters

	// Fetch Schema Registry clusters
	schemaRegistries := getC3SchemaRegistries(client, config)
	report.SchemaRegistries = schemaRegistries

	// Fetch ksqlDB clusters
	ksqlClusters := getC3KsqlClusters(client, config)
	report.KsqlClusters = ksqlClusters

	// Fetch consumer lag metrics
	totalLag := getC3ConsumerLag(client, config)
	report.TotalConsumerLag = totalLag
}

func getC3KafkaClusters(client *http.Client, config model.ControlCenterConfig) []model.C3ClusterInfo {
	clustersURL := fmt.Sprintf("%s/2.0/clusters/kafka", config.URL)
	req, err := http.NewRequest("GET", clustersURL, nil)
	if err != nil {
		return []model.C3ClusterInfo{}
	}

	httpauth.ApplyControlCenterAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.C3ClusterInfo{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []model.C3ClusterInfo{}
	}

	body, _ := io.ReadAll(resp.Body)
	var clusters []C3KafkaCluster
	if json.Unmarshal(body, &clusters) != nil {
		return []model.C3ClusterInfo{}
	}

	clusterInfos := make([]model.C3ClusterInfo, 0, len(clusters))
	for _, cluster := range clusters {
		clusterInfo := model.C3ClusterInfo{
			ClusterID:   cluster.ClusterID,
			ClusterName: cluster.ClusterName,
			BrokerCount: cluster.BrokerCount,
			TopicCount:  cluster.TopicCount,
			PartitionCount: cluster.PartitionCount,
		}

		// Try to get health status for the cluster
		healthURL := fmt.Sprintf("%s/2.0/clusters/kafka/%s/health", config.URL, cluster.ClusterID)
		healthReq, err := http.NewRequest("GET", healthURL, nil)
		if err == nil {
			httpauth.ApplyControlCenterAuth(healthReq, config)

			healthResp, err := client.Do(healthReq)
			if err == nil {
				defer healthResp.Body.Close()
				if healthResp.StatusCode == http.StatusOK {
					healthBody, _ := io.ReadAll(healthResp.Body)
					var healthData map[string]interface{}
					if json.Unmarshal(healthBody, &healthData) == nil {
						if status, ok := healthData["status"].(string); ok {
							clusterInfo.HealthStatus = status
						}
					}
				}
			}
		}

		clusterInfos = append(clusterInfos, clusterInfo)
	}

	return clusterInfos
}

func getC3ConnectClusters(client *http.Client, config model.ControlCenterConfig) []model.C3ConnectClusterInfo {
	connectURL := fmt.Sprintf("%s/2.0/clusters/connect", config.URL)
	req, err := http.NewRequest("GET", connectURL, nil)
	if err != nil {
		return []model.C3ConnectClusterInfo{}
	}

	httpauth.ApplyControlCenterAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.C3ConnectClusterInfo{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []model.C3ConnectClusterInfo{}
	}

	body, _ := io.ReadAll(resp.Body)
	var clusters []C3ConnectCluster
	if json.Unmarshal(body, &clusters) != nil {
		return []model.C3ConnectClusterInfo{}
	}

	connectInfos := make([]model.C3ConnectClusterInfo, 0, len(clusters))
	for _, cluster := range clusters {
		connectInfo := model.C3ConnectClusterInfo{
			ClusterName:    cluster.ClusterName,
			ClusterID:      cluster.ClusterID,
			KafkaClusterID: cluster.KafkaClusterID,
		}

		// Try to get detailed connector information
		detailURL := fmt.Sprintf("%s/2.0/clusters/connect/%s/connectors", config.URL, cluster.ClusterID)
		detailReq, err := http.NewRequest("GET", detailURL, nil)
		if err == nil {
			httpauth.ApplyControlCenterAuth(detailReq, config)

			detailResp, err := client.Do(detailReq)
			if err == nil {
				defer detailResp.Body.Close()
				if detailResp.StatusCode == http.StatusOK {
					detailBody, _ := io.ReadAll(detailResp.Body)
					var connectors []struct {
						Name  string `json:"name"`
						Type  string `json:"type"`
						State string `json:"state"`
						Tasks int    `json:"tasks"`
					}
					if json.Unmarshal(detailBody, &connectors) == nil {
						connectInfo.ConnectorCount = len(connectors)
						failedCount := 0
						sourceCount := 0
						sinkCount := 0
						runningCount := 0
						connectorDetails := make([]model.C3ConnectorInfo, 0, len(connectors))

						for _, conn := range connectors {
							// Count by state
							if conn.State == "FAILED" {
								failedCount++
							} else if conn.State == "RUNNING" {
								runningCount++
							}

							// Count by type
							if conn.Type == "source" {
								sourceCount++
							} else if conn.Type == "sink" {
								sinkCount++
							}

							// Add to detailed list
							connectorDetails = append(connectorDetails, model.C3ConnectorInfo{
								Name:  conn.Name,
								Type:  conn.Type,
								State: conn.State,
								Tasks: conn.Tasks,
							})
						}

						connectInfo.FailedConnectors = failedCount
						connectInfo.SourceConnectors = sourceCount
						connectInfo.SinkConnectors = sinkCount
						connectInfo.RunningConnectors = runningCount
						connectInfo.Connectors = connectorDetails
					}
				}
			}
		}

		// Try to get worker count
		workersURL := fmt.Sprintf("%s/2.0/clusters/connect/%s/workers", config.URL, cluster.ClusterID)
		workersReq, err := http.NewRequest("GET", workersURL, nil)
		if err == nil {
			httpauth.ApplyControlCenterAuth(workersReq, config)

			workersResp, err := client.Do(workersReq)
			if err == nil {
				defer workersResp.Body.Close()
				if workersResp.StatusCode == http.StatusOK {
					workersBody, _ := io.ReadAll(workersResp.Body)
					var workers []map[string]interface{}
					if json.Unmarshal(workersBody, &workers) == nil {
						connectInfo.WorkerCount = len(workers)
					}
				}
			}
		}

		connectInfos = append(connectInfos, connectInfo)
	}

	return connectInfos
}

func getC3SchemaRegistries(client *http.Client, config model.ControlCenterConfig) []model.C3SchemaRegistryInfo {
	srURL := fmt.Sprintf("%s/2.0/clusters/schema-registry", config.URL)
	req, err := http.NewRequest("GET", srURL, nil)
	if err != nil {
		return []model.C3SchemaRegistryInfo{}
	}

	httpauth.ApplyControlCenterAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.C3SchemaRegistryInfo{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []model.C3SchemaRegistryInfo{}
	}

	body, _ := io.ReadAll(resp.Body)
	var clusters []C3SchemaRegistryCluster
	if json.Unmarshal(body, &clusters) != nil {
		return []model.C3SchemaRegistryInfo{}
	}

	srInfos := make([]model.C3SchemaRegistryInfo, 0, len(clusters))
	for _, cluster := range clusters {
		srInfo := model.C3SchemaRegistryInfo{
			ClusterName:    cluster.ClusterName,
			ClusterID:      cluster.ClusterID,
			KafkaClusterID: cluster.KafkaClusterID,
		}

		// Try to get schema count and version
		detailURL := fmt.Sprintf("%s/2.0/clusters/schema-registry/%s", config.URL, cluster.ClusterID)
		detailReq, err := http.NewRequest("GET", detailURL, nil)
		if err == nil {
			httpauth.ApplyControlCenterAuth(detailReq, config)

			detailResp, err := client.Do(detailReq)
			if err == nil {
				defer detailResp.Body.Close()
				if detailResp.StatusCode == http.StatusOK {
					detailBody, _ := io.ReadAll(detailResp.Body)
					var detail struct {
						Version      string   `json:"version"`
						SubjectCount int      `json:"subjectCount"`
						Mode         string   `json:"mode"`
						Subjects     []string `json:"subjects"`
					}
					if json.Unmarshal(detailBody, &detail) == nil {
						srInfo.Version = detail.Version
						srInfo.SchemaCount = detail.SubjectCount
						srInfo.Mode = detail.Mode
						srInfo.Subjects = detail.Subjects
					}
				}
			}
		}

		srInfos = append(srInfos, srInfo)
	}

	return srInfos
}

func getC3KsqlClusters(client *http.Client, config model.ControlCenterConfig) []model.C3KsqlClusterInfo {
	ksqlURL := fmt.Sprintf("%s/2.0/clusters/ksql", config.URL)
	req, err := http.NewRequest("GET", ksqlURL, nil)
	if err != nil {
		return []model.C3KsqlClusterInfo{}
	}

	httpauth.ApplyControlCenterAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.C3KsqlClusterInfo{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []model.C3KsqlClusterInfo{}
	}

	body, _ := io.ReadAll(resp.Body)
	var clusters []C3KsqlCluster
	if json.Unmarshal(body, &clusters) != nil {
		return []model.C3KsqlClusterInfo{}
	}

	ksqlInfos := make([]model.C3KsqlClusterInfo, 0, len(clusters))
	for _, cluster := range clusters {
		ksqlInfo := model.C3KsqlClusterInfo{
			ClusterName:    cluster.ClusterName,
			ClusterID:      cluster.ClusterID,
			KafkaClusterID: cluster.KafkaClusterID,
		}

		// Try to get query, stream, and table counts
		detailURL := fmt.Sprintf("%s/2.0/clusters/ksql/%s", config.URL, cluster.ClusterID)
		detailReq, err := http.NewRequest("GET", detailURL, nil)
		if err == nil {
			httpauth.ApplyControlCenterAuth(detailReq, config)

			detailResp, err := client.Do(detailReq)
			if err == nil {
				defer detailResp.Body.Close()
				if detailResp.StatusCode == http.StatusOK {
					detailBody, _ := io.ReadAll(detailResp.Body)
					var detail struct {
						QueryCount  int `json:"queryCount"`
						StreamCount int `json:"streamCount"`
						TableCount  int `json:"tableCount"`
					}
					if json.Unmarshal(detailBody, &detail) == nil {
						ksqlInfo.QueryCount = detail.QueryCount
						ksqlInfo.StreamCount = detail.StreamCount
						ksqlInfo.TableCount = detail.TableCount
					}
				}
			}
		}

		ksqlInfos = append(ksqlInfos, ksqlInfo)
	}

	return ksqlInfos
}

func getC3ConsumerLag(client *http.Client, config model.ControlCenterConfig) int64 {
	lagURL := fmt.Sprintf("%s/2.0/monitoring/consumer-groups/lag", config.URL)
	req, err := http.NewRequest("GET", lagURL, nil)
	if err != nil {
		return 0
	}

	httpauth.ApplyControlCenterAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	body, _ := io.ReadAll(resp.Body)
	var lagData struct {
		TotalLag int64 `json:"totalLag"`
	}
	if json.Unmarshal(body, &lagData) == nil {
		return lagData.TotalLag
	}

	// Alternative: sum up lag from all consumer groups
	var groups []C3ConsumerGroupLag
	if json.Unmarshal(body, &groups) == nil {
		totalLag := int64(0)
		for _, group := range groups {
			totalLag += group.Lag
		}
		return totalLag
	}

	return 0
}
