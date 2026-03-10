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

type ConnectInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

// ConnectorExpandedStatus represents the response from /connectors?expand=status
type ConnectorExpandedStatus struct {
	Status ConnectorStatusDetail `json:"status"`
}

type ConnectorStatusDetail struct {
	Name      string             `json:"name"`
	Connector ConnectorStateInfo `json:"connector"`
	Tasks     []TaskStatusInfo   `json:"tasks"`
	Type      string             `json:"type"` // "source" or "sink"
}

type ConnectorStateInfo struct {
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Version  string `json:"version"`
}

type TaskStatusInfo struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Version  string `json:"version"`
}

// ConnectorExpandedInfo represents the response from /connectors?expand=info
type ConnectorExpandedInfo struct {
	Info ConnectorInfoDetail `json:"info"`
}

type ConnectorInfoDetail struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
	Tasks  []TaskInfo             `json:"tasks"`
	Type   string                 `json:"type"` // "source" or "sink"
}

type TaskInfo struct {
	Connector string `json:"connector"`
	Task      int    `json:"task"`
}

func DiscoverKafkaConnect(config model.KafkaConnectConfig, detailed bool) (model.KafkaConnectReport, error) {
	report := model.KafkaConnectReport{
		Available:        false,
		Connectors:       make([]model.ConnectorInfo, 0),
		Replicators:      make([]model.ReplicatorInfo, 0),
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

	// First call: Get connectors with status information
	statusURL := fmt.Sprintf("%s/connectors?expand=status", config.URL)
	req, err = http.NewRequest("GET", statusURL, nil)
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

	var connectorsStatus map[string]ConnectorExpandedStatus
	body, _ = io.ReadAll(resp.Body)
	if json.Unmarshal(body, &connectorsStatus) != nil {
		return report, nil
	}

	// Second call: Get connectors with info/config information
	infoURL := fmt.Sprintf("%s/connectors?expand=info", config.URL)
	req, err = http.NewRequest("GET", infoURL, nil)
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

	var connectorsInfo map[string]ConnectorExpandedInfo
	body, _ = io.ReadAll(resp.Body)
	if json.Unmarshal(body, &connectorsInfo) != nil {
		return report, nil
	}

	report.TotalConnectors = len(connectorsStatus)

	// Process each connector by merging status and info data
	for connectorName, statusData := range connectorsStatus {
		// Get connector type from status
		connectorType := statusData.Status.Type

		if connectorType == "source" {
			report.SourceConnectors++
		} else if connectorType == "sink" {
			report.SinkConnectors++
		}

		// Get connector state
		state := "unknown"
		if statusData.Status.Connector.State != "" {
			state = statusData.Status.Connector.State
		}

		// Get task count from the tasks array
		taskCount := len(statusData.Status.Tasks)

		// Extract connector.class and quickstart from info data (if available)
		connectorClass := ""
		quickstart := ""
		var connectorConfig map[string]interface{}
		if infoData, ok := connectorsInfo[connectorName]; ok {
			connectorConfig = infoData.Info.Config
			if class, ok := infoData.Info.Config["connector.class"].(string); ok {
				connectorClass = class
			}
			if qs, ok := infoData.Info.Config["quickstart"].(string); ok {
				quickstart = qs
			}
		}

		// Always include connector info (name, type, tasks, state, connector.class, quickstart)
		info := model.ConnectorInfo{
			Name:           connectorName,
			Type:           connectorType,
			State:          state,
			Tasks:          taskCount,
			ConnectorClass: connectorClass,
			Quickstart:     quickstart,
		}

		report.Connectors = append(report.Connectors, info)

		// Check if this is a Replicator connector
		if isReplicatorConnector(connectorClass) && connectorConfig != nil {
			replicatorInfo := extractReplicatorInfo(connectorName, state, taskCount, connectorConfig)
			report.Replicators = append(report.Replicators, replicatorInfo)
		}
	}

	report.ReplicatorCount = len(report.Replicators)

	return report, nil
}

func isReplicatorConnector(connectorClass string) bool {
	if connectorClass == "" {
		return false
	}

	// Check if connector class is a Confluent Replicator
	replicatorClasses := []string{
		"io.confluent.connect.replicator.ReplicatorSourceConnector",
		"com.confluent.connect.replicator.ReplicatorSourceConnector",
	}

	for _, replicatorClass := range replicatorClasses {
		if connectorClass == replicatorClass {
			return true
		}
	}

	// Also check if "Replicator" is in the class name
	if len(connectorClass) > 0 && contains(connectorClass, "Replicator") {
		return true
	}

	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexOfString(s, substr) >= 0)
}

func indexOfString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func extractReplicatorInfo(name, state string, tasks int, config map[string]interface{}) model.ReplicatorInfo {
	replicatorInfo := model.ReplicatorInfo{
		Name:  name,
		State: state,
		Tasks: tasks,
	}

	// Extract source and destination cluster information
	if srcCluster, ok := config["src.kafka.bootstrap.servers"].(string); ok {
		replicatorInfo.SourceCluster = srcCluster
	}
	if destCluster, ok := config["dest.kafka.bootstrap.servers"].(string); ok {
		replicatorInfo.DestinationCluster = destCluster
	}

	// Extract topic whitelist/blacklist
	if whitelist, ok := config["topic.whitelist"].(string); ok {
		replicatorInfo.TopicWhitelist = whitelist
	}
	if blacklist, ok := config["topic.blacklist"].(string); ok {
		replicatorInfo.TopicBlacklist = blacklist
	}

	// Extract topic rename format
	if renameFormat, ok := config["topic.rename.format"].(string); ok {
		replicatorInfo.TopicRenameFormat = renameFormat
	}

	return replicatorInfo
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
