package store

import (
	"fmt"
	"log"

	"groupie/models"
	"groupie/repository"
)

// DataStore coordinates access to all application data repositories
type DataStore struct {
	apiRepo         repository.APIRepository
	artistRepo      repository.ArtistRepository
	coordinatesRepo repository.CoordinatesRepository
}

// New creates a new DataStore instance with initialized repositories
func New() *DataStore {
	return &DataStore{
		apiRepo:         repository.NewAPIRepository(""),
		artistRepo:      repository.NewArtistRepository(),
		coordinatesRepo: repository.NewCoordinatesRepository(),
	}
}

// Initialize fetches and processes API data for all repositories
func (ds *DataStore) Initialize() error {
	// 1. Get API index to find all endpoints
	index, err := ds.apiRepo.GetAPIIndex()
	if err != nil {
		return fmt.Errorf("failed to fetch API index: %v", err)
	}

	// 2. Initialize the artist repository with the API data
	if err := ds.artistRepo.LoadData(index); err != nil {
		return fmt.Errorf("failed to load artist data: %v", err)
	}

	// 3. Start prefetching coordinates for all unique locations
	uniqueLocations := ds.artistRepo.GetUniqueLocations()
	ds.coordinatesRepo.PrefetchLocations(uniqueLocations)

	return nil
}

// SwapData safely replaces the data store contents with data from a new store
func (ds *DataStore) SwapData(newStore *DataStore) {
	// Create new repositories
	newArtistRepo := repository.NewArtistRepository()
	newCoordinatesRepo := repository.NewCoordinatesRepository()

	// Get the API index
	index, err := ds.apiRepo.GetAPIIndex()
	if err != nil {
		log.Printf("Error fetching API index during swap: %v", err)
		return
	}

	// Initialize the new artist repository
	if err := newArtistRepo.LoadData(index); err != nil {
		log.Printf("Error loading artist data during swap: %v", err)
		return
	}

	// Import coordinate cache from new store
	newCoordinatesRepo.ImportCache(newStore.coordinatesRepo)

	// Atomically swap the repositories
	ds.artistRepo = newArtistRepo
	ds.coordinatesRepo = newCoordinatesRepo
}

// GetArtistCards delegates to ArtistRepository
func (ds *DataStore) GetArtistCards() []models.ArtistCard {
	return ds.artistRepo.GetArtistCards()
}

// GetArtist delegates to ArtistRepository
func (ds *DataStore) GetArtist(id int) (models.Artist, error) {
	return ds.artistRepo.GetArtistByID(id)
}

// GetAllArtists delegates to ArtistRepository
func (ds *DataStore) GetAllArtists() []models.Artist {
	return ds.artistRepo.GetAllArtists()
}

// GetMinYears delegates to ArtistRepository
func (ds *DataStore) GetMinYears() (minCreation, minAlbum int) {
	return ds.artistRepo.GetMinYears()
}

// GetLocationCoordinates delegates to CoordinatesRepository
func (ds *DataStore) GetLocationCoordinates(location string) (models.Coordinates, error) {
	return ds.coordinatesRepo.Get(location)
}

// UniqueLocations returns all unique concert locations
func (ds *DataStore) UniqueLocations() []string {
	return ds.artistRepo.GetUniqueLocations()
}
