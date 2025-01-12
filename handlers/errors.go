// handlers/errors.go
package handlers

import (
	"html/template"
	"net/http"
)

type ErrorType struct {
	Status  int
	Message string
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

type ErrorData struct {
	Status      int
	Message     string
	Description string
}

func ErrorHandler(w http.ResponseWriter, errType ErrorType, description string) {
	w.WriteHeader(errType.Status)

	data := ErrorData{
		Status:      errType.Status,
		Message:     errType.Message,
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
