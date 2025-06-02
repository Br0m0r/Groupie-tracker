package utils

import (
	"net/http"

	"groupie/models"
)

func ExtractFilterParams(r *http.Request) models.FilterParams {
	return models.FilterParams{
		MemberCounts:   GetMemberCounts(r),
		Locations:      r.Form["location"],
		CreationStart:  ParseIntDefault(r.FormValue("creation_start"), 1950),
		CreationEnd:    ParseIntDefault(r.FormValue("creation_end"), 2025),
		AlbumStartYear: ParseIntDefault(r.FormValue("album_start"), 1950),
		AlbumEndYear:   ParseIntDefault(r.FormValue("album_end"), 2025),
	}
}
