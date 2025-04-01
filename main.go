package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"groupie/handlers"
	"groupie/store"
)

// API base URL and server port
const (
	baseURL = "https://groupietrackers.herokuapp.com/api"
	port    = ":8080"
)

// setupServer configures and returns the HTTP router
func setupServer() *http.ServeMux {
	mux := http.NewServeMux()

	// Configure routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/artist", handlers.ArtistHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.HandleFunc("/filter", handlers.FilterHandler)
	mux.HandleFunc("/api/coordinates", handlers.GetLocationCoordinates)

	// Static file handling
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}

func main() {
	// Initialize data store
	fmt.Println("Initializing data store...")
	startTime := time.Now()

	dataStore := store.New()
	if err := dataStore.Initialize(); err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}

	handlers.Initialize(dataStore)
	fmt.Printf("Data store initialized in %v\n", time.Since(startTime))

	// Start background refresh of datastore periodically
	go func() {
		ticker := time.NewTicker(4 * time.Minute) // adjust refresh interval as needed
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Refreshing datastore...")
			if err := dataStore.Initialize(); err != nil {
				log.Printf("Error refreshing datastore: %v", err)
			} else {
				log.Println("Datastore refreshed successfully.")
			}
		}
	}()

	// Configure and start HTTP server
	mux := setupServer()
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
