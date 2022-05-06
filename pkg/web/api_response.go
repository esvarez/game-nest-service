package web

import (
	"encoding/json"
	"net/http"
)

type ResponseAPI struct {
	Success bool `json:"-"`
	Status  int  `json:"-"`
	Result  any  `json:"result"`
}

var (
	NotContent = ResponseAPI{Success: true, Status: http.StatusNoContent}
)

func Success(result any, status int) *ResponseAPI {
	return &ResponseAPI{Success: true, Status: status, Result: result}
}

func (r *ResponseAPI) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)

	return json.NewEncoder(w).Encode(r.Result)
}
