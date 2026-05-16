// Package controllers provides HTTP request handlers (MVC controllers) for the API.
package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"grubbin-data/backend/internal/models"
	"grubbin-data/backend/internal/utilities"
)

// AuthService defines the authentication operations the controller depends on.
// This interface allows mocking the service layer in controller tests.
type AuthService interface {
	RegisterUser(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error)
	AuthenticateUser(ctx context.Context, username string, password string) (*models.User, error)
}

// GreetingService defines the greeting operations the controller depends on.
type GreetingService interface {
	GenerateGreeting(ctx context.Context, name string) (*models.Greeting, error)
}

// Controller handles HTTP requests for all API endpoints.
type Controller struct {
	responseBuilder  *utilities.ResponseBuilder
	authService      AuthService
	greetingService  GreetingService
}

// NewController creates a new API controller.
func NewController(
	responseBuilder *utilities.ResponseBuilder,
	authService AuthService,
	greetingService GreetingService,
) *Controller {
	return &Controller{
		responseBuilder: responseBuilder,
		authService:     authService,
		greetingService: greetingService,
	}
}

// RegisterUser handles POST /auth/register requests.
func (c *Controller) RegisterUser(response http.ResponseWriter, request *http.Request) {
	var req models.RegisterUserRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	if req.Username == "" || req.Password == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("username and password are required"))
		return
	}

	if req.FirstName == "" || req.LastName == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("first name and last name are required"))
		return
	}

	user, err := c.authService.RegisterUser(request.Context(), req.Username, req.Password, req.FirstName, req.LastName, req.VehicleName, req.VehicleMPG)
	if err != nil {
		if isValidationError(err) {
			c.responseBuilder.Error(response, http.StatusBadRequest, err)
			return
		}
		if isDuplicateError(err) {
			c.responseBuilder.Error(response, http.StatusConflict, errors.New("username already exists"))
			return
		}
		c.responseBuilder.Error(response, http.StatusInternalServerError, errors.New("failed to create user"))
		return
	}

	c.responseBuilder.JSON(response, http.StatusCreated, toUserResponse(user))
}

// GetGreeting handles GET /hello requests.
func (c *Controller) GetGreeting(response http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")

	greeting, err := c.greetingService.GenerateGreeting(request.Context(), name)
	if err != nil {
		c.responseBuilder.Error(response, http.StatusInternalServerError, err)
		return
	}

	c.responseBuilder.JSON(response, http.StatusOK, greeting)
}

// Login handles POST /auth/login requests.
func (c *Controller) Login(response http.ResponseWriter, request *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	if req.Username == "" || req.Password == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("username and password are required"))
		return
	}

	user, err := c.authService.AuthenticateUser(request.Context(), req.Username, req.Password)
	if err != nil {
		c.responseBuilder.Error(response, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	c.responseBuilder.JSON(response, http.StatusOK, toUserResponse(user))
}

// toUserResponse converts a models.User to a models.UserResponse,
// ensuring sensitive fields like Password are never exposed.
func toUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Created:  user.Created,
		Updated:  user.Updated,
		Profile:  user.Profile,
	}
}

// isValidationError checks if an error is a password validation failure.
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "password validation failed")
}

// isDuplicateError checks if an error indicates a unique constraint violation.
func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "already exists")
}
