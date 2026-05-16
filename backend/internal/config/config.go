// Package config handles application configuration loading from environment variables.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	Env         string
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080"
	}
	portNumber, portError := strconv.Atoi(portString)
	if portError != nil {
		return nil, portError
	}

	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "development"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/grubbin?sslmode=disable"
	}

	return &Config{
		Port:        portNumber,
		Env:         environment,
		DatabaseURL: databaseURL,
	}, nil
}
