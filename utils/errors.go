package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SendError(w http.ResponseWriter, statuscode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	json.NewEncoder(w).Encode(ErrorResponse{Status: statuscode, Message: message})
}
