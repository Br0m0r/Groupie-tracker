// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"groupie/models"
	"groupie/store"
	"groupie/utils"
)

// dataStore holds the application's data layer
var dataStore *store.DataStore

// Initialize sets up the handlers package with a data store instance
func Initialize(ds *store.DataStore) {
	dataStore = ds
}

// HomeHandler serves the main page with artist listings
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for valid path
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

	// Get data for initial page load
	allArtists := dataStore.GetAllArtists()

	// Create initial filter data with default values
	data := models.FilterData{
		Artists:         dataStore.GetArtistCards(),
		UniqueLocations: utils.GetUniqueLocations(allArtists),
		SelectedFilters: getDefaultFilterParams(),
		TotalResults:    len(allArtists),
		CurrentPath:     r.URL.Path,
	}

	// Use the existing executeFilterTemplate function
	if err := executeFilterTemplate(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to process template")
		return
	}
}

// Return the default filter parameters.
func getDefaultFilterParams() models.FilterParams {
	return models.FilterParams{
		CreationStart:  1950,
		CreationEnd:    2024,
		AlbumStartYear: 1950,
		AlbumEndYear:   2024,
		MemberCounts:   []int{},    // Empty slice - no members selected
		Locations:      []string{}, // Empty slice - no locations selected
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
