package config

import (
	"fmt"
	"os"
)

type Config struct {
	ApiHost         string
	ApiPort         string
	ApiPortTraining string
	Environment     string
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

func LoadConfig() *Config {
	return &Config{
		ApiHost:         getEnv("API_HOST", "localhost"),
		ApiPort:         getEnv("API_PORT", "3000"),
		ApiPortTraining: getEnv("API_PORT_TRAINING", "3001"),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}
}

var EnvConfigInstance *Config

func Init() {
	EnvConfigInstance = LoadConfig()

	if EnvConfigInstance.Environment != "development" &&
		EnvConfigInstance.Environment != "staging" &&
		EnvConfigInstance.Environment != "production" {
		fmt.Printf("Warning: Unknown environment: %s\n", EnvConfigInstance.Environment)
	}
}
