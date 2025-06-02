package utils

import (
	"net/http"
	"text/template"

	"groupie/models"
)

// executeFilterTemplate helper function for rendering filter templates
func ExecuteFilterTemplate(w http.ResponseWriter, data models.FilterData) error {
	funcMap := template.FuncMap{
		"iterate": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
