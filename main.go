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
	port = ":8080" // Port on which the server will listen
)

// setupServer configures and returns an HTTP multiplexer (ServeMux).
// A multiplexer automatically routes incoming HTTP requests to their designated handlers based on URL patterns.
// This approach is more modular and maintainable than the classic approach where a single handler manually checks URL paths.
func setupServer() *http.ServeMux {
	// Create a new HTTP request multiplexer.
	mux := http.NewServeMux()

	// Register individual routes with corresponding handler functions from the handlers package.
	// Compared to the classic approach (a single handler function with a switch/case),
	// this approach allows each route to be independently managed.
	mux.HandleFunc("/", handlers.HomeHandler)                           // Home page displaying artist cards and filters.
	mux.HandleFunc("/artist", handlers.ArtistHandler)                   // Detailed view for a selected artist.
	mux.HandleFunc("/search", handlers.SearchHandler)                   // Search functionality for artists, members, etc.
	mux.HandleFunc("/filter", handlers.FilterHandler)                   // Filtering of artists based on various criteria.
	mux.HandleFunc("/api/coordinates", handlers.GetLocationCoordinates) // API endpoint for retrieving location coordinates.

	// Serve static files (CSS, JavaScript, images, etc.):
	// http.FileServer returns a handler that serves files from the given directory.
	// http.StripPrefix is used to remove the "/static/" segment from the URL so that the file paths
	// correctly map to the file system (e.g., "/static/css/style.css" becomes "css/style.css" within the "static" directory).
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}

func main() {
	// Log the start of data store initialization.
	fmt.Println("Initializing data store...")
	startTime := time.Now() // Capture the start time to measure initialization duration.

	// Create a new DataStore instance.
	dataStore := store.New()

	// Initialize the DataStore:
	// This function fetches and decodes data from the external API,
	// populating the DataStore with full artist details and additional data.
	if err := dataStore.Initialize(); err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}

	// Pass the initialized DataStore to the handlers package,
	// enabling all HTTP handlers to access and serve the stored data.
	handlers.Initialize(dataStore)

	// Log how long it took to initialize the data store.
	fmt.Printf("Data store initialized in %v\n", time.Since(startTime))

	// Set up the HTTP server using the multiplexer from setupServer().
	// The multiplexer automatically dispatches requests to the correct handler.
	mux := setupServer()
	server := &http.Server{
		Addr:         port,             // The server will listen on the defined port.
		Handler:      mux,              // Use the multiplexer for request routing.
		ReadTimeout:  15 * time.Second, // Maximum time for reading the entire request.
		WriteTimeout: 15 * time.Second, // Maximum time for writing the response.
		IdleTimeout:  60 * time.Second, // Maximum idle time for keeping a connection open.
	}

	// Inform that the server is starting.
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	// Start the HTTP server.
	// If an error occurs (e.g., port already in use), it will be logged and the server will exit.
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
