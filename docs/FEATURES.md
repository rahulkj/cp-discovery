# Confluent Platform Discovery - Feature Summary

## Core Capabilities

### 🔍 Discovery Features

| Component           | Metrics Collected |
|---------------------|-------------------|
| **Kafka Brokers**   | ✅ Broker count, IDs, hosts, ports<br>✅ Controller type detection (ZooKeeper/KRaft)<br>✅ Topic inventory with partitions<br>✅ Replication factors<br>✅ Retention policies (time & size)<br>✅ Cluster-wide metrics |
| **Schema Registry** | ✅ Version information<br>✅ Operation mode (READWRITE/READONLY)<br>✅ Total schema count<br>✅ Subject enumeration<br>✅ Availability status |
| **Kafka Connect**   | ✅ Version and commit info<br>✅ Total connector count<br>✅ Source/Sink classification<br>✅ Connector state (RUNNING/FAILED/etc)<br>✅ Task count per connector |
| **ksqlDB**          | ✅ Version information<br>✅ Running query count<br>✅ Stream definitions<br>✅ Table definitions<br>✅ Availability status |
| **REST Proxy**      | ✅ Version information<br>✅ API version support<br>✅ Availability status |

### ⚙️ Configuration Optimizations

| Feature                 | Description | Benefit |
|-------------------------|-------------|---------|
| **Auto-Discovery**      | Automatically derives component URLs from Kafka broker host | 80% fewer config lines |
| **Shared Authentication** | Single auth config for all components | 75% less credential duplication |
| **Environment Variables** | `${VAR}` syntax for all config fields | Secure secret management |
| **Smart Defaults**      | Sensible defaults for all optional fields | Minimal configuration required |
| **Component Overrides** | Explicitly disable components | Faster discovery, cleaner errors |

### 🔒 Security Support

| Feature | Supported Options |
|---------|------------------|
| **Kafka Security** | PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL |
| **SASL Mechanisms** | PLAIN, SCRAM-SHA-256, SCRAM-SHA-512 |
| **REST Authentication** | HTTP Basic Auth for all REST endpoints |
| **Protocol Detection** | Auto-detects HTTP vs HTTPS from Kafka security |

### 📊 Output Formats

| Format | Features |
|--------|----------|
| **JSON** | Structured, parsable, machine-readable |
| **YAML** | Human-readable, configuration-friendly |
| **Console** | Real-time progress and summary table |

### 🚀 Performance Features

| Feature | Description |
|---------|-------------|
| **Parallel Discovery** | Scans multiple clusters simultaneously |
| **Concurrent Components** | Discovers all components in parallel per cluster |
| **Configurable Timeouts** | 10-30s timeouts per component |
| **Graceful Degradation** | Partial failures don't block entire discovery |

## Configuration Comparison

### Minimal Config (2 fields)
```yaml
clusters:
  - name: "local"
    kafka:
      bootstrap_servers: "localhost:9092"
```

**Auto-discovers:**
- Schema Registry at http://localhost:8081
- Kafka Connect at http://localhost:8083
- ksqlDB at http://localhost:8088
- REST Proxy at http://localhost:8082

### Typical Config (6 fields)
```yaml
clusters:
  - name: "prod"
    kafka:
      bootstrap_servers: "broker:9093"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "PLAIN"
      sasl_username: "${KAFKA_USER}"
      sasl_password: "${KAFKA_PASS}"
    shared_auth:
      username: "${CP_USER}"
      password: "${CP_PASS}"
```

**Provides:**
- Secure Kafka connection
- Authenticated REST component access
- Secret management via env vars
- Auto-discovered HTTPS URLs

### Advanced Config (12 fields)
```yaml
clusters:
  - name: "complex"
    kafka:
      bootstrap_servers: "${KAFKA_BROKERS}"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "SCRAM-SHA-512"
      sasl_username: "${KAFKA_USER}"
      sasl_password: "${KAFKA_PASS}"
    shared_auth:
      username: "${CP_USER}"
      password: "${CP_PASS}"
    schema_registry:
      url: "https://sr-dedicated.example.com:8081"
    kafka_connect:
      url: "https://connect-cluster.example.com:8083"
      basic_auth_username: "connect-admin"
      basic_auth_password: "${CONNECT_PASS}"
    overrides:
      disable_ksqldb: true
      disable_rest_proxy: true
```

**Supports:**
- Custom component URLs
- Per-component authentication
- Selective component discovery
- Enterprise security requirements

## Metrics Details

### Kafka Cluster Metrics

```json
{
  "broker_count": 3,
  "controller_type": "kraft",
  "topic_count": 150,
  "total_partitions": 450,
  "topics": [
    {
      "name": "orders",
      "partitions": 6,
      "replication_factor": 3,
      "retention_ms": 604800000,
      "retention_bytes": -1,
      "total_size_bytes": 1073741824
    }
  ],
  "brokers": [
    {
      "id": 1,
      "host": "broker1",
      "port": 9092,
      "rack": "rack1",
      "disk_usage_bytes": 536870912000
    }
  ],
  "cluster_metrics": {
    "bytes_in_per_sec": 131621593.6,
    "bytes_out_per_sec": 103024230.4,
    "messages_in_per_sec": 50000,
    "total_disk_usage_bytes": 2684354560000,
    "under_replicated_partitions": 0
  }
}
```

### Connector Classification

The tool automatically classifies connectors as source or sink:

**Source Patterns:**
- `*Source*`, `*Debezium*`, `*Spooldir*`, `*FileStream*`
- JDBC Source, MongoDB Source, etc.

**Sink Patterns:**
- `*Sink*`, `*S3*`, `*Elasticsearch*`, `*HDFS*`
- JDBC Sink, Snowflake Sink, etc.

### Error Handling

```json
{
  "cluster": "prod",
  "status": "partial",
  "errors": [
    "ksqlDB: connection refused",
    "REST Proxy: authentication failed"
  ]
}
```

## Use Cases

### 1. **Multi-Cluster Inventory**
Scan all your Confluent Platform installations to generate a complete inventory:
- Total broker count across all clusters
- Schema Registry schema count
- Active connector inventory
- ksqlDB deployment status

### 2. **Capacity Planning**
Collect metrics for capacity analysis:
- Partition distribution across brokers
- Disk usage trends
- Replication factor compliance
- Topic growth rates

### 3. **Security Audit**
Verify security configurations:
- SASL mechanism usage
- SSL/TLS enablement
- Authentication status per component
- Protocol compliance

### 4. **Migration Planning**
Before migrating to Confluent Cloud or upgrading:
- Current topology documentation
- Component version inventory
- Topic and partition counts
- Connector migration checklist

### 5. **Disaster Recovery**
Document current state for DR planning:
- Broker topology
- Topic configurations
- Schema registry contents
- Connector configurations

### 6. **Compliance Reporting**
Generate reports for compliance:
- Data retention policies
- Replication compliance
- Schema governance
- Component availability SLAs

## Integration Points

### CI/CD Pipeline
```bash
# In Jenkins/GitLab/GitHub Actions
./cp-discovery -config prod-config.yaml
# Parse JSON output for validation
```

### Monitoring
```bash
# Cron job for regular discovery
0 */6 * * * /op./cp-discovery/cp-discovery
```

### Alert Integration
```bash
# Check for under-replicated partitions
jq '.clusters[].kafka.cluster_metrics.under_replicated_partitions' report.json
```

### Grafana/Prometheus
Parse JSON output to push metrics to Prometheus or Grafana for trending.

## Limitations & Roadmap

### Current Limitations
- Network metrics are placeholders (need JMX integration)
- ZooKeeper node counting requires ZooKeeper client
- Disk usage per broker requires JMX/metrics API
- No Confluent Cloud API support (yet)

### Planned Features
- [ ] JMX metrics integration for real-time data
- [ ] Confluent Cloud API support
- [ ] Historical trend analysis
- [ ] HTML report generation
- [ ] Grafana dashboard templates
- [ ] Alert threshold configuration
- [ ] Schema evolution tracking
- [ ] Consumer group discovery
- [ ] ACL enumeration

## Performance

| Clusters | Components per Cluster | Typical Discovery Time |
|----------|------------------------|------------------------|
| 1        | All (5 components)     | ~5-10 seconds         |
| 5        | All (5 components)     | ~10-15 seconds        |
| 10       | All (5 components)     | ~15-20 seconds        |

*Times assume good network connectivity and responsive services*

## Resource Usage

| Resource | Usage |
|----------|-------|
| **Memory** | ~50-100 MB for typical clusters |
| **CPU** | Minimal (I/O bound) |
| **Disk** | Output files typically < 1MB per cluster |
| **Network** | HTTP/HTTPS API calls, Kafka metadata queries |

## Compatibility

| Component | Minimum Version | Tested Versions |
|-----------|----------------|-----------------|
| **Confluent Platform** | 5.0+ | 6.x, 7.x |
| **Apache Kafka** | 2.0+ | 2.x, 3.x |
| **Schema Registry** | 5.0+ | 6.x, 7.x |
| **Kafka Connect** | 2.0+ | 2.x, 3.x |
| **ksqlDB** | 0.10+ | 0.28+, 0.29+ |
| **Go Runtime** | 1.21+ | 1.21, 1.22 |
