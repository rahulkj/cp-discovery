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

type AlertmanagerStatus struct {
	Cluster struct {
		Status string   `json:"status"`
		Peers  []Peer   `json:"peers"`
	} `json:"cluster"`
	VersionInfo struct {
		Version   string `json:"version"`
		Revision  string `json:"revision"`
		Branch    string `json:"branch"`
		GoVersion string `json:"goVersion"`
	} `json:"versionInfo"`
	Uptime string `json:"uptime"`
}

type Peer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type AlertmanagerAlertsResponse struct {
	Status string  `json:"status"`
	Data   []Alert `json:"data"`
}

type Alert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	State        string            `json:"state"`
	ActiveAt     string            `json:"activeAt"`
	Value        string            `json:"value"`
}

func DiscoverAlertmanager(config model.AlertmanagerConfig, detailed bool) (model.AlertmanagerReport, error) {
	report := model.AlertmanagerReport{
		Available:    false,
		ClusterPeers: make([]string, 0),
	}

	if config.URL == "" {
		return report, fmt.Errorf("alertmanager URL not configured")
	}

	report.URL = config.URL

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if Alertmanager is available and get status
	statusURL := fmt.Sprintf("%s/api/v2/status", config.URL)
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return report, fmt.Errorf("creating request: %w", err)
	}

	httpauth.ApplyAlertmanagerAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return report, fmt.Errorf("connecting to alertmanager: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("alertmanager returned status: %d", resp.StatusCode)
	}

	report.Available = true

	// Parse status information
	body, _ := io.ReadAll(resp.Body)
	var status AlertmanagerStatus
	if json.Unmarshal(body, &status) == nil {
		report.Version = status.VersionInfo.Version
		report.ClusterSize = len(status.Cluster.Peers)

		// Extract peer addresses
		for _, peer := range status.Cluster.Peers {
			report.ClusterPeers = append(report.ClusterPeers, peer.Address)
		}

		// If no peers but Alertmanager is running, it's a single node
		if report.ClusterSize == 0 {
			report.ClusterSize = 1
		}
	}

	// Get active alerts
	alertsURL := fmt.Sprintf("%s/api/v2/alerts", config.URL)
	req, err = http.NewRequest("GET", alertsURL, nil)
	if err == nil {
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var alerts []Alert
				if json.Unmarshal(body, &alerts) == nil {
					// Count only active/firing alerts
					for _, alert := range alerts {
						if alert.State == "active" || alert.State == "firing" {
							report.ActiveAlerts++
						}
					}
				}
			}
		}
	}

	return report, nil
}
