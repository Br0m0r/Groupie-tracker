package config

import (
	"fmt"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig
	API    APIConfig
	Cache  CacheConfig
	App    AppConfig
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Port         string        // :8080
	ReadTimeout  time.Duration // 15s
	WriteTimeout time.Duration // 15s
	IdleTimeout  time.Duration // 60s
}

// APIConfig holds external API settings
type APIConfig struct {
	BaseURL string        // https://groupietrackers.herokuapp.com/api
	Timeout time.Duration // 10s
}

// CacheConfig holds caching and refresh settings
type CacheConfig struct {
	RefreshInterval      time.Duration // 4m
	CoordinatesRateLimit time.Duration // 1s
}

// AppConfig holds general app settings
type AppConfig struct {
	Environment string // development, production
	LogLevel    string // info, debug, error
}

// Global config instance
var appConfig *Config

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore if missing)
	if err := loadEnvFile(".env"); err != nil {
		// Log warning but continue with defaults
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Create config with values from environment (with defaults)
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvWithDefault("PORT", ":8080"),
			ReadTimeout:  parseDurationWithDefault("SERVER_READ_TIMEOUT", "15s"),
			WriteTimeout: parseDurationWithDefault("SERVER_WRITE_TIMEOUT", "15s"),
			IdleTimeout:  parseDurationWithDefault("SERVER_IDLE_TIMEOUT", "60s"),
		},
		API: APIConfig{
			BaseURL: getEnvWithDefault("API_BASE_URL", "https://groupietrackers.herokuapp.com/api"),
			Timeout: parseDurationWithDefault("API_TIMEOUT", "10s"),
		},
		Cache: CacheConfig{
			RefreshInterval:      parseDurationWithDefault("REFRESH_INTERVAL", "4m"),
			CoordinatesRateLimit: parseDurationWithDefault("COORDINATES_RATE_LIMIT", "1s"),
		},
		App: AppConfig{
			Environment: getEnvWithDefault("ENVIRONMENT", "development"),
			LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
		},
	}

	// Validate the configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Set global config
	appConfig = config
	return config, nil
}

// Get returns the global config instance
func Get() *Config {
	if appConfig == nil {
		panic("Config not loaded. Call LoadConfig() first")
	}
	return appConfig
}

// validate performs basic validation on config values
func (c *Config) validate() error {
	// Validate port format
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}
	if c.Server.Port[0] != ':' {
		return fmt.Errorf("server port must start with ':'")
	}

	// Validate port number
	portNum, err := strconv.Atoi(c.Server.Port[1:])
	if err != nil {
		return fmt.Errorf("invalid port number: %s", c.Server.Port[1:])
	}
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", portNum)
	}

	// Validate durations are positive
	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("server read timeout must be positive")
	}
	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("server write timeout must be positive")
	}
	if c.Cache.RefreshInterval <= 0 {
		return fmt.Errorf("refresh interval must be positive")
	}

	// Validate API URL is not empty
	if c.API.BaseURL == "" {
		return fmt.Errorf("API base URL cannot be empty")
	}

	// Validate environment
	if c.App.Environment != "development" && c.App.Environment != "production" {
		return fmt.Errorf("environment must be 'development' or 'production', got '%s'", c.App.Environment)
	}

	return nil
}
