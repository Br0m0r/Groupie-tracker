// Package handlers provides HTTP request handlers for the web application.
package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"groupie/models"
)

// SearchHandler handles both AJAX search suggestions and full search page requests.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the "q" query parameter from the URL and remove extra whitespace.
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	// Check if the request is an AJAX call (for live search suggestions) by inspecting the custom header.
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// For AJAX requests, perform the search using the provided query.
		results := searchAllData(query)
		// Set the Content-Type header to "application/json" so that the client knows the format.
		w.Header().Set("Content-Type", "application/json")
		// Encode the search results into JSON and write them to the response.
		json.NewEncoder(w).Encode(results)
		// Exit the handler to avoid further processing.
		return
	}

	// If this is a full page search (non-AJAX) and the query is not empty:
	if query != "" {
		// Perform the search and collect the results.
		results := searchAllData(query)

		// If exactly one result is found, redirect the user directly to the artist detail page.
		if len(results) == 1 {
			http.Redirect(w, r, fmt.Sprintf("/artist?id=%d", results[0].ArtistId), http.StatusSeeOther)
			return
		}

		// Otherwise, load the search results page template.
		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			// If the template fails to load, render an error page.
			ErrorHandler(w, ErrInternalServer, "Failed to load template")
			return
		}

		// Create a SearchData struct containing the original query and the search results.
		data := models.SearchData{
			Query:   query,
			Results: results,
		}

		// Execute the template with the search data to generate the full search results page.
		if err := tmpl.Execute(w, data); err != nil {
			// If the template execution fails, render an error page.
			ErrorHandler(w, ErrInternalServer, "Failed to execute template")
			return
		}
		return
	}

	// If the query is empty, simply render an empty search page.
	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	// Prepare a SearchData struct with no query and no results.
	data := models.SearchData{
		Query:   "",
		Results: nil,
	}

	// Execute the template to render the search page.
	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
	}
}

// searchAllData performs the actual search across the entire dataset of artists.
// It searches for matches in various fields (artist name, members, locations, creation date, and first album)
// and returns a slice of SearchResult objects.
func searchAllData(query string) []models.SearchResult {
	var results []models.SearchResult
	// Define separate slices to hold results from different search categories.
	var (
		artistResults   []models.SearchResult
		memberResults   []models.SearchResult
		locationResults []models.SearchResult
		dateResults     []models.SearchResult
		albumResults    []models.SearchResult
	)

	// Retrieve all artist data from the datastore.
	artists := dataStore.GetAllArtists()
	// Convert the query to lowercase to enable case-insensitive searching.
	query = strings.ToLower(query)
	// Check if the query is a single character (special handling for single-letter searches).
	isSingleLetter := len([]rune(query)) == 1

	// Loop over each artist in the datastore.
	for _, artist := range artists {
		// ------------------- Artist Name Search -------------------
		// Convert the artist's name to lowercase.
		artistNameLower := strings.ToLower(artist.Name)
		// Check if the artist's name starts with the query (for single letters) or contains the query.
		if (isSingleLetter && strings.HasPrefix(artistNameLower, query)) ||
			(!isSingleLetter && strings.Contains(artistNameLower, query)) {
			// Append a search result for the artist name.
			artistResults = append(artistResults, models.SearchResult{
				Text:        artist.Name,
				Type:        "artist/band",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
				ArtistId:    artist.ID,
			})
		}

		// ------------------- Members Search -------------------
		// Loop over each member of the current artist.
		for _, member := range artist.Members {
			memberLower := strings.ToLower(member)
			// Check if the member's name matches the query.
			if (isSingleLetter && strings.HasPrefix(memberLower, query)) ||
				(!isSingleLetter && strings.Contains(memberLower, query)) {
				// Append a search result for the member.
				memberResults = append(memberResults, models.SearchResult{
					Text:        member,
					Type:        "member",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Member of %s", artist.Name),
					ArtistId:    artist.ID,
				})
			}
		}

		// ------------------- Locations Search -------------------
		// Loop over each location in the artist's formatted locations list.
		for _, location := range artist.LocationsList {
			locationLower := strings.ToLower(location)
			// Check if the location matches the query.
			if (isSingleLetter && strings.HasPrefix(locationLower, query)) ||
				(!isSingleLetter && strings.Contains(locationLower, query)) {
				// Append a search result for the location.
				locationResults = append(locationResults, models.SearchResult{
					Text:        location,
					Type:        "location",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Concert location for %s", artist.Name),
					ArtistId:    artist.ID,
				})
			}
		}

		// ------------------- Creation Date Search -------------------
		// Only process creation date search if the query has more than one character.
		if !isSingleLetter {
			// Convert the creation date to a string.
			creationStr := fmt.Sprintf("%d", artist.CreationDate)
			// Check if the creation date string contains the query.
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

		// ------------------- First Album Search -------------------
		// Convert the first album information to lowercase.
		albumLower := strings.ToLower(artist.FirstAlbum)
		// Check if the first album field matches the query.
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

	// Combine the results from all categories in a specific order:
	// Artist names first, followed by members, locations, creation dates, and album details.
	results = append(results, artistResults...)
	results = append(results, memberResults...)
	results = append(results, locationResults...)
	results = append(results, dateResults...)
	results = append(results, albumResults...)

	// Return the combined slice of search results.
	return results
}
