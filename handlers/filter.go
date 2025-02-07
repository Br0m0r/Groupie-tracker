// Package handlers provides HTTP request handlers for the web application
package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"groupie/models"
)

// FilterData represents all the data needed for the filter page
type FilterData struct {
	Artists         []models.ArtistCard
	UniqueLocations []string
	SelectedFilters FilterParams
	TotalResults    int
}

// FilterParams holds all possible filter parameters
type FilterParams struct {
	MemberCounts   []int    // Selected member counts
	Locations      []string // Selected locations
	CreationStart  int      // Creation date range start
	CreationEnd    int      // Creation date range end
	AlbumStartYear int      // First album year range start
	AlbumEndYear   int      // First album year range end
}

// getUniqueLocations extracts unique locations from all artists
func getUniqueLocations(artists []models.Artist) []string {
	locationMap := make(map[string]bool)

	for _, artist := range artists {
		for _, location := range artist.LocationsList {
			// Split location into parts (e.g., "London, UK" -> ["London", "UK"])
			parts := strings.Split(location, ", ")
			// Add each part as a unique location
			for _, part := range parts {
				locationMap[part] = true
			}
		}
	}

	// Convert map to sorted slice
	var locations []string
	for location := range locationMap {
		locations = append(locations, location)
	}
	return locations
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		ErrorHandler(w, ErrBadRequest, "Invalid form data")
		return
	}

	// Debug: Print all form values
	fmt.Printf("Received form data: %+v\n", r.Form)

	// Get all filter parameters
	params := FilterParams{
		MemberCounts:   getMemberCounts(r),
		Locations:      r.Form["location"],
		CreationStart:  parseIntDefault(r.FormValue("creation_start"), 1950),
		CreationEnd:    parseIntDefault(r.FormValue("creation_end"), 2024),
		AlbumStartYear: parseIntDefault(r.FormValue("album_start"), 1950),
		AlbumEndYear:   parseIntDefault(r.FormValue("album_end"), 2024),
	}

	// Debug: Print parsed parameters
	fmt.Printf("Parsed filter params: %+v\n", params)

	// Get all artists and apply filters
	allArtists := dataStore.GetAllArtists()
	filteredArtists := filterArtists(allArtists, params)

	// Debug: Print counts
	fmt.Printf("Total artists: %d, Filtered artists: %d\n", len(allArtists), len(filteredArtists))

	// Prepare data for template
	data := FilterData{
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

// Helper function to get selected member counts from form
func getMemberCounts(r *http.Request) []int {
	var counts []int
	for i := 1; i <= 8; i++ {
		if r.FormValue(fmt.Sprintf("members_%d", i)) != "" {
			counts = append(counts, i)
		}
	}
	return counts
}

// Helper function to parse int with default value
func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return val
}

// extractYear gets the year from a date string in format "DD-MM-YYYY"
func extractYear(date string) int {
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return 0
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}
	return year
}

// filterArtists applies all filters to the artist list
func filterArtists(artists []models.Artist, params FilterParams) []models.Artist {
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
func matchesFilters(artist models.Artist, params FilterParams) bool {
	// Debug print
	fmt.Printf("Checking artist %s against filters\n", artist.Name)

	// Check member count
	if len(params.MemberCounts) > 0 {
		memberCount := len(artist.Members)
		if memberCount > 8 {
			memberCount = 8 // Cap at 8 for "8+" option
		}
		found := false
		for _, count := range params.MemberCounts {
			if memberCount == count {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Artist %s filtered out by member count\n", artist.Name)
			return false
		}
	}

	// Check creation date range
	if artist.CreationDate < params.CreationStart || artist.CreationDate > params.CreationEnd {
		fmt.Printf("Artist %s filtered out by creation date\n", artist.Name)
		return false
	}

	// Check locations
	if len(params.Locations) > 0 {
		found := false
		for _, artistLocation := range artist.LocationsList {
			for _, filterLocation := range params.Locations {
				if strings.Contains(strings.ToLower(artistLocation), strings.ToLower(filterLocation)) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			fmt.Printf("Artist %s filtered out by location\n", artist.Name)
			return false
		}
	}

	// Check first album year
	albumYear := extractYear(artist.FirstAlbum)
	if albumYear < params.AlbumStartYear || albumYear > params.AlbumEndYear {
		fmt.Printf("Artist %s filtered out by album year\n", artist.Name)
		return false
	}

	fmt.Printf("Artist %s passed all filters\n", artist.Name)
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
