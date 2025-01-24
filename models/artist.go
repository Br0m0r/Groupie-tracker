package models

// Artist represents a music artist or band with all their details
type Artist struct {
	ID            int                 `json:"id"`
	Image         string              `json:"image"`
	Name          string              `json:"name"`
	Members       []string            `json:"members"`
	CreationDate  int                 `json:"creationDate"`
	FirstAlbum    string              `json:"firstAlbum"`
	Locations     string              `json:"locations"`    // URL for locations data
	ConcertDates  string              `json:"concertDates"` // URL for dates data
	Relations     string              `json:"relations"`    // URL for relations data
	LocationsList []string            `json:"-"`            // Actual locations data after fetching
	DatesList     []string            `json:"-"`            // Actual dates after fetching
	RelationsList map[string][]string `json:"-"`            // Actual relations after fetching
}

// ArtistCard represents the minimal artist info for grid display
type ArtistCard struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// Location represents the concert locations data structure
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Index     []int    `json:"index,omitempty"`
}

// Date represents the concert dates data structure
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Relation represents the relation between dates and locations
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// SearchResult represents a single search result item
type SearchResult struct {
	Text        string `json:"text"`
	Type        string `json:"type"`
	ArtistName  string `json:"artistName"`
	Description string `json:"description"`
	ArtistId    int    `json:"artistId,omitempty"`
}

// SearchData represents the complete search response structure
type SearchData struct {
	Query   string
	Results []SearchResult
	ShowAll bool
}
