package utils

import (
	"encoding/base64"
	"testing"
)

func TestGenerateSessionID(t *testing.T) {
	t.Run("should generate non-empty session ID", func(t *testing.T) {
		sessionID := GenerateSessionID()

		if sessionID == "" {
			t.Error("Expected non-empty session ID, got empty string")
		}
	})

	t.Run("should generate session ID with correct length", func(t *testing.T) {
		sessionID := GenerateSessionID()

		// 32 bytes encoded in base64 RawURL = 43 characters
		expectedLength := 43
		if len(sessionID) != expectedLength {
			t.Errorf("Expected session ID length to be %d, got %d", expectedLength, len(sessionID))
		}
	})

	t.Run("should generate unique session IDs", func(t *testing.T) {
		sessionID1 := GenerateSessionID()
		sessionID2 := GenerateSessionID()

		if sessionID1 == sessionID2 {
			t.Error("Expected unique session IDs, got identical values")
		}
	})

	t.Run("should generate valid base64 RawURL encoded string", func(t *testing.T) {
		sessionID := GenerateSessionID()

		decoded, err := base64.RawURLEncoding.DecodeString(sessionID)
		if err != nil {
			t.Errorf("Expected valid base64 RawURL encoded string, got error: %v", err)
		}

		if len(decoded) != 32 {
			t.Errorf("Expected decoded session ID to be 32 bytes, got %d", len(decoded))
		}
	})
}
