// Package main seeds the database with a test user for local development.
// Usage: go run cmd/seed/main.go
//
// This command is intended for local development only. It creates a single
// test user with linked profile and vehicle if they don't already exist.
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"grubbin-data/backend/internal/config"
	"grubbin-data/backend/internal/core"
	"grubbin-data/backend/internal/utilities"
)

func main() {
	ctx := context.Background()

	// Load configuration
	appConfig, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Connect to database
	pool, err := pgxpool.New(ctx, appConfig.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	// Check if seed user already exists
	var existingID string
	err = pool.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", "test").Scan(&existingID)
	if err == nil {
		slog.Info("Seed user 'test' already exists, skipping", "user_id", existingID)
		return
	}

	// Generate ULID (shared across user, profile, vehicle)
	id, err := utilities.GenerateULID()
	if err != nil {
		slog.Error("Failed to generate ULID", "error", err)
		os.Exit(1)
	}

	// Hash password with pepper
	hash, err := core.HashPassword("T3stPass135!", appConfig.PasswordPepper, appConfig.BcryptCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		os.Exit(1)
	}

	now := time.Now().UTC()

	// Insert user, profile, and vehicle atomically
	tx, err := pool.Begin(ctx)
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		os.Exit(1)
	}
	defer tx.Rollback(ctx)

	// Insert user
	_, err = tx.Exec(ctx,
		"INSERT INTO users (id, username, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		id, "test", hash, now,
	)
	if err != nil {
		slog.Error("Failed to insert user", "error", err)
		os.Exit(1)
	}

	// Insert profile (shares same ID as user)
	_, err = tx.Exec(ctx,
		"INSERT INTO profiles (id, first_name, last_name) VALUES ($1, $2, $3)",
		id, "Test", "User",
	)
	if err != nil {
		slog.Error("Failed to insert profile", "error", err)
		os.Exit(1)
	}

	// Insert vehicle (shares same ID as user)
	_, err = tx.Exec(ctx,
		"INSERT INTO vehicles (id, name, average_mpg) VALUES ($1, $2, $3)",
		id, "Ford Focus", 30.1,
	)
	if err != nil {
		slog.Error("Failed to insert vehicle", "error", err)
		os.Exit(1)
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		os.Exit(1)
	}

	slog.Info("Seed user created successfully",
		"user_id", id,
		"username", "test",
		"vehicle", "Ford Focus",
	)
}
