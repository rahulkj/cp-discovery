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

type SchemaRegistryInfo struct {
	Version string `json:"version"`
}

type SchemaRegistryMode struct {
	Mode string `json:"mode"`
}

type SchemaRegistryClusterStatus struct {
	Leader    string                 `json:"leader"`
	IsLeader  bool                   `json:"isLeader"`
	Members   []string               `json:"members"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func DiscoverSchemaRegistry(config model.SchemaRegistryConfig, detailed bool) (model.SchemaRegistryReport, error) {
	report := model.SchemaRegistryReport{
		Available: false,
		Subjects:  make([]string, 0),
	}

	if config.URL == "" {
		return report, fmt.Errorf("schema registry URL not configured")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if Schema Registry is available and get version
	versionURL := fmt.Sprintf("%s/", config.URL)
	req, err := http.NewRequest("GET", versionURL, nil)
	if err != nil {
		return report, fmt.Errorf("creating request: %w", err)
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return report, fmt.Errorf("connecting to schema registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("schema registry returned status: %d", resp.StatusCode)
	}

	report.Available = true

	// Get version info
	var versionInfo SchemaRegistryInfo
	body, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(body, &versionInfo) == nil {
		report.Version = versionInfo.Version
	}

	// Get mode
	modeURL := fmt.Sprintf("%s/mode", config.URL)
	req, err = http.NewRequest("GET", modeURL, nil)
	if err == nil {
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}

		resp, err = client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var mode SchemaRegistryMode
			body, _ := io.ReadAll(resp.Body)
			if json.Unmarshal(body, &mode) == nil {
				report.Mode = mode.Mode
			}
		}
	}

	// Try to get cluster status to count nodes
	clusterURL := fmt.Sprintf("%s/v1/metadata/id", config.URL)
	req, err = http.NewRequest("GET", clusterURL, nil)
	if err == nil {
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				// Schema Registry is available, default to 1 node
				report.NodeCount = 1
			}
		}
	}

	// Try newer cluster status endpoint for multi-node deployments
	clusterStatusURL := fmt.Sprintf("%s/clusterStatus", config.URL)
	req, err = http.NewRequest("GET", clusterStatusURL, nil)
	if err == nil {
		if config.BasicAuthUsername != "" {
			req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var clusterStatus SchemaRegistryClusterStatus
				if json.Unmarshal(body, &clusterStatus) == nil {
					if len(clusterStatus.Members) > 0 {
						report.NodeCount = len(clusterStatus.Members)
					}
				}
			}
		}
	}

	// Get subjects (schema count)
	subjectsURL := fmt.Sprintf("%s/subjects", config.URL)
	req, err = http.NewRequest("GET", subjectsURL, nil)
	if err != nil {
		return report, nil
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err = client.Do(req)
	if err != nil {
		return report, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var subjects []string
		body, _ := io.ReadAll(resp.Body)
		if json.Unmarshal(body, &subjects) == nil {
			report.TotalSchemas = len(subjects)
			if detailed {
				report.Subjects = subjects
			}
		}
	}

	return report, nil
}
