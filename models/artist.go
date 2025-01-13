// models/artist.go
package models

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	// These are URLs in the initial response
	Locations    string `json:"locations"`
	ConcertDates string `json:"concertDates"`
	Relations    string `json:"relations"`

	// These will store the actual data after fetching from URLs
	LocationsList []string            `json:"-"`
	DatesList     []string            `json:"-"`
	RelationsList map[string][]string `json:"-"`
}
type ArtistCard struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Image         string   `json:"image"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	Members       []string `json:"members"`
	LocationsList []string `json:"locations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Index     []int    `json:"index,omitempty"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
