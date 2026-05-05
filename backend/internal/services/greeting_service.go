// Package services orchestrates business logic and side effects, composing
// pure functions from core with external dependencies like databases and time.
package services

import (
	"context"
	"time"

	"grubbin-data/backend/internal/core"
	"grubbin-data/backend/internal/models"
)

// GreetingService handles greeting business logic
type GreetingService struct{}

// NewGreetingService creates a new greeting service
func NewGreetingService() *GreetingService {
	return &GreetingService{}
}

// GenerateGreeting creates a greeting for the given name
func (s *GreetingService) GenerateGreeting(ctx context.Context, name string) (*models.Greeting, error) {
	// Pure business logic from core
	message := core.GenerateGreeting(name)

	// Side effect: getting current time (imperative shell)
	return &models.Greeting{
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
