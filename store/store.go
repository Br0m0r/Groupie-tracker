package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"groupie/models"
	"groupie/utils"
)

type DataStore struct {
	Artists         []models.Artist
	UniqueLocations []string
	mu              sync.RWMutex
	CoordinateCache struct {
		data map[string]models.Coordinates
		mu   sync.RWMutex
	}
}

func New() *DataStore { // Constructor for DataStore: initializes Artists slice
	return &DataStore{
		Artists: make([]models.Artist, 0),
	}
}

func (ds *DataStore) Initialize() error {
	// 1. Fetch the API index (updates the variable 'index' of type models.ApiIndex)
	var index models.ApiIndex
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		return fmt.Errorf("failed to fetch API index: %v", err)
	}
	defer resp.Body.Close()

	// Decode JSON response into index (models.ApiIndex from artistModel.go)
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil { // updates variable "index"
		return fmt.Errorf("failed to decode API index: %v", err)
	}

	// 2. Fetch artists data (updates the variable 'artists' of type []models.Artist)
	var artists []models.Artist
	resp, err = http.Get(index.Artists)
	if err != nil {
		return fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer resp.Body.Close()

	// Decode JSON response into artists (updates slice of models.Artist defined in artistModel.go)
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil { // updates "artists" ([]models.Artist)
		return fmt.Errorf("failed to decode artists: %v", err)
	}

	// Create wait group for concurrent fetching
	var wg sync.WaitGroup
	errChan := make(chan error, len(artists))

	// 3. Fetch additional data for each artist concurrently
	for i := range artists {
		wg.Add(1)
		go func(artist *models.Artist) {
			defer wg.Done()
			artist.LocationStatesCities = make(map[string][]string)

			// Fetch locations for the artist (uses the URL stored in artist.Locations)
			var location models.Location
			resp, err := http.Get(artist.Locations)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch locations for artist %d: %v", artist.ID, err)
				return
			}
			defer resp.Body.Close()

			// Decode JSON response into location (models.Location); data used to update artist.LocationsList and artist.LocationStatesCities
			if err := json.NewDecoder(resp.Body).Decode(&location); err != nil { // updates variable "location" (models.Location)
				errChan <- err
				return
			}

			// Process each location string and update artist.LocationsList
			for _, loc := range location.Locations {
				formattedLoc := utils.FormatLocation(loc)
				artist.LocationsList = append(artist.LocationsList, formattedLoc)

				// Check and map locations to their respective states using utils.StateCityMap
				for state, cities := range utils.StateCityMap {
					for _, city := range cities {
						if formattedLoc == city {
							artist.LocationStatesCities[state] = append(artist.LocationStatesCities[state], city)
						}
					}
				}
			}

			// Fetch concert dates for the artist (uses the URL stored in artist.ConcertDates)
			var date models.Date
			resp, err = http.Get(artist.ConcertDates)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch dates for artist %d: %v", artist.ID, err)
				return
			}
			defer resp.Body.Close()

			// Decode JSON response into date (models.Date); data used to update artist.DatesList
			if err := json.NewDecoder(resp.Body).Decode(&date); err != nil { // updates variable "date" (models.Date)
				errChan <- err
				return
			}
			for _, d := range date.Dates {
				artist.DatesList = append(artist.DatesList, utils.FormatDate(d))
			}

			// Fetch relations for the artist (uses the URL stored in artist.Relations)
			var relation models.Relation
			resp, err = http.Get(artist.Relations)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch relations for artist %d: %v", artist.ID, err)
				return
			}
			defer resp.Body.Close()

			// Decode JSON response into relation (models.Relation); data used to update artist.RelationsList
			if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil { // updates variable "relation" (models.Relation)
				errChan <- err
				return
			}
			artist.RelationsList = utils.FormatRelation(relation.DatesLocations)
		}(&artists[i])
	}

	// Wait for all goroutines to complete and close error channel
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for any errors from concurrent requests
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// 4. Store the data into the DataStore
	ds.mu.Lock()
	ds.Artists = artists

	// Calculate and store unique locations across all artists
	locationMap := make(map[string]bool)
	for _, artist := range artists {
		for _, loc := range artist.LocationsList {
			locationMap[loc] = true
		}
	}

	// Convert map to sorted slice for UniqueLocations
	ds.UniqueLocations = make([]string, 0, len(locationMap))
	for loc := range locationMap {
		ds.UniqueLocations = append(ds.UniqueLocations, loc)
	}
	sort.Strings(ds.UniqueLocations)
	ds.mu.Unlock()

	// Load additional coordinates data in the background
	ds.loadCoordinatesInBackground()

	return nil
}

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
