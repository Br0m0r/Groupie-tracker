package store

import (
	"fmt"
	"sync"

	"groupie/client"
	"groupie/models"
)

type DataStore struct {
	Artists []models.Artist
	mu      sync.RWMutex
}

func New() *DataStore {
	return &DataStore{
		Artists: make([]models.Artist, 0),
	}
}

func (ds *DataStore) Initialize(baseURL string) error {
	// Create new client
	apiClient := client.New() // Now matches the client.New() signature

	// Fetch all data
	artists, err := apiClient.FetchAllData(baseURL)
	if err != nil {
		return err
	}

	// Store the data
	ds.mu.Lock()
	ds.Artists = artists
	ds.mu.Unlock()

	return nil
}

func (ds *DataStore) GetArtistCards() []models.ArtistCard {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	cards := make([]models.ArtistCard, len(ds.Artists))
	for i, artist := range ds.Artists {
		cards[i] = models.ArtistCard{
			ID:            artist.ID,
			Name:          artist.Name,
			Image:         artist.Image,
			CreationDate:  artist.CreationDate,
			FirstAlbum:    artist.FirstAlbum,
			Members:       artist.Members,
			LocationsList: artist.LocationsList,
		}
	}
	return cards
}

func (ds *DataStore) GetArtist(id int) (models.Artist, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	for _, artist := range ds.Artists {
		if artist.ID == id {
			return artist, nil
		}
	}
	return models.Artist{}, fmt.Errorf("artist with ID %d not found", id)
}

func (ds *DataStore) GetAllArtists() []models.Artist {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Create a copy of the artists slice to prevent concurrent access issues
	artists := make([]models.Artist, len(ds.Artists))
	copy(artists, ds.Artists)

	return artists
}
