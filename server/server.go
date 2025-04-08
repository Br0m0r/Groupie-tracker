package server

import (
	"net/http"

	"groupie/handlers"
)

// setupServer configures and returns the HTTP router
func SetupServer() *http.ServeMux {
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
