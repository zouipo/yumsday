package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	configPathEnvVar = "YUMSDAY_CONFIG_PATH"
)

type Config struct {
	Host     string     `yaml:"host"`
	Port     int        `yaml:"port"`
	DBPath   string     `yaml:"db_path"`
	LogLevel slog.Level `yaml:"log_level"`
}

func newConfig() Config {
	return Config{
		Host:     "[::0]",
		Port:     8080,
		DBPath:   "yumsday.db",
		LogLevel: slog.LevelInfo,
	}
}

func LoadConfig() (Config, error) {
	yamlPath := configPath()

	cfg, err := readConfigFile(yamlPath)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func configPath() string {
	yamlPath := "./config.yaml"

	if yamlPathEnvVar := os.Getenv(configPathEnvVar); yamlPathEnvVar != "" {
		yamlPath = yamlPathEnvVar
	}

	return yamlPath
}

// readConfigFile parses the yaml config file and overwrite the default configuration with its content.
// If the config file doesn't exist, then the default config is returned.
func readConfigFile(yamlPath string) (Config, error) {
	cfg := newConfig()

	buf, err := os.ReadFile(yamlPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		} else {
			return cfg, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse yaml config: %w", err)
	}

	return cfg, nil
}
