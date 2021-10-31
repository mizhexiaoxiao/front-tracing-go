package config

import (
	"os"

	"github.com/mizhexiaoxiao/front-tracing-go/logger"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(configPath())
	viper.SetConfigName(configName())
	viper.SetConfigType(configType())
	logger.InfoLogger().Println("Configuration initialization")
}

func Parse() error {
	return viper.ReadInConfig()
}

func GetString(key string) string {
	return viper.GetString(key)
}

func configPath() string {
	if configPath := os.Getenv("CONFIG_PATH"); configPath == "" {
		return "."
	} else {
		return configPath
	}
}

func configName() string {
	if configName := os.Getenv("CONFIG_NAME"); configName == "" {
		return "config"
	} else {
		return configName
	}
}

func configType() string {
	if configType := os.Getenv("CONFIG_TYPE"); configType == "" {
		return "yaml"
	} else {
		return configType
	}
}
