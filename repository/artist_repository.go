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
	uniqueLocations []string
	mu              sync.RWMutex
}

// NewArtistRepository creates a new ArtistRepository instance
func NewArtistRepository() *ArtistRepositoryImpl {
	return &ArtistRepositoryImpl{artists: make([]models.Artist, 0)}
}

// LoadData fetches and processes all artist-related data from the API
func (ar *ArtistRepositoryImpl) LoadData(apiIndex models.ApiIndex) error {
	// Create a temporary API repository to fetch the data
	apiRepo := NewAPIRepository("")

	// 1. Fetch artists data
	artists, err := apiRepo.FetchArtists(apiIndex.Artists)
	if err != nil {
		return fmt.Errorf("failed to fetch artists: %v", err)
	}

	// 2. Fetch locations data
	locationsData, err := apiRepo.FetchLocations(apiIndex.Locations)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %v", err)
	}

	// 3. Fetch dates data
	datesData, err := apiRepo.FetchDates(apiIndex.Dates)
	if err != nil {
		return fmt.Errorf("failed to fetch dates: %v", err)
	}

	// 4. Fetch relation data
	relationData, err := apiRepo.FetchRelations(apiIndex.Relation)
	if err != nil {
		return fmt.Errorf("failed to fetch relation: %v", err)
	}

	// 5. Combine data into artists
	locationMap := make(map[string]bool) // For tracking unique locations

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

	// Store the data
	ar.mu.Lock()
	ar.artists = artists

	// Convert map to sorted slice for unique locations
	ar.uniqueLocations = make([]string, 0, len(locationMap))
	for location := range locationMap {
		ar.uniqueLocations = append(ar.uniqueLocations, location)
	}
	sort.Strings(ar.uniqueLocations)
	ar.mu.Unlock()

	return nil
}

// GetArtistByID retrieves an artist by their ID
func (ar *ArtistRepositoryImpl) GetArtistByID(id int) (models.Artist, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	for _, artist := range ar.artists {
		if artist.ID == id {
			return artist, nil
		}
	}
	return models.Artist{}, fmt.Errorf("artist with ID %d not found", id)
}

// GetArtistCards returns minimal information for all artists
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

	artists := make([]models.Artist, len(ar.artists))
	copy(artists, ar.artists)
	return artists
}

// GetUniqueLocations returns all unique concert locations
func (ar *ArtistRepositoryImpl) GetUniqueLocations() []string {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	// Create a copy to avoid external code modifying our internal data
	locations := make([]string, len(ar.uniqueLocations))
	copy(locations, ar.uniqueLocations)
	return locations
}

// GetMinYears returns the minimum year for creation date and minimum year for first albums
func (ar *ArtistRepositoryImpl) GetMinYears() (minCreation, minAlbum int) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	// Set initial values
	if len(ar.artists) > 0 {
		minCreation = ar.artists[0].CreationDate
		minAlbum = utils.ExtractYear(ar.artists[0].FirstAlbum)
	} else {
		// Default values if no artists found
		return 1950, 1950
	}

	// Find min values
	for _, artist := range ar.artists {
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
