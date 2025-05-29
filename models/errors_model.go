package models

import "net/http"

// ErrorType represents a basic HTTP error with status code and message
type ErrorType struct {
	Status  int
	Message string
}

// ErrorData combines an ErrorType with additional description for template rendering
type ErrorData struct {
	ErrorType
	Description string
}

// Predefined application errors
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
