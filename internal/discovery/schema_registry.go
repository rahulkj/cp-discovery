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

	// Fetch additional detailed information
	if detailed {
		additionalInfo := fetchSchemaRegistryAdditionalInfo(client, config, report.Subjects)
		report.AdditionalInfo = &additionalInfo
	}

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

// fetchSchemaRegistryAdditionalInfo collects extended information about Schema Registry
func fetchSchemaRegistryAdditionalInfo(client *http.Client, config model.SchemaRegistryConfig, subjects []string) model.SchemaRegistryAdditionalInfo {
	info := model.SchemaRegistryAdditionalInfo{
		Subjects:            make([]model.SubjectDetail, 0),
		CompatibilityLevels: make(map[string]string),
		Contexts:            make([]string, 0),
		Config:              make(map[string]string),
	}

	// Get global compatibility level
	info.GlobalCompatibility = getGlobalCompatibility(client, config)

	// Get cluster information
	info.ClusterInfo = getClusterInfo(client, config)

	// Get global config
	info.Config = getGlobalConfig(client, config)

	// Get contexts
	info.Contexts = getContexts(client, config)

	// For each subject, get detailed information
	for _, subject := range subjects {
		subjectDetail := getSubjectDetail(client, config, subject)
		if subjectDetail.Subject != "" {
			info.Subjects = append(info.Subjects, subjectDetail)
			// Store compatibility level
			if subjectDetail.Compatibility != "" {
				info.CompatibilityLevels[subject] = subjectDetail.Compatibility
			}
		}
	}

	return info
}

// getGlobalCompatibility gets the global compatibility level
func getGlobalCompatibility(client *http.Client, config model.SchemaRegistryConfig) string {
	url := fmt.Sprintf("%s/config", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, _ := io.ReadAll(resp.Body)
	var configResp struct {
		CompatibilityLevel string `json:"compatibilityLevel"`
	}
	if json.Unmarshal(body, &configResp) == nil {
		return configResp.CompatibilityLevel
	}

	return ""
}

// getClusterInfo gets Schema Registry cluster information
func getClusterInfo(client *http.Client, config model.SchemaRegistryConfig) model.SRClusterInfo {
	clusterInfo := model.SRClusterInfo{}

	url := fmt.Sprintf("%s/clusterStatus", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return clusterInfo
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return clusterInfo
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return clusterInfo
	}

	body, _ := io.ReadAll(resp.Body)
	var clusterStatus SchemaRegistryClusterStatus
	if json.Unmarshal(body, &clusterStatus) == nil {
		clusterInfo.Leader = clusterStatus.Leader
		clusterInfo.IsLeader = clusterStatus.IsLeader
		clusterInfo.Members = clusterStatus.Members
	}

	return clusterInfo
}

// getGlobalConfig gets global Schema Registry config
func getGlobalConfig(client *http.Client, config model.SchemaRegistryConfig) map[string]string {
	configMap := make(map[string]string)

	url := fmt.Sprintf("%s/config", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return configMap
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return configMap
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return configMap
	}

	body, _ := io.ReadAll(resp.Body)
	var rawConfig map[string]interface{}
	if json.Unmarshal(body, &rawConfig) == nil {
		for key, val := range rawConfig {
			configMap[key] = fmt.Sprintf("%v", val)
		}
	}

	return configMap
}

// getContexts gets all Schema Registry contexts
func getContexts(client *http.Client, config model.SchemaRegistryConfig) []string {
	url := fmt.Sprintf("%s/contexts", config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []string{}
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{}
	}

	body, _ := io.ReadAll(resp.Body)
	var contexts []string
	if json.Unmarshal(body, &contexts) == nil {
		return contexts
	}

	return []string{}
}

// getSubjectDetail gets detailed information for a subject
func getSubjectDetail(client *http.Client, config model.SchemaRegistryConfig, subject string) model.SubjectDetail {
	detail := model.SubjectDetail{
		Subject:  subject,
		Versions: make([]int, 0),
	}

	// Get all versions for the subject
	versionsURL := fmt.Sprintf("%s/subjects/%s/versions", config.URL, subject)
	req, err := http.NewRequest("GET", versionsURL, nil)
	if err != nil {
		return detail
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return detail
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var versions []int
		if json.Unmarshal(body, &versions) == nil {
			detail.Versions = versions
			if len(versions) > 0 {
				detail.LatestVersion = versions[len(versions)-1]
			}
		}
	}

	// Get latest schema
	if detail.LatestVersion > 0 {
		detail.LatestSchema = getSchemaVersion(client, config, subject, detail.LatestVersion)
	}

	// Get subject-level compatibility
	compatURL := fmt.Sprintf("%s/config/%s", config.URL, subject)
	req, err = http.NewRequest("GET", compatURL, nil)
	if err == nil {
		httpauth.ApplySchemaRegistryAuth(req, config)
		resp, err = client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var compatResp struct {
					CompatibilityLevel string `json:"compatibilityLevel"`
				}
				if json.Unmarshal(body, &compatResp) == nil {
					detail.Compatibility = compatResp.CompatibilityLevel
				}
			}
		}
	}

	return detail
}

// getSchemaVersion gets a specific version of a schema
func getSchemaVersion(client *http.Client, config model.SchemaRegistryConfig, subject string, version int) model.SchemaDetail {
	schemaDetail := model.SchemaDetail{
		Version:    version,
		References: make([]model.SchemaReference, 0),
		Metadata:   make(map[string]interface{}),
	}

	url := fmt.Sprintf("%s/subjects/%s/versions/%d", config.URL, subject, version)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return schemaDetail
	}

	httpauth.ApplySchemaRegistryAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return schemaDetail
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return schemaDetail
	}

	body, _ := io.ReadAll(resp.Body)
	var schemaResp struct {
		Subject    string                   `json:"subject"`
		ID         int                      `json:"id"`
		Version    int                      `json:"version"`
		Schema     string                   `json:"schema"`
		SchemaType string                   `json:"schemaType"`
		References []model.SchemaReference  `json:"references"`
		Metadata   map[string]interface{}   `json:"metadata"`
	}

	if json.Unmarshal(body, &schemaResp) == nil {
		schemaDetail.ID = schemaResp.ID
		schemaDetail.Version = schemaResp.Version
		schemaDetail.Schema = schemaResp.Schema
		schemaDetail.SchemaType = schemaResp.SchemaType
		schemaDetail.References = schemaResp.References
		schemaDetail.Metadata = schemaResp.Metadata
	}

	return schemaDetail
}
