package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"groupie/models"
)

// APIRepositoryImpl handles fetching data from external API endpoints
type APIRepositoryImpl struct {
	baseURL string
}

// NewAPIRepository creates a new API repository instance
func NewAPIRepository(baseURL string) APIRepository {
	if baseURL == "" {
		baseURL = "https://groupietrackers.herokuapp.com/api"
	}
	return &APIRepositoryImpl{
		baseURL: baseURL,
	}
}

// GetAPIIndex fetches the main API index with all endpoints
func (api *APIRepositoryImpl) GetAPIIndex() (models.ApiIndex, error) {
	var index models.ApiIndex
	resp, err := http.Get(api.baseURL)
	if err != nil {
		return index, fmt.Errorf("failed to fetch API index: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return index, fmt.Errorf("failed to decode API index: %v", err)
	}

	return index, nil
}

// FetchArtists retrieves artist data from the API
func (api *APIRepositoryImpl) FetchArtists(url string) ([]models.Artist, error) {
	var artists []models.Artist
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("failed to decode artists: %v", err)
	}

	return artists, nil
}

// FetchLocations retrieves location data from the API
func (api *APIRepositoryImpl) FetchLocations(url string) ([]struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}, error,
) {
	var response struct {
		Index []struct {
			ID        int      `json:"id"`
			Locations []string `json:"locations"`
		} `json:"index"`
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode locations: %v", err)
	}

	return response.Index, nil
}

// FetchDates retrieves concert dates from the API
func (api *APIRepositoryImpl) FetchDates(url string) ([]struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}, error,
) {
	var response struct {
		Index []struct {
			ID    int      `json:"id"`
			Dates []string `json:"dates"`
		} `json:"index"`
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dates: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode dates: %v", err)
	}

	return response.Index, nil
}

// FetchRelations retrieves date-location relations from the API
func (api *APIRepositoryImpl) FetchRelations(url string) ([]struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}, error,
) {
	var response struct {
		Index []struct {
			ID             int                 `json:"id"`
			DatesLocations map[string][]string `json:"datesLocations"`
		} `json:"index"`
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode relations: %v", err)
	}

	return response.Index, nil
}
