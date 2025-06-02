package models

import (
	"strings"
)

// Artist represents a complete artist with all their information
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`

	// Processed data
	LocationsList        []string            `json:"-"`
	LocationStatesCities map[string][]string `json:"-"`
	DatesList            []string            `json:"-"`
	RelationsList        map[string][]string `json:"-"`

	// Runtime data
	CurrentYear int `json:"-"`
}

// GetMemberCount returns the number of members, capped at 8 for filtering
func (a *Artist) GetMemberCount() int {
	count := len(a.Members)
	if count > 8 {
		return 8
	}
	return count
}

// HasLocation checks if the artist performed at a specific location
func (a *Artist) HasLocation(location string) bool {
	for _, loc := range a.LocationsList {
		if strings.EqualFold(loc, location) {
			return true
		}
	}
	return false
}

// ArtistCard contains minimal info for list views (memory optimized)
