// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"html/template"
	"net/http"
	"sort"
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

// HomeHandler serves the main page with artist listings
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

	// Get all locations
	locations := make(map[string]bool)
	for _, artist := range dataStore.GetAllArtists() {
		for _, loc := range artist.LocationsList {
			locations[loc] = true
		}
	}
	uniqueLocations := make([]string, 0, len(locations))
	for loc := range locations {
		uniqueLocations = append(uniqueLocations, loc)
	}
	sort.Strings(uniqueLocations)

	// Create template with custom functions
	funcMap := template.FuncMap{
		"intRange": func(min, max int) []int {
			a := make([]int, max-min)
			for i := range a {
				a[i] = min + i
			}
			return a
		},
	}

	tmpl := template.New("index.html").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	data := struct {
		Artists   []models.ArtistCard
		Locations []string
	}{
		Artists:   dataStore.GetArtistCards(),
		Locations: uniqueLocations,
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
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
