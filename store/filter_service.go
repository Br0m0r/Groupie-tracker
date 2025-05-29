package store

import "groupie/models"

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

// IsDefaultParams checks if user submitted filters are all defaults
func IsDefaultParams(params, defaultParams models.FilterParams) bool {
	return len(params.MemberCounts) == 0 &&
		len(params.Locations) == 0 &&
		params.CreationStart == defaultParams.CreationStart &&
		params.CreationEnd == defaultParams.CreationEnd &&
		params.AlbumStartYear == defaultParams.AlbumStartYear &&
		params.AlbumEndYear == defaultParams.AlbumEndYear
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
