package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
			Path:   "chat.db",
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
	// Try to load from demo/config.yaml first
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

	// Try to load from current directory
	configPath = "config.yaml"
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

	// Try to load from absolute path
	absPath, err := filepath.Abs("demo/config.yaml")
	if err == nil {
		if _, err := os.Stat(absPath); err == nil {
			data, err := ioutil.ReadFile(absPath)
			if err != nil {
				log.Printf("Warning: Could not read config file %s: %v", absPath, err)
				return DefaultConfig()
			}

			var config Config
			if err := yaml.Unmarshal(data, &config); err != nil {
				log.Printf("Warning: Could not parse config file %s: %v", absPath, err)
				return DefaultConfig()
			}

			log.Printf("Loaded configuration from %s", absPath)
			return &config
		}
	}

	log.Printf("No config file found, using default configuration")
	return DefaultConfig()
}
