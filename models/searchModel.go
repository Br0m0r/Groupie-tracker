package models

// SearchResult is a struct that holds the search result.
type SearchResult struct {
	Text        string `json:"text"`
	Type        string `json:"type"`
	ArtistName  string `json:"artistName"`
	Description string `json:"description"`
	ArtistId    int    `json:"artistId,omitempty"`
}

// SearchData is a struct that holds the search query and the results.
type SearchData struct {
	Query   string
	Results []SearchResult
}
