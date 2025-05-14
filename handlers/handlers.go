// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"groupie/models"
	"groupie/store"
)

// HomeHandler serves the main page with artist listings
func HomeHandler(dataStore *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			ErrorHandler(w, ErrNotFound, "Page not exist")
			return
		}

		data := models.FilterData{
			Artists:         dataStore.GetArtistCards(),
			UniqueLocations: dataStore.UniqueLocations(),
			SelectedFilters: getDefaultFilterParams(dataStore),
			TotalResults:    len(dataStore.GetArtistCards()),
			CurrentPath:     r.URL.Path,
			CurrentYear:     time.Now().Year(),
		}

		if err := executeFilterTemplate(w, data); err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to process template")
			return
		}
	}
}

// ArtistHandler serves individual artist details pages
func ArtistHandler(dataStore *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		artist, err := dataStore.GetArtist(id)
		if err != nil {
			ErrorHandler(w, ErrNotFound, "Artist not found")
			return
		}
		artist.CurrentYear = time.Now().Year()

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
}
