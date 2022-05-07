package web

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Type       string `json:"type"`
	Message    string `json:"message"`
}

var (
	ErrInvalidJSON      = AppError{StatusCode: http.StatusBadRequest, Type: "invalid_json", Message: "Invalid or malformed JSON"}
	ErrInternalServer   = AppError{StatusCode: http.StatusInternalServerError, Type: "internal_server_error", Message: "System fail internally"}
	ErrResourceNotFound = AppError{StatusCode: http.StatusNotFound, Type: "resource_not_found", Message: "Resource not found"}
	ErrInvalidRequest   = AppError{StatusCode: http.StatusConflict, Type: "invalid_request", Message: "Invalid request"}
)

func InvalidBody(err error) AppError {
	return AppError{
		StatusCode: http.StatusBadRequest,
		Type:       "invalid_entity",
		Message:    err.Error(),
	}
}

func (e AppError) Send(w http.ResponseWriter) error {
	statusCode := e.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(e)
}
