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

type ConnectInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

type ConnectorStatus struct {
	Name      string                 `json:"name"`
	Connector map[string]interface{} `json:"connector"`
	Tasks     []map[string]interface{} `json:"tasks"`
}

type ConnectorConfig struct {
	ConnectorClass string `json:"connector.class"`
	TasksMax       string `json:"tasks.max"`
}

type ConnectClusterInfo struct {
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	KafkaClusterID string `json:"kafka_cluster_id"`
}

type ConnectWorkerInfo struct {
	WorkerID string `json:"worker_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func DiscoverKafkaConnect(config model.KafkaConnectConfig, detailed bool) (model.KafkaConnectReport, error) {
	report := model.KafkaConnectReport{
		Available:        false,
		Connectors:       make([]model.ConnectorInfo, 0),
		SinkConnectors:   0,
		SourceConnectors: 0,
	}

	if config.URL == "" {
		return report, fmt.Errorf("kafka connect URL not configured")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if Kafka Connect is available and get version
	versionURL := fmt.Sprintf("%s/", config.URL)
	req, err := http.NewRequest("GET", versionURL, nil)
	if err != nil {
		return report, fmt.Errorf("creating request: %w", err)
	}

	httpauth.ApplyKafkaConnectAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return report, fmt.Errorf("connecting to kafka connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("kafka connect returned status: %d", resp.StatusCode)
	}

	report.Available = true

	// Get version info
	var connectInfo ConnectInfo
	body, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(body, &connectInfo) == nil {
		report.Version = connectInfo.Version
	}

	// Get worker count (for distributed mode)
	workerCount := getConnectWorkerCount(client, config)
	report.WorkerCount = workerCount

	// Get connectors list
	connectorsURL := fmt.Sprintf("%s/connectors", config.URL)
	req, err = http.NewRequest("GET", connectorsURL, nil)
	if err != nil {
		return report, nil
	}

	httpauth.ApplyKafkaConnectAuth(req, config)

	resp, err = client.Do(req)
	if err != nil {
		return report, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, nil
	}

	var connectorNames []string
	body, _ = io.ReadAll(resp.Body)
	if json.Unmarshal(body, &connectorNames) != nil {
		return report, nil
	}

	report.TotalConnectors = len(connectorNames)

	// Get detailed information for each connector
	for _, name := range connectorNames {
		connectorInfo, err := getConnectorInfo(client, config, name)
		if err != nil {
			continue
		}

		// Determine connector type
		connectorType := determineConnectorType(connectorInfo)
		if connectorType == "source" {
			report.SourceConnectors++
		} else if connectorType == "sink" {
			report.SinkConnectors++
		}

		if detailed {
			// Get connector status
			status, _ := getConnectorStatus(client, config, name)

			info := model.ConnectorInfo{
				Name:  name,
				Type:  connectorType,
				State: status,
				Tasks: getTaskCount(connectorInfo),
			}
			report.Connectors = append(report.Connectors, info)
		}
	}

	return report, nil
}

func getConnectorInfo(client *http.Client, config model.KafkaConnectConfig, name string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/connectors/%s/config", config.URL, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpauth.ApplyKafkaConnectAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	var connectorConfig map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &connectorConfig); err != nil {
		return nil, err
	}

	return connectorConfig, nil
}

func getConnectorStatus(client *http.Client, config model.KafkaConnectConfig, name string) (string, error) {
	url := fmt.Sprintf("%s/connectors/%s/status", config.URL, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "unknown", err
	}

	httpauth.ApplyKafkaConnectAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return "unknown", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "unknown", fmt.Errorf("status: %d", resp.StatusCode)
	}

	var status ConnectorStatus
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &status); err != nil {
		return "unknown", err
	}

	if state, ok := status.Connector["state"].(string); ok {
		return state, nil
	}

	return "unknown", nil
}

func determineConnectorType(config map[string]interface{}) string {
	if connectorClass, ok := config["connector.class"].(string); ok {
		lowerClass := strings.ToLower(connectorClass)

		// Common source connector patterns
		sourcePatterns := []string{"source", "debezium", "jdbc", "mongodb", "spooldir", "filestream"}
		for _, pattern := range sourcePatterns {
			if strings.Contains(lowerClass, pattern) && !strings.Contains(lowerClass, "sink") {
				return "source"
			}
		}

		// Common sink connector patterns
		sinkPatterns := []string{"sink", "s3", "elasticsearch", "jdbc", "hdfs"}
		for _, pattern := range sinkPatterns {
			if strings.Contains(lowerClass, pattern) {
				return "sink"
			}
		}
	}

	return "unknown"
}

func getTaskCount(config map[string]interface{}) int {
	if tasksMax, ok := config["tasks.max"].(string); ok {
		var count int
		fmt.Sscanf(tasksMax, "%d", &count)
		return count
	}
	return 1 // Default
}

func getConnectWorkerCount(client *http.Client, config model.KafkaConnectConfig) int {
	// Try to get worker information from admin/cluster endpoint
	clusterURL := fmt.Sprintf("%s/admin/cluster", config.URL)
	req, err := http.NewRequest("GET", clusterURL, nil)
	if err != nil {
		return 1 // Default to single worker
	}

	httpauth.ApplyKafkaConnectAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 1 // Default to single worker
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var clusterInfo map[string]interface{}
		if json.Unmarshal(body, &clusterInfo) == nil {
			// Try to get worker count from workers array
			if workers, ok := clusterInfo["workers"].([]interface{}); ok {
				return len(workers)
			}
		}
	}

	// Fallback: check connectors endpoint which might give us a hint
	// If we can connect, there's at least 1 worker
	connectorsURL := fmt.Sprintf("%s/connectors", config.URL)
	req, err = http.NewRequest("GET", connectorsURL, nil)
	if err == nil {
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return 1 // At least 1 worker is running
			}
		}
	}

	return 1 // Default to single worker
}
