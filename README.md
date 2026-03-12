# Confluent Platform Discovery Tool

A **pure Go** tool to scan and discover multiple Confluent Platform installations, gathering detailed metrics and component information.

🎯 **Single binary** • 🚀 **No dependencies** • 📊 **Interactive viewer** • ⚡ **Fast parallel discovery**

```bash
# Build once, run anywhere
goreleaser build --clean --snapshot

# Discover your clusters
./dist/cp-discovery-Darwin-arm64 -config my-config.yaml -output data.json

# View results in browser
open viewer.html
```

## ✨ Why This Tool?

### Easy to Use
- **Minimal Configuration**: Only 2 fields required! Auto-discovers all component URLs
- **Smart Defaults**: Reduces config from 50+ fields to as few as 2 fields (88% reduction!)
- **Secure by Default**: Support for all Kafka security protocols and shared authentication
- **Environment Variables**: Keep secrets safe with `${VAR}` syntax

### Pure Go Implementation
- **No Dependencies**: Single static binary - no librdkafka or C libraries required
- **Cross-Platform**: Build for Linux, macOS, Windows without Docker
- **Fast Compilation**: Quick builds with true cross-compilation
- **Easy Distribution**: Just copy the binary - it works everywhere

### Powerful Features
- **Real-Time Progress**: Interactive progress bar shows discovery status as it happens
- **Parallel Discovery**: Scans multiple clusters and components simultaneously
- **Flexible Output**: JSON, YAML, console summary, and interactive web viewer
- **Comprehensive Metrics**: Topics, partitions, storage, connectors, consumer groups, and more
- **Security Support**: SASL/PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, SSL/TLS, mTLS

## 🆕 Recent Enhancements

### v2.0 - Pure Go Migration (Latest)
- ✅ **Pure Go Implementation**: Migrated from `confluent-kafka-go` to `kafka-go` (pure Go)
  - No more CGO dependencies or librdkafka requirement
  - True cross-compilation without Docker
  - Smaller, faster builds with single static binary
- ✅ **Interactive Progress Bar**: Real-time visual feedback during discovery
  - Shows current cluster and component being discovered
  - Color-coded progress with percentage and completion count
  - Three-stage process indicator (Discovery → Report → Summary)
- ✅ **Simplified Distribution**: Binary + `viewer.html` = complete solution
- ✅ **Standalone Web Viewer**: Interactive HTML viewer with file upload (no server required)
- ✅ **Cleaner CLI**: Removed embedded web server, simplified flags

### Previous Enhancements
- ✅ **SSL/TLS Support**: Full SSL certificate configuration for secure Kafka connections
- ✅ **Topic Storage Tracking**: Calculates per-topic and cluster-wide storage from partition offsets
- ✅ **Schema Mapping**: Auto-links Schema Registry subjects to Kafka topics
- ✅ **Enhanced Connector Discovery**:
  - Uses `/connectors?expand=status` and `/connectors?expand=info` API calls
  - Captures `connector.class` and `quickstart` template information
  - Accurate source/sink classification from Kafka Connect API
- ✅ **Controller Count**: Displays KRaft controller count or ZooKeeper controller status
- ✅ **Always-On Details**: Topic and connector information always included

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
  - OAuth/SSO (Client Credentials flow with automatic token management)
  - LDAP (Active Directory, OpenLDAP)
  - Basic Authentication (username/password)
  - Bearer Token (pre-configured OAuth, JWT)
  - API Key with custom headers
  - Priority-based auth selection (OAuth > LDAP > Bearer > API Key > Basic)
  - Per-component or shared authentication
  - Token caching for OAuth to minimize requests

- **Detailed Metrics Collection**:
  - **Kafka Cluster**:
    - Broker count, controller detection (KRaft/ZooKeeper), controller node count
    - Topic inventory with storage size, retention, and partition details
    - Cluster-wide storage metrics calculated from partition offsets
  - **Topics**:
    - Internal vs external categorization
    - Partitions, replication factors, retention policies
    - Per-topic storage size (calculated from offsets)
    - Associated Schema Registry subjects auto-linked
  - **Consumer Groups**: Total count, active groups, member counts, lag metrics
  - **Security**:
    - ACLs, authentication mechanisms, SSL/TLS detection
    - Full SSL/TLS support with client certificates
  - **Cluster Configuration**: Important cluster-level settings
  - **Schema Registry**:
    - Schema counts, subjects, node counts, version
    - Subject-to-topic mapping for data lineage
  - **Kafka Connect**:
    - Connector inventory with accurate source/sink classification
    - Connector class and quickstart template detection
    - Worker counts, connector states, task counts
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

- Go 1.23 or later
- Access to Confluent Platform installations

**No external dependencies required!** This tool uses pure Go libraries and compiles to a single binary.

### Quick Build

```bash
# Using Make (recommended)
make build

# Or manually with Go
go mod download
go build -o bin/cp-discovery ./cmd/cp-discovery

# Or use GoReleaser for production builds
goreleaser build --clean --snapshot
```

The binary will be created at:
- **Make**: `bin/cp-discovery`
- **Manual**: `cp-discovery` (current directory)
- **GoReleaser**: `dist/cp-discovery-<OS>-<ARCH>`

For multi-platform builds and releases, see the [Building from Source](#building-from-source) section below.

## Building from Source

### Using GoReleaser (Recommended)

GoReleaser provides the easiest way to build cross-platform binaries:

```bash
# Install GoReleaser (macOS)
brew install goreleaser

# Build for all platforms (no Docker required!)
goreleaser build --clean --snapshot

# Or create a full release with archives
goreleaser release --clean --snapshot
```

**Output** (in `dist/` directory):
- `cp-discovery-Darwin-x86_64` - macOS Intel
- `cp-discovery-Darwin-arm64` - macOS Apple Silicon
- `cp-discovery-Linux-x86_64` - Linux Intel/AMD
- `cp-discovery-Linux-arm64` - Linux ARM

Plus compressed `.tar.gz` archives and checksums.

### Using Make

For quick local development builds:

```bash
# Build for your current platform
make build

# Output: bin/cp-discovery

# Other useful targets
make test        # Run tests
make fmt         # Format code
make clean       # Remove build artifacts
make help        # Show all targets
```

### Manual Build

```bash
# Install dependencies
go mod download

# Build for current platform
go build -o cp-discovery ./cmd/cp-discovery

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o cp-discovery-linux ./cmd/cp-discovery
```

### Build for Release

Create a tagged release with GoReleaser:

```bash
# Tag the release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Build release artifacts
goreleaser release --clean
```

### Why No Docker Required?

This tool uses **pure Go libraries** (specifically `kafka-go` instead of `confluent-kafka-go`), which means:
- ✅ No CGO dependencies
- ✅ No external C libraries (like librdkafka)
- ✅ True cross-compilation without Docker
- ✅ Smaller, faster builds
- ✅ Single static binary

## Getting Started

### Quick Start (5 Minutes)

**1. Download or Build**
```bash
# Using GoReleaser (recommended)
goreleaser build --clean --snapshot

# Your binary will be at:
# dist/cp-discovery-Darwin-arm64 (or Linux-x86_64, etc.)

# Or use Make
make build
# Binary at: bin/cp-discovery
```

**2. Create Configuration**
```bash
# Copy the template
cp configs/config-template.yaml my-config.yaml

# Edit and add your Kafka broker address (minimum required)
cat > my-config.yaml << EOF
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "localhost:9092"
EOF
```

**3. Run Discovery**
```bash
# Using GoReleaser binary
./dist/cp-discovery-Darwin-arm64 -config my-config.yaml -output data.json

# Or using Make binary
./bin/cp-discovery -config my-config.yaml -output data.json
```

**4. View Results**

**Option A: Console** (instant feedback)
```bash
./bin/cp-discovery -config my-config.yaml
# See summary directly in terminal
```

**Option B: Web Viewer** (interactive dashboard)
```bash
# 1. Generate JSON report
./bin/cp-discovery -config my-config.yaml -output data.json

# 2. Open viewer in browser
open viewer.html  # macOS
# or
xdg-open viewer.html  # Linux

# 3. Load data.json file using the "Choose File" button
```

**Option C: File Export** (for automation)
```bash
# JSON format
./bin/cp-discovery -output discovery-report.json

# YAML format
./bin/cp-discovery -format yaml -output discovery-report.yaml
```

### Step-by-Step Configuration

**Step 1: Choose Your Configuration Level**

Pick the configuration file that matches your needs:

| File | Use Case | Fields Required |
|------|----------|-----------------|
| `config-minimal.yaml` | Local development | 2 (name, bootstrap_servers) |
| `config.yaml` | Simple production | 6 (+ security) |
| `config-production.yaml` | Multi-environment | 10+ (+ all components) |
| `config-complete.yaml` | Reference | All available options |

**Step 2: Edit Configuration**

Minimal setup (2 fields):
```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "broker:9092"
```

With security (6 fields):
```yaml
clusters:
  - name: "prod-cluster"
    kafka:
      bootstrap_servers: "broker:9093"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "PLAIN"
      sasl_username: "${KAFKA_USER}"      # Use env vars!
      sasl_password: "${KAFKA_PASSWORD}"
    shared_auth:
      username: "${API_USER}"
      password: "${API_PASSWORD}"
```

**Step 3: Set Environment Variables** (optional but recommended)
```bash
export KAFKA_USER="my-user"
export KAFKA_PASSWORD="my-secret"
export API_USER="admin"
export API_PASSWORD="admin-secret"
```

**Step 4: Run Discovery**
```bash
./bin/cp-discovery -config my-config.yaml -output data.json
```

**Step 5: Explore Results**

See the [Viewing Your Results](#viewing-your-results) section for details on console output, file exports, and the interactive web viewer.

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

**Getting Started:**
- **[config-template.yaml](configs/config-template.yaml)** - 📋 **START HERE** - Copy and fill in your values
- **[config.yaml](configs/config.yaml)** - Simple working examples
- **[config-minimal.yaml](configs/config-minimal.yaml)** - Minimal configurations (2-6 fields)

**Reference & Examples:**
- **[config-complete.yaml](configs/config-complete.yaml)** - Complete reference with all available properties
- **[config-auth-examples.yaml](configs/config-auth-examples.yaml)** - Authentication patterns (Basic Auth, Bearer Token, API Key)
- **[config-ssl-examples.yaml](configs/config-ssl-examples.yaml)** - SSL/TLS configuration examples (One-way SSL, mTLS, SASL_SSL)
- **[config-production.yaml](configs/config-production.yaml)** - Production-ready multi-environment setup

**Documentation:**
- **[CONFIG_REFERENCE.md](docs/CONFIG_REFERENCE.md)** - Comprehensive configuration documentation

### Technical Documentation

- **[API_ENDPOINTS.md](docs/API_ENDPOINTS.md)** - Complete list of REST Proxy and Control Center API endpoints used
- **[ENHANCEMENTS.md](docs/ENHANCEMENTS.md)** - Detailed guide to enhanced discovery capabilities
- **[CONTROL_CENTER_AS_SOURCE.md](docs/CONTROL_CENTER_AS_SOURCE.md)** - Using Control Center as primary discovery source
- **[PROMETHEUS_METRICS.md](docs/PROMETHEUS_METRICS.md)** - Comprehensive cluster metrics from Prometheus

## Authentication Methods

The tool supports five authentication types with priority ordering (OAuth SSO > LDAP > Bearer Token > API Key > Basic Auth):

### 1. OAuth/SSO Authentication (Highest Priority)
Automatic token retrieval using OAuth 2.0 Client Credentials flow:
```yaml
schema_registry:
  url: "https://sr.example.com:8081"
  oauth_enabled: true
  oauth_client_id: "${SR_OAUTH_CLIENT_ID}"
  oauth_client_secret: "${SR_OAUTH_CLIENT_SECRET}"
  oauth_token_url: "https://auth.example.com/oauth/token"
  oauth_scopes: "schema-registry.read schema-registry.write"
```

**Features:**
- Automatic token retrieval and refresh
- Token caching to minimize auth requests
- Supports standard OAuth 2.0 client credentials flow
- Compatible with Keycloak, Okta, Auth0, and other OAuth providers

### 2. LDAP Authentication
Enterprise directory authentication with LDAP:
```yaml
kafka_connect:
  url: "https://connect.example.com:8083"
  ldap_enabled: true
  ldap_server: "ldaps://ldap.example.com:636"  # or ldap://ldap.example.com:389
  ldap_username: "${LDAP_USER}"
  ldap_password: "${LDAP_PASS}"
  ldap_base_dn: "ou=users,dc=example,dc=com"
```

**Features:**
- Support for both LDAP and LDAPS (LDAP over SSL)
- Falls back to Basic Auth with LDAP credentials
- Configurable base DN for user lookup
- Compatible with Active Directory and OpenLDAP

### 3. Bearer Token (Pre-configured Token)
For pre-obtained OAuth tokens or JWT:
```yaml
kafka_connect:
  url: "https://connect.example.com:8083"
  bearer_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 4. API Key Authentication
```yaml
ksqldb:
  url: "https://ksqldb.example.com:8088"
  api_key: "my-api-key-12345"
  api_key_header: "X-API-Key"  # Optional, defaults to "X-API-Key"
```

### 5. Basic Authentication (Most Common)
```yaml
schema_registry:
  url: "https://sr.example.com:8081"
  basic_auth_username: "admin"
  basic_auth_password: "secret"
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

## Configuration Properties Reference

### Quick Reference - All Available Properties

<details>
<summary><b>Click to expand full configuration properties</b></summary>

#### Kafka Configuration (Required)
| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `bootstrap_servers` | string | ✅ Yes | Comma-separated broker addresses |
| `security_protocol` | string | No | PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL |
| `sasl_mechanism` | string | No | PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, GSSAPI, OAUTHBEARER |
| `sasl_username` | string | No | SASL username (use env vars) |
| `sasl_password` | string | No | SASL password (use env vars) |
| `ssl_ca_location` | string | No | Path to CA certificate |
| `ssl_cert_location` | string | No | Path to client certificate (mTLS) |
| `ssl_key_location` | string | No | Path to client private key (mTLS) |
| `ssl_key_password` | string | No | Password for encrypted private key |
| `ssl_endpoint_identification` | string | No | Hostname verification: `https` or `none` |

#### Shared Authentication (Optional)
| Property | Type | Description |
|----------|------|-------------|
| `username` | string | Username for all REST components |
| `password` | string | Password for all REST components |

#### Component Configuration (All Optional - Auto-discovered)
Each component supports:
- `url` - Component URL (auto-discovered if not specified)
- **OAuth/SSO Authentication:**
  - `oauth_enabled` - Enable OAuth authentication
  - `oauth_client_id` - OAuth client ID
  - `oauth_client_secret` - OAuth client secret
  - `oauth_token_url` - OAuth token endpoint URL
  - `oauth_scopes` - OAuth scopes (space-separated)
- **LDAP Authentication:**
  - `ldap_enabled` - Enable LDAP authentication
  - `ldap_server` - LDAP server URL (ldap:// or ldaps://)
  - `ldap_username` - LDAP username
  - `ldap_password` - LDAP password
  - `ldap_base_dn` - LDAP base DN for user lookup
- **Other Authentication:**
  - `basic_auth_username` / `basic_auth_password` - Basic authentication
  - `bearer_token` - Bearer token authentication (pre-configured OAuth/JWT)
  - `api_key` / `api_key_header` - API key authentication

**Supported Components:**
- `schema_registry` - Default: http://broker-host:8081
- `kafka_connect` - Default: http://broker-host:8083
- `ksqldb` - Default: http://broker-host:8088
- `rest_proxy` - Default: http://broker-host:8082
- `control_center` - Default: http://broker-host:9021
- `prometheus` - Default: http://broker-host:9090
- `alertmanager` - Default: http://broker-host:9093

#### Component Overrides (Optional)
| Property | Type | Description |
|----------|------|-------------|
| `disable_schema_registry` | boolean | Skip Schema Registry discovery |
| `disable_kafka_connect` | boolean | Skip Kafka Connect discovery |
| `disable_ksqldb` | boolean | Skip ksqlDB discovery |
| `disable_rest_proxy` | boolean | Skip REST Proxy discovery |
| `disable_control_center` | boolean | Skip Control Center discovery |
| `disable_prometheus` | boolean | Skip Prometheus discovery |
| `disable_alertmanager` | boolean | Skip Alertmanager discovery |

#### Output Configuration
| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `format` | string | json | Output format: `json` or `yaml` |
| `file` | string | discovery-report.json | Output file path |
| `detailed` | boolean | false | Include detailed topic/connector info |

</details>

## Usage

### Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-config` | string | `configs/config.yaml` | Path to configuration file |
| `-output` | string | (from config) | Output file path |
| `-format` | string | `json` | Output format: `json` or `yaml` |
| `-detailed` | bool | false | Enable detailed discovery mode |

### Basic Usage

```bash
# Run with default config
./bin/cp-discovery

# Use custom config file
./bin/cp-discovery -config /path/to/config.yaml

# Save to specific file
./bin/cp-discovery -output /tmp/my-report.json

# Change output format to YAML
./bin/cp-discovery -format yaml -output report.yaml

# Enable detailed mode
./bin/cp-discovery -detailed

# Combine options
./bin/cp-discovery \
  -config configs/production.yaml \
  -output reports/$(date +%Y%m%d)-discovery.json \
  -format json \
  -detailed
```

### Common Usage Patterns

**1. Quick Local Check:**
```bash
./bin/cp-discovery
# Uses default config (configs/config.yaml)
# Outputs to console
```

**2. Production Scan with Report:**
```bash
./bin/cp-discovery \
  -config configs/production.yaml \
  -output reports/prod-$(date +%Y%m%d).json
```

**3. Multi-Cluster Discovery:**
```bash
# Your config.yaml has multiple clusters
./bin/cp-discovery -config multi-cluster.yaml -output all-clusters.json
```

**4. YAML Output for GitOps:**
```bash
./bin/cp-discovery \
  -format yaml \
  -output kafka-state.yaml
# Commit to Git for change tracking
```

**5. Automated Monitoring:**
```bash
#!/bin/bash
# Run discovery every hour and save with timestamp
while true; do
  ./bin/cp-discovery \
    -config production.yaml \
    -output "reports/$(date +%Y%m%d-%H%M).json"
  sleep 3600
done
```

**6. View in Browser:**
```bash
# Generate report
./bin/cp-discovery -output data.json

# Open viewer
open viewer.html

# Load data.json in the browser UI
```

**7. Compare Environments:**
```bash
# Discover dev
./bin/cp-discovery -config dev.yaml -output dev.json

# Discover prod
./bin/cp-discovery -config prod.yaml -output prod.json

# Compare (using jq or diff tools)
diff <(jq -S . dev.json) <(jq -S . prod.json)
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
- OAuth/SSO authentication (Client Credentials flow) with automatic token management
- LDAP authentication (Active Directory, OpenLDAP)
- HTTP Basic Authentication (username/password)
- Bearer Token authentication (pre-configured OAuth, JWT)
- API Key authentication with custom headers
- Shared authentication across all components
- Per-component authentication override
- Priority-based auth selection (OAuth > LDAP > Bearer > API Key > Basic)

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

### SSL/TLS Configuration

The tool supports SSL/TLS encryption for Kafka connections:

```yaml
kafka:
  bootstrap_servers: "broker:9093"
  security_protocol: "SSL"  # or "SASL_SSL"

  # SSL Certificate Configuration
  ssl_ca_location: "/path/to/ca-cert.pem"
  ssl_cert_location: "/path/to/client-cert.pem"
  ssl_key_location: "/path/to/client-key.pem"
  ssl_key_password: "${SSL_KEY_PASSWORD}"  # Optional, if key is encrypted
  ssl_endpoint_identification: "https"  # Optional, default is "https"
```

**SSL Configuration Options:**

| Option | Description | Required |
|--------|-------------|----------|
| `ssl_ca_location` | Path to CA certificate file | Yes (for SSL) |
| `ssl_cert_location` | Path to client certificate file | Yes (for mTLS) |
| `ssl_key_location` | Path to client private key file | Yes (for mTLS) |
| `ssl_key_password` | Password for encrypted private key | No |
| `ssl_endpoint_identification` | Hostname verification algorithm (`https` or `none`) | No |

**Common SSL Scenarios:**

1. **SSL without client authentication:**
```yaml
kafka:
  bootstrap_servers: "broker:9093"
  security_protocol: "SSL"
  ssl_ca_location: "/path/to/ca-cert.pem"
```

2. **SSL with mutual TLS (mTLS):**
```yaml
kafka:
  bootstrap_servers: "broker:9093"
  security_protocol: "SSL"
  ssl_ca_location: "/path/to/ca-cert.pem"
  ssl_cert_location: "/path/to/client-cert.pem"
  ssl_key_location: "/path/to/client-key.pem"
```

3. **SASL_SSL (SASL with TLS encryption):**
```yaml
kafka:
  bootstrap_servers: "broker:9093"
  security_protocol: "SASL_SSL"
  sasl_mechanism: "PLAIN"
  sasl_username: "kafka-user"
  sasl_password: "${KAFKA_PASSWORD}"
  ssl_ca_location: "/path/to/ca-cert.pem"
```

**Troubleshooting SSL Issues:**

If you encounter SSL certificate errors:

1. **Verify certificate paths** - Ensure all certificate files exist and are readable
2. **Check certificate validity** - Certificates must not be expired
3. **Validate certificate chain** - CA certificate must match the broker's certificate
4. **Disable hostname verification** (not recommended for production):
   ```yaml
   ssl_endpoint_identification: "none"
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

## Viewing Your Results

### Real-Time Progress

Watch discovery progress in real-time with an interactive progress bar:

```
🔍 Starting discovery for 2 cluster(s)...

[1/3] Discovering multi-cluster → Kafka Connect  75% [=================================>      ] (12/16)
```

The progress bar shows:
- **Stage indicator** (1/3, 2/3, 3/3): Discovery → Report → Summary
- **Current cluster** being scanned (color-coded in yellow)
- **Component** being discovered (Kafka, Schema Registry, Connect, ksqlDB, etc.)
- **Progress percentage** and completion count
- **Visual progress bar** with green indicators

### Console Output

After discovery completes, view a comprehensive summary in your terminal:

```
================================================================================
CONFLUENT PLATFORM DISCOVERY SUMMARY
================================================================================
Timestamp: 2026-03-05T12:41:46-06:00
Total Clusters: 2

--------------------------------------------------------------------------------
NODE COUNT SUMMARY (Across All Clusters)
--------------------------------------------------------------------------------
  Kafka Brokers:           4
  KRaft Controllers:       4
  Schema Registry Nodes:   2
  Kafka Connect Workers:   2
  ksqlDB Nodes:            2
  REST Proxy Instances:    2
  Control Center Instances: 1

Cluster: multi-cluster [healthy]
--------------------------------------------------------------------------------
  Kafka:
    Brokers: 3
    Controller: kraft (Controllers: 3)
    Topics: 16 (Internal: 10, External: 6)
    Total Partitions: 246
    Storage:
      Total Cluster Storage: 24.50 GB
...
```

### JSON/YAML Files

Export detailed reports for automation and archival:

```bash
# JSON format (default)
./bin/cp-discovery -output discovery-report.json

# YAML format
./bin/cp-discovery -format yaml -output discovery-report.yaml
```

### Interactive Web Viewer

View reports in a beautiful, interactive dashboard:

1. **Generate Report**:
   ```bash
   ./bin/cp-discovery -output data.json
   ```

2. **Open Viewer**:
   ```bash
   # macOS
   open viewer.html

   # Linux
   xdg-open viewer.html

   # Windows
   start viewer.html
   ```

3. **Load Your Data**:
   - Click "Choose File" button
   - Select your `data.json` file
   - Explore your clusters!

**Viewer Features**:
- 📊 **Overview Tab**: Summary cards with cluster-wide metrics
- 🎯 **Kafka Details Tab**: Complete topic and broker information
- 🔌 **Components Tab**: Connectors, consumer groups, and more
- 📝 **Raw JSON Tab**: View the complete data structure
- 🎨 **Modern UI**: Responsive design with color-coded statuses
- 🚀 **Standalone**: No server required - works offline!

## Viewer Distribution

The `viewer.html` file is automatically included in release archives when using GoReleaser:

```bash
# Build a release
goreleaser release --clean --snapshot

# Archives include:
# ├── cp-discovery-<OS>-<ARCH>
# ├── viewer.html          ← Interactive web viewer
# ├── configs/             ← Sample configurations
# └── README.md
```

This makes it easy to distribute the tool with its viewer to users who don't have the source code.

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

