package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"grubbin-data/backend/internal/config"
	"grubbin-data/backend/internal/controllers"
	"grubbin-data/backend/internal/routes"
	"grubbin-data/backend/internal/services"
	"grubbin-data/backend/internal/utilities"
)

func main() {
	// Load configuration
	appConfig, configError := config.Load()
	if configError != nil {
		slog.Error("Failed to load configuration", "error", configError)
		os.Exit(1)
	}

	// Initialize database connection pool
	dbPool, dbPoolError := pgxpool.New(context.Background(), appConfig.DatabaseURL)
	if dbPoolError != nil {
		slog.Error("Failed to connect to database", "error", dbPoolError)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Response builder (utilities)
	responseBuilder := utilities.NewResponseBuilder()

	// Service layer (business logic + orchestration)
	greetingService := services.NewGreetingService()

	// Controller layer (HTTP handlers)
	greetingController := controllers.NewGreetingController(
		responseBuilder,
		greetingService,
	)

	// Routes (Chi router setup)
	router := routes.SetupRoutes(greetingController)

	// Create HTTP server
	serverAddress := fmt.Sprintf(":%d", appConfig.Port)
	httpServer := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Starting server", "address", serverAddress, "environment", appConfig.Env)
		if listenError := httpServer.ListenAndServe(); listenError != nil && listenError != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", listenError)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)
	<-shutdownChannel

	// Graceful shutdown
	slog.Info("Shutting down server...")
	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if shutdownError := httpServer.Shutdown(shutdownContext); shutdownError != nil {
		slog.Error("Server forced to shutdown", "error", shutdownError)
		os.Exit(1)
	}

	slog.Info("Server exited properly")
}
