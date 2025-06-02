package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"groupie/models"
)

// CoordinatesRepository manages geolocation data for concert locations
type CoordinatesRepository struct {
	cache map[string]*models.Coordinates // Store pointers for consistency
	mu    sync.RWMutex
}

// NewCoordinatesRepository creates a new instance
func NewCoordinatesRepository() *CoordinatesRepository {
	return &CoordinatesRepository{
		cache: make(map[string]*models.Coordinates),
	}
}

// Has reports whether the given location is already cached
func (cr *CoordinatesRepository) Has(location string) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	_, ok := cr.cache[location]
	return ok
}

// CacheSize returns the number of cached entries
func (cr *CoordinatesRepository) CacheSize() int {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	return len(cr.cache)
}

// Get retrieves coordinates for a location, fetching from API if not cached
func (cr *CoordinatesRepository) Get(location string) (*models.Coordinates, error) {
	// Check cache first
	cr.mu.RLock()
	coords, exists := cr.cache[location]
	cr.mu.RUnlock()

	if exists {
		return coords, nil
	}

	// Not in cache, fetch from API
	coords, err := cr.fetchFromAPI(location)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coordinates for %s: %w", location, err)
	}

	// Validate before caching

	// Cache the result
	cr.mu.Lock()
	cr.cache[location] = coords
	cr.mu.Unlock()

	return coords, nil
}

// PrefetchLocations starts background loading of coordinates
func (cr *CoordinatesRepository) PrefetchLocations(locations []string) {
	go func() {
		// Rate limiter for Nominatim API (1 request per second)
		rateLimiter := time.NewTicker(1 * time.Second)
		defer rateLimiter.Stop()

		for _, location := range locations {
			// Wait for rate limiter
			<-rateLimiter.C

			// Skip if already cached
			if cr.Has(location) {
				continue
			}

			// Fetch coordinates
			coords, err := cr.fetchFromAPI(location)
			if err != nil {
				//		log.Printf("Failed to prefetch coordinates for %s: %v", location, err)
				continue
			}

			cr.mu.Lock()
			cr.cache[location] = coords
			cr.mu.Unlock()

			//	log.Printf("Prefetched coordinates for: %s", location)
		}
		log.Println("Background coordinate prefetching completed")
	}()
}

// fetchFromAPI retrieves coordinates from Nominatim API
func (cr *CoordinatesRepository) fetchFromAPI(location string) (*models.Coordinates, error) {
	// Build API URL
	encodedLocation := url.QueryEscape(location)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1", encodedLocation)

	// Create request with headers
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	// Make request with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse response
	var nominatimResp []models.NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&nominatimResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(nominatimResp) == 0 {
		return nil, fmt.Errorf("no coordinates found for location: %s", location)
	}

	// Convert to our Coordinates model
	coords, err := nominatimResp[0].ToCoordinates(location)
	if err != nil {
		return nil, fmt.Errorf("failed to convert coordinates: %w", err)
	}

	return coords, nil
}

// GetAllCached returns all cached coordinates (for debugging/monitoring)
func (cr *CoordinatesRepository) GetAllCached() map[string]*models.Coordinates {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	// Return a copy of the map to prevent external modifications
	cached := make(map[string]*models.Coordinates, len(cr.cache))
	for k, v := range cr.cache {
		cached[k] = v
	}
	return cached
}

// ClearCache clears all cached coordinates (useful for testing)
func (cr *CoordinatesRepository) ClearCache() {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.cache = make(map[string]*models.Coordinates)
}
