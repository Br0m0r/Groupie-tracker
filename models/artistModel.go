// models/artist.go
package models

// ApiIndex holds the structure of the main API response
type ApiIndex struct {
	Artists string `json:"artists"`
	// Locations string `json:"locations"` }
	// Dates     string `json:"dates`      } those are not needed!!!
	// Relations string `json:"relation"`  }
}

// The basic artsit model that contains all  the data for an artist
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
	// These store the actual data after fetching from URLs
	LocationsList        []string            `json:"-"`
	LocationStatesCities map[string][]string `json:"-"` // maps states to their cities from geography
	DatesList            []string            `json:"-"`
	RelationsList        map[string][]string `json:"-"`
}

// The artist card model that contains only the basic data for an artist for index page
type ArtistCard struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Location struct {
	Locations []string `json:"locations"`
}

type Date struct {
	Dates []string `json:"dates"`
}

type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}
