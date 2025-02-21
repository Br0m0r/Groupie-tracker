package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"groupie/models"
)

// ParseIntDefault safely parses a string to int with a default value
func ParseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return val
}

// ExtractYear gets the year from a date string in format "DD-MM-YYYY"
func ExtractYear(date string) int {
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return 0
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}
	return year
}

// GetMemberCounts extracts selected member counts from form
func GetMemberCounts(r *http.Request) []int {
	var counts []int
	for i := 1; i <= 8; i++ {
		if r.FormValue(fmt.Sprintf("members_%d", i)) != "" {
			counts = append(counts, i)
		}
	}
	return counts
}

// Helper function to convert Artists to ArtistCards
func ConvertToCards(artists []models.Artist) []models.ArtistCard {
	cards := make([]models.ArtistCard, len(artists))
	for i, artist := range artists {
		cards[i] = models.ArtistCard{
			ID:    artist.ID,
			Name:  artist.Name,
			Image: artist.Image,
		}
	}
	return cards
}

// Return the default filter parameters.
func GetDefaultFilterParams() models.FilterParams {
	return models.FilterParams{
		CreationStart:  1950,
		CreationEnd:    2024,
		AlbumStartYear: 1950,
		AlbumEndYear:   2024,
		MemberCounts:   []int{},    // Empty slice - no members selected
		Locations:      []string{}, // Empty slice - no locations selected
	}
}
