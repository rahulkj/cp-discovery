# Configuration Reference

Complete reference for all configuration options in the Confluent Platform Discovery tool.

## Table of Contents

- [Quick Start](#quick-start)
- [Cluster Configuration](#cluster-configuration)
- [Component Configuration](#component-configuration)
- [Authentication Methods](#authentication-methods)
- [Output Configuration](#output-configuration)
- [Environment Variables](#environment-variables)
- [Examples](#examples)

---

## Quick Start

### Minimal Configuration

```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "localhost:9092"
```

This auto-discovers all components using default ports.

---

## Cluster Configuration

### Required Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `name` | string | Cluster identifier | `"production-cluster"` |
| `kafka.bootstrap_servers` | string | Comma-separated Kafka brokers | `"broker1:9092,broker2:9092"` |

### Kafka Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `security_protocol` | string | Security protocol: `PLAINTEXT`, `SSL`, `SASL_PLAINTEXT`, `SASL_SSL` | `PLAINTEXT` |
| `sasl_mechanism` | string | SASL mechanism: `PLAIN`, `SCRAM-SHA-256`, `SCRAM-SHA-512`, `GSSAPI`, `OAUTHBEARER` | - |
| `sasl_username` | string | SASL username | - |
| `sasl_password` | string | SASL password | - |

---

## Component Configuration

All components support the same configuration structure:

```yaml
component_name:
  url: "https://host:port"
  # Authentication (choose one)
  basic_auth_username: "user"
  basic_auth_password: "pass"
  # OR
  bearer_token: "token"
  # OR
  api_key: "key"
  api_key_header: "X-API-Key"  # optional
```

### Supported Components

| Component | Default Port | Auto-Discovery |
|-----------|--------------|----------------|
| `schema_registry` | 8081 | ✅ |
| `kafka_connect` | 8083 | ✅ |
| `ksqldb` | 8088 | ✅ |
| `rest_proxy` | 8082 | ✅ |
| `control_center` | 9021 | ✅ |
| `prometheus` | 9090 | ✅ |
| `alertmanager` | 9093 | ✅ |

### Auto-Discovery

If component URLs are not specified, they are auto-discovered using:
- Host from Kafka `bootstrap_servers` (first broker)
- Default port for each component
- Protocol based on Kafka `security_protocol` (HTTP for `PLAINTEXT`, HTTPS for `SSL`/`SASL_SSL`)

Example:
```yaml
kafka:
  bootstrap_servers: "broker1.example.com:9092"
  security_protocol: "SASL_SSL"

# Auto-discovered URLs:
# schema_registry: https://broker1.example.com:8081
# kafka_connect: https://broker1.example.com:8083
# ksqldb: https://broker1.example.com:8088
# rest_proxy: https://broker1.example.com:8082
# control_center: https://broker1.example.com:9021
# prometheus: http://broker1.example.com:9090
# alertmanager: http://broker1.example.com:9093
```

---

## Authentication Methods

### Priority Order

When multiple auth methods are configured, this priority is used:
1. Bearer Token
2. API Key
3. Basic Auth

### Basic Authentication

Most common method for Confluent Platform components:

```yaml
schema_registry:
  url: "https://sr.example.com:8081"
  basic_auth_username: "admin"
  basic_auth_password: "secret123"
```

### Bearer Token (OAuth/JWT)

For OAuth 2.0 or JWT-based authentication:

```yaml
kafka_connect:
  url: "https://connect.example.com:8083"
  bearer_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

Header sent: `Authorization: Bearer <token>`

### API Key

For API key-based authentication:

```yaml
ksqldb:
  url: "https://ksqldb.example.com:8088"
  api_key: "my-api-key-12345"
  api_key_header: "X-API-Key"  # Optional, defaults to "X-API-Key"
```

Custom header example:
```yaml
control_center:
  url: "https://c3.example.com:9021"
  api_key: "my-key"
  api_key_header: "Authorization"  # Some systems use this
```

### Shared Authentication

Apply the same credentials to all components:

```yaml
kafka:
  bootstrap_servers: "broker:9092"

shared_auth:
  username: "admin"
  password: "admin-secret"

# All components will use shared_auth unless they specify their own
```

Component-specific auth overrides shared auth:

```yaml
shared_auth:
  username: "general-user"
  password: "general-pass"

kafka_connect:
  url: "https://connect.example.com:8083"
  basic_auth_username: "connect-admin"  # Overrides shared_auth
  basic_auth_password: "connect-secret"
```

---

## Output Configuration

```yaml
output:
  format: "json"  # or "yaml"
  file: "discovery-report.json"
  detailed: true  # Include detailed information
```

### Output Formats

**JSON** (default):
```yaml
output:
  format: "json"
  file: "report.json"
```

**YAML**:
```yaml
output:
  format: "yaml"
  file: "report.yaml"
```

### Detailed Mode

When `detailed: true`:
- Individual topic configurations (replication factor, retention)
- Full broker list with roles
- Connector details
- Schema subjects
- Prometheus targets
- Alertmanager peers

When `detailed: false`:
- Summary counts only
- Faster discovery
- Smaller output file

---

## Environment Variables

Use `${VAR_NAME}` syntax for environment variables:

```yaml
kafka:
  bootstrap_servers: "${KAFKA_BROKERS}"
  sasl_username: "${KAFKA_USER}"
  sasl_password: "${KAFKA_PASSWORD}"

schema_registry:
  url: "${SR_URL}"
  basic_auth_username: "${SR_USER}"
  basic_auth_password: "${SR_PASSWORD}"
```

Example `.env` file:
```bash
KAFKA_BROKERS=broker1:9092,broker2:9092,broker3:9092
KAFKA_USER=admin
KAFKA_PASSWORD=secret123

SR_URL=https://schema-registry.example.com:8081
SR_USER=sr-admin
SR_PASSWORD=sr-secret
```

---

## Component Overrides

Disable specific components:

```yaml
overrides:
  disable_schema_registry: false
  disable_kafka_connect: false
  disable_ksqldb: false
  disable_rest_proxy: false
  disable_control_center: false
  disable_prometheus: false
  disable_alertmanager: false
```

Use cases:
- Skip unavailable components
- Reduce discovery time
- Prevent error messages for non-deployed components

---

## Examples

### Confluent Cloud

```yaml
clusters:
  - name: "ccloud-prod"
    kafka:
      bootstrap_servers: "${CCLOUD_BOOTSTRAP_SERVERS}"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "PLAIN"
      sasl_username: "${CCLOUD_API_KEY}"
      sasl_password: "${CCLOUD_API_SECRET}"

    schema_registry:
      url: "${CCLOUD_SR_URL}"
      basic_auth_username: "${CCLOUD_SR_KEY}"
      basic_auth_password: "${CCLOUD_SR_SECRET}"

    # Disable on-prem components
    overrides:
      disable_kafka_connect: true
      disable_ksqldb: true
      disable_control_center: true
      disable_prometheus: true
      disable_alertmanager: true
```

### Multi-Datacenter Setup

```yaml
clusters:
  - name: "dc1-production"
    kafka:
      bootstrap_servers: "dc1-broker1:9092,dc1-broker2:9092"
    shared_auth:
      username: "${DC1_USER}"
      password: "${DC1_PASS}"

  - name: "dc2-production"
    kafka:
      bootstrap_servers: "dc2-broker1:9092,dc2-broker2:9092"
    shared_auth:
      username: "${DC2_USER}"
      password: "${DC2_PASS}"

  - name: "dc3-disaster-recovery"
    kafka:
      bootstrap_servers: "dc3-broker1:9092"
    overrides:
      disable_ksqldb: true
      disable_control_center: true

output:
  format: "json"
  file: "multi-dc-report.json"
  detailed: true
```

### Local Development

```yaml
clusters:
  - name: "docker-compose"
    kafka:
      bootstrap_servers: "localhost:9092"

    # All components auto-discovered at localhost with default ports
    # No authentication required

output:
  format: "json"
  file: "local-dev-report.json"
  detailed: false
```

---

## Complete Field Reference

```yaml
clusters:
  - name: string (required)

    kafka: (required)
      bootstrap_servers: string (required)
      security_protocol: string
      sasl_mechanism: string
      sasl_username: string
      sasl_password: string

    shared_auth:
      username: string
      password: string

    schema_registry:
      url: string
      basic_auth_username: string
      basic_auth_password: string
      bearer_token: string
      api_key: string
      api_key_header: string

    kafka_connect:
      url: string
      basic_auth_username: string
      basic_auth_password: string
      bearer_token: string
      api_key: string
      api_key_header: string

    ksqldb:
      url: string
      basic_auth_username: string
      basic_auth_password: string
      bearer_token: string
      api_key: string
      api_key_header: string

    rest_proxy:
      url: string
      basic_auth_username: string
      basic_auth_password: string
      bearer_token: string
      api_key: string
      api_key_header: string

    control_center:
      url: string
      basic_auth_username: string
      basic_auth_password: string
      bearer_token: string
      api_key: string
      api_key_header: string

    prometheus:
      url: string
      basic_auth_username: string
      basic_auth_password: string
      bearer_token: string
      api_key: string
      api_key_header: string

    alertmanager:
      url: string
      basic_auth_username: string
      basic_auth_password: string
      bearer_token: string
      api_key: string
      api_key_header: string

    overrides:
      disable_schema_registry: boolean
      disable_kafka_connect: boolean
      disable_ksqldb: boolean
      disable_rest_proxy: boolean
      disable_control_center: boolean
      disable_prometheus: boolean
      disable_alertmanager: boolean

output:
  format: string ("json" | "yaml")
  file: string
  detailed: boolean
```
