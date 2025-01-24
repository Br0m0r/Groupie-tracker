package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"groupie/models"
)

type FilterParams struct {
	Creation struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"creation"`
	Album struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"album"`
	Members   []int    `json:"members"`
	Locations []string `json:"locations"`
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, ErrBadRequest, "Only POST method is allowed")
		return
	}

	var params FilterParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		ErrorHandler(w, ErrBadRequest, "Invalid filter parameters")
		return
	}

	// Set default ranges if not provided
	if params.Creation.Min == 0 && params.Creation.Max == 0 {
		params.Creation.Min = 1950
		params.Creation.Max = 2024
	}
	if params.Album.Min == 0 && params.Album.Max == 0 {
		params.Album.Min = 1950
		params.Album.Max = 2024
	}

	allArtists := dataStore.GetAllArtists()
	filteredArtists := make([]models.ArtistCard, 0)

	for _, artist := range allArtists {
		if !isInRange(artist.CreationDate, params.Creation.Min, params.Creation.Max) {
			continue
		}

		albumYear := parseAlbumYear(artist.FirstAlbum)
		if !isInRange(albumYear, params.Album.Min, params.Album.Max) {
			continue
		}

		if len(params.Members) > 0 && !contains(params.Members, len(artist.Members)) {
			continue
		}

		if len(params.Locations) > 0 && !hasMatchingLocation(artist.LocationsList, params.Locations) {
			continue
		}

		filteredArtists = append(filteredArtists, models.ArtistCard{
			ID:    artist.ID,
			Name:  artist.Name,
			Image: artist.Image,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredArtists); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to encode response")
		return
	}
}

func isInRange(value, min, max int) bool {
	return value >= min && value <= max
}

func contains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func hasMatchingLocation(artistLocations, filterLocations []string) bool {
	for _, filterLoc := range filterLocations {
		for _, artistLoc := range artistLocations {
			if strings.Contains(strings.ToLower(artistLoc), strings.ToLower(filterLoc)) {
				return true
			}
		}
	}
	return false
}

func parseAlbumYear(albumDate string) int {
	if len(albumDate) < 4 {
		return 1960
	}
	yearStr := albumDate[len(albumDate)-4:]
	var year int
	if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
		return 1960
	}
	return year
}
