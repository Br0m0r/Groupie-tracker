package store

import (
	"fmt"
	"log"

	"groupie/models"
	"groupie/repository"
)

var DefaultAPIBaseURL = "https://groupietrackers.herokuapp.com/api"

type DataStore struct {
	apiRepo         *repository.APIRepositoryImpl
	artistRepo      *repository.ArtistRepositoryImpl
	coordinatesRepo *repository.CoordinatesRepositoryImpl
}

func New() *DataStore {
	return &DataStore{
		apiRepo:         repository.NewAPIRepository(DefaultAPIBaseURL),
		artistRepo:      repository.NewArtistRepository(),
		coordinatesRepo: repository.NewCoordinatesRepository(),
	}
}

// Initialize fetches and processes API data for all repositories
func (ds *DataStore) Initialize() error {
	// 1. Get API index to find all endpoints
	index, err := ds.apiRepo.GetAPIIndex()
	if err != nil {
		return fmt.Errorf("failed to fetch API index: %v", err)
	}

	// 2. Initialize the artist repository with the API data
	if err := ds.artistRepo.LoadData(index); err != nil {
		return fmt.Errorf("failed to load artist data: %v", err)
	}

	// 3. Start prefetching coordinates for all unique locations
	uniqueLocations := ds.artistRepo.GetUniqueLocations()
	ds.coordinatesRepo.PrefetchLocations(uniqueLocations)

	return nil
}

func (ds *DataStore) SwapData(newStore *DataStore) {
	log.Println("Data refresh started")

	// Snapshot old stats
	oldArtistCount := len(ds.artistRepo.GetAllArtists())
	oldCoordCount := ds.coordinatesRepo.CacheSize()

	// Swap in the new artists
	ds.artistRepo = newStore.artistRepo

	// Now discover any locations we didn’t have yet
	newLocations := ds.artistRepo.GetUniqueLocations()
	added := 0
	for _, loc := range newLocations {
		if !ds.coordinatesRepo.Has(loc) {
			// this will fetch & cache under the hood
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

// GetArtistCards delegates to ArtistRepository
func (ds *DataStore) GetArtistCards() []models.ArtistCard {
	return ds.artistRepo.GetArtistCards()
}

// GetArtist delegates to ArtistRepository
func (ds *DataStore) GetArtist(id int) (models.Artist, error) {
	return ds.artistRepo.GetArtistByID(id)
}

// GetAllArtists delegates to ArtistRepository
func (ds *DataStore) GetAllArtists() []models.Artist {
	return ds.artistRepo.GetAllArtists()
}

// GetMinYears delegates to ArtistRepository
func (ds *DataStore) GetMinYears() (minCreation, minAlbum int) {
	return ds.artistRepo.GetMinYears()
}

// GetLocationCoordinates delegates to CoordinatesRepository
func (ds *DataStore) GetLocationCoordinates(location string) (models.Coordinates, error) {
	return ds.coordinatesRepo.Get(location)
}

// UniqueLocations returns all unique concert locations
func (ds *DataStore) UniqueLocations() []string {
	return ds.artistRepo.GetUniqueLocations()
}

// GetArtistsByLocation delegates to ArtistRepository
func (ds *DataStore) GetArtistsByLocation(location string) []models.Artist {
	return ds.artistRepo.GetArtistsByLocation(location)
}

// GetArtistsByMemberCount delegates to ArtistRepository
func (ds *DataStore) GetArtistsByMemberCount(count int) []models.Artist {
	return ds.artistRepo.GetArtistsByMemberCount(count)
}

// GetArtistsByCreationYear delegates to ArtistRepository
func (ds *DataStore) GetArtistsByCreationYear(year int) []models.Artist {
	return ds.artistRepo.GetArtistsByCreationYear(year)
}

// GetArtistsByAlbumYear delegates to ArtistRepository
func (ds *DataStore) GetArtistsByAlbumYear(year int) []models.Artist {
	return ds.artistRepo.GetArtistsByAlbumYear(year)
}

// RefreshData fetches the latest artist list and merges in only the new entries.
// It then fetches coordinates for any new locations, returning counts of new
// artists and new coords pulled in.
func (ds *DataStore) RefreshData() (newArtists, newCoords int, err error) {
    // 1) Get fresh API index
    index, err := ds.apiRepo.GetAPIIndex()
    if err != nil {
        return 0, 0, fmt.Errorf("failed to fetch API index: %v", err)
    }

    // 2) Fetch the full artist list from the API
    fresh, err := ds.apiRepo.FetchArtists(index.Artists)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to fetch artists: %v", err)
    }

    // 3) Merge in only the truly new artists
    added := ds.artistRepo.AddArtists(fresh)
    newArtists = len(added)

    // 4) For each new artist, fetch any coords we don't already have
    for _, a := range added {
        for _, loc := range a.LocationsList {
            if !ds.coordinatesRepo.Has(loc) {
                if _, err := ds.coordinatesRepo.Get(loc); err == nil {
                    newCoords++
                }
            }
        }
    }

    return newArtists, newCoords, nil
}