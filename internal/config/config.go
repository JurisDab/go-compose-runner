package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Service describes one compose-managed service that devtool orchestrates.
type Service struct {
	Name                 string   `yaml:"name"`
	WorkDir              string   `yaml:"workDir"`
	HealthURL            string   `yaml:"healthUrl,omitempty"`
	HealthTimeoutSeconds int      `yaml:"healthTimeoutSeconds,omitempty"`
	Port                 int      `yaml:"port,omitempty"`
	StartAfter           string   `yaml:"startAfter,omitempty"`
	RequiredEnv          []string `yaml:"requiredEnv,omitempty"`
}

// DefaultHealthTimeoutSeconds is used when a service defines a HealthURL
// but no explicit HealthTimeoutSeconds.
const DefaultHealthTimeoutSeconds = 30

// Config is the parsed .devtool.yaml for a project.
type Config struct {
	Project  string    `yaml:"project"`
	Services []Service `yaml:"services"`
}

// Load reads and parses a .devtool.yaml file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	if len(cfg.Services) == 0 {
		return nil, fmt.Errorf("config %s defines no services", path)
	}

	return &cfg, nil
}
