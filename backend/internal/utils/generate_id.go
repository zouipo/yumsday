package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func GenerateSessionID() string {
	id := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		panic("Failed to generate session ID: " + err.Error())
	}

	return base64.RawURLEncoding.EncodeToString(id)
}
