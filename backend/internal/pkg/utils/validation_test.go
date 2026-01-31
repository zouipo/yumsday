package utils

import (
	"strings"
	"testing"
)

// Tests for IsUsernameValid
func TestIsUsernameValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid-single letter", "A", true},
		{"valid-many letters", "Ab", true},
		{"valid-starts with lowercase", "popo", true},
		{"valid-starts with uppercase", "Popo", true},
		{"valid-alphanumeric", "PoPo123", true},
		{"valid-symbols", "A.2_b-", true},
		{"invalid-empty", "", false},
		{"invalid-starts with digit", "1popo", false},
		{"invalid-starts with valid symbol", "_popo", false},
		{"invalid-contains space", "Popo Zoui", false},
		{"invalid-invalid char", "Po!po.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUsernameValid(tt.input)
			if got != tt.want {
				t.Fatalf("IsUsernameValid(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// Tests for IsPasswordValid
func TestIsPasswordValid(t *testing.T) {
	// build a string of a specific length (indifferent character)
	make := func(n int) string { return strings.Repeat("x", n) }

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid-min length (12)", make(12), true},
		{"valid-min length (35)", make(35), true},
		{"valid-max length (72)", make(72), true},
		{"invalid-empty", "", false},
		{"invalid-too short (11)", make(11), false},
		{"invalid-too long (73)", make(73), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPasswordValid(tt.input)
			if got != tt.want {
				t.Fatalf("IsPasswordValid(len=%d) = %v, want %v", len(tt.input), got, tt.want)
			}
		})
	}
}
