package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/rahulkj/cp-discovery/internal/model"
	httpauth "github.com/rahulkj/cp-discovery/internal/http"
)

type KsqlDBServerInfo struct {
	Version        string `json:"version"`
	KafkaClusterID string `json:"kafkaClusterId"`
	KsqlServiceID  string `json:"ksqlServiceId"`
}

type KsqlDBResponse struct {
	StatementText      string                   `json:"statementText"`
	Queries            []map[string]interface{} `json:"queries"`
	SourceDescriptions []map[string]interface{} `json:"sourceDescriptions"`
}

type KsqlDBClusterStatusResponse struct {
	ClusterStatus map[string]KsqlDBServerStatus `json:"clusterStatus"`
}

type KsqlDBServerStatus struct {
	HostAlive     bool   `json:"hostAlive"`
	LastStatusUpdate int64 `json:"lastStatusUpdate"`
	HostStoreLags map[string]interface{} `json:"hostStoreLags"`
}

func DiscoverKsqlDB(config model.KsqlDBConfig, detailed bool) (model.KsqlDBReport, error) {
	report := model.KsqlDBReport{
		Available: false,
	}

	if config.URL == "" {
		return report, fmt.Errorf("ksqlDB URL not configured")
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Check if ksqlDB is available and get version
	infoURL := fmt.Sprintf("%s/info", config.URL)
	req, err := http.NewRequest("GET", infoURL, nil)
	if err != nil {
		return report, fmt.Errorf("creating request: %w", err)
	}

	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return report, fmt.Errorf("connecting to ksqlDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("ksqlDB returned status: %d", resp.StatusCode)
	}

	report.Available = true

	// Get version info
	var serverInfo KsqlDBServerInfo
	body, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(body, &serverInfo) == nil {
		report.Version = serverInfo.Version
	}

	// Get cluster status to count nodes
	nodeCount := getKsqlDBClusterNodeCount(client, config)
	report.NodeCount = nodeCount

	// Get queries
	queriesCount, _ := getKsqlDBQueries(client, config)
	report.Queries = queriesCount

	// Get streams and tables
	streamsCount, tablesCount, _ := getKsqlDBStreamsAndTables(client, config)
	report.Streams = streamsCount
	report.Tables = tablesCount

	// Fetch additional detailed information
	if detailed {
		additionalInfo := fetchKsqlDBAdditionalInfo(client, config)
		report.AdditionalInfo = &additionalInfo
	}

	return report, nil
}

func getKsqlDBClusterNodeCount(client *http.Client, config model.KsqlDBConfig) int {
	url := fmt.Sprintf("%s/clusterStatus", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 1 // Default to single node
	}

	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 1 // Default to single node
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 1 // Default to single node
	}

	body, _ := io.ReadAll(resp.Body)
	var clusterStatus KsqlDBClusterStatusResponse
	if json.Unmarshal(body, &clusterStatus) == nil {
		// Count servers that are alive
		aliveCount := 0
		for _, status := range clusterStatus.ClusterStatus {
			if status.HostAlive {
				aliveCount++
			}
		}
		if aliveCount > 0 {
			return aliveCount
		}
		// If no alive status info, return total count
		if len(clusterStatus.ClusterStatus) > 0 {
			return len(clusterStatus.ClusterStatus)
		}
	}

	return 1 // Default to single node
}

func getKsqlDBQueries(client *http.Client, config model.KsqlDBConfig) (int, error) {
	url := fmt.Sprintf("%s/ksql", config.URL)

	payload := map[string]interface{}{
		"ksql": "SHOW QUERIES;",
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/vnd.ksql.v1+json")
	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, nil
	}

	body, _ := io.ReadAll(resp.Body)

	// Parse response - ksqlDB returns array of responses
	var responses []map[string]interface{}
	if err := json.Unmarshal(body, &responses); err != nil {
		return 0, err
	}

	// Count queries from response
	queryCount := 0
	for _, response := range responses {
		if queries, ok := response["queries"].([]interface{}); ok {
			queryCount = len(queries)
		}
	}

	return queryCount, nil
}

func getKsqlDBStreamsAndTables(client *http.Client, config model.KsqlDBConfig) (int, int, error) {
	url := fmt.Sprintf("%s/ksql", config.URL)

	// Get streams
	streamsPayload := map[string]interface{}{
		"ksql": "SHOW STREAMS;",
	}

	streamsCount := 0
	jsonData, _ := json.Marshal(streamsPayload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err == nil {
		req.Header.Set("Content-Type", "application/vnd.ksql.v1+json")
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var responses []map[string]interface{}
				if json.Unmarshal(body, &responses) == nil {
					for _, response := range responses {
						if streams, ok := response["streams"].([]interface{}); ok {
							streamsCount = len(streams)
						}
					}
				}
			}
		}
	}

	// Get tables
	tablesPayload := map[string]interface{}{
		"ksql": "SHOW TABLES;",
	}

	tablesCount := 0
	jsonData, _ = json.Marshal(tablesPayload)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err == nil {
		req.Header.Set("Content-Type", "application/vnd.ksql.v1+json")
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var responses []map[string]interface{}
				if json.Unmarshal(body, &responses) == nil {
					for _, response := range responses {
						if tables, ok := response["tables"].([]interface{}); ok {
							// Filter out internal tables
							for _, table := range tables {
								if tableMap, ok := table.(map[string]interface{}); ok {
									if name, ok := tableMap["name"].(string); ok {
										if !strings.HasPrefix(name, "KSQL_") {
											tablesCount++
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return streamsCount, tablesCount, nil
}

// fetchKsqlDBAdditionalInfo collects extended information about ksqlDB
func fetchKsqlDBAdditionalInfo(client *http.Client, config model.KsqlDBConfig) model.KsqlDBAdditionalInfo {
	info := model.KsqlDBAdditionalInfo{
		Queries:    make([]model.KsqlQueryDetail, 0),
		Streams:    make([]model.KsqlStreamDetail, 0),
		Tables:     make([]model.KsqlTableDetail, 0),
		Topics:     make([]string, 0),
		Connectors: make([]string, 0),
	}

	// Get server info
	info.ServerInfo = getKsqlServerInfo(client, config)

	// Get cluster status
	info.ClusterStatus = getKsqlClusterStatus(client, config)

	// Get detailed queries
	info.Queries = getDetailedKsqlQueries(client, config)

	// Get detailed streams and tables
	streams, tables := getDetailedKsqlStreamsAndTables(client, config)
	info.Streams = streams
	info.Tables = tables

	// Get topics
	info.Topics = getKsqlTopics(client, config)

	// Get connectors
	info.Connectors = getKsqlConnectors(client, config)

	return info
}

// getKsqlServerInfo retrieves ksqlDB server information
func getKsqlServerInfo(client *http.Client, config model.KsqlDBConfig) model.KsqlServerInfo {
	serverInfo := model.KsqlServerInfo{}

	url := fmt.Sprintf("%s/info", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return serverInfo
	}

	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return serverInfo
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return serverInfo
	}

	body, _ := io.ReadAll(resp.Body)
	var info struct {
		Version         string `json:"version"`
		KafkaClusterID  string `json:"kafkaClusterId"`
		KsqlServiceID   string `json:"ksqlServiceId"`
		ServerStatus    string `json:"serverStatus"`
	}

	if json.Unmarshal(body, &info) == nil {
		serverInfo.Version = info.Version
		serverInfo.KafkaClusterID = info.KafkaClusterID
		serverInfo.KsqlServiceID = info.KsqlServiceID
		serverInfo.ServerStatus = info.ServerStatus
	}

	return serverInfo
}

// getKsqlClusterStatus retrieves ksqlDB cluster status
func getKsqlClusterStatus(client *http.Client, config model.KsqlDBConfig) model.KsqlClusterStatus {
	clusterStatus := model.KsqlClusterStatus{
		Hosts: make([]model.KsqlHostInfo, 0),
	}

	url := fmt.Sprintf("%s/clusterStatus", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return clusterStatus
	}

	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return clusterStatus
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return clusterStatus
	}

	body, _ := io.ReadAll(resp.Body)
	var status struct {
		ClusterStatus map[string]interface{} `json:"clusterStatus"`
	}

	if json.Unmarshal(body, &status) == nil {
		// Parse cluster status structure
		for _, hostData := range status.ClusterStatus {
			if hostMap, ok := hostData.(map[string]interface{}); ok {
				hostInfo := model.KsqlHostInfo{}
				if hostInfoData, ok := hostMap["hostInfo"].(map[string]interface{}); ok {
					hostInfo.HostInfo.Host, _ = hostInfoData["host"].(string)
					if port, ok := hostInfoData["port"].(float64); ok {
						hostInfo.HostInfo.Port = int(port)
					}
				}
				if isActive, ok := hostMap["hostAlive"].(bool); ok {
					hostInfo.IsActiveHost = isActive
				}
				clusterStatus.Hosts = append(clusterStatus.Hosts, hostInfo)
			}
		}
	}

	return clusterStatus
}

// getDetailedKsqlQueries retrieves detailed query information
func getDetailedKsqlQueries(client *http.Client, config model.KsqlDBConfig) []model.KsqlQueryDetail {
	queries := make([]model.KsqlQueryDetail, 0)

	payload := map[string]string{"ksql": "SHOW QUERIES;"}
	jsonData, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/ksql", config.URL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return queries
	}

	req.Header.Set("Content-Type", "application/vnd.ksql.v1+json; charset=utf-8")
	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return queries
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var responses []map[string]interface{}
		if json.Unmarshal(body, &responses) == nil {
			for _, response := range responses {
				if queriesData, ok := response["queries"].([]interface{}); ok {
					for _, q := range queriesData {
						if queryMap, ok := q.(map[string]interface{}); ok {
							query := model.KsqlQueryDetail{}
							query.ID, _ = queryMap["id"].(string)
							query.QueryString, _ = queryMap["queryString"].(string)
							query.State, _ = queryMap["state"].(string)
							query.QueryType, _ = queryMap["queryType"].(string)
							queries = append(queries, query)
						}
					}
				}
			}
		}
	}

	return queries
}

// getDetailedKsqlStreamsAndTables retrieves detailed stream and table information
func getDetailedKsqlStreamsAndTables(client *http.Client, config model.KsqlDBConfig) ([]model.KsqlStreamDetail, []model.KsqlTableDetail) {
	streams := make([]model.KsqlStreamDetail, 0)
	tables := make([]model.KsqlTableDetail, 0)

	// Get streams
	streamsPayload := map[string]string{"ksql": "SHOW STREAMS EXTENDED;"}
	jsonData, _ := json.Marshal(streamsPayload)

	url := fmt.Sprintf("%s/ksql", config.URL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err == nil {
		req.Header.Set("Content-Type", "application/vnd.ksql.v1+json; charset=utf-8")
		httpauth.ApplyKsqlDBAuth(req, config)

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var responses []map[string]interface{}
				if json.Unmarshal(body, &responses) == nil {
					for _, response := range responses {
						if streamsData, ok := response["streams"].([]interface{}); ok {
							for _, s := range streamsData {
								if streamMap, ok := s.(map[string]interface{}); ok {
									stream := model.KsqlStreamDetail{}
									stream.Name, _ = streamMap["name"].(string)
									stream.Topic, _ = streamMap["topic"].(string)
									stream.KeyFormat, _ = streamMap["keyFormat"].(string)
									stream.ValueFormat, _ = streamMap["valueFormat"].(string)
									if !strings.HasPrefix(stream.Name, "KSQL_") {
										streams = append(streams, stream)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Get tables
	tablesPayload := map[string]string{"ksql": "SHOW TABLES EXTENDED;"}
	jsonData, _ = json.Marshal(tablesPayload)

	req, err = http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err == nil {
		req.Header.Set("Content-Type", "application/vnd.ksql.v1+json; charset=utf-8")
		httpauth.ApplyKsqlDBAuth(req, config)

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var responses []map[string]interface{}
				if json.Unmarshal(body, &responses) == nil {
					for _, response := range responses {
						if tablesData, ok := response["tables"].([]interface{}); ok {
							for _, t := range tablesData {
								if tableMap, ok := t.(map[string]interface{}); ok {
									table := model.KsqlTableDetail{}
									table.Name, _ = tableMap["name"].(string)
									table.Topic, _ = tableMap["topic"].(string)
									table.KeyFormat, _ = tableMap["keyFormat"].(string)
									table.ValueFormat, _ = tableMap["valueFormat"].(string)
									if !strings.HasPrefix(table.Name, "KSQL_") {
										tables = append(tables, table)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return streams, tables
}

// getKsqlTopics retrieves topics used by ksqlDB
func getKsqlTopics(client *http.Client, config model.KsqlDBConfig) []string {
	topics := make([]string, 0)

	payload := map[string]string{"ksql": "SHOW TOPICS;"}
	jsonData, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/ksql", config.URL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return topics
	}

	req.Header.Set("Content-Type", "application/vnd.ksql.v1+json; charset=utf-8")
	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return topics
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var responses []map[string]interface{}
		if json.Unmarshal(body, &responses) == nil {
			for _, response := range responses {
				if topicsData, ok := response["topics"].([]interface{}); ok {
					for _, t := range topicsData {
						if topicMap, ok := t.(map[string]interface{}); ok {
							if name, ok := topicMap["name"].(string); ok {
								topics = append(topics, name)
							}
						}
					}
				}
			}
		}
	}

	return topics
}

// getKsqlConnectors retrieves connectors managed by ksqlDB
func getKsqlConnectors(client *http.Client, config model.KsqlDBConfig) []string {
	connectors := make([]string, 0)

	payload := map[string]string{"ksql": "SHOW CONNECTORS;"}
	jsonData, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/ksql", config.URL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return connectors
	}

	req.Header.Set("Content-Type", "application/vnd.ksql.v1+json; charset=utf-8")
	httpauth.ApplyKsqlDBAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return connectors
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var responses []map[string]interface{}
		if json.Unmarshal(body, &responses) == nil {
			for _, response := range responses {
				if connectorsData, ok := response["connectors"].([]interface{}); ok {
					for _, c := range connectorsData {
						if connectorMap, ok := c.(map[string]interface{}); ok {
							if name, ok := connectorMap["name"].(string); ok {
								connectors = append(connectors, name)
							}
						}
					}
				}
			}
		}
	}

	return connectors
}
