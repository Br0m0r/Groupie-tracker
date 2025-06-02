package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"groupie/config"
	"groupie/models"
)

// APIRepositoryImpl handles fetching data from external API endpoints
type APIRepositoryImpl struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIRepository creates a new API repository instance
func NewAPIRepository(baseURL string) *APIRepositoryImpl {
	return &APIRepositoryImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: config.SERVER_READ_TIMEOUT, // Configurable timeout
		},
	}
}

// GetAPIIndex fetches the main API index with all endpoints
func (api *APIRepositoryImpl) GetAPIIndex() (*models.ApiIndex, error) {
	var index models.ApiIndex

	resp, err := api.httpClient.Get(api.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API index request failed with status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, fmt.Errorf("failed to decode API index: %w", err)
	}

	// Validate the response
	if err := index.Validate(); err != nil {
		return nil, fmt.Errorf("invalid API index: %w", err)
	}

	return &index, nil
}

// FetchArtists retrieves artist data from the API
func (api *APIRepositoryImpl) FetchArtists(url string) ([]models.Artist, error) {
	var artists []models.Artist

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artists request failed with status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("failed to decode artists: %w", err)
	}

	return artists, nil
}

// LocationData represents location data from API
type LocationData struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// FetchLocations retrieves location data from the API
func (api *APIRepositoryImpl) FetchLocations(url string) ([]LocationData, error) {
	var response struct {
		Index []LocationData `json:"index"`
	}

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("locations request failed with status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode locations: %w", err)
	}

	return response.Index, nil
}

// DateData represents date data from API
type DateData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// FetchDates retrieves concert dates from the API
func (api *APIRepositoryImpl) FetchDates(url string) ([]DateData, error) {
	var response struct {
		Index []DateData `json:"index"`
	}

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dates request failed with status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode dates: %w", err)
	}

	return response.Index, nil
}

// RelationData represents relation data from API
type RelationData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// FetchRelations retrieves date-location relations from the API
func (api *APIRepositoryImpl) FetchRelations(url string) ([]RelationData, error) {
	var response struct {
		Index []RelationData `json:"index"`
	}

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("relations request failed with status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode relations: %w", err)
	}

	return response.Index, nil
}
