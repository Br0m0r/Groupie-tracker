package repository

import (
	"groupie/models"
)

// DataRepository is a base interface for all data access operations
// This helps establish a common pattern for all repositories
type DataRepository interface {
	// Each repository should have some form of initialization
	Initialize() error
}

// ArtistRepository defines the interface for accessing artist data
type ArtistRepository interface {
	// Read operations
	GetArtistByID(id int) (models.Artist, error)
	GetAllArtists() []models.Artist
	GetArtistCards() []models.ArtistCard
	GetUniqueLocations() []string
	GetMinYears() (minCreation, minAlbum int)

	// Data loading operations
	LoadData(apiIndex models.ApiIndex) error
}

// CoordinatesRepository defines the interface for accessing geographic coordinates
type CoordinatesRepository interface {
	// Read operations
	Get(location string) (models.Coordinates, error)

	// Background/optimization operations
	PrefetchLocations(locations []string)
	ImportCache(other CoordinatesRepository)
}

// APIRepository defines the interface for accessing external API data
// This is the source of truth for all application data
type APIRepository interface {
	// Meta operations
	GetAPIIndex() (models.ApiIndex, error)

	// Data fetching operations
	FetchArtists(url string) ([]models.Artist, error)

	// Location data operations
	FetchLocations(url string) ([]struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	}, error)

	// Date data operations
	FetchDates(url string) ([]struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}, error)

	// Relation data operations
	FetchRelations(url string) ([]struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}, error)
}
