package services

import (
	"errors"
	"log/slog"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
	"github.com/zouipo/yumsday/backend/internal/repositories"
)

// UserServiceInterface defines the contract for user service operations
type UserServiceInterface interface {
	GetAll() ([]models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) (int64, error)
	Update(user *models.User) error
	UpdateAdminRole(userID int64, role bool) error
	UpdatePassword(userID int64, oldPassword string, newPassword string) error
	Delete(id int64) error
}

type UserService struct {
	repo repositories.UserRepositoryInterface
}

// NewUserService creates a new UserService using the provided UserRepository.
func NewUserService(repo repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}

/*** READ OPERATIONS ***/

// GetAll returns all users from the repository or an error if the fetch fails.
func (s *UserService) GetAll() ([]models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID returns the user identified by id or an error if not found.
func (s *UserService) GetByID(id int64) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByUsername returns the user that matches the provided username or an error.
func (s *UserService) GetByUsername(username string) (*models.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

/*** CREATE OPERATIONS ***/

// Create validates and creates a new user, returning the new user ID or an error.
func (s *UserService) Create(user *models.User) (int64, error) {
	user.CreatedAt = time.Now()

	if !utils.IsUsernameValid(user.Username) {
		slog.Debug("Invalid username format", "username", user.Username)
		return 0, customErrors.NewValidationError("username", "Invalid username format", nil)
	}

	if !utils.IsPasswordValid(user.Password) {
		slog.Debug("Invalid password length")
		return 0, customErrors.NewValidationError("password", "Invalid password length", nil)
	}

	var err error
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		slog.Error("CreateUser: failed to hash password", "error", err)
		return 0, customErrors.NewInternalServerError("Failed to hash password", err)
	}

	id, err := s.repo.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

/*** UPDATE OPERATIONS ***/

// Update updates mutable fields (username, avatar, language, theme) of the given user after validation.
func (s *UserService) Update(user *models.User) error {
	currentUser, err := s.GetByID(user.ID)
	if err != nil {
		return err
	}

	// Check if the username is being updated to an already existing one
	if user.Username != currentUser.Username {
		if !utils.IsUsernameValid(user.Username) {
			return customErrors.NewValidationError("username", "Invalid username format", nil)
		}

		currentUser.Username = user.Username
	}

	if user.Avatar != currentUser.Avatar {
		currentUser.Avatar = user.Avatar
	}

	if user.Language != currentUser.Language {
		currentUser.Language = user.Language
	}

	if user.AppTheme != currentUser.AppTheme {
		currentUser.AppTheme = user.AppTheme
	}

	if err = s.repo.Update(currentUser); err != nil {
		return err
	}

	return nil
}

// UpdateAdminRole sets or clears the admin flag for the user with the given ID.
func (s *UserService) UpdateAdminRole(userID int64, role bool) error {
	if err := s.repo.UpdateAdminRole(userID, role); err != nil {
		return err
	}

	return nil
}

// UpdatePassword verifies the old password and updates to the new password after validation.
func (s *UserService) UpdatePassword(userID int64, oldPassword string, newPassword string) error {
	if oldPassword == newPassword {
		slog.Debug("Same old and new passwords")
		return nil
	}

	if oldPassword == "" || newPassword == "" {
		slog.Debug("Old or new password empty")
		return customErrors.NewValidationError("password", "Old and new passwords must be provided", nil)
	}

	if !utils.IsPasswordValid(newPassword) {
		slog.Debug("Invalid password length")
		return customErrors.NewValidationError("password", "Invalid password length", nil)
	}

	currentUser, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	// Comparing current old password with the one provided by the user
	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(oldPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			slog.Debug("Wrong old password")
			return customErrors.NewValidationError("password", "Old password is incorrect for user "+currentUser.Username, err)
		}
		return err
	}

	currentUser.Password, err = utils.HashPassword(newPassword)
	if err != nil {
		slog.Error("UpdatePassword: failed to hash new password", "error", err)
		return customErrors.NewInternalServerError("Failed to hash new password", err)
	}

	if err = s.repo.Update(currentUser); err != nil {
		return err
	}

	return nil
}

/*** DELETE OPERATIONS ***/

// Delete removes the user with the specified ID from the repository.
func (s *UserService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
