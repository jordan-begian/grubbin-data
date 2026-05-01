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

	httprouter "grubbin-data/backend/internal/adapters/http"
	"grubbin-data/backend/internal/adapters/http/handlers"
	"grubbin-data/backend/internal/config"
	"grubbin-data/backend/internal/services"
)

func main() {
	appConfig, configError := config.Load()
	if configError != nil {
		slog.Error("Failed to load configuration", "error", configError)
		os.Exit(1)
	}

	helloWorldService := services.NewHelloWorldService()
	helloWorldHandler := handlers.NewHelloWorldHandler(helloWorldService)
	router := httprouter.NewRouter(helloWorldHandler)

	serverAddress := fmt.Sprintf(":%d", appConfig.Port)
	httpServer := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	go func() {
		slog.Info("Starting server", "address", serverAddress, "environment", appConfig.Env)
		if listenError := httpServer.ListenAndServe(); listenError != nil && listenError != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", listenError)
			os.Exit(1)
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)
	<-shutdownChannel

	slog.Info("Shutting down server...")
	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if shutdownError := httpServer.Shutdown(shutdownContext); shutdownError != nil {
		slog.Error("Server forced to shutdown", "error", shutdownError)
		os.Exit(1)
	}

	slog.Info("Server exited properly")
}
