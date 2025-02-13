// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"groupie/models"
	"groupie/utils"
)

// getUniqueLocations extracts unique locations from all artists
func getUniqueLocations(artists []models.Artist) []string {
	locationMap := make(map[string]bool)

	for _, artist := range artists {
		for _, location := range artist.LocationsList {
			locationMap[location] = true
		}
	}

	// Convert map to sorted slice
	var locations []string
	for location := range locationMap {
		locations = append(locations, location)
	}
	sort.Strings(locations)
	return locations
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		ErrorHandler(w, ErrBadRequest, "Invalid form data")
		return
	}

	// Get all filter parameters
	params := models.FilterParams{
		MemberCounts:   utils.GetMemberCounts(r),
		Locations:      r.Form["location"],
		CreationStart:  utils.ParseIntDefault(r.FormValue("creation_start"), 1950),
		CreationEnd:    utils.ParseIntDefault(r.FormValue("creation_end"), 2024),
		AlbumStartYear: utils.ParseIntDefault(r.FormValue("album_start"), 1950),
		AlbumEndYear:   utils.ParseIntDefault(r.FormValue("album_end"), 2024),
	}

	// Get all artists and apply filters
	allArtists := dataStore.GetAllArtists()
	filteredArtists := filterArtists(allArtists, params)

	// Prepare data for template
	data := models.FilterData{
		Artists:         convertToCards(filteredArtists),
		UniqueLocations: getUniqueLocations(allArtists),
		SelectedFilters: params,
		TotalResults:    len(filteredArtists),
	}
	// Parse and execute template with functions
	funcMap := template.FuncMap{
		"iterate": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"contains": func(slice interface{}, item interface{}) bool {
			switch slice := slice.(type) {
			case []int:
				itemInt, ok := item.(int)
				if !ok {
					return false
				}
				for _, s := range slice {
					if s == itemInt {
						return true
					}
				}
			case []string:
				itemStr, ok := item.(string)
				if !ok {
					return false
				}
				for _, s := range slice {
					if s == itemStr {
						return true
					}
				}
			}
			return false
		},
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		return
	}
}

// filterArtists applies all filters to the artist list
func filterArtists(artists []models.Artist, params models.FilterParams) []models.Artist {
	var filtered []models.Artist

	for _, artist := range artists {
		if !matchesFilters(artist, params) {
			continue
		}
		filtered = append(filtered, artist)
	}

	return filtered
}

// matchesFilters checks if an artist matches all filter criteria
func matchesFilters(artist models.Artist, params models.FilterParams) bool {
	fmt.Printf("\nChecking artist: %s\n", artist.Name)

	// Member count check remains the same
	if len(params.MemberCounts) > 0 {
		memberCount := len(artist.Members)
		if memberCount > 8 {
			memberCount = 8
		}
		found := false
		for _, count := range params.MemberCounts {
			if memberCount == count {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Failed member count check\n")
			return false
		}
	}

	// Creation date check remains the same
	if artist.CreationDate < params.CreationStart || artist.CreationDate > params.CreationEnd {
		return false
	}

	// Location check - with debug prints
	if len(params.Locations) > 0 {
		found := false
		for _, filterLocation := range params.Locations {
			// First check LocationsList
			for _, artistLocation := range artist.LocationsList {
				if strings.Contains(strings.ToLower(artistLocation), strings.ToLower(filterLocation)) {
					found = true
					break
				}
			}

			// Then check LocationData if not found yet
			if !found {
				// Check if it's a state
				if _, exists := utils.StateCityMap[filterLocation]; exists {
					if cities, ok := artist.LocationData[filterLocation]; ok && len(cities) > 0 {
						found = true
					}
				}
			}

			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	// Album year check remains the same
	albumYear := utils.ExtractYear(artist.FirstAlbum)
	if albumYear < params.AlbumStartYear || albumYear > params.AlbumEndYear {
		return false
	}

	return true
}

// Helper function to convert Artists to ArtistCards
func convertToCards(artists []models.Artist) []models.ArtistCard {
	cards := make([]models.ArtistCard, len(artists))
	for i, artist := range artists {
		cards[i] = models.ArtistCard{
			ID:    artist.ID,
			Name:  artist.Name,
			Image: artist.Image,
		}
	}
	return cards
}
