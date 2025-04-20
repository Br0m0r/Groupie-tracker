package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"groupie/models"
)

// CoordinatesRepositoryImpl manages the geolocation data for concert locations
type CoordinatesRepositoryImpl struct {
	cache map[string]models.Coordinates
	mu    sync.RWMutex
}

// NewCoordinatesRepository creates a new instance of the CoordinatesRepository
func NewCoordinatesRepository() CoordinatesRepository {
	return &CoordinatesRepositoryImpl{
		cache: make(map[string]models.Coordinates),
	}
}

// Get retrieves coordinates for a location, fetching from API if not cached
func (cr *CoordinatesRepositoryImpl) Get(location string) (models.Coordinates, error) {
	// First, try to get from cache
	cr.mu.RLock()
	coords, exists := cr.cache[location]
	cr.mu.RUnlock()

	if exists {
		return coords, nil
	}

	// Not in cache, need to fetch
	coords, err := cr.fetchFromAPI(location)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to fetch coordinates: %v", err)
	}

	// Store in cache
	cr.mu.Lock()
	cr.cache[location] = coords
	cr.mu.Unlock()

	return coords, nil
}

// PrefetchLocations starts background loading of coordinates for a set of locations
func (cr *CoordinatesRepositoryImpl) PrefetchLocations(locations []string) {
	// Start background loading
	go func() {
		// Create rate limiter for Nominatim API (1 request per second)
		rateLimiter := time.NewTicker(1 * time.Second)
		defer rateLimiter.Stop()

		for _, location := range locations {
			// Wait for rate limiter
			<-rateLimiter.C

			// Check if we already have this location (could have been added by a request)
			cr.mu.RLock()
			_, exists := cr.cache[location]
			cr.mu.RUnlock()

			if exists {
				continue
			}

			// Try to get coordinates
			coords, err := cr.fetchFromAPI(location)
			if err != nil {
				log.Printf("Failed to fetch coordinates for %s: %v", location, err)
				continue
			}

			// Store the coordinates
			cr.mu.Lock()
			cr.cache[location] = coords
			cr.mu.Unlock()
			log.Println("New location coordinates added:", location)
		}
		log.Println("Background coordinate loading completed")
	}()
}

// Helper function to fetch coordinates from Nominatim
func (cr *CoordinatesRepositoryImpl) fetchFromAPI(location string) (models.Coordinates, error) {
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

// ImportCache merges coordinates from another repository
func (cr *CoordinatesRepositoryImpl) ImportCache(other CoordinatesRepository) {
	if other == nil {
		return
	}

	// We need to cast to the concrete type to access the cache
	otherImpl, ok := other.(*CoordinatesRepositoryImpl)
	if !ok {
		return // Not the same implementation type
	}

	otherImpl.mu.RLock()
	defer otherImpl.mu.RUnlock()

	cr.mu.Lock()
	defer cr.mu.Unlock()

	// Copy coordinates from other repository
	for loc, coord := range otherImpl.cache {
		cr.cache[loc] = coord
	}
}
