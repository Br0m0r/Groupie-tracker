// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"html/template"
	"net/http"
	"strconv"

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
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

	// Get all artists for unique locations
	allArtists := dataStore.GetAllArtists()

	// Prepare data for template including filter data
	data := FilterData{
		Artists:         dataStore.GetArtistCards(),
		UniqueLocations: getUniqueLocations(allArtists),
		SelectedFilters: FilterParams{
			CreationStart:  1950,
			CreationEnd:    2024,
			AlbumStartYear: 1950,
			AlbumEndYear:   2024,
			MemberCounts:   []int{},    // Empty slice for initial state
			Locations:      []string{}, // Empty slice for initial state
		},
		TotalResults: len(dataStore.GetArtistCards()),
	}

	// Define template functions
	funcMap := template.FuncMap{
		"iterate": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"contains": func(slice interface{}, item interface{}) bool {
			switch slice := slice.(type) {
			case []int:
				itemInt, ok := item.(int)
				if !ok {
					return false
				}
				for _, s := range slice {
					if s == itemInt {
						return true
					}
				}
			case []string:
				itemStr, ok := item.(string)
				if !ok {
					return false
				}
				for _, s := range slice {
					if s == itemStr {
						return true
					}
				}
			}
			return false
		},
	}

	// Parse template with functions
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	// Execute template with data
	if err := tmpl.Execute(w, data); err != nil {
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
	// Format the data
	artist.LocationsList = utils.FormatLocationsList(artist.LocationsList)
	artist.DatesList = utils.FormatDatesList(artist.DatesList)
	artist.RelationsList = utils.FormatRelation(artist.RelationsList)
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
