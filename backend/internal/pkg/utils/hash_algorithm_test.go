package utils

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{"valid-simple password", "password123", false},
		{"valid-max bcrypt length (72)", strings.Repeat("x", 72), false},
		{"valid-special characters", "P@ssw0rd!#$%", false},
		{"valid-unicode", "пароль密码🔐", false},
		{"valid-whitespace", "pass word test", false},
		{"valid-empty", "", false},
		{"invalid-long password", strings.Repeat("a", 73), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantError {
				if err == nil {
					t.Fatalf("HashPassword(%q) expected error instead of nil", tt.password)
				}
				return
			}

			if err != nil {
				t.Fatalf("HashPassword(%q) unexpected error: %v", tt.password, err)
			}

			if hash == "" {
				t.Fatalf("HashPassword(%q) returned empty hash", tt.password)
			}

			if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password)); err != nil {
				t.Fatalf("HashPassword(%q) produced invalid hash: %v", tt.password, err)
			}

			// Verify wrong password doesn't match.
			// Prepend "wrong" to avoid bcrypt's 72-byte truncation issue:
			// if the password is already 72 character, wrong won't be included in the hash comparison, leading to a false positive match.
			wrongPassword := "wrong" + tt.password
			if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(wrongPassword)); err == nil {
				t.Fatalf("HashPassword(%q) hash incorrectly matched wrong password", tt.password)
			}
		})
	}
}

func TestHashPassword_ProducesDifferentHashes(t *testing.T) {
	password := "testpassword123"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword(%q) unexpected error: %v", password, err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword(%q) unexpected error: %v", password, err)
	}

	if hash1 == hash2 {
		t.Fatalf("HashPassword(%q) produced identical hashes, expected different due to salt", password)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash1), []byte(password)); err != nil {
		t.Fatalf("hash1 failed verification: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash2), []byte(password)); err != nil {
		t.Fatalf("hash2 failed verification: %v", err)
	}
}

func TestHashPassword_ConsistentVerification(t *testing.T) {
	password := "mySecurePassword456"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword(%q) unexpected error: %v", password, err)
	}

	// Verify multiple times to ensure consistency
	for i := 0; i < 10; i++ {
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			t.Fatalf("verification attempt %d failed: %v", i+1, err)
		}
	}
}
