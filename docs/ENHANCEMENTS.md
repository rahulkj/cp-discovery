# Enhanced REST Proxy and Control Center Discovery

## Overview

The Confluent Platform Discovery tool has been significantly enhanced to fetch comprehensive details from REST Proxy and Control Center API endpoints. This document summarizes all the new capabilities.

## Latest Enhancements (March 2026)

### 1. Controller Count Fix
- **Fixed**: Now shows total number of controller nodes, not just the active controller
- **New Field**: `controller_count` in RestProxyReport
- **Behavior**: Counts all brokers with controller role (KRaft mode)
- **Console Output**: Displays "Controllers: X" to show total controller nodes

### 2. Component Details via Control Center
- **Enhancement**: Control Center now fetches comprehensive component details
- **Benefit**: Can use Control Center as single source of truth for all components
- **Components**: Schema Registry, Kafka Connect, ksqlDB details now richer

**Control Center Now Provides:**
- **Connect**: Source/sink breakdown, running connector count, per-connector details
- **Schema Registry**: Mode (READWRITE/READONLY), subject list, cluster IDs
- **ksqlDB**: Cluster ID and Kafka cluster ID associations

### 3. Prometheus Cluster Metrics
- **New Feature**: Fetch comprehensive cluster metrics from Prometheus
- **Benefit**: Real-time operational visibility without additional load on Kafka
- **Activation**: Enabled in detailed mode when Prometheus is configured

**Prometheus Now Provides:**
- **Throughput**: Bytes in/out per second, messages per second
- **Controller**: Active controller count
- **Partitions**: Total, under-replicated, offline, leader counts
- **Brokers**: Total and online broker counts
- **Consumers**: Total lag and consumer group counts
- **JVM**: Average heap and CPU usage across brokers

---

## REST Proxy Enhancements

### 1. Consumer Groups Discovery

**New Fields:**
- `consumer_groups` - Array of consumer group details
- `consumer_group_count` - Total number of consumer groups
- `active_consumer_groups` - Count of active/stable groups

**Per Consumer Group:**
```json
{
  "group_id": "my-consumer-group",
  "state": "Stable",
  "partition_assignor": "range",
  "member_count": 3,
  "lag": 1250
}
```

**States Tracked:**
- `Stable` - Active and consuming
- `PreparingRebalance` - Rebalancing in progress
- `Empty` - No members
- `Dead` - Inactive group

**API Endpoints Used:**
- `GET /v3/clusters/{cluster_id}/consumer-groups` - List all groups
- `GET /v3/clusters/{cluster_id}/consumer-groups/{group_id}/consumers` - Member count

### 2. Access Control Lists (ACLs)

**New Fields:**
- `acls` - Array of ACL details
- `acl_count` - Total number of ACLs

**Per ACL:**
```json
{
  "resource_type": "TOPIC",
  "resource_name": "orders",
  "pattern_type": "LITERAL",
  "principal": "User:alice",
  "operation": "READ",
  "permission": "ALLOW"
}
```

**Resource Types:**
- `TOPIC` - Topic-level permissions
- `GROUP` - Consumer group permissions
- `CLUSTER` - Cluster-level permissions
- `TRANSACTIONAL_ID` - Transaction permissions

**API Endpoint Used:**
- `GET /v3/clusters/{cluster_id}/acls`

### 3. Cluster Configuration

**New Field:**
- `cluster_config` - Map of important cluster-level settings

**Configurations Tracked:**
```json
{
  "compression.type": "producer",
  "log.retention.hours": "168",
  "log.retention.bytes": "-1",
  "message.max.bytes": "1048588",
  "min.insync.replicas": "2",
  "default.replication.factor": "3",
  "auto.create.topics.enable": "false",
  "delete.topic.enable": "true",
  "offsets.retention.minutes": "10080"
}
```

**Only Non-Default Settings:** To reduce noise, only configurations that differ from defaults are included.

**API Endpoint Used:**
- `GET /v3/clusters/{cluster_id}/broker-configs`

### 4. Enhanced Existing Fields

**Security Configuration:**
- Expanded to include all SASL mechanisms
- Full security protocol mapping
- Listener-specific security detection

**Broker Information:**
- Controller role detection (KRaft vs ZooKeeper)
- Active controller identification
- Rack awareness

**Topic Statistics:**
- Internal vs external categorization
- Average replication factor calculation
- Per-topic retention policies

---

## Control Center Enhancements

### 1. Monitored Kafka Clusters

**New Field:**
- `clusters` - Array of monitored Kafka cluster details

**Per Kafka Cluster:**
```json
{
  "cluster_id": "lkc-12345",
  "cluster_name": "production-kafka",
  "broker_count": 5,
  "topic_count": 150,
  "partition_count": 900,
  "health_status": "HEALTHY"
}
```

**Health States:**
- `HEALTHY` - All brokers online, no issues
- `WARNING` - Some degradation
- `CRITICAL` - Major issues detected
- `UNKNOWN` - Cannot determine health

**API Endpoints Used:**
- `GET /2.0/clusters/kafka` - List clusters
- `GET /2.0/clusters/kafka/{cluster_id}/health` - Health status

### 2. Connect Clusters

**New Field:**
- `connect_clusters` - Array of Connect cluster details

**Per Connect Cluster:**
```json
{
  "cluster_name": "production-connect",
  "connector_count": 25,
  "worker_count": 3,
  "failed_connectors": 1
}
```

**Tracks:**
- Total connectors (source + sink)
- Worker node count
- Failed connector count for alerting

**API Endpoints Used:**
- `GET /2.0/clusters/connect` - List Connect clusters
- `GET /2.0/clusters/connect/{cluster_id}/connectors` - Connector details
- `GET /2.0/clusters/connect/{cluster_id}/workers` - Worker count

### 3. Schema Registry Clusters

**New Field:**
- `schema_registries` - Array of Schema Registry cluster details

**Per Schema Registry:**
```json
{
  "cluster_name": "production-sr",
  "version": "7.6.0",
  "schema_count": 75
}
```

**API Endpoints Used:**
- `GET /2.0/clusters/schema-registry` - List SR clusters
- `GET /2.0/clusters/schema-registry/{cluster_id}` - SR details

### 4. ksqlDB Clusters

**New Field:**
- `ksql_clusters` - Array of ksqlDB cluster details

**Per ksqlDB Cluster:**
```json
{
  "cluster_name": "production-ksql",
  "query_count": 12,
  "stream_count": 8,
  "table_count": 5
}
```

**API Endpoints Used:**
- `GET /2.0/clusters/ksql` - List ksqlDB clusters
- `GET /2.0/clusters/ksql/{cluster_id}` - Query/stream/table counts

### 5. Consumer Lag Monitoring

**New Field:**
- `total_consumer_lag` - Total lag across all monitored consumer groups

**Provides:**
- Cluster-wide consumer lag visibility
- Aggregated metric for monitoring dashboards
- Early warning for consumer performance issues

**API Endpoint Used:**
- `GET /2.0/monitoring/consumer-groups/lag`

---

## Console Output Enhancements

### REST Proxy Output

```
  REST Proxy:
    Version: v3+
    Cluster ID: lkc-abc123
    Brokers: 3
    Controller Mode: kraft-combined
    Consumer Groups: 25 (Active: 22)
    ACLs: 45
    Cluster Configs: 12 custom settings
```

### Control Center Output

```
  Control Center:
    Version: 7.6.0
    URL: https://control-center:9021
    Monitored Clusters: 3

    Kafka Clusters:
      - production-kafka: 5 brokers, 150 topics, 900 partitions [HEALTHY]
      - staging-kafka: 3 brokers, 75 topics, 450 partitions [HEALTHY]
      - dev-kafka: 1 brokers, 20 topics, 60 partitions [WARNING]

    Connect Clusters:
      - production-connect: 25 connectors, 3 workers (1 failed)
      - staging-connect: 10 connectors, 2 workers

    Schema Registries:
      - production-sr: 75 schemas (v7.6.0)
      - staging-sr: 40 schemas (v7.6.0)

    ksqlDB Clusters:
      - production-ksql: 12 queries, 8 streams, 5 tables
      - staging-ksql: 5 queries, 3 streams, 2 tables

    Total Consumer Lag: 15,234 messages
```

---

## JSON Report Structure

### Complete REST Proxy Report Example

```json
{
  "available": true,
  "version": "v3+",
  "broker_count": 3,
  "controller_id": 1,
  "controller_mode": "kraft-combined",
  "cluster_id": "lkc-abc123",
  "topic_count": 150,
  "internal_topics": 12,
  "external_topics": 138,
  "partition_count": 900,
  "avg_replication_factor": 3.0,
  "security_config": {
    "sasl_mechanisms": ["SCRAM-SHA-256", "SCRAM-SHA-512"],
    "security_protocols": ["SASL_SSL", "SSL"],
    "ssl_enabled": true,
    "sasl_enabled": true,
    "authentication_method": "SASL/SCRAM-SHA-256 + SSL/TLS"
  },
  "consumer_groups": [
    {
      "group_id": "order-processor",
      "state": "Stable",
      "partition_assignor": "range",
      "member_count": 3
    }
  ],
  "consumer_group_count": 25,
  "active_consumer_groups": 22,
  "acls": [
    {
      "resource_type": "TOPIC",
      "resource_name": "orders",
      "pattern_type": "LITERAL",
      "principal": "User:alice",
      "operation": "READ",
      "permission": "ALLOW"
    }
  ],
  "acl_count": 45,
  "cluster_config": {
    "min.insync.replicas": "2",
    "default.replication.factor": "3",
    "auto.create.topics.enable": "false"
  }
}
```

### Complete Control Center Report Example

```json
{
  "available": true,
  "version": "7.6.0",
  "url": "https://control-center:9021",
  "monitored_clusters": 3,
  "clusters": [
    {
      "cluster_id": "lkc-prod",
      "cluster_name": "production-kafka",
      "broker_count": 5,
      "topic_count": 150,
      "partition_count": 900,
      "health_status": "HEALTHY"
    }
  ],
  "connect_clusters": [
    {
      "cluster_name": "production-connect",
      "connector_count": 25,
      "worker_count": 3,
      "failed_connectors": 1
    }
  ],
  "schema_registries": [
    {
      "cluster_name": "production-sr",
      "version": "7.6.0",
      "schema_count": 75
    }
  ],
  "ksql_clusters": [
    {
      "cluster_name": "production-ksql",
      "query_count": 12,
      "stream_count": 8,
      "table_count": 5
    }
  ],
  "total_consumer_lag": 15234
}
```

---

## Performance Impact

### Discovery Time

| Component | Minimal Mode | Detailed Mode |
|-----------|-------------|---------------|
| REST Proxy (before) | ~2s | ~3s |
| REST Proxy (after) | ~2s | ~8s |
| Control Center (before) | ~1s | ~2s |
| Control Center (after) | ~1s | ~5s |

*Times for a typical 3-broker cluster with 100 topics, 25 consumer groups, 10 connectors*

### Output File Size

| Mode | Before | After | Increase |
|------|--------|-------|----------|
| Minimal | ~500 KB | ~600 KB | +20% |
| Detailed | ~2 MB | ~10 MB | +400% |

*Detailed mode now includes significantly more information*

---

## Use Cases Enabled

### 1. Consumer Group Monitoring
- Identify stale/empty consumer groups
- Track active consumer distribution
- Monitor consumer group states
- Plan consumer group cleanup

### 2. Security Auditing
- Review all ACL configurations
- Identify overly permissive ACLs
- Audit principal access levels
- Plan least-privilege migrations

### 3. Configuration Management
- Discover non-default cluster settings
- Ensure configuration consistency
- Identify configuration drift
- Plan standardization initiatives

### 4. Multi-Cluster Monitoring (via Control Center)
- Aggregate health across all clusters
- Identify troubled clusters quickly
- Monitor Connect/SR/ksqlDB ecosystem
- Track consumer lag trends

### 5. Capacity Planning
- Analyze broker/topic/partition distribution
- Identify connector worker needs
- Plan schema registry scaling
- Forecast ksqlDB resource requirements

---

## Configuration

### Enable Detailed Discovery

In `config.yaml`:

```yaml
output:
  format: "json"
  file: "discovery-report.json"
  detailed: true  # Enable comprehensive data collection
```

### Disable Components

Skip specific discovery to reduce time:

```yaml
overrides:
  disable_rest_proxy: false    # Keep enabled for cluster details
  disable_control_center: true # Skip if not deployed
```

---

## Future Enhancements

### Potential Additions

1. **REST Proxy:**
   - Partition-level lag per consumer group
   - Topic-level ACL summaries
   - Broker-level metrics (JMX integration)
   - Partition leader distribution analysis

2. **Control Center:**
   - Historical lag trends
   - Alert configuration extraction
   - Connector error details
   - Query performance metrics

3. **Performance:**
   - Caching layer for repeated queries
   - Batch API calls where supported
   - Incremental discovery mode

---

## Migration Guide

### Updating from Previous Versions

No configuration changes required! The tool is backward compatible:

- **Minimal Mode**: Same performance, slightly more data
- **Detailed Mode**: Much more comprehensive data

### Accessing New Data

All new fields are optional and won't break existing parsers:

```python
# Python example
import json

with open('discovery-report.json') as f:
    report = json.load(f)

# Access new consumer group data (safe even if not present)
for cluster in report['clusters']:
    rest_proxy = cluster.get('rest_proxy', {})
    consumer_groups = rest_proxy.get('consumer_groups', [])

    for group in consumer_groups:
        print(f"Group: {group['group_id']}, State: {group['state']}")
```

---

## Summary

The enhanced discovery tool provides **unprecedented visibility** into Confluent Platform deployments:

✅ **Consumer Groups**: Full lifecycle monitoring
✅ **Security**: Complete ACL inventory
✅ **Configuration**: Cluster settings discovery
✅ **Multi-Component**: Unified Control Center view
✅ **Consumer Lag**: Aggregate lag tracking

**Result**: A single command now provides enterprise-grade observability across your entire Confluent Platform ecosystem.
