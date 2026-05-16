package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"grubbin-data/backend/internal/models"
	"grubbin-data/backend/internal/utilities"
)

// mockAuthService implements the AuthService interface for testing.
type mockAuthService struct {
	registerUserFunc     func(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error)
	authenticateUserFunc func(ctx context.Context, username string, password string) (*models.User, error)
}

func (m *mockAuthService) RegisterUser(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error) {
	return m.registerUserFunc(ctx, username, password, firstName, lastName, vehicleName, vehicleMPG)
}

func (m *mockAuthService) AuthenticateUser(ctx context.Context, username string, password string) (*models.User, error) {
	return m.authenticateUserFunc(ctx, username, password)
}

// mockGreetingService implements the GreetingService interface for testing.
type mockGreetingService struct {
	generateGreetingFunc func(ctx context.Context, name string) (*models.Greeting, error)
}

func (m *mockGreetingService) GenerateGreeting(ctx context.Context, name string) (*models.Greeting, error) {
	return m.generateGreetingFunc(ctx, name)
}

func newTestController(authService AuthService, greetingService GreetingService) *Controller {
	return NewController(utilities.NewResponseBuilder(), authService, greetingService)
}

func TestController_RegisterUser_Success(t *testing.T) {
	created := time.Now().UTC()
	authService := &mockAuthService{
		registerUserFunc: func(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error) {
			return &models.User{
				ID:       "01HV8J3K2M4N5P6Q7R8S9T0UV1",
				Username: username,
				Created:  created,
				Profile: &models.Profile{
					ID:        "01HV8J3K2M4N5P6Q7R8S9T0UV1",
					FirstName: firstName,
					LastName:  lastName,
				},
			}, nil
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"testuser","password":"Xk9#mP2$vL7@qR4!","first_name":"John","last_name":"Doe"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var resp models.UserResponse
	if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", resp.Username)
	}
	if resp.ID == "" {
		t.Error("Expected user ID in response")
	}
	if resp.Profile == nil {
		t.Fatal("Expected profile in response")
	}
	if resp.Profile.FirstName != "John" {
		t.Errorf("Expected first name 'John', got %q", resp.Profile.FirstName)
	}
	if resp.Profile.LastName != "Doe" {
		t.Errorf("Expected last name 'Doe', got %q", resp.Profile.LastName)
	}
}

func TestController_RegisterUser_WithVehicle(t *testing.T) {
	created := time.Now().UTC()
	authService := &mockAuthService{
		registerUserFunc: func(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error) {
			return &models.User{
				ID:       "01HV8J3K2M4N5P6Q7R8S9T0UV1",
				Username: username,
				Created:  created,
				Profile: &models.Profile{
					ID:        "01HV8J3K2M4N5P6Q7R8S9T0UV1",
					FirstName: firstName,
					LastName:  lastName,
					Vehicle: &models.Vehicle{
						ID:         "01HV8J3K2M4N5P6Q7R8S9T0UV1",
						Name:       *vehicleName,
						AverageMPG: *vehicleMPG,
					},
				},
			}, nil
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"testuser","password":"Xk9#mP2$vL7@qR4!","first_name":"John","last_name":"Doe","vehicle_name":"Honda Civic","vehicle_mpg":25.5}`
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var resp models.UserResponse
	if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.Profile == nil {
		t.Fatal("Expected profile in response")
	}
	if resp.Profile.Vehicle == nil {
		t.Fatal("Expected vehicle in response")
	}
	if resp.Profile.Vehicle.Name != "Honda Civic" {
		t.Errorf("Expected vehicle name 'Honda Civic', got %q", resp.Profile.Vehicle.Name)
	}
	if resp.Profile.Vehicle.AverageMPG != 25.5 {
		t.Errorf("Expected vehicle MPG 25.5, got %f", resp.Profile.Vehicle.AverageMPG)
	}
}

func TestController_RegisterUser_InvalidPassword(t *testing.T) {
	authService := &mockAuthService{
		registerUserFunc: func(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error) {
			return nil, errors.New("password validation failed: too short")
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"testuser","password":"short","first_name":"John","last_name":"Doe"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_RegisterUser_DuplicateUsername(t *testing.T) {
	authService := &mockAuthService{
		registerUserFunc: func(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error) {
			return nil, errors.New("create user: unique constraint violation: duplicate key value")
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"existinguser","password":"Xk9#mP2$vL7@qR4!","first_name":"John","last_name":"Doe"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, response.Code)
	}
}

func TestController_RegisterUser_RepoError(t *testing.T) {
	authService := &mockAuthService{
		registerUserFunc: func(ctx context.Context, username string, password string, firstName string, lastName string, vehicleName *string, vehicleMPG *float64) (*models.User, error) {
			return nil, errors.New("database connection failed")
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"testuser","password":"Xk9#mP2$vL7@qR4!","first_name":"John","last_name":"Doe"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}

func TestController_RegisterUser_MalformedJSON(t *testing.T) {
	controller := newTestController(nil, nil)
	body := `{"username": "testuser", "password":` // malformed
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_RegisterUser_MissingFields(t *testing.T) {
	controller := newTestController(nil, nil)
	body := `{"username":"","password":"","first_name":"","last_name":""}`
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_RegisterUser_MissingNames(t *testing.T) {
	controller := newTestController(nil, nil)
	body := `{"username":"testuser","password":"Xk9#mP2$vL7@qR4!"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.RegisterUser(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_Login_Success(t *testing.T) {
	created := time.Now().UTC()
	authService := &mockAuthService{
		authenticateUserFunc: func(ctx context.Context, username string, password string) (*models.User, error) {
			return &models.User{
				ID:       "01HV8J3K2M4N5P6Q7R8S9T0UV1",
				Username: username,
				Created:  created,
			}, nil
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"testuser","password":"MySecureP@ssw0rd123"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.Login(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	var resp models.UserResponse
	if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", resp.Username)
	}
}

func TestController_Login_WrongPassword(t *testing.T) {
	authService := &mockAuthService{
		authenticateUserFunc: func(ctx context.Context, username string, password string) (*models.User, error) {
			return nil, errors.New("invalid credentials")
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"testuser","password":"WrongPassword456!"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.Login(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestController_Login_UserNotFound(t *testing.T) {
	authService := &mockAuthService{
		authenticateUserFunc: func(ctx context.Context, username string, password string) (*models.User, error) {
			return nil, errors.New("authentication failed: user not found")
		},
	}

	controller := newTestController(authService, nil)
	body := `{"username":"unknownuser","password":"Xk9#mP2$vL7@qR4!"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.Login(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestController_Login_MalformedJSON(t *testing.T) {
	controller := newTestController(nil, nil)
	body := `{"username": "testuser", "password":` // malformed
	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.Login(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_Login_MissingFields(t *testing.T) {
	controller := newTestController(nil, nil)
	body := `{"password":"secret"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.Login(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_GetGreeting_Success(t *testing.T) {
	greetingService := &mockGreetingService{
		generateGreetingFunc: func(ctx context.Context, name string) (*models.Greeting, error) {
			return &models.Greeting{
				Message:   "Hello, Test!",
				Timestamp: "2024-01-01T00:00:00Z",
			}, nil
		},
	}

	controller := newTestController(nil, greetingService)
	request := httptest.NewRequest(http.MethodGet, "/hello?name=Test", nil)
	response := httptest.NewRecorder()

	controller.GetGreeting(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	var resp models.Greeting
	if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.Message != "Hello, Test!" {
		t.Errorf("Expected message 'Hello, Test!', got %q", resp.Message)
	}
}

func TestController_GetGreeting_ServiceError(t *testing.T) {
	greetingService := &mockGreetingService{
		generateGreetingFunc: func(ctx context.Context, name string) (*models.Greeting, error) {
			return nil, errors.New("service unavailable")
		},
	}

	controller := newTestController(nil, greetingService)
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	response := httptest.NewRecorder()

	controller.GetGreeting(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}

func TestController_toUserResponse(t *testing.T) {
	created := time.Now().UTC()
	user := &models.User{
		ID:       "01HV8J3K2M4N5P6Q7R8S9T0UV1",
		Username: "testuser",
		Password: func() *string { s := "should-not-appear"; return &s }(),
		Created:  created,
	}

	resp := toUserResponse(user)

	if resp.ID != user.ID {
		t.Errorf("Expected ID %q, got %q", user.ID, resp.ID)
	}
	if resp.Username != user.Username {
		t.Errorf("Expected username %q, got %q", user.Username, resp.Username)
	}
	if resp.Created != user.Created {
		t.Errorf("Expected created %v, got %v", user.Created, resp.Created)
	}
	// Verify Password field is not present on UserResponse struct
	// (compile-time check via struct definition)
}
