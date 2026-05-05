// Package core provides pure business logic functions with no side effects.
// All functions in this package are deterministic and free of I/O,
// making them trivial to unit test.
package core

// GenerateGreeting returns a greeting message for the given name.
// If name is empty, returns a default greeting.
func GenerateGreeting(name string) string {
	if name == "" {
		return "Hello from Go"
	}
	return "Hello, " + name + "!"
}
