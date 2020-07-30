package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config wraps all configs that are used in this project
type Config struct {
	Node  NodeConfig  `mapstructure:"node"`
	Redis RedisConfig `mapstructure:"redis"`
	Web   WebConfig   `mapstructure:"web"`
}

// NodeConfig wraps both endpoints for Tendermint RPC Node and REST API Server
type NodeConfig struct {
	RPC string `mapstructure:"rpc"`
	LCD string `mapstructure:"lcd"`
}

// RedisConfig is configure for redis
type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// WebConfig wraps port number of this project
type WebConfig struct {
	Port string `mapstructure:"port"`
}

// ParseConfig attempts to read and parse config.yaml from the given path
// An error reading or parsing the config results in a panic.
func ParseConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../") // for test cases
	// viper.AddConfigPath("/home/ubuntu/xxxx") // for production

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	if viper.GetString("active") == "" {
		log.Fatal("define active param in your config file.")
	}

	var config Config
	sub := viper.Sub(viper.GetString("active"))
	sub.Unmarshal(&config)

	return &config
}
