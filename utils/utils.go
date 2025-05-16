package utils

import (
	"strings"
	"time"
)

const DefaultMinYear = 1950

// Format Locations Dates Relations and LocationsList !!!!

// FormatLocation formats a location string to be more readable
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
