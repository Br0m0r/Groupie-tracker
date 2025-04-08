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

func New() *DataStore { // constructor for new Struct Datastore { Artist []model.Artist }
	return &DataStore{
		Artists: make([]models.Artist, 0),
	}
}

func (ds *DataStore) Initialize() error {
	// 1. Get API index to find all endpoints
	var index models.ApiIndex
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		return fmt.Errorf("failed to fetch API index: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return fmt.Errorf("failed to decode API index: %v", err)
	}

	// 2. Fetch artists data
	var artists []models.Artist
	resp, err = http.Get(index.Artists)
	if err != nil {
		return fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return fmt.Errorf("failed to decode artists: %v", err)
	}

	// 3. Fetch locations data
	var locationsData struct {
		Index []struct {
			ID        int      `json:"id"`
			Locations []string `json:"locations"`
		} `json:"index"`
	}
	resp, err = http.Get(index.Locations)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&locationsData); err != nil {
		return fmt.Errorf("failed to decode locations: %v", err)
	}

	// 4. Fetch dates data
	var datesData struct {
		Index []struct {
			ID    int      `json:"id"`
			Dates []string `json:"dates"`
		} `json:"index"`
	}
	resp, err = http.Get(index.Dates)
	if err != nil {
		return fmt.Errorf("failed to fetch dates: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&datesData); err != nil {
		return fmt.Errorf("failed to decode dates: %v", err)
	}

	// 5. Fetch relation data
	var relationData struct {
		Index []struct {
			ID             int                 `json:"id"`
			DatesLocations map[string][]string `json:"datesLocations"`
		} `json:"index"`
	}
	resp, err = http.Get(index.Relation)
	if err != nil {
		return fmt.Errorf("failed to fetch relation: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&relationData); err != nil {
		return fmt.Errorf("failed to decode relation: %v", err)
	}

	// 6. Combine data into artists
	locationMap := make(map[string]bool) // For tracking unique locations

	for i := range artists {
		artist := &artists[i]
		artist.LocationStatesCities = make(map[string][]string)

		// Find and add locations
		for _, locItem := range locationsData.Index {
			if locItem.ID == artist.ID {
				for _, loc := range locItem.Locations {
					formattedLoc := utils.FormatLocation(loc)
					artist.LocationsList = append(artist.LocationsList, formattedLoc)
					locationMap[formattedLoc] = true

					// Check if this location is in our StateCityMap
					for state, cities := range utils.StateCityMap {
						for _, city := range cities {
							if formattedLoc == city {
								artist.LocationStatesCities[state] = append(artist.LocationStatesCities[state], city)
							}
						}
					}
				}
				break
			}
		}

		// Find and add dates
		for _, dateItem := range datesData.Index {
			if dateItem.ID == artist.ID {
				for _, date := range dateItem.Dates {
					artist.DatesList = append(artist.DatesList, utils.FormatDate(date))
				}
				break
			}
		}

		// Find and add relations
		for _, relItem := range relationData.Index {
			if relItem.ID == artist.ID {
				artist.RelationsList = utils.FormatRelation(relItem.DatesLocations)
				break
			}
		}
	}

	// Store the data
	ds.mu.Lock()
	ds.Artists = artists

	// Convert map to sorted slice for unique locations
	ds.UniqueLocations = make([]string, 0, len(locationMap))
	for location := range locationMap {
		ds.UniqueLocations = append(ds.UniqueLocations, location)
	}
	sort.Strings(ds.UniqueLocations)
	ds.mu.Unlock()

	// Start loading coordinates in the background
	ds.loadCoordinatesInBackground()

	return nil
}

// SwapData safely replaces the data store contents with data from a new store
func (ds *DataStore) SwapData(newStore *DataStore) {
	// Lock both data stores to ensure thread safety
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Replace the main data
	ds.Artists = newStore.Artists
	ds.UniqueLocations = newStore.UniqueLocations

	// Merge the coordinate cache - preserve existing coordinates
	// and add any new ones from the new store
	ds.CoordinateCache.mu.Lock()
	defer ds.CoordinateCache.mu.Unlock()

	// If we haven't initialized the map yet, do it now
	if ds.CoordinateCache.data == nil {
		ds.CoordinateCache.data = make(map[string]models.Coordinates)
	}

	// Copy coordinates from new store (if any)
	if newStore.CoordinateCache.data != nil {
		newStore.CoordinateCache.mu.RLock()
		for loc, coord := range newStore.CoordinateCache.data {
			ds.CoordinateCache.data[loc] = coord
		}
		newStore.CoordinateCache.mu.RUnlock()
	}
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

// GetMinYears returns the minimum year for creation date and minimum year for first albums
func (ds *DataStore) GetMinYears() (minCreation, minAlbum int) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Set initial values
	if len(ds.Artists) > 0 {
		minCreation = ds.Artists[0].CreationDate
		minAlbum = utils.ExtractYear(ds.Artists[0].FirstAlbum)
	} else {
		// Default values if no artists found
		return 1950, 1950
	}

	// Find min values only
	for _, artist := range ds.Artists {
		// Min creation date
		if artist.CreationDate < minCreation {
			minCreation = artist.CreationDate
		}

		// Min first album year
		albumYear := utils.ExtractYear(artist.FirstAlbum)
		if albumYear < minAlbum && albumYear > 0 {
			minAlbum = albumYear
		}
	}

	return minCreation, minAlbum
}
