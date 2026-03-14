package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

const (
	CONFIG_PATH_ENV_VAR = "YUMSDAY_CONFIG_PATH"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBPath   string `mapstructure:"db_path"`
	LogLevel string `mapstructure:"log_level"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv(CONFIG_PATH_ENV_VAR))
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); !ok {
			return nil, err
		}
	}

	viper.SetEnvPrefix("YUMSDAY")
	// Check for environment variables following
	// the pattern YUMSDAY_<viper>_<key>_<name>
	// e.g. YUMSDAY_DB_PATH
	viper.AutomaticEnv()

	var config Config
	// Assign values to the Config struct following the names
	// given by the mapstructure attributes
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
