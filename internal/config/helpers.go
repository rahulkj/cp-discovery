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
	c.Kafka.BootstrapServers = expandEnv(c.Kafka.BootstrapServers)
	c.Kafka.SaslUsername = expandEnv(c.Kafka.SaslUsername)
	c.Kafka.SaslPassword = expandEnv(c.Kafka.SaslPassword)

	c.SchemaRegistry.URL = expandEnv(c.SchemaRegistry.URL)
	c.SchemaRegistry.BasicAuthUsername = expandEnv(c.SchemaRegistry.BasicAuthUsername)
	c.SchemaRegistry.BasicAuthPassword = expandEnv(c.SchemaRegistry.BasicAuthPassword)

	c.KafkaConnect.URL = expandEnv(c.KafkaConnect.URL)
	c.KafkaConnect.BasicAuthUsername = expandEnv(c.KafkaConnect.BasicAuthUsername)
	c.KafkaConnect.BasicAuthPassword = expandEnv(c.KafkaConnect.BasicAuthPassword)

	c.KsqlDB.URL = expandEnv(c.KsqlDB.URL)
	c.KsqlDB.BasicAuthUsername = expandEnv(c.KsqlDB.BasicAuthUsername)
	c.KsqlDB.BasicAuthPassword = expandEnv(c.KsqlDB.BasicAuthPassword)

	c.RestProxy.URL = expandEnv(c.RestProxy.URL)
	c.RestProxy.BasicAuthUsername = expandEnv(c.RestProxy.BasicAuthUsername)
	c.RestProxy.BasicAuthPassword = expandEnv(c.RestProxy.BasicAuthPassword)

	c.ControlCenter.URL = expandEnv(c.ControlCenter.URL)
	c.ControlCenter.BasicAuthUsername = expandEnv(c.ControlCenter.BasicAuthUsername)
	c.ControlCenter.BasicAuthPassword = expandEnv(c.ControlCenter.BasicAuthPassword)

	c.Prometheus.URL = expandEnv(c.Prometheus.URL)
	c.Prometheus.BasicAuthUsername = expandEnv(c.Prometheus.BasicAuthUsername)
	c.Prometheus.BasicAuthPassword = expandEnv(c.Prometheus.BasicAuthPassword)

	c.Alertmanager.URL = expandEnv(c.Alertmanager.URL)
	c.Alertmanager.BasicAuthUsername = expandEnv(c.Alertmanager.BasicAuthUsername)
	c.Alertmanager.BasicAuthPassword = expandEnv(c.Alertmanager.BasicAuthPassword)
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
