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
