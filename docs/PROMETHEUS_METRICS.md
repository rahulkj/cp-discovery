# Prometheus Cluster Metrics

## Overview

When Prometheus is enabled and configured, the discovery tool can fetch comprehensive cluster metrics directly from Prometheus, providing real-time operational visibility into your Kafka cluster.

---

## Prerequisites

### 1. Prometheus Setup

Prometheus must be scraping Kafka metrics using one of these exporters:
- **JMX Exporter** - Most common for Kafka
- **Confluent Metrics Reporter** - Exports to Confluent Control Center and Prometheus
- **Kafka Exporter** - Standalone exporter for Kafka metrics

### 2. Required Metrics

The tool queries the following Prometheus metrics:

#### Throughput Metrics
- `kafka_server_brokertopicmetrics_bytesin_total` - Total bytes received
- `kafka_server_brokertopicmetrics_bytesout_total` - Total bytes sent
- `kafka_server_brokertopicmetrics_messagesin_total` - Total messages received

#### Controller Metrics
- `kafka_controller_kafkacontroller_activecontrollercount` - Active controller count

#### Partition Metrics
- `kafka_server_replicamanager_underreplicatedpartitions` - Under-replicated partitions
- `kafka_server_replicamanager_offlinepartitionscount` - Offline partitions
- `kafka_server_replicamanager_partitioncount` - Total partitions
- `kafka_server_replicamanager_leadercount` - Leader partition count

#### Consumer Metrics
- `kafka_consumergroup_lag` - Consumer group lag (per group)
- `kafka_consumergroup_lag_sum` - Total consumer lag

#### JVM Metrics
- `jvm_memory_bytes_used{area="heap"}` - JVM heap memory used
- `jvm_memory_bytes_max{area="heap"}` - JVM heap memory max
- `process_cpu_seconds_total` - CPU usage

#### Broker Availability
- `up{job=~".*kafka.*"}` - Broker up/down status

---

## Configuration

### Enable Prometheus Discovery

In your `config.yaml`:

```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "localhost:9092"

    prometheus:
      url: "http://localhost:9090"
      # Authentication (if needed)
      basic_auth_username: "${PROM_USER}"
      basic_auth_password: "${PROM_PASS}"

# Enable detailed mode to fetch cluster metrics
output:
  detailed: true  # REQUIRED for cluster metrics
```

**Important**: Cluster metrics are only fetched when `detailed: true` is set in the output configuration.

---

## Metrics Collected

### 1. Throughput Metrics

| Metric | Description | Unit | Calculation |
|--------|-------------|------|-------------|
| `bytes_in_per_sec` | Data ingress rate | Bytes/sec | `rate(kafka_server_brokertopicmetrics_bytesin_total[5m])` |
| `bytes_out_per_sec` | Data egress rate | Bytes/sec | `rate(kafka_server_brokertopicmetrics_bytesout_total[5m])` |
| `messages_in_per_sec` | Message ingress rate | Messages/sec | `rate(kafka_server_brokertopicmetrics_messagesin_total[5m])` |

**Note**: Rates are calculated over a 5-minute window using PromQL's `rate()` function.

### 2. Controller Metrics

| Metric | Description | Expected Value |
|--------|-------------|----------------|
| `active_controller_count` | Number of active controllers | 1 (normal), 0 (election in progress), >1 (split-brain) |

### 3. Partition Metrics

| Metric | Description | Healthy Value |
|--------|-------------|---------------|
| `under_replicated_partitions` | Partitions not fully replicated | 0 |
| `offline_partitions` | Partitions with no leader | 0 |
| `total_partitions` | Total partition count | N/A |
| `leader_count` | Total leader partitions | N/A |

### 4. Broker Metrics

| Metric | Description | Calculation |
|--------|-------------|-------------|
| `total_brokers` | Total broker count | Count of unique instances with Kafka metrics |
| `online_brokers` | Brokers currently up | Count of `up{job=~".*kafka.*"} == 1` |

### 5. Consumer Metrics

| Metric | Description | Healthy Value |
|--------|-------------|---------------|
| `total_consumer_lag` | Aggregate lag across all groups | Low (< 1000) |
| `consumer_groups` | Number of consumer groups | N/A |

### 6. JVM Metrics (Optional)

| Metric | Description | Healthy Range |
|--------|-------------|---------------|
| `avg_heap_used_percent` | Average JVM heap usage | < 75% |
| `avg_cpu_used_percent` | Average CPU usage | < 80% |

---

## Console Output Example

When Prometheus metrics are available:

```
Prometheus:
  Version: 2.45.0
  URL: http://localhost:9090
  Targets: 6 up, 0 down

  Cluster Metrics:
    Brokers: 3 total, 3 online
    Active Controllers: 1
    Throughput: 125.50 MB/s in, 98.25 MB/s out
    Messages/sec: 50000
    Partitions: 450 total, 450 leaders
    Under-replicated: 0
    Consumer Groups: 25 (lag: 1250 messages)
    Avg Heap Usage: 68.5%
    Avg CPU Usage: 45.2%
```

---

## JSON Output Example

```json
{
  "prometheus": {
    "available": true,
    "version": "2.45.0",
    "url": "http://localhost:9090",
    "targets_up": 6,
    "targets_down": 0,
    "cluster_metrics": {
      "bytes_in_per_sec": 131621593.6,
      "bytes_out_per_sec": 103024230.4,
      "messages_in_per_sec": 50000.0,
      "active_controller_count": 1,
      "under_replicated_partitions": 0,
      "offline_partitions": 0,
      "total_partitions": 450,
      "total_brokers": 3,
      "online_brokers": 3,
      "leader_count": 450,
      "total_consumer_lag": 1250,
      "consumer_groups": 25,
      "avg_heap_used_percent": 68.5,
      "avg_cpu_used_percent": 45.2
    }
  }
}
```

---

## Health Indicators

### Critical Issues

🔴 **Immediate Attention Required:**
- `offline_partitions > 0` - Data unavailable
- `active_controller_count != 1` - Controller election issue
- `online_brokers < total_brokers` - Broker(s) down
- `avg_heap_used_percent > 90` - JVM memory pressure

### Warnings

⚠️ **Monitor Closely:**
- `under_replicated_partitions > 0` - Replication lag
- `total_consumer_lag > 10000` - Consumer falling behind
- `avg_heap_used_percent > 75` - Approaching memory limits
- `avg_cpu_used_percent > 80` - High CPU utilization

### Healthy

✅ **Normal Operation:**
- `offline_partitions == 0`
- `under_replicated_partitions == 0`
- `active_controller_count == 1`
- `online_brokers == total_brokers`
- `avg_heap_used_percent < 75`

---

## Prometheus Queries Used

### Throughput (Rate over 5 minutes)
```promql
# Bytes in per second
sum(rate(kafka_server_brokertopicmetrics_bytesin_total[5m]))

# Bytes out per second
sum(rate(kafka_server_brokertopicmetrics_bytesout_total[5m]))

# Messages in per second
sum(rate(kafka_server_brokertopicmetrics_messagesin_total[5m]))
```

### Controller Status
```promql
# Active controller count
sum(kafka_controller_kafkacontroller_activecontrollercount)
```

### Partition Health
```promql
# Under-replicated partitions
sum(kafka_server_replicamanager_underreplicatedpartitions)

# Offline partitions
sum(kafka_server_replicamanager_offlinepartitionscount)

# Total partitions
sum(kafka_server_replicamanager_partitioncount)

# Leader count
sum(kafka_server_replicamanager_leadercount)
```

### Broker Count
```promql
# Total brokers
count(count by (instance) (kafka_server_brokertopicmetrics_bytesin_total))

# Online brokers
count(up{job=~".*kafka.*"} == 1)
```

### Consumer Lag
```promql
# Total consumer lag
sum(kafka_consumergroup_lag_sum)

# Consumer group count
count(count by (consumergroup) (kafka_consumergroup_lag))
```

### JVM Metrics
```promql
# Average heap usage percentage
avg(jvm_memory_bytes_used{area="heap"} / jvm_memory_bytes_max{area="heap"} * 100)

# Average CPU usage percentage
avg(rate(process_cpu_seconds_total{job=~".*kafka.*"}[5m]) * 100)
```

---

## Troubleshooting

### Metric Not Found

**Problem**: Prometheus query returns no data

**Solutions**:
1. Verify JMX Exporter is configured and running
2. Check Prometheus scrape targets are healthy
3. Confirm metric names match your exporter version
4. Review Prometheus logs for scrape errors

### Incorrect Values

**Problem**: Metrics show unexpected values

**Solutions**:
1. Verify time range (tool uses 5-minute rate)
2. Check for clock skew between Kafka and Prometheus
3. Confirm exporters are properly configured
4. Compare with Kafka broker JMX metrics directly

### Missing Metrics

**Problem**: Some metrics unavailable

**Solutions**:
1. Ensure `detailed: true` in config
2. Verify all required exporters are running
3. Check Prometheus retention period
4. Confirm metric labels match expected format

### Authentication Errors

**Problem**: Cannot connect to Prometheus

**Solutions**:
1. Verify Prometheus URL is correct
2. Check authentication credentials
3. Test Prometheus API manually: `curl http://prometheus:9090/api/v1/query?query=up`
4. Review firewall rules and network connectivity

---

## Customization

### Custom Metric Names

If your Prometheus uses different metric names, you can modify the queries in `prometheus.go`:

```go
// Example: Change metric prefix from kafka_server to your custom prefix
metrics.BytesInPerSec = queryPrometheusRate(client, config, "custom_prefix_brokertopicmetrics_bytesin_total")
```

### Additional Metrics

To add custom metrics:

1. Add field to `PrometheusClusterMetrics` struct in `main.go`
2. Query the metric in `getClusterMetricsFromPrometheus()` in `prometheus.go`
3. Display in console output in `main.go`

Example:
```go
// In main.go - Add to PrometheusClusterMetrics
RequestsPerSec float64 `json:"requests_per_sec"`

// In prometheus.go - Query the metric
metrics.RequestsPerSec = queryPrometheusRate(client, config, "kafka_network_requestmetrics_requests_total")

// In main.go - Display in console
if metrics.RequestsPerSec > 0 {
    fmt.Printf("      Requests/sec: %.0f\n", metrics.RequestsPerSec)
}
```

---

## Performance Considerations

### Query Overhead

Each metric requires a Prometheus query:
- **Total Queries**: ~15 queries in detailed mode
- **Query Time**: ~50-100ms per query
- **Total Time**: ~1-2 seconds for all metrics

### Prometheus Load

The tool's queries are lightweight:
- No heavy aggregations
- Short time windows (5 minutes)
- Optimized queries with `sum()` and `count()`

### Caching

Prometheus caches query results:
- Repeated queries are faster
- Recent data is in memory
- First query may be slower

---

## Integration with Other Tools

### Grafana Dashboards

Use the same queries in Grafana:
1. Import Prometheus as data source
2. Create panels with the PromQL queries above
3. Set refresh interval (e.g., 30 seconds)

### Alertmanager

Create alerts based on these metrics:
```yaml
groups:
  - name: kafka_alerts
    rules:
      - alert: OfflinePartitions
        expr: sum(kafka_server_replicamanager_offlinepartitionscount) > 0
        for: 5m
        annotations:
          summary: "Kafka has offline partitions"

      - alert: UnderReplicatedPartitions
        expr: sum(kafka_server_replicamanager_underreplicatedpartitions) > 0
        for: 10m
        annotations:
          summary: "Kafka has under-replicated partitions"
```

---

## Benefits

✅ **Real-Time Monitoring**: Live metrics from Prometheus
✅ **No Additional Load**: Metrics already collected by Prometheus
✅ **Historical Context**: Prometheus provides time-series data
✅ **Aggregated View**: Cluster-wide metrics in one place
✅ **Performance Insights**: Throughput, lag, and JVM metrics
✅ **Health Checks**: Automated detection of issues

---

## Limitations

❌ **Requires Prometheus**: Metrics unavailable without Prometheus
❌ **Exporter Dependency**: Needs JMX exporter or equivalent
❌ **Metric Name Variations**: Different exporters use different names
❌ **Detailed Mode Only**: Metrics not collected in minimal mode
❌ **No Historical Data**: Tool shows current state, not trends

---

## Summary

Prometheus integration provides **comprehensive cluster visibility** without additional load on Kafka brokers. When enabled, the discovery tool becomes a powerful monitoring and health check tool, combining:

- **Discovery Data**: From Kafka Admin API and REST Proxy
- **Operational Metrics**: From Prometheus
- **Health Indicators**: Automated analysis

This makes it ideal for:
- 🔍 Cluster health checks
- 📊 Performance analysis
- 🚨 Incident response
- 📈 Capacity planning
- 🔬 Troubleshooting
