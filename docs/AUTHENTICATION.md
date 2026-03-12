# Authentication Guide

This guide covers all supported authentication mechanisms in the CP Discovery tool.

## Overview

The tool supports **five authentication methods** with automatic priority-based selection:

```
Priority Order: OAuth/SSO > LDAP > Bearer Token > API Key > Basic Auth
```

## Authentication Methods

### 1. OAuth/SSO Authentication (Highest Priority)

OAuth 2.0 Client Credentials flow with automatic token management.

**Features:**
- Automatic token retrieval from OAuth provider
- Token caching to minimize authentication requests
- Automatic token refresh before expiration
- Compatible with Keycloak, Okta, Auth0, Azure AD, and other OAuth 2.0 providers

**Configuration:**

```yaml
schema_registry:
  url: "https://sr.example.com:8081"
  oauth_enabled: true
  oauth_client_id: "${SR_OAUTH_CLIENT_ID}"
  oauth_client_secret: "${SR_OAUTH_CLIENT_SECRET}"
  oauth_token_url: "https://auth.example.com/oauth/token"
  oauth_scopes: "schema-registry.read schema-registry.write"
```

**OAuth Token URL Examples:**

| Provider | Token URL Format |
|----------|-----------------|
| Keycloak | `https://keycloak.example.com/realms/{realm}/protocol/openid-connect/token` |
| Okta | `https://{domain}.okta.com/oauth2/default/v1/token` |
| Auth0 | `https://{domain}.auth0.com/oauth/token` |
| Azure AD | `https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token` |

**Environment Variables:**

```bash
export SR_OAUTH_CLIENT_ID="my-client-id"
export SR_OAUTH_CLIENT_SECRET="my-client-secret"
```

### 2. LDAP Authentication

Enterprise directory authentication with Active Directory or OpenLDAP.

**Features:**
- Support for both LDAP (port 389) and LDAPS (port 636)
- Falls back to Basic Auth with LDAP credentials
- Configurable base DN for user lookup
- Compatible with Active Directory, OpenLDAP, and other LDAP servers

**Configuration:**

```yaml
kafka_connect:
  url: "https://connect.example.com:8083"
  ldap_enabled: true
  ldap_server: "ldaps://ldap.example.com:636"  # LDAP over SSL
  ldap_username: "${LDAP_USER}"
  ldap_password: "${LDAP_PASS}"
  ldap_base_dn: "ou=users,dc=example,dc=com"
```

**LDAP Server URL Formats:**

| Protocol | Port | URL Format | Use Case |
|----------|------|------------|----------|
| LDAP | 389 | `ldap://ldap.example.com:389` | Non-encrypted (dev/test) |
| LDAPS | 636 | `ldaps://ldap.example.com:636` | SSL/TLS encrypted (production) |

**Active Directory Example:**

```yaml
control_center:
  url: "https://c3.example.com:9021"
  ldap_enabled: true
  ldap_server: "ldaps://ad.corp.example.com:636"
  ldap_username: "${AD_USERNAME}"  # e.g., "john.doe" or "john.doe@corp.example.com"
  ldap_password: "${AD_PASSWORD}"
  ldap_base_dn: "dc=corp,dc=example,dc=com"
```

**Environment Variables:**

```bash
export LDAP_USER="john.doe"
export LDAP_PASS="my-ldap-password"
```

### 3. Bearer Token (Pre-configured)

For pre-obtained OAuth tokens, JWT tokens, or other bearer tokens.

**Configuration:**

```yaml
kafka_connect:
  url: "https://connect.example.com:8083"
  bearer_token: "${CONNECT_JWT_TOKEN}"
```

**Use Cases:**
- Pre-obtained JWT tokens from external authentication
- Service account tokens
- Personal access tokens (PATs)
- Session tokens

**Environment Variables:**

```bash
export CONNECT_JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 4. API Key Authentication

Custom API key authentication with configurable headers.

**Configuration:**

```yaml
ksqldb:
  url: "https://ksqldb.example.com:8088"
  api_key: "${KSQLDB_API_KEY}"
  api_key_header: "X-API-Key"  # Optional, defaults to "X-API-Key"
```

**Common Header Names:**
- `X-API-Key` (default)
- `X-Api-Key`
- `Authorization` (some systems use this for API keys)
- `X-Custom-Auth`

**Environment Variables:**

```bash
export KSQLDB_API_KEY="my-api-key-12345"
```

### 5. Basic Authentication

Standard HTTP Basic Authentication with username and password.

**Configuration:**

```yaml
schema_registry:
  url: "https://sr.example.com:8081"
  basic_auth_username: "${SR_USER}"
  basic_auth_password: "${SR_PASS}"
```

**Environment Variables:**

```bash
export SR_USER="admin"
export SR_PASS="admin-secret"
```

## Authentication Priority

When multiple authentication methods are configured, the tool uses this priority order:

1. **OAuth/SSO** - If `oauth_enabled: true` and credentials are provided
2. **LDAP** - If `ldap_enabled: true` and credentials are provided
3. **Bearer Token** - If `bearer_token` is provided
4. **API Key** - If `api_key` is provided
5. **Basic Auth** - If `basic_auth_username` is provided

**Example:** If both OAuth and Basic Auth are configured, OAuth will be used:

```yaml
schema_registry:
  url: "https://sr.example.com:8081"
  # OAuth will be used (highest priority)
  oauth_enabled: true
  oauth_client_id: "${OAUTH_CLIENT_ID}"
  oauth_client_secret: "${OAUTH_CLIENT_SECRET}"
  oauth_token_url: "https://auth.example.com/oauth/token"
  # Basic Auth is ignored (lower priority)
  basic_auth_username: "admin"
  basic_auth_password: "secret"
```

## Shared Authentication

Apply the same credentials to all components using `shared_auth`:

```yaml
clusters:
  - name: "my-cluster"
    kafka:
      bootstrap_servers: "broker:9092"

    # Shared auth applies to all components
    shared_auth:
      username: "${API_USER}"
      password: "${API_PASSWORD}"

    # Components use shared_auth unless they specify their own
    schema_registry:
      url: "https://sr.example.com:8081"
      # Will use shared_auth

    kafka_connect:
      url: "https://connect.example.com:8083"
      # Override with OAuth for this component only
      oauth_enabled: true
      oauth_client_id: "${CONNECT_OAUTH_CLIENT_ID}"
      oauth_client_secret: "${CONNECT_OAUTH_CLIENT_SECRET}"
      oauth_token_url: "https://auth.example.com/oauth/token"
```

## Complete Examples

### Example 1: Full OAuth/SSO Setup

```yaml
clusters:
  - name: "oauth-cluster"
    kafka:
      bootstrap_servers: "broker:9092"
      security_protocol: "SASL_SSL"
      sasl_mechanism: "OAUTHBEARER"
      sasl_username: "${KAFKA_USER}"
      sasl_password: "${KAFKA_PASS}"

    schema_registry:
      url: "https://sr.example.com:8081"
      oauth_enabled: true
      oauth_client_id: "${SR_OAUTH_CLIENT_ID}"
      oauth_client_secret: "${SR_OAUTH_CLIENT_SECRET}"
      oauth_token_url: "https://keycloak.example.com/realms/confluent/protocol/openid-connect/token"
      oauth_scopes: "schema-registry.read schema-registry.write"

    kafka_connect:
      url: "https://connect.example.com:8083"
      oauth_enabled: true
      oauth_client_id: "${CONNECT_OAUTH_CLIENT_ID}"
      oauth_client_secret: "${CONNECT_OAUTH_CLIENT_SECRET}"
      oauth_token_url: "https://keycloak.example.com/realms/confluent/protocol/openid-connect/token"
      oauth_scopes: "connect.admin"
```

### Example 2: Full LDAP Setup

```yaml
clusters:
  - name: "ldap-cluster"
    kafka:
      bootstrap_servers: "broker:9092"

    schema_registry:
      url: "https://sr.example.com:8081"
      ldap_enabled: true
      ldap_server: "ldaps://ldap.corp.example.com:636"
      ldap_username: "${LDAP_USER}"
      ldap_password: "${LDAP_PASS}"
      ldap_base_dn: "ou=users,dc=corp,dc=example,dc=com"

    kafka_connect:
      url: "https://connect.example.com:8083"
      ldap_enabled: true
      ldap_server: "ldaps://ldap.corp.example.com:636"
      ldap_username: "${LDAP_USER}"
      ldap_password: "${LDAP_PASS}"
      ldap_base_dn: "ou=users,dc=corp,dc=example,dc=com"
```

### Example 3: Hybrid Authentication with Fallback

Use OAuth as primary, LDAP as fallback:

```yaml
clusters:
  - name: "hybrid-cluster"
    kafka:
      bootstrap_servers: "broker:9092"

    schema_registry:
      url: "https://sr.example.com:8081"
      # Primary: OAuth
      oauth_enabled: true
      oauth_client_id: "${SR_OAUTH_CLIENT_ID}"
      oauth_client_secret: "${SR_OAUTH_CLIENT_SECRET}"
      oauth_token_url: "https://auth.example.com/oauth/token"
      oauth_scopes: "schema-registry.read"
      # Fallback: LDAP (used if OAuth fails)
      ldap_enabled: true
      ldap_server: "ldaps://ldap.example.com:636"
      ldap_username: "${LDAP_USER}"
      ldap_password: "${LDAP_PASS}"
      ldap_base_dn: "dc=example,dc=com"
```

### Example 4: Mixed Authentication per Component

```yaml
clusters:
  - name: "mixed-auth-cluster"
    kafka:
      bootstrap_servers: "broker:9092"

    # Schema Registry with OAuth
    schema_registry:
      url: "https://sr.example.com:8081"
      oauth_enabled: true
      oauth_client_id: "${SR_OAUTH_CLIENT_ID}"
      oauth_client_secret: "${SR_OAUTH_CLIENT_SECRET}"
      oauth_token_url: "https://auth.example.com/oauth/token"

    # Kafka Connect with LDAP
    kafka_connect:
      url: "https://connect.example.com:8083"
      ldap_enabled: true
      ldap_server: "ldaps://ldap.example.com:636"
      ldap_username: "${LDAP_USER}"
      ldap_password: "${LDAP_PASS}"
      ldap_base_dn: "dc=example,dc=com"

    # ksqlDB with Bearer Token
    ksqldb:
      url: "https://ksqldb.example.com:8088"
      bearer_token: "${KSQLDB_JWT_TOKEN}"

    # Control Center with API Key
    control_center:
      url: "https://c3.example.com:9021"
      api_key: "${C3_API_KEY}"

    # Prometheus with Basic Auth
    prometheus:
      url: "http://prometheus.example.com:9090"
      basic_auth_username: "admin"
      basic_auth_password: "${PROM_PASS}"
```

## OAuth Token Caching

The tool automatically caches OAuth tokens to minimize authentication requests:

- Tokens are cached in memory per unique client ID + token URL + scopes combination
- Cached tokens are reused if they have more than 5 minutes until expiration
- Automatic token refresh when cache expires
- Thread-safe caching with concurrent access protection

**Cache Key Format:**
```
{client_id}:{token_url}:{scopes}
```

## Security Best Practices

### 1. Use Environment Variables

Never hardcode credentials in configuration files:

```yaml
# ❌ BAD - Credentials in config file
schema_registry:
  oauth_client_id: "my-client-id"
  oauth_client_secret: "super-secret-key"

# ✅ GOOD - Use environment variables
schema_registry:
  oauth_client_id: "${SR_OAUTH_CLIENT_ID}"
  oauth_client_secret: "${SR_OAUTH_CLIENT_SECRET}"
```

### 2. Use LDAPS Instead of LDAP

Always use encrypted LDAPS in production:

```yaml
# ❌ BAD - Unencrypted LDAP
ldap_server: "ldap://ldap.example.com:389"

# ✅ GOOD - Encrypted LDAPS
ldap_server: "ldaps://ldap.example.com:636"
```

### 3. Limit OAuth Scopes

Request only the scopes you need:

```yaml
# ❌ BAD - Requesting admin scope when not needed
oauth_scopes: "admin full-access"

# ✅ GOOD - Request minimal scopes
oauth_scopes: "schema-registry.read"
```

### 4. Rotate Credentials Regularly

- Rotate OAuth client secrets regularly
- Update LDAP passwords according to your organization's policy
- Revoke and regenerate API keys periodically

### 5. Use Separate Credentials per Environment

Don't reuse the same credentials across dev/staging/production:

```bash
# Development
export SR_OAUTH_CLIENT_ID="dev-client-id"

# Staging
export SR_OAUTH_CLIENT_ID="staging-client-id"

# Production
export SR_OAUTH_CLIENT_ID="prod-client-id"
```

## Troubleshooting

### OAuth Authentication Failures

**Problem:** `OAuth token request failed with status 401`

**Solutions:**
1. Verify client ID and secret are correct
2. Check token URL is accessible
3. Ensure scopes are valid for the OAuth provider
4. Verify network connectivity to auth server

**Debug Steps:**
```bash
# Test OAuth endpoint manually
curl -X POST "https://auth.example.com/oauth/token" \
  -d "grant_type=client_credentials" \
  -d "client_id=${OAUTH_CLIENT_ID}" \
  -d "client_secret=${OAUTH_CLIENT_SECRET}" \
  -d "scope=schema-registry.read"
```

### LDAP Authentication Failures

**Problem:** `LDAP bind failed`

**Solutions:**
1. Verify LDAP server URL and port
2. Check username format (may need full DN or UPN)
3. Ensure base DN is correct
4. Verify LDAPS certificate if using SSL

**Debug Steps:**
```bash
# Test LDAP connectivity
ldapsearch -H ldaps://ldap.example.com:636 \
  -D "uid=username,ou=users,dc=example,dc=com" \
  -W -b "dc=example,dc=com"
```

### Token Expiration Issues

**Problem:** Requests fail after some time

**Solutions:**
1. OAuth tokens are automatically refreshed - check logs for refresh errors
2. Verify token expiration time from OAuth provider
3. Check system clock synchronization (important for JWT validation)

### Mixed Authentication Not Working

**Problem:** Wrong authentication method being used

**Solutions:**
1. Check authentication priority order
2. Verify boolean flags (`oauth_enabled`, `ldap_enabled`)
3. Review logs to see which auth method was selected

## Additional Resources

- [OAuth 2.0 Client Credentials Flow](https://oauth.net/2/grant-types/client-credentials/)
- [LDAP Authentication Overview](https://ldap.com/basic-ldap-concepts/)
- [Configuration Examples](../configs/config-auth-examples.yaml)
- [Complete Configuration Reference](CONFIG_REFERENCE.md)
