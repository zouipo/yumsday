package utils

import (
	"testing"
)

func TestGenerateSessionID(t *testing.T) {
	sessionID := GenerateSessionID()

	if sessionID == "" {
		t.Error("Expected non-empty session ID, got empty string")
	}
}

func TestGenerateSessionIDWithCorrectLength(t *testing.T) {
	sessionID := GenerateSessionID()

	// 32 bytes encoded in base64 RawURL = 43 characters
	expectedLength := 43
	if len(sessionID) != expectedLength {
		t.Errorf("Expected session ID length to be %d, got %d", expectedLength, len(sessionID))
	}
}

func TestGenerateUniqueSessionIDs(t *testing.T) {
	sessionID1 := GenerateSessionID()
	sessionID2 := GenerateSessionID()

	if sessionID1 == sessionID2 {
		t.Error("Expected unique session IDs, got identical values")
	}
}
