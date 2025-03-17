// models/artist.go
package models

// ApiIndex holds the structure of the main API response
type ApiIndex struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

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

	// Processed data (not from direct JSON parsing)
	LocationsList        []string            `json:"-"`
	LocationStatesCities map[string][]string `json:"-"` // Maps states to cities
	DatesList            []string            `json:"-"`
	RelationsList        map[string][]string `json:"-"` // Maps locations to dates
}

// ArtistCard contains minimal info for list views
type ArtistCard struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
