package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"groupie/models"
	"groupie/store"
	"groupie/utils"
)

var dataStore *store.DataStore

func Initialize(ds *store.DataStore) {
	dataStore = ds
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, ErrNotFound, "Page not exist")
		return
	}

	data := models.FilterData{
		Artists:         dataStore.GetArtistCards(),
		UniqueLocations: dataStore.UniqueLocations,
		SelectedFilters: utils.GetDefaultFilterParams(),
		TotalResults:    len(dataStore.GetArtistCards()),
		CurrentPath:     r.URL.Path,
	}

	if err := executeFilterTemplate(w, data); err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to process template")
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorHandler(w, ErrBadRequest, "Artist ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, ErrInvalidID, "Invalid artist ID format")
		return
	}

	artist, err := dataStore.GetArtist(id)
	if err != nil {
		ErrorHandler(w, ErrNotFound, "Artist not found")
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to load template")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, artist)
	if err != nil {
		ErrorHandler(w, ErrInternalServer, "Failed to execute template")
		return
	}
}
