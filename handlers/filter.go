package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"groupie/models"
)

type FilterParams struct {
	Creation struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"creation"`
	Album struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"album"`
	Members   []int    `json:"members"`
	Locations []string `json:"locations"`
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, ErrBadRequest, "Only POST method is allowed")
		return
	}

	if err := r.ParseForm(); err != nil {
		ErrorHandler(w, ErrBadRequest, "Failed to parse form data")
		return
	}

	// Parse form values
	creationMin, _ := strconv.Atoi(r.FormValue("creation_min"))
	creationMax, _ := strconv.Atoi(r.FormValue("creation_max"))
	albumMin, _ := strconv.Atoi(r.FormValue("album_min"))
	albumMax, _ := strconv.Atoi(r.FormValue("album_max"))

	members := make([]int, 0)
	for _, m := range r.Form["members"] {
		if val, err := strconv.Atoi(m); err == nil {
			members = append(members, val)
		}
	}

	locations := r.Form["locations"]

	// Filter artists
	allArtists := dataStore.GetAllArtists()
	filteredArtists := make([]models.Artist, 0)

	for _, artist := range allArtists {
		if !isInRange(artist.CreationDate, creationMin, creationMax) {
			continue
		}

		albumYear := parseAlbumYear(artist.FirstAlbum)
		if !isInRange(albumYear, albumMin, albumMax) {
			continue
		}

		if len(members) > 0 && !contains(members, len(artist.Members)) {
			continue
		}

		if len(locations) > 0 && !hasMatchingLocation(artist.LocationsList, locations) {
			continue
		}

		filteredArtists = append(filteredArtists, artist)
	}

	// Re-render the index page with filtered results
	data := struct {
		Artists   []models.Artist
		Locations []string
	}{
		Artists:   filteredArtists,
		Locations: getUniqueLocations(),
	}

	tmpl := template.New("index.html").Funcs(template.FuncMap{
		"intRange": func(min, max int) []int {
			a := make([]int, max-min)
			for i := range a {
				a[i] = min + i
			}
			return a
		},
	})

	tmpl, err := tmpl.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
	}
}

func isInRange(value, min, max int) bool {
	return value >= min && value <= max
}

func contains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func hasMatchingLocation(artistLocations, filterLocations []string) bool {
	for _, filterLoc := range filterLocations {
		for _, artistLoc := range artistLocations {
			if strings.Contains(strings.ToLower(artistLoc), strings.ToLower(filterLoc)) {
				return true
			}
		}
	}
	return false
}

func parseAlbumYear(albumDate string) int {
	if len(albumDate) < 4 {
		return 1960
	}
	yearStr := albumDate[len(albumDate)-4:]
	var year int
	if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
		return 1960
	}
	return year
}

func getUniqueLocations() []string {
	locations := make(map[string]bool)
	for _, artist := range dataStore.GetAllArtists() {
		for _, loc := range artist.LocationsList {
			locations[loc] = true
		}
	}
	uniqueLocations := make([]string, 0, len(locations))
	for loc := range locations {
		uniqueLocations = append(uniqueLocations, loc)
	}
	sort.Strings(uniqueLocations)
	return uniqueLocations
}
