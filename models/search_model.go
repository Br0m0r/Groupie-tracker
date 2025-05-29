package models

import (
	"errors"
	"fmt"
	"strings"
)

// SearchResult represents a single search result
type SearchResult struct {
	Text        string // The text to display
	Type        string // "artist", "member", "location", etc.
	ArtistName  string // Name of the artist this result belongs to
	Description string // Additional description
	ArtistId    int    // ID of the artist to link to
}

// Validate performs basic validation on SearchResult
func (sr *SearchResult) Validate() error {
	if strings.TrimSpace(sr.Text) == "" {
		return errors.New("search result text is required")
	}

	if strings.TrimSpace(sr.Type) == "" {
		return errors.New("search result type is required")
	}

	if strings.TrimSpace(sr.ArtistName) == "" {
		return errors.New("artist name is required")
	}

	if sr.ArtistId <= 0 {
		return fmt.Errorf("artist ID must be positive, got %d", sr.ArtistId)
	}

	return nil
}

// SearchData holds the search query and results
type SearchData struct {
	Query       string
	Results     []SearchResult
	CurrentYear int
}
