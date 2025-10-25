package config

import (
	"os"
	"strconv"
)

// Config contains application configuration
type Config struct {
	Port string
}

// Load loads configuration from environment variables
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	return &Config{
		Port: port,
	}
}
