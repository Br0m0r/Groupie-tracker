package repository

import (
	"groupie/models"
)

// ArtistRepository defines the interface for accessing artist data
type ArtistRepository interface {
	// Get artist by ID
	GetArtistByID(id int) (models.Artist, error)
	
	// Get all artists
	GetAllArtists() []models.Artist
	
	// Get minimal artist information for listing
	GetArtistCards() []models.ArtistCard
	
	// Get unique concert locations from all artists
	GetUniqueLocations() []string
	
	// Get minimum creation year and album year across all artists
	GetMinYears() (minCreation, minAlbum int)
	
	// Load artist data from API
	LoadData(apiIndex models.ApiIndex) error
}

// CoordinatesRepository defines the interface for accessing geographic coordinates
type CoordinatesRepository interface {
	// Get coordinates for a location, fetching if needed
	Get(location string) (models.Coordinates, error)
	
	// Start background loading of coordinates for a set of locations
	PrefetchLocations(locations []string)
	
	// Import coordinates from another repository (useful for data refresh)
	ImportCache(other CoordinatesRepository)
}

// APIRepository defines the interface for accessing external API data
type APIRepository interface {
	// Get the API index with URLs for all endpoints
	GetAPIIndex() (models.ApiIndex, error)
	
	// Fetch artist data
	FetchArtists(url string) ([]models.Artist, error)
	
	// Fetch locations data
	FetchLocations(url string) ([]struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	}, error)
	
	// Fetch concert dates
	FetchDates(url string) ([]struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}, error)
	
	// Fetch relationships between dates and locations
	FetchRelations(url string) ([]struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}, error)
}