package main

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON sends a JSON response with the given status code and data.
// w: the HTTP response writer
// code: the HTTP status code (e.g., 200, 404)
// payload: the data to encode as JSON and send in the response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json") // Set response content type to JSON
	w.WriteHeader(code)                                // Set the HTTP status code
	json.NewEncoder(w).Encode(payload)                 // Encode the payload as JSON and write to response
}

// HandleError sends a JSON response for an error with the given status code.
// w: the HTTP response writer
// code: the HTTP status code (e.g., 400, 500)
// message: the error message to send in the response
func HandleError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message}) // Send error message as JSON
}
