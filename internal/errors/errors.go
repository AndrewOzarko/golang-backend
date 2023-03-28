package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err error, status int) {
	response := ErrorResponse{
		Error: err.Error(),
	}
	responseBody, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(responseBody)
}
