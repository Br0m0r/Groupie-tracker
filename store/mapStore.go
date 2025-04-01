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
	ds.mu.Lock()
	if ds.coordinateLoaderRunning {
		ds.mu.Unlock()
		return
	}
	ds.coordinateLoaderRunning = true
	ds.mu.Unlock()

	go func() {
		newCount := 0
		rateLimiter := time.NewTicker(1 * time.Second)
		defer rateLimiter.Stop()
		for _, location := range ds.UniqueLocations {
			<-rateLimiter.C
			ds.CoordinateCache.mu.RLock()
			_, exists := ds.CoordinateCache.data[location]
			ds.CoordinateCache.mu.RUnlock()
			if exists {
				continue
			}
			coords, err := ds.fetchCoordinatesFromAPI(location)
			if err != nil {
				log.Printf("Failed to fetch coordinates for %s: %v", location, err)
				continue
			}
			ds.CoordinateCache.mu.Lock()
			ds.CoordinateCache.data[location] = coords
			ds.CoordinateCache.mu.Unlock()
			log.Println("New location coordinates added:", location)
			newCount++
		}
		if newCount == 0 {
			log.Println("No new coordinates were fetched.")
		} else {
			log.Printf("Total new coordinates fetched: %d", newCount)
		}
		log.Println("Background coordinate loading completed")
		ds.mu.Lock()
		ds.coordinateLoaderRunning = false
		ds.mu.Unlock()
	}()
}

func (ds *DataStore) GetLocationCoordinates(location string) (models.Coordinates, error) {
	ds.CoordinateCache.mu.RLock()
	coords, exists := ds.CoordinateCache.data[location]
	ds.CoordinateCache.mu.RUnlock()
	if exists {
		return coords, nil
	}
	coords, err := ds.fetchCoordinatesFromAPI(location)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to fetch coordinates: %v", err)
	}
	ds.CoordinateCache.mu.Lock()
	ds.CoordinateCache.data[location] = coords
	ds.CoordinateCache.mu.Unlock()
	return coords, nil
}

func (ds *DataStore) fetchCoordinatesFromAPI(location string) (models.Coordinates, error) {
	encodedLocation := url.QueryEscape(location)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1", encodedLocation)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return models.Coordinates{}, err
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return models.Coordinates{}, err
	}
	defer resp.Body.Close()
	var nominatimResp []models.NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&nominatimResp); err != nil {
		return models.Coordinates{}, err
	}
	if len(nominatimResp) == 0 {
		return models.Coordinates{}, fmt.Errorf("no coordinates found for location")
	}
	lat, _ := strconv.ParseFloat(nominatimResp[0].Lat, 64)
	lon, _ := strconv.ParseFloat(nominatimResp[0].Lon, 64)
	return models.Coordinates{
		Lat:     lat,
		Lon:     lon,
		Address: location,
	}, nil
}
