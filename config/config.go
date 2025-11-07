package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server      ServerConfig      `json:"server"`
	Backends    []BackendConfig   `json:"backends"`
	HealthCheck HealthCheckConfig `json:"healthCheck"`
	Strategy    StrategyConfig    `json:"strategy"`
	Logging     LoggingConfig     `json:"logging"`
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
}

// BackendConfig holds backend server configuration
type BackendConfig struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

// HealthCheckConfig holds health check settings
type HealthCheckConfig struct {
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
	Path     string        `json:"path"`
}

// StrategyConfig holds load balancing strategy settings
type StrategyConfig struct {
	Type string `json:"type"` // roundrobin, leastconnections, random, weighted
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string `json:"level"`  // debug, info, warn, error
	Format string `json:"format"` // text, json
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Backends: []BackendConfig{
			{URL: "http://localhost:8081", Weight: 1},
			{URL: "http://localhost:8082", Weight: 1},
			{URL: "http://localhost:8083", Weight: 1},
		},
		HealthCheck: HealthCheckConfig{
			Interval: 10 * time.Second,
			Timeout:  5 * time.Second,
			Path:     "/health",
		},
		Strategy: StrategyConfig{
			Type: "roundrobin",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
