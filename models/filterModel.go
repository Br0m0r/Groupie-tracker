package models

type FilterData struct {
	Artists         []ArtistCard
	UniqueLocations []string
	SelectedFilters FilterParams
	TotalResults    int
	CurrentPath     string
}

type FilterParams struct {
	MemberCounts   []int
	Locations      []string
	CreationStart  int
	CreationEnd    int
	AlbumStartYear int
	AlbumEndYear   int
}
