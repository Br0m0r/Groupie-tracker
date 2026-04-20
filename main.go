package main

import (
	"log"
	"net/http"
	"time"

	"groupie/handlers"
	"groupie/store"
)

const port = ":8080"

func setupServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/artist", handlers.ArtistHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.HandleFunc("/filter", handlers.FilterHandler)
	mux.HandleFunc("/api/coordinates", handlers.GetLocationCoordinates)

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}

func main() {
	dataStore := store.New()
	if err := dataStore.Initialize(); err != nil {
		log.Fatalf("Failed to initialize data store: %v", err)
	}

	handlers.Initialize(dataStore)

	mux := setupServer()
	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on http://localhost%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
