// Package config handles application configuration loading from environment variables.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           int
	Env            string
	DatabaseURL    string
	BcryptCost     int
	PasswordPepper string
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

	bcryptCostString := os.Getenv("BCRYPT_COST")
	if bcryptCostString == "" {
		bcryptCostString = "12"
	}
	bcryptCost, bcryptCostError := strconv.Atoi(bcryptCostString)
	if bcryptCostError != nil {
		return nil, bcryptCostError
	}

	passwordPepper := os.Getenv("PASSWORD_PEPPER")

	return &Config{
		Port:           portNumber,
		Env:            environment,
		DatabaseURL:    databaseURL,
		BcryptCost:     bcryptCost,
		PasswordPepper: passwordPepper,
	}, nil
}
