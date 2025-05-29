package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Coordinates represents geographic coordinates
type Coordinates struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Address string  `json:"address"`
}

// Validate performs basic validation on Coordinates
func (c *Coordinates) Validate() error {
	if c.Lat < -90 || c.Lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", c.Lat)
	}

	if c.Lon < -180 || c.Lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %f", c.Lon)
	}

	if strings.TrimSpace(c.Address) == "" {
		return errors.New("address is required")
	}

	return nil
}

// IsValid quickly checks if coordinates are valid (without error details)
func (c *Coordinates) IsValid() bool {
	return c.Lat >= -90 && c.Lat <= 90 &&
		c.Lon >= -180 && c.Lon <= 180 &&
		strings.TrimSpace(c.Address) != ""
}

// NominatimResponse represents the response from Nominatim API
type NominatimResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// ToCoordinates converts NominatimResponse to Coordinates
func (nr *NominatimResponse) ToCoordinates(address string) (*Coordinates, error) {
	lat, err := strconv.ParseFloat(nr.Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude format: %s", nr.Lat)
	}

	lon, err := strconv.ParseFloat(nr.Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude format: %s", nr.Lon)
	}

	coords := &Coordinates{
		Lat:     lat,
		Lon:     lon,
		Address: address,
	}

	if err := coords.Validate(); err != nil {
		return nil, err
	}

	return coords, nil
}
