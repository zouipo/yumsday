package service

import (
	"errors"
	"log/slog"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Authenticate(session *model.Session, username, password string) error
	Logout(session *model.Session) error
}

type AuthService struct {
	sessionService SessionServiceInterface
	userService    UserServiceInterface
}

func NewAuthService(sessionService SessionServiceInterface, userService UserServiceInterface) *AuthService {
	return &AuthService{
		sessionService: sessionService,
		userService:    userService,
	}
}

// Checks if the password is valid for this username.
// Assigns the user carrying this username to the session.
func (s *AuthService) Authenticate(session *model.Session, username, password string) error {
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		return err
	}

	slog.Debug("Checking password", "username", username)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			slog.Debug("Wrong credentials", "username", username)
			return customErrors.NewUnauthorizedError(err, "invalid credentials")
		}
		return customErrors.NewInternalServerError("an error occurred while checking credentials", err)
	}

	session.UserID = user.ID
	s.sessionService.Save(session)
	slog.Debug("User authenticated successfully", "username", username)
	return nil
}

func (s *AuthService) Logout(session *model.Session) error {
	err := s.sessionService.Remove(session)
	if err != nil {
		slog.Error("Failed to remove session", "error", err)
		return customErrors.NewInternalServerError("An error occurred while logging out", err)
	}
	slog.Debug("User logged out successfully", "sessionID", session.ID)
	return nil
}
