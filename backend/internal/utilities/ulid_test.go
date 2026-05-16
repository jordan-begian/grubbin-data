package utilities

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateULID_Format(t *testing.T) {
	id, err := GenerateULID()
	if err != nil {
		t.Fatalf("GenerateULID failed: %v", err)
	}

	// ULIDs are exactly 26 characters of Crockford's base32
	if len(id) != 26 {
		t.Errorf("Expected ULID length 26, got %d", len(id))
	}

	// Should only contain valid base32 characters
	for _, r := range id {
		if !strings.ContainsRune("0123456789ABCDEFGHJKMNPQRSTVWXYZ", r) {
			t.Errorf("Invalid character in ULID: %q", r)
		}
	}
}

func TestGenerateULID_Uniqueness(t *testing.T) {
	const count = 1000
	seen := make(map[string]struct{}, count)

	for i := range count {
		id, err := GenerateULID()
		if err != nil {
			t.Fatalf("GenerateULID failed on iteration %d: %v", i, err)
		}
		if _, exists := seen[id]; exists {
			t.Fatalf("Duplicate ULID generated: %s", id)
		}
		seen[id] = struct{}{}
	}
}

func TestGenerateULID_Sortable(t *testing.T) {
	id1, err := GenerateULID()
	if err != nil {
		t.Fatalf("GenerateULID failed: %v", err)
	}

	// Sleep to ensure a different timestamp (ULIDs are sortable by time)
	time.Sleep(5 * time.Millisecond)

	id2, err := GenerateULID()
	if err != nil {
		t.Fatalf("GenerateULID failed: %v", err)
	}

	// Later ULID should be lexicographically greater
	if id2 <= id1 {
		t.Errorf("Expected later ULID %s to be greater than earlier %s", id2, id1)
	}
}
