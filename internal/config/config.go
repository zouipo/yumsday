package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"go.yaml.in/yaml/v4"
)

const (
	envVarPrefix     = "YUMSDAY_"
	configPathEnvVar = envVarPrefix + "CONFIG_PATH"
	logLevelEnvVar   = envVarPrefix + "LOG_LEVEL"
	hostEnvVar       = envVarPrefix + "HOST"
	portEnvVar       = envVarPrefix + "PORT"
	dbPathEnvVar     = envVarPrefix + "DB_PATH"
)

var (
	// create a logger just to get the same formatting as the rest of the program
	// when logging that we parse the yaml config in readConfigFile()
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
)

type Config struct {
	LogLevel slog.Level `yaml:"log_level"`
	Host     string     `yaml:"host"`
	Port     uint       `yaml:"port"`
	DBPath   string     `yaml:"db_path"`
}

func newConfig() Config {
	return Config{
		LogLevel: slog.LevelInfo,
		Host:     "[::0]",
		Port:     8080,
		DBPath:   "yumsday.db",
	}
}

// LoadConfig returns the configuration, with the overwrites from the yaml file and environment variables applied.
// The sources predence is as follow, in decreasing order:
// environment variable > yaml > default, i.e. env vars overwrite yaml values which overwrite defaults.
func LoadConfig() (Config, error) {
	yamlPath := configPath()

	cfg := newConfig()

	if err := readConfigFile(&cfg, yamlPath); err != nil {
		return cfg, err
	}

	if err := readEnvVars(&cfg); err != nil {
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
func readConfigFile(cfg *Config, yamlPath string) error {
	buf, err := os.ReadFile(yamlPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	logger.Info("loading yaml configuration", "path", yamlPath)

	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return fmt.Errorf("failed to parse yaml config: %w", err)
	}

	return nil
}

func readEnvVars(cfg *Config) error {
	if logLevel := os.Getenv(logLevelEnvVar); logLevel != "" {
		cfg.LogLevel.UnmarshalText([]byte(logLevel))
	}

	if host := os.Getenv(hostEnvVar); host != "" {
		cfg.Host = host
	}

	if portStr := os.Getenv(portEnvVar); portStr != "" {
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return fmt.Errorf("failed to parse port '%s' from environment variable: %w", portStr, err)
		}

		cfg.Port = uint(port)
	}

	if dbPath := os.Getenv(dbPathEnvVar); dbPath != "" {
		cfg.DBPath = dbPath
	}

	return nil
}
