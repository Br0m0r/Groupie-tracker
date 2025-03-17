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

// search all data function base on query string
func searchAllData(query string) []models.SearchResult {
	var results []models.SearchResult
	var (
		artistResults   []models.SearchResult
		memberResults   []models.SearchResult
		locationResults []models.SearchResult
		dateResults     []models.SearchResult
		albumResults    []models.SearchResult
	)

	artists := dataStore.GetAllArtists()
	query = strings.ToLower(query)
	isSingleLetter := len([]rune(query)) == 1

	for _, artist := range artists {
		// Artist name search
		artistNameLower := strings.ToLower(artist.Name)
		if (isSingleLetter && strings.HasPrefix(artistNameLower, query)) ||
			(!isSingleLetter && strings.Contains(artistNameLower, query)) {
			artistResults = append(artistResults, models.SearchResult{
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
				memberResults = append(memberResults, models.SearchResult{
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
				locationResults = append(locationResults, models.SearchResult{
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
			creationStr := fmt.Sprintf("%d", artist.CreationDate) // creation date is int so convert to string
			if strings.Contains(creationStr, query) {
				dateResults = append(dateResults, models.SearchResult{
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
			albumResults = append(albumResults, models.SearchResult{
				Text:        fmt.Sprintf("%s - %s", artist.Name, artist.FirstAlbum),
				Type:        "first album",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("First album by %s", artist.Name),
				ArtistId:    artist.ID,
			})
		}
	}

	// Combine results in desired order
	results = append(results, artistResults...)   // Artists first
	results = append(results, memberResults...)   // Then members
	results = append(results, locationResults...) // Then locations
	results = append(results, dateResults...)     // Then creation dates
	results = append(results, albumResults...)    // Finally albums

	return results
}
