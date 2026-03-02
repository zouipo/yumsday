package service

import (
	"log/slog"
	"time"

	"github.com/zouipo/yumsday/backend/internal/repository"
)

type Service struct {
	repo       *repository.SessionRepository
	cookieName string
	expiration time.Duration
}

func NewSessionService(repo *repository.SessionRepository, cookieName string, expiration time.Duration) *Service {
	s := &Service{
		repo:       repo,
		cookieName: cookieName,
		expiration: expiration,
	}
	go s.cleanUp(expiration)
	return s
}

func (s *Service) CookieName() string {
	return s.cookieName
}

func (s *Service) Expiration() time.Duration {
	return s.expiration
}

/*** PRIVATE METHODS ***/

// cleanUp periodically removes expired sessions from the database.
func (s *Service) cleanUp(interval time.Duration) {
	impl := func() {
		slog.Debug("Cleaning up expired sessions")
		removed := s.repo.CleanUp(s.expiration)
		if removed > 0 {
			slog.Info("Removed expired sessions", "removed", removed)
		}
	}

	impl()
	// Run the cleanup every hour
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		impl()
	}
}
