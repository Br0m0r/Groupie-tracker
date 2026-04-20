package handlers

import (
	"html/template"
	"net/http"
)

type ErrorType struct {
	Status  int
	Message string
}

type ErrorData struct {
	ErrorType
	Description string
}

var (
	ErrBadRequest = ErrorType{
		Status:  http.StatusBadRequest,
		Message: "Bad Request",
	}
	ErrNotFound = ErrorType{
		Status:  http.StatusNotFound,
		Message: "Page Not Found",
	}
	ErrInternalServer = ErrorType{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
	}
	ErrInvalidID = ErrorType{
		Status:  http.StatusBadRequest,
		Message: "Invalid ID Format",
	}
)

// ErrorHandler renders the error page template with provided error information
// If template processing fails, falls back to basic HTTP error response
func ErrorHandler(w http.ResponseWriter, errType ErrorType, description string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(errType.Status)

	data := ErrorData{
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
