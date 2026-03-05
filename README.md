# Confluent Platform Discovery Tool

A comprehensive Go-based tool to scan and discover multiple Confluent Platform installations, gathering detailed metrics and component information.

## ✨ Why This Tool?

- **Minimal Configuration**: Only 2 fields required! Auto-discovers all component URLs
- **Smart Defaults**: Reduces config from 50+ fields to as few as 2 fields (88% reduction!)
- **Secure by Default**: Support for all Kafka security protocols and shared authentication
- **Environment Variables**: Keep secrets safe with `${VAR}` syntax
- **Parallel Discovery**: Scans multiple clusters and components simultaneously
- **Flexible Output**: JSON, YAML, and console summary formats

## Features

- **Multi-Cluster Discovery**: Scan multiple Confluent Platform clusters in parallel
- **Auto-Discovery**: Automatically detects component URLs from Kafka broker host
- **Shared Authentication**: Single auth configuration applies to all components
- **Comprehensive Component Detection**:
  - Kafka Brokers (ZooKeeper and KRaft modes)
  - Schema Registry
  - Kafka Connect (with source/sink connector classification)
  - ksqlDB
  - REST Proxy
  - Confluent Control Center
  - Prometheus
  - Alertmanager

- **Multiple Authentication Methods**:
  - Basic Authentication (username/password)
  - Bearer Token (OAuth, JWT)
  - API Key with custom headers
  - Priority-based auth selection
  - Per-component or shared authentication

- **Detailed Metrics Collection**:
  - **Kafka Cluster**: Broker count, controller detection, topic inventory
  - **Topics**: Internal vs external categorization, partitions, replication factors, retention policies
  - **Consumer Groups**: Total count, active groups, member counts, lag metrics
  - **Security**: ACLs, authentication mechanisms, SSL/TLS detection
  - **Cluster Configuration**: Important cluster-level settings
  - **Schema Registry**: Schema counts, subjects, node counts, version
  - **Kafka Connect**: Connector inventory (source/sink), worker counts, connector states
  - **ksqlDB**: Queries, streams, tables, node counts
  - **Control Center**:
    - Monitored Kafka clusters with broker/topic/partition counts
    - Connect clusters with connector and worker counts
    - Schema Registry clusters with schema counts
    - ksqlDB clusters with query/stream/table counts
    - Total consumer lag across all monitored clusters
  - **Prometheus**: Targets (up/down), HA status, node count
  - **Alertmanager**: Cluster size, peers, active alerts
  - **Network & Storage**: Throughput metrics, disk usage statistics

## Installation

### Prerequisites

- Go 1.21 or later
- Access to Confluent Platform installations
- librdkafka (for confluent-kafka-go)

### Install librdkafka

**macOS:**
```bash
brew install librdkafka
```

**Ubuntu/Debian:**
```bash
sudo apt-get install librdkafka-dev
```

**RHEL/CentOS:**
```bash
sudo yum install librdkafka-devel
```

### Build the Tool

```bash
go mod download
go build -o cp-discovery
```

## Quick Start Configuration

### Minimal Configuration (Just 2 Fields!)

```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "localhost:9092"
```

That's it! The tool auto-discovers:
- Schema Registry at http://localhost:8081
- Kafka Connect at http://localhost:8083
- ksqlDB at http://localhost:8088
- REST Proxy at http://localhost:8082

### With Security (6 Fields)

```yaml
clusters:
  - name: "prod-cluster"
    kafka:
      bootstrap_servers: "broker:9093"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "PLAIN"
      sasl_username: "kafka-user"
      sasl_password: "secret"
    shared_auth:  # Applied to all REST components
      username: "admin"
      password: "admin-secret"
```

## Documentation

### Configuration Files

- **[config.yaml](config.yaml)** - Simple examples to get started
- **[config-minimal.yaml](config-minimal.yaml)** - Minimal configurations (2-6 fields)
- **[config-auth-examples.yaml](config-auth-examples.yaml)** - Authentication patterns (Basic Auth, Bearer Token, API Key)
- **[config-production.yaml](config-production.yaml)** - Production-ready multi-environment setup
- **[config-complete.yaml](config-complete.yaml)** - Complete reference with all available options
- **[CONFIG_REFERENCE.md](docs/CONFIG_REFERENCE.md)** - Comprehensive configuration documentation

### Technical Documentation

- **[API_ENDPOINTS.md](docs/API_ENDPOINTS.md)** - Complete list of REST Proxy and Control Center API endpoints used
- **[ENHANCEMENTS.md](docs/ENHANCEMENTS.md)** - Detailed guide to enhanced discovery capabilities
- **[CONTROL_CENTER_AS_SOURCE.md](docs/CONTROL_CENTER_AS_SOURCE.md)** - Using Control Center as primary discovery source
- **[PROMETHEUS_METRICS.md](docs/PROMETHEUS_METRICS.md)** - Comprehensive cluster metrics from Prometheus

## Authentication Methods

The tool supports three authentication types with priority ordering (Bearer Token > API Key > Basic Auth):

### 1. Basic Authentication (Most Common)
```yaml
schema_registry:
  url: "https://sr.example.com:8081"
  basic_auth_username: "admin"
  basic_auth_password: "secret"
```

### 2. Bearer Token (OAuth/JWT)
```yaml
kafka_connect:
  url: "https://connect.example.com:8083"
  bearer_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. API Key
```yaml
ksqldb:
  url: "https://ksqldb.example.com:8088"
  api_key: "my-api-key-12345"
  api_key_header: "X-API-Key"  # Optional, defaults to "X-API-Key"
```

### Shared Authentication
Apply the same credentials to all components:
```yaml
kafka:
  bootstrap_servers: "broker:9092"

shared_auth:
  username: "admin"
  password: "admin-secret"

# All components (SR, Connect, ksqlDB, REST Proxy, Control Center) will use shared_auth
# unless they specify their own authentication
```

See [config-auth-examples.yaml](config-auth-examples.yaml) for complete authentication examples.

### Security Protocols

Supported values for `security_protocol`:
- `PLAINTEXT` - No encryption or authentication
- `SSL` - SSL/TLS encryption
- `SASL_PLAINTEXT` - SASL authentication, no encryption
- `SASL_SSL` - SASL authentication with SSL/TLS encryption

Supported values for `sasl_mechanism`:
- `PLAIN`
- `SCRAM-SHA-256`
- `SCRAM-SHA-512`

## Usage

### Basic Usage

```bash
./cp-discovery
```

### Custom Configuration File

```bash
./cp-discovery -config /path/to/config.yaml
```

## Output

The tool generates two types of output:

### 1. Console Summary

Real-time progress and summary displayed in the terminal:

```
Starting Confluent Platform Discovery for 2 cluster(s)...
Discovering cluster: production-cluster...
Discovering cluster: staging-cluster...

================================================================================
CONFLUENT PLATFORM DISCOVERY SUMMARY
================================================================================
Timestamp: 2026-03-04T10:30:00Z
Total Clusters: 2

Cluster: production-cluster [success]
--------------------------------------------------------------------------------
  Kafka:
    Brokers: 3
    Controller: kraft
    Topics: 150
    Total Partitions: 450
    Throughput: 125.50 MB/s in, 98.25 MB/s out
    Total Disk Usage: 2500.00 GB
  Schema Registry:
    Version: 7.6.0
    Schemas: 75
  Kafka Connect:
    Version: 7.6.0
    Total Connectors: 12
    Source Connectors: 5
    Sink Connectors: 7
  ksqlDB:
    Version: 0.29.0
    Queries: 8
    Streams: 15
    Tables: 10
```

### 2. Detailed JSON/YAML Report

Complete discovery data saved to the configured output file:

```json
{
  "timestamp": "2026-03-04T10:30:00Z",
  "total_clusters": 2,
  "clusters": [
    {
      "name": "production-cluster",
      "status": "success",
      "kafka": {
        "available": true,
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
      },
      "schema_registry": {
        "available": true,
        "version": "7.6.0",
        "mode": "READWRITE",
        "total_schemas": 75,
        "subjects": ["orders-value", "customers-value"]
      },
      "kafka_connect": {
        "available": true,
        "version": "7.6.0",
        "total_connectors": 12,
        "sink_connectors": 7,
        "source_connectors": 5,
        "connectors": [
          {
            "name": "postgres-source",
            "type": "source",
            "state": "RUNNING",
            "tasks": 1
          }
        ]
      },
      "ksqldb": {
        "available": true,
        "version": "0.29.0",
        "queries": 8,
        "streams": 15,
        "tables": 10
      },
      "rest_proxy": {
        "available": true,
        "version": "v3+"
      }
    }
  ]
}
```

## Metrics Collected

### Kafka Cluster
- Broker count and IDs
- Controller type (ZooKeeper vs KRaft)
- ZooKeeper node count (for ZooKeeper-based clusters)
- Topic count and details
- Partition distribution
- Replication factors
- Retention policies (time and size)
- Network throughput (bytes in/out, messages in)
- Disk usage per broker and total
- Under-replicated partitions

### Schema Registry
- Service version
- Operation mode (READWRITE, READONLY, etc.)
- Total schema count
- Subject list
- Availability status

### Kafka Connect
- Service version
- Total connector count
- Source connector count
- Sink connector count
- Individual connector details (name, type, state, task count)
- Connector classification

### ksqlDB
- Service version
- Running query count
- Stream count
- Table count
- Availability status

### REST Proxy
- Service version and API support
- Cluster ID and metadata
- Broker count and controller information (KRaft/ZooKeeper mode detection)
- Topic counts (total, internal, external)
- Partition counts and average replication factor
- **Consumer Groups**:
  - Total consumer group count
  - Active consumer group count
  - Per-group state and partition assignor
  - Member counts per group
- **Security & Access Control**:
  - ACL count and details (resource type, principal, operation, permission)
  - SASL mechanisms configured
  - Security protocols in use
  - SSL/TLS enabled status
- **Cluster Configuration**:
  - Important cluster-level settings
  - Retention policies
  - Replication settings
  - Topic defaults
- Detailed broker information (host, port, controller role)

### Confluent Control Center
- Service version and health status
- Total monitored cluster count
- **Monitored Kafka Clusters** (per cluster):
  - Cluster ID and name
  - Broker count
  - Topic count
  - Partition count
  - Health status
- **Connect Clusters** (per cluster):
  - Connector count
  - Worker count
  - Failed connector count
- **Schema Registry Clusters** (per cluster):
  - Schema count
  - Version information
- **ksqlDB Clusters** (per cluster):
  - Running query count
  - Stream count
  - Table count
- **Consumer Lag Monitoring**:
  - Total consumer lag across all monitored clusters

### Prometheus
- Service version
- Node count
- Targets (up/down)
- High Availability status
- Availability status

### Alertmanager
- Service version
- Cluster size and peers
- Active alert count
- Availability status

## Advanced Features

### Parallel Discovery

The tool discovers all clusters in parallel for optimal performance:
- Each cluster is scanned independently
- Components within a cluster are also discovered concurrently
- Safe concurrent access to shared resources

### Error Handling

- Graceful degradation: partial failures don't stop the entire discovery
- Detailed error reporting per component
- Cluster status indicators (success, partial, failed)

### Security Support

**Kafka Authentication:**
- Multiple SASL mechanisms (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, GSSAPI, OAUTHBEARER)
- SSL/TLS encryption support
- All security protocols (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL)

**REST API Authentication:**
- HTTP Basic Authentication (username/password)
- Bearer Token authentication (OAuth, JWT)
- API Key authentication with custom headers
- Shared authentication across all components
- Per-component authentication override
- Priority-based auth selection

**Credential Management:**
- Environment variable support with `${VAR}` syntax
- Secure credential handling
- No plaintext secrets in configs

## Limitations and Future Enhancements

### Current Limitations

1. **Real-time Metrics**: Network throughput and disk usage metrics are currently placeholders. In production, these should be fetched from:
   - JMX metrics
   - Confluent Metrics Reporter
   - Prometheus exporters
   - Confluent Control Center API

2. **ZooKeeper Detection**: ZooKeeper node counting requires additional implementation with ZooKeeper client library

### Planned Enhancements

- [ ] JMX integration for real-time metrics
- [ ] Prometheus metrics support
- [ ] ZooKeeper client for accurate node counting
- [ ] Confluent Cloud API support
- [ ] HTML report generation
- [ ] Trend analysis and historical comparison
- [ ] Alert configuration for thresholds
- [ ] Grafana dashboard export

## Troubleshooting

### Connection Errors

If you encounter connection errors:

1. Verify network connectivity to all endpoints
2. Check firewall rules
3. Validate credentials
4. Ensure services are running

### SSL/TLS Issues

For SSL certificate problems:

```yaml
# Add SSL configuration (future enhancement)
kafka:
  ssl_ca_location: "/path/to/ca-cert"
  ssl_certificate_location: "/path/to/client-cert"
  ssl_key_location: "/path/to/client-key"
```

### Permission Errors

Ensure the user has necessary permissions:
- Kafka: Describe/Read permissions on cluster and topics
- Schema Registry: Read access to schemas
- Connect: Read access to connectors
- ksqlDB: Execute permissions for SHOW statements

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This tool is provided as-is for Confluent Platform discovery and monitoring purposes.

## Command-Line Arguments (New!)

The tool now supports flexible command-line arguments that override configuration file settings:

### Available Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-config` | string | `configs/config.yaml` | Path to configuration file |
| `-output` | string | (from config) | Output file path (overrides config) |
| `-format` | string | (from config) | Output format: `json` or `yaml` |
| `-detailed` | bool | false | Enable detailed discovery mode |
| `-view` | bool | false | Open report in web browser after discovery |
| `-view-file` | string | "" | View existing report file in browser (skip discovery) |
| `-port` | int | 8080 | Port for web view server (used with -view) |

### Usage Examples

```bash
# Use custom output file
./bin/cp-discovery -output /tmp/my-report.json

# Change output format
./bin/cp-discovery -format yaml -output report.yaml

# Enable detailed mode
./bin/cp-discovery -detailed

# Combine all options
./bin/cp-discovery \
  -config configs/config-production.yaml \
  -output /var/reports/kafka-$(date +%Y%m%d).json \
  -format json \
  -detailed
```

For more examples, see [USAGE_EXAMPLES.md](docs/USAGE_EXAMPLES.md)

## Web Viewer (New!)

View discovery reports in your browser with a modern, interactive HTML interface:

```bash
# View an existing report
./bin/cp-discovery -view-file test-report.json

# Run discovery and view results (uses temporary file, auto-cleanup)
./bin/cp-discovery -view

# Run discovery and save report for later
./bin/cp-discovery -view -output my-report.json

# Use custom port
./bin/cp-discovery -view-file report.json -port 8888
```

**Features:**
- 🎨 Modern gradient UI with responsive design
- 📊 Tabbed interface (Overview, Clusters, Raw JSON)
- 📈 Summary cards with key metrics
- 🎯 Component cards for each Confluent Platform service
- 🚦 Color-coded status badges
- 🌐 No external dependencies - fully self-contained
- 🗑️ **Auto-cleanup** - When using `-view` without `-output`, creates temporary file and cleans up on exit

For complete documentation, see [WEB_VIEWER.md](docs/WEB_VIEWER.md)

## Enhanced Output (New!)

### Network Throughput
The tool now displays network throughput metrics in the console output:

```
  Kafka:
    Network Throughput:
      Bytes In: 125.50 MB/s
      Bytes Out: 256.75 MB/s
      Messages In: 50000.00 msg/s
```

### Storage Details
Storage information is now prominently displayed:

```
  Kafka:
    Storage:
      Total Disk Usage: 1250.50 GB
```

### Health Metrics
Cluster health indicators are clearly shown:

```
  Kafka:
    Health:
      Under-Replicated Partitions: 5
```

For complete feature documentation, see [NEW_FEATURES.md](docs/NEW_FEATURES.md)

