package services

import (
	"context"
	"time"

	"grubbin-data/backend/internal/core"
)

type GreetingResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type HelloWorldService struct{}

func NewHelloWorldService() *HelloWorldService {
	return &HelloWorldService{}
}

func (s *HelloWorldService) GenerateGreeting(ctx context.Context, name string) (*GreetingResponse, error) {
	greetingMessage := core.GenerateGreeting(name)
	return &GreetingResponse{
		Message:   greetingMessage,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
