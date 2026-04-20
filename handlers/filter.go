package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"groupie/models"
	"groupie/utils"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ErrorHandler(w, ErrBadRequest, "Invalid form data")
		return
	}

	// Extract filter parameters
	params := extractFilterParams(r)

	// Check if params match default params
	defaultParams := utils.GetDefaultFilterParams()
	if isDefaultParams(params, defaultParams) {
		// Redirect to home page instead of processing the filter
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Continue with filtering only if params differ from default
	allArtists := dataStore.GetAllArtists()
	filteredArtists := NewArtistFilter(params).Filter(allArtists)

	data := models.FilterData{
		Artists:         utils.ConvertToCards(filteredArtists),
		UniqueLocations: dataStore.UniqueLocations,
		SelectedFilters: params,
		TotalResults:    len(filteredArtists),
		CurrentPath:     r.URL.Path,
	}

	if err := executeFilterTemplate(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to process template")
		return
	}
}

func isDefaultParams(params, defaultParams models.FilterParams) bool {
	return len(params.MemberCounts) == 0 &&
		len(params.Locations) == 0 &&
		params.CreationStart == defaultParams.CreationStart &&
		params.CreationEnd == defaultParams.CreationEnd &&
		params.AlbumStartYear == defaultParams.AlbumStartYear &&
		params.AlbumEndYear == defaultParams.AlbumEndYear
}

type ArtistFilter struct {
	params models.FilterParams
}

func NewArtistFilter(params models.FilterParams) *ArtistFilter {
	return &ArtistFilter{params: params}
}

func (af *ArtistFilter) Filter(artists []models.Artist) []models.Artist {
	var filtered []models.Artist
	for _, artist := range artists {
		if af.matches(artist) {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

func (af *ArtistFilter) matches(artist models.Artist) bool {
	return af.matchesMemberCount(artist) &&
		af.matchesCreationDate(artist) &&
		af.matchesLocation(artist) &&
		af.matchesAlbumYear(artist)
}

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

func (af *ArtistFilter) matchesCreationDate(artist models.Artist) bool {
	return artist.CreationDate >= af.params.CreationStart &&
		artist.CreationDate <= af.params.CreationEnd
}

func (af *ArtistFilter) matchesLocation(artist models.Artist) bool {
	if len(af.params.Locations) == 0 {
		return true
	}

	for _, filterLocation := range af.params.Locations {
		filterLocationLower := strings.ToLower(filterLocation)

		for _, artistLocation := range artist.LocationsList {
			if strings.Contains(strings.ToLower(artistLocation), filterLocationLower) {
				return true
			}
		}

		if citiesInState, exists := artist.LocationStatesCities[filterLocation]; exists && len(citiesInState) > 0 {
			return true
		}
	}

	return false
}

func (af *ArtistFilter) matchesAlbumYear(artist models.Artist) bool {
	albumYear := utils.ExtractYear(artist.FirstAlbum)
	return albumYear >= af.params.AlbumStartYear &&
		albumYear <= af.params.AlbumEndYear
}

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, data)
}
