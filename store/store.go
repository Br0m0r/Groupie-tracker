package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

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

func New() *DataStore {
	return &DataStore{
		Artists: make([]models.Artist, 0),
	}
}

func (ds *DataStore) Initialize() error {
	client := &http.Client{Timeout: 10 * time.Second}
	fetchJSON := func(url string, target interface{}) error {
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("get %s: %w", url, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("get %s: unexpected status %s", url, resp.Status)
		}
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("decode %s: %w", url, err)
		}
		return nil
	}

	var index models.ApiIndex
	if err := fetchJSON("https://groupietrackers.herokuapp.com/api", &index); err != nil {
		return fmt.Errorf("failed to fetch API index: %w", err)
	}

	var artists []models.Artist
	if err := fetchJSON(index.Artists, &artists); err != nil {
		return fmt.Errorf("failed to fetch artists: %w", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(artists))

	for i := range artists {
		wg.Add(1)
		go func(artist *models.Artist) {
			defer wg.Done()
			artist.LocationStatesCities = make(map[string][]string)

					var location models.Location
			if err := fetchJSON(artist.Locations, &location); err != nil {
				errChan <- fmt.Errorf("failed to fetch locations for artist %d: %w", artist.ID, err)
				return
			}

			for _, loc := range location.Locations {
				formattedLoc := utils.FormatLocation(loc)
				artist.LocationsList = append(artist.LocationsList, formattedLoc)

				for state, cities := range utils.StateCityMap {
					for _, city := range cities {
						if formattedLoc == city {
							artist.LocationStatesCities[state] = append(artist.LocationStatesCities[state], city)
						}
					}
				}
			}
			var date models.Date
			if err := fetchJSON(artist.ConcertDates, &date); err != nil {
				errChan <- fmt.Errorf("failed to fetch dates for artist %d: %w", artist.ID, err)
				return
			}
			for _, date := range date.Dates {
				artist.DatesList = append(artist.DatesList, utils.FormatDate(date))
			}

			var relation models.Relation
			if err := fetchJSON(artist.Relations, &relation); err != nil {
				errChan <- fmt.Errorf("failed to fetch relations for artist %d: %w", artist.ID, err)
				return
			}
			artist.RelationsList = utils.FormatRelation(relation.DatesLocations)
		}(&artists[i])
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	ds.mu.Lock()
	ds.Artists = artists

	locationMap := make(map[string]bool)
	for _, artist := range artists {
		for _, location := range artist.LocationsList {
			locationMap[location] = true
		}
	}

	ds.UniqueLocations = make([]string, 0, len(locationMap))
	for location := range locationMap {
		ds.UniqueLocations = append(ds.UniqueLocations, location)
	}
	sort.Strings(ds.UniqueLocations)
	ds.mu.Unlock()
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
