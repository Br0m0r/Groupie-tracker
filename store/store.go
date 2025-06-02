package store

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"groupie/config"
	"groupie/models"
	"groupie/repository"
)

// DataStore manages all data repositories and provides unified access
type DataStore struct {
	apiRepo         *repository.APIRepositoryImpl
	artistRepo      *repository.ArtistRepository
	coordinatesRepo *repository.CoordinatesRepository
}

// New creates a new DataStore instance
func New() *DataStore {
	return &DataStore{
		apiRepo:         repository.NewAPIRepository(config.API_BASE_URL),
		artistRepo:      repository.NewArtistRepository(),
		coordinatesRepo: repository.NewCoordinatesRepository(),
	}
}

// Initialize fetches and processes API data for all repositories
func (ds *DataStore) Initialize() error {
	// Get API index to find all endpoints
	index, err := ds.apiRepo.GetAPIIndex()
	if err != nil {
		return fmt.Errorf("failed to fetch API index: %w", err)
	}

	// Initialize the artist repository with the API data
	if err := ds.artistRepo.LoadData(*index); err != nil {
		return fmt.Errorf("failed to load artist data: %w", err)
	}

	// Start prefetching coordinates for all unique locations
	uniqueLocations := ds.artistRepo.GetUniqueLocations()
	ds.coordinatesRepo.PrefetchLocations(uniqueLocations)

	return nil
}

// SwapData refreshes data by swapping in new artist repository while preserving coordinates cache
func (ds *DataStore) SwapData(newStore *DataStore) {
	log.Println("Data refresh started")

	// Snapshot old stats
	oldArtistCount := len(ds.artistRepo.GetAllArtists())
	oldCoordCount := ds.coordinatesRepo.CacheSize()

	// Swap in the new artists repository
	ds.artistRepo = newStore.artistRepo

	// Discover any new locations we didn't have yet
	newLocations := ds.artistRepo.GetUniqueLocations()
	added := 0
	for _, loc := range newLocations {
		if !ds.coordinatesRepo.Has(loc) {
			// Fetch and cache new location
			if _, err := ds.coordinatesRepo.Get(loc); err == nil {
				added++
			}
		}
	}

	// Final stats
	newArtistCount := len(ds.artistRepo.GetAllArtists())
	newCoordCount := ds.coordinatesRepo.CacheSize()

	log.Printf(
		"Data refresh complete: artists %d→%d; coords %d→%d (added %d new)\n",
		oldArtistCount, newArtistCount,
		oldCoordCount, newCoordCount,
		added,
	)
}

// GetArtist retrieves a single artist by ID (returns pointer now)
func (ds *DataStore) GetArtist(id int) (*models.Artist, error) {
	return ds.artistRepo.GetArtistByID(id)
}

// GetAllArtists returns all artists (returns pointers now)
func (ds *DataStore) GetAllArtists() []*models.Artist {
	return ds.artistRepo.GetAllArtists()
}

// GetMinYears returns cached minimum creation and album years
func (ds *DataStore) GetMinYears() (minCreation, minAlbum int) {
	return ds.artistRepo.GetMinYears()
}

// GetLocationCoordinates retrieves coordinates for a location
func (ds *DataStore) GetLocationCoordinates(location string) (*models.Coordinates, error) {
	return ds.coordinatesRepo.Get(location)
}

// UniqueLocations returns all unique concert locations
func (ds *DataStore) UniqueLocations() []string {
	return ds.artistRepo.GetUniqueLocations()
}

// GetArtistsByLocation retrieves artists who performed at the given location (returns pointers)
func (ds *DataStore) GetArtistsByLocation(location string) []*models.Artist {
	return ds.artistRepo.GetArtistsByLocation(location)
}

// GetArtistsByMemberCount retrieves artists matching the given member count (returns pointers)
func (ds *DataStore) GetArtistsByMemberCount(count int) []*models.Artist {
	return ds.artistRepo.GetArtistsByMemberCount(count)
}

// GetArtistsByCreationYear retrieves artists matching the given creation year (returns pointers)
func (ds *DataStore) GetArtistsByCreationYear(year int) []*models.Artist {
	return ds.artistRepo.GetArtistsByCreationYear(year)
}

// GetArtistsByAlbumYear retrieves artists matching the given first-album year (returns pointers)
func (ds *DataStore) GetArtistsByAlbumYear(year int) []*models.Artist {
	return ds.artistRepo.GetArtistsByAlbumYear(year)
}

// RefreshData fetches the latest artist list and reloads all data
func (ds *DataStore) RefreshData() (newArtists, newCoords int, err error) {
	// Get fresh API index
	index, err := ds.apiRepo.GetAPIIndex()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch API index: %w", err)
	}

	// Create new repository and load fresh data
	newRepo := repository.NewArtistRepository()
	if err := newRepo.LoadData(*index); err != nil {
		return 0, 0, fmt.Errorf("failed to load fresh artist data: %w", err)
	}

	// Count changes (simplified - just compare total counts)
	oldCount := len(ds.artistRepo.GetAllArtists())
	newCount := len(newRepo.GetAllArtists())
	newArtists = newCount - oldCount
	if newArtists < 0 {
		newArtists = 0 // Handle case where count decreased
	}

	// Swap the repository
	ds.artistRepo = newRepo

	// Check for new locations to cache
	newLocations := ds.artistRepo.GetUniqueLocations()
	for _, loc := range newLocations {
		if !ds.coordinatesRepo.Has(loc) {
			if _, err := ds.coordinatesRepo.Get(loc); err == nil {
				newCoords++
			}
		}
	}

	log.Printf("Refresh complete: %d new artists, %d new coordinates", newArtists, newCoords)
	return newArtists, newCoords, nil
}

// store/store.go - add these methods
func (ds *DataStore) SearchArtists(query string) []models.SearchResult {
	// Early return for empty queries
	if strings.TrimSpace(query) == "" {
		return []models.SearchResult{}
	}

	artists := ds.GetAllArtists()
	query = strings.ToLower(strings.TrimSpace(query))
	isSingleLetter := len([]rune(query)) == 1

	// Pre-allocate results slice
	results := make([]models.SearchResult, 0, len(artists)*2)

	// Process each artist
	for _, artist := range artists {
		results = append(results, ds.searchInArtist(artist, query, isSingleLetter)...)
	}

	// Sort and return
	ds.sortResultsByType(results)
	return results
}

// Helper methods (private to store)
func (ds *DataStore) searchInArtist(artist *models.Artist, query string, isSingleLetter bool) []models.SearchResult {
	var results []models.SearchResult

	// Artist name search
	if ds.searchMatch(artist.Name, query, isSingleLetter) {
		results = append(results, models.SearchResult{
			Text:        artist.Name,
			Type:        "artist/band",
			ArtistName:  artist.Name,
			Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
			ArtistId:    artist.ID,
		})
	}

	// Members search
	for _, member := range artist.Members {
		if ds.searchMatch(member, query, isSingleLetter) {
			results = append(results, models.SearchResult{
				Text:        member,
				Type:        "member",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("Member of %s", artist.Name),
				ArtistId:    artist.ID,
			})
		}
	}

	// Locations search
	for _, location := range artist.LocationsList {
		if ds.searchMatch(location, query, isSingleLetter) {
			results = append(results, models.SearchResult{
				Text:        location,
				Type:        "location",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("Concert location for %s", artist.Name),
				ArtistId:    artist.ID,
			})
		}
	}

	// Creation date search - only for non-single letter queries
	if !isSingleLetter {
		creationStr := fmt.Sprintf("%d", artist.CreationDate)
		if strings.Contains(creationStr, query) {
			results = append(results, models.SearchResult{
				Text:        fmt.Sprintf("%s (%d)", artist.Name, artist.CreationDate),
				Type:        "creation date",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
				ArtistId:    artist.ID,
			})
		}
	}

	// First album search
	if ds.searchMatch(artist.FirstAlbum, query, isSingleLetter) {
		results = append(results, models.SearchResult{
			Text:        fmt.Sprintf("%s - %s", artist.Name, artist.FirstAlbum),
			Type:        "first album",
			ArtistName:  artist.Name,
			Description: fmt.Sprintf("First album by %s", artist.Name),
			ArtistId:    artist.ID,
		})
	}

	return results
}

func (ds *DataStore) searchMatch(value, query string, isSingleLetter bool) bool {
	if isSingleLetter {
		return strings.HasPrefix(strings.ToLower(value), query)
	}
	return strings.Contains(strings.ToLower(value), query)
}

func (ds *DataStore) sortResultsByType(results []models.SearchResult) {
	// Define type priorities
	typePriority := map[string]int{
		"artist/band":   1,
		"member":        2,
		"location":      3,
		"creation date": 4,
		"first album":   5,
	}

	// Sort by type priority
	sort.Slice(results, func(i, j int) bool {
		return typePriority[results[i].Type] < typePriority[results[j].Type]
	})
}

// Update the function signature
func GetDefaultFilterParams(dataStore *DataStore) models.FilterParams {
	minCreation, minAlbum := dataStore.GetMinYears()

	return models.FilterParams{
		CreationStart:  minCreation, // Real minimum from data!
		CreationEnd:    time.Now().Year(),
		AlbumStartYear: minAlbum, // Real minimum from data!
		AlbumEndYear:   time.Now().Year(),
		MemberCounts:   []int{},
		Locations:      []string{},
	}
}
