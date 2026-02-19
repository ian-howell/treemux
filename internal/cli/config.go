package cli

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds treemux configuration options.
type Config struct {
	// TODO
}

// LoadConfig loads treemux configuration from the specified file path.
func LoadConfig(configFilePath string) (Config, error) {
	yamlData, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}
