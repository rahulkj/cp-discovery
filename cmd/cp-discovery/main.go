package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	configPkg "github.com/rahulkj/cp-discovery/internal/config"
	"github.com/rahulkj/cp-discovery/internal/discovery"
	"github.com/rahulkj/cp-discovery/internal/model"
	"gopkg.in/yaml.v3"
)

func main() {
	configFile := flag.String("config", "configs/config.yaml", "Path to configuration file")
	outputFile := flag.String("output", "", "Output file path (overrides config file setting)")
	outputFormat := flag.String("format", "", "Output format: json or yaml (overrides config file setting)")
	detailed := flag.Bool("detailed", false, "Enable detailed discovery (overrides config file setting)")
	viewReport := flag.Bool("view", false, "Open report in web browser after discovery")
	viewPort := flag.Int("port", 8080, "Port for web view server (used with -view)")
	viewOnly := flag.String("view-file", "", "View existing report file in browser (skip discovery)")
	flag.Parse()

	// View-only mode: just open existing file in browser
	if *viewOnly != "" {
		if err := viewReportFile(*viewOnly, *viewPort); err != nil {
			log.Fatalf("Failed to view report: %v", err)
		}
		return
	}

	cfg, err := configPkg.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// If -view is specified without -output, create a temporary file
	useTempFile := false
	var tempFilePath string
	if *viewReport && *outputFile == "" && cfg.Output.File == "" {
		// Create a temporary file for viewing
		tempFile, err := os.CreateTemp("", "cp-discovery-*.json")
		if err != nil {
			log.Fatalf("Failed to create temporary file: %v", err)
		}
		tempFilePath = tempFile.Name()
		tempFile.Close()
		useTempFile = true
		cfg.Output.File = tempFilePath
		// Default to JSON format for temp files
		if *outputFormat == "" && cfg.Output.Format == "" {
			cfg.Output.Format = "json"
		}
		fmt.Printf("Using temporary report file: %s\n", tempFilePath)
	}

	// Override config with command-line flags if provided
	if *outputFile != "" {
		cfg.Output.File = *outputFile
	}
	if *outputFormat != "" {
		cfg.Output.Format = *outputFormat
	}
	if *detailed {
		cfg.Output.Detailed = true
	}

	report := discoverClusters(cfg)

	if err := outputReport(report, cfg.Output); err != nil {
		log.Fatalf("Failed to output report: %v", err)
	}

	printSummary(report)

	// Open in browser if requested
	if *viewReport {
		reportFile := cfg.Output.File
		if reportFile == "" {
			reportFile = "discovery-report.json" // Default
		}

		// Set up signal handling for cleanup if using temp file
		if useTempFile {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-sigChan
				fmt.Printf("\n\nReceived interrupt signal\n")
				fmt.Printf("Cleaning up temporary file: %s\n", tempFilePath)
				os.Remove(tempFilePath)
				os.Exit(0)
			}()
		}

		if err := viewReportFile(reportFile, *viewPort); err != nil {
			// Clean up temp file if server fails to start
			if useTempFile {
				os.Remove(tempFilePath)
			}
			log.Fatalf("Failed to view report: %v", err)
		}
	}
}

func discoverClusters(cfg *model.Config) *model.DiscoveryReport {
	report := &model.DiscoveryReport{
		Timestamp:     time.Now().Format(time.RFC3339),
		Clusters:      make([]model.ClusterReport, 0),
		TotalClusters: len(cfg.Clusters),
	}

	detailed := cfg.Output.Detailed

	var wg sync.WaitGroup
	clusterReports := make([]model.ClusterReport, len(cfg.Clusters))

	for i, cluster := range cfg.Clusters {
		wg.Add(1)
		go func(idx int, clusterConfig model.ClusterConfig) {
			defer wg.Done()
			clusterReports[idx] = discoverCluster(clusterConfig, detailed)
		}(i, cluster)
	}

	wg.Wait()

	report.Clusters = clusterReports
	return report
}

func discoverCluster(config model.ClusterConfig, detailed bool) model.ClusterReport {
	report := model.ClusterReport{
		Name:   config.Name,
		Status: "healthy",
		Errors: make([]string, 0),
	}

	var wg sync.WaitGroup

	// Discover Kafka (always required)
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafkaReport, err := discovery.DiscoverKafka(config.Kafka, detailed)
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("Kafka: %v", err))
		} else {
			report.Kafka = kafkaReport
		}
	}()

	// Discover Schema Registry (only if enabled)
	if configPkg.ShouldDiscoverSchemaRegistry(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			srReport, err := discovery.DiscoverSchemaRegistry(config.SchemaRegistry, detailed)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("Schema Registry: %v", err))
			} else {
				report.SchemaRegistry = srReport
			}
		}()
	}

	// Discover Kafka Connect (only if enabled)
	if configPkg.ShouldDiscoverKafkaConnect(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connectReport, err := discovery.DiscoverKafkaConnect(config.KafkaConnect, detailed)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("Kafka Connect: %v", err))
			} else {
				report.KafkaConnect = connectReport
			}
		}()
	}

	// Discover ksqlDB (only if enabled)
	if configPkg.ShouldDiscoverKsqlDB(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ksqlReport, err := discovery.DiscoverKsqlDB(config.KsqlDB, detailed)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("ksqlDB: %v", err))
			} else {
				report.KsqlDB = ksqlReport
			}
		}()
	}

	// Discover REST Proxy (only if enabled)
	if configPkg.ShouldDiscoverRestProxy(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			restReport, err := discovery.DiscoverRestProxy(config.RestProxy, detailed)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("REST Proxy: %v", err))
			} else {
				report.RestProxy = restReport
			}
		}()
	}

	// Discover Control Center (only if enabled)
	if configPkg.ShouldDiscoverControlCenter(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ccReport, err := discovery.DiscoverControlCenter(config.ControlCenter, detailed)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("Control Center: %v", err))
			} else {
				report.ControlCenter = ccReport
			}
		}()
	}

	// Discover Prometheus (only if enabled)
	if configPkg.ShouldDiscoverPrometheus(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			promReport, err := discovery.DiscoverPrometheus(config.Prometheus, detailed)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("Prometheus: %v", err))
			} else {
				report.Prometheus = promReport
			}
		}()
	}

	// Discover Alertmanager (only if enabled)
	if configPkg.ShouldDiscoverAlertmanager(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			amReport, err := discovery.DiscoverAlertmanager(config.Alertmanager, detailed)
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("Alertmanager: %v", err))
			} else {
				report.Alertmanager = amReport
			}
		}()
	}

	wg.Wait()

	if len(report.Errors) > 0 {
		report.Status = "partial"
	}

	return report
}

func outputReport(report *model.DiscoveryReport, config model.OutputConfig) error {
	var data []byte
	var err error

	switch config.Format {
	case "yaml":
		data, err = yaml.Marshal(report)
	case "json":
		data, err = json.MarshalIndent(report, "", "  ")
	default:
		data, err = json.MarshalIndent(report, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("marshaling report: %w", err)
	}

	if config.File != "" {
		if err := os.WriteFile(config.File, data, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Printf("Report written to %s\n", config.File)
	} else {
		fmt.Println(string(data))
	}

	return nil
}

func formatRetention(ms int64) string {
	if ms < 0 {
		return "unlimited"
	}
	hours := ms / (1000 * 60 * 60)
	if hours < 24 {
		return fmt.Sprintf("%dh", hours)
	}
	days := hours / 24
	return fmt.Sprintf("%dd", days)
}

func printSummary(report *model.DiscoveryReport) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("CONFLUENT PLATFORM DISCOVERY SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Timestamp: %s\n", report.Timestamp)
	fmt.Printf("Total Clusters: %d\n\n", report.TotalClusters)

	for _, cluster := range report.Clusters {
		fmt.Printf("Cluster: %s [%s]\n", cluster.Name, cluster.Status)
		fmt.Println(strings.Repeat("-", 80))

		// Kafka
		if cluster.Kafka.Available {
			fmt.Printf("  Kafka:\n")
			fmt.Printf("    Brokers: %d\n", cluster.Kafka.BrokerCount)
			fmt.Printf("    Controller: %s\n", cluster.Kafka.ControllerType)
			if cluster.Kafka.ZookeeperNodes > 0 {
				fmt.Printf("    ZooKeeper Nodes: %d\n", cluster.Kafka.ZookeeperNodes)
			}
			fmt.Printf("    Topics: %d (Internal: %d, External: %d)\n",
				cluster.Kafka.TopicCount, cluster.Kafka.InternalTopics, cluster.Kafka.ExternalTopics)
			fmt.Printf("    Total Partitions: %d\n", cluster.Kafka.TotalPartitions)
			if cluster.Kafka.SecurityConfig.AuthenticationMethod != "" {
				fmt.Printf("    Security: %s\n", cluster.Kafka.SecurityConfig.AuthenticationMethod)
			}

			// Network throughput
			metrics := cluster.Kafka.ClusterMetrics
			if metrics.BytesInPerSec > 0 || metrics.BytesOutPerSec > 0 {
				fmt.Printf("    Network Throughput:\n")
				fmt.Printf("      Bytes In: %.2f MB/s\n", metrics.BytesInPerSec/1024/1024)
				fmt.Printf("      Bytes Out: %.2f MB/s\n", metrics.BytesOutPerSec/1024/1024)
				if metrics.MessagesInPerSec > 0 {
					fmt.Printf("      Messages In: %.2f msg/s\n", metrics.MessagesInPerSec)
				}
			}

			// Storage
			if metrics.TotalDiskUsageBytes > 0 {
				fmt.Printf("    Storage:\n")
				fmt.Printf("      Total Disk Usage: %.2f GB\n", float64(metrics.TotalDiskUsageBytes)/1024/1024/1024)
			}

			// Broker-level details (if available in detailed mode)
			if len(cluster.Kafka.Brokers) > 0 {
				totalBrokerStorage := int64(0)
				for _, broker := range cluster.Kafka.Brokers {
					if broker.DiskUsageBytes > 0 {
						totalBrokerStorage += broker.DiskUsageBytes
					}
				}
				if totalBrokerStorage > 0 && metrics.TotalDiskUsageBytes == 0 {
					fmt.Printf("    Storage:\n")
					fmt.Printf("      Total Disk Usage: %.2f GB (from brokers)\n", float64(totalBrokerStorage)/1024/1024/1024)
				}
			}

			// Cluster health
			if metrics.UnderReplicatedPartitions > 0 {
				fmt.Printf("    Health:\n")
				fmt.Printf("      Under-Replicated Partitions: %d\n", metrics.UnderReplicatedPartitions)
			}
		}

		// Schema Registry
		if cluster.SchemaRegistry.Available {
			fmt.Printf("  Schema Registry:\n")
			fmt.Printf("    Version: %s\n", cluster.SchemaRegistry.Version)
			fmt.Printf("    Mode: %s\n", cluster.SchemaRegistry.Mode)
			if cluster.SchemaRegistry.NodeCount > 0 {
				fmt.Printf("    Nodes: %d\n", cluster.SchemaRegistry.NodeCount)
			}
			fmt.Printf("    Schemas: %d\n", cluster.SchemaRegistry.TotalSchemas)
		}

		// Kafka Connect
		if cluster.KafkaConnect.Available {
			fmt.Printf("  Kafka Connect:\n")
			fmt.Printf("    Version: %s\n", cluster.KafkaConnect.Version)
			if cluster.KafkaConnect.WorkerCount > 0 {
				fmt.Printf("    Workers: %d\n", cluster.KafkaConnect.WorkerCount)
			}
			fmt.Printf("    Connectors: %d (Source: %d, Sink: %d)\n",
				cluster.KafkaConnect.TotalConnectors,
				cluster.KafkaConnect.SourceConnectors,
				cluster.KafkaConnect.SinkConnectors)
		}

		// ksqlDB
		if cluster.KsqlDB.Available {
			fmt.Printf("  ksqlDB:\n")
			fmt.Printf("    Version: %s\n", cluster.KsqlDB.Version)
			fmt.Printf("    Queries: %d\n", cluster.KsqlDB.Queries)
			fmt.Printf("    Streams: %d\n", cluster.KsqlDB.Streams)
			fmt.Printf("    Tables: %d\n", cluster.KsqlDB.Tables)
		}

		// REST Proxy
		if cluster.RestProxy.Available {
			fmt.Printf("  REST Proxy:\n")
			fmt.Printf("    Version: %s\n", cluster.RestProxy.Version)
			fmt.Printf("    Cluster ID: %s\n", cluster.RestProxy.ClusterID)
			fmt.Printf("    Brokers: %d\n", cluster.RestProxy.BrokerCount)
			fmt.Printf("    Controller Mode: %s\n", cluster.RestProxy.ControllerMode)
			if cluster.RestProxy.ControllerCount > 0 {
				fmt.Printf("    Controllers: %d\n", cluster.RestProxy.ControllerCount)
			}
			if cluster.RestProxy.ConsumerGroupCount > 0 {
				fmt.Printf("    Consumer Groups: %d (Active: %d)\n",
					cluster.RestProxy.ConsumerGroupCount, cluster.RestProxy.ActiveConsumerGroups)
			}
			if cluster.RestProxy.AclCount > 0 {
				fmt.Printf("    ACLs: %d\n", cluster.RestProxy.AclCount)
			}

			// Security details
			if cluster.RestProxy.SecurityConfig.AuthenticationMethod != "" {
				fmt.Printf("    Security: %s\n", cluster.RestProxy.SecurityConfig.AuthenticationMethod)
			}
		}

		// Control Center
		if cluster.ControlCenter.Available {
			fmt.Printf("  Control Center:\n")
			fmt.Printf("    Version: %s\n", cluster.ControlCenter.Version)
			fmt.Printf("    URL: %s\n", cluster.ControlCenter.URL)
			if cluster.ControlCenter.MonitoredClusters > 0 {
				fmt.Printf("    Monitored Clusters: %d\n", cluster.ControlCenter.MonitoredClusters)
			}
			if cluster.ControlCenter.TotalConsumerLag > 0 {
				fmt.Printf("    Total Consumer Lag: %d\n", cluster.ControlCenter.TotalConsumerLag)
			}
		}

		// Prometheus
		if cluster.Prometheus.Available {
			fmt.Printf("  Prometheus:\n")
			fmt.Printf("    Version: %s\n", cluster.Prometheus.Version)
			fmt.Printf("    URL: %s\n", cluster.Prometheus.URL)
			fmt.Printf("    Targets Up: %d\n", cluster.Prometheus.TargetsUp)
			fmt.Printf("    Targets Down: %d\n", cluster.Prometheus.TargetsDown)

			// Display cluster metrics if available
			metrics := &cluster.Prometheus.ClusterMetrics
			if metrics.TotalBrokers > 0 || metrics.TotalPartitions > 0 {
				fmt.Printf("    Cluster Metrics:\n")
				if metrics.BytesInPerSec > 0 || metrics.BytesOutPerSec > 0 {
					fmt.Printf("      Throughput: %.2f MB/s in, %.2f MB/s out\n",
						metrics.BytesInPerSec/1024/1024, metrics.BytesOutPerSec/1024/1024)
				}
				if metrics.MessagesInPerSec > 0 {
					fmt.Printf("      Messages: %.2f msg/s in\n", metrics.MessagesInPerSec)
				}
				if metrics.ActiveControllerCount > 0 {
					fmt.Printf("      Active Controllers: %d\n", metrics.ActiveControllerCount)
				}
				if metrics.TotalBrokers > 0 {
					fmt.Printf("      Brokers: %d online / %d total\n",
						metrics.OnlineBrokers, metrics.TotalBrokers)
				}
				if metrics.TotalPartitions > 0 {
					fmt.Printf("      Partitions: %d total", metrics.TotalPartitions)
					if metrics.UnderReplicatedPartitions > 0 {
						fmt.Printf(" (%d under-replicated)", metrics.UnderReplicatedPartitions)
					}
					if metrics.OfflinePartitions > 0 {
						fmt.Printf(" (%d offline)", metrics.OfflinePartitions)
					}
					fmt.Println()
				}
				if metrics.TotalConsumerLag > 0 || metrics.ConsumerGroups > 0 {
					fmt.Printf("      Consumers: %d groups, lag: %d\n",
						metrics.ConsumerGroups, metrics.TotalConsumerLag)
				}
				if metrics.AvgHeapUsedPercent > 0 || metrics.AvgCPUUsedPercent > 0 {
					fmt.Printf("      JVM: %.1f%% heap, %.1f%% CPU (avg across brokers)\n",
						metrics.AvgHeapUsedPercent, metrics.AvgCPUUsedPercent)
				}
			}
		}

		// Alertmanager
		if cluster.Alertmanager.Available {
			fmt.Printf("  Alertmanager:\n")
			fmt.Printf("    Version: %s\n", cluster.Alertmanager.Version)
			if cluster.Alertmanager.ClusterSize > 0 {
				fmt.Printf("    Cluster Size: %d\n", cluster.Alertmanager.ClusterSize)
			}
			fmt.Printf("    Active Alerts: %d\n", cluster.Alertmanager.ActiveAlerts)
		}

		// Errors
		if len(cluster.Errors) > 0 {
			fmt.Printf("  Errors:\n")
			for _, err := range cluster.Errors {
				fmt.Printf("    - %s\n", err)
			}
		}

		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 80))
}

// viewReportFile opens the report in a web browser
func viewReportFile(reportPath string, port int) error {
	// Read the report file
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return fmt.Errorf("failed to read report file: %w", err)
	}

	// Parse JSON to check validity
	var report model.DiscoveryReport
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("failed to parse report file: %w", err)
	}

	// Start HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveReportHTML(w, r, &report)
	})

	http.HandleFunc("/api/report", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	addr := fmt.Sprintf("localhost:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("🌐 Web Report Viewer\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n")
	fmt.Printf("Report: %s\n", reportPath)
	fmt.Printf("Server: %s\n", url)
	fmt.Printf("Press Ctrl+C to stop the server\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")

	// Open browser
	if err := openBrowser(url); err != nil {
		fmt.Printf("⚠️  Could not open browser automatically: %v\n", err)
		fmt.Printf("Please open manually: %s\n", url)
	} else {
		fmt.Printf("✅ Opening browser...\n\n")
	}

	// Start server
	return http.ListenAndServe(addr, nil)
}

// openBrowser opens the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}

// serveReportHTML serves the HTML page
func serveReportHTML(w http.ResponseWriter, r *http.Request, report *model.DiscoveryReport) {
	html := generateReportHTML(report)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// generateReportHTML generates the HTML content
func generateReportHTML(report *model.DiscoveryReport) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Confluent Platform Discovery Report</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            min-height: 100vh;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        .header .subtitle {
            font-size: 1.1em;
            opacity: 0.9;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            padding: 30px;
            background: #f8f9fa;
            border-bottom: 2px solid #e9ecef;
        }
        .summary-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            text-align: center;
        }
        .summary-card .number {
            font-size: 2.5em;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 5px;
        }
        .summary-card .label {
            color: #6c757d;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .clusters {
            padding: 30px;
        }
        .cluster {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 25px;
            margin-bottom: 20px;
            border-left: 5px solid #667eea;
        }
        .cluster-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }
        .cluster-name {
            font-size: 1.8em;
            font-weight: bold;
            color: #2d3748;
        }
        .status-badge {
            padding: 8px 16px;
            border-radius: 20px;
            font-size: 0.9em;
            font-weight: bold;
            text-transform: uppercase;
        }
        .status-healthy {
            background: #d4edda;
            color: #155724;
        }
        .status-partial {
            background: #fff3cd;
            color: #856404;
        }
        .status-error {
            background: #f8d7da;
            color: #721c24;
        }
        .components {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        .component {
            background: white;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .component-header {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }
        .component-icon {
            width: 40px;
            height: 40px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 8px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            margin-right: 12px;
        }
        .component-title {
            font-size: 1.2em;
            font-weight: bold;
            color: #2d3748;
        }
        .component-detail {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #e9ecef;
        }
        .component-detail:last-child {
            border-bottom: none;
        }
        .detail-label {
            color: #6c757d;
            font-size: 0.9em;
        }
        .detail-value {
            font-weight: 600;
            color: #2d3748;
        }
        .available-yes {
            color: #28a745;
        }
        .available-no {
            color: #dc3545;
        }
        .tabs {
            display: flex;
            gap: 10px;
            padding: 20px 30px 0 30px;
            border-bottom: 2px solid #e9ecef;
        }
        .tab {
            padding: 10px 20px;
            cursor: pointer;
            border-radius: 8px 8px 0 0;
            transition: all 0.3s;
        }
        .tab:hover {
            background: #f8f9fa;
        }
        .tab.active {
            background: #667eea;
            color: white;
            font-weight: bold;
        }
        .tab-content {
            display: none;
        }
        .tab-content.active {
            display: block;
        }
        .json-viewer {
            background: #2d3748;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            margin: 20px;
        }
        .json-viewer pre {
            margin: 0;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 0.9em;
            line-height: 1.5;
        }
        .footer {
            text-align: center;
            padding: 20px;
            color: #6c757d;
            font-size: 0.9em;
            border-top: 2px solid #e9ecef;
        }
        .metric-good {
            color: #28a745;
            font-weight: bold;
        }
        .metric-warning {
            color: #ffc107;
            font-weight: bold;
        }
        .metric-error {
            color: #dc3545;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Confluent Platform Discovery</h1>
            <div class="subtitle" id="timestamp"></div>
        </div>

        <div class="tabs">
            <div class="tab active" onclick="switchTab('overview')">Overview</div>
            <div class="tab" onclick="switchTab('clusters')">Clusters</div>
            <div class="tab" onclick="switchTab('json')">Raw JSON</div>
        </div>

        <div id="overview-content" class="tab-content active">
            <div class="summary" id="summary"></div>
        </div>

        <div id="clusters-content" class="tab-content">
            <div class="clusters" id="clusters"></div>
        </div>

        <div id="json-content" class="tab-content">
            <div class="json-viewer">
                <pre id="json-data"></pre>
            </div>
        </div>

        <div class="footer">
            Generated by Confluent Discovery Tool v2.0.0
        </div>
    </div>

    <script>
        let reportData = null;

        fetch('/api/report')
            .then(response => response.json())
            .then(data => {
                reportData = data;
                renderReport(data);
            })
            .catch(error => {
                console.error('Error loading report:', error);
            });

        function renderReport(data) {
            document.getElementById('timestamp').textContent = 'Generated: ' + new Date(data.timestamp).toLocaleString();
            renderSummary(data);
            renderClusters(data);
            document.getElementById('json-data').textContent = JSON.stringify(data, null, 2);
        }

        function renderSummary(data) {
            var summary = document.getElementById('summary');
            var totalClusters = data.total_clusters || 0;
            var healthyClusters = data.clusters.filter(function(c) { return c.status === 'healthy'; }).length;
            var totalBrokers = data.clusters.reduce(function(sum, c) { return sum + (c.kafka && c.kafka.broker_count || 0); }, 0);
            var totalTopics = data.clusters.reduce(function(sum, c) { return sum + (c.kafka && c.kafka.topic_count || 0); }, 0);

            summary.innerHTML = '<div class="summary-card"><div class="number">' + totalClusters + '</div><div class="label">Total Clusters</div></div>' +
                '<div class="summary-card"><div class="number">' + healthyClusters + '</div><div class="label">Healthy Clusters</div></div>' +
                '<div class="summary-card"><div class="number">' + totalBrokers + '</div><div class="label">Total Brokers</div></div>' +
                '<div class="summary-card"><div class="number">' + totalTopics + '</div><div class="label">Total Topics</div></div>';
        }

        function renderClusters(data) {
            var clustersDiv = document.getElementById('clusters');
            clustersDiv.innerHTML = data.clusters.map(function(cluster) { return renderCluster(cluster); }).join('');
        }

        function renderCluster(cluster) {
            var statusClass = 'status-' + cluster.status;
            var componentsHTML = '';

            if (cluster.kafka && cluster.kafka.available) {
                componentsHTML += renderComponent('Kafka', 'K', {
                    'Brokers': cluster.kafka.broker_count,
                    'Controller': cluster.kafka.controller_type,
                    'Topics': cluster.kafka.topic_count,
                    'Partitions': cluster.kafka.total_partitions,
                    'Security': (cluster.kafka.security_config && cluster.kafka.security_config.authentication_method) || 'None'
                });
            }

            if (cluster.schema_registry && cluster.schema_registry.available) {
                componentsHTML += renderComponent('Schema Registry', 'SR', {
                    'Version': cluster.schema_registry.version,
                    'Mode': cluster.schema_registry.mode,
                    'Schemas': cluster.schema_registry.total_schemas
                });
            }

            if (cluster.kafka_connect && cluster.kafka_connect.available) {
                componentsHTML += renderComponent('Kafka Connect', 'KC', {
                    'Version': cluster.kafka_connect.version,
                    'Connectors': cluster.kafka_connect.total_connectors,
                    'Source': cluster.kafka_connect.source_connectors,
                    'Sink': cluster.kafka_connect.sink_connectors
                });
            }

            if (cluster.ksqldb && cluster.ksqldb.available) {
                componentsHTML += renderComponent('ksqlDB', 'KS', {
                    'Version': cluster.ksqldb.version,
                    'Queries': cluster.ksqldb.queries,
                    'Streams': cluster.ksqldb.streams,
                    'Tables': cluster.ksqldb.tables
                });
            }

            if (cluster.rest_proxy && cluster.rest_proxy.available) {
                componentsHTML += renderComponent('REST Proxy', 'RP', {
                    'Version': cluster.rest_proxy.version,
                    'Cluster ID': cluster.rest_proxy.cluster_id,
                    'Controller Mode': cluster.rest_proxy.controller_mode
                });
            }

            if (cluster.control_center && cluster.control_center.available) {
                componentsHTML += renderComponent('Control Center', 'C3', {
                    'Version': cluster.control_center.version,
                    'URL': cluster.control_center.url,
                    'Monitored Clusters': cluster.control_center.monitored_clusters
                });
            }

            if (cluster.prometheus && cluster.prometheus.available) {
                componentsHTML += renderComponent('Prometheus', 'PM', {
                    'Version': cluster.prometheus.version,
                    'URL': cluster.prometheus.url,
                    'Targets Up': cluster.prometheus.targets_up,
                    'Targets Down': cluster.prometheus.targets_down
                });
            }

            if (cluster.alertmanager && cluster.alertmanager.available) {
                componentsHTML += renderComponent('Alertmanager', 'AM', {
                    'Version': cluster.alertmanager.version,
                    'Cluster Size': cluster.alertmanager.cluster_size,
                    'Active Alerts': cluster.alertmanager.active_alerts
                });
            }

            return '<div class="cluster"><div class="cluster-header">' +
                '<div class="cluster-name">' + cluster.name + '</div>' +
                '<div class="status-badge ' + statusClass + '">' + cluster.status + '</div>' +
                '</div><div class="components">' + componentsHTML + '</div></div>';
        }

        function renderComponent(title, icon, details) {
            var detailsHTML = Object.keys(details).map(function(key) {
                return '<div class="component-detail">' +
                    '<span class="detail-label">' + key + ':</span>' +
                    '<span class="detail-value">' + details[key] + '</span>' +
                    '</div>';
            }).join('');

            return '<div class="component"><div class="component-header">' +
                '<div class="component-icon">' + icon + '</div>' +
                '<div class="component-title">' + title + '</div>' +
                '</div>' + detailsHTML + '</div>';
        }

        function switchTab(tabName) {
            document.querySelectorAll('.tab-content').forEach(function(content) {
                content.classList.remove('active');
            });
            document.querySelectorAll('.tab').forEach(function(tab) {
                tab.classList.remove('active');
            });
            document.getElementById(tabName + '-content').classList.add('active');
            event.target.classList.add('active');
        }
    </script>
</body>
</html>`
}
