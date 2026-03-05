package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"github.com/rahulkj/cp-discovery/internal/model"
	httpauth "github.com/rahulkj/cp-discovery/internal/http"
)

type PrometheusBuildInfo struct {
	Status string `json:"status"`
	Data   struct {
		Version   string `json:"version"`
		Revision  string `json:"revision"`
		GoVersion string `json:"goVersion"`
	} `json:"data"`
}

type PrometheusTargetsResponse struct {
	Status string `json:"status"`
	Data   struct {
		ActiveTargets []PrometheusTarget `json:"activeTargets"`
	} `json:"data"`
}

type PrometheusTarget struct {
	Health string `json:"health"`
	Labels map[string]string `json:"labels"`
}

type PrometheusRuntimeInfo struct {
	Status string `json:"status"`
	Data   struct {
		ReloadConfigSuccess bool   `json:"reloadConfigSuccess"`
		LastConfigTime      string `json:"lastConfigTime"`
	} `json:"data"`
}

type PrometheusQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string           `json:"resultType"`
		Result     []PrometheusResult `json:"result"`
	} `json:"data"`
}

type PrometheusResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"` // [timestamp, value]
}

func DiscoverPrometheus(config model.PrometheusConfig, detailed bool) (model.PrometheusReport, error) {
	report := model.PrometheusReport{
		Available: false,
	}

	if config.URL == "" {
		return report, fmt.Errorf("prometheus URL not configured")
	}

	report.URL = config.URL

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if Prometheus is available and get version
	buildInfoURL := fmt.Sprintf("%s/api/v1/status/buildinfo", config.URL)
	req, err := http.NewRequest("GET", buildInfoURL, nil)
	if err != nil {
		return report, fmt.Errorf("creating request: %w", err)
	}

	httpauth.ApplyPrometheusAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return report, fmt.Errorf("connecting to prometheus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("prometheus returned status: %d", resp.StatusCode)
	}

	report.Available = true

	// Parse version info
	body, _ := io.ReadAll(resp.Body)
	var buildInfo PrometheusBuildInfo
	if json.Unmarshal(body, &buildInfo) == nil {
		report.Version = buildInfo.Data.Version
	}

	// Get targets information
	targetsURL := fmt.Sprintf("%s/api/v1/targets", config.URL)
	req, err = http.NewRequest("GET", targetsURL, nil)
	if err == nil {
		httpauth.ApplyPrometheusAuth(req, config)

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var targetsResp PrometheusTargetsResponse
				if json.Unmarshal(body, &targetsResp) == nil {
					for _, target := range targetsResp.Data.ActiveTargets {
						if target.Health == "up" {
							report.TargetsUp++
						} else {
							report.TargetsDown++
						}
					}
				}
			}
		}
	}

	// Check for HA setup by querying runtime info
	runtimeURL := fmt.Sprintf("%s/api/v1/status/runtimeinfo", config.URL)
	req, err = http.NewRequest("GET", runtimeURL, nil)
	if err == nil {
		httpauth.ApplyPrometheusAuth(req, config)

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			// If runtime info is available, check for HA indicators
			// This is a simplified check - in production, you'd query multiple Prometheus instances
			report.NodeCount = 1
			report.HighAvailability = false
		}
	}

	// Fetch cluster metrics if detailed mode is enabled
	if detailed {
		clusterMetrics := getClusterMetricsFromPrometheus(client, config)
		report.ClusterMetrics = clusterMetrics
	}

	return report, nil
}

func getClusterMetricsFromPrometheus(client *http.Client, config model.PrometheusConfig) model.PrometheusClusterMetrics {
	metrics := model.PrometheusClusterMetrics{}

	// Query throughput metrics
	metrics.BytesInPerSec = queryPrometheusRate(client, config, "kafka_server_brokertopicmetrics_bytesin_total")
	metrics.BytesOutPerSec = queryPrometheusRate(client, config, "kafka_server_brokertopicmetrics_bytesout_total")
	metrics.MessagesInPerSec = queryPrometheusRate(client, config, "kafka_server_brokertopicmetrics_messagesin_total")

	// Query controller metrics
	metrics.ActiveControllerCount = int(queryPrometheusGauge(client, config, "kafka_controller_kafkacontroller_activecontrollercount"))

	// Query partition metrics
	metrics.UnderReplicatedPartitions = int(queryPrometheusGauge(client, config, "kafka_server_replicamanager_underreplicatedpartitions"))
	metrics.OfflinePartitions = int(queryPrometheusGauge(client, config, "kafka_server_replicamanager_offlinepartitionscount"))
	metrics.TotalPartitions = int(queryPrometheusGauge(client, config, "kafka_server_replicamanager_partitioncount"))
	metrics.LeaderCount = int(queryPrometheusGauge(client, config, "kafka_server_replicamanager_leadercount"))

	// Query broker metrics - count unique brokers from targets
	metrics.TotalBrokers = queryPrometheusCountBrokers(client, config)
	metrics.OnlineBrokers = queryPrometheusOnlineBrokers(client, config)

	// Query consumer lag metrics
	metrics.TotalConsumerLag = int64(queryPrometheusGauge(client, config, "kafka_consumergroup_lag_sum"))
	metrics.ConsumerGroups = queryPrometheusCountConsumerGroups(client, config)

	// Query JVM metrics
	metrics.AvgHeapUsedPercent = queryPrometheusAvgHeapUsage(client, config)
	metrics.AvgCPUUsedPercent = queryPrometheusAvgCPUUsage(client, config)

	return metrics
}

func queryPrometheusRate(client *http.Client, config model.PrometheusConfig, metric string) float64 {
	// Query rate over last 5 minutes
	query := fmt.Sprintf("rate(%s[5m])", metric)
	return queryPrometheusSum(client, config, query)
}

func queryPrometheusGauge(client *http.Client, config model.PrometheusConfig, metric string) float64 {
	return queryPrometheusSum(client, config, metric)
}

func queryPrometheusSum(client *http.Client, config model.PrometheusConfig, query string) float64 {
	encodedQuery := url.QueryEscape(fmt.Sprintf("sum(%s)", query))
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", config.URL, encodedQuery)
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return 0
	}

	httpauth.ApplyPrometheusAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	body, _ := io.ReadAll(resp.Body)
	var queryResp PrometheusQueryResponse
	if json.Unmarshal(body, &queryResp) != nil {
		return 0
	}

	if len(queryResp.Data.Result) == 0 {
		return 0
	}

	// Extract value from result
	if len(queryResp.Data.Result[0].Value) >= 2 {
		if valueStr, ok := queryResp.Data.Result[0].Value[1].(string); ok {
			var value float64
			fmt.Sscanf(valueStr, "%f", &value)
			return value
		}
	}

	return 0
}

func queryPrometheusCountBrokers(client *http.Client, config model.PrometheusConfig) int {
	// Count unique Kafka brokers by counting unique instances with kafka_server metrics
	query := "count(count by (instance) (kafka_server_brokertopicmetrics_bytesin_total))"
	encodedQuery := url.QueryEscape(query)
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", config.URL, encodedQuery)
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return 0
	}

	httpauth.ApplyPrometheusAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	body, _ := io.ReadAll(resp.Body)
	var queryResp PrometheusQueryResponse
	if json.Unmarshal(body, &queryResp) != nil {
		return 0
	}

	if len(queryResp.Data.Result) == 0 {
		return 0
	}

	if len(queryResp.Data.Result[0].Value) >= 2 {
		if valueStr, ok := queryResp.Data.Result[0].Value[1].(string); ok {
			var count int
			fmt.Sscanf(valueStr, "%d", &count)
			return count
		}
	}

	return 0
}

func queryPrometheusOnlineBrokers(client *http.Client, config model.PrometheusConfig) int {
	// Count brokers that are up (have recent metrics)
	query := "count(up{job=~\".*kafka.*\"} == 1)"
	encodedQuery := url.QueryEscape(query)
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", config.URL, encodedQuery)
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return 0
	}

	httpauth.ApplyPrometheusAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	body, _ := io.ReadAll(resp.Body)
	var queryResp PrometheusQueryResponse
	if json.Unmarshal(body, &queryResp) != nil {
		return 0
	}

	if len(queryResp.Data.Result) == 0 {
		return 0
	}

	if len(queryResp.Data.Result[0].Value) >= 2 {
		if valueStr, ok := queryResp.Data.Result[0].Value[1].(string); ok {
			var count int
			fmt.Sscanf(valueStr, "%d", &count)
			return count
		}
	}

	return 0
}

func queryPrometheusCountConsumerGroups(client *http.Client, config model.PrometheusConfig) int {
	// Count unique consumer groups
	query := "count(count by (consumergroup) (kafka_consumergroup_lag))"
	encodedQuery := url.QueryEscape(query)
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", config.URL, encodedQuery)
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return 0
	}

	httpauth.ApplyPrometheusAuth(req, config)

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	body, _ := io.ReadAll(resp.Body)
	var queryResp PrometheusQueryResponse
	if json.Unmarshal(body, &queryResp) != nil {
		return 0
	}

	if len(queryResp.Data.Result) == 0 {
		return 0
	}

	if len(queryResp.Data.Result[0].Value) >= 2 {
		if valueStr, ok := queryResp.Data.Result[0].Value[1].(string); ok {
			var count int
			fmt.Sscanf(valueStr, "%d", &count)
			return count
		}
	}

	return 0
}

func queryPrometheusAvgHeapUsage(client *http.Client, config model.PrometheusConfig) float64 {
	// Average JVM heap usage across all Kafka brokers
	query := "avg(jvm_memory_bytes_used{area=\"heap\"} / jvm_memory_bytes_max{area=\"heap\"} * 100)"
	return queryPrometheusSum(client, config, query)
}

func queryPrometheusAvgCPUUsage(client *http.Client, config model.PrometheusConfig) float64 {
	// Average CPU usage across all Kafka brokers
	query := "avg(rate(process_cpu_seconds_total{job=~\".*kafka.*\"}[5m]) * 100)"
	return queryPrometheusSum(client, config, query)
}
