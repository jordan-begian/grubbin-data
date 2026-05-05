// Package controllers provides HTTP request handlers (MVC controllers) for the API.
package controllers

import (
	"encoding/json"
	"net/http"

	"grubbin-data/backend/internal/services"
)

// GreetingController handles HTTP requests for greetings
type GreetingController struct {
	greetingService *services.GreetingService
}

// NewGreetingController creates a new greeting controller
func NewGreetingController(greetingService *services.GreetingService) *GreetingController {
	return &GreetingController{
		greetingService: greetingService,
	}
}

// GetGreeting handles GET /hello requests
func (c *GreetingController) GetGreeting(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	greeting, err := c.greetingService.GenerateGreeting(r.Context(), name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, greeting)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, statusCode int, err error) {
	respondWithJSON(w, statusCode, map[string]string{"error": err.Error()})
}
