package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"groupie/models"
	"groupie/utils"
)

// ApiIndex holds the structure of the main API response
type ApiIndex struct {
	Artists string `json:"artists"`
	// Locations string `json:"locations"` }
	// Dates     string `json:"dates`      } those are not needed!!!
	// Relations string `json:"relation"`  }
}

type DataStore struct {
	Artists []models.Artist
	mu      sync.RWMutex
}

func New() *DataStore { // constructor for new Struct Datastore { Artist []model.Artist }
	return &DataStore{
		Artists: make([]models.Artist, 0),
	}
}

func (ds *DataStore) Initialize() error {
	// First get the API index
	var index ApiIndex
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		return fmt.Errorf("failed to fetch API index: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return fmt.Errorf("failed to decode API index: %v", err)
	}

	// Fetch artists data
	var artists []models.Artist
	resp, err = http.Get(index.Artists)
	if err != nil {
		return fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return fmt.Errorf("failed to decode artists: %v", err)
	}

	// Create wait group for concurrent fetching
	var wg sync.WaitGroup
	errChan := make(chan error, len(artists))

	// Fetch additional data for each artist concurrently
	for i := range artists {
		wg.Add(1)
		go func(artist *models.Artist) {
			defer wg.Done()

			// Fetch locations
			var location models.Location
			resp, err := http.Get(artist.Locations)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch locations for artist %d: %v", artist.ID, err)
				return
			}
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
				errChan <- err
				return
			}
			for _, loc := range location.Locations {
				artist.LocationsList = append(artist.LocationsList, utils.FormatLocation(loc))
			}

			// Fetch dates
			var date models.Date
			resp, err = http.Get(artist.ConcertDates)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch dates for artist %d: %v", artist.ID, err)
				return
			}
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(&date); err != nil {
				errChan <- err
				return
			}
			for _, date := range date.Dates {
				artist.DatesList = append(artist.DatesList, utils.FormatDate(date))
			}

			// Fetch relations
			var relation models.Relation
			resp, err = http.Get(artist.Relations)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch relations for artist %d: %v", artist.ID, err)
				return
			}
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
				errChan <- err
				return
			}
			artist.RelationsList = utils.FormatRelation(relation.DatesLocations)
		}(&artists[i])
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Store the data
	ds.mu.Lock()
	ds.Artists = artists
	ds.mu.Unlock()

	return nil
}

// The rest of your DataStore methods remain the same
func (ds *DataStore) GetArtistCards() []models.ArtistCard {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	cards := make([]models.ArtistCard, len(ds.Artists))
	for i, artist := range ds.Artists {
		cards[i] = models.ArtistCard{
			ID:    artist.ID,
			Name:  artist.Name,
			Image: artist.Image,
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

	artists := make([]models.Artist, len(ds.Artists))
	copy(artists, ds.Artists)
	return artists
}
