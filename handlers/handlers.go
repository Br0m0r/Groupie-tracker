// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"groupie/models"
	"groupie/store"
)

// dataStore holds the application's data layer
var dataStore *store.DataStore

// Initialize sets up the handlers package with a data store instance
func Initialize(ds *store.DataStore) {
	dataStore = ds
}

type HomePageData struct {
	Artists   []models.Artist
	Locations []string // Unique locations for filter dropdown
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

	artists := dataStore.GetAllArtists()

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	err = tmpl.Execute(w, artists) // Pass artists directly
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		return
	}
}

// ArtistHandler serves individual artist details pages
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Extract and validate artist ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorHandler(w, ErrBadRequest, "Artist ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, ErrInvalidID, "Invalid artist ID format")
		return
	}

	// Fetch artist data
	artist, err := dataStore.GetArtist(id)
	if err != nil {
		ErrorHandler(w, ErrNotFound, "Artist not found")
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	err = tmpl.Execute(w, artist)
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		return
	}
}
