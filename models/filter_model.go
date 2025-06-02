package models

// FilterParams holds all possible filter parameters
type FilterParams struct {
	MemberCounts   []int    // Selected member counts
	Locations      []string // Selected locations
	CreationStart  int      // Creation date range start
	CreationEnd    int      // Creation date range end
	AlbumStartYear int      // First album year range start
	AlbumEndYear   int      // First album year range end
}

// FilterData represents all the data needed for the filter page
type FilterData struct {
	Artists         []*Artist
	UniqueLocations []string
	SelectedFilters FilterParams
	TotalResults    int
	CurrentPath     string
	CurrentYear     int
}
