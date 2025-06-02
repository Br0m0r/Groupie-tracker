package models

import (
	"fmt"
	"strconv"
)

// Coordinates represents geographic coordinates
type Coordinates struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Address string  `json:"address"`
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

	return coords, nil
}
