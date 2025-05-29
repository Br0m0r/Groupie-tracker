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

// ConvertToCards converts artist pointers to ArtistCards (UPDATED for pointers)
func ConvertToCards(artists []*models.Artist) []models.ArtistCard {
	cards := make([]models.ArtistCard, len(artists))
	for i, artist := range artists {
		cards[i] = artist.ToArtistCard() // artist is now a pointer
	}
	return cards
}

// ConvertToCardsFromValues converts artist values to ArtistCards (for backward compatibility)
func ConvertToCardsFromValues(artists []models.Artist) []models.ArtistCard {
	cards := make([]models.ArtistCard, len(artists))
	for i, artist := range artists {
		cards[i] = artist.ToArtistCard() // artist is a value
	}
	return cards
}
