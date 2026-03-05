package config

import (
	"fmt"
	"os"

	"github.com/rahulkj/cp-discovery/internal/model"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from a YAML file
func LoadConfig(filename string) (*model.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config model.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// Apply defaults and validate each cluster
	for i := range config.Clusters {
		ApplyDefaults(&config.Clusters[i])
		if err := Validate(&config.Clusters[i]); err != nil {
			return nil, fmt.Errorf("validating cluster config: %w", err)
		}
	}

	return &config, nil
}
