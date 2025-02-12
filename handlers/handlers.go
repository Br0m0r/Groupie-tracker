// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"groupie/store"
)

// dataStore holds the application's data layer
var dataStore *store.DataStore

// Initialize sets up the handlers package with a data store instance
func Initialize(ds *store.DataStore) {
	dataStore = ds
}

<<<<<<< HEAD
type ArtistTemplateData struct {
	models.Artist
	RefererQuery string
}

// FuncMap for template functions
var funcMap = template.FuncMap{
	"join": strings.Join,
	"parseDate": func(date string) int {
		// Assuming date is in format "DD-MM-YYYY"
		parts := strings.Split(date, "-")
		if len(parts) >= 3 {
			year, _ := strconv.Atoi(parts[2])
			return year
		}
		return 0
	},
}

=======
// HomeHandler serves the main page with artist listings
>>>>>>> giannis
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

<<<<<<< HEAD
	// Parse template
=======
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
>>>>>>> giannis
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

<<<<<<< HEAD
	// Get data
	artists := dataStore.GetArtistCards()

	// Execute template
	err = tmpl.Execute(w, artists)
	if err != nil {
		log.Printf("Template execution error: %v", err)
=======
	// Execute template with data
	if err := tmpl.Execute(w, data); err != nil {
>>>>>>> giannis
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
