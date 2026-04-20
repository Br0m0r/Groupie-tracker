package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"groupie/models"
)

func GetLocationCoordinates(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		ErrorHandler(w, ErrBadRequest, "Artist ID is required")
		return
	}

	id, err := strconv.Atoi(artistID)
	if err != nil {
		ErrorHandler(w, ErrInvalidID, "Invalid artist ID format")
		return
	}

	artist, err := dataStore.GetArtist(id)
	if err != nil {
		ErrorHandler(w, ErrNotFound, "Artist not found")
		return
	}

	var coordinates []models.Coordinates

	for _, location := range artist.LocationsList {
		// Try to get coordinates from cache/API
		coords, err := dataStore.GetLocationCoordinates(location)
		if err != nil {
			log.Printf("Error getting coordinates for %s: %v", location, err)
			continue
		}
		coordinates = append(coordinates, coords)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coordinates)
}
