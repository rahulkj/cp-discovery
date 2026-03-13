package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/rahulkj/cp-discovery/internal/model"
)

// ApplyDefaults applies smart defaults and conventions to cluster config
func ApplyDefaults(c *model.ClusterConfig) {
	// Extract base host from Kafka bootstrap servers if available
	baseHost := extractBaseHost(c.Kafka.BootstrapServers)

	// Apply defaults for Schema Registry
	if c.SchemaRegistry.URL == "" && baseHost != "" {
		c.SchemaRegistry.URL = deriveURL(baseHost, 8081, c.Kafka.SecurityProtocol)
	}
	applySharedAuth(&c.SchemaRegistry.BasicAuthUsername, &c.SchemaRegistry.BasicAuthPassword, c)

	// Apply defaults for Kafka Connect
	if c.KafkaConnect.URL == "" && baseHost != "" {
		c.KafkaConnect.URL = deriveURL(baseHost, 8083, c.Kafka.SecurityProtocol)
	}
	applySharedAuth(&c.KafkaConnect.BasicAuthUsername, &c.KafkaConnect.BasicAuthPassword, c)

	// Apply defaults for ksqlDB
	if c.KsqlDB.URL == "" && baseHost != "" {
		c.KsqlDB.URL = deriveURL(baseHost, 8088, c.Kafka.SecurityProtocol)
	}
	applySharedAuth(&c.KsqlDB.BasicAuthUsername, &c.KsqlDB.BasicAuthPassword, c)

	// Apply defaults for REST Proxy
	if c.RestProxy.URL == "" && baseHost != "" {
		c.RestProxy.URL = deriveURL(baseHost, 8082, c.Kafka.SecurityProtocol)
	}
	applySharedAuth(&c.RestProxy.BasicAuthUsername, &c.RestProxy.BasicAuthPassword, c)

	// Apply defaults for Control Center
	if c.ControlCenter.URL == "" && baseHost != "" {
		c.ControlCenter.URL = deriveURL(baseHost, 9021, c.Kafka.SecurityProtocol)
	}
	applySharedAuth(&c.ControlCenter.BasicAuthUsername, &c.ControlCenter.BasicAuthPassword, c)

	// Apply defaults for Prometheus
	if c.Prometheus.URL == "" && baseHost != "" {
		c.Prometheus.URL = deriveURL(baseHost, 9090, "PLAINTEXT")
	}
	applySharedAuth(&c.Prometheus.BasicAuthUsername, &c.Prometheus.BasicAuthPassword, c)

	// Apply defaults for Alertmanager
	if c.Alertmanager.URL == "" && baseHost != "" {
		c.Alertmanager.URL = deriveURL(baseHost, 9093, "PLAINTEXT")
	}
	applySharedAuth(&c.Alertmanager.BasicAuthUsername, &c.Alertmanager.BasicAuthPassword, c)

	// Expand environment variables in all fields
	expandEnvVars(c)
}

// extractBaseHost extracts the first host from bootstrap servers
func extractBaseHost(bootstrapServers string) string {
	if bootstrapServers == "" {
		return ""
	}

	// Split by comma and take first server
	servers := strings.Split(bootstrapServers, ",")
	if len(servers) == 0 {
		return ""
	}

	// Remove port if present
	hostPort := strings.TrimSpace(servers[0])
	parts := strings.Split(hostPort, ":")
	return parts[0]
}

// deriveURL creates a URL from host, port, and security protocol
func deriveURL(host string, port int, securityProtocol string) string {
	scheme := "http"
	if strings.Contains(strings.ToUpper(securityProtocol), "SSL") {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

// applySharedAuth applies shared authentication if component auth is not set
func applySharedAuth(username, password *string, config *model.ClusterConfig) {
	// If shared auth is configured and component auth is not, use shared
	if config.SharedAuth != nil {
		if *username == "" && config.SharedAuth.Username != "" {
			*username = config.SharedAuth.Username
		}
		if *password == "" && config.SharedAuth.Password != "" {
			*password = config.SharedAuth.Password
		}
	}
}

// expandEnvVars expands environment variables in configuration
func expandEnvVars(c *model.ClusterConfig) {
	// Expand shared auth
	if c.SharedAuth != nil {
		c.SharedAuth.Username = expandEnv(c.SharedAuth.Username)
		c.SharedAuth.Password = expandEnv(c.SharedAuth.Password)
	}

	// Expand Kafka config
	c.Kafka.BootstrapServers = expandEnv(c.Kafka.BootstrapServers)
	c.Kafka.SecurityProtocol = expandEnv(c.Kafka.SecurityProtocol)
	c.Kafka.SaslMechanism = expandEnv(c.Kafka.SaslMechanism)
	c.Kafka.SaslUsername = expandEnv(c.Kafka.SaslUsername)
	c.Kafka.SaslPassword = expandEnv(c.Kafka.SaslPassword)
	c.Kafka.SslCaLocation = expandEnv(c.Kafka.SslCaLocation)
	c.Kafka.SslCertLocation = expandEnv(c.Kafka.SslCertLocation)
	c.Kafka.SslKeyLocation = expandEnv(c.Kafka.SslKeyLocation)
	c.Kafka.SslKeyPassword = expandEnv(c.Kafka.SslKeyPassword)
	c.Kafka.SslEndpointIdentification = expandEnv(c.Kafka.SslEndpointIdentification)

	// Expand Schema Registry config
	c.SchemaRegistry.URL = expandEnv(c.SchemaRegistry.URL)
	c.SchemaRegistry.BasicAuthUsername = expandEnv(c.SchemaRegistry.BasicAuthUsername)
	c.SchemaRegistry.BasicAuthPassword = expandEnv(c.SchemaRegistry.BasicAuthPassword)
	c.SchemaRegistry.BearerToken = expandEnv(c.SchemaRegistry.BearerToken)
	c.SchemaRegistry.APIKey = expandEnv(c.SchemaRegistry.APIKey)
	c.SchemaRegistry.APIKeyHeader = expandEnv(c.SchemaRegistry.APIKeyHeader)
	c.SchemaRegistry.LDAPServer = expandEnv(c.SchemaRegistry.LDAPServer)
	c.SchemaRegistry.LDAPUsername = expandEnv(c.SchemaRegistry.LDAPUsername)
	c.SchemaRegistry.LDAPPassword = expandEnv(c.SchemaRegistry.LDAPPassword)
	c.SchemaRegistry.LDAPBaseDN = expandEnv(c.SchemaRegistry.LDAPBaseDN)
	c.SchemaRegistry.OAuthClientID = expandEnv(c.SchemaRegistry.OAuthClientID)
	c.SchemaRegistry.OAuthClientSecret = expandEnv(c.SchemaRegistry.OAuthClientSecret)
	c.SchemaRegistry.OAuthTokenURL = expandEnv(c.SchemaRegistry.OAuthTokenURL)
	c.SchemaRegistry.OAuthScopes = expandEnv(c.SchemaRegistry.OAuthScopes)

	// Expand Kafka Connect config
	c.KafkaConnect.URL = expandEnv(c.KafkaConnect.URL)
	c.KafkaConnect.BasicAuthUsername = expandEnv(c.KafkaConnect.BasicAuthUsername)
	c.KafkaConnect.BasicAuthPassword = expandEnv(c.KafkaConnect.BasicAuthPassword)
	c.KafkaConnect.BearerToken = expandEnv(c.KafkaConnect.BearerToken)
	c.KafkaConnect.APIKey = expandEnv(c.KafkaConnect.APIKey)
	c.KafkaConnect.APIKeyHeader = expandEnv(c.KafkaConnect.APIKeyHeader)
	c.KafkaConnect.LDAPServer = expandEnv(c.KafkaConnect.LDAPServer)
	c.KafkaConnect.LDAPUsername = expandEnv(c.KafkaConnect.LDAPUsername)
	c.KafkaConnect.LDAPPassword = expandEnv(c.KafkaConnect.LDAPPassword)
	c.KafkaConnect.LDAPBaseDN = expandEnv(c.KafkaConnect.LDAPBaseDN)
	c.KafkaConnect.OAuthClientID = expandEnv(c.KafkaConnect.OAuthClientID)
	c.KafkaConnect.OAuthClientSecret = expandEnv(c.KafkaConnect.OAuthClientSecret)
	c.KafkaConnect.OAuthTokenURL = expandEnv(c.KafkaConnect.OAuthTokenURL)
	c.KafkaConnect.OAuthScopes = expandEnv(c.KafkaConnect.OAuthScopes)

	// Expand ksqlDB config
	c.KsqlDB.URL = expandEnv(c.KsqlDB.URL)
	c.KsqlDB.BasicAuthUsername = expandEnv(c.KsqlDB.BasicAuthUsername)
	c.KsqlDB.BasicAuthPassword = expandEnv(c.KsqlDB.BasicAuthPassword)
	c.KsqlDB.BearerToken = expandEnv(c.KsqlDB.BearerToken)
	c.KsqlDB.APIKey = expandEnv(c.KsqlDB.APIKey)
	c.KsqlDB.APIKeyHeader = expandEnv(c.KsqlDB.APIKeyHeader)
	c.KsqlDB.LDAPServer = expandEnv(c.KsqlDB.LDAPServer)
	c.KsqlDB.LDAPUsername = expandEnv(c.KsqlDB.LDAPUsername)
	c.KsqlDB.LDAPPassword = expandEnv(c.KsqlDB.LDAPPassword)
	c.KsqlDB.LDAPBaseDN = expandEnv(c.KsqlDB.LDAPBaseDN)
	c.KsqlDB.OAuthClientID = expandEnv(c.KsqlDB.OAuthClientID)
	c.KsqlDB.OAuthClientSecret = expandEnv(c.KsqlDB.OAuthClientSecret)
	c.KsqlDB.OAuthTokenURL = expandEnv(c.KsqlDB.OAuthTokenURL)
	c.KsqlDB.OAuthScopes = expandEnv(c.KsqlDB.OAuthScopes)

	// Expand REST Proxy config
	c.RestProxy.URL = expandEnv(c.RestProxy.URL)
	c.RestProxy.BasicAuthUsername = expandEnv(c.RestProxy.BasicAuthUsername)
	c.RestProxy.BasicAuthPassword = expandEnv(c.RestProxy.BasicAuthPassword)
	c.RestProxy.BearerToken = expandEnv(c.RestProxy.BearerToken)
	c.RestProxy.APIKey = expandEnv(c.RestProxy.APIKey)
	c.RestProxy.APIKeyHeader = expandEnv(c.RestProxy.APIKeyHeader)
	c.RestProxy.LDAPServer = expandEnv(c.RestProxy.LDAPServer)
	c.RestProxy.LDAPUsername = expandEnv(c.RestProxy.LDAPUsername)
	c.RestProxy.LDAPPassword = expandEnv(c.RestProxy.LDAPPassword)
	c.RestProxy.LDAPBaseDN = expandEnv(c.RestProxy.LDAPBaseDN)
	c.RestProxy.OAuthClientID = expandEnv(c.RestProxy.OAuthClientID)
	c.RestProxy.OAuthClientSecret = expandEnv(c.RestProxy.OAuthClientSecret)
	c.RestProxy.OAuthTokenURL = expandEnv(c.RestProxy.OAuthTokenURL)
	c.RestProxy.OAuthScopes = expandEnv(c.RestProxy.OAuthScopes)

	// Expand Control Center config
	c.ControlCenter.URL = expandEnv(c.ControlCenter.URL)
	c.ControlCenter.BasicAuthUsername = expandEnv(c.ControlCenter.BasicAuthUsername)
	c.ControlCenter.BasicAuthPassword = expandEnv(c.ControlCenter.BasicAuthPassword)
	c.ControlCenter.BearerToken = expandEnv(c.ControlCenter.BearerToken)
	c.ControlCenter.APIKey = expandEnv(c.ControlCenter.APIKey)
	c.ControlCenter.APIKeyHeader = expandEnv(c.ControlCenter.APIKeyHeader)
	c.ControlCenter.LDAPServer = expandEnv(c.ControlCenter.LDAPServer)
	c.ControlCenter.LDAPUsername = expandEnv(c.ControlCenter.LDAPUsername)
	c.ControlCenter.LDAPPassword = expandEnv(c.ControlCenter.LDAPPassword)
	c.ControlCenter.LDAPBaseDN = expandEnv(c.ControlCenter.LDAPBaseDN)
	c.ControlCenter.OAuthClientID = expandEnv(c.ControlCenter.OAuthClientID)
	c.ControlCenter.OAuthClientSecret = expandEnv(c.ControlCenter.OAuthClientSecret)
	c.ControlCenter.OAuthTokenURL = expandEnv(c.ControlCenter.OAuthTokenURL)
	c.ControlCenter.OAuthScopes = expandEnv(c.ControlCenter.OAuthScopes)

	// Expand Prometheus config
	c.Prometheus.URL = expandEnv(c.Prometheus.URL)
	c.Prometheus.BasicAuthUsername = expandEnv(c.Prometheus.BasicAuthUsername)
	c.Prometheus.BasicAuthPassword = expandEnv(c.Prometheus.BasicAuthPassword)
	c.Prometheus.BearerToken = expandEnv(c.Prometheus.BearerToken)
	c.Prometheus.APIKey = expandEnv(c.Prometheus.APIKey)
	c.Prometheus.APIKeyHeader = expandEnv(c.Prometheus.APIKeyHeader)
	c.Prometheus.LDAPServer = expandEnv(c.Prometheus.LDAPServer)
	c.Prometheus.LDAPUsername = expandEnv(c.Prometheus.LDAPUsername)
	c.Prometheus.LDAPPassword = expandEnv(c.Prometheus.LDAPPassword)
	c.Prometheus.LDAPBaseDN = expandEnv(c.Prometheus.LDAPBaseDN)
	c.Prometheus.OAuthClientID = expandEnv(c.Prometheus.OAuthClientID)
	c.Prometheus.OAuthClientSecret = expandEnv(c.Prometheus.OAuthClientSecret)
	c.Prometheus.OAuthTokenURL = expandEnv(c.Prometheus.OAuthTokenURL)
	c.Prometheus.OAuthScopes = expandEnv(c.Prometheus.OAuthScopes)

	// Expand Alertmanager config
	c.Alertmanager.URL = expandEnv(c.Alertmanager.URL)
	c.Alertmanager.BasicAuthUsername = expandEnv(c.Alertmanager.BasicAuthUsername)
	c.Alertmanager.BasicAuthPassword = expandEnv(c.Alertmanager.BasicAuthPassword)
	c.Alertmanager.BearerToken = expandEnv(c.Alertmanager.BearerToken)
	c.Alertmanager.APIKey = expandEnv(c.Alertmanager.APIKey)
	c.Alertmanager.APIKeyHeader = expandEnv(c.Alertmanager.APIKeyHeader)
	c.Alertmanager.LDAPServer = expandEnv(c.Alertmanager.LDAPServer)
	c.Alertmanager.LDAPUsername = expandEnv(c.Alertmanager.LDAPUsername)
	c.Alertmanager.LDAPPassword = expandEnv(c.Alertmanager.LDAPPassword)
	c.Alertmanager.LDAPBaseDN = expandEnv(c.Alertmanager.LDAPBaseDN)
	c.Alertmanager.OAuthClientID = expandEnv(c.Alertmanager.OAuthClientID)
	c.Alertmanager.OAuthClientSecret = expandEnv(c.Alertmanager.OAuthClientSecret)
	c.Alertmanager.OAuthTokenURL = expandEnv(c.Alertmanager.OAuthTokenURL)
	c.Alertmanager.OAuthScopes = expandEnv(c.Alertmanager.OAuthScopes)
}

// expandEnv expands environment variables in format ${VAR} or $VAR
func expandEnv(s string) string {
	return os.ExpandEnv(s)
}

// Validate validates the configuration
func Validate(c *model.ClusterConfig) error {
	if c.Name == "" {
		return fmt.Errorf("cluster name is required")
	}

	if c.Kafka.BootstrapServers == "" {
		return fmt.Errorf("kafka bootstrap_servers is required for cluster %s", c.Name)
	}

	// Validate URLs if provided
	if c.SchemaRegistry.URL != "" {
		if _, err := url.Parse(c.SchemaRegistry.URL); err != nil {
			return fmt.Errorf("invalid schema registry URL for cluster %s: %w", c.Name, err)
		}
	}

	if c.KafkaConnect.URL != "" {
		if _, err := url.Parse(c.KafkaConnect.URL); err != nil {
			return fmt.Errorf("invalid kafka connect URL for cluster %s: %w", c.Name, err)
		}
	}

	if c.KsqlDB.URL != "" {
		if _, err := url.Parse(c.KsqlDB.URL); err != nil {
			return fmt.Errorf("invalid ksqlDB URL for cluster %s: %w", c.Name, err)
		}
	}

	if c.RestProxy.URL != "" {
		if _, err := url.Parse(c.RestProxy.URL); err != nil {
			return fmt.Errorf("invalid REST proxy URL for cluster %s: %w", c.Name, err)
		}
	}

	return nil
}

// ShouldDiscover checks if a component should be discovered
func ShouldDiscoverSchemaRegistry(c *model.ClusterConfig) bool {
	if c.Overrides != nil && c.Overrides.DisableSchemaRegistry {
		return false
	}
	return c.SchemaRegistry.URL != ""
}

func ShouldDiscoverKafkaConnect(c *model.ClusterConfig) bool {
	if c.Overrides != nil && c.Overrides.DisableKafkaConnect {
		return false
	}
	return c.KafkaConnect.URL != ""
}

func ShouldDiscoverKsqlDB(c *model.ClusterConfig) bool {
	if c.Overrides != nil && c.Overrides.DisableKsqlDB {
		return false
	}
	return c.KsqlDB.URL != ""
}

func ShouldDiscoverRestProxy(c *model.ClusterConfig) bool {
	if c.Overrides != nil && c.Overrides.DisableRestProxy {
		return false
	}
	return c.RestProxy.URL != ""
}

func ShouldDiscoverControlCenter(c *model.ClusterConfig) bool {
	if c.Overrides != nil && c.Overrides.DisableControlCenter {
		return false
	}
	return c.ControlCenter.URL != ""
}

func ShouldDiscoverPrometheus(c *model.ClusterConfig) bool {
	if c.Overrides != nil && c.Overrides.DisablePrometheus {
		return false
	}
	return c.Prometheus.URL != ""
}

func ShouldDiscoverAlertmanager(c *model.ClusterConfig) bool {
	if c.Overrides != nil && c.Overrides.DisableAlertmanager {
		return false
	}
	return c.Alertmanager.URL != ""
}
