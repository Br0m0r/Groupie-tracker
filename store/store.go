package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sort"
	"sync"

	"groupie/models"
	"groupie/utils"
)

// DataStore holds the artists, unique locations, and coordinate cache.
type DataStore struct {
	Artists         []models.Artist
	UniqueLocations []string
	mu              sync.RWMutex
	CoordinateCache struct {
		data map[string]models.Coordinates
		mu   sync.RWMutex
	}
	// Flag to prevent running multiple coordinate loaders concurrently.
	coordinateLoaderRunning bool
}

// New creates and returns a new DataStore.
func New() *DataStore {
	return &DataStore{
		Artists: make([]models.Artist, 0),
	}
}

// Initialize fetches API data and updates the datastore using caching.
func (ds *DataStore) Initialize() error {
	// 1. Get API index.
	var index models.ApiIndex
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		return fmt.Errorf("failed to fetch API index: %v", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return fmt.Errorf("failed to decode API index: %v", err)
	}

	// 2. Fetch artists data.
	var artists []models.Artist
	resp, err = http.Get(index.Artists)
	if err != nil {
		return fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return fmt.Errorf("failed to decode artists: %v", err)
	}

	// 3. Fetch locations data.
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

	// 4. Fetch dates data.
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

	// 5. Fetch relation data.
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

	// 6. Combine data into artists and build unique locations.
	locationMap := make(map[string]bool)
	for i := range artists {
		artist := &artists[i]
		artist.LocationStatesCities = make(map[string][]string)
		// Process locations.
		for _, locItem := range locationsData.Index {
			if locItem.ID == artist.ID {
				for _, loc := range locItem.Locations {
					formattedLoc := utils.FormatLocation(loc)
					artist.LocationsList = append(artist.LocationsList, formattedLoc)
					locationMap[formattedLoc] = true
					// Map states and cities using our StateCityMap.
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
		// Process dates.
		for _, dateItem := range datesData.Index {
			if dateItem.ID == artist.ID {
				for _, date := range dateItem.Dates {
					artist.DatesList = append(artist.DatesList, utils.FormatDate(date))
				}
				break
			}
		}
		// Process relations.
		for _, relItem := range relationData.Index {
			if relItem.ID == artist.ID {
				artist.RelationsList = utils.FormatRelation(relItem.DatesLocations)
				break
			}
		}
	}

	newUniqueLocations := make([]string, 0, len(locationMap))
	for loc := range locationMap {
		newUniqueLocations = append(newUniqueLocations, loc)
	}
	sort.Strings(newUniqueLocations)

	// Save previous unique locations.
	ds.mu.Lock()
	oldUniqueLocations := ds.UniqueLocations
	// Build a map of existing artists.
	oldArtistsMap := make(map[int]models.Artist)
	for _, a := range ds.Artists {
		oldArtistsMap[a.ID] = a
	}
	var updatedArtists []models.Artist
	newCount, updatedCount, unchangedCount := 0, 0, 0
	for _, newArtist := range artists {
		if old, exists := oldArtistsMap[newArtist.ID]; exists {
			if reflect.DeepEqual(old, newArtist) {
				updatedArtists = append(updatedArtists, old)
				unchangedCount++
			} else {
				updatedArtists = append(updatedArtists, newArtist)
				updatedCount++
			}
		} else {
			updatedArtists = append(updatedArtists, newArtist)
			newCount++
		}
	}
	ds.Artists = updatedArtists
	ds.UniqueLocations = newUniqueLocations
	ds.mu.Unlock()

	// Log a summary of changes.
	if newCount+updatedCount == 0 {
		log.Println("Initialize: No new data was fetched; all records are unchanged.")
	} else {
		log.Printf("Initialize: %d new records added, %d records updated, and %d records unchanged.",
			newCount, updatedCount, unchangedCount)
	}

	// Ensure the coordinate cache is preserved (create it only if nil).
	ds.CoordinateCache.mu.Lock()
	if ds.CoordinateCache.data == nil {
		ds.CoordinateCache.data = make(map[string]models.Coordinates)
	}
	ds.CoordinateCache.mu.Unlock()

	// Only launch the coordinate loader if unique locations have changed.
	if reflect.DeepEqual(oldUniqueLocations, newUniqueLocations) {
		log.Println("Unique locations unchanged; skipping coordinate loader.")
	} else {
		ds.loadCoordinatesInBackground()
	}

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

func (ds *DataStore) GetMinYears() (minCreation, minAlbum int) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if len(ds.Artists) > 0 {
		minCreation = ds.Artists[0].CreationDate
		minAlbum = utils.ExtractYear(ds.Artists[0].FirstAlbum)
	} else {
		return 1950, 1950
	}
	for _, artist := range ds.Artists {
		if artist.CreationDate < minCreation {
			minCreation = artist.CreationDate
		}
		albumYear := utils.ExtractYear(artist.FirstAlbum)
		if albumYear < minAlbum && albumYear > 0 {
			minAlbum = albumYear
		}
	}
	return minCreation, minAlbum
}
