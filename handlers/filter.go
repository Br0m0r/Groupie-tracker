package handlers

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"groupie/models"
	"groupie/store"
	"groupie/utils"
)

// FilterHandler processes filter requests and returns filtered artists
func FilterHandler(dataStore *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			ErrorHandler(w, ErrBadRequest, "Invalid form data")
			return
		}

		// Extract filter parameters
		params := extractFilterParams(r)

		// Determine default params (for redirects and fast-path checks)
		defaultParams := getDefaultFilterParams(dataStore)

		// If nothing changed, redirect home
		if isDefaultParams(params, defaultParams) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Pull all artists once
		allArtists := dataStore.GetAllArtists()
		var filteredArtists []models.Artist

		// Fast‐path: ONLY locations selected (everything else at defaults)
		if len(params.Locations) > 0 &&
			len(params.MemberCounts) == 0 &&
			params.CreationStart == defaultParams.CreationStart &&
			params.CreationEnd == defaultParams.CreationEnd &&
			params.AlbumStartYear == defaultParams.AlbumStartYear &&
			params.AlbumEndYear == defaultParams.AlbumEndYear {

			// Use the locationMap for O(1) lookups
			artistSet := make(map[int]models.Artist)
			for _, loc := range params.Locations {
				for _, a := range dataStore.GetArtistsByLocation(loc) {
					artistSet[a.ID] = a
				}
			}

			// Flatten to slice
			filteredArtists = make([]models.Artist, 0, len(artistSet))
			for _, a := range artistSet {
				filteredArtists = append(filteredArtists, a)
			}

		} else if len(params.MemberCounts) > 0 &&
			len(params.Locations) == 0 &&
			params.CreationStart == defaultParams.CreationStart &&
			params.CreationEnd == defaultParams.CreationEnd &&
			params.AlbumStartYear == defaultParams.AlbumStartYear &&
			params.AlbumEndYear == defaultParams.AlbumEndYear {

			// Use memberCountMap for O(1) lookups
			artistSet := make(map[int]models.Artist)
			for _, c := range params.MemberCounts {
				for _, a := range dataStore.GetArtistsByMemberCount(c) {
					artistSet[a.ID] = a
				}
			}

			// Flatten to slice
			filteredArtists = make([]models.Artist, 0, len(artistSet))
			for _, a := range artistSet {
				filteredArtists = append(filteredArtists, a)
			}

		} else if len(params.MemberCounts) == 0 &&
			len(params.Locations) == 0 &&
			params.CreationStart == params.CreationEnd &&
			(params.CreationStart != defaultParams.CreationStart || params.CreationEnd != defaultParams.CreationEnd) &&
			params.AlbumStartYear == defaultParams.AlbumStartYear &&
			params.AlbumEndYear == defaultParams.AlbumEndYear {

			// Use creationYearMap for O(1) lookups
			filteredArtists = dataStore.GetArtistsByCreationYear(params.CreationStart)

		} else if len(params.MemberCounts) == 0 &&
			len(params.Locations) == 0 &&
			params.CreationStart == defaultParams.CreationStart &&
			params.CreationEnd == defaultParams.CreationEnd &&
			params.AlbumStartYear == params.AlbumEndYear &&
			(params.AlbumStartYear != defaultParams.AlbumStartYear || params.AlbumEndYear != defaultParams.AlbumEndYear) {

			// Use albumYearMap for O(1) lookups
			filteredArtists = dataStore.GetArtistsByAlbumYear(params.AlbumStartYear)

		} else {
			// Fallback: full in-memory filter scan
			filteredArtists = NewArtistFilter(params).Filter(allArtists)
		}

		// Prepare template data
		data := models.FilterData{
			Artists:         utils.ConvertToCards(filteredArtists),
			UniqueLocations: dataStore.UniqueLocations(),
			SelectedFilters: params,
			TotalResults:    len(filteredArtists),
			CurrentPath:     r.URL.Path,
			CurrentYear:     time.Now().Year(),
		}

		// Render
		if err := executeFilterTemplate(w, data); err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to process template")
			return
		}
	}
}

// isDefaultParams checks if user submitted filters are all defaults
func isDefaultParams(params, defaultParams models.FilterParams) bool {
	return len(params.MemberCounts) == 0 &&
		len(params.Locations) == 0 &&
		params.CreationStart == defaultParams.CreationStart &&
		params.CreationEnd == defaultParams.CreationEnd &&
		params.AlbumStartYear == defaultParams.AlbumStartYear &&
		params.AlbumEndYear == defaultParams.AlbumEndYear
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

// matches applies each criterion
func (af *ArtistFilter) matches(artist models.Artist) bool {
	return af.matchesMemberCount(artist) &&
		af.matchesCreationDate(artist) &&
		af.matchesLocation(artist) &&
		af.matchesAlbumYear(artist)
}

// matchesMemberCount checks member count filter
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

// matchesCreationDate checks creation date range
func (af *ArtistFilter) matchesCreationDate(artist models.Artist) bool {
	return artist.CreationDate >= af.params.CreationStart &&
		artist.CreationDate <= af.params.CreationEnd
}

// matchesLocation checks location filters
func (af *ArtistFilter) matchesLocation(artist models.Artist) bool {
	if len(af.params.Locations) == 0 {
		return true
	}

	for _, loc := range af.params.Locations {
		for _, artistLoc := range artist.LocationsList {
			if strings.EqualFold(artistLoc, loc) {
				return true
			}
		}
	}
	return false
}

// matchesAlbumYear checks album year range
func (af *ArtistFilter) matchesAlbumYear(artist models.Artist) bool {
	year := utils.ExtractYear(artist.FirstAlbum)
	return year >= af.params.AlbumStartYear &&
		year <= af.params.AlbumEndYear
}

// extractFilterParams reads FilterParams from the HTTP request
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
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// Return the default filter parameters.
func getDefaultFilterParams(dataStore *store.DataStore) models.FilterParams {
	minCreationYear, minFirstAlbum := dataStore.GetMinYears()
	return models.FilterParams{
		CreationStart:  minCreationYear,
		CreationEnd:    time.Now().Year(),
		AlbumStartYear: minFirstAlbum,
		AlbumEndYear:   time.Now().Year(),
		MemberCounts:   []int{},    // Empty slice - no members selected
		Locations:      []string{}, // Empty slice - no locations selected
	}
}
