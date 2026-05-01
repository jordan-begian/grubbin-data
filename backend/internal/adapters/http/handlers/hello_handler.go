package handlers

import (
	"encoding/json"
	"net/http"

	"grubbin-data/backend/internal/services"
)

type HelloWorldHandler struct {
	helloWorldService *services.HelloWorldService
}

func NewHelloWorldHandler(helloWorldService *services.HelloWorldService) *HelloWorldHandler {
	return &HelloWorldHandler{
		helloWorldService: helloWorldService,
	}
}

func (h *HelloWorldHandler) GetGreeting(responseWriter http.ResponseWriter, request *http.Request) {
	requestedName := request.URL.Query().Get("name")

	greetingResponse, greetingError := h.helloWorldService.GenerateGreeting(request.Context(), requestedName)
	if greetingError != nil {
		h.respondWithError(responseWriter, http.StatusInternalServerError, greetingError)
		return
	}

	h.respondWithJSON(responseWriter, http.StatusOK, greetingResponse)
}

func (h *HelloWorldHandler) respondWithJSON(responseWriter http.ResponseWriter, statusCode int, payload interface{}) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(payload)
}

func (h *HelloWorldHandler) respondWithError(responseWriter http.ResponseWriter, statusCode int, err error) {
	h.respondWithJSON(responseWriter, statusCode, map[string]string{"error": err.Error()})
}
