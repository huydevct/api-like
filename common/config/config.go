package config

import (
	"fmt"
	"os"

	"app/common/adapter"

	Viper "github.com/spf13/viper"
)

// Config struct ..
type Config struct {
	Name     string            `mapstructure:"name"`
	Port     map[string]string `mapstructure:"port"`
	Version  string            `mapstructure:"version"`
	Debug    bool              `mapstructure:"debug"`
	Mongo    adapter.Mongos
	Redis    adapter.Redises
	RabbitMQ adapter.Rabbits
	API      adapter.APIs
	Other    adapter.Others
	Cache    adapter.Caches
}

var config *Config

func init() {
	var folder string

	env := os.Getenv("APPLICATION_ENV")

	switch env {
	case "master", "dev", "uat", "localhost":
		folder = env
	default:
		folder = "dev"
	}

	path := fmt.Sprintf("config/%v", folder)

	//Get base config
	config = new(Config)
	fetchDataToConfig(path, "base", config)

	//Get all sub config
	fetchDataToConfig(path, "mongo", &(config.Mongo))
	fetchDataToConfig(path, "redis", &(config.Redis))
	fetchDataToConfig(path, "rabbit", &(config.RabbitMQ))
	fetchDataToConfig(path, "api", &(config.API))
	fetchDataToConfig(path, "other", &(config.Other))
	fetchDataToConfig(path, "cache", &(config.Cache))
}

func fetchDataToConfig(configPath, configName string, result interface{}) {
	viper := Viper.New()
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	err := viper.ReadInConfig() // Find and read the config file
	if err == nil {             // Handle errors reading the config file
		err = viper.Unmarshal(result)
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}
}

// GetConfig func
func GetConfig() *Config {
	return config
}
