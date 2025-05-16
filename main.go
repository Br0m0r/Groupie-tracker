package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"groupie/server"
)

// API base URL and server port
const (
	port = ":8080"
)

func main() {
	// Configure and start HTTP server
	mux := server.SetupServer()
	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
