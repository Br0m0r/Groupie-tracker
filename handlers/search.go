// handlers/search.go
package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"

	"groupie/models"
)

// SearchHandler handles both AJAX search suggestions and full search page requests
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	// Handle AJAX requests for search suggestions
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		results := searchAllData(query)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	// Handle search page with non-empty query
	if query != "" {
		results := searchAllData(query)

		// Redirect to artist page if exactly one result
		if len(results) == 1 {
			http.Redirect(w, r, fmt.Sprintf("/artist?id=%d", results[0].ArtistId), http.StatusSeeOther)
			return
		}

		// Render search results page
		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to load template")
			return
		}

		data := models.SearchData{
			Query:       query,
			Results:     results,
			CurrentYear: time.Now().Year(),
		}

		if err := tmpl.Execute(w, data); err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to execute template")
			return
		}
		return
	}

	// Render empty search page
	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	data := models.SearchData{
		Query:       "",
		Results:     nil,
		CurrentYear: time.Now().Year(),
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
	}
}

func searchAllData(query string) []models.SearchResult {
	// Early return for empty queries
	if query == "" {
		return nil
	}

	artists := dataStore.GetAllArtists()
	query = strings.ToLower(query)
	isSingleLetter := len([]rune(query)) == 1

	// Pre-allocate results slice with estimated capacity
	// Average estimated size based on expected matches per artist
	estimatedSize := len(artists) * 2
	results := make([]models.SearchResult, 0, estimatedSize)

	// Create a single searchMatch helper function
	searchMatch := func(value, query string, isSingleLetter bool) bool {
		if isSingleLetter {
			return strings.HasPrefix(strings.ToLower(value), query)
		}
		return strings.Contains(strings.ToLower(value), query)
	}

	// Process each artist
	for _, artist := range artists {
		// Artist name search
		if searchMatch(artist.Name, query, isSingleLetter) {
			results = append(results, models.SearchResult{
				Text:        artist.Name,
				Type:        "artist/band",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
				ArtistId:    artist.ID,
			})
		}

		// Members search
		for _, member := range artist.Members {
			if searchMatch(member, query, isSingleLetter) {
				results = append(results, models.SearchResult{
					Text:        member,
					Type:        "member",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Member of %s", artist.Name),
					ArtistId:    artist.ID,
				})
			}
		}

		// Locations search
		for _, location := range artist.LocationsList {
			if searchMatch(location, query, isSingleLetter) {
				results = append(results, models.SearchResult{
					Text:        location,
					Type:        "location",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Concert location for %s", artist.Name),
					ArtistId:    artist.ID,
				})
			}
		}

		// Creation date search - only for non-single letter queries
		if !isSingleLetter {
			creationStr := fmt.Sprintf("%d", artist.CreationDate)
			if strings.Contains(creationStr, query) {
				results = append(results, models.SearchResult{
					Text:        fmt.Sprintf("%s (%d)", artist.Name, artist.CreationDate),
					Type:        "creation date",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
					ArtistId:    artist.ID,
				})
			}
		}

		// First album search
		if searchMatch(artist.FirstAlbum, query, isSingleLetter) {
			results = append(results, models.SearchResult{
				Text:        fmt.Sprintf("%s - %s", artist.Name, artist.FirstAlbum),
				Type:        "first album",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("First album by %s", artist.Name),
				ArtistId:    artist.ID,
			})
		}
	}

	// Sort the results by type if needed
	sortResultsByType(results)

	return results
}

// Helper function to sort results by type
func sortResultsByType(results []models.SearchResult) {
	// Define type priorities
	typePriority := map[string]int{
		"artist/band":   1,
		"member":        2,
		"location":      3,
		"creation date": 4,
		"first album":   5,
	}

	// Sort by type priority
	sort.Slice(results, func(i, j int) bool {
		return typePriority[results[i].Type] < typePriority[results[j].Type]
	})
}
