package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/zouipo/yumsday/backend/internal/models"
	validation "github.com/zouipo/yumsday/backend/internal/pkg/utils"
	"github.com/zouipo/yumsday/backend/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

// NewUserService creates a new UserService using the provided UserRepository.
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

/*** READ OPERATIONS ***/

// GetAll returns all users from the repository or an error if the fetch fails.
func (s *UserService) GetAll() ([]models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.New("Failed to fetch users: " + err.Error())
	}

	return users, nil
}

// GetByID returns the user identified by id or an error if not found.
func (s *UserService) GetByID(id int64) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, fmt.Errorf("User ID %v doesn't exist: %v", id, err.Error())
	} else if err != nil {
		return nil, fmt.Errorf("Failed to fetch user by ID %v: %v", id, err.Error())
	}

	return user, nil
}

// GetByUsername returns the user that matches the provided username or an error.
func (s *UserService) GetByUsername(username string) (*models.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch user by username %v: %v", username, err.Error())
	}

	return user, nil
}

/*** CREATE OPERATIONS ***/

// Create validates and creates a new user, returning the new user ID or an error.
func (s *UserService) Create(user *models.User) (int64, error) {
	user.CreatedAt = time.Now()

	usernameExists, err := s.usernameExists(user.Username)
	if usernameExists == true && err == nil {
		return 0, fmt.Errorf("Username %v already exists", user.Username)
	} else if err != nil {
		return 0, err
	}
	if !validation.IsUsernameValid(user.Username) {
		return 0, fmt.Errorf("Invalid username format")
	}

	if !validation.IsPasswordValid(user.Password) {
		return 0, fmt.Errorf("Invalid password length")
	}
	// TODO: Hash password before saving user in database

	id, err := s.repo.Create(user)
	if err != nil {
		return 0, fmt.Errorf("Failed to create user %v: %v", user.Username, err.Error())
	}

	return id, nil
}

/*** UPDATE OPERATIONS ***/

// Update updates mutable fields (username, avatar, language, theme) of the given user after validation.
func (s *UserService) Update(user *models.User) error {
	// Fetch the current user data in database
	currentUser, err := s.GetByID(user.ID)
	if err != nil {
		return err
	}

	// Check if the username is being updated to an already existing one
	if user.Username != currentUser.Username {
		usernameExists, err := s.usernameExists(user.Username)
		if err != nil {
			return err
		}
		if usernameExists == true {
			return fmt.Errorf("Username %v already exists", user.Username)
		}
		if !validation.IsUsernameValid(user.Username) {
			return fmt.Errorf("Invalid username format")
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
		return fmt.Errorf("Failed to update user %v: %v", user.Username, err.Error())
	}

	return nil
}

// UpdateAdminRole sets or clears the admin flag for the user with the given ID.
func (s *UserService) UpdateAdminRole(userID int64, role bool) error {
	// Fetch the current user data in database
	currentUser, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	currentUser.AppAdmin = role

	if err = s.repo.Update(currentUser); err != nil {
		return fmt.Errorf("Failed to update user %v: %v", currentUser.Username, err.Error())
	}

	return nil
}

// UpdatePassword verifies the old password and updates to the new password after validation.
func (s *UserService) UpdatePassword(userID int64, oldPassword string, newPassword string) error {
	// Fetch the current user data in database
	currentUser, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	if currentUser.Password != oldPassword { // TODO: Hash the old password before comparing
		return fmt.Errorf("Old password is incorrect for user %v", currentUser.Username)
	}

	if !validation.IsPasswordValid(newPassword) {
		return fmt.Errorf("Invalid password length")
	}

	currentUser.Password = newPassword // TODO: Hash the new password before saving

	if err = s.repo.Update(currentUser); err != nil {
		return fmt.Errorf("Failed to update password for user %v: %v", currentUser.Username, err.Error())
	}

	return nil
}

/*** DELETE OPERATIONS ***/

// Delete removes the user with the specified ID from the repository.
func (s *UserService) Delete(id int64) error {
	userExists, err := s.userExists(id)
	if userExists == false && err == nil {
		return fmt.Errorf("User ID %v doesn't exist", id)
	} else if err != nil {
		return err
	}

	if err = s.repo.Delete(id); err != nil {
		return fmt.Errorf("Failed to delete user %v: %v", id, err.Error())
	}

	return nil
}

/*** HELPER FUNCTIONS ***/

// userExists checks whether a user with the given id exists in the repository.
func (s *UserService) userExists(id int64) (bool, error) {
	user, err := s.repo.GetByID(id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("Failed to check if user %v exists: %v", id, err.Error())
	}

	return user != nil, nil
}

// usernameExists checks whether the provided username is already taken.
func (s *UserService) usernameExists(username string) (bool, error) {
	userByUsername, err := s.repo.GetByUsername(username)
	// TODO: implement a custom error "already existing username"
	// and check if the function correctly returns this specific error
	if err == nil && userByUsername != nil {
		return true, nil
	} else if err != nil && err.Error() != "sql: no rows in result set" {
		return false, fmt.Errorf("Failed to check existing username %v: %v", username, err.Error())
	}

	return false, nil
}
