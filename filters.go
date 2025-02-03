// handlers/filters.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"groupie/models"
)

// FilterHandler handles AJAX filter requests.
// It reads filter parameters from the query string, filters the data,
// and returns the results as JSON.
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is made via AJAX.
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "AJAX requests only", http.StatusBadRequest)
		return
	}

	// Retrieve range filter parameters (as strings).
	creationMinStr := strings.TrimSpace(r.URL.Query().Get("creation_min"))
	creationMaxStr := strings.TrimSpace(r.URL.Query().Get("creation_max"))
	albumMinStr := strings.TrimSpace(r.URL.Query().Get("album_min"))
	albumMaxStr := strings.TrimSpace(r.URL.Query().Get("album_max"))
	membersMinStr := strings.TrimSpace(r.URL.Query().Get("members_min"))
	membersMaxStr := strings.TrimSpace(r.URL.Query().Get("members_max"))
	// Retrieve checkbox filter parameters (multiple selections allowed).
	locations := r.URL.Query()["locations"]

	// Apply the filters to the dataset.
	results := filterAllData(creationMinStr, creationMaxStr, albumMinStr, albumMaxStr, membersMinStr, membersMaxStr, locations)

	// Return the filtered results as JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// filterAllData applies the provided filters to the artists data and returns matching results.
func filterAllData(creationMinStr, creationMaxStr, albumMinStr, albumMaxStr, membersMinStr, membersMaxStr string, locations []string) []models.SearchResult {
	var (
		creationMin, creationMax int
		albumMin, albumMax       int
		membersMin, membersMax   int
		err                      error
	)

	// Flags to check whether a given filter parameter is set.
	creationMinSet := false
	creationMaxSet := false
	albumMinSet := false
	albumMaxSet := false
	membersMinSet := false
	membersMaxSet := false

	// Parse creation date range.
	if creationMinStr != "" {
		creationMin, err = strconv.Atoi(creationMinStr)
		if err == nil {
			creationMinSet = true
		}
	}
	if creationMaxStr != "" {
		creationMax, err = strconv.Atoi(creationMaxStr)
		if err == nil {
			creationMaxSet = true
		}
	}

	// Parse first album date range.
	if albumMinStr != "" {
		albumMin, err = strconv.Atoi(albumMinStr)
		if err == nil {
			albumMinSet = true
		}
	}
	if albumMaxStr != "" {
		albumMax, err = strconv.Atoi(albumMaxStr)
		if err == nil {
			albumMaxSet = true
		}
	}

	// Parse number of members range.
	if membersMinStr != "" {
		membersMin, err = strconv.Atoi(membersMinStr)
		if err == nil {
			membersMinSet = true
		}
	}
	if membersMaxStr != "" {
		membersMax, err = strconv.Atoi(membersMaxStr)
		if err == nil {
			membersMaxSet = true
		}
	}

	// Prepare the locations filters (trim and convert to lowercase for case-insensitive matching).
	for i, loc := range locations {
		locations[i] = strings.ToLower(strings.TrimSpace(loc))
	}

	var results []models.SearchResult
	artists := dataStore.GetAllArtists() // Assume dataStore is a global variable providing the artist data.
	for _, artist := range artists {
		// --- Creation Date Filter ---
		if creationMinSet && artist.CreationDate < creationMin {
			continue
		}
		if creationMaxSet && artist.CreationDate > creationMax {
			continue
		}

		// --- First Album Date Filter ---
		// Assume artist.FirstAlbum is a string representing the year.
		albumYear, err := strconv.Atoi(artist.FirstAlbum)
		if err != nil {
			// If the album date cannot be parsed and a filter is set, skip this artist.
			if albumMinSet || albumMaxSet {
				continue
			}
		} else {
			if albumMinSet && albumYear < albumMin {
				continue
			}
			if albumMaxSet && albumYear > albumMax {
				continue
			}
		}

		// --- Number of Members Filter ---
		numMembers := len(artist.Members)
		if membersMinSet && numMembers < membersMin {
			continue
		}
		if membersMaxSet && numMembers > membersMax {
			continue
		}

		// --- Concert Locations Filter ---
		// If location filters are provided, the artist must have at least one matching location.
		if len(locations) > 0 {
			matched := false
			for _, loc := range artist.LocationsList {
				if containsIgnoreCase(loc, locations) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// Artist passes all filters; add it to the results.
		results = append(results, models.SearchResult{
			Text:        artist.Name,
			Type:        "artist",
			ArtistName:  artist.Name,
			Description: fmt.Sprintf("Band formed in %d with %d members", artist.CreationDate, numMembers),
			ArtistId:    artist.ID,
		})
	}

	return results
}

// containsIgnoreCase checks if the target string (in lowercase) matches any of the strings in arr.
func containsIgnoreCase(target string, arr []string) bool {
	targetLower := strings.ToLower(target)
	for _, s := range arr {
		if targetLower == s {
			return true
		}
	}
	return false
}
