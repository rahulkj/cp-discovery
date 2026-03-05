# Quick Start Guide

## Prerequisites

Make sure you have the following installed:

1. **Go 1.21+**
   ```bash
   go version
   ```

2. **librdkafka** (required for confluent-kafka-go)

   **macOS:**
   ```bash
   brew install librdkafka pkg-config
   ```

   **Ubuntu/Debian:**
   ```bash
   sudo apt-get update
   sudo apt-get install -y librdkafka-dev pkg-config
   ```

   **RHEL/CentOS:**
   ```bash
   sudo yum install -y librdkafka-devel pkgconfig
   ```

## Installation

### Option 1: Build from Source

```bash
# Clone or navigate to the project directory
cd cp-discovery

# Download dependencies
go mod download

# Build the binary
go build -o cp-discovery .
```

### Option 2: Use Make

```bash
# Install dependencies
make install

# Build
make build
```

## Configuration

### Quick Setup for Local Docker Environment

If you're running Confluent Platform locally with Docker:

```bash
# Copy the example configuration
cp example-local.yaml config.yaml

# Edit if needed
nano config.yaml
```

### Full Configuration

Create or edit `config.yaml`:

```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "localhost:9092"
      security_protocol: "PLAINTEXT"  # or SASL_SSL, SASL_PLAINTEXT, SSL
      sasl_mechanism: ""              # PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
      sasl_username: ""
      sasl_password: ""

    schema_registry:
      url: "http://localhost:8081"
      basic_auth_username: ""
      basic_auth_password: ""

    kafka_connect:
      url: "http://localhost:8083"
      basic_auth_username: ""
      basic_auth_password: ""

    ksqldb:
      url: "http://localhost:8088"
      basic_auth_username: ""
      basic_auth_password: ""

    rest_proxy:
      url: "http://localhost:8082"
      basic_auth_username: ""
      basic_auth_password: ""

output:
  format: "json"  # or "yaml"
  file: "discovery-report.json"
  detailed: true
```

## Running the Tool

### Basic Usage

```bash
./cp-discovery
```

This will:
1. Read `config.yaml` from the current directory
2. Scan all configured clusters
3. Display a summary in the console
4. Save detailed report to `discovery-report.json`

### Custom Configuration File

```bash
./cp-discovery -config /path/to/my-config.yaml
```

### Using Make

```bash
# Run with default config
make run

# Run with custom config
make run-config CONFIG=/path/to/config.yaml
```

## Testing with Docker Compose

If you want to test with a local Confluent Platform:

### 1. Start Confluent Platform

```bash
# Create a docker-compose.yaml or use an existing one
docker-compose up -d
```

### 2. Wait for Services to Start

```bash
# Check if services are ready
docker-compose ps

# Wait for all services to be healthy (may take 1-2 minutes)
```

### 3. Run Discovery

```bash
./cp-discovery
```

### 4. View Results

```bash
# View JSON report
cat discovery-report.json | jq .

# Or use any JSON viewer
```

## Example Output

### Console Summary

```
Starting Confluent Platform Discovery for 1 cluster(s)...
Discovering cluster: my-cluster...

================================================================================
CONFLUENT PLATFORM DISCOVERY SUMMARY
================================================================================
Timestamp: 2026-03-04T16:15:00Z
Total Clusters: 1

Cluster: my-cluster [success]
--------------------------------------------------------------------------------
  Kafka:
    Brokers: 3
    Controller: kraft
    Topics: 10
    Total Partitions: 30
    Throughput: 0.00 MB/s in, 0.00 MB/s out
    Total Disk Usage: 0.00 GB
  Schema Registry:
    Version: 7.6.0
    Schemas: 5
  Kafka Connect:
    Version: 7.6.0
    Total Connectors: 2
    Source Connectors: 1
    Sink Connectors: 1
  ksqlDB:
    Version: 0.29.0
    Queries: 0
    Streams: 0
    Tables: 0
  REST Proxy:
    Version: v3+

Discovery completed successfully!
Report saved to: discovery-report.json
```

### JSON Report

The detailed JSON report includes:
- Complete broker information
- All topics with partition counts and retention settings
- All schemas from Schema Registry
- Connector details (name, type, state, tasks)
- ksqlDB queries, streams, and tables
- Error messages for any failed components

## Troubleshooting

### "Cannot connect to Kafka"

Check:
1. Bootstrap servers address is correct
2. Kafka is running: `docker ps` or `systemctl status kafka`
3. Network connectivity: `telnet localhost 9092`
4. Security settings match (PLAINTEXT vs SASL_SSL)

### "Schema Registry not available"

Check:
1. Schema Registry URL is correct
2. Service is running
3. Authentication credentials (if required)

### "Build failed - librdkafka not found"

Install librdkafka:
```bash
# macOS
brew install librdkafka pkg-config

# Ubuntu/Debian
sudo apt-get install librdkafka-dev pkg-config
```

### "Permission denied"

Make the binary executable:
```bash
chmod +x cp-discovery
```

## Next Steps

1. **Add More Clusters**: Edit `config.yaml` to add additional cluster configurations
2. **Schedule Regular Scans**: Use cron to run discovery periodically
3. **Integrate with Monitoring**: Parse the JSON output with your monitoring tools
4. **Customize Output**: Change output format to YAML or implement custom reporters

## Scheduled Discovery (Cron)

To run discovery every hour:

```bash
# Edit crontab
crontab -e

# Add this line (adjust path)
0 * * * * cd /path/t./cp-discovery && ./cp-discovery
```

## Docker Usage (Optional)

You can also run the tool in Docker:

```dockerfile
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc g++ make musl-dev librdkafka-dev pkgconfig

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o cp-discovery .

FROM alpine:latest
RUN apk add --no-cache librdkafka
COPY --from=builder /ap./cp-discovery /usr/local/bin/
COPY config.yaml /config.yaml

CMD ["confluent-discovery", "-config", "/config.yaml"]
```

Build and run:
```bash
docker build -t cp-discovery .
docker run -v $(pwd)/config.yaml:/config.yaml cp-discovery
```

## Support

For issues, questions, or contributions, please refer to the README.md file.
