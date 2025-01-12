// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"groupie/handlers"
	"groupie/store"
)

const (
	baseURL = "https://groupietrackers.herokuapp.com/api"
	port    = ":8080"
)

func setupServer() *http.ServeMux {
	mux := http.NewServeMux()

	// Route handlers
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/artist", handlers.ArtistHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)

	// Serve static files
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}

func main() {
	// Initialize data store
	fmt.Println("Initializing data store...")
	startTime := time.Now()

	dataStore := store.New()
	if err := dataStore.Initialize(baseURL); err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}

	// Initialize handlers with the data store
	handlers.Initialize(dataStore) // This is now correct

	fmt.Printf("Data store initialized in %v\n", time.Since(startTime))

	// Setup server
	mux := setupServer()
	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
