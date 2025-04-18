// handlers/search.go
package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

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
			Query:   query,
			Results: results,
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
		Query:   "",
		Results: nil,
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
	}
}

// handlers/search.go
func searchAllData(query string) []models.SearchResult {
	// Convert query to lowercase once
	query = strings.ToLower(query)
	isSingleLetter := len([]rune(query)) == 1

	// Create a struct to hold categorized results
	type categorizedResults struct {
		artists   []models.SearchResult
		members   []models.SearchResult
		locations []models.SearchResult
		dates     []models.SearchResult
		albums    []models.SearchResult
	}

	// Initialize with reasonable capacities
	results := categorizedResults{
		artists:   make([]models.SearchResult, 0, 10),
		members:   make([]models.SearchResult, 0, 10),
		locations: make([]models.SearchResult, 0, 10),
		dates:     make([]models.SearchResult, 0, 10),
		albums:    make([]models.SearchResult, 0, 10),
	}

	// Get all artists
	artists := dataStore.GetAllArtists()

	// Single pass through all artists
	for _, artist := range artists {
		// Cache frequently used lowercase values
		artistNameLower := strings.ToLower(artist.Name)

		// Artist name search	
		if (isSingleLetter && strings.HasPrefix(artistNameLower, query)) ||
			(!isSingleLetter && strings.Contains(artistNameLower, query)) {
			results.artists = append(results.artists, models.SearchResult{
				Text:        artist.Name,
				Type:        "artist/band",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
				ArtistId:    artist.ID,
			})
		}

		// Members search
		for _, member := range artist.Members {
			memberLower := strings.ToLower(member)
			if (isSingleLetter && strings.HasPrefix(memberLower, query)) ||
				(!isSingleLetter && strings.Contains(memberLower, query)) {
				results.members = append(results.members, models.SearchResult{
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
			locationLower := strings.ToLower(location)
			if (isSingleLetter && strings.HasPrefix(locationLower, query)) ||
				(!isSingleLetter && strings.Contains(locationLower, query)) {
				results.locations = append(results.locations, models.SearchResult{
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
				results.dates = append(results.dates, models.SearchResult{
					Text:        fmt.Sprintf("%s (%d)", artist.Name, artist.CreationDate),
					Type:        "creation date",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
					ArtistId:    artist.ID,
				})
			}
		}

		// First album search
		albumLower := strings.ToLower(artist.FirstAlbum)
		if (isSingleLetter && strings.HasPrefix(albumLower, query)) ||
			(!isSingleLetter && strings.Contains(albumLower, query)) {
			results.albums = append(results.albums, models.SearchResult{
				Text:        fmt.Sprintf("%s - %s", artist.Name, artist.FirstAlbum),
				Type:        "first album",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("First album by %s", artist.Name),
				ArtistId:    artist.ID,
			})
		}
	}

	// Calculate total capacity for combined results
	totalCapacity := len(results.artists) + len(results.members) +
		len(results.locations) + len(results.dates) + len(results.albums)

	// Pre-allocate final slice with exact capacity needed
	finalResults := make([]models.SearchResult, 0, totalCapacity)

	// Combine in the desired order - no sorting needed!
	finalResults = append(finalResults, results.artists...)
	finalResults = append(finalResults, results.members...)
	finalResults = append(finalResults, results.locations...)
	finalResults = append(finalResults, results.dates...)
	finalResults = append(finalResults, results.albums...)

	return finalResults
}
