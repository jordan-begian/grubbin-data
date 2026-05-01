package core

// GenerateGreeting returns a greeting message for the given name.
// If name is empty, returns a default greeting.
func GenerateGreeting(name string) string {
	if name == "" {
		return "Hello from Go"
	}
	return "Hello, " + name + "!"
}
