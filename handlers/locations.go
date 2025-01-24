// handlers/locations.go
package handlers

import (
	"encoding/json"
	"net/http"
)

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, ErrBadRequest, "Only GET method is allowed")
		return
	}

	locations := make(map[string]bool)
	for _, artist := range dataStore.GetAllArtists() {
		for _, loc := range artist.LocationsList {
			locations[loc] = true
		}
	}

	uniqueLocations := make([]string, 0, len(locations))
	for loc := range locations {
		uniqueLocations = append(uniqueLocations, loc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uniqueLocations)
}
