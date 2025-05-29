// Package handlers provides HTTP request handlers and error handling for the application
package handlers

import (
	"html/template"
	"net/http"

	"groupie/models"
)

// ErrorHandler renders the error page template with provided error information
// If template processing fails, falls back to basic HTTP error response
func ErrorHandler(w http.ResponseWriter, errType models.ErrorType, description string) {
	w.WriteHeader(errType.Status)

	data := models.ErrorData{
		ErrorType:   errType,
		Description: description,
	}

	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
