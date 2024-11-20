package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// RespondWithError sends an error response with a specific message and status code
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response with the given status code
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// ParsePrice converts a string to a float, returns an error if invalid
func ParsePrice(price string) (float64, error) {
	price = strings.TrimSpace(price)
	return strconv.ParseFloat(price, 64)
}
