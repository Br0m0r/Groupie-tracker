package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"groupie/config"
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

func periodicDataRefresh(dataStore *store.DataStore) {
	// Delay first refresh so it won’t collide with initial Initialize()
	time.Sleep(config.REFRESH_INTERVAL)

	for {
		log.Println("Performing scheduled incremental refresh...")
		start := time.Now()

		na, nc, err := dataStore.RefreshData()
		if err != nil {
			log.Printf("  ERROR during refresh: %v\n", err)
		} else {
			log.Printf(
				"  Incremental refresh complete: added %d artists, %d coords (took %v)\n",
				na, nc, time.Since(start),
			)
		}

		time.Sleep(config.REFRESH_INTERVAL)
	}
}
