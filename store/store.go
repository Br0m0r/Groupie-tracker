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

// GetArtistCards returns a slice of ArtistCard structures,
// which are simplified representations (ID, Name, Image) of the full Artist data.
// This function is typically called by the HomeHandler (in handlers.go) to render
// the homepage with a grid or list of artist cards.
func (ds *DataStore) GetArtistCards() []models.ArtistCard {
	// Acquire a read lock to safely access ds.Artists (a slice of models.Artist)
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Create a slice of ArtistCard (defined in artistModel.go) with the same length as ds.Artists.
	cards := make([]models.ArtistCard, len(ds.Artists))
	// Iterate over each Artist (models.Artist) in ds.Artists and extract key fields.
	for i, artist := range ds.Artists {
		cards[i] = models.ArtistCard{
			ID:    artist.ID,    // ID from models.Artist is used to set ArtistCard.ID
			Name:  artist.Name,  // Name from models.Artist is used to set ArtistCard.Name
			Image: artist.Image, // Image URL from models.Artist is used to set ArtistCard.Image
		}
	}
	// Return the slice of ArtistCard objects for rendering on the homepage.
	return cards
}

// GetArtist returns a complete Artist structure for a given artist ID.
// This function is typically called by the ArtistHandler (in handlers.go) when a user
// clicks on an artist card or navigates to an artist's detail page.
// It retrieves the full Artist record (as defined in artistModel.go) with detailed information.
func (ds *DataStore) GetArtist(id int) (models.Artist, error) {
	// Acquire a read lock to safely access ds.Artists.
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Iterate over the slice of models.Artist to find the artist with the matching ID.
	for _, artist := range ds.Artists {
		if artist.ID == id {
			// Return the complete models.Artist record when a match is found.
			return artist, nil
		}
	}
	// If no matching artist is found, return an empty models.Artist and an error.
	return models.Artist{}, fmt.Errorf("artist with ID %d not found", id)
}

// GetAllArtists returns a copy of the complete list of Artist structures.
// This function is typically called by the SearchHandler (in search.go) for processing
// search queries and by the FilterHandler (in filter.go) to filter the artist list.
// It provides the full collection of models.Artist (defined in artistModel.go) for these operations.
func (ds *DataStore) GetAllArtists() []models.Artist {
	// Acquire a read lock to safely access ds.Artists.
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Create a new slice of models.Artist with the same length as ds.Artists.
	artists := make([]models.Artist, len(ds.Artists))
	// Copy the content of ds.Artists (slice of models.Artist) into the new slice.
	copy(artists, ds.Artists)
	// Return the copied slice to ensure the original data remains unmodified during searches/filters.
	return artists
}
