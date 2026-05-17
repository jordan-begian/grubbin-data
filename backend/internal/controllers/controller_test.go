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

	"github.com/go-chi/chi/v5"

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

// mockDeliveryService implements the DeliveryService interface for testing.
type mockDeliveryService struct {
	createDeliveryFunc   func(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error)
	getDeliveryFunc      func(ctx context.Context, userID, deliveryID string) (*models.DeliveryResponse, error)
	getDeliveriesFunc    func(ctx context.Context, userID string, startDate, endDate *time.Time) (*models.DeliveryListResponse, error)
	updateDeliveriesFunc func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error
	deleteDeliveriesFunc func(ctx context.Context, userID string, ids []string) error
}

func (m *mockDeliveryService) CreateDelivery(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error) {
	return m.createDeliveryFunc(ctx, userID, req)
}

func (m *mockDeliveryService) GetDelivery(ctx context.Context, userID, deliveryID string) (*models.DeliveryResponse, error) {
	return m.getDeliveryFunc(ctx, userID, deliveryID)
}

func (m *mockDeliveryService) GetDeliveries(ctx context.Context, userID string, startDate, endDate *time.Time) (*models.DeliveryListResponse, error) {
	return m.getDeliveriesFunc(ctx, userID, startDate, endDate)
}

func (m *mockDeliveryService) UpdateDeliveries(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
	return m.updateDeliveriesFunc(ctx, userID, updates)
}

func (m *mockDeliveryService) DeleteDeliveries(ctx context.Context, userID string, ids []string) error {
	return m.deleteDeliveriesFunc(ctx, userID, ids)
}

func newTestController(authService AuthService, greetingService GreetingService) *Controller {
	return NewController(utilities.NewResponseBuilder(), authService, greetingService, &mockDeliveryService{})
}

// withChiParams adds Chi URL parameters to the request context for testing.
func withChiParams(request *http.Request, params map[string]string) *http.Request {
	routeCtx := chi.NewRouteContext()
	for key, value := range params {
		routeCtx.URLParams.Add(key, value)
	}
	return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
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

// --- Delivery Controller Tests ---

func TestController_CreateDelivery_Success(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	deliveryService := &mockDeliveryService{
		createDeliveryFunc: func(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error) {
			if userID != "user-1" {
				t.Errorf("UserID = %s, want user-1", userID)
			}
			return &models.DeliveryResponse{
				ID:     "delivery-1",
				UserID: userID,
				Start:  req.Start,
				End:    req.End,
				Pickup: models.PickupResponse{ID: "pickup-1", Name: req.Pickup.Name, Lat: req.Pickup.Lat, Lon: req.Pickup.Lon},
				Dropoff: models.DropoffResponse{ID: "dropoff-1", Lat: req.Dropoff.Lat, Lon: req.Dropoff.Lon},
				Earnings: models.EarningsResponse{ID: "earnings-1", Tip: req.Earnings.Tip, Base: req.Earnings.Base},
			}, nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	body := `{
		"start": "` + now.Format(time.RFC3339) + `",
		"end": "` + later.Format(time.RFC3339) + `",
		"pickup": {"name": "Restaurant", "lat": 40.7, "lon": -74.0},
		"dropoff": {"lat": 40.8, "lon": -73.9},
		"earnings": {"tip": 300, "base": 500}
	}`
	request := httptest.NewRequest(http.MethodPost, "/users/user-1/deliveries", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.CreateDelivery(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var resp models.DeliveryResponse
	if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.ID != "delivery-1" {
		t.Errorf("ID = %s, want delivery-1", resp.ID)
	}
	if resp.Pickup.Name != "Restaurant" {
		t.Errorf("Pickup.Name = %s, want Restaurant", resp.Pickup.Name)
	}
}

func TestController_CreateDelivery_MissingUserID(t *testing.T) {
	deliveryService := &mockDeliveryService{
		createDeliveryFunc: func(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error) {
			t.Error("CreateDelivery should not be called without user ID")
			return nil, nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	body := `{"start": "2024-01-01T00:00:00Z", "end": "2024-01-01T00:30:00Z", "pickup": {"name": "R", "lat": 1, "lon": 1}, "dropoff": {"lat": 2, "lon": 2}, "earnings": {"tip": 0, "base": 0}}`
	request := httptest.NewRequest(http.MethodPost, "/users//deliveries", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	controller.CreateDelivery(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_CreateDelivery_InvalidJSON(t *testing.T) {
	deliveryService := &mockDeliveryService{
		createDeliveryFunc: func(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error) {
			t.Error("CreateDelivery should not be called with invalid JSON")
			return nil, nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	body := `{"start": "invalid`
	request := httptest.NewRequest(http.MethodPost, "/users/user-1/deliveries", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.CreateDelivery(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_CreateDelivery_ValidationError(t *testing.T) {
	deliveryService := &mockDeliveryService{
		createDeliveryFunc: func(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error) {
			return nil, errors.New("validation failed: [start time is required]")
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	body := `{"end": "2024-01-01T00:30:00Z", "pickup": {"name": "R", "lat": 1, "lon": 1}, "dropoff": {"lat": 2, "lon": 2}, "earnings": {"tip": 0, "base": 0}}`
	request := httptest.NewRequest(http.MethodPost, "/users/user-1/deliveries", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.CreateDelivery(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_GetDelivery_Success(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	deliveryService := &mockDeliveryService{
		getDeliveryFunc: func(ctx context.Context, userID, deliveryID string) (*models.DeliveryResponse, error) {
			if userID != "user-1" || deliveryID != "delivery-1" {
				t.Errorf("UserID = %s, DeliveryID = %s", userID, deliveryID)
			}
			return &models.DeliveryResponse{
				ID:       "delivery-1",
				UserID:   userID,
				Start:    now,
				End:      later,
				Pickup:   models.PickupResponse{ID: "pickup-1", Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff:  models.DropoffResponse{ID: "dropoff-1", Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsResponse{ID: "earnings-1", Tip: 300, Base: 500},
			}, nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	request := httptest.NewRequest(http.MethodGet, "/users/user-1/deliveries/delivery-1", nil)
	request = withChiParams(request, map[string]string{"userId": "user-1", "deliveryId": "delivery-1"})
	response := httptest.NewRecorder()

	controller.GetDelivery(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	var resp models.DeliveryResponse
	if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.ID != "delivery-1" {
		t.Errorf("ID = %s, want delivery-1", resp.ID)
	}
}

func TestController_GetDelivery_NotFound(t *testing.T) {
	deliveryService := &mockDeliveryService{
		getDeliveryFunc: func(ctx context.Context, userID, deliveryID string) (*models.DeliveryResponse, error) {
			return nil, errors.New("no rows in result set")
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	request := httptest.NewRequest(http.MethodGet, "/users/user-1/deliveries/nonexistent", nil)
	request = withChiParams(request, map[string]string{"userId": "user-1", "deliveryId": "nonexistent"})
	response := httptest.NewRecorder()

	controller.GetDelivery(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestController_GetDeliveries_Success(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	deliveryService := &mockDeliveryService{
		getDeliveriesFunc: func(ctx context.Context, userID string, startDate, endDate *time.Time) (*models.DeliveryListResponse, error) {
			return &models.DeliveryListResponse{
				Deliveries: []models.DeliveryResponse{
					{
						ID:       "delivery-1",
						UserID:   userID,
						Start:    now,
						End:      later,
						Pickup:   models.PickupResponse{ID: "pickup-1", Name: "Restaurant", Lat: 40.7, Lon: -74.0},
						Dropoff:  models.DropoffResponse{ID: "dropoff-1", Lat: 40.8, Lon: -73.9},
						Earnings: models.EarningsResponse{ID: "earnings-1", Tip: 300, Base: 500},
					},
				},
				Stats: models.DeliveryStats{TotalTime: 1800, TotalMiles: 5.2},
			}, nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	request := httptest.NewRequest(http.MethodGet, "/users/user-1/deliveries", nil)
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.GetDeliveries(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	var resp models.DeliveryListResponse
	if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp.Deliveries) != 1 {
		t.Errorf("Deliveries count = %d, want 1", len(resp.Deliveries))
	}
	if resp.Stats.TotalTime != 1800 {
		t.Errorf("Stats.TotalTime = %d, want 1800", resp.Stats.TotalTime)
	}
}

func TestController_GetDeliveries_WithDateRange(t *testing.T) {
	deliveryService := &mockDeliveryService{
		getDeliveriesFunc: func(ctx context.Context, userID string, startDate, endDate *time.Time) (*models.DeliveryListResponse, error) {
			if startDate == nil || endDate == nil {
				t.Error("Expected startDate and endDate to be set")
			}
			return &models.DeliveryListResponse{
				Deliveries: []models.DeliveryResponse{},
				Stats:      models.DeliveryStats{},
			}, nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	request := httptest.NewRequest(http.MethodGet, "/users/user-1/deliveries?start_date=2024-01-01T00:00:00Z&end_date=2024-12-31T23:59:59Z", nil)
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.GetDeliveries(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestController_GetDeliveries_InvalidDateFormat(t *testing.T) {
	deliveryService := &mockDeliveryService{
		getDeliveriesFunc: func(ctx context.Context, userID string, startDate, endDate *time.Time) (*models.DeliveryListResponse, error) {
			t.Error("GetDeliveries should not be called with invalid date format")
			return nil, nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	request := httptest.NewRequest(http.MethodGet, "/users/user-1/deliveries?start_date=not-a-date", nil)
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.GetDeliveries(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_UpdateDeliveries_Success(t *testing.T) {
	deliveryService := &mockDeliveryService{
		updateDeliveriesFunc: func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
			if len(updates) != 2 {
				t.Errorf("Updates count = %d, want 2", len(updates))
			}
			return nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	body := `[{"id": "delivery-1", "note": "Updated"}, {"id": "delivery-2", "note": "Another"}]`
	request := httptest.NewRequest(http.MethodPatch, "/users/user-1/deliveries", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.UpdateDeliveries(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestController_UpdateDeliveries_EmptyArray(t *testing.T) {
	deliveryService := &mockDeliveryService{
		updateDeliveriesFunc: func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
			t.Error("UpdateDeliveries should not be called with empty array")
			return nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	body := `[]`
	request := httptest.NewRequest(http.MethodPatch, "/users/user-1/deliveries", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.UpdateDeliveries(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_UpdateDeliveries_InvalidJSON(t *testing.T) {
	deliveryService := &mockDeliveryService{
		updateDeliveriesFunc: func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
			t.Error("UpdateDeliveries should not be called with invalid JSON")
			return nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	body := `[{"id": "delivery-1"`
	request := httptest.NewRequest(http.MethodPatch, "/users/user-1/deliveries", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.UpdateDeliveries(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestController_DeleteDeliveries_Success(t *testing.T) {
	deliveryService := &mockDeliveryService{
		deleteDeliveriesFunc: func(ctx context.Context, userID string, ids []string) error {
			if len(ids) != 2 {
				t.Errorf("IDs count = %d, want 2", len(ids))
			}
			return nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	request := httptest.NewRequest(http.MethodDelete, "/users/user-1/deliveries?id=delivery-1&id=delivery-2", nil)
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.DeleteDeliveries(response, request)

	if response.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, response.Code)
	}
}

func TestController_DeleteDeliveries_MissingIDs(t *testing.T) {
	deliveryService := &mockDeliveryService{
		deleteDeliveriesFunc: func(ctx context.Context, userID string, ids []string) error {
			t.Error("DeleteDeliveries should not be called without IDs")
			return nil
		},
	}

	controller := NewController(utilities.NewResponseBuilder(), nil, nil, deliveryService)
	request := httptest.NewRequest(http.MethodDelete, "/users/user-1/deliveries", nil)
	request = withChiParams(request, map[string]string{"userId": "user-1"})
	response := httptest.NewRecorder()

	controller.DeleteDeliveries(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}
