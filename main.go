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
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Configure and start HTTP server using config values
	mux := server.SetupServer()
	httpServer := &http.Server{
		Addr:         config.GetPort(),
		Handler:      mux,
		ReadTimeout:  config.GetServerReadTimeout(),
		WriteTimeout: config.GetServerWriteTimeout(),
		IdleTimeout:  config.GetServerIdleTimeout(),
	}

	fmt.Printf("Server starting on http://localhost%s\n", config.GetPort())
	fmt.Println("Press Ctrl+C to stop the server")

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}