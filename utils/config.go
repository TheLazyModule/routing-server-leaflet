package utils

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	DBUrl         string `mapstructure:"DATABASE_URL"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv() // Bind environment variables

	err = viper.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &configFileNotFoundError) {
		// Config file not found; ignore error if we can use environment variables
		err = nil
	}

	err = viper.Unmarshal(&config)
	return
}
