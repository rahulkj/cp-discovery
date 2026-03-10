# Original Project Prompt (Reconstructed)

**Project Name:** Confluent Platform Discovery Tool

**Date:** March 2026

---

## Project Overview

Create a production-ready command-line tool in Go that can scan and discover multiple Confluent Platform installations across different environments, collecting comprehensive metrics and component information from all platform components.

## Core Requirements

### 1. Primary Functionality

**Objective:** Build a discovery tool that can:
- Connect to multiple Confluent Platform clusters simultaneously
- Auto-discover all components in each cluster
- Collect detailed metrics and configuration information
- Generate comprehensive reports in multiple formats
- Provide real-time progress feedback during discovery

### 2. Supported Components

The tool must discover and collect metrics from:

1. **Kafka Cluster**
   - Broker count, IDs, hosts, and ports
   - Controller detection (ZooKeeper vs KRaft mode)
   - Controller node counts
   - Topic inventory with detailed information
   - Partition distribution and replication factors
   - Retention policies (time and size-based)
   - Storage metrics (per-topic and cluster-wide)
   - Network throughput metrics
   - Cluster health indicators (under-replicated partitions)

2. **Schema Registry**
   - Service version and operation mode
   - Total schema count and subject enumeration
   - Subject-to-topic mapping for data lineage
   - Node count in the cluster
   - Availability status

3. **Kafka Connect**
   - Service version and worker counts
   - Connector inventory with accurate source/sink classification
   - Connector class and quickstart template detection
   - Connector states and task counts
   - Use `/connectors?expand=status` and `/connectors?expand=info` API calls

4. **ksqlDB**
   - Service version and node counts
   - Running query count
   - Stream and table counts
   - Availability status

5. **REST Proxy**
   - Service version and API support level
   - Consumer group information (total, active, state, members)
   - ACL counts and security configurations
   - SASL mechanisms and security protocols
   - SSL/TLS status
   - Cluster configuration settings

6. **Confluent Control Center**
   - Service version and health status
   - Monitored Kafka clusters with metrics (brokers, topics, partitions)
   - Connect clusters (connectors, workers, failures)
   - Schema Registry clusters (schema counts)
   - ksqlDB clusters (queries, streams, tables)
   - Total consumer lag across all monitored clusters

7. **Prometheus** (Optional)
   - Service version and node count
   - Target status (up/down)
   - High Availability status

8. **Alertmanager** (Optional)
   - Service version and cluster size
   - Peer information
   - Active alert counts

### 3. Configuration Requirements

**Goal:** Minimize configuration complexity while maintaining flexibility

**Configuration Optimization:**
- Support minimal configuration (as few as 2 fields: cluster name + bootstrap servers)
- Implement smart auto-discovery for component URLs based on Kafka broker host
- Use standard Confluent Platform ports (8081, 8083, 8088, 8082, 9021, 9090, 9093)
- Support shared authentication across all REST components
- Allow per-component authentication overrides
- Support environment variable expansion with `${VAR_NAME}` syntax
- Provide component disable flags for faster targeted discovery

**Configuration Examples Needed:**
- `config-minimal.yaml` - Minimal 2-field configuration
- `config.yaml` - Simple production setup (6 fields)
- `config-production.yaml` - Multi-environment production setup
- `config-complete.yaml` - Complete reference with all options
- `config-auth-examples.yaml` - Various authentication patterns
- `config-ssl-examples.yaml` - SSL/TLS configuration examples

### 4. Security and Authentication

**Kafka Security:**
- Support all security protocols: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL
- Support SASL mechanisms: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, GSSAPI, OAUTHBEARER
- Full SSL/TLS support with client certificates (mTLS)
- SSL certificate configuration (CA, client cert, client key, key password)
- Optional hostname verification control

**REST API Authentication:**
- HTTP Basic Authentication (username/password)
- Bearer Token authentication (OAuth, JWT)
- API Key authentication with custom headers
- Priority-based authentication selection (Bearer Token > API Key > Basic Auth)
- Shared authentication configuration with per-component overrides

**Credential Management:**
- Environment variable support for secrets
- No plaintext secrets in configuration files
- Secure credential handling throughout

### 5. Technical Requirements

**Language and Architecture:**
- Pure Go implementation (Go 1.21+)
- **Critical:** Use pure Go Kafka client library (e.g., `kafka-go`) - NO CGO dependencies
- No external C libraries (like librdkafka)
- Single static binary compilation
- True cross-compilation support (Linux, macOS, Windows)
- Support for multiple architectures (x86_64, ARM64)

**Performance:**
- Parallel cluster discovery (all clusters scanned simultaneously)
- Concurrent component discovery within each cluster
- Safe concurrent access to shared resources
- Configurable timeouts
- Graceful degradation on component failures

**Code Organization:**
- Clean package structure with separation of concerns
- `cmd/cp-discovery` - Main application entry point
- `internal/model` - All data structures and types
- `internal/config` - Configuration loading and validation
- `internal/http` - HTTP authentication helpers
- `internal/discovery` - Component discovery implementations

### 6. Output Requirements

**Console Output:**
- Real-time interactive progress bar showing:
  - Current cluster being discovered
  - Current component being scanned
  - Progress percentage and completion count
  - Color-coded visual indicators
  - Three-stage process (Discovery → Report → Summary)
- Comprehensive summary after discovery completion:
  - Cluster-wide node count summary
  - Per-cluster detailed metrics
  - Component versions and health status
  - Storage metrics and throughput
  - Clear status indicators (healthy, partial, failed)

**File Export:**
- JSON format (default)
- YAML format (optional)
- Configurable output file path
- Detailed mode for comprehensive information
- Always include topic and connector details

**Web Viewer:**
- Create a standalone HTML viewer (`viewer.html`)
- No server required - runs entirely in browser
- File upload capability for loading JSON reports
- Multiple tabs for different views:
  - Overview tab with summary cards
  - Kafka Details tab with topics and brokers
  - Components tab with connectors and consumer groups
  - Raw JSON tab for debugging
- Modern, responsive UI with color-coded statuses
- Works offline

### 7. Build and Distribution

**Build System:**
- Makefile for local development builds
- GoReleaser configuration for production builds
- Multi-platform binary generation:
  - Darwin (macOS): x86_64, ARM64
  - Linux: x86_64, ARM64
  - Windows: x86_64, ARM64
- Compressed archives with checksums
- Include viewer.html in release archives

**Release Artifacts:**
- Compiled binaries for all platforms
- Configuration file examples
- Interactive web viewer (viewer.html)
- README and documentation

**Docker Support:**
- Multi-stage Dockerfile for optimized image size
- Non-root user execution
- Docker Compose configuration for local testing

### 8. Documentation Requirements

Create comprehensive documentation including:

1. **README.md** - Main project documentation with:
   - Feature highlights and benefits
   - Quick start guide (5-minute setup)
   - Installation instructions
   - Usage examples
   - Configuration reference
   - Authentication examples
   - Output format documentation
   - Troubleshooting guide

2. **docs/QUICKSTART.md** - Step-by-step quick start guide

3. **docs/CONFIG_REFERENCE.md** - Complete configuration property reference

4. **docs/USAGE_EXAMPLES.md** - Comprehensive usage patterns and examples

5. **docs/FEATURES.md** - Detailed feature documentation

6. **docs/API_ENDPOINTS.md** - List of all REST API endpoints used

7. **docs/ENHANCEMENTS.md** - Enhancement history and capabilities

8. **docs/PROJECT_STRUCTURE.md** - Codebase organization

9. **docs/CHANGELOG.md** - Version history

10. **docs/RELEASE_NOTES.md** - Major release highlights

### 9. Error Handling and Resilience

- Graceful degradation when components are unavailable
- Continue discovery even if individual components fail
- Detailed error reporting per component
- Cluster status indicators (success, partial, failed)
- Timeout handling for unresponsive endpoints
- Clear error messages for troubleshooting

### 10. Enhanced Features

**Storage Tracking:**
- Calculate per-topic storage from partition offsets
- Aggregate cluster-wide storage metrics
- Display in human-readable format (GB, TB)

**Schema Mapping:**
- Auto-link Schema Registry subjects to Kafka topics
- Show which topics use schemas

**Controller Information:**
- Display KRaft controller count
- Show ZooKeeper controller status
- Detect cluster mode automatically

**Node Counting:**
- Count nodes for all multi-node components
- Display in summary output
- Include in JSON/YAML reports

**Network Metrics:**
- Bytes in/out per second
- Messages in per second
- Throughput calculations

## Success Criteria

### Configuration Optimization
- Achieve 88% reduction in required configuration fields
- Reduce from 17+ fields to as few as 2 fields for basic setup
- Maintain backward compatibility with existing configurations

### User Experience
- 5-minute setup from download to first discovery
- Clear, informative progress indicators
- Intuitive configuration with helpful examples
- Beautiful, easy-to-use web viewer

### Technical Excellence
- Pure Go with zero CGO dependencies
- Single static binary (no runtime dependencies)
- Fast compilation and build times
- Cross-platform without Docker or special tooling
- Production-ready error handling

### Performance
- Complete discovery in 5-20 seconds (depending on cluster count)
- Memory usage under 100 MB
- Minimal CPU usage (I/O bound operations)
- Parallel execution for optimal speed

### Deliverables
- 8 Go source files with clean architecture
- 4+ configuration examples covering all scenarios
- 10+ documentation files
- Interactive web viewer
- Build automation (Makefile, GoReleaser, Docker)
- GitHub Actions CI/CD pipeline

## Future Enhancements (Nice to Have)

- JMX integration for real-time broker metrics
- Confluent Cloud API support
- Historical trend analysis
- Prometheus metrics export endpoint
- Grafana dashboard templates
- Alert threshold configuration
- Schema evolution tracking
- Comparison between environments

## Technical Constraints

1. **No CGO Dependencies:** Must use pure Go libraries only
2. **No External Binaries:** Single self-contained executable
3. **Backward Compatible:** Existing configurations must continue to work
4. **Security First:** No plaintext secrets, secure credential handling
5. **Production Ready:** Comprehensive error handling, logging, timeouts
6. **Well Tested:** Validated against real Confluent Platform installations

## Version 2.0 Specific Requirements

**Pure Go Migration:**
- Migrate from `confluent-kafka-go` (uses librdkafka/CGO) to `kafka-go` (pure Go)
- Remove all CGO dependencies
- Maintain feature parity with v1.x
- Improve build speed and cross-compilation

**Interactive Progress:**
- Add real-time progress bar during discovery
- Show current cluster and component being discovered
- Color-coded visual indicators
- Three-stage process tracking

**Standalone Viewer:**
- Remove embedded web server
- Create standalone HTML file that works offline
- File upload for loading JSON reports
- No server dependencies

**Cleaner CLI:**
- Simplify command-line flags
- Remove web server related flags
- Focus on discovery and reporting

---

## Expected Output Example

```
🔍 Starting discovery for 2 cluster(s)...

[1/3] Discovering production-cluster → Kafka  100% [========================================] (8/8)

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

Cluster: production-cluster [healthy]
--------------------------------------------------------------------------------
  Kafka:
    Brokers: 3
    Controller: kraft (Controllers: 3)
    Topics: 150 (Internal: 45, External: 105)
    Total Partitions: 450
    Storage:
      Total Cluster Storage: 2500.00 GB
    Network Throughput:
      Bytes In: 125.50 MB/s
      Bytes Out: 98.25 MB/s
  Schema Registry:
    Version: 7.6.0
    Schemas: 75
    Nodes: 1
  Kafka Connect:
    Version: 7.6.0
    Total Connectors: 12 (Source: 5, Sink: 7)
    Workers: 1
```

---

**End of Prompt**

This prompt represents the comprehensive requirements and vision for the Confluent Platform Discovery Tool as reconstructed from the implemented project.
