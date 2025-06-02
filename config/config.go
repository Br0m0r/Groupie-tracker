// config/config.go
package config

import "time"

// Configuration constants (replacing .env file values)
const (
	// Server Configuration
	PORT                 = ":8080"
	SERVER_READ_TIMEOUT  = 15 * time.Second
	SERVER_WRITE_TIMEOUT = 15 * time.Second
	SERVER_IDLE_TIMEOUT  = 60 * time.Second

	// API Configuration
	API_BASE_URL = "https://groupietrackers.herokuapp.com/api"
	API_TIMEOUT  = 10 * time.Second

	// Cache Configuration
	REFRESH_INTERVAL       = 60 * time.Minute
	COORDINATES_RATE_LIMIT = 1 * time.Second

	
)
