package core

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "MySecureP@ssw0rd"
	pepper := "test-pepper"
	cost := 10 // Use lower cost for faster tests

	hash1, err := HashPassword(password, pepper, cost)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash1 == "" {
		t.Error("HashPassword returned empty string")
	}

	// Same password should produce different hashes (due to salt)
	hash2, err := HashPassword(password, pepper, cost)
	if err != nil {
		t.Fatalf("HashPassword second call failed: %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword should produce different hashes for the same password")
	}

	// Different pepper should produce different hash
	hash3, err := HashPassword(password, "different-pepper", cost)
	if err != nil {
		t.Fatalf("HashPassword with different pepper failed: %v", err)
	}

	if hash1 == hash3 {
		t.Error("HashPassword should produce different hashes for different peppers")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "MySecureP@ssw0rd"
	pepper := "test-pepper"
	cost := 10

	hash, err := HashPassword(password, pepper, cost)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Correct password should verify
	valid, err := VerifyPassword(password, pepper, hash)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}
	if !valid {
		t.Error("VerifyPassword should return true for correct password")
	}

	// Wrong password should fail
	valid, err = VerifyPassword("WrongPassword123!", pepper, hash)
	if err != nil {
		t.Fatalf("VerifyPassword with wrong password failed: %v", err)
	}
	if valid {
		t.Error("VerifyPassword should return false for wrong password")
	}

	// Wrong pepper should fail
	valid, err = VerifyPassword(password, "wrong-pepper", hash)
	if err != nil {
		t.Fatalf("VerifyPassword with wrong pepper failed: %v", err)
	}
	if valid {
		t.Error("VerifyPassword should return false for wrong pepper")
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"Valid password", "MySecureP@ssw0rd", true},
		{"Too short", "Short1!", false},
		{"No uppercase", "mysecurep@ssw0rd", false},
		{"No lowercase", "MYSECUREP@SSW0RD", false},
		{"No digit", "MySecureP@ssword", false},
		{"No special", "MySecurePassw0rd", false},
		{"Common password", "Password123!", false},
		{"Sequential", "Abcdefgh1!", false},
		{"Minimum valid", "A1b2C3d4E5f!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidatePassword(tt.password)
			isValid := len(errors) == 0

			if isValid != tt.valid {
				t.Errorf("ValidatePassword(%q) valid = %v, want %v, errors: %v", tt.password, isValid, tt.valid, errors)
			}
		})
	}
}

func TestIsPasswordValid(t *testing.T) {
	if !IsPasswordValid("MySecureP@ssw0rd") {
		t.Error("IsPasswordValid should return true for valid password")
	}

	if IsPasswordValid("short") {
		t.Error("IsPasswordValid should return false for invalid password")
	}
}

func TestHashPasswordEmptyPepper(t *testing.T) {
	password := "MySecureP@ssw0rd"
	cost := 10

	hash, err := HashPassword(password, "", cost)
	if err != nil {
		t.Fatalf("HashPassword with empty pepper failed: %v", err)
	}

	valid, err := VerifyPassword(password, "", hash)
	if err != nil {
		t.Fatalf("VerifyPassword with empty pepper failed: %v", err)
	}
	if !valid {
		t.Error("VerifyPassword should work with empty pepper")
	}
}
