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

type SchemaExporter struct {
	Name          string                 `json:"name"`
	Subjects      []string               `json:"subjects"`
	SubjectFormat string                 `json:"subjectFormat"`
	ContextType   string                 `json:"contextType"`
	Context       string                 `json:"context"`
	Config        map[string]interface{} `json:"config"`
}

func DiscoverSchemaRegistry(config model.SchemaRegistryConfig, detailed bool) (model.SchemaRegistryReport, error) {
	report := model.SchemaRegistryReport{
		Available:       false,
		Subjects:        make([]string, 0),
		SchemaExporters: make([]model.SchemaLinkInfo, 0),
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

	// Get schema exporters (for schema linking)
	exporters := getSchemaExporters(client, config, detailed)
	report.SchemaExporters = exporters
	report.ExporterCount = len(exporters)

	return report, nil
}

func getSchemaExporters(client *http.Client, config model.SchemaRegistryConfig, detailed bool) []model.SchemaLinkInfo {
	exportersURL := fmt.Sprintf("%s/exporters", config.URL)
	req, err := http.NewRequest("GET", exportersURL, nil)
	if err != nil {
		return []model.SchemaLinkInfo{}
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []model.SchemaLinkInfo{}
	}
	defer resp.Body.Close()

	// If endpoint doesn't exist (404) or not supported, return empty list
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return []model.SchemaLinkInfo{}
	}

	if resp.StatusCode != http.StatusOK {
		return []model.SchemaLinkInfo{}
	}

	body, _ := io.ReadAll(resp.Body)
	var exporterNames []string
	if json.Unmarshal(body, &exporterNames) != nil {
		return []model.SchemaLinkInfo{}
	}

	exporters := make([]model.SchemaLinkInfo, 0, len(exporterNames))

	for _, exporterName := range exporterNames {
		// Get detailed exporter information
		exporterDetail := getExporterDetail(client, config, exporterName)
		if exporterDetail != nil {
			schemaLink := model.SchemaLinkInfo{
				ExporterName:  exporterDetail.Name,
				Subjects:      exporterDetail.Subjects,
				SubjectFormat: exporterDetail.SubjectFormat,
				ContextType:   exporterDetail.ContextType,
				Context:       exporterDetail.Context,
			}

			// Convert config to map[string]string
			if exporterDetail.Config != nil && detailed {
				configs := make(map[string]string)
				for key, val := range exporterDetail.Config {
					if strVal, ok := val.(string); ok {
						configs[key] = strVal
					} else {
						configs[key] = fmt.Sprintf("%v", val)
					}
				}
				schemaLink.Config = configs
			}

			exporters = append(exporters, schemaLink)
		}
	}

	return exporters
}

func getExporterDetail(client *http.Client, config model.SchemaRegistryConfig, exporterName string) *SchemaExporter {
	exporterURL := fmt.Sprintf("%s/exporters/%s", config.URL, exporterName)
	req, err := http.NewRequest("GET", exporterURL, nil)
	if err != nil {
		return nil
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	var exporter SchemaExporter
	if json.Unmarshal(body, &exporter) == nil {
		return &exporter
	}

	return nil
}
