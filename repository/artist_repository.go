package repository

import (
	"fmt"
	"sort"
	"sync"

	"groupie/models"
	"groupie/utils"
)

// ArtistRepository manages all artist-related data and operations
type ArtistRepository struct {
	// Core data storage - use pointers to save memory
	artists []*models.Artist

	// Fast lookup maps - optimized for memory
	artistMap       map[int]*models.Artist // ID -> pointer (not full struct)
	locationMap     map[string][]int       // location -> artist IDs (not full structs)
	memberCountMap  map[int][]int          // count -> artist IDs
	creationYearMap map[int][]int          // year -> artist IDs
	albumYearMap    map[int][]int          // year -> artist IDs

	// Cached values
	minCreationYear int
	minAlbumYear    int
	uniqueLocations []string

	// Thread safety
	mu sync.RWMutex
}

// NewArtistRepository creates a new ArtistRepository instance
func NewArtistRepository() *ArtistRepository {
	return &ArtistRepository{
		artists:         make([]*models.Artist, 0),
		artistMap:       make(map[int]*models.Artist),
		locationMap:     make(map[string][]int),
		memberCountMap:  make(map[int][]int),
		creationYearMap: make(map[int][]int),
		albumYearMap:    make(map[int][]int),
		minCreationYear: utils.DefaultMinYear,
		minAlbumYear:    utils.DefaultMinYear,
		uniqueLocations: []string{},
	}
}

// LoadData fetches and processes all artist-related data from the API
func (ar *ArtistRepository) LoadData(apiIndex models.ApiIndex) error {
	apiRepo := NewAPIRepository("")

	// Fetch all data
	artists, err := apiRepo.FetchArtists(apiIndex.Artists)
	if err != nil {
		return fmt.Errorf("failed to fetch artists: %w", err)
	}

	locationsData, err := apiRepo.FetchLocations(apiIndex.Locations)
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}

	datesData, err := apiRepo.FetchDates(apiIndex.Dates)
	if err != nil {
		return fmt.Errorf("failed to fetch dates: %w", err)
	}

	relationData, err := apiRepo.FetchRelations(apiIndex.Relation)
	if err != nil {
		return fmt.Errorf("failed to fetch relation: %w", err)
	}

	// Process and enrich artists
	locationSet := make(map[string]bool)

	for i := range artists {
		artist := &artists[i]
		artist.LocationStatesCities = make(map[string][]string)

		// Add locations
		for _, locItem := range locationsData {
			if locItem.ID == artist.ID {
				for _, loc := range locItem.Locations {
					formattedLoc := utils.FormatLocation(loc)
					artist.LocationsList = append(artist.LocationsList, formattedLoc)
					locationSet[formattedLoc] = true

					// Check state/city mapping
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

		// Add dates
		for _, dateItem := range datesData {
			if dateItem.ID == artist.ID {
				for _, date := range dateItem.Dates {
					artist.DatesList = append(artist.DatesList, utils.FormatDate(date))
				}
				break
			}
		}

		// Add relations
		for _, relItem := range relationData {
			if relItem.ID == artist.ID {
				artist.RelationsList = utils.FormatRelation(relItem.DatesLocations)
				break
			}
		}
	}

	// Build optimized data structures
	ar.buildOptimizedMaps(artists, locationSet)

	return nil
}

// buildOptimizedMaps creates memory-efficient lookup maps
func (ar *ArtistRepository) buildOptimizedMaps(artists []models.Artist, locationSet map[string]bool) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	// Convert to pointer slice and build artist map
	ar.artists = make([]*models.Artist, len(artists))
	ar.artistMap = make(map[int]*models.Artist, len(artists))

	for i := range artists {
		ar.artists[i] = &artists[i]                 // Store pointer
		ar.artistMap[artists[i].ID] = ar.artists[i] // Map ID to pointer
	}

	// Build location map (store IDs, not full artists)
	ar.locationMap = make(map[string][]int)
	for _, artist := range ar.artists {
		for _, loc := range artist.LocationsList {
			ar.locationMap[loc] = append(ar.locationMap[loc], artist.ID)
		}
	}

	// Build member count map (store IDs)
	ar.memberCountMap = make(map[int][]int)
	for _, artist := range ar.artists {
		count := artist.GetMemberCount()
		ar.memberCountMap[count] = append(ar.memberCountMap[count], artist.ID)
	}

	// Build creation year map (store IDs)
	ar.creationYearMap = make(map[int][]int)
	for _, artist := range ar.artists {
		year := artist.CreationDate
		ar.creationYearMap[year] = append(ar.creationYearMap[year], artist.ID)
	}

	// Build album year map (store IDs)
	ar.albumYearMap = make(map[int][]int)
	for _, artist := range ar.artists {
		year := artist.GetFirstAlbumYear()
		if year > 0 {
			ar.albumYearMap[year] = append(ar.albumYearMap[year], artist.ID)
		}
	}

	// Cache minimum years
	ar.updateMinYears()

	// Build unique locations slice
	ar.uniqueLocations = make([]string, 0, len(locationSet))
	for loc := range locationSet {
		ar.uniqueLocations = append(ar.uniqueLocations, loc)
	}
	sort.Strings(ar.uniqueLocations)
}

// updateMinYears calculates and caches minimum years
func (ar *ArtistRepository) updateMinYears() {
	if len(ar.artists) == 0 {
		ar.minCreationYear = utils.DefaultMinYear
		ar.minAlbumYear = utils.DefaultMinYear
		return
	}

	ar.minCreationYear = ar.artists[0].CreationDate
	ar.minAlbumYear = ar.artists[0].GetFirstAlbumYear()

	for _, artist := range ar.artists {
		if artist.CreationDate < ar.minCreationYear {
			ar.minCreationYear = artist.CreationDate
		}
		if albumYear := artist.GetFirstAlbumYear(); albumYear > 0 && albumYear < ar.minAlbumYear {
			ar.minAlbumYear = albumYear
		}
	}
}

// GetArtistByID retrieves an artist by ID - O(1) lookup
func (ar *ArtistRepository) GetArtistByID(id int) (*models.Artist, error) {
	ar.mu.RLock()
	artist, ok := ar.artistMap[id]
	ar.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("artist with ID %d not found", id)
	}
	return artist, nil
}

// GetArtistsByLocation retrieves artists by location - optimized lookup
func (ar *ArtistRepository) GetArtistsByLocation(location string) []*models.Artist {
	ar.mu.RLock()
	artistIDs, ok := ar.locationMap[location]
	ar.mu.RUnlock()

	if !ok {
		return []*models.Artist{}
	}

	// Convert IDs to artists
	return ar.idsToArtists(artistIDs)
}

// GetArtistsByMemberCount retrieves artists by member count
func (ar *ArtistRepository) GetArtistsByMemberCount(count int) []*models.Artist {
	ar.mu.RLock()
	artistIDs := ar.memberCountMap[count]
	ar.mu.RUnlock()

	return ar.idsToArtists(artistIDs)
}

// GetArtistsByCreationYear retrieves artists by creation year
func (ar *ArtistRepository) GetArtistsByCreationYear(year int) []*models.Artist {
	ar.mu.RLock()
	artistIDs := ar.creationYearMap[year]
	ar.mu.RUnlock()

	return ar.idsToArtists(artistIDs)
}

// GetArtistsByAlbumYear retrieves artists by album year
func (ar *ArtistRepository) GetArtistsByAlbumYear(year int) []*models.Artist {
	ar.mu.RLock()
	artistIDs := ar.albumYearMap[year]
	ar.mu.RUnlock()

	return ar.idsToArtists(artistIDs)
}

// idsToArtists converts artist IDs to artist pointers
func (ar *ArtistRepository) idsToArtists(ids []int) []*models.Artist {
	if len(ids) == 0 {
		return []*models.Artist{}
	}

	ar.mu.RLock()
	defer ar.mu.RUnlock()

	artists := make([]*models.Artist, 0, len(ids))
	for _, id := range ids {
		if artist, ok := ar.artistMap[id]; ok {
			artists = append(artists, artist)
		}
	}
	return artists
}

// GetArtistCards returns lightweight cards for all artists
func (ar *ArtistRepository) GetArtistCards() []models.ArtistCard {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	cards := make([]models.ArtistCard, len(ar.artists))
	for i, artist := range ar.artists {
		cards[i] = artist.ToArtistCard()
	}
	return cards
}

// GetAllArtists returns all artists (returns pointers for efficiency)
func (ar *ArtistRepository) GetAllArtists() []*models.Artist {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	// Return copy of slice (pointers are shared, but slice is independent)
	artists := make([]*models.Artist, len(ar.artists))
	copy(artists, ar.artists)
	return artists
}

// GetUniqueLocations returns all unique concert locations
func (ar *ArtistRepository) GetUniqueLocations() []string {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	locations := make([]string, len(ar.uniqueLocations))
	copy(locations, ar.uniqueLocations)
	return locations
}

// GetMinYears returns cached minimum years
func (ar *ArtistRepository) GetMinYears() (minCreation, minAlbum int) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	return ar.minCreationYear, ar.minAlbumYear
}
