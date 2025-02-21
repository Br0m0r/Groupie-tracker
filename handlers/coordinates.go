package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Coordinates struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Address string  `json:"address"`
}

type NominatimResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func GetLocationCoordinates(w http.ResponseWriter, r *http.Request) {
	// Extract artist ID from query params
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		ErrorHandler(w, ErrBadRequest, "Artist ID is required")
		return
	}

	// Get artist data
	fmt.Printf("Geocoding request received for artist ID: %s\n", artistID)
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

	// Create a slice to store all coordinates
	var coordinates []Coordinates

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Process each location
	fmt.Printf("Processing %d locations for artist\n", len(artist.LocationsList))
	for _, location := range artist.LocationsList {
		// Build Nominatim API URL
		encodedLocation := url.QueryEscape(location)
		apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1", encodedLocation)

		// Create request
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			continue
		}

		// Set User-Agent as required by Nominatim
		req.Header.Set("User-Agent", "GroupieTracker/1.0")

		// Make request
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		// Parse response
		var nominatimResp []NominatimResponse
		if err := json.NewDecoder(resp.Body).Decode(&nominatimResp); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// If we got coordinates, add them to our result
		if len(nominatimResp) > 0 {
			lat, _ := strconv.ParseFloat(nominatimResp[0].Lat, 64)
			lon, _ := strconv.ParseFloat(nominatimResp[0].Lon, 64)
			coordinates = append(coordinates, Coordinates{
				Lat:     lat,
				Lon:     lon,
				Address: location,
			})
		}

		// Respect Nominatim's usage policy with a delay
		time.Sleep(1 * time.Second)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coordinates)
}
