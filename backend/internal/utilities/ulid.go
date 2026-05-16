package utilities

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/oklog/ulid/v2"
)

// ulidEntropyReader allows test injection of the entropy source.
var ulidEntropyReader io.Reader = rand.Reader

// GenerateULID creates a new cryptographically secure, sortable ULID.
// Uses the current time for the timestamp component.
func GenerateULID() (string, error) {
	entropy := ulid.Monotonic(ulidEntropyReader, 0)
	id, err := ulid.New(ulid.Timestamp(time.Now().UTC()), entropy)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
