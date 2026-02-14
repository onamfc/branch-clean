package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for branch-clean
type Config struct {
	StaleDays int      `yaml:"stale_days"`
	Protected []string `yaml:"protected"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		StaleDays: 30,
		Protected: []string{"main", "master", "develop", "release/*"},
	}
}

// LoadConfig loads configuration from a file
// If the file doesn't exist, returns the default configuration
func LoadConfig() (*Config, error) {
	// Try to load from ~/.branch-clean.yaml
	home, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), nil
	}

	configPath := filepath.Join(home, ".branch-clean.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file doesn't exist, use defaults
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing values
	if config.StaleDays == 0 {
		config.StaleDays = 30
	}
	if len(config.Protected) == 0 {
		config.Protected = []string{"main", "master", "develop", "release/*"}
	}

	return &config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".branch-clean.yaml")
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
