package models

// FilterData represents all the data needed for the filter page
type FilterData struct {
	Artists         []ArtistCard
	UniqueLocations []string
	SelectedFilters FilterParams
	TotalResults    int
	CurrentPath     string
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
