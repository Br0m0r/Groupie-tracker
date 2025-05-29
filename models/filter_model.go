package models

import (
	"fmt"
)

// FilterParams holds all possible filter parameters
type FilterParams struct {
	MemberCounts   []int    // Selected member counts
	Locations      []string // Selected locations  
	CreationStart  int      // Creation date range start
	CreationEnd    int      // Creation date range end
	AlbumStartYear int      // First album year range start
	AlbumEndYear   int      // First album year range end
}

// Validate performs basic validation on FilterParams
func (fp *FilterParams) Validate() error {
	// Validate member counts
	for _, count := range fp.MemberCounts {
		if count < 1 || count > 8 {
			return fmt.Errorf("member count must be between 1 and 8, got %d", count)
		}
	}
	
	// Validate creation date range
	if fp.CreationStart > fp.CreationEnd {
		return fmt.Errorf("creation start year (%d) cannot be greater than end year (%d)", fp.CreationStart, fp.CreationEnd)
	}
	
	// Validate album year range
	if fp.AlbumStartYear > fp.AlbumEndYear {
		return fmt.Errorf("album start year (%d) cannot be greater than end year (%d)", fp.AlbumStartYear, fp.AlbumEndYear)
	}
	
	return nil
}

// IsEmpty checks if all filter parameters are at their default/empty state
func (fp *FilterParams) IsEmpty() bool {
	return len(fp.MemberCounts) == 0 &&
		len(fp.Locations) == 0 &&
		fp.CreationStart == fp.CreationEnd &&
		fp.AlbumStartYear == fp.AlbumEndYear
}

// FilterData represents all the data needed for the filter page
type FilterData struct {
	Artists         []ArtistCard
	UniqueLocations []string
	SelectedFilters FilterParams
	TotalResults    int
	CurrentPath     string
	CurrentYear     int
}