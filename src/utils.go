package main

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON sends a JSON response with the given status code and data.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// HandleError sends a JSON response for an error with the given status code.
func HandleError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}