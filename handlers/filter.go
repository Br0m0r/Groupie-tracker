package handlers

import (
	"net/http"
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

		// Validate filter parameters
		if err := params.Validate(); err != nil {
			ErrorHandler(w, ErrBadRequest, err.Error())
			return
		}

		// Get default params for comparison
		defaultParams := store.GetDefaultFilterParams(dataStore)

		// If nothing changed, redirect home
		if isDefaultParams(params, defaultParams) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Get all artists (now returns []*models.Artist)
		allArtists := dataStore.GetAllArtists()
		var filteredArtists []*models.Artist

		// Fast-path optimizations using our optimized repository lookups
		if len(params.Locations) > 0 &&
			len(params.MemberCounts) == 0 &&
			params.CreationStart == defaultParams.CreationStart &&
			params.CreationEnd == defaultParams.CreationEnd &&
			params.AlbumStartYear == defaultParams.AlbumStartYear &&
			params.AlbumEndYear == defaultParams.AlbumEndYear {

			// Fast-path: ONLY locations selected - use O(1) lookups
			artistSet := make(map[int]*models.Artist)
			for _, loc := range params.Locations {
				for _, artist := range dataStore.GetArtistsByLocation(loc) {
					artistSet[artist.ID] = artist
				}
			}

			// Convert map to slice
			filteredArtists = make([]*models.Artist, 0, len(artistSet))
			for _, artist := range artistSet {
				filteredArtists = append(filteredArtists, artist)
			}

		} else if len(params.MemberCounts) > 0 &&
			len(params.Locations) == 0 &&
			params.CreationStart == defaultParams.CreationStart &&
			params.CreationEnd == defaultParams.CreationEnd &&
			params.AlbumStartYear == defaultParams.AlbumStartYear &&
			params.AlbumEndYear == defaultParams.AlbumEndYear {

			// Fast-path: ONLY member counts selected
			artistSet := make(map[int]*models.Artist)
			for _, count := range params.MemberCounts {
				for _, artist := range dataStore.GetArtistsByMemberCount(count) {
					artistSet[artist.ID] = artist
				}
			}

			// Convert map to slice
			filteredArtists = make([]*models.Artist, 0, len(artistSet))
			for _, artist := range artistSet {
				filteredArtists = append(filteredArtists, artist)
			}

		} else if len(params.MemberCounts) == 0 &&
			len(params.Locations) == 0 &&
			params.CreationStart == params.CreationEnd &&
			(params.CreationStart != defaultParams.CreationStart || params.CreationEnd != defaultParams.CreationEnd) &&
			params.AlbumStartYear == defaultParams.AlbumStartYear &&
			params.AlbumEndYear == defaultParams.AlbumEndYear {

			// Fast-path: ONLY creation year selected
			filteredArtists = dataStore.GetArtistsByCreationYear(params.CreationStart)
		} else if len(params.MemberCounts) == 0 &&
			len(params.Locations) == 0 &&
			params.CreationStart == defaultParams.CreationStart &&
			params.CreationEnd == defaultParams.CreationEnd &&
			params.AlbumStartYear == params.AlbumEndYear &&
			(params.AlbumStartYear != defaultParams.AlbumStartYear || params.AlbumEndYear != defaultParams.AlbumEndYear) {

			// Fast-path: ONLY album year selected
			filteredArtists = dataStore.GetArtistsByAlbumYear(params.AlbumStartYear)
		} else {
			// Fallback: full in-memory filter scan
			filteredArtists = NewArtistFilter(params).Filter(allArtists)
		}

		// Prepare template data
		data := models.FilterData{
			Artists:         utils.ConvertToCards(filteredArtists), // Updated function call
			UniqueLocations: dataStore.UniqueLocations(),
			SelectedFilters: params,
			TotalResults:    len(filteredArtists),
			CurrentPath:     r.URL.Path,
			CurrentYear:     time.Now().Year(),
		}

		// Render template
		if err := utils.ExecuteFilterTemplate(w, data); err != nil {
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

// ArtistFilter encapsulates filtering logic for fallback filtering
type ArtistFilter struct {
	params models.FilterParams
}

// NewArtistFilter creates a new filter with given parameters
func NewArtistFilter(params models.FilterParams) *ArtistFilter {
	return &ArtistFilter{params: params}
}

// Filter applies all filters to the artist list (updated for pointers)
func (af *ArtistFilter) Filter(artists []*models.Artist) []*models.Artist {
	var filtered []*models.Artist
	for _, artist := range artists {
		if af.matches(artist) {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// matches applies each criterion (updated for pointers)
func (af *ArtistFilter) matches(artist *models.Artist) bool {
	return af.matchesMemberCount(artist) &&
		af.matchesCreationDate(artist) &&
		af.matchesLocation(artist) &&
		af.matchesAlbumYear(artist)
}

// matchesMemberCount checks member count filter (updated for pointers)
func (af *ArtistFilter) matchesMemberCount(artist *models.Artist) bool {
	if len(af.params.MemberCounts) == 0 {
		return true
	}

	memberCount := artist.GetMemberCount()
	for _, count := range af.params.MemberCounts {
		if memberCount == count {
			return true
		}
	}
	return false
}

// matchesCreationDate checks creation date range (updated for pointers)
func (af *ArtistFilter) matchesCreationDate(artist *models.Artist) bool {
	return artist.CreationDate >= af.params.CreationStart &&
		artist.CreationDate <= af.params.CreationEnd
}

// matchesLocation checks location filters (updated for pointers)
func (af *ArtistFilter) matchesLocation(artist *models.Artist) bool {
	if len(af.params.Locations) == 0 {
		return true
	}

	for _, loc := range af.params.Locations {
		if artist.HasLocation(loc) {
			return true
		}
	}
	return false
}

// matchesAlbumYear checks album year range (updated for pointers)
func (af *ArtistFilter) matchesAlbumYear(artist *models.Artist) bool {
	year := artist.GetFirstAlbumYear()
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
