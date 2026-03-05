# Usage Examples

This document provides practical examples of using the cp-discovery tool with various command-line options.

## Basic Usage

### Default Configuration
Run with default config file (`configs/config.yaml`):
```bash
./bin/cp-discovery
```

### Custom Configuration File
Specify a different configuration file:
```bash
./bin/cp-discovery -config configs/config-production.yaml
```

## Output Control

### Save to Specific File
Override the output file from command line:
```bash
./bin/cp-discovery -output /tmp/my-discovery.json
```

### Change Output Format
Override output format (json or yaml):
```bash
# JSON output
./bin/cp-discovery -format json -output report.json

# YAML output
./bin/cp-discovery -format yaml -output report.yaml
```

### Detailed Mode
Enable detailed discovery to get comprehensive information:
```bash
./bin/cp-discovery -detailed -output detailed-report.json
```

## Combined Options

### Full Custom Discovery
Combine all options for complete control:
```bash
./bin/cp-discovery \
  -config configs/config-production.yaml \
  -output /var/reports/discovery-$(date +%Y%m%d).json \
  -format json \
  -detailed
```

### Quick JSON Export
Get a quick JSON report with default config:
```bash
./bin/cp-discovery -output report.json -format json
```

### Detailed YAML Report
Get comprehensive YAML report:
```bash
./bin/cp-discovery -detailed -format yaml -output detailed.yaml
```

## Practical Scenarios

### Daily Production Monitoring
```bash
#!/bin/bash
DATE=$(date +%Y-%m-%d)
./bin/cp-discovery \
  -config configs/config-production.yaml \
  -output /var/logs/kafka-discovery-${DATE}.json \
  -format json \
  -detailed
```

### Quick Dev Check
```bash
./bin/cp-discovery -config configs/example-local.yaml -output /tmp/dev-check.json
```

### Multi-Environment Discovery
```bash
# Production
./bin/cp-discovery -config configs/config-production.yaml -output prod-report.json

# Staging
./bin/cp-discovery -config configs/config-staging.yaml -output staging-report.json

# Development
./bin/cp-discovery -config configs/example-local.yaml -output dev-report.json
```

### CI/CD Integration
```bash
#!/bin/bash
# Run discovery and check for failures
./bin/cp-discovery \
  -config configs/config.yaml \
  -output discovery-report.json \
  -detailed

# Process the report (example with jq)
UNHEALTHY=$(jq '[.clusters[] | select(.status != "healthy")] | length' discovery-report.json)

if [ "$UNHEALTHY" -gt 0 ]; then
  echo "Warning: $UNHEALTHY unhealthy clusters detected"
  exit 1
fi
```

### Minimal vs Detailed Comparison
```bash
# Minimal (fast, basic info)
./bin/cp-discovery -config configs/config.yaml -output minimal.json

# Detailed (slower, comprehensive)
./bin/cp-discovery -config configs/config.yaml -output detailed.json -detailed

# Compare file sizes
ls -lh minimal.json detailed.json
```

## Output Examples

### Console Output Features

The tool displays comprehensive information including:

**Network Throughput:**
```
  Kafka:
    Network Throughput:
      Bytes In: 125.50 MB/s
      Bytes Out: 256.75 MB/s
      Messages In: 50000.00 msg/s
```

**Storage Information:**
```
  Kafka:
    Storage:
      Total Disk Usage: 1250.50 GB
```

**Health Metrics:**
```
  Kafka:
    Health:
      Under-Replicated Partitions: 5
```

**Prometheus Metrics (if enabled):**
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

## Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-config` | string | `configs/config.yaml` | Path to configuration file |
| `-output` | string | (from config) | Output file path (overrides config) |
| `-format` | string | (from config) | Output format: json or yaml (overrides config) |
| `-detailed` | bool | false | Enable detailed discovery (overrides config) |

## Configuration File vs Command-Line

Command-line flags take precedence over configuration file settings:

**Priority Order (highest to lowest):**
1. Command-line flags (`-output`, `-format`, `-detailed`)
2. Configuration file settings
3. Built-in defaults

**Example:**
```yaml
# config.yaml
output:
  file: "default-report.json"
  format: "json"
  detailed: false
```

```bash
# This will use detailed mode and save to custom-output.json
./bin/cp-discovery -output custom-output.json -detailed

# Result:
# - File: custom-output.json (from CLI)
# - Format: json (from config)
# - Detailed: true (from CLI)
```

## Best Practices

### 1. Use Detailed Mode for Troubleshooting
```bash
./bin/cp-discovery -detailed -output troubleshooting-$(date +%Y%m%d-%H%M%S).json
```

### 2. Regular Monitoring Schedule
```bash
# Crontab entry: run every 6 hours
0 */6 * * * /path/to/bin/cp-discovery -config /etc/kafka/discovery-config.yaml -output /var/log/kafka/discovery-$(date +\%Y\%m\%d-\%H\%M).json
```

### 3. Environment-Specific Configs
```bash
# Use environment variable for config selection
CONFIG_ENV=${KAFKA_ENV:-production}
./bin/cp-discovery -config configs/config-${CONFIG_ENV}.yaml
```

### 4. Retention Management
```bash
# Keep only last 7 days of reports
find /var/log/kafka/discovery-* -mtime +7 -delete
```

## Troubleshooting

### Permission Issues
```bash
# Make binary executable
chmod +x bin/cp-discovery

# Run with proper permissions for config file
chmod 600 configs/config-production.yaml  # Contains secrets
```

### Output File Issues
```bash
# Ensure output directory exists
mkdir -p /var/reports/kafka

# Check disk space before detailed mode
df -h /var/reports/kafka
```

### Configuration Validation
```bash
# Test config without running full discovery
./bin/cp-discovery -config configs/config.yaml -output /tmp/test.json
echo $?  # Should be 0 for success
```

## Integration Examples

### With Monitoring Systems

**Prometheus Alert:**
```bash
# Export metrics format (future enhancement)
./bin/cp-discovery -format json | jq -r '.clusters[].kafka.cluster_metrics'
```

**Grafana Dashboard:**
```bash
# Periodic update for dashboard
while true; do
  ./bin/cp-discovery -output /var/www/grafana/kafka-status.json
  sleep 300  # Every 5 minutes
done
```

**Slack Notifications:**
```bash
#!/bin/bash
./bin/cp-discovery -output /tmp/discovery.json -detailed

ISSUES=$(jq '[.clusters[] | select(.status != "healthy")] | length' /tmp/discovery.json)
if [ "$ISSUES" -gt 0 ]; then
  curl -X POST -H 'Content-type: application/json' \
    --data "{\"text\":\"⚠️  $ISSUES Kafka clusters need attention\"}" \
    $SLACK_WEBHOOK_URL
fi
```
