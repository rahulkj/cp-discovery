package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
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
	flag.Parse()

	cfg, err := configPkg.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
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

	fmt.Printf("\n[cyan][2/3][reset] Generating report...\n")
	if err := outputReport(report, cfg.Output); err != nil {
		log.Fatalf("Failed to output report: %v", err)
	}

	fmt.Printf("[cyan][3/3][reset] Displaying summary...\n\n")
	printSummary(report)

	fmt.Printf("\n✅ [green]Discovery completed successfully![reset]\n")
}

func discoverClusters(cfg *model.Config) *model.DiscoveryReport {
	report := &model.DiscoveryReport{
		Timestamp:     time.Now().Format(time.RFC3339),
		Clusters:      make([]model.ClusterReport, 0),
		TotalClusters: len(cfg.Clusters),
	}

	detailed := cfg.Output.Detailed

	// Calculate total steps for progress bar
	// Each cluster has: Kafka (1) + up to 7 optional components
	totalSteps := 0
	for _, cluster := range cfg.Clusters {
		steps := 1 // Kafka is always discovered
		if configPkg.ShouldDiscoverSchemaRegistry(&cluster) {
			steps++
		}
		if configPkg.ShouldDiscoverKafkaConnect(&cluster) {
			steps++
		}
		if configPkg.ShouldDiscoverKsqlDB(&cluster) {
			steps++
		}
		if configPkg.ShouldDiscoverRestProxy(&cluster) {
			steps++
		}
		if configPkg.ShouldDiscoverControlCenter(&cluster) {
			steps++
		}
		if configPkg.ShouldDiscoverPrometheus(&cluster) {
			steps++
		}
		if configPkg.ShouldDiscoverAlertmanager(&cluster) {
			steps++
		}
		totalSteps += steps
	}

	// Create progress bar
	bar := progressbar.NewOptions(totalSteps,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Discovering clusters..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	fmt.Printf("\n🔍 Starting discovery for %d cluster(s)...\n\n", len(cfg.Clusters))

	var wg sync.WaitGroup
	clusterReports := make([]model.ClusterReport, len(cfg.Clusters))

	for i, cluster := range cfg.Clusters {
		wg.Add(1)
		go func(idx int, clusterConfig model.ClusterConfig) {
			defer wg.Done()
			clusterReports[idx] = discoverCluster(clusterConfig, detailed, bar)
		}(i, cluster)
	}

	wg.Wait()

	fmt.Println()
	report.Clusters = clusterReports
	return report
}

func discoverCluster(config model.ClusterConfig, detailed bool, bar *progressbar.ProgressBar) model.ClusterReport {
	report := model.ClusterReport{
		Name:   config.Name,
		Status: "healthy",
		Errors: make([]string, 0),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Discover Kafka (always required)
	wg.Add(1)
	go func() {
		defer wg.Done()
		bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → Kafka", config.Name))
		kafkaReport, err := discovery.DiscoverKafka(config.Kafka, detailed)
		bar.Add(1)
		if err != nil {
			mu.Lock()
			report.Errors = append(report.Errors, fmt.Sprintf("Kafka: %v", err))
			mu.Unlock()
		} else {
			mu.Lock()
			report.Kafka = kafkaReport
			mu.Unlock()
		}
	}()

	// Discover Schema Registry (only if enabled)
	if configPkg.ShouldDiscoverSchemaRegistry(&config) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → Schema Registry", config.Name))
			srReport, err := discovery.DiscoverSchemaRegistry(config.SchemaRegistry, detailed)
			bar.Add(1)
			mu.Lock()
			defer mu.Unlock()
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
			bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → Kafka Connect", config.Name))
			connectReport, err := discovery.DiscoverKafkaConnect(config.KafkaConnect, detailed)
			bar.Add(1)
			mu.Lock()
			defer mu.Unlock()
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
			bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → ksqlDB", config.Name))
			ksqlReport, err := discovery.DiscoverKsqlDB(config.KsqlDB, detailed)
			bar.Add(1)
			mu.Lock()
			defer mu.Unlock()
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
			bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → REST Proxy", config.Name))
			restReport, err := discovery.DiscoverRestProxy(config.RestProxy, detailed)
			bar.Add(1)
			mu.Lock()
			defer mu.Unlock()
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
			bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → Control Center", config.Name))
			ccReport, err := discovery.DiscoverControlCenter(config.ControlCenter, detailed)
			bar.Add(1)
			mu.Lock()
			defer mu.Unlock()
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
			bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → Prometheus", config.Name))
			promReport, err := discovery.DiscoverPrometheus(config.Prometheus, detailed)
			bar.Add(1)
			mu.Lock()
			defer mu.Unlock()
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
			bar.Describe(fmt.Sprintf("[cyan][1/3][reset] Discovering [yellow]%s[reset] → Alertmanager", config.Name))
			amReport, err := discovery.DiscoverAlertmanager(config.Alertmanager, detailed)
			bar.Add(1)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("Alertmanager: %v", err))
			} else {
				report.Alertmanager = amReport
			}
		}()
	}

	wg.Wait()

	// Enrich topics with associated schemas (if Schema Registry is available)
	if report.SchemaRegistry.Available && len(report.SchemaRegistry.Subjects) > 0 && len(report.Kafka.Topics) > 0 {
		enrichTopicsWithSchemas(&report)
	}

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

func printNodeCountSummary(report *model.DiscoveryReport) {
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("NODE COUNT SUMMARY (Across All Clusters)")
	fmt.Println(strings.Repeat("-", 80))

	// Aggregate counts
	totalBrokers := 0
	totalControllers := 0
	totalSchemaRegistryNodes := 0
	totalConnectWorkers := 0
	totalKsqlDBNodes := 0
	totalRestProxyInstances := 0
	totalControlCenterInstances := 0

	for _, cluster := range report.Clusters {
		if cluster.Kafka.Available {
			totalBrokers += cluster.Kafka.BrokerCount
			totalControllers += cluster.Kafka.ControllerCount
		}
		if cluster.SchemaRegistry.Available {
			totalSchemaRegistryNodes += cluster.SchemaRegistry.NodeCount
		}
		if cluster.KafkaConnect.Available {
			totalConnectWorkers += cluster.KafkaConnect.WorkerCount
		}
		if cluster.KsqlDB.Available {
			totalKsqlDBNodes += cluster.KsqlDB.NodeCount
		}
		if cluster.RestProxy.Available {
			totalRestProxyInstances++ // REST Proxy doesn't track node count, so count instances
		}
		if cluster.ControlCenter.Available {
			totalControlCenterInstances++ // Control Center doesn't track node count, so count instances
		}
	}

	fmt.Printf("  Kafka Brokers:           %d\n", totalBrokers)
	if totalControllers > 0 {
		fmt.Printf("  KRaft Controllers:       %d\n", totalControllers)
	}
	if totalSchemaRegistryNodes > 0 {
		fmt.Printf("  Schema Registry Nodes:   %d\n", totalSchemaRegistryNodes)
	}
	if totalConnectWorkers > 0 {
		fmt.Printf("  Kafka Connect Workers:   %d\n", totalConnectWorkers)
	}
	if totalKsqlDBNodes > 0 {
		fmt.Printf("  ksqlDB Nodes:            %d\n", totalKsqlDBNodes)
	}
	if totalRestProxyInstances > 0 {
		fmt.Printf("  REST Proxy Instances:    %d\n", totalRestProxyInstances)
	}
	if totalControlCenterInstances > 0 {
		fmt.Printf("  Control Center Instances: %d\n", totalControlCenterInstances)
	}
	fmt.Println()
}

func printSummary(report *model.DiscoveryReport) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("CONFLUENT PLATFORM DISCOVERY SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Timestamp: %s\n", report.Timestamp)
	fmt.Printf("Total Clusters: %d\n\n", report.TotalClusters)

	// Print aggregate node counts across all clusters
	printNodeCountSummary(report)

	for _, cluster := range report.Clusters {
		fmt.Printf("Cluster: %s [%s]\n", cluster.Name, cluster.Status)
		fmt.Println(strings.Repeat("-", 80))

		// Kafka
		if cluster.Kafka.Available {
			fmt.Printf("  Kafka:\n")
			fmt.Printf("    Brokers: %d\n", cluster.Kafka.BrokerCount)
			fmt.Printf("    Controller: %s", cluster.Kafka.ControllerType)
			if cluster.Kafka.ControllerCount > 0 {
				fmt.Printf(" (Controllers: %d)", cluster.Kafka.ControllerCount)
			}
			fmt.Println()
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

			// Storage - always show if we have topic data
			totalBrokerStorage := int64(0)
			if len(cluster.Kafka.Brokers) > 0 {
				for _, broker := range cluster.Kafka.Brokers {
					if broker.DiskUsageBytes > 0 {
						totalBrokerStorage += broker.DiskUsageBytes
					}
				}
			}

			// Display storage if we have data from either metrics or brokers
			if metrics.TotalDiskUsageBytes > 0 || totalBrokerStorage > 0 {
				fmt.Printf("    Storage:\n")
				if metrics.TotalDiskUsageBytes > 0 {
					fmt.Printf("      Total Cluster Storage: %.2f GB\n", float64(metrics.TotalDiskUsageBytes)/1024/1024/1024)
				} else if totalBrokerStorage > 0 {
					fmt.Printf("      Total Cluster Storage: %.2f GB (from brokers)\n", float64(totalBrokerStorage)/1024/1024/1024)
				}

				// Show topic count with storage
				topicsWithStorage := 0
				for _, topic := range cluster.Kafka.Topics {
					if topic.SizeBytes > 0 {
						topicsWithStorage++
					}
				}
				if topicsWithStorage > 0 {
					fmt.Printf("      Topics with Storage Data: %d\n", topicsWithStorage)
				}
			}

			// Cluster health
			if metrics.UnderReplicatedPartitions > 0 {
				fmt.Printf("    Health:\n")
				fmt.Printf("      Under-Replicated Partitions: %d\n", metrics.UnderReplicatedPartitions)
			}

			// Topic details (if detailed mode)
			if len(cluster.Kafka.Topics) > 0 {
				fmt.Printf("\n    Topics:\n")
				for _, topic := range cluster.Kafka.Topics {
					if !topic.IsInternal {
						fmt.Printf("      • %s\n", topic.Name)
						fmt.Printf("        Partitions: %d | Replication: %d", topic.Partitions, topic.ReplicationFactor)
						if topic.RetentionMs > 0 {
							fmt.Printf(" | Retention: %s", formatRetention(topic.RetentionMs))
						}
						if topic.SizeBytes > 0 {
							fmt.Printf(" | Size: %.2f MB", float64(topic.SizeBytes)/1024/1024)
						}
						fmt.Println()
						if len(topic.AssociatedSchemas) > 0 {
							fmt.Printf("        Schemas: %s\n", strings.Join(topic.AssociatedSchemas, ", "))
						}
					}
				}
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

// enrichTopicsWithSchemas maps Schema Registry subjects to Kafka topics
func enrichTopicsWithSchemas(report *model.ClusterReport) {
	// Create a map of topic names for quick lookup
	topicMap := make(map[string]*model.TopicInfo)
	for i := range report.Kafka.Topics {
		topicMap[report.Kafka.Topics[i].Name] = &report.Kafka.Topics[i]
	}

	// Match subjects to topics
	// Schema naming conventions:
	// - <topic-name>-key
	// - <topic-name>-value
	// - <topic-name>
	for _, subject := range report.SchemaRegistry.Subjects {
		// Try to extract topic name from subject
		topicName := extractTopicFromSubject(subject)

		if topic, exists := topicMap[topicName]; exists {
			// Add schema to topic's associated schemas
			topic.AssociatedSchemas = append(topic.AssociatedSchemas, subject)
		}
	}
}

// extractTopicFromSubject extracts the topic name from a schema subject
func extractTopicFromSubject(subject string) string {
	// Common patterns:
	// - topic-name-key -> topic-name
	// - topic-name-value -> topic-name
	// - topic-name -> topic-name

	subject = strings.TrimSuffix(subject, "-key")
	subject = strings.TrimSuffix(subject, "-value")

	return subject
}
