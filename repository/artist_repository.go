package repository

import (
	"fmt"
	"sort"
	"sync"

	"groupie/models"
	"groupie/utils"
)

// ArtistRepositoryImpl manages all artist-related data and operations
type ArtistRepositoryImpl struct {
	artists         []models.Artist
	artistMap       map[int]models.Artist      // map for fast ID lookups
	locationMap     map[string][]models.Artist // map for fast location-based lookups
	uniqueLocations []string
	mu              sync.RWMutex
}

// NewArtistRepository creates a new ArtistRepository instance
func NewArtistRepository() *ArtistRepositoryImpl {
	return &ArtistRepositoryImpl{
		artists:         make([]models.Artist, 0),
		artistMap:       make(map[int]models.Artist),
		locationMap:     make(map[string][]models.Artist),
		uniqueLocations: []string{},
	}
}

// LoadData fetches and processes all artist-related data from the API
func (ar *ArtistRepositoryImpl) LoadData(apiIndex models.ApiIndex) error {
	// Create a temporary API repository to fetch the data
	apiRepo := NewAPIRepository("")

	// Fetch artists, locations, dates, and relations
	artists, err := apiRepo.FetchArtists(apiIndex.Artists)
	if err != nil {
		return fmt.Errorf("failed to fetch artists: %v", err)
	}
	locationsData, err := apiRepo.FetchLocations(apiIndex.Locations)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %v", err)
	}
	datesData, err := apiRepo.FetchDates(apiIndex.Dates)
	if err != nil {
		return fmt.Errorf("failed to fetch dates: %v", err)
	}
	relationData, err := apiRepo.FetchRelations(apiIndex.Relation)
	if err != nil {
		return fmt.Errorf("failed to fetch relation: %v", err)
	}

	// Temporary map to gather unique locations
	locationMap := make(map[string]bool)

	for i := range artists {
		artist := &artists[i]
		artist.LocationStatesCities = make(map[string][]string)

		// Find and add locations
		for _, locItem := range locationsData {
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
		for _, dateItem := range datesData {
			if dateItem.ID == artist.ID {
				for _, date := range dateItem.Dates {
					artist.DatesList = append(artist.DatesList, utils.FormatDate(date))
				}
				break
			}
		}

		// Find and add relations
		for _, relItem := range relationData {
			if relItem.ID == artist.ID {
				artist.RelationsList = utils.FormatRelation(relItem.DatesLocations)
				break
			}
		}
	}

	ar.mu.Lock()
	defer ar.mu.Unlock()

	// Store the complete slice
	ar.artists = artists

	// Build the ID lookup map
	ar.artistMap = make(map[int]models.Artist, len(artists))
	for _, a := range artists {
		ar.artistMap[a.ID] = a
	}

	// Build the location lookup map
	ar.locationMap = make(map[string][]models.Artist)
	for _, a := range artists {
		for _, loc := range a.LocationsList {
			ar.locationMap[loc] = append(ar.locationMap[loc], a)
		}
	}

	// Build uniqueLocations slice from map keys
	ar.uniqueLocations = make([]string, 0, len(ar.locationMap))
	for loc := range ar.locationMap {
		ar.uniqueLocations = append(ar.uniqueLocations, loc)
	}
	sort.Strings(ar.uniqueLocations)

	return nil
}

// GetArtistByID retrieves an artist by their ID using constant-time map lookup
func (ar *ArtistRepositoryImpl) GetArtistByID(id int) (models.Artist, error) {
	ar.mu.RLock()
	artist, ok := ar.artistMap[id]
	ar.mu.RUnlock()
	if !ok {
		return models.Artist{}, fmt.Errorf("artist with ID %d not found", id)
	}
	return artist, nil
}

// GetArtistsByLocation retrieves artists who have performed at the given location
func (ar *ArtistRepositoryImpl) GetArtistsByLocation(location string) []models.Artist {
	ar.mu.RLock()
	list, ok := ar.locationMap[location]
	ar.mu.RUnlock()
	if !ok {
		return []models.Artist{}
	}
	// Return a copy to avoid mutation
	result := make([]models.Artist, len(list))
	copy(result, list)
	return result
}

// GetArtistCards retrieves minimal information for all artists
func (ar *ArtistRepositoryImpl) GetArtistCards() []models.ArtistCard {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	cards := make([]models.ArtistCard, len(ar.artists))
	for i, artist := range ar.artists {
		cards[i] = models.ArtistCard{
			ID:    artist.ID,
			Name:  artist.Name,
			Image: artist.Image,
		}
	}
	return cards
}

// GetAllArtists returns a copy of all artist data
func (ar *ArtistRepositoryImpl) GetAllArtists() []models.Artist {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	copySlice := make([]models.Artist, len(ar.artists))
	copy(copySlice, ar.artists)
	return copySlice
}

// GetUniqueLocations returns all unique concert locations
func (ar *ArtistRepositoryImpl) GetUniqueLocations() []string {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	locs := make([]string, len(ar.uniqueLocations))
	copy(locs, ar.uniqueLocations)
	return locs
}

// GetMinYears returns the minimum creation year and first album year across all artists
func (ar *ArtistRepositoryImpl) GetMinYears() (minCreation, minAlbum int) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	if len(ar.artists) == 0 {
		return 1950, 1950
	}

	minCreation = ar.artists[0].CreationDate
	minAlbum = utils.ExtractYear(ar.artists[0].FirstAlbum)

	for _, artist := range ar.artists {
		if artist.CreationDate < minCreation {
			minCreation = artist.CreationDate
		}
		if year := utils.ExtractYear(artist.FirstAlbum); year > 0 && year < minAlbum {
			minAlbum = year
		}
	}
	return minCreation, minAlbum
}
