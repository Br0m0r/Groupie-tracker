package models

// SearchResult represents a single search result
type SearchResult struct {
	Text        string // The text to display
	Type        string // "artist", "member", "location", etc.
	ArtistName  string // Name of the artist this result belongs to
	Description string // Additional description
	ArtistId    int    // ID of the artist to link to
}

// SearchData holds the search query and results
type SearchData struct {
	Query       string
	Results     []SearchResult
	CurrentYear int
}
