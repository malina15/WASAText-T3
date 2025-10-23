package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Path   string `yaml:"path"`
		Driver string `yaml:"driver"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Port string `yaml:"port"`
			Host string `yaml:"host"`
		}{
			Port: "8080",
			Host: "localhost",
		},
		Database: struct {
			Path   string `yaml:"path"`
			Driver string `yaml:"driver"`
		}{
			Path:   "./data/chat.db", // Use relative path for local development
			Driver: "sqlite3",
		},
		Logging: struct {
			Level  string `yaml:"level"`
			Format string `yaml:"format"`
		}{
			Level:  "info",
			Format: "text",
		},
	}
}

// LoadConfig loads configuration from file or returns default
func LoadConfig() *Config {
	// For local development, use default config to avoid Docker paths
	// Only use demo/config.yaml in Docker environment
	if os.Getenv("DOCKER_ENV") == "true" {
		// Try to load from demo/config.yaml first (for Docker)
		configPath := "demo/config.yaml"
		if _, err := os.Stat(configPath); err == nil {
			data, err := ioutil.ReadFile(configPath)
			if err != nil {
				log.Printf("Warning: Could not read config file %s: %v", configPath, err)
				return DefaultConfig()
			}

			var config Config
			if err := yaml.Unmarshal(data, &config); err != nil {
				log.Printf("Warning: Could not parse config file %s: %v", configPath, err)
				return DefaultConfig()
			}

			log.Printf("Loaded configuration from %s", configPath)
			return &config
		}
	}

	// For local development, use default config
	log.Printf("Using default configuration for local development")
	return DefaultConfig()
}
