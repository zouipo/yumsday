package services

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/mattn/go-sqlite3"
	customErrors "github.com/zouipo/yumsday/backend/internal/errors"

	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/models/enums"
)

// Variables for test data
var (
	testUser1 = createTestUser(1, "user1", "password123")
	testUser2 = createTestUser(2, "user2", "password456")
	testUser3 = createTestUser(3, "user3", "password789")

	validUsername = "validuser"
	validPassword = "ValidPass123"

	invalidId       = -1
	invalidUsername = "_"
	invalidPassword = "tooshort"

	notFoundIdErr = customErrors.NewEntityNotFoundError("User", strconv.FormatInt(int64(invalidId), 10), nil)
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users        []models.User
	nextID       int64
	getAllErr    error
	getByIDErr   error
	getByNameErr error
	createErr    error
	updateErr    error
	deleteErr    error
}

// NewMockUserRepository creates a new mock repository with some test data
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make([]models.User, 0),
		nextID: 1,
	}
}

/*** USERREPOSITORY IMPLEMENTATION ***/

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}

	return m.users, nil
}

func (m *MockUserRepository) GetByID(id int64) (*models.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.users {
		if m.users[i].ID == id {
			return &m.users[i], nil
		}
	}
	return nil, customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), nil)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}

	for i := range m.users {
		if m.users[i].Username == username {
			return &m.users[i], nil
		}
	}
	return nil, customErrors.NewEntityNotFoundError("User", username, nil)
}

func (m *MockUserRepository) Create(user *models.User) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}

	for i := range m.users {
		if m.users[i].Username == user.Username {
			return 0, customErrors.NewConflictError("User", "already exists", nil)
		}
	}

	id := m.nextID
	m.nextID++

	user.ID = id
	m.users = append(m.users, *user)

	return id, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}

	// Check if another user already has this username
	for _, existingUser := range m.users {
		if existingUser.Username == user.Username && existingUser.ID != user.ID {
			return customErrors.NewConflictError("User", "already exists", sqlite3.ErrConstraintUnique)
		}
	}

	// Find and update the user
	for i, existingUser := range m.users {
		if existingUser.ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(user.ID, 10), nil)
}

func (m *MockUserRepository) UpdateAdminRole(userID int64, role bool) error {
	if m.updateErr != nil {
		return m.updateErr
	}

	for i, existingUser := range m.users {
		if existingUser.ID == userID {
			m.users[i].AppAdmin = role
			return nil
		}
	}
	return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(userID, 10), nil)
}

func (m *MockUserRepository) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}

	for i, user := range m.users {
		if user.ID == id {
			// Remove user from slice by recreating the slice with users
			// before and after the index to delete
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), nil)
}

/*** HELPER FUNCTIONS ***/

// Helper methods for setting up test scenarios
func (m *MockUserRepository) addUser(user *models.User) {
	user.ID = m.nextID
	m.nextID++
	m.users = append(m.users, *user)
}

// Helper function to create a test user
func createTestUser(id int64, username string, password string) *models.User {
	avatar := enums.Avatar1
	return &models.User{
		ID:        id,
		Username:  username,
		Password:  password,
		AppAdmin:  false,
		CreatedAt: time.Now(),
		Avatar:    &avatar,
		Language:  enums.English,
		AppTheme:  enums.Light,
	}
}

// setupTestData creates a fresh mock repository with predefined test users for test independence
func setupTestData() *MockUserRepository {
	mockRepo := NewMockUserRepository()
	mockRepo.addUser(testUser1)
	mockRepo.addUser(testUser2)
	mockRepo.addUser(testUser3)
	return mockRepo
}

func copyUser(user models.User) models.User {
	copy := user
	if user.Avatar != nil {
		avatarCopy := *user.Avatar
		copy.Avatar = &avatarCopy
	}
	if user.LastVisitedGroup != nil {
		groupCopy := *user.LastVisitedGroup
		copy.LastVisitedGroup = &groupCopy
	}
	return copy
}

// compareUsers compares two User objects and returns an error if any field does not match
func compareUsers(actual, expected *models.User) error {
	if actual.ID != expected.ID {
		return fmt.Errorf("ID = %d instead of %d", actual.ID, expected.ID)
	}
	if actual.Username != expected.Username {
		return fmt.Errorf("username = %s instead of %s", actual.Username, expected.Username)
	}
	if actual.Password != expected.Password {
		return fmt.Errorf("password = %s instead of %s", actual.Password, expected.Password)
	}
	if actual.AppAdmin != expected.AppAdmin {
		return fmt.Errorf("appAdmin ='%v'instead of %v", actual.AppAdmin, expected.AppAdmin)
	}
	// Verify both dates are within the last 2 minutes from now
	now := time.Now()
	threshold := now.Add(-2 * time.Minute)
	if actual.CreatedAt.Before(threshold) || actual.CreatedAt.After(now) {
		return fmt.Errorf("createdAt ='%v'is not within the last 2 minutes (threshold: %v, now: %v)", actual.CreatedAt, threshold, now)
	}
	if expected.CreatedAt.Before(threshold) || expected.CreatedAt.After(now) {
		return fmt.Errorf("expected createdAt ='%v'is not within the last 2 minutes (threshold: %v, now: %v)", expected.CreatedAt, threshold, now)
	}
	// Check Avatar with nil handling
	if (actual.Avatar == nil) != (expected.Avatar == nil) {
		return fmt.Errorf("avatar ='%v'instead of %v", actual.Avatar, expected.Avatar)
	}
	if actual.Avatar != nil && expected.Avatar != nil && *actual.Avatar != *expected.Avatar {
		return fmt.Errorf("avatar ='%v'instead of %v", *actual.Avatar, *expected.Avatar)
	}
	if actual.Language != expected.Language {
		return fmt.Errorf("language ='%v'instead of %v", actual.Language, expected.Language)
	}
	if actual.AppTheme != expected.AppTheme {
		return fmt.Errorf("appTheme ='%v'instead of %v", actual.AppTheme, expected.AppTheme)
	}
	return nil
}

func compareErrors(actual, expected error) bool {
	if actual == nil && expected == nil {
		return true
	}

	if (actual == nil) != (expected == nil) {
		return false
	}

	// Cast both to *AppError
	actualAppErr, actualIsAppErr := actual.(*customErrors.AppError)
	expectedAppErr, expectedIsAppErr := expected.(*customErrors.AppError)

	if actualIsAppErr && expectedIsAppErr {
		if actualAppErr.Message != expectedAppErr.Message && actualAppErr.StatusCode != expectedAppErr.StatusCode && actualAppErr.Err.Error() != expectedAppErr.Err.Error() {
			return false
		}
		return true
	}

	// If not AppErrors, compare their error messages
	return actual.Error() == expected.Error()
}

/*** TEST CONSTRUCTOR ***/

func TestNewUserService(t *testing.T) {
	mockRepo := NewMockUserRepository()

	service := NewUserService(mockRepo)

	if service == nil {
		t.Fatal("NewUserService() returned nil")
	}

	if service.repo == nil {
		t.Error("NewUserService() repo is nil")
	}
}

/*** READ OPERATIONS TESTS ***/

func TestGetAll_Success(t *testing.T) {
	// Arrange
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	// Act
	users, err := service.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error ='%v'instead of nil", err)
	}

	// Assert
	if len(users) != len(mockRepo.users) {
		t.Errorf("GetAll() returned %d users instead of %d", len(users), len(mockRepo.users))
	}
}

func TestGetAll_RepositoryError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.getAllErr = customErrors.NewInternalServerError("Failed to fetch users", nil)

	service := &UserService{repo: mockRepo}

	users, err := service.GetAll()
	if users != nil {
		t.Error("GetAll() expected error instead of non-nil users")
	}

	if !compareErrors(err, mockRepo.getAllErr) {
		t.Errorf("GetAll() error ='%v'instead of expected %v", err, mockRepo.getAllErr)
	}
}

func TestGetByID_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	expectedUser := mockRepo.users[0]

	user, err := service.GetByID(expectedUser.ID)
	if err != nil {
		t.Fatalf("GetByID() error ='%v'instead of nil", err)
	}

	if err := compareUsers(user, &expectedUser); err != nil {
		t.Error("GetByID() returned user with mismatched fields: " + err.Error())
	}
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user, err := service.GetByID(int64(invalidId))
	if user != nil {
		t.Error("GetByID() expected error instead of non-nil user")
	}

	if !compareErrors(err, notFoundIdErr) {
		t.Errorf("GetByID() error ='%v'instead of expected %v", err, notFoundIdErr)
	}
}

func TestGetByID_RepositoryError(t *testing.T) {
	mockRepo := setupTestData()
	mockRepo.getByIDErr = customErrors.NewInternalServerError("Failed to fetch user by ID", nil)

	service := &UserService{repo: mockRepo}

	expectedUser := mockRepo.users[0]

	user, err := service.GetByID(expectedUser.ID)
	if user != nil {
		t.Error("GetByID() expected error instead of non-nil user")
	}

	if !compareErrors(err, mockRepo.getByIDErr) {
		t.Errorf("GetByID() error ='%v'instead of expected %v", err, mockRepo.getByIDErr)
	}
}

func TestGetByUsername_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	expectedUser := mockRepo.users[1]

	user, err := service.GetByUsername(expectedUser.Username)
	if err != nil {
		t.Fatalf("GetByUsername() error ='%v'instead of nil", err)
	}

	if err := compareUsers(user, &expectedUser); err != nil {
		t.Error("GetByUsername() returned user with mismatched fields: " + err.Error())
	}
}

func TestGetByUsername_NotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	// A améliorer en renvoyant une erreur NotFound custom et en vérifiant que c'est bien cette erreur qui est renvoyée
	user, err := service.GetByUsername(invalidUsername)
	if user != nil {
		t.Error("GetByUsername() expected error instead of non-nil user")
	}

	notFoundUsernameErr := customErrors.NewEntityNotFoundError("User", invalidUsername, nil)
	if !compareErrors(err, notFoundUsernameErr) {
		t.Errorf("GetByUsername() error ='%v'instead of expected %v", err, notFoundUsernameErr)
	}
}

/*** CREATE OPERATIONS TESTS ***/

func TestCreate_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	usersNb := len(mockRepo.users)

	newUser := createTestUser(0, validUsername, validPassword)

	id, err := service.Create(newUser)
	if err != nil {
		t.Fatalf("Create() error ='%v'instead of nil", err)
	}

	if id == 0 {
		t.Error("Create() returned ID 0, expected non-zero ID")
	}

	// Verify the user was added to the repository
	createdUser, err := mockRepo.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID() error ='%v'instead of nil", err)
	}

	if len(mockRepo.users) != usersNb+1 {
		t.Errorf("Create() did not add user to repository, expected %d users instead of %d", usersNb+1, len(mockRepo.users))
	}

	err = compareUsers(createdUser, &models.User{
		ID:        id,
		Username:  newUser.Username,
		Password:  newUser.Password,
		AppAdmin:  newUser.AppAdmin,
		CreatedAt: newUser.CreatedAt,
		Avatar:    newUser.Avatar,
		Language:  newUser.Language,
		AppTheme:  newUser.AppTheme,
	})

	if err != nil {
		t.Error("Create() returned user with mismatched fields: " + err.Error())
	}
}

func TestCreate_DuplicateUsername(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	existingUser := mockRepo.users[0]

	newUser := createTestUser(0, existingUser.Username, validPassword)

	id, err := service.Create(newUser)

	if id != 0 {
		t.Errorf("Create() expected ID 0 for duplicate username, got %d", id)
	}

	errConflict := customErrors.NewConflictError("User", "already exists", nil)
	if !compareErrors(err, errConflict) {
		t.Error("Create() expected error for duplicate username, got nil")
	}
}

func TestCreate_InvalidUsername(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	usersNb := len(mockRepo.users)

	newUser := createTestUser(0, invalidUsername, validPassword)

	id, err := service.Create(newUser)
	if id != 0 {
		t.Errorf("Create() expected ID 0 for invalid username, got %d", id)
	}

	validationErr := customErrors.NewValidationError("username", "Invalid username format", nil)
	if !compareErrors(err, validationErr) {
		t.Errorf("Create() error ='%v'instead of expected %v", err, validationErr)
	}

	if len(mockRepo.users) != usersNb {
		t.Errorf("Create() should not add user to repository on invalid username, expected %d users instead of %d", usersNb, len(mockRepo.users))
	}
}

func TestCreate_InvalidPassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	usersNb := len(mockRepo.users)

	newUser := createTestUser(0, validUsername, invalidPassword)

	id, err := service.Create(newUser)
	if id != 0 {
		t.Error("Create() expected error for invalid password instead of nil")
	}

	validationErr := customErrors.NewValidationError("password", "Old and new passwords must be provided", nil)
	if !compareErrors(err, validationErr) {
		t.Errorf("Create() error ='%v'instead of expected %v", err, validationErr)
	}

	if len(mockRepo.users) != usersNb {
		t.Errorf("Create() should not add user to repository on invalid password, expected %d users instead of %d", usersNb, len(mockRepo.users))
	}
}

func TestCreate_RepositoryError(t *testing.T) {
	mockRepo := setupTestData()
	mockRepo.createErr = customErrors.NewInternalServerError("Failed to retrieve created user", nil)

	usersNb := len(mockRepo.users)

	service := &UserService{repo: mockRepo}

	newUser := createTestUser(0, validUsername, validPassword)

	id, err := service.Create(newUser)
	if id != 0 {
		t.Errorf("Create() expected ID 0 for database error, got %d", id)
	}

	if !compareErrors(err, mockRepo.createErr) {
		t.Errorf("Create() error ='%v'instead of expected %v", err, mockRepo.createErr)
	}

	if len(mockRepo.users) != usersNb {
		t.Errorf("Create() should not add user to repository on error, expected %d users instead of %d", usersNb, len(mockRepo.users))
	}
}

/*** UPDATE OPERATIONS TESTS ***/

func TestUpdate_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	existingUser := copyUser(mockRepo.users[0])
	existingUser.Username = validUsername
	existingUser.Avatar = nil
	existingUser.Language = enums.French
	existingUser.AppTheme = enums.Dark

	err := service.Update(&existingUser)
	if err != nil {
		t.Fatalf("Update() error = '%v' instead of nil", err)
	}

	// Verify the username was updated
	user, _ := mockRepo.GetByID(existingUser.ID)
	if user == nil {
		t.Error("Update() failed to update username")
	}

	err = compareUsers(user, &existingUser)

	if err != nil {
		t.Error("Update() returned user with mismatched fields: " + err.Error())
	}
}

func TestUpdate_UserNotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	updatedUser := createTestUser(int64(invalidId), validUsername, validPassword)

	err := service.Update(updatedUser)
	if !compareErrors(err, notFoundIdErr) {
		t.Error("Update() expected error for non-existent user, got nil")
	}
}

func TestUpdate_DuplicateUsername(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	firstUser := mockRepo.users[0]
	secondUser := mockRepo.users[1]

	updatedUser := createTestUser(firstUser.ID, secondUser.Username, validPassword)

	conflictErr := customErrors.NewConflictError("User", "already exists", sqlite3.ErrConstraintUnique)
	err := service.Update(updatedUser)
	if !compareErrors(err, conflictErr) {
		t.Errorf("Update() expected error '%v' for duplicate username instead of '%v'", conflictErr, err)
	}
}

func TestUpdate_InvalidUsername(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	updatedUser := createTestUser(1, invalidUsername, validPassword)

	err := service.Update(updatedUser)
	validationErr := customErrors.NewValidationError("username", "Invalid username format", nil)
	if !compareErrors(err, validationErr) {
		t.Errorf("Update() expected error '%v' for invalid username instead of %v", validationErr, err)
	}

	if user, _ := mockRepo.GetByID(1); user.Username == invalidUsername {
		t.Error("Update() should not update username to invalid value")
	}
}

func TestUpdateAdminRole_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	err := service.UpdateAdminRole(1, true)
	if err != nil {
		t.Fatalf("UpdateAdminRole() error ='%v'instead of nil", err)
	}

	user, _ := mockRepo.GetByID(1)
	if !user.AppAdmin {
		t.Error("UpdateAdminRole() failed to set admin role")
	}
}

func TestUpdateAdminRole_UserNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := &UserService{repo: mockRepo}

	err := service.UpdateAdminRole(int64(invalidId), true)

	notFoundIdErr = customErrors.NewEntityNotFoundError("User", strconv.FormatInt(int64(invalidId), 10), nil)
	if !compareErrors(err, notFoundIdErr) {
		t.Errorf("UpdateAdminRole() expected error '%v' for non-existent user instead of '%v'", notFoundIdErr, err)
	}
}

func TestUpdateAdminRole_RepositoryError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := &UserService{repo: mockRepo}
	mockRepo.updateErr = customErrors.NewInternalServerError("Failed to update user admin role", nil)

	err := service.UpdateAdminRole(1, true)
	if !compareErrors(err, mockRepo.updateErr) {
		t.Errorf("UpdateAdminRole() expected error '%v' for repository error instead of '%v'", mockRepo.updateErr, err)
	}
}

func TestUpdatePassword_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, user.Password, validPassword)
	if err != nil {
		t.Fatalf("UpdatePassword() error ='%v'instead of nil", err)
	}

	// Verify password was updated
	updatedUser, _ := mockRepo.GetByID(user.ID)
	if updatedUser.Password != validPassword {
		t.Error("UpdatePassword() failed to update password")
	}
}

func TestUpdatePassword_EmptyPasswords(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, "", validPassword)

	validationErr := customErrors.NewValidationError("password", "Old and new passwords must be provided", nil)
	if !compareErrors(err, validationErr) {
		t.Errorf("UpdatePassword() expected error '%v' for empty old password instead of '%v'", validationErr, err)
	}

	err = service.UpdatePassword(user.ID, user.Password, "")
	if !compareErrors(err, validationErr) {
		t.Errorf("UpdatePassword() expected error '%v' for empty new password instead of '%v'", validationErr, err)
	}

	if updatedUser, _ := mockRepo.GetByID(user.ID); updatedUser.Password != user.Password {
		t.Error("UpdatePassword() should not update password when old or new password is empty")
	}
}

func TestUpdatePassword_IncorrectOldPassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, validPassword, validPassword+"123")
	validationErr := customErrors.NewValidationError("password", "Old password is incorrect for user "+user.Password, nil)
	if !compareErrors(err, validationErr) {
		t.Errorf("UpdatePassword() expected error '%v' for incorrect old password instead of '%v'", validationErr, err)
	}

	if updatedUser, _ := mockRepo.GetByID(user.ID); updatedUser.Password != user.Password {
		t.Error("UpdatePassword() should not update password when old password is incorrect")
	}
}

func TestUpdatePassword_SamePassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, user.Password, user.Password)
	if err != nil {
		t.Errorf("UpdatePassword() error ='%v'instead of nil for same password", err)
	}

	if updatedUser, _ := mockRepo.GetByID(user.ID); updatedUser.Password != user.Password {
		t.Error("UpdatePassword() should not change password when new password is the same as old password")
	}
}

func TestUpdatePassword_InvalidNewPassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, user.Password, invalidPassword)

	validationErr := customErrors.NewValidationError("password", "Invalid password length", nil)
	if !compareErrors(err, validationErr) {
		t.Errorf("UpdatePassword() expected error '%v' for invalid new password instead of '%v'", validationErr, err)
	}

	if updatedUser, _ := mockRepo.GetByID(user.ID); updatedUser.Password != user.Password {
		t.Error("UpdatePassword() should not update password when new password is invalid")
	}
}

func TestUpdatePassword_UserNotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(int64(invalidId), user.Password, validPassword)

	notFoundIdErr = customErrors.NewEntityNotFoundError("User", strconv.FormatInt(int64(invalidId), 10), nil)
	if !compareErrors(err, notFoundIdErr) {
		t.Errorf("UpdatePassword() expected error '%v' for non-existent user instead of '%v'", notFoundIdErr, err)
	}
}

func TestUpdatePassword_RepositoryError(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}
	mockRepo.updateErr = customErrors.NewInternalServerError("Failed to update user", nil)

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, user.Password, validPassword)
	if !compareErrors(err, mockRepo.updateErr) {
		t.Errorf("UpdatePassword() expected error '%v' for repository error instead of '%v'", mockRepo.updateErr, err)
	}
}

/*** DELETE OPERATIONS TESTS ***/

func TestDelete_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	id := mockRepo.users[0].ID

	err := service.Delete(id)
	if err != nil {
		t.Fatalf("Delete() error ='%v'instead of nil", err)
	}

	// Verify user was deleted
	_, err = mockRepo.GetByID(id)
	if err == nil {
		t.Error("Delete() failed to delete user")
	}
}

func TestDelete_UserNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := &UserService{repo: mockRepo}

	err := service.Delete(int64(invalidId))

	notFoundIdErr = customErrors.NewEntityNotFoundError("User", strconv.FormatInt(int64(invalidId), 10), nil)
	if !compareErrors(err, notFoundIdErr) {
		t.Errorf("Delete() expected error '%v' for non-existent user instead of '%v'", notFoundIdErr, err)
	}
}

func TestDelete_RepositoryError(t *testing.T) {
	mockRepo := setupTestData()
	mockRepo.deleteErr = customErrors.NewInternalServerError("Failed to delete user", nil)

	service := &UserService{repo: mockRepo}

	err := service.Delete(1)
	if !compareErrors(err, mockRepo.deleteErr) {
		t.Errorf("Delete() expected error '%v' for repository error instead of '%v'", mockRepo.deleteErr, err)
	}

	// Verify user was not deleted
	_, err = mockRepo.GetByID(1)
	if err != nil {
		t.Error("User should not be deleted when repository returns an error")
	}
}
