package handlers

import (
	"fmt"
	"net/http"
	"os"
)

// WriteJSON writes JSON data to the response with the specified status code
func WriteJSON(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

// LoadFixture loads a JSON fixture file from the spec-examples directory
func LoadFixture(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load fixture %s: %w", path, err)
	}
	return data, nil
}

// WriteError writes an error response in JSON format
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
