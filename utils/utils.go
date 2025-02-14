package utils

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"groupie/models"
)

// In this package we have a few functions that are used to format data for search and filter results.

// FormatLocation formats a location string to be more readable
// Example: "los_angeles-usa" -> "Los Angeles, USA"
func FormatLocation(location string) string {
	// Split by hyphen first (for country separation)
	parts := strings.Split(location, "-")

	// Format each part (city and country)
	for i, part := range parts {
		// Replace underscores with spaces
		part = strings.ReplaceAll(part, "_", " ")

		// Capitalize each word
		words := strings.Split(part, " ")
		for j, word := range words {
			words[j] = strings.Title(strings.ToLower(word))
		}
		parts[i] = strings.Join(words, " ")
	}

	// Join with comma if there's a country specified
	if len(parts) > 1 {
		return parts[0] + ", " + parts[1]
	}
	return parts[0]
}

// FormatDate formats a concert date string to be more readable
// Example: "02-01-2019" -> "January 2, 2019"
func FormatDate(date string) string {
	if strings.HasPrefix(date, "*") {
		date = strings.TrimSpace(date[1:])
	}
	t, err := time.Parse("02-01-2006", date)
	if err != nil {
		return date // Return original if parsing fails
	}

	// Format in a more readable way
	return t.Format("January 2, 2006")
}

// FormatRelation formats the relation data (date-location pairs)
// Takes a map of locations to dates and returns a formatted string
func FormatRelation(relations map[string][]string) map[string][]string {
	formatted := make(map[string][]string)

	for loc, dates := range relations {
		// Format the location
		formattedLoc := FormatLocation(loc)

		// Format each date
		formattedDates := make([]string, len(dates))
		for i, date := range dates {
			formattedDates[i] = FormatDate(date)
		}

		formatted[formattedLoc] = formattedDates
	}

	return formatted
}

// FormatLocationsList formats a slice of locations
func FormatLocationsList(locations []string) []string {
	formatted := make([]string, len(locations))
	for i, loc := range locations {
		formatted[i] = FormatLocation(loc)
	}
	return formatted
}

// extractYear gets the year from a date string in format "DD-MM-YYYY"
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

// Helper function to parse int with default value
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

// Helper function to get selected member counts from form
func GetMemberCounts(r *http.Request) []int {
	var counts []int
	for i := 1; i <= 8; i++ {
		if r.FormValue(fmt.Sprintf("members_%d", i)) != "" {
			counts = append(counts, i)
		}
	}
	return counts
}

// getUniqueLocations extracts unique locations from all artists
func GetUniqueLocations(artists []models.Artist) []string {
	locationMap := make(map[string]bool)

	for _, artist := range artists {
		for _, location := range artist.LocationsList {
			locationMap[location] = true
		}
	}

	// Convert map to sorted slice
	var locations []string
	for location := range locationMap {
		locations = append(locations, location)
	}
	sort.Strings(locations)
	return locations
}
