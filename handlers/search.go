package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"groupie/models"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		results := searchAllData(query)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	if query != "" {
		results := searchAllData(query)

		if len(results) == 1 {
			http.Redirect(w, r, fmt.Sprintf("/artist?id=%d", results[0].ArtistId), http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to load template")
			return
		}

		data := models.SearchData{
			Query:   query,
			Results: results,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		}
		return
	}

	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, models.SearchData{}); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
	}
}

func searchAllData(query string) []models.SearchResult {
	var results []models.SearchResult
	var (
		artistResults   []models.SearchResult
		memberResults   []models.SearchResult
		locationResults []models.SearchResult
		dateResults     []models.SearchResult
		albumResults    []models.SearchResult
	)

	artists := dataStore.GetAllArtists()
	query = strings.ToLower(query)
	isSingleLetter := len([]rune(query)) == 1

	for _, artist := range artists {
		artistNameLower := strings.ToLower(artist.Name)
		if (isSingleLetter && strings.HasPrefix(artistNameLower, query)) ||
			(!isSingleLetter && strings.Contains(artistNameLower, query)) {
			artistResults = append(artistResults, models.SearchResult{
				Text:        artist.Name,
				Type:        "artist/band",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
				ArtistId:    artist.ID,
			})
		}

		for _, member := range artist.Members {
			memberLower := strings.ToLower(member)
			if (isSingleLetter && strings.HasPrefix(memberLower, query)) ||
				(!isSingleLetter && strings.Contains(memberLower, query)) {
				memberResults = append(memberResults, models.SearchResult{
					Text:        member,
					Type:        "member",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Member of %s", artist.Name),
					ArtistId:    artist.ID,
				})
			}
		}

		for _, location := range artist.LocationsList {
			locationLower := strings.ToLower(location)
			if (isSingleLetter && strings.HasPrefix(locationLower, query)) ||
				(!isSingleLetter && strings.Contains(locationLower, query)) {
				locationResults = append(locationResults, models.SearchResult{
					Text:        location,
					Type:        "location",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Concert location for %s", artist.Name),
					ArtistId:    artist.ID,
				})
			}
		}

		if !isSingleLetter {
			creationStr := fmt.Sprintf("%d", artist.CreationDate)
			if strings.Contains(creationStr, query) {
				dateResults = append(dateResults, models.SearchResult{
					Text:        fmt.Sprintf("%s (%d)", artist.Name, artist.CreationDate),
					Type:        "creation date",
					ArtistName:  artist.Name,
					Description: fmt.Sprintf("Band formed in %d", artist.CreationDate),
					ArtistId:    artist.ID,
				})
			}
		}

		albumLower := strings.ToLower(artist.FirstAlbum)
		if (isSingleLetter && strings.HasPrefix(albumLower, query)) ||
			(!isSingleLetter && strings.Contains(albumLower, query)) {
			albumResults = append(albumResults, models.SearchResult{
				Text:        fmt.Sprintf("%s - %s", artist.Name, artist.FirstAlbum),
				Type:        "first album",
				ArtistName:  artist.Name,
				Description: fmt.Sprintf("First album by %s", artist.Name),
				ArtistId:    artist.ID,
			})
		}
	}

	results = append(results, artistResults...)
	results = append(results, memberResults...)
	results = append(results, locationResults...)
	results = append(results, dateResults...)
	results = append(results, albumResults...)

	return results
}
