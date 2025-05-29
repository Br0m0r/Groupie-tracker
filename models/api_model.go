package models

import (
	"errors"
	"strings"
)

// ApiIndex holds the structure of the main API response
type ApiIndex struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

// Validate checks if the API index has all required endpoints
func (ai *ApiIndex) Validate() error {
	if strings.TrimSpace(ai.Artists) == "" {
		return errors.New("artists endpoint URL is required")
	}
	if strings.TrimSpace(ai.Locations) == "" {
		return errors.New("locations endpoint URL is required")
	}
	if strings.TrimSpace(ai.Dates) == "" {
		return errors.New("dates endpoint URL is required")
	}
	if strings.TrimSpace(ai.Relation) == "" {
		return errors.New("relation endpoint URL is required")
	}
	return nil
}
