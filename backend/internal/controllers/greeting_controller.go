// Package controllers provides HTTP request handlers (MVC controllers) for the API.
package controllers

import (
	"net/http"

	"grubbin-data/backend/internal/services"
	"grubbin-data/backend/internal/utilities"
)

// GreetingController handles HTTP requests for greetings
type GreetingController struct {
	responseBuilder *utilities.ResponseBuilder
	greetingService *services.GreetingService
}

// NewGreetingController creates a new greeting controller
func NewGreetingController(
	responseBuilder *utilities.ResponseBuilder,
	greetingService *services.GreetingService,
) *GreetingController {
	return &GreetingController{
		responseBuilder: responseBuilder,
		greetingService: greetingService,
	}
}

// GetGreeting handles GET /hello requests
func (controller *GreetingController) GetGreeting(
	response http.ResponseWriter,
	request *http.Request,
) {
	name := request.URL.Query().Get("name")

	greeting, err := controller.greetingService.GenerateGreeting(request.Context(), name)
	if err != nil {
		controller.responseBuilder.Error(response, http.StatusInternalServerError, err)
		return
	}

	controller.responseBuilder.JSON(response, http.StatusOK, greeting)
}
