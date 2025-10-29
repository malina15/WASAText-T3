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
	// Try to load from demo/config.yaml first (for Docker)
	configPath := "demo/config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		data, err := ioutil.ReadFile(configPath)
		if err == nil {
			var config Config
			if err := yaml.Unmarshal(data, &config); err == nil {
				log.Printf("Loaded configuration from %s", configPath)
				return &config
			}
		}
	}

	// Check if we're in Docker - if so, use 0.0.0.0, otherwise localhost
	config := DefaultConfig()
	if os.Getenv("DOCKER_ENV") == "true" || os.Getenv("CONTAINER") == "true" {
		config.Server.Host = "0.0.0.0"
		log.Printf("Using Docker configuration (0.0.0.0)")
	} else {
		log.Printf("Using default configuration for local development")
	}
	return config
}
