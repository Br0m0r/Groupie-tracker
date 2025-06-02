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
			ErrorHandler(w, models.ErrBadRequest, "Invalid form data")
			return
		}

		// Extract filter parameters
		params := utils.ExtractFilterParams(r)

		// Get default params for comparison
		defaultParams := store.GetDefaultFilterParams(dataStore)

		// If nothing changed, redirect home
		if store.IsDefaultParams(params, defaultParams) {
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
			filteredArtists = store.NewArtistFilter(params).Filter(allArtists)
		}

		// Prepare template data
		data := models.FilterData{
			Artists:         filteredArtists,
			UniqueLocations: dataStore.UniqueLocations(),
			SelectedFilters: params,
			TotalResults:    len(filteredArtists),
			CurrentPath:     r.URL.Path,
			CurrentYear:     time.Now().Year(),
		}

		// Render template
		if err := utils.ExecuteFilterTemplate(w, data); err != nil {
			ErrorHandler(w, models.ErrInternalServer, "Failed to process template")
			return
		}
	}
}
