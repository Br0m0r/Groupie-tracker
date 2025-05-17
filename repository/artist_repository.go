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
	memberCountMap  map[int][]models.Artist    // map for fast member-count lookups
	creationYearMap map[int][]models.Artist    // map for fast creation-year lookups
	albumYearMap    map[int][]models.Artist    // map for fast first-album-year lookups
	minCreationYear int                        // cached minimum creation year
	minAlbumYear    int                        // cached minimum first-album year
	uniqueLocations []string
	mu              sync.RWMutex
}

// NewArtistRepository creates a new ArtistRepository instance
func NewArtistRepository() *ArtistRepositoryImpl {
	return &ArtistRepositoryImpl{
		artists:         make([]models.Artist, 0),
		artistMap:       make(map[int]models.Artist),
		locationMap:     make(map[string][]models.Artist),
		memberCountMap:  make(map[int][]models.Artist),
		creationYearMap: make(map[int][]models.Artist),
		albumYearMap:    make(map[int][]models.Artist),
		minCreationYear: utils.DefaultMinYear,
		minAlbumYear:    utils.DefaultMinYear,
		uniqueLocations: []string{},
	}
}

// LoadData fetches and processes all artist-related data from the API
func (ar *ArtistRepositoryImpl) LoadData(apiIndex models.ApiIndex) error {
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

	// Store slice
	ar.artists = artists

	// Build ID lookup map
	ar.artistMap = make(map[int]models.Artist, len(artists))
	for _, a := range artists {
		ar.artistMap[a.ID] = a
	}

	// Build location lookup map
	ar.locationMap = make(map[string][]models.Artist)
	for _, a := range artists {
		for _, loc := range a.LocationsList {
			ar.locationMap[loc] = append(ar.locationMap[loc], a)
		}
	}

	// Build member-count map
	ar.memberCountMap = make(map[int][]models.Artist)
	for _, a := range artists {
		count := len(a.Members)
		ar.memberCountMap[count] = append(ar.memberCountMap[count], a)
	}

	// Build creation-year map
	ar.creationYearMap = make(map[int][]models.Artist)
	for _, a := range artists {
		year := a.CreationDate
		ar.creationYearMap[year] = append(ar.creationYearMap[year], a)
	}

	// Build first-album-year map
	ar.albumYearMap = make(map[int][]models.Artist)
	for _, a := range artists {
		year := utils.ExtractYear(a.FirstAlbum)
		if year > 0 {
			ar.albumYearMap[year] = append(ar.albumYearMap[year], a)
		}
	}

	// Cache minimum years
	if len(ar.artists) > 0 {
		ar.minCreationYear = ar.artists[0].CreationDate
		ar.minAlbumYear = utils.ExtractYear(ar.artists[0].FirstAlbum)
		for _, a := range ar.artists {
			if a.CreationDate < ar.minCreationYear {
				ar.minCreationYear = a.CreationDate
			}
			if y := utils.ExtractYear(a.FirstAlbum); y > 0 && y < ar.minAlbumYear {
				ar.minAlbumYear = y
			}
		}
	} else {
		ar.minCreationYear = utils.DefaultMinYear
		ar.minAlbumYear = utils.DefaultMinYear
	}

	// Build uniqueLocations slice
	ar.uniqueLocations = make([]string, 0, len(locationMap))
	for loc := range locationMap {
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

// GetArtistsByLocation retrieves artists who performed at the given location
func (ar *ArtistRepositoryImpl) GetArtistsByLocation(location string) []models.Artist {
	ar.mu.RLock()
	list, ok := ar.locationMap[location]
	ar.mu.RUnlock()
	if !ok {
		return []models.Artist{}
	}
	result := make([]models.Artist, len(list))
	copy(result, list)
	return result
}

// GetArtistsByMemberCount retrieves artists matching the given member count
func (ar *ArtistRepositoryImpl) GetArtistsByMemberCount(count int) []models.Artist {
	ar.mu.RLock()
	list := ar.memberCountMap[count]
	ar.mu.RUnlock()
	result := make([]models.Artist, len(list))
	copy(result, list)
	return result
}

// GetArtistsByCreationYear retrieves artists matching the given creation year
func (ar *ArtistRepositoryImpl) GetArtistsByCreationYear(year int) []models.Artist {
	ar.mu.RLock()
	list := ar.creationYearMap[year]
	ar.mu.RUnlock()
	result := make([]models.Artist, len(list))
	copy(result, list)
	return result
}

// GetArtistsByAlbumYear retrieves artists matching the given first-album year
func (ar *ArtistRepositoryImpl) GetArtistsByAlbumYear(year int) []models.Artist {
	ar.mu.RLock()
	list := ar.albumYearMap[year]
	ar.mu.RUnlock()
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

// GetMinYears returns the cached minimum creation year and first album year
func (ar *ArtistRepositoryImpl) GetMinYears() (minCreation, minAlbum int) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	return ar.minCreationYear, ar.minAlbumYear
}

// AddArtists merges the given slice of artists into the repository,
// returning only those that were not already present.
func (ar *ArtistRepositoryImpl) AddArtists(artists []models.Artist) []models.Artist {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	var added []models.Artist
	for _, a := range artists {
		if _, exists := ar.artistMap[a.ID]; exists {
			continue
		}
		// 1) append to slice & map
		ar.artists = append(ar.artists, a)
		ar.artistMap[a.ID] = a

		// 2) update locationMap & uniqueLocations
		for _, loc := range a.LocationsList {
			ar.locationMap[loc] = append(ar.locationMap[loc], a)
			// if first time seeing loc, add to uniqueLocations
			found := false
			for _, u := range ar.uniqueLocations {
				if u == loc {
					found = true
					break
				}
			}
			if !found {
				ar.uniqueLocations = append(ar.uniqueLocations, loc)
			}
		}

		// 3) update your other maps (memberCountMap, creationYearMap, albumYearMap)
		count := len(a.Members)
		ar.memberCountMap[count] = append(ar.memberCountMap[count], a)

		ar.creationYearMap[a.CreationDate] = append(ar.creationYearMap[a.CreationDate], a)

		if y := utils.ExtractYear(a.FirstAlbum); y > 0 {
			ar.albumYearMap[y] = append(ar.albumYearMap[y], a)
		}

		// 4) update minima
		if a.CreationDate < ar.minCreationYear {
			ar.minCreationYear = a.CreationDate
		}
		if y := utils.ExtractYear(a.FirstAlbum); y > 0 && y < ar.minAlbumYear {
			ar.minAlbumYear = y
		}

		added = append(added, a)
	}

	// (re-sort uniqueLocations if you care about order)
	sort.Strings(ar.uniqueLocations)
	return added
}
