package models

import (
	"errors"
	"fmt"
	"strconv"
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

	// URLs from the initial API response
	Locations    string `json:"locations"`
	ConcertDates string `json:"concertDates"`
	Relations    string `json:"relations"`

	// Processed data
	LocationsList        []string            `json:"-"`
	LocationStatesCities map[string][]string `json:"-"`
	DatesList            []string            `json:"-"`
	RelationsList        map[string][]string `json:"-"`

	// Runtime data
	CurrentYear int `json:"-"`
}

// Validate performs basic validation on an Artist
func (a *Artist) Validate() error {
	if a.ID <= 0 {
		return fmt.Errorf("artist ID must be positive, got %d", a.ID)
	}

	if strings.TrimSpace(a.Name) == "" {
		return errors.New("artist name is required")
	}

	if a.CreationDate <= 0 {
		return fmt.Errorf("creation date must be positive, got %d", a.CreationDate)
	}

	if len(a.Members) == 0 {
		return errors.New("artist must have at least one member")
	}

	return nil
}

// GetMemberCount returns the number of members, capped at 8 for filtering
func (a *Artist) GetMemberCount() int {
	count := len(a.Members)
	if count > 8 {
		return 8
	}
	return count
}

// GetFirstAlbumYear extracts the year from the FirstAlbum field
func (a *Artist) GetFirstAlbumYear() int {
	parts := strings.Split(a.FirstAlbum, "-")
	if len(parts) != 3 {
		return 0
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}
	return year
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

