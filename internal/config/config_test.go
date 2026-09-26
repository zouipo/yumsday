package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"testing"

	"go.yaml.in/yaml/v4"
)

const (
	yamlPath = "./config.yaml"
	logLevel = slog.LevelError
	host     = "localhost"
	port     = 1234
	dbPath   = "./yumsday.db"
)

func TestLoadConfig(t *testing.T) {
	yamlConfig := Config{
		Host:   host,
		Port:   port,
		DBPath: "some_directory/yumsday.db",
	}
	yamlContent, err := yaml.Marshal(yamlConfig)
	if err != nil {
		t.Fatalf("unexpected error while marshalling yaml config: %s", err)
	}

	if err := os.Setenv(configPathEnvVar, yamlPath); err != nil {
		t.Fatalf("unexpected error while setting environment variable: %s", err)
	}
	defer os.Unsetenv(configPathEnvVar)

	if err := os.Setenv(dbPathEnvVar, dbPath); err != nil {
		t.Fatalf("unexpected error while setting environment variable: %s", err)
	}
	defer os.Unsetenv(dbPathEnvVar)

	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("unexpected error while writing config yaml: %s", err)
	}
	defer func() {
		if err := os.Remove(yamlPath); err != nil {
			panic(fmt.Errorf("unexpected error while removing test yaml config: %w", err))
		}
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error when loading configuration: %s", err)
	}

	def := newConfig()

	// should be default value as nothing touches it
	if cfg.LogLevel != def.LogLevel {
		t.Errorf("expected default log level (%s), got %s", def.LogLevel, cfg.LogLevel)
	}

	// should have values from yaml
	if cfg.Host != yamlConfig.Host {
		t.Errorf("expected host to match yaml value (%s), got %s", yamlConfig.Host, cfg.Host)
	}
	if cfg.Port != yamlConfig.Port {
		t.Errorf("expected port to match yaml value (%d), got %d", yamlConfig.Port, cfg.Port)
	}

	// should have value from env var
	if cfg.DBPath != dbPath {
		t.Errorf("expected db path to match env (%s), got %s", dbPath, cfg.DBPath)
	}
}

func TestEnvVars(t *testing.T) {
	if err := os.Setenv(logLevelEnvVar, logLevel.String()); err != nil {
		t.Fatalf("unexpected error while setting environment variable: %s", err)
	}
	if err := os.Setenv(hostEnvVar, host); err != nil {
		t.Fatalf("unexpected error while setting environment variable: %s", err)
	}
	if err := os.Setenv(portEnvVar, strconv.FormatUint(port, 10)); err != nil {
		t.Fatalf("unexpected error while setting environment variable: %s", err)
	}
	if err := os.Setenv(dbPathEnvVar, dbPath); err != nil {
		t.Fatalf("unexpected error while setting environment variable: %s", err)
	}
	defer func() {
		os.Unsetenv(logLevelEnvVar)
		os.Unsetenv(hostEnvVar)
		os.Unsetenv(portEnvVar)
		os.Unsetenv(dbPathEnvVar)
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error when loading config: %s", err)
	}

	if cfg.LogLevel != logLevel {
		t.Errorf("expected log level %s, got %s", logLevel, cfg.LogLevel)
	}

	if cfg.Host != host {
		t.Errorf("expected host %s, got %s", host, cfg.Host)
	}

	if cfg.Port != port {
		t.Errorf("expected port %d, got %d", port, cfg.Port)
	}

	if cfg.DBPath != dbPath {
		t.Errorf("expected db path %s, got %s", dbPath, cfg.DBPath)
	}
}

func TestReadConfigFile_ReadError(t *testing.T) {
	yamlContent := "log_level: info"
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0111); err != nil {
		t.Fatalf("unexpected error while writing yaml: %s", err)
	}
	defer func() {
		if err := os.Remove(yamlPath); err != nil {
			panic(fmt.Errorf("unexpected error when removing file: %w", err))
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Errorf("expected error, got none")
	}
	if !errors.Is(err, os.ErrPermission) {
		t.Errorf("expected permission error, got %s", err)
	}
}

func TestReadConfigFile_UnmarshalError(t *testing.T) {
	yamlContent := "log_level: invalid_log_level"
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("unexpected error while writing yaml: %s", err)
	}
	defer func() {
		if err := os.Remove(yamlPath); err != nil {
			panic(fmt.Errorf("unexpected error when removing file: %w", err))
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Errorf("expected error, got none")
	}
}

func TestEnvVars_InvalidLogLevel(t *testing.T) {
	os.Setenv(logLevelEnvVar, "invalid_log_level")
	defer os.Unsetenv(logLevelEnvVar)

	_, err := LoadConfig()
	if err == nil {
		t.Errorf("expected error, got none")
	}
}

func TestEnvVars_InvalidPort(t *testing.T) {
	os.Setenv(portEnvVar, "invalid_port")
	defer os.Unsetenv(portEnvVar)

	_, err := LoadConfig()
	if err == nil {
		t.Errorf("expected error, got none")
	}
}
