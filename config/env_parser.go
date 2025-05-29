package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format at line %d: %s", lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = removeQuotes(value)

		// Set environment variable
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", filename, err)
	}

	return nil
}

// removeQuotes removes surrounding quotes from a string
func removeQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// getEnvWithDefault gets an environment variable with a default fallback
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDurationWithDefault parses a duration from environment with default fallback
func parseDurationWithDefault(key, defaultValue string) time.Duration {
	value := getEnvWithDefault(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		fmt.Printf("Warning: Invalid duration for %s (%s), using default %s\n", key, value, defaultValue)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}

// parseIntWithDefault parses an integer from environment with default fallback
func parseIntWithDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// You can implement integer parsing here if needed
	// For now, we'll use the default approach
	return defaultValue
}
