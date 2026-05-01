package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"grubbin-data/backend/internal/adapters/http/handlers"
)

func NewRouter(helloHandler *handlers.HelloWorldHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Mount v1 routes
	router.Mount("/api/v1", helloRoutes(helloHandler))

	// Mount latest routes (alias to v1)
	router.Mount("/api", helloRoutes(helloHandler))

	return router
}

func helloRoutes(helloHandler *handlers.HelloWorldHandler) http.Handler {
	router := chi.NewRouter()
	router.Get("/hello", helloHandler.GetGreeting)
	return router
}
