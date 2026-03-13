# Using Environment Variables

This guide explains how to use environment variables in your CP Discovery configuration files.

## Overview

All credential and configuration fields in CP Discovery support environment variable expansion. This allows you to:

- Keep sensitive credentials out of configuration files
- Use the same configuration across different environments
- Easily integrate with CI/CD pipelines and container orchestration
- Follow security best practices

## Environment Variable Syntax

CP Discovery supports two environment variable formats:

1. **`${VAR_NAME}`** - Expands to the value of `VAR_NAME` environment variable
2. **`${VAR_NAME:-default_value}`** - Expands to the value of `VAR_NAME`, or `default_value` if not set

### Examples

```yaml
# Simple variable expansion
sasl_username: "${KAFKA_USER}"

# With default value
security_protocol: "${KAFKA_SECURITY_PROTOCOL:-SASL_SSL}"

# Cluster name with default
name: "${CLUSTER_NAME:-production-cluster}"
```

## Supported Fields

Environment variable expansion is supported for **all** configuration fields, including:

### Kafka Configuration
- `bootstrap_servers`
- `security_protocol`
- `sasl_mechanism`
- `sasl_username`
- `sasl_password`
- `ssl_ca_location`
- `ssl_cert_location`
- `ssl_key_location`
- `ssl_key_password`
- `ssl_endpoint_identification`

### Component Authentication
For all components (Schema Registry, Kafka Connect, ksqlDB, REST Proxy, Control Center, Prometheus, Alertmanager):

- `url`
- `basic_auth_username`
- `basic_auth_password`
- `bearer_token`
- `api_key`
- `api_key_header`

### LDAP Configuration
- `ldap_enabled`
- `ldap_server`
- `ldap_username`
- `ldap_password`
- `ldap_base_dn`

### OAuth Configuration
- `oauth_enabled`
- `oauth_client_id`
- `oauth_client_secret`
- `oauth_token_url`
- `oauth_scopes`

### Output Configuration
- `format`
- `file`

## Usage Methods

### Method 1: Export Environment Variables

```bash
# Export variables
export KAFKA_BOOTSTRAP_SERVERS="broker1:9092,broker2:9092"
export KAFKA_USER="admin"
export KAFKA_PASS="secret-password"
export SR_USER="sr-admin"
export SR_PASS="sr-password"

# Run CP Discovery
cp-discovery --config configs/config-env-vars.yaml
```

### Method 2: Inline Environment Variables

```bash
KAFKA_USER=admin KAFKA_PASS=secret cp-discovery --config config.yaml
```

### Method 3: Using a .env File

```bash
# Create .env file
cat > .env << 'EOF'
KAFKA_BOOTSTRAP_SERVERS=broker1:9092
KAFKA_USER=admin
KAFKA_PASS=secret-password
EOF

# Source the file
source .env

# Run CP Discovery
cp-discovery --config configs/config-env-vars.yaml
```

### Method 4: Using direnv (Recommended for Development)

Install [direnv](https://direnv.net/) and create a `.envrc` file:

```bash
# Install direnv
# macOS: brew install direnv
# Linux: apt-get install direnv

# Create .envrc
cat > .envrc << 'EOF'
export KAFKA_BOOTSTRAP_SERVERS=broker1:9092
export KAFKA_USER=admin
export KAFKA_PASS=secret-password
EOF

# Allow direnv
direnv allow

# Variables are automatically loaded when you cd into the directory
cp-discovery --config configs/config-env-vars.yaml
```

### Method 5: Docker/Kubernetes

```yaml
# docker-compose.yml
services:
  cp-discovery:
    image: cp-discovery:latest
    environment:
      - KAFKA_BOOTSTRAP_SERVERS=broker1:9092
      - KAFKA_USER=admin
      - KAFKA_PASS=${KAFKA_PASS}  # From host environment
    volumes:
      - ./config.yaml:/config.yaml
```

```yaml
# Kubernetes Secret
apiVersion: v1
kind: Secret
metadata:
  name: cp-discovery-credentials
type: Opaque
stringData:
  kafka-user: admin
  kafka-pass: secret-password
---
# Kubernetes Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cp-discovery
spec:
  template:
    spec:
      containers:
      - name: cp-discovery
        image: cp-discovery:latest
        env:
        - name: KAFKA_USER
          valueFrom:
            secretKeyRef:
              name: cp-discovery-credentials
              key: kafka-user
        - name: KAFKA_PASS
          valueFrom:
            secretKeyRef:
              name: cp-discovery-credentials
              key: kafka-pass
```

## Example Configuration Files

### Basic Example

See [`configs/config-env-vars.yaml`](../configs/config-env-vars.yaml) for a comprehensive example using environment variables.

### Minimal Example

```yaml
clusters:
  - name: "${CLUSTER_NAME:-my-cluster}"
    kafka:
      bootstrap_servers: "${KAFKA_BOOTSTRAP_SERVERS}"
      security_protocol: "${KAFKA_SECURITY_PROTOCOL:-SASL_SSL}"
      sasl_mechanism: "${KAFKA_SASL_MECHANISM:-PLAIN}"
      sasl_username: "${KAFKA_USER}"
      sasl_password: "${KAFKA_PASS}"

    shared_auth:
      username: "${CP_USER}"
      password: "${CP_PASS}"

output:
  format: "json"
  file: "${OUTPUT_FILE:-discovery-report.json}"
```

## Best Practices

### 1. Use Specific Variable Names

```yaml
# Good - specific and clear
sasl_username: "${KAFKA_SASL_USERNAME}"
sasl_password: "${KAFKA_SASL_PASSWORD}"

# Avoid - too generic
sasl_username: "${USER}"
sasl_password: "${PASS}"
```

### 2. Provide Defaults for Non-Sensitive Values

```yaml
# Good - defaults for protocol and mechanism
security_protocol: "${KAFKA_SECURITY_PROTOCOL:-SASL_SSL}"
sasl_mechanism: "${KAFKA_SASL_MECHANISM:-PLAIN}"

# Don't provide defaults for credentials
sasl_username: "${KAFKA_USER}"  # No default
sasl_password: "${KAFKA_PASS}"  # No default
```

### 3. Use Shared Auth for Consistency

```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "${KAFKA_BOOTSTRAP_SERVERS}"

    # Define once
    shared_auth:
      username: "${CP_USER}"
      password: "${CP_PASS}"

    # Automatically applied to all components
    schema_registry:
      url: "${SR_URL}"
    kafka_connect:
      url: "${CONNECT_URL}"
```

### 4. Keep .env Files Out of Version Control

```bash
# Add to .gitignore
echo ".env" >> .gitignore
echo ".envrc" >> .gitignore
```

### 5. Use Environment-Specific Configurations

```bash
# dev.env
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_SECURITY_PROTOCOL=PLAINTEXT

# prod.env
KAFKA_BOOTSTRAP_SERVERS=prod-broker1:9092,prod-broker2:9092
KAFKA_SECURITY_PROTOCOL=SASL_SSL
```

## Security Considerations

### 1. Never Commit Credentials

- Use `.env.example` with placeholder values
- Keep actual `.env` files in `.gitignore`
- Use secret management tools in production

### 2. Use Secret Management

For production deployments, use proper secret management:

- **Kubernetes**: Use Secrets and External Secrets Operator
- **Docker Swarm**: Use Docker Secrets
- **Cloud**: Use AWS Secrets Manager, Azure Key Vault, GCP Secret Manager
- **HashiCorp Vault**: Integration via environment variables

### 3. Limit Permissions

```bash
# Restrict .env file permissions
chmod 600 .env

# Only readable by owner
ls -la .env
# -rw------- 1 user user 1234 Jan 01 12:00 .env
```

### 4. Audit Environment Variables

```bash
# Check what variables are set (don't log values!)
env | grep -E "KAFKA|SR_|CONNECT_" | cut -d= -f1

# Verify variables are set before running
[ -z "$KAFKA_USER" ] && echo "ERROR: KAFKA_USER not set" && exit 1
```

## Troubleshooting

### Variable Not Expanding

**Problem**: Variable appears as `${VAR_NAME}` in logs

**Solution**: Ensure the variable is exported:
```bash
# Wrong
VAR_NAME=value

# Correct
export VAR_NAME=value
```

### Empty Values

**Problem**: Configuration field is empty

**Solution**: Check if the variable is set:
```bash
echo "KAFKA_USER: ${KAFKA_USER}"
echo "KAFKA_PASS: ${KAFKA_PASS}"
```

### Special Characters in Values

**Problem**: Password with special characters not working

**Solution**: Quote the value when exporting:
```bash
# Correct
export KAFKA_PASS='p@ssw0rd!#$'

# Also correct
export KAFKA_PASS="p@ssw0rd!#$"
```

## Examples by Use Case

### Local Development

```bash
# Simple local setup
export KAFKA_BOOTSTRAP_SERVERS="localhost:9092"
export KAFKA_SECURITY_PROTOCOL="PLAINTEXT"
cp-discovery --config configs/config-env-vars.yaml
```

### CI/CD Pipeline

```yaml
# GitHub Actions
- name: Run CP Discovery
  env:
    KAFKA_BOOTSTRAP_SERVERS: ${{ secrets.KAFKA_BOOTSTRAP_SERVERS }}
    KAFKA_USER: ${{ secrets.KAFKA_USER }}
    KAFKA_PASS: ${{ secrets.KAFKA_PASS }}
  run: |
    cp-discovery --config config.yaml
```

### Confluent Cloud

```bash
export CCLOUD_BOOTSTRAP_SERVERS="pkc-xxxxx.us-east-1.aws.confluent.cloud:9092"
export CCLOUD_API_KEY="your-api-key"
export CCLOUD_API_SECRET="your-api-secret"
export CCLOUD_SR_URL="https://psrc-xxxxx.us-east-2.aws.confluent.cloud"
export CCLOUD_SR_KEY="your-sr-key"
export CCLOUD_SR_SECRET="your-sr-secret"

cp-discovery --config configs/config-env-vars.yaml
```

## See Also

- [Configuration Guide](CONFIG_GUIDE.md)
- [Authentication Guide](AUTHENTICATION.md)
- [Example Configurations](../configs/)
- [.env.example](../.env.example)
