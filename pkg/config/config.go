// Package config loads and exposes tagguard configuration from .tagguard.yaml.
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all tagguard configuration options.
type Config struct {
	// ExtraKnownKeys whitelists additional tag keys beyond the built-in list.
	// Useful for internal or niche library tags.
	//   extra-known-keys:
	//     - mytag
	//     - internal
	ExtraKnownKeys []string `yaml:"extra-known-keys"`

	// Disable turns off specific rules by name.
	// Valid values: "unknown-key", "validate-rules", "naming-consistency"
	//   disable:
	//     - naming-consistency
	Disable []string `yaml:"disable"`

	// NamingStyle enforces a single naming convention across all serialization
	// tags in the entire project. When set, any tag value that doesn't match
	// this style is flagged.
	// Valid values: "snake_case", "camelCase", "PascalCase", "kebab-case"
	//   naming-style: snake_case
	NamingStyle string `yaml:"naming-style"`
}

// IsDisabled reports whether a rule has been disabled in config.
func (c *Config) IsDisabled(rule string) bool {
	for _, d := range c.Disable {
		if d == rule {
			return true
		}
	}
	return false
}

// Load searches for .tagguard.yaml starting from dir and walking up to the
// filesystem root. Returns an empty Config (not an error) if no file is found.
func Load(dir string) (*Config, error) {
	path, found := findConfigFile(dir)
	if !found {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// findConfigFile walks up from dir looking for .tagguard.yaml.
func findConfigFile(dir string) (string, bool) {
	for {
		candidate := filepath.Join(dir, ".tagguard.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			return "", false
		}
		dir = parent
	}
}
