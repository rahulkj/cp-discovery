# API Endpoints Reference

This document lists all the API endpoints used by the Confluent Platform Discovery tool to fetch detailed information from REST Proxy and Control Center.

## REST Proxy v3 API Endpoints

The tool uses the Kafka REST Proxy v3 API to gather comprehensive cluster information.

### Cluster Information

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /v3/clusters` | List clusters | Cluster IDs and metadata |
| `GET /v3/clusters/{cluster_id}` | Cluster details | Cluster ID, controller info, metadata |
| `GET /v3/clusters/{cluster_id}/brokers` | List brokers | Broker IDs, hosts, ports, racks |
| `GET /v3/clusters/{cluster_id}/brokers/{broker_id}` | Broker details | Individual broker information |
| `GET /v3/clusters/{cluster_id}/brokers/{broker_id}/configs` | Broker configs | All broker configurations including process.roles |

### Controller Information

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET {controller.related}` | Active controller | Controller broker details |
| `GET /v3/clusters/{cluster_id}/brokers/{broker_id}/configs/process.roles` | KRaft roles | Controller role detection |

### Topic Information

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /v3/clusters/{cluster_id}/topics` | List topics | Topic names, partition counts, replication factors |
| `GET /v3/clusters/{cluster_id}/topics/{topic_name}` | Topic details | Detailed topic configuration, partitions, replication |
| `GET /v3/clusters/{cluster_id}/topics/{topic_name}/configs` | Topic configs | Retention policies, compression, etc. |

### Consumer Groups

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /v3/clusters/{cluster_id}/consumer-groups` | List consumer groups | Group IDs, states, partition assignors |
| `GET /v3/clusters/{cluster_id}/consumer-groups/{group_id}` | Consumer group details | Group state, coordinator, members |
| `GET /v3/clusters/{cluster_id}/consumer-groups/{group_id}/consumers` | Group members | Consumer IDs, instance IDs, client IDs |
| `GET /v3/clusters/{cluster_id}/consumer-groups/{group_id}/lag-summary` | Consumer lag | Total lag for consumer group |

### Access Control Lists (ACLs)

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /v3/clusters/{cluster_id}/acls` | List ACLs | Resource types, principals, operations, permissions |

### Cluster Configuration

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /v3/clusters/{cluster_id}/broker-configs` | Cluster configs | Cluster-level configuration settings |

### Security Configuration

Retrieved from broker configs:
- `sasl.enabled.mechanisms` - SASL mechanisms (PLAIN, SCRAM-SHA-256, etc.)
- `listener.security.protocol.map` - Protocol mappings
- `security.inter.broker.protocol` - Inter-broker security
- `listeners` / `advertised.listeners` - Listener configurations

---

## Confluent Control Center API Endpoints

The tool uses Control Center's monitoring API to gather information about monitored components.

### Health & Version

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /health` | Health check | Control Center availability |
| `GET /api/version` | Version info | Control Center version |

### Kafka Clusters

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /2.0/clusters/kafka` | List Kafka clusters | Cluster IDs, names, broker counts, topic counts |
| `GET /2.0/clusters/kafka/{cluster_id}` | Cluster details | Detailed cluster information |
| `GET /2.0/clusters/kafka/{cluster_id}/health` | Cluster health | Health status of Kafka cluster |
| `GET /2.0/clusters/kafka/{cluster_id}/brokers` | Brokers | Broker details for cluster |
| `GET /2.0/clusters/kafka/{cluster_id}/topics` | Topics | Topic list and partition counts |

### Kafka Connect Clusters

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /2.0/clusters/connect` | List Connect clusters | Connect cluster IDs and names |
| `GET /2.0/clusters/connect/{cluster_id}/connectors` | Connectors | Connector names, states, types |
| `GET /2.0/clusters/connect/{cluster_id}/workers` | Workers | Worker node information |

### Schema Registry Clusters

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /2.0/clusters/schema-registry` | List SR clusters | Schema Registry cluster IDs |
| `GET /2.0/clusters/schema-registry/{cluster_id}` | SR details | Version, schema count |

### ksqlDB Clusters

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /2.0/clusters/ksql` | List ksqlDB clusters | ksqlDB cluster IDs and names |
| `GET /2.0/clusters/ksql/{cluster_id}` | ksqlDB details | Query, stream, and table counts |

### Consumer Lag Monitoring

| Endpoint | Purpose | Data Retrieved |
|----------|---------|----------------|
| `GET /2.0/monitoring/consumer-groups/lag` | Consumer lag | Total lag across all consumer groups |

---

## Data Enrichment

### REST Proxy Data Enrichment

1. **Controller Mode Detection**: Multi-method approach
   - Check `process.roles` config for explicit KRaft roles
   - Compare active controller ID with broker IDs
   - Analyze broker ID patterns (sequential IDs indicate KRaft)
   - Determine: `kraft-combined`, `kraft-separated`, or `zookeeper`

2. **Topic Categorization**: Internal vs External
   - Internal topics: Start with `_` or match patterns like `connect-*`, `_confluent*`
   - External topics: All other topics

3. **Security Configuration Aggregation**:
   - Parse SASL mechanisms from configs
   - Detect security protocols from listener configs
   - Determine SSL/TLS enablement
   - Build authentication method string

4. **Consumer Group Activity**:
   - Count active groups (state = "Stable" or "PreparingRebalance")
   - Fetch member counts for each group
   - Calculate total consumer groups

### Control Center Data Enrichment

1. **Multi-Component Aggregation**:
   - Kafka clusters: Aggregate broker, topic, partition counts
   - Connect clusters: Count connectors, workers, failures
   - Schema Registry: Collect schema counts and versions
   - ksqlDB: Aggregate queries, streams, tables

2. **Health Monitoring**:
   - Fetch health status for each Kafka cluster
   - Identify failed connectors in Connect clusters

3. **Consumer Lag Analysis**:
   - Aggregate total consumer lag across all monitored clusters
   - Provide cluster-wide lag visibility

---

## Authentication

All API calls support three authentication methods (in priority order):

1. **Bearer Token** - `Authorization: Bearer <token>`
2. **API Key** - Custom header (default: `X-API-Key`)
3. **Basic Auth** - `Authorization: Basic <base64(user:pass)>`

Applied via unified authentication helper functions:
- `ApplyRestProxyAuth()`
- `ApplyControlCenterAuth()`

---

## Error Handling

- **Graceful Degradation**: If an endpoint fails, the tool continues with available data
- **Timeout**: All HTTP requests have a 10-15 second timeout
- **Non-Blocking**: Consumer groups, ACLs, and detailed configs are fetched only in detailed mode
- **Fallback Detection**: Multiple methods for controller mode detection if primary method fails

---

## Performance Considerations

### Detailed Mode

When `detailed: true` in config:
- **REST Proxy**: Fetches consumer group members, ACLs, cluster configs
- **Control Center**: Fetches detailed info for all monitored components
- **Trade-off**: More comprehensive data vs. longer discovery time

### Minimal Mode

When `detailed: false`:
- **REST Proxy**: Skips consumer group members, ACLs, and cluster configs
- **Control Center**: Fetches only summary information
- **Benefit**: Faster discovery, smaller output files

### Parallelization

- All component discoveries run in parallel (goroutines)
- Within REST Proxy: Sequential fetching to respect rate limits
- Within Control Center: Sequential fetching per component type

---

## Version Compatibility

### REST Proxy
- **v3 API**: Required for full functionality
- **Minimum Version**: Confluent Platform 6.0+
- **Fallback**: Basic info available from v2 API

### Control Center
- **API Version**: 2.0
- **Minimum Version**: Confluent Platform 5.3+
- **Features**: Full monitoring API support

---

## Example Output Sizes

| Mode | REST Proxy Data | Control Center Data | Typical File Size |
|------|----------------|---------------------|-------------------|
| Minimal (`detailed: false`) | ~500 KB | ~200 KB | ~1 MB |
| Detailed (`detailed: true`) | ~5 MB | ~2 MB | ~10 MB |

*For a 3-broker cluster with 100 topics, 50 consumer groups, 10 connectors*
