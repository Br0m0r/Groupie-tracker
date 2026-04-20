package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"groupie/models"
)

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

func GetMemberCounts(r *http.Request) []int {
	var counts []int
	for i := 1; i <= 8; i++ {
		if r.FormValue(fmt.Sprintf("members_%d", i)) != "" {
			counts = append(counts, i)
		}
	}
	return counts
}

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
