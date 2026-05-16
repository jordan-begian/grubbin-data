// Package core provides pure business logic functions with no side effects.
// All functions in this package are deterministic and free of I/O,
// making them trivial to unit test.
package core

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password with bcrypt after applying a server-side pepper.
// The pepper is appended to the password before hashing, providing defense-in-depth
// if the database is compromised without the pepper.
func HashPassword(password string, pepper string, cost int) (string, error) {
	pepperedPassword := password + pepper
	hash, err := bcrypt.GenerateFromPassword([]byte(pepperedPassword), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword checks if a password matches a bcrypt hash, applying the same pepper.
func VerifyPassword(password string, pepper string, hash string) (bool, error) {
	pepperedPassword := password + pepper
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pepperedPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ValidatePassword checks if a password meets complexity requirements.
// Returns a slice of validation error messages. An empty slice means the password is valid.
func ValidatePassword(password string) []string {
	var errors []string

	if len(password) < 12 {
		errors = append(errors, "Password must be at least 12 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	if !hasDigit {
		errors = append(errors, "Password must contain at least one digit")
	}
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}

	// Check for common weak patterns
	if isCommonPassword(password) {
		errors = append(errors, "Password is too common or easily guessable")
	}

	return errors
}

// IsPasswordValid is a convenience function that returns true if the password is valid.
func IsPasswordValid(password string) bool {
	return len(ValidatePassword(password)) == 0
}

// commonPasswords is a small list of commonly used weak passwords to reject.
// In production, consider using a larger list or an external service.
var commonPasswords = []string{
	"password", "123456", "qwerty", "admin", "letmein",
	"welcome", "monkey", "dragon", "master", "sunshine",
}

var commonPasswordPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^password`),
	regexp.MustCompile(`(?i)^123456`),
	regexp.MustCompile(`(?i)^qwerty`),
	regexp.MustCompile(`(?i)^admin`),
	regexp.MustCompile(`(?i)^letmein`),
}

func isCommonPassword(password string) bool {
	for _, common := range commonPasswords {
		if regexp.MustCompile(`(?i)`+regexp.QuoteMeta(common)).MatchString(password) {
			return true
		}
	}

	// Check for sequential patterns
	if regexp.MustCompile(`(012|123|234|345|456|567|678|789|890|abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)`).MatchString(password) {
		return true
	}

	return false
}
