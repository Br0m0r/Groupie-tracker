package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"groupie/models"
	"groupie/store"
	"groupie/utils"
)

// dataStore is a package-level variable that holds the application's data layer,
// which is injected from main.go via the Initialize function.
// It provides access to all artist data and related methods.
var dataStore *store.DataStore

// Initialize sets up the handlers package with a data store instance.
// This dependency injection makes the datastore available to all handler functions.
func Initialize(ds *store.DataStore) {
	dataStore = ds
}

// HomeHandler serves the main page of the application.
// It displays a list of artist cards along with filter options.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Verify that the requested URL path is exactly "/".
	// If not, return a 404 error page using the ErrorHandler.
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

	// Prepare the data for the template using a FilterData struct.
	// The .Artists field here is filled with a slice of ArtistCard objects,
	// each containing only the fields needed to display the artist card (ID, Name, Image).
	data := models.FilterData{
		Artists:         dataStore.GetArtistCards(),      // Retrieves simplified artist data.
		UniqueLocations: dataStore.UniqueLocations,       // List of unique concert locations.
		SelectedFilters: utils.GetDefaultFilterParams(),  // Default filter parameters.
		TotalResults:    len(dataStore.GetArtistCards()), // Count of artist cards.
		CurrentPath:     r.URL.Path,                      // Current URL ("/").
	}

	// Render the main page using the filter data.
	// The executeFilterTemplate helper function loads and executes the index.html template.
	if err := executeFilterTemplate(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to process template")
		return
	}
}

// ArtistHandler serves the detailed view for an individual artist.
// It reads the artist ID from the query parameters, fetches the artist's full data,
// and renders the artist detail page using the artist.html template.
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" query parameter from the URL (e.g., /artist?id=3).
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// If no artist ID is provided, return a 400 Bad Request error.
		ErrorHandler(w, ErrBadRequest, "Artist ID is required")
		return
	}

	// Convert the extracted ID from string to integer.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If conversion fails, return an error indicating invalid artist ID format.
		ErrorHandler(w, ErrInvalidID, "Invalid artist ID format")
		return
	}

	// Use the datastore to fetch the full artist data corresponding to the given ID.
	artist, err := dataStore.GetArtist(id)
	if err != nil {
		// If no matching artist is found, return a 404 error.
		ErrorHandler(w, ErrNotFound, "Artist not found")
		return
	}

	// Parse the artist.html template which defines the layout for the artist detail page.
	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		// If template parsing fails, return a 500 Internal Server Error.
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	// Execute the template with the fetched artist data.
	// This generates the final HTML page that displays all the artist's details.
	err = tmpl.Execute(w, artist)
	if err != nil {
		// If template execution fails, return a 500 Internal Server Error.
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		return
	}
}
