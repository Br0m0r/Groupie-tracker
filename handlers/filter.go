package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"groupie/models"
	"groupie/utils"
)

// FilterHandler processes filter requests and returns filtered artists
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ErrorHandler(w, ErrBadRequest, "Invalid form data")
		return
	}

	// Extract filter parameters
	params := extractFilterParams(r)

	// Get and filter artists
	allArtists := dataStore.GetAllArtists()
	filteredArtists := NewArtistFilter(params).Filter(allArtists)

	// Prepare template data
	data := models.FilterData{
		Artists:         utils.ConvertToCards(filteredArtists),
		UniqueLocations: utils.GetUniqueLocations(allArtists),
		SelectedFilters: params,
		TotalResults:    len(filteredArtists),
		CurrentPath:     r.URL.Path,
	}

	// Execute template
	if err := executeFilterTemplate(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to process template")
		return
	}
}

// ArtistFilter encapsulates filtering logic
type ArtistFilter struct {
	params models.FilterParams
}

// NewArtistFilter creates a new filter with given parameters
func NewArtistFilter(params models.FilterParams) *ArtistFilter {
	return &ArtistFilter{params: params}
}

// Filter applies all filters to the artist list
func (af *ArtistFilter) Filter(artists []models.Artist) []models.Artist {
	var filtered []models.Artist
	for _, artist := range artists {
		if af.matches(artist) {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// matches checks if an artist matches all filter criteria
func (af *ArtistFilter) matches(artist models.Artist) bool {
	return af.matchesMemberCount(artist) &&
		af.matchesCreationDate(artist) &&
		af.matchesLocation(artist) &&
		af.matchesAlbumYear(artist)
}

// matchesMemberCount checks if artist matches member count filter
func (af *ArtistFilter) matchesMemberCount(artist models.Artist) bool {
	if len(af.params.MemberCounts) == 0 {
		return true
	}

	memberCount := len(artist.Members)
	if memberCount > 8 {
		memberCount = 8
	}

	for _, count := range af.params.MemberCounts {
		if memberCount == count {
			return true
		}
	}
	return false
}

// matchesCreationDate checks if artist matches creation date range
func (af *ArtistFilter) matchesCreationDate(artist models.Artist) bool {
	return artist.CreationDate >= af.params.CreationStart &&
		artist.CreationDate <= af.params.CreationEnd
}

// matchesLocation checks if artist matches location filters
func (af *ArtistFilter) matchesLocation(artist models.Artist) bool {
	if len(af.params.Locations) == 0 {
		return true
	}

	for _, filterLocation := range af.params.Locations {
		// Check direct location matches
		for _, artistLocation := range artist.LocationsList {
			if strings.Contains(
				strings.ToLower(artistLocation),
				strings.ToLower(filterLocation),
			) {
				return true
			}
		}

		// Check state-city relationships
		if cities, isState := utils.StateCityMap[filterLocation]; isState {
			for _, city := range cities {
				if af.hasLocation(artist, city) {
					return true
				}
			}
		}
	}
	return false
}

// hasLocation checks if artist has a specific location
func (af *ArtistFilter) hasLocation(artist models.Artist, location string) bool {
	for _, artistLocation := range artist.LocationsList {
		if artistLocation == location {
			return true
		}
	}
	return false
}

// matchesAlbumYear checks if artist matches album year range
func (af *ArtistFilter) matchesAlbumYear(artist models.Artist) bool {
	albumYear := utils.ExtractYear(artist.FirstAlbum)
	return albumYear >= af.params.AlbumStartYear &&
		albumYear <= af.params.AlbumEndYear
}

// Helper function to extract filter parameters from request
func extractFilterParams(r *http.Request) models.FilterParams {
	return models.FilterParams{
		MemberCounts:   utils.GetMemberCounts(r),
		Locations:      r.Form["location"],
		CreationStart:  utils.ParseIntDefault(r.FormValue("creation_start"), 1950),
		CreationEnd:    utils.ParseIntDefault(r.FormValue("creation_end"), 2024),
		AlbumStartYear: utils.ParseIntDefault(r.FormValue("album_start"), 1950),
		AlbumEndYear:   utils.ParseIntDefault(r.FormValue("album_end"), 2024),
	}
}

// Helper function to execute the filter template
func executeFilterTemplate(w http.ResponseWriter, data models.FilterData) error {
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
		return err
	}

	return tmpl.Execute(w, data)
}
