package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// SecretMapping maps an environment variable name to a Vault secret path.
type SecretMapping struct {
	EnvVar string `yaml:"env"`
	Path   string `yaml:"path"`
}

// Config holds the top-level vaultpipe configuration.
type Config struct {
	Secrets []SecretMapping `yaml:"secrets"`
}

// LoadFile reads and parses a YAML config file from the given path.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}
	return parse(data)
}

func parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing yaml: %w", err)
	}
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	for i, s := range cfg.Secrets {
		if strings.TrimSpace(s.EnvVar) == "" {
			return fmt.Errorf("config: secret[%d]: env must not be empty", i)
		}
		if strings.TrimSpace(s.Path) == "" {
			return fmt.Errorf("config: secret[%d]: path must not be empty", i)
		}
	}
	return nil
}
