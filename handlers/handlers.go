// handlers/handlers.go
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"groupie/models"
	"groupie/store" // Add this import
)

var dataStore *store.DataStore // Change to store.DataStore

func Initialize(ds *store.DataStore) { // Change parameter type
	dataStore = ds
}

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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

	// Parse template
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	// Get data
	artists := dataStore.GetArtistCards()

	// Execute template
	err = tmpl.Execute(w, artists)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get the referer query if coming from search
	var refererQuery string
	referer := r.Header.Get("Referer")
	if referer != "" {
		if refURL, err := url.Parse(referer); err == nil {
			if refURL.Path == "/search" {
				refererQuery = refURL.Query().Get("q")
			}
		}
	}

	// Create template data
	data := ArtistTemplateData{
		Artist:       artist,
		RefererQuery: refererQuery,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		return
	}
}
