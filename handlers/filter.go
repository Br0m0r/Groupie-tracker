package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"groupie/models"
)

type FilterParams struct {
	CreationYearMin int   `json:"creationYearMin"`
	CreationYearMax int   `json:"creationYearMax"`
	AlbumYearMin    int   `json:"albumYearMin"`
	AlbumYearMax    int   `json:"albumYearMax"`
	Members         []int `json:"members"`
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		ErrorHandler(w, ErrBadRequest, "Invalid request method")
		return
	}

	// Parse the filter parameters from the request body
	var params FilterParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		ErrorHandler(w, ErrBadRequest, "Invalid filter parameters")
		return
	}

	// Get all artists
	allArtists := dataStore.GetAllArtists()

	// Filter artists
	filteredArtists := make([]models.Artist, 0)
	for _, artist := range allArtists {
		// Check creation date range
		if artist.CreationDate < params.CreationYearMin ||
			artist.CreationDate > params.CreationYearMax {
			continue
		}

		// Parse and check first album year
		firstAlbumYear := parseYear(artist.FirstAlbum)
		if firstAlbumYear == 0 ||
			firstAlbumYear < params.AlbumYearMin ||
			firstAlbumYear > params.AlbumYearMax {
			continue
		}

		// Check member count if any member filters are selected
		if len(params.Members) > 0 {
			memberCount := len(artist.Members)
			isMatchingMembers := false
			for _, count := range params.Members {
				if count == 8 {
					// 8+ members
					if memberCount >= 8 {
						isMatchingMembers = true
						break
					}
				} else if memberCount == count {
					isMatchingMembers = true
					break
				}
			}
			if !isMatchingMembers {
				continue
			}
		}

		filteredArtists = append(filteredArtists, artist)
	}

	// Return filtered artists as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredArtists); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to encode response")
		return
	}
}

// Helper function to parse year from string
func parseYear(date string) int {
	// Extract 4-digit year from the date string
	yearStr := ""
	foundDigits := 0
	for _, char := range date {
		if char >= '0' && char <= '9' {
			yearStr += string(char)
			foundDigits++
			if foundDigits == 4 {
				break
			}
		}
	}

	// Convert to integer if we found a 4-digit year
	if len(yearStr) == 4 {
		year, err := strconv.Atoi(yearStr)
		if err == nil && year >= 1900 && year <= 2024 {
			return year
		}
	}

	return 0
}
