package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Structs for storing concert info
type Concert struct {
	Artist   string `json:"artist"`
	Location string `json:"location"`
	Lat      string `json:"lat,omitempty"`
	Lon      string `json:"lon,omitempty"`
}

type GeocodeResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

var concerts []Concert

func main() {
	// Load concert data
	data, err := ioutil.ReadFile("concerts.json")
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
	}

	if err := json.Unmarshal(data, &concerts); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Fetch coordinates
	for i, concert := range concerts {
		lat, lon, err := getCoordinates(concert.Location)
		if err != nil {
			log.Println("Error fetching coordinates for", concert.Location, ":", err)
			continue
		}
		concerts[i].Lat = lat
		concerts[i].Lon = lon
		time.Sleep(1 * time.Second) // Delay to avoid rate-limiting
	}

	// Save updated JSON
	jsonOutput, _ := json.MarshalIndent(concerts, "", "  ")
	err = ioutil.WriteFile("concerts_with_coordinates.json", jsonOutput, 0644)
	if err != nil {
		log.Fatal("Error writing output JSON file:", err)
	}

	fmt.Println("Concert data updated!")

	// Start a simple server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/api/concerts", serveConcertData)
	fmt.Println("Server started at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Serve JSON data via API
func serveConcertData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(concerts)
}

// Fetch coordinates from Nominatim API
func getCoordinates(location string) (string, string, error) {
	encodedLocation := url.QueryEscape(location)
	requestURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s", encodedLocation)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", "", err
	}
	// Set a proper User-Agent as required by Nominatim
	req.Header.Set("User-Agent", "concert-map-demo/1.0 (your_email@example.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil || len(results) == 0 {
		return "", "", fmt.Errorf("location not found")
	}

	return results[0].Lat, results[0].Lon, nil
}
