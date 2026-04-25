package service

import (
	"errors"
	"log/slog"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Authenticate(session *model.Session, username, password string) (*model.User, error)
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
func (s *AuthService) Authenticate(session *model.Session, username, password string) (*model.User, error) {
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	slog.Debug("Checking password", "username", username)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, customErrors.NewUnauthorizedError("invalid credentials", err)
		}
		return nil, customErrors.NewInternalError("an error occurred while checking credentials", err)
	}

	session.UserID = user.ID
	err = s.sessionService.Save(session)
	if err != nil {
		return nil, err
	}
	slog.Debug("User authenticated successfully", "username", username)
	return user, nil
}

// Logout removes the session from the session store, effectively logging out the user.
func (s *AuthService) Logout(session *model.Session) error {
	err := s.sessionService.Remove(session)
	if err != nil {
		return err
	}
	slog.Debug("User logged out successfully", "sessionID", session.ID)
	return nil
}
