package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"groupie/models"
)

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) fetch(url string, target interface{}) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *Client) FetchAllData(baseURL string) ([]models.Artist, error) {
	// First get the initial artist data
	var artists []models.Artist
	if err := c.fetch(baseURL+"/artists", &artists); err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %v", err)
	}

	// Create wait group for concurrent fetching
	var wg sync.WaitGroup
	errChan := make(chan error, len(artists))

	// Fetch additional data for each artist concurrently
	for i := range artists {
		wg.Add(1)
		go func(artist *models.Artist) {
			defer wg.Done()

			// Fetch locations
			var location models.Location
			if err := c.fetch(artist.Locations, &location); err != nil {
				errChan <- fmt.Errorf("failed to fetch locations for artist %d: %v", artist.ID, err)
				return
			}
			artist.LocationsList = location.Locations

			// Fetch dates
			var date models.Date
			if err := c.fetch(artist.ConcertDates, &date); err != nil {
				errChan <- fmt.Errorf("failed to fetch dates for artist %d: %v", artist.ID, err)
				return
			}
			artist.DatesList = date.Dates

			// Fetch relations
			var relation models.Relation
			if err := c.fetch(artist.Relations, &relation); err != nil {
				errChan <- fmt.Errorf("failed to fetch relations for artist %d: %v", artist.ID, err)
				return
			}
			artist.RelationsList = relation.DatesLocations
		}(&artists[i])
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return artists, nil
}
