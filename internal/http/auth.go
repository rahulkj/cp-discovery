package http

import (
	"net/http"

	"github.com/rahulkj/cp-discovery/internal/model"
)

// AuthConfig represents authentication configuration for HTTP clients
type AuthConfig interface {
	ApplyAuth(req *http.Request)
}

// ApplySchemaRegistryAuth applies authentication to Schema Registry HTTP requests
func ApplySchemaRegistryAuth(req *http.Request, config model.SchemaRegistryConfig) {
	applyAuth(req, config.BasicAuthUsername, config.BasicAuthPassword, config.BearerToken, config.APIKey, config.APIKeyHeader)
}

// ApplyKafkaConnectAuth applies authentication to Kafka Connect HTTP requests
func ApplyKafkaConnectAuth(req *http.Request, config model.KafkaConnectConfig) {
	applyAuth(req, config.BasicAuthUsername, config.BasicAuthPassword, config.BearerToken, config.APIKey, config.APIKeyHeader)
}

// ApplyKsqlDBAuth applies authentication to ksqlDB HTTP requests
func ApplyKsqlDBAuth(req *http.Request, config model.KsqlDBConfig) {
	applyAuth(req, config.BasicAuthUsername, config.BasicAuthPassword, config.BearerToken, config.APIKey, config.APIKeyHeader)
}

// ApplyRestProxyAuth applies authentication to REST Proxy HTTP requests
func ApplyRestProxyAuth(req *http.Request, config model.RestProxyConfig) {
	applyAuth(req, config.BasicAuthUsername, config.BasicAuthPassword, config.BearerToken, config.APIKey, config.APIKeyHeader)
}

// ApplyControlCenterAuth applies authentication to Control Center HTTP requests
func ApplyControlCenterAuth(req *http.Request, config model.ControlCenterConfig) {
	applyAuth(req, config.BasicAuthUsername, config.BasicAuthPassword, config.BearerToken, config.APIKey, config.APIKeyHeader)
}

// ApplyPrometheusAuth applies authentication to Prometheus HTTP requests
func ApplyPrometheusAuth(req *http.Request, config model.PrometheusConfig) {
	applyAuth(req, config.BasicAuthUsername, config.BasicAuthPassword, config.BearerToken, config.APIKey, config.APIKeyHeader)
}

// ApplyAlertmanagerAuth applies authentication to Alertmanager HTTP requests
func ApplyAlertmanagerAuth(req *http.Request, config model.AlertmanagerConfig) {
	applyAuth(req, config.BasicAuthUsername, config.BasicAuthPassword, config.BearerToken, config.APIKey, config.APIKeyHeader)
}

// applyAuth applies authentication headers to an HTTP request based on configured auth type
func applyAuth(req *http.Request, username, password, bearerToken, apiKey, apiKeyHeader string) {
	// Priority: Bearer Token > API Key > Basic Auth

	// Bearer Token Authentication (OAuth, JWT)
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
		return
	}

	// API Key Authentication
	if apiKey != "" {
		headerName := "X-API-Key" // Default header
		if apiKeyHeader != "" {
			headerName = apiKeyHeader
		}
		req.Header.Set(headerName, apiKey)
		return
	}

	// Basic Authentication
	if username != "" {
		req.SetBasicAuth(username, password)
		return
	}

	// No authentication configured
}
