package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"groupie/models"
)

func (ds *DataStore) loadCoordinatesInBackground() {
	// Initialize the map
	ds.CoordinateCache.data = make(map[string]models.Coordinates)

	// Start background loading
	go func() {
		// Create rate limiter for Nominatim API (1 request per second)
		rateLimiter := time.NewTicker(1 * time.Second)
		defer rateLimiter.Stop()

		for _, location := range ds.UniqueLocations {
			// Wait for rate limiter
			<-rateLimiter.C

			// Check if we already have this location (could have been added by a request)
			ds.CoordinateCache.mu.RLock()
			_, exists := ds.CoordinateCache.data[location]
			ds.CoordinateCache.mu.RUnlock()

			if exists {
				continue
			}

			// Try to get coordinates
			coords, err := ds.fetchCoordinatesFromAPI(location)
			if err != nil {
				log.Printf("Failed to fetch coordinates for %s: %v", location, err)
				continue
			}

			// Store the coordinates
			ds.CoordinateCache.mu.Lock()
			ds.CoordinateCache.data[location] = coords
			ds.CoordinateCache.mu.Unlock()
		}
		log.Println("Background coordinate loading completed")
	}()
}

func (ds *DataStore) GetLocationCoordinates(location string) (models.Coordinates, error) {
	// First, try to get from cache
	ds.CoordinateCache.mu.RLock()
	coords, exists := ds.CoordinateCache.data[location]
	ds.CoordinateCache.mu.RUnlock()

	if exists {
		return coords, nil
	}

	// Not in cache, need to fetch
	coords, err := ds.fetchCoordinatesFromAPI(location)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to fetch coordinates: %v", err)
	}

	// Store in cache
	ds.CoordinateCache.mu.Lock()
	ds.CoordinateCache.data[location] = coords
	ds.CoordinateCache.mu.Unlock()

	return coords, nil
}

// Helper function to fetch coordinates from Nominatim
func (ds *DataStore) fetchCoordinatesFromAPI(location string) (models.Coordinates, error) {
	// Create the URL
	encodedLocation := url.QueryEscape(location)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1", encodedLocation)

	// Create request with required headers
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return models.Coordinates{}, err
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	// Make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return models.Coordinates{}, err
	}
	defer resp.Body.Close()

	// Parse response
	var nominatimResp []models.NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&nominatimResp); err != nil {
		return models.Coordinates{}, err
	}

	if len(nominatimResp) == 0 {
		return models.Coordinates{}, fmt.Errorf("no coordinates found for location")
	}

	// Convert response to our format
	lat, _ := strconv.ParseFloat(nominatimResp[0].Lat, 64)
	lon, _ := strconv.ParseFloat(nominatimResp[0].Lon, 64)

	return models.Coordinates{
		Lat:     lat,
		Lon:     lon,
		Address: location,
	}, nil
}
