package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"groupie/handlers"
	"groupie/server"
	"groupie/store"
)

// API base URL and server port
const (
	port = ":8080"
)

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

	// Start periodic data refresh
	go periodicDataRefresh(dataStore)

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

// periodicDataRefresh refreshes the API data every hour
func periodicDataRefresh(dataStore *store.DataStore) {
	refreshInterval := 1 * time.Hour

	for {
		// Sleep for the refresh interval
		time.Sleep(refreshInterval)

		// Log refresh attempt
		fmt.Println("Performing scheduled data refresh...")
		refreshStartTime := time.Now()

		// Create a temporary data store
		tempStore := store.New()

		// Initialize the temporary store - this is a blocking call
		// that will fully complete before continuing
		if err := tempStore.Initialize(); err != nil {
			log.Printf("Error refreshing data: %v", err)
			continue // Skip this refresh cycle on error
		}

		// Now that initialization is complete, swap the data atomically
		dataStore.SwapData(tempStore)

		fmt.Printf("Data refreshed successfully in %v\n", time.Since(refreshStartTime))
	}
}
