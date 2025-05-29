package main

import (
	"fmt"
	"log"
	"net/http"

	"groupie/config"
	"groupie/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Starting server with config:\n")
	fmt.Printf("  Environment: %s\n", cfg.App.Environment)
	fmt.Printf("  Port: %s\n", cfg.Server.Port)
	fmt.Printf("  API URL: %s\n", cfg.API.BaseURL)
	fmt.Printf("  Refresh Interval: %v\n", cfg.Cache.RefreshInterval)

	// Configure and start HTTP server using config values
	mux := server.SetupServer()
	httpServer := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	fmt.Printf("Server starting on http://localhost%s\n", cfg.Server.Port)
	fmt.Println("Press Ctrl+C to stop the server")

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
