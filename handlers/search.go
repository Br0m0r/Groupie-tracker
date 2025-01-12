// handlers/search.go
package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type SearchResult struct {
	Text        string `json:"text"`
	Type        string `json:"type"`
	ArtistName  string `json:"artistName"`
	Description string `json:"description"`
	ArtistId    int    `json:"artistId,omitempty"`
}

type SearchData struct {
	Query   string
	Results []SearchResult
	ShowAll bool
}

// handlers/search.go
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	// Handle AJAX requests for suggestions
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		results := searchAllData(query)
		// Limit to first 5 suggestions only
		if len(results) > 5 {
			results = results[:5]
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	// Handle regular form submission
	if query != "" {
		results := searchAllData(query) // Full results for search page

		// If exactly one result, redirect to artist page
		if len(results) == 1 {
			http.Redirect(w, r, fmt.Sprintf("/artist?id=%d", results[0].ArtistId), http.StatusSeeOther)
			return
		}

		// Otherwise render search results page
		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to load template")
			return
		}

		data := SearchData{
			Query:   query,
			Results: results,
			ShowAll: false,
		}

		if err := tmpl.Execute(w, data); err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to execute template")
			return
		}
		return
	}

	// Handle empty query
	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	data := SearchData{
		Query:   "",
		Results: nil,
		ShowAll: true,
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
	}
}

func searchAllData(query string) []SearchResult {
	var results []SearchResult
	artists := dataStore.GetAllArtists()
	query = strings.ToLower(query)

	// Check if it's a single letter query
	isSingleLetter := len([]rune(query)) == 1

	for _, artist := range artists {
		// Artist name search
		artistNameLower := strings.ToLower(artist.Name)
		if (isSingleLetter && strings.HasPrefix(artistNameLower, query)) ||
			(!isSingleLetter && strings.Contains(artistNameLower, query)) {
			results = append(results, SearchResult{
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
				results = append(results, SearchResult{
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
				results = append(results, SearchResult{
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
				results = append(results, SearchResult{
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
			results = append(results, SearchResult{
				Text:        fmt.Sprintf("%s - %s", artist.Name, artist.FirstAlbum),
				Type:        "first album",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("First album by %s", artist.Name),
				ArtistId:    artist.ID,
			})
		}
	}

	return results
}
