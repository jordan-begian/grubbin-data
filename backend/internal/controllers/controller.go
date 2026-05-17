// Package controllers provides HTTP request handlers (MVC controllers) for the API.
package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

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

// DeliveryService defines the delivery operations the controller depends on.
type DeliveryService interface {
	CreateDelivery(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error)
	GetDelivery(ctx context.Context, userID, deliveryID string) (*models.DeliveryResponse, error)
	GetDeliveries(ctx context.Context, userID string, startDate, endDate *time.Time) (*models.DeliveryListResponse, error)
	UpdateDeliveries(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error
	DeleteDeliveries(ctx context.Context, userID string, ids []string) error
}

// Controller handles HTTP requests for all API endpoints.
type Controller struct {
	responseBuilder  *utilities.ResponseBuilder
	authService      AuthService
	greetingService  GreetingService
	deliveryService  DeliveryService
}

// NewController creates a new API controller.
func NewController(
	responseBuilder *utilities.ResponseBuilder,
	authService AuthService,
	greetingService GreetingService,
	deliveryService DeliveryService,
) *Controller {
	return &Controller{
		responseBuilder: responseBuilder,
		authService:     authService,
		greetingService: greetingService,
		deliveryService: deliveryService,
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

// CreateDelivery handles POST /users/{userId}/deliveries requests.
func (c *Controller) CreateDelivery(response http.ResponseWriter, request *http.Request) {
	userID := chi.URLParam(request, "userId")
	if userID == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("user id is required"))
		return
	}

	var req models.CreateDeliveryRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	delivery, err := c.deliveryService.CreateDelivery(request.Context(), userID, req)
	if err != nil {
		if isValidationError(err) {
			c.responseBuilder.Error(response, http.StatusBadRequest, err)
			return
		}
		c.responseBuilder.Error(response, http.StatusInternalServerError, errors.New("failed to create delivery"))
		return
	}

	c.responseBuilder.JSON(response, http.StatusCreated, delivery)
}

// GetDelivery handles GET /users/{userId}/deliveries/{deliveryId} requests.
func (c *Controller) GetDelivery(response http.ResponseWriter, request *http.Request) {
	userID := chi.URLParam(request, "userId")
	deliveryID := chi.URLParam(request, "deliveryId")

	if userID == "" || deliveryID == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("user id and delivery id are required"))
		return
	}

	delivery, err := c.deliveryService.GetDelivery(request.Context(), userID, deliveryID)
	if err != nil {
		if isNotFoundError(err) {
			c.responseBuilder.Error(response, http.StatusNotFound, errors.New("delivery not found"))
			return
		}
		c.responseBuilder.Error(response, http.StatusInternalServerError, errors.New("failed to get delivery"))
		return
	}

	c.responseBuilder.JSON(response, http.StatusOK, delivery)
}

// GetDeliveries handles GET /users/{userId}/deliveries requests.
func (c *Controller) GetDeliveries(response http.ResponseWriter, request *http.Request) {
	userID := chi.URLParam(request, "userId")
	if userID == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("user id is required"))
		return
	}

	var startDate, endDate *time.Time
	startStr := request.URL.Query().Get("start_date")
	endStr := request.URL.Query().Get("end_date")

	if startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("invalid start_date format, use ISO 8601 (e.g. 2024-01-01T00:00:00Z)"))
			return
		}
		startDate = &t
	}

	if endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("invalid end_date format, use ISO 8601 (e.g. 2024-12-31T23:59:59Z)"))
			return
		}
		endDate = &t
	}

	result, err := c.deliveryService.GetDeliveries(request.Context(), userID, startDate, endDate)
	if err != nil {
		c.responseBuilder.Error(response, http.StatusInternalServerError, errors.New("failed to get deliveries"))
		return
	}

	c.responseBuilder.JSON(response, http.StatusOK, result)
}

// UpdateDeliveries handles PATCH /users/{userId}/deliveries requests.
func (c *Controller) UpdateDeliveries(response http.ResponseWriter, request *http.Request) {
	userID := chi.URLParam(request, "userId")
	if userID == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("user id is required"))
		return
	}

	var updates []models.UpdateDeliveryRequest
	if err := json.NewDecoder(request.Body).Decode(&updates); err != nil {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("invalid request body, expected array of updates"))
		return
	}

	if len(updates) == 0 {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("at least one delivery update is required"))
		return
	}

	if err := c.deliveryService.UpdateDeliveries(request.Context(), userID, updates); err != nil {
		if isValidationError(err) {
			c.responseBuilder.Error(response, http.StatusBadRequest, err)
			return
		}
		if isNotFoundError(err) {
			c.responseBuilder.Error(response, http.StatusNotFound, errors.New("one or more deliveries not found"))
			return
		}
		c.responseBuilder.Error(response, http.StatusInternalServerError, errors.New("failed to update deliveries"))
		return
	}

	c.responseBuilder.JSON(response, http.StatusOK, map[string]string{"message": "deliveries updated"})
}

// DeleteDeliveries handles DELETE /users/{userId}/deliveries?id=id1&id=id2 requests.
func (c *Controller) DeleteDeliveries(response http.ResponseWriter, request *http.Request) {
	userID := chi.URLParam(request, "userId")
	if userID == "" {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("user id is required"))
		return
	}

	ids := request.URL.Query()["id"]
	if len(ids) == 0 {
		c.responseBuilder.Error(response, http.StatusBadRequest, errors.New("at least one delivery id is required (use ?id=...)"))
		return
	}

	if err := c.deliveryService.DeleteDeliveries(request.Context(), userID, ids); err != nil {
		if isNotFoundError(err) {
			c.responseBuilder.Error(response, http.StatusNotFound, errors.New("one or more deliveries not found"))
			return
		}
		c.responseBuilder.Error(response, http.StatusInternalServerError, errors.New("failed to delete deliveries"))
		return
	}

	response.WriteHeader(http.StatusNoContent)
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

// isValidationError checks if an error is a validation failure.
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "validation failed")
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

// isNotFoundError checks if an error indicates a record was not found.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "no rows") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "ownership")
}
