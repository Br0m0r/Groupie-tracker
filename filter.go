package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"groupie/models"
)

// FilterHandler processes filtering of artists based on user input
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	creationFrom := r.URL.Query().Get("creation_from")
	creationTo := r.URL.Query().Get("creation_to")
	firstAlbumFrom := r.URL.Query().Get("album_from")
	firstAlbumTo := r.URL.Query().Get("album_to")
	memberCount := r.URL.Query().Get("members")
	locations := r.URL.Query()["locations"] // Can be multiple values

	// Get all artists from the data store
	artists := dataStore.GetAllArtists()
	var filteredArtists []models.Artist

	// Iterate over artists and apply filters
	for _, artist := range artists {
		if !filterByCreationDate(artist, creationFrom, creationTo) {
			continue
		}
		if !filterByFirstAlbum(artist, firstAlbumFrom, firstAlbumTo) {
			continue
		}
		if !filterByMembers(artist, memberCount) {
			continue
		}
		if !filterByLocations(artist, locations) {
			continue
		}

		// If artist passes all filters, add to results
		filteredArtists = append(filteredArtists, artist)
	}

	// Respond with filtered results in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredArtists)
}

// Helper function to filter by creation date (range)
func filterByCreationDate(artist models.Artist, from, to string) bool {
	if from == "" && to == "" {
		return true
	}
	creationYear := artist.CreationDate
	fromYear, _ := strconv.Atoi(from)
	toYear, _ := strconv.Atoi(to)

	if from != "" && creationYear < fromYear {
		return false
	}
	if to != "" && creationYear > toYear {
		return false
	}
	return true
}

// Helper function to filter by first album date (range)
func filterByFirstAlbum(artist models.Artist, from, to string) bool {
	if from == "" && to == "" {
		return true
	}

	// Extract the album year
	albumYear := extractYear(artist.FirstAlbum)
	fromYear, _ := strconv.Atoi(from)
	toYear, _ := strconv.Atoi(to)

	if from != "" && albumYear < fromYear {
		return false
	}
	if to != "" && albumYear > toYear {
		return false
	}
	return true
}

// Helper function to extract the year from album release string (e.g., "1995-06-20" -> 1995)
func extractYear(album string) int {
	parts := strings.Split(album, "-")
	year, _ := strconv.Atoi(parts[0])
	return year
}

// Helper function to filter by number of members
func filterByMembers(artist models.Artist, count string) bool {
	if count == "" {
		return true
	}
	expectedCount, _ := strconv.Atoi(count)
	return len(artist.Members) == expectedCount
}

// Helper function to filter by locations (checkbox filter)
func filterByLocations(artist models.Artist, selectedLocations []string) bool {
	if len(selectedLocations) == 0 {
		return true
	}

	for _, location := range artist.LocationsList {
		for _, selected := range selectedLocations {
			if strings.EqualFold(location, selected) {
				return true
			}
		}
	}
	return false
}
