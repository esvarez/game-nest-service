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

func (e AppError) Send(w http.ResponseWriter) error {
	statusCode := e.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(e)
}
