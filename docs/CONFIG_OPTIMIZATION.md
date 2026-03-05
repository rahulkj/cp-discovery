# Configuration Optimization Guide

## Overview

The configuration has been optimized to reduce the number of required inputs from **50+ fields** to as few as **2 fields** for simple setups!

## What Changed?

### ✨ New Features

1. **Smart URL Auto-Discovery**
   - Automatically derives component URLs from Kafka broker host
   - Uses standard Confluent Platform ports (8081, 8083, 8088, 8082)
   - Detects HTTP vs HTTPS based on Kafka security protocol

2. **Shared Authentication**
   - Single `shared_auth` section applies to all components
   - Override per-component only when needed
   - Reduces credential duplication

3. **Environment Variable Support**
   - Use `${VAR_NAME}` syntax anywhere in config
   - Keep secrets out of config files
   - CI/CD friendly

4. **Component Overrides**
   - Explicitly disable components you don't have
   - Faster discovery by skipping unavailable services
   - Cleaner error reporting

5. **Smart Defaults**
   - Only `bootstrap_servers` is required
   - Everything else has sensible defaults
   - `omitempty` tags reduce YAML verbosity

## Before vs After

### Before (Old Format)
```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "localhost:9092"
      security_protocol: "PLAINTEXT"
      sasl_mechanism: ""
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
```
**Lines: 23 | Fields: 17**

### After (New Format)
```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "localhost:9092"
```
**Lines: 4 | Fields: 2**

That's **82% fewer lines** and **88% fewer fields**!

## Configuration Levels

### Level 1: Absolute Minimum (Local Development)

```yaml
clusters:
  - name: "local"
    kafka:
      bootstrap_servers: "localhost:9092"
```

**What happens:**
- ✅ Discovers Kafka at localhost:9092
- ✅ Auto-discovers Schema Registry at http://localhost:8081
- ✅ Auto-discovers Kafka Connect at http://localhost:8083
- ✅ Auto-discovers ksqlDB at http://localhost:8088
- ✅ Auto-discovers REST Proxy at http://localhost:8082
- ✅ No authentication assumed (PLAINTEXT)

### Level 2: With Security (Typical Production)

```yaml
clusters:
  - name: "prod"
    kafka:
      bootstrap_servers: "broker.example.com:9093"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "PLAIN"
      sasl_username: "kafka-user"
      sasl_password: "kafka-pass"
    shared_auth:
      username: "cp-admin"
      password: "cp-secret"
```

**What happens:**
- ✅ Connects to Kafka with SASL_SSL
- ✅ Auto-discovers all components at `https://broker.example.com:XXXX`
- ✅ Uses `shared_auth` for all REST components
- ✅ Proper security for production

**Lines: 10 | Fields: 6**

### Level 3: Custom URLs + Overrides (Complex Production)

```yaml
clusters:
  - name: "prod"
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
      url: "https://sr.example.com:8081"

    kafka_connect:
      url: "https://connect.example.com:8083"
      basic_auth_username: "connect-specific-user"
      basic_auth_password: "connect-specific-pass"

    overrides:
      disable_ksqldb: true
      disable_rest_proxy: true
```

**What happens:**
- ✅ Uses environment variables for secrets
- ✅ Custom URLs for specific components
- ✅ Per-component auth override for Kafka Connect
- ✅ Disables ksqlDB and REST Proxy discovery
- ✅ Production-ready security

**Lines: 24 | Fields: 12**

## Auto-Discovery Rules

### URL Generation

Based on the first broker in `bootstrap_servers`:

| Component        | Default Port | HTTP/HTTPS             |
|------------------|--------------|------------------------|
| Schema Registry  | 8081         | Based on Kafka protocol|
| Kafka Connect    | 8083         | Based on Kafka protocol|
| ksqlDB           | 8088         | Based on Kafka protocol|
| REST Proxy       | 8082         | Based on Kafka protocol|

**Protocol Detection:**
- If Kafka uses `PLAINTEXT` or `SASL_PLAINTEXT` → HTTP
- If Kafka uses `SSL` or `SASL_SSL` → HTTPS

**Example:**
```yaml
kafka:
  bootstrap_servers: "broker1.prod.com:9093,broker2.prod.com:9093"
  security_protocol: "SASL_SSL"
```

Auto-discovered URLs:
- Schema Registry: `https://broker1.prod.com:8081`
- Kafka Connect: `https://broker1.prod.com:8083`
- ksqlDB: `https://broker1.prod.com:8088`
- REST Proxy: `https://broker1.prod.com:8082`

### Authentication Inheritance

```yaml
shared_auth:
  username: "admin"
  password: "secret"
```

Applied to all components **unless** overridden:

```yaml
kafka_connect:
  basic_auth_username: "special-user"  # Overrides shared_auth for Connect only
  basic_auth_password: "special-pass"
```

## Environment Variables

Use `${VAR_NAME}` or `$VAR_NAME` syntax:

```yaml
clusters:
  - name: "prod"
    kafka:
      bootstrap_servers: "${KAFKA_BOOTSTRAP_SERVERS}"
      sasl_username: "${KAFKA_USERNAME}"
      sasl_password: "${KAFKA_PASSWORD}"
    shared_auth:
      username: "${CP_USERNAME}"
      password: "${CP_PASSWORD}"
```

Set environment variables:
```bash
export KAFKA_BOOTSTRAP_SERVERS="broker:9092"
export KAFKA_USERNAME="admin"
export KAFKA_PASSWORD="secret"
export CP_USERNAME="cp-admin"
export CP_PASSWORD="cp-secret"

./cp-discovery
```

## Component Overrides

Explicitly disable components to:
- Skip discovery for unavailable services
- Reduce discovery time
- Avoid error noise

```yaml
overrides:
  disable_schema_registry: true   # Don't discover Schema Registry
  disable_kafka_connect: true     # Don't discover Kafka Connect
  disable_ksqldb: true            # Don't discover ksqlDB
  disable_rest_proxy: true        # Don't discover REST Proxy
```

## Migration Guide

### Old Config → New Config

**Old:**
```yaml
clusters:
  - name: "cluster1"
    kafka:
      bootstrap_servers: "localhost:9092"
      security_protocol: "PLAINTEXT"
      sasl_mechanism: ""
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
```

**New:**
```yaml
clusters:
  - name: "cluster1"
    kafka:
      bootstrap_servers: "localhost:9092"
```

### Secure Cluster Migration

**Old:**
```yaml
clusters:
  - name: "prod"
    kafka:
      bootstrap_servers: "broker:9092"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "PLAIN"
      sasl_username: "kafka-user"
      sasl_password: "kafka-pass"
    schema_registry:
      url: "https://broker:8081"
      basic_auth_username: "admin"
      basic_auth_password: "secret"
    kafka_connect:
      url: "https://broker:8083"
      basic_auth_username: "admin"
      basic_auth_password: "secret"
    # ... (repeated for ksqldb and rest_proxy)
```

**New:**
```yaml
clusters:
  - name: "prod"
    kafka:
      bootstrap_servers: "broker:9092"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "PLAIN"
      sasl_username: "kafka-user"
      sasl_password: "kafka-pass"
    shared_auth:
      username: "admin"
      password: "secret"
```

## Benefits Summary

| Feature                  | Lines Saved | Complexity Reduction |
|--------------------------|-------------|----------------------|
| Auto-discovery URLs      | ~15 lines   | 80%                  |
| Shared authentication    | ~12 lines   | 75%                  |
| Environment variables    | Variable    | Infinite (secrets)   |
| Smart defaults           | ~8 lines    | 60%                  |
| Component overrides      | N/A         | Clearer intent       |

**Total Reduction:** Up to **88% fewer fields** for typical configurations!

## Backward Compatibility

✅ The old format **still works**! You can:
- Keep existing configs unchanged
- Mix old and new styles
- Gradually migrate at your own pace

The new optimizations are **additive**, not breaking changes.

## Examples Repository

See the included example files:
- `config-minimal.yaml` - Simplest possible configurations
- `config-advanced.yaml` - Advanced features and overrides
- `example-local.yaml` - Local Docker Compose setup
- `config.yaml` - Recommended balanced approach
