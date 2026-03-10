package model

import (
	"time"

	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

type Session struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	UserID       int64     `json:"user_id"`
}

// NewSession creates a new session with a unique session ID and sets the creation and expiration times.
func NewSession() *Session {
	return &Session{
		ID:           utils.GenerateSessionID(),
		CreatedAt:    time.Now().UTC(),
		LastActivity: time.Now().UTC(),
	}
}
