package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// LoadConfig loads configuration from .env file - .env file is REQUIRED
func LoadConfig() error {
	// Load .env file - REQUIRED, no fallback
	if err := loadEnvFile(".env"); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	// Validate all required environment variables exist
	requiredVars := []string{
		"PORT",
		"SERVER_READ_TIMEOUT",
		"SERVER_WRITE_TIMEOUT", 
		"SERVER_IDLE_TIMEOUT",
		"API_BASE_URL",
		"API_TIMEOUT",
		"REFRESH_INTERVAL",
		"COORDINATES_RATE_LIMIT",
		"ENVIRONMENT",
		"LOG_LEVEL",
	}

	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			return fmt.Errorf("required environment variable %s is not set", varName)
		}
	}

	// Validate the configuration
	if err := validateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

// validateConfig performs basic validation on config values
func validateConfig() error {
	// Validate port
	port := os.Getenv("PORT")
	if !strings.HasPrefix(port, ":") {
		return fmt.Errorf("PORT must start with ':'")
	}
	
	portNum, err := strconv.Atoi(port[1:])
	if err != nil {
		return fmt.Errorf("invalid port number: %s", port[1:])
	}
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", portNum)
	}

	// Validate durations
	durationVars := []string{
		"SERVER_READ_TIMEOUT",
		"SERVER_WRITE_TIMEOUT", 
		"SERVER_IDLE_TIMEOUT",
		"API_TIMEOUT",
		"REFRESH_INTERVAL",
		"COORDINATES_RATE_LIMIT",
	}

	for _, varName := range durationVars {
		if _, err := time.ParseDuration(os.Getenv(varName)); err != nil {
			return fmt.Errorf("invalid duration for %s: %s", varName, os.Getenv(varName))
		}
	}

	// Validate environment
	env := os.Getenv("ENVIRONMENT")
	if env != "development" && env != "production" {
		return fmt.Errorf("ENVIRONMENT must be 'development' or 'production', got '%s'", env)
	}

	// Validate API URL
	if os.Getenv("API_BASE_URL") == "" {
		return fmt.Errorf("API_BASE_URL cannot be empty")
	}

	return nil
}

// Helper functions to get typed values
func GetPort() string {
	return os.Getenv("PORT")
}

func GetAPIBaseURL() string {
	return os.Getenv("API_BASE_URL")
}

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func GetServerReadTimeout() time.Duration {
	duration, _ := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	return duration
}

func GetServerWriteTimeout() time.Duration {
	duration, _ := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	return duration
}

func GetServerIdleTimeout() time.Duration {
	duration, _ := time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	return duration
}

func GetAPITimeout() time.Duration {
	duration, _ := time.ParseDuration(os.Getenv("API_TIMEOUT"))
	return duration
}

func GetRefreshInterval() time.Duration {
	duration, _ := time.ParseDuration(os.Getenv("REFRESH_INTERVAL"))
	return duration
}

func GetCoordinatesRateLimit() time.Duration {
	duration, _ := time.ParseDuration(os.Getenv("COORDINATES_RATE_LIMIT"))
	return duration
}