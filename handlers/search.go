package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"groupie/models"
	"groupie/store"
)

// handlers/search.go - much simpler now
func SearchHandler(dataStore *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))

		// Handle AJAX requests for search suggestions
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			results := dataStore.SearchArtists(query)

			// Set proper headers for JSON response
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "no-cache")

			if err := json.NewEncoder(w).Encode(results); err != nil {
				fmt.Printf("JSON encoding error: %v\n", err)
				w.Write([]byte("[]"))
			}
			return
		}

		// Handle search page with non-empty query
		if query != "" {
			results := dataStore.SearchArtists(query)

			// Redirect to artist page if exactly one result
			if len(results) == 1 {
				http.Redirect(w, r, fmt.Sprintf("/artist?id=%d", results[0].ArtistId), http.StatusSeeOther)
				return
			}

			// Render search results page
			if err := renderSearchResults(w, query, results); err != nil {
				ErrorHandler(w, models.ErrInternalServer, "Failed to execute template")
				return
			}
			return
		}

		// Render empty search page
		if err := renderEmptySearch(w); err != nil {
			ErrorHandler(w, models.ErrInternalServer, "Failed to execute template")
		}
	}
}

// Helper functions for rendering (could move to utils if reused)
func renderSearchResults(w http.ResponseWriter, query string, results []models.SearchResult) error {
	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		return err
	}

	data := models.SearchData{
		Query:       query,
		Results:     results,
		CurrentYear: time.Now().Year(),
	}

	return tmpl.Execute(w, data)
}

func renderEmptySearch(w http.ResponseWriter) error {
	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		return err
	}

	data := models.SearchData{
		Query:       "",
		Results:     nil,
		CurrentYear: time.Now().Year(),
	}

	return tmpl.Execute(w, data)
}
