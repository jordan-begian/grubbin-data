// Package models defines domain types and data structures used throughout
// the application.
package models

// Greeting represents a greeting response
type Greeting struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
