package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"groupie/models"
)

func GetLocationCoordinates(w http.ResponseWriter, r *http.Request) {
	// Extract artist ID from query params
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		ErrorHandler(w, ErrBadRequest, "Artist ID is required")
		return
	}

	// Convert artist ID to integer
	id, err := strconv.Atoi(artistID)
	if err != nil {
		ErrorHandler(w, ErrInvalidID, "Invalid artist ID format")
		return
	}

	// Get artist data
	artist, err := dataStore.GetArtist(id)
	if err != nil {
		ErrorHandler(w, ErrNotFound, "Artist not found")
		return
	}

	// Create a slice to store all coordinates
	var coordinates []models.Coordinates

	// Process each location
	for _, location := range artist.LocationsList {
		// Try to get coordinates from repository through the dataStore
		coords, err := dataStore.GetLocationCoordinates(location)
		if err != nil {
			fmt.Printf("Error getting coordinates for %s: %v\n", location, err)
			continue // Skip this location if coordinates cannot be retrieved
		}
		coordinates = append(coordinates, coords)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coordinates)
}
