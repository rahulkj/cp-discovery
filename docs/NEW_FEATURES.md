# New Features Summary

## Command-Line Arguments Enhancement

### Overview
Added flexible command-line arguments that allow users to override configuration file settings without modifying the config file.

### New Flags

#### 1. `-output` flag
**Purpose:** Specify output file path from command line  
**Type:** string  
**Default:** Uses value from config file  

**Example:**
```bash
./bin/cp-discovery -output /tmp/my-report.json
```

**Use Cases:**
- Dynamic file naming with timestamps
- Different output locations per run
- CI/CD pipelines with specific paths
- Environment-based output directories

#### 2. `-format` flag
**Purpose:** Override output format  
**Type:** string  
**Values:** `json` or `yaml`  
**Default:** Uses value from config file  

**Example:**
```bash
./bin/cp-discovery -format yaml -output report.yaml
```

**Use Cases:**
- Quick format switching without config changes
- Different formats for different tools
- Human-readable YAML for debugging
- Machine-parseable JSON for automation

#### 3. `-detailed` flag
**Purpose:** Enable detailed discovery mode  
**Type:** boolean  
**Default:** false (uses config file setting if not specified)  

**Example:**
```bash
./bin/cp-discovery -detailed
```

**Use Cases:**
- Comprehensive cluster analysis
- Troubleshooting and diagnostics
- Audit and compliance reporting
- Performance analysis

### Combined Usage

All flags can be combined for maximum flexibility:

```bash
./bin/cp-discovery \
  -config configs/config-production.yaml \
  -output /var/reports/kafka-$(date +%Y%m%d-%H%M%S).json \
  -format json \
  -detailed
```

### Priority Order

When both config file and command-line flags are specified:

1. **Command-line flags** (highest priority)
2. **Configuration file** settings
3. **Built-in defaults**

**Example:**
```yaml
# config.yaml
output:
  file: "default.json"
  format: "json"
  detailed: false
```

```bash
# Command overrides config
./bin/cp-discovery -output custom.json -detailed

# Results in:
# - Output: custom.json (from CLI)
# - Format: json (from config, not overridden)
# - Detailed: true (from CLI)
```

---

## Network Throughput Display

### Overview
Enhanced console output to display network throughput metrics for Kafka clusters.

### Metrics Displayed

#### From Kafka Discovery
```
  Kafka:
    Network Throughput:
      Bytes In: 125.50 MB/s
      Bytes Out: 256.75 MB/s
      Messages In: 50000.00 msg/s
```

**Fields:**
- **Bytes In/s**: Incoming network traffic in MB/s
- **Bytes Out/s**: Outgoing network traffic in MB/s
- **Messages In/s**: Incoming message rate

**Data Source:**
- `ClusterMetrics.BytesInPerSec`
- `ClusterMetrics.BytesOutPerSec`
- `ClusterMetrics.MessagesInPerSec`

#### From Prometheus (Enhanced)
```
  Prometheus:
    Cluster Metrics:
      Throughput: 150.25 MB/s in, 300.50 MB/s out
      Messages: 75000.00 msg/s in
```

**Data Source:**
- Prometheus metrics queries
- Real-time cluster monitoring

### Display Logic

Network throughput is displayed when:
- Any throughput metric > 0
- Data is available from cluster metrics
- Prometheus metrics are enabled (in detailed mode)

---

## Storage Details Display

### Overview
Enhanced console output to display storage information for Kafka clusters.

### Metrics Displayed

#### Cluster-Level Storage
```
  Kafka:
    Storage:
      Total Disk Usage: 1250.50 GB
```

**Data Source:**
- `ClusterMetrics.TotalDiskUsageBytes`
- Aggregated from all brokers

#### Broker-Level Storage (Future)
```
  Kafka:
    Storage:
      Total Disk Usage: 1250.50 GB (from brokers)
      Broker 1: 250.10 GB
      Broker 2: 250.10 GB
      Broker 3: 250.10 GB
```

**Model Update:**
- Added `DiskUsageBytes` field to `BrokerInfo`
- Enables per-broker storage tracking

### Display Logic

Storage information is displayed when:
- `TotalDiskUsageBytes > 0`
- Individual broker storage data is available
- Detailed mode is enabled

### Storage Calculation

Two methods for storage calculation:

1. **Cluster Metrics** (primary)
   - Direct from `ClusterMetrics.TotalDiskUsageBytes`
   - Most accurate for overall storage

2. **Broker Aggregation** (fallback)
   - Sum of `BrokerInfo.DiskUsageBytes` across all brokers
   - Used when cluster-level metrics unavailable

---

## Health Metrics Display

### Overview
Enhanced display of cluster health indicators.

### Metrics Displayed

```
  Kafka:
    Health:
      Under-Replicated Partitions: 5
```

**Data Source:**
- `ClusterMetrics.UnderReplicatedPartitions`

**Display Conditions:**
- Only shown when > 0 (indicates issues)
- Helps identify replication problems quickly

---

## Enhanced Prometheus Metrics Display

### Overview
Comprehensive Prometheus metrics display showing detailed cluster health.

### Full Metrics Output

```
  Prometheus:
    Version: 2.45.0
    URL: http://prometheus:9090
    Targets Up: 15
    Targets Down: 0

    Cluster Metrics:
      Throughput: 150.25 MB/s in, 300.50 MB/s out
      Messages: 75000.00 msg/s in
      Active Controllers: 3
      Brokers: 5 online / 5 total
      Partitions: 1000 total (5 under-replicated) (0 offline)
      Consumers: 25 groups, lag: 12500
      JVM: 65.5% heap, 45.2% CPU (avg across brokers)
```

### Metric Categories

1. **Throughput**
   - Bytes in/out per second
   - Messages per second

2. **Broker Health**
   - Total brokers
   - Online brokers
   - Active controllers

3. **Partition Health**
   - Total partitions
   - Under-replicated partitions
   - Offline partitions

4. **Consumer Health**
   - Consumer group count
   - Total lag across all groups

5. **JVM Metrics**
   - Average heap usage %
   - Average CPU usage %

---

## JSON/YAML Report Structure

### Network Throughput in Reports

```json
{
  "kafka": {
    "cluster_metrics": {
      "bytes_in_per_sec": 131621888.5,
      "bytes_out_per_sec": 268435456.0,
      "messages_in_per_sec": 50000.0
    }
  }
}
```

### Storage in Reports

```json
{
  "kafka": {
    "cluster_metrics": {
      "total_disk_usage_bytes": 1342177280000
    },
    "brokers": [
      {
        "id": 1,
        "host": "broker-1",
        "port": 9092,
        "disk_usage_bytes": 268435456000
      }
    ]
  }
}
```

---

## Use Cases

### 1. Production Monitoring
```bash
# Daily detailed report with custom naming
./bin/cp-discovery \
  -config configs/config-production.yaml \
  -output /var/reports/prod-$(date +%Y%m%d).json \
  -detailed
```

### 2. Performance Troubleshooting
```bash
# High-detail snapshot for debugging
./bin/cp-discovery \
  -detailed \
  -output /tmp/perf-analysis-$(date +%H%M%S).json
```

### 3. CI/CD Integration
```bash
# Quick health check in pipeline
./bin/cp-discovery \
  -config configs/ci-config.yaml \
  -output $CI_WORKSPACE/kafka-status.json \
  -format json
```

### 4. Multi-Format Reporting
```bash
# JSON for automation
./bin/cp-discovery -output report.json -format json

# YAML for human review
./bin/cp-discovery -output report.yaml -format yaml
```

### 5. Capacity Planning
```bash
# Detailed storage and throughput analysis
./bin/cp-discovery \
  -detailed \
  -output capacity-analysis.json

# Extract storage metrics
jq '.clusters[].kafka.cluster_metrics.total_disk_usage_bytes' capacity-analysis.json
```

---

## Implementation Details

### Code Changes

**File:** `cm./cp-discovery/main.go`

1. **Added command-line flags:**
   ```go
   outputFile := flag.String("output", "", "Output file path")
   outputFormat := flag.String("format", "", "Output format: json or yaml")
   detailed := flag.Bool("detailed", false, "Enable detailed discovery")
   ```

2. **Override logic:**
   ```go
   if *outputFile != "" {
       cfg.Output.File = *outputFile
   }
   if *outputFormat != "" {
       cfg.Output.Format = *outputFormat
   }
   if *detailed {
       cfg.Output.Detailed = true
   }
   ```

3. **Enhanced console output:**
   - Network throughput section
   - Storage details section
   - Health metrics section
   - Broker-level details

**File:** `internal/model/models.go`

1. **Enhanced BrokerInfo:**
   ```go
   type BrokerInfo struct {
       ID              int    `json:"id"`
       Host            string `json:"host"`
       Port            int    `json:"port"`
       Rack            string `json:"rack,omitempty"`
       DiskUsageBytes  int64  `json:"disk_usage_bytes,omitempty"`
   }
   ```

---

## Benefits

### For Operators
- ✅ **Flexible Execution**: Override settings without editing configs
- ✅ **Quick Troubleshooting**: Enable detailed mode on-demand
- ✅ **Better Visibility**: See network and storage metrics at a glance
- ✅ **Custom Workflows**: Integrate with scripts and automation

### For Automation
- ✅ **Dynamic Paths**: Generate timestamped reports
- ✅ **CI/CD Ready**: Easy integration with pipelines
- ✅ **Format Flexibility**: Choose output format per use case
- ✅ **Scriptable**: All options available via command line

### For Monitoring
- ✅ **Network Insights**: Track throughput trends
- ✅ **Capacity Planning**: Monitor storage usage
- ✅ **Health Indicators**: Identify issues quickly
- ✅ **Comprehensive Data**: All metrics in one place

---

## Migration from Previous Version

### No Breaking Changes
All existing configurations and usage patterns continue to work:

```bash
# Still works exactly as before
./bin/cp-discovery -config configs/config.yaml
```

### Enhanced Capabilities
New features are additive and optional:

```bash
# Use new features as needed
./bin/cp-discovery -config configs/config.yaml -output custom.json
```

### Backward Compatibility
- ✅ All existing configs work unchanged
- ✅ Default behavior preserved
- ✅ Optional enhancements only
- ✅ Graceful handling of missing data
