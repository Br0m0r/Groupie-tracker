package models

type SearchResult struct {
	Text        string `json:"text"`
	Type        string `json:"type"`
	ArtistName  string `json:"artistName"`
	Description string `json:"description"`
	ArtistId    int    `json:"artistId,omitempty"`
}

type SearchData struct {
	Query   string
	Results []SearchResult
}
