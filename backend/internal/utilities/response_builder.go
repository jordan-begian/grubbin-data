// Package utilities provides helper functions and utilities for the API, such as response builders, error handling, and other common tasks.
package utilities

import (
	"encoding/json"
	"net/http"
)

type ResponseBuilder struct{}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

func (*ResponseBuilder) JSON(
	response http.ResponseWriter,
	statusCode int,
	payload any,
) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	json.NewEncoder(response).Encode(payload)
}

func (builder *ResponseBuilder) Error(
	response http.ResponseWriter,
	statusCode int,
	err error,
) {
	builder.JSON(response, statusCode, map[string]string{"error": err.Error()})
}
