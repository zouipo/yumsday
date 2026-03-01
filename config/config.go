package config

import "github.com/spf13/viper"

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBPath   string `mapstructure:"db_path"`
	LogLevel string `mapstructure:"log_level"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("YUMSDAY")
	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
