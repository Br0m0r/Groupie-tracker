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
	"groupie/store"
)

// SearchHandler handles both AJAX search suggestions and full search page requests
func SearchHandler(dataStore *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))

		// Handle AJAX requests for search suggestions
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			results := searchAllData(dataStore, query)

			// DEBUG: Log what we're returning
			fmt.Printf("AJAX Search: query='%s', results=%d\n", query, len(results))
			if len(results) > 0 {
				fmt.Printf("First result: %+v\n", results[0])
			}

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
			results := searchAllData(dataStore, query)

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
}

// searchAllData performs search across all artist data (UPDATED for pointers)
func searchAllData(dataStore *store.DataStore, query string) []models.SearchResult {
	// Early return for empty queries
	if query == "" {
		return []models.SearchResult{}
	}

	// Get all artists (now returns []*models.Artist)
	artists := dataStore.GetAllArtists()
	query = strings.ToLower(query)
	isSingleLetter := len([]rune(query)) == 1

	// Pre-allocate results slice
	results := make([]models.SearchResult, 0, len(artists)*2)

	// Helper function for matching
	searchMatch := func(value, query string, isSingleLetter bool) bool {
		if isSingleLetter {
			return strings.HasPrefix(strings.ToLower(value), query)
		}
		return strings.Contains(strings.ToLower(value), query)
	}

	// Process each artist (artist is now a pointer)
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

	// Sort results by type priority
	sortResultsByType(results)

	return results
}

// sortResultsByType sorts results with a priority order
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
