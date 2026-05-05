// Package config handles application configuration loading from environment variables.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
	Env  string
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

	return &Config{
		Port: portNumber,
		Env:  environment,
	}, nil
}
