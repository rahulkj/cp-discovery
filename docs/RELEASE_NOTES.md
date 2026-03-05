# Release Notes - Version 2.0.0

## 🎉 Major Feature Release

### Released: March 4, 2026

---

## 🚀 What's New

### 1. Command-Line Arguments

**Flexible execution without modifying config files!**

```bash
# Override output location
./bin/cp-discovery -output /tmp/report.json

# Change format on the fly
./bin/cp-discovery -format yaml

# Enable detailed mode
./bin/cp-discovery -detailed

# Combine everything
./bin/cp-discovery \
  -config configs/config-production.yaml \
  -output /var/reports/kafka-$(date +%Y%m%d).json \
  -format json \
  -detailed
```

**New Flags:**
- ✅ `-output <path>` - Override output file location
- ✅ `-format <json|yaml>` - Override output format
- ✅ `-detailed` - Enable comprehensive discovery mode
- ✅ `-config <path>` - Configuration file (existing, now more useful)

**Benefits:**
- Dynamic file naming with timestamps
- CI/CD integration made easy
- No need to modify config files for one-off runs
- Perfect for automation and scripting

---

### 2. Network Throughput Display

**See your cluster's network performance at a glance!**

```
  Kafka:
    Network Throughput:
      Bytes In: 125.50 MB/s
      Bytes Out: 256.75 MB/s
      Messages In: 50000.00 msg/s
```

**What's tracked:**
- Incoming bytes per second (MB/s)
- Outgoing bytes per second (MB/s)
- Incoming messages per second

**Data sources:**
- Kafka ClusterMetrics
- Prometheus (when enabled)

**Use cases:**
- Performance monitoring
- Capacity planning
- Trend analysis
- Bottleneck identification

---

### 3. Storage Details Display

**Monitor disk usage across your cluster!**

```
  Kafka:
    Storage:
      Total Disk Usage: 1250.50 GB
```

**Features:**
- Cluster-level total disk usage
- Per-broker storage tracking (in JSON/YAML)
- Formatted in GB for easy reading
- Supports capacity planning

**Model Enhancement:**
```go
type BrokerInfo struct {
    ID              int    `json:"id"`
    Host            string `json:"host"`
    Port            int    `json:"port"`
    Rack            string `json:"rack,omitempty"`
    DiskUsageBytes  int64  `json:"disk_usage_bytes,omitempty"`  // NEW!
}
```

---

### 4. Health Metrics Display

**Quickly identify cluster issues!**

```
  Kafka:
    Health:
      Under-Replicated Partitions: 5
```

**Shown when:**
- Under-replicated partitions > 0
- Other health issues detected
- Helps prioritize troubleshooting

---

## 📊 Enhanced Prometheus Metrics

**Comprehensive cluster insights from Prometheus!**

```
  Prometheus:
    Cluster Metrics:
      Throughput: 150.25 MB/s in, 300.50 MB/s out
      Messages: 75000.00 msg/s in
      Active Controllers: 3
      Brokers: 5 online / 5 total
      Partitions: 1000 total (5 under-replicated) (0 offline)
      Consumers: 25 groups, lag: 12500
      JVM: 65.5% heap, 45.2% CPU (avg across brokers)
```

**Metrics Categories:**
1. **Network Throughput** - Bytes and messages per second
2. **Broker Health** - Total, online, controllers
3. **Partition Health** - Total, under-replicated, offline
4. **Consumer Health** - Groups and lag
5. **JVM Metrics** - Heap and CPU usage

---

## 🏗️ Project Restructuring

**Clean, maintainable Go project structure!**

### New Organization
```
cp-discovery/
├── cm./cp-discovery/    # Main application
├── internal/
│   ├── model/                  # Data structures
│   ├── config/                 # Configuration
│   ├── http/                   # Authentication
│   └── discovery/              # Component discovery
├── configs/                    # All YAML configs
└── bin/                        # Compiled binary
```

### Benefits
- ✅ Standard Go layout
- ✅ Clear separation of concerns
- ✅ Easier to navigate
- ✅ Better maintainability
- ✅ Follows community best practices

---

## 📚 New Documentation

### USAGE_EXAMPLES.md
- Comprehensive usage guide
- Practical examples for every flag
- CI/CD integration patterns
- Best practices

### NEW_FEATURES.md
- Detailed feature documentation
- Implementation details
- Use cases and benefits
- Migration guide

### PROJECT_STRUCTURE.md
- Complete project organization
- Package descriptions
- Build instructions
- Docker usage

### CHANGELOG.md
- Version history
- Breaking changes (none!)
- Upgrade guide
- Future roadmap

---

## 🔧 Use Cases

### Production Monitoring
```bash
# Daily detailed report
./bin/cp-discovery \
  -config configs/config-production.yaml \
  -output /var/reports/prod-$(date +%Y%m%d).json \
  -detailed
```

### Performance Troubleshooting
```bash
# Quick detailed snapshot
./bin/cp-discovery \
  -detailed \
  -output /tmp/perf-analysis.json
```

### CI/CD Integration
```bash
# Pipeline health check
./bin/cp-discovery \
  -config configs/ci-config.yaml \
  -output $CI_WORKSPACE/kafka-status.json \
  -format json
```

### Capacity Planning
```bash
# Extract storage metrics
./bin/cp-discovery -detailed -output capacity.json
jq '.clusters[].kafka.cluster_metrics.total_disk_usage_bytes' capacity.json
```

---

## ⚡ Quick Start

### Installation
```bash
# Build
go build -o bin/cp-discovery ./cm./cp-discovery

# Or use make
make build
```

### Basic Usage
```bash
# Default run
./bin/cp-discovery

# Custom output
./bin/cp-discovery -output my-report.json

# Detailed mode
./bin/cp-discovery -detailed
```

### Help
```bash
./bin/cp-discovery --help
```

---

## 🔄 Upgrade from 1.x

### Zero Breaking Changes!

All existing configurations and scripts work unchanged:

```bash
# This still works exactly as before
./bin/cp-discovery -config configs/config.yaml
```

### New Capabilities Are Additive

Use new features when you need them:

```bash
# Add new flags as needed
./bin/cp-discovery \
  -config configs/config.yaml \
  -output custom.json \
  -detailed
```

### What Continues to Work
- ✅ All existing config files
- ✅ All existing workflows
- ✅ All existing integrations
- ✅ All existing outputs

### What's Enhanced
- ✅ More flexible execution
- ✅ Better console output
- ✅ More detailed metrics
- ✅ Cleaner code structure

---

## 📈 What This Means for You

### For Operators
- **Faster troubleshooting** with network and storage visibility
- **Flexible execution** without config file changes
- **Better insights** with enhanced metrics display
- **Easier automation** with CLI flags

### For Automation
- **Dynamic file naming** for timestamped reports
- **CI/CD ready** with format and output control
- **Scriptable** with all options on command line
- **Integration friendly** with flexible output

### For Monitoring
- **Capacity planning** with storage metrics
- **Performance tracking** with throughput data
- **Health monitoring** with issue indicators
- **Trend analysis** with detailed metrics

---

## 🎯 Next Steps

1. **Try the new flags:**
   ```bash
   ./bin/cp-discovery -output test.json -detailed
   ```

2. **Review the new output sections** - Network, storage, health

3. **Check out USAGE_EXAMPLES.md** for integration patterns

4. **Update your automation scripts** to use new features

5. **Provide feedback** on what works and what to improve

---

## 🐛 Bug Fixes

- Fixed controller count to show total controllers (not just active)
- Improved type safety with proper model prefixes
- Enhanced error handling in discovery functions

---

## 🙏 Thank You

This release represents significant enhancements to make cp-discovery more powerful and flexible for production use.

**Feedback welcome!**
- GitHub Issues: https://github.com/rahulkj/cp-discovery/issues
- Feature Requests: Same URL

**Happy Discovering! 🚀**

---

## 📦 Download

```bash
# Clone repository
git clone https://github.com/rahulkj/cp-discovery.git

# Build
cd cp-discovery
go build -o bin/cp-discovery ./cm./cp-discovery

# Run
./bin/cp-discovery --help
```

---

**Version:** 2.0.0  
**Release Date:** March 4, 2026  
**License:** Apache 2.0
