package service

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

// SessionServiceInterface defines the contract for session operations used by middlewares.
type SessionServiceInterface interface {
	GetSession(r *http.Request) *model.Session
	CookieName() string
	Expiration() time.Duration
	Save(session *model.Session)
	Remove(session *model.Session) error
}

type SessionService struct {
	repo       repository.SessionRepositoryInterface
	cookieName string
	expiration time.Duration
}

func NewSessionService(repo repository.SessionRepositoryInterface, cookieName string, expiration time.Duration) *SessionService {
	s := &SessionService{
		repo:       repo,
		cookieName: cookieName,
		expiration: expiration,
	}
	go s.cleanUp()
	return s
}

func (s *SessionService) CookieName() string {
	return s.cookieName
}

func (s *SessionService) Expiration() time.Duration {
	return s.expiration
}

// GetSession gets the session with the ID from the request's session cookie.
// It returns the session if found, or a new session if not found.
func (s *SessionService) GetSession(r *http.Request) *model.Session {
	cookie, err := r.Cookie(s.cookieName)
	var session *model.Session
	if err == nil {
		session, err = s.repo.GetByID(cookie.Value)
		if err != nil {
			appErr, ok := errors.AsType[*customErrors.AppError](err)
			if !ok || appErr.StatusCode != http.StatusNotFound {
				slog.Error(
					"Failed to read session from repo, generating new session",
					"error", err.Error(),
				)
			}
		}
		if err == nil {
			slog.Debug("Found session", "id", cookie.Value)
		}
	}

	// Check if the session is expired
	if session != nil && time.Since(session.LastActivity) > s.expiration {
		err = s.repo.Delete(session.ID)
		if err != nil {
			panic(err)
		}
		session = nil
		slog.Debug("Session expired, removed from repo")
	}

	if session == nil {
		session = model.NewSession()
		slog.Debug("Generated new session", "id", session.ID)
	}

	return session
}

func (s *SessionService) Save(session *model.Session) {
	session.LastActivity = time.Now().UTC()
	if err := s.repo.Write(session); err != nil {
		slog.Error("Failed to write session to repository", "error", err)
	}
}

func (s *SessionService) Remove(session *model.Session) error {
	slog.Debug("Removing session", "id", session.ID)
	return s.repo.Delete(session.ID)
}

/*** PRIVATE METHODS ***/

// cleanUp periodically removes expired sessions from the database.
func (s *SessionService) cleanUp() {
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

	for range ticker.C {
		impl()
	}
}
