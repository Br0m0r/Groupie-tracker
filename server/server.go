package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"groupie/handlers"
	"groupie/store"
)

// setupServer configures and returns the HTTP router
// server.go
func SetupServer() *http.ServeMux {
	// Initialize data store
	fmt.Println("Initializing data store...")
	startTime := time.Now()

	dataStore := store.New()
	if err := dataStore.Initialize(); err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}

	fmt.Printf("Data store initialized in %v\n", time.Since(startTime))

	mux := http.NewServeMux()
	go periodicDataRefresh(dataStore) // Start the periodic data refresh in a goroutine
	// Configure routes with handler factories
	mux.HandleFunc("/", handlers.HomeHandler(dataStore))
	mux.HandleFunc("/artist", handlers.ArtistHandler(dataStore))
	mux.HandleFunc("/search", handlers.SearchHandler(dataStore))
	mux.HandleFunc("/filter", handlers.FilterHandler(dataStore))
	mux.HandleFunc("/api/coordinates", handlers.GetLocationCoordinates(dataStore))

	// Static file handling remains the same
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}

// periodicDataRefresh refreshes the API data every hour
func periodicDataRefresh(dataStore *store.DataStore) {
	refreshInterval := 1 * time.Minute

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
