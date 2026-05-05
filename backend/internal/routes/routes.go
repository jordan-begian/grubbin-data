// Package routes configures HTTP routing using the Chi router.
package routes

import (
	"net/http"
	"time"

	"grubbin-data/backend/internal/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes configures the Chi router with all routes
func SetupRoutes(greetingController *controllers.GreetingController) *chi.Mux {
	router := chi.NewRouter()

	// Middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Mount v1 routes
	router.Mount("/api/v1", apiV1Routes(greetingController))

	// Mount latest routes (alias to v1)
	router.Mount("/api", apiV1Routes(greetingController))

	return router
}

// apiV1Routes defines all v1 API routes
func apiV1Routes(greetingController *controllers.GreetingController) http.Handler {
	router := chi.NewRouter()
	router.Get("/hello", greetingController.GetGreeting)
	return router
}
