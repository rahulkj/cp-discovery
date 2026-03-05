# Project Structure

This document describes the organization of the cp-discovery project.

## Directory Layout

```
cp-discovery/
├── cmd/
│   └── cp-discovery/
│       └── main.go              # Application entry point
├── internal/
│   ├── model/
│   │   └── models.go           # All data structures and types
│   ├── config/
│   │   ├── config.go           # Configuration loading
│   │   └── helpers.go          # Configuration defaults and validation
│   ├── http/
│   │   └── auth.go             # HTTP authentication helpers
│   └── discovery/
│       ├── kafka.go            # Kafka discovery
│       ├── schema_registry.go  # Schema Registry discovery
│       ├── kafka_connect.go    # Kafka Connect discovery
│       ├── ksqldb.go           # ksqlDB discovery
│       ├── rest_proxy.go       # REST Proxy discovery
│       ├── control_center.go   # Control Center discovery
│       ├── prometheus.go       # Prometheus discovery
│       └── alertmanager.go     # Alertmanager discovery
├── configs/
│   ├── config.yaml             # Default configuration
│   ├── config-minimal.yaml     # Minimal example
│   ├── config-production.yaml  # Production example
│   ├── config-complete.yaml    # Complete example with all options
│   ├── config-auth-examples.yaml # Authentication examples
│   ├── config-advanced.yaml    # Advanced configuration
│   ├── example-local.yaml      # Local development example
│   └── rj-config.yml           # Custom configuration
├── bin/
│   └── cp-discovery     # Compiled binary
├── docs/
│   ├── API_ENDPOINTS.md        # API endpoint documentation
│   ├── CHANGELOG.md            # Version history
│   ├── CLEANUP_SUMMARY.md      # Project restructuring summary
│   ├── CONFIG_REFERENCE.md     # Configuration reference
│   ├── CONFIG_OPTIMIZATION.md  # Configuration optimization guide
│   ├── CONTROL_CENTER_AS_SOURCE.md # Using Control Center as source
│   ├── ENHANCEMENTS.md         # Feature enhancements log
│   ├── FEATURES.md             # Feature documentation
│   ├── INDEX.md                # Documentation index
│   ├── NEW_FEATURES.md         # v2.0.0 new features
│   ├── PROMETHEUS_METRICS.md   # Prometheus metrics guide
│   ├── PROJECT_STRUCTURE.md    # This file
│   ├── QUICKSTART.md           # Quick start guide
│   ├── RELEASE_NOTES.md        # v2.0.0 release notes
│   ├── SUMMARY.md              # Project summary
│   └── USAGE_EXAMPLES.md       # Usage examples and best practices
├── docker-compose.yaml         # Docker Compose configuration
├── Dockerfile                  # Docker build configuration
├── Makefile                    # Build automation
├── README.md                   # Main documentation
├── go.mod                      # Go module definition
└── go.sum                      # Go dependency checksums
```

## Package Organization

### cm./cp-discovery
The main application entry point. Contains:
- Command-line argument parsing
- Orchestration of discovery across all components
- Output formatting (JSON/YAML)
- Console summary printing

### internal/model
All data structures and type definitions:
- Configuration structures (Config, ClusterConfig, component configs)
- Report structures (DiscoveryReport, component reports)
- Supporting types (TopicInfo, BrokerInfo, ConnectorInfo, etc.)

### internal/config
Configuration management:
- `config.go`: YAML configuration loading
- `helpers.go`: Default values, validation, environment variable expansion

### internal/http
HTTP client utilities:
- `auth.go`: Authentication helpers for different auth methods (Basic Auth, Bearer Token, API Key)

### internal/discovery
Component discovery implementations:
- Each file handles discovery for a specific Confluent Platform component
- All functions are exported (e.g., `DiscoverKafka`, `DiscoverSchemaRegistry`)
- Uses model types for all data structures

## Build and Run

### Build
```bash
go build -o bin/cp-discovery ./cm./cp-discovery
```

### Run
```bash
# Using default config
./bin/cp-discovery

# Using custom config
./bin/cp-discovery -config configs/config-production.yaml
```

### Docker
```bash
# Build
docker build -t cp-discovery .

# Run
docker run -v $(pwd)/configs:/app/configs cp-discovery
```

## Configuration Files

All configuration files are located in the `configs/` directory:

- **config.yaml**: Default configuration (used when no -config flag provided)
- **config-minimal.yaml**: Minimal working example
- **config-production.yaml**: Production-ready configuration
- **config-complete.yaml**: All available options documented
- **config-auth-examples.yaml**: Various authentication method examples
- **config-advanced.yaml**: Advanced features and optimizations

## Output

Discovery reports are written to files specified in the configuration:
```yaml
output:
  format: json        # or yaml
  file: discovery-report.json
  detailed: true      # Include detailed information
```

## Documentation

Main documentation in root, detailed docs in `docs/` directory:
- **README.md**: Main project documentation (root directory)
- **docs/INDEX.md**: Documentation index and navigation
- **docs/QUICKSTART.md**: Quick start guide
- **docs/USAGE_EXAMPLES.md**: Comprehensive usage examples
- **docs/CONFIG_REFERENCE.md**: Complete configuration reference
- **docs/FEATURES.md**: Feature documentation
- **docs/NEW_FEATURES.md**: v2.0.0 new features
- **docs/ENHANCEMENTS.md**: Enhancement history
- **docs/PROMETHEUS_METRICS.md**: Prometheus metrics guide
- **docs/CONTROL_CENTER_AS_SOURCE.md**: Control Center integration guide
- **docs/CHANGELOG.md**: Version history
- **docs/RELEASE_NOTES.md**: v2.0.0 release highlights

## Module Information

- **Module Path**: github.com/rahulkj/cp-discovery
- **Go Version**: 1.21+
- **Main Dependencies**:
  - github.com/confluentinc/confluent-kafka-go/v2
  - gopkg.in/yaml.v3
