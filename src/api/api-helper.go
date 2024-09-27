package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func AllowCors(w http.ResponseWriter) {
	// Specify Content Type to receive as Json Format
	w.Header().Set("Content-Type", "application/json")
	// Set CORS headers to allow requests from all origins - Different Ports
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Allow Content-Type header
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func HandleOptionsPreflightRequests(w http.ResponseWriter, r *http.Request) bool {
	AllowCors(w)

    // Check if it's a preflight request (OPTIONS method)
    if r.Method == http.MethodOptions {
		// Respond with 200 OK for preflight requests
        w.WriteHeader(http.StatusOK)
        return true
    }
    return false
}

func WriteAsJson(w http.ResponseWriter, mappedResponse map[string]interface{}) error {
	jsonResponse, err := json.Marshal(mappedResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to encode response to frontend: %s", err.Error()), http.StatusInternalServerError)
		return fmt.Errorf("unable to encode response to frontend: %w", err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		return fmt.Errorf("unable to write response to frontend: %w", err)
	}

	return nil
}
