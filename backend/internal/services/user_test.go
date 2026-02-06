package services

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

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
		nextID: 4,
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
	return nil, sql.ErrNoRows
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
	return nil, sql.ErrNoRows
}

func (m *MockUserRepository) Create(user *models.User) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
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

	for i, existingUser := range m.users {
		if existingUser.ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return sql.ErrNoRows
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
	return sql.ErrNoRows
}

/*** HELPER FUNCTIONS ***/

// Helper methods for setting up test scenarios
func (m *MockUserRepository) addUser(user *models.User) {
	m.users = append(m.users, *user)
	if user.ID >= m.nextID {
		m.nextID = user.ID + 1
	}
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

// compareUsers compares two User objects and returns an error if any field does not match
func compareUsers(actual *models.User, expected *models.User) error {
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
		return fmt.Errorf("appAdmin = %v instead of %v", actual.AppAdmin, expected.AppAdmin)
	}
	// Compare only the date (year, month, day), ignoring time
	actualDate := actual.CreatedAt.Truncate(24 * time.Hour)
	expectedDate := expected.CreatedAt.Truncate(24 * time.Hour)
	if actualDate != expectedDate {
		return fmt.Errorf("createdAt date = %v instead of %v", actualDate, expectedDate)
	}
	if *actual.Avatar != *expected.Avatar {
		return fmt.Errorf("avatar = %v instead of %v", *actual.Avatar, *expected.Avatar)
	}
	if actual.Language != expected.Language {
		return fmt.Errorf("language = %v instead of %v", actual.Language, expected.Language)
	}
	if actual.AppTheme != expected.AppTheme {
		return fmt.Errorf("appTheme = %v instead of %v", actual.AppTheme, expected.AppTheme)
	}
	return nil
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
		t.Fatalf("GetAll() error = %v instead of nil", err)
	}

	// Assert
	if len(users) != len(mockRepo.users) {
		t.Errorf("GetAll() returned %d users instead of %d", len(users), len(mockRepo.users))
	}
}

func TestGetAll_RepositoryError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.getAllErr = errors.New("database error")

	service := &UserService{repo: mockRepo}

	_, err := service.GetAll()
	if err == nil {
		t.Error("GetAll() expected error, got nil")
	}
}

func TestGetByID_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	expectedUser := mockRepo.users[0]

	user, err := service.GetByID(expectedUser.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v instead of nil", err)
	}

	err = compareUsers(user, &expectedUser)
	if err != nil {
		t.Error("GetByID() returned user with mismatched fields: " + err.Error())
	}
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	_, err := service.GetByID(int64(invalidId))
	// A améliorer en renvoyant une erreur NotFound custom et en vérifiant que c'est bien cette erreur qui est renvoyée
	if err == nil {
		t.Error("GetByID() expected error for non-existent user, got nil")
	}
}

func TestGetByID_RepositoryError(t *testing.T) {
	mockRepo := setupTestData()
	mockRepo.getByIDErr = errors.New("database error")

	service := &UserService{repo: mockRepo}

	expectedUser := mockRepo.users[0]

	_, err := service.GetByID(expectedUser.ID)
	if err == nil {
		t.Error("GetByID() expected error, got nil")
	}
}

func TestGetByUsername_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	expectedUser := mockRepo.users[1]

	user, err := service.GetByUsername(expectedUser.Username)
	if err != nil {
		t.Fatalf("GetByUsername() error = %v instead of nil", err)
	}

	err = compareUsers(user, &expectedUser)
	if err != nil {
		t.Error("GetByUsername() returned user with mismatched fields: " + err.Error())
	}
}

func TestGetByUsername_NotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	// A améliorer en renvoyant une erreur NotFound custom et en vérifiant que c'est bien cette erreur qui est renvoyée
	_, err := service.GetByUsername(invalidUsername)
	if err == nil {
		t.Error("GetByUsername() expected error for non-existent user, got nil")
	}
}

/*** CREATE OPERATIONS TESTS ***/

func TestCreate_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	newUser := createTestUser(0, validUsername, validPassword)

	id, err := service.Create(newUser)
	if err != nil {
		t.Fatalf("Create() error = %v instead of nil", err)
	}

	if id == 0 {
		t.Error("Create() returned ID 0, expected non-zero ID")
	}

	// Verify the user was added to the repository
	createdUser, err := mockRepo.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID() error = %v instead of nil", err)
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

	_, err := service.Create(newUser)
	if err == nil {
		t.Error("Create() expected error for duplicate username, got nil")
	}
}

func TestCreate_InvalidUsername(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	newUser := createTestUser(0, invalidUsername, validPassword)

	_, err := service.Create(newUser)
	if err == nil {
		t.Error("Create() expected error for invalid username, got nil")
	}
}

func TestCreate_InvalidPassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	newUser := createTestUser(0, validUsername, invalidPassword)

	_, err := service.Create(newUser)
	if err == nil {
		t.Error("Create() expected error for invalid password, got nil")
	}
}

func TestCreate_RepositoryError(t *testing.T) {
	mockRepo := setupTestData()
	mockRepo.createErr = errors.New("database error")

	service := &UserService{repo: mockRepo}

	newUser := createTestUser(0, validUsername, validPassword)

	_, err := service.Create(newUser)
	if err == nil {
		t.Error("Create() expected error from repository, got nil")
	}
}

/*** UPDATE OPERATIONS TESTS ***/

func TestUpdate_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	existingUser := mockRepo.users[0]

	updatedUser := createTestUser(existingUser.ID, validUsername, existingUser.Password)

	err := service.Update(updatedUser)
	if err != nil {
		t.Fatalf("Update() error = %v instead of nil", err)
	}

	// Verify the username was updated
	user, _ := mockRepo.GetByID(existingUser.ID)
	if user == nil {
		t.Error("Update() failed to update username")
	}

	err = compareUsers(user, &models.User{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Password:  updatedUser.Password,
		AppAdmin:  updatedUser.AppAdmin,
		CreatedAt: updatedUser.CreatedAt,
		Avatar:    updatedUser.Avatar,
		Language:  updatedUser.Language,
		AppTheme:  updatedUser.AppTheme,
	})

	if err != nil {
		t.Error("Update() returned user with mismatched fields: " + err.Error())
	}
}

func TestUpdate_UserNotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	updatedUser := createTestUser(int64(invalidId), validUsername, validPassword)

	err := service.Update(updatedUser)
	if err == nil {
		t.Error("Update() expected error for non-existent user, got nil")
	}
}

func TestUpdate_DuplicateUsername(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	firstUser := mockRepo.users[0]
	secondUser := mockRepo.users[1]

	updatedUser := createTestUser(firstUser.ID, secondUser.Username, validPassword)

	err := service.Update(updatedUser)
	if err == nil {
		t.Error("Update() expected error for duplicate username, got nil")
	}
}

func TestUpdate_InvalidUsername(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	updatedUser := createTestUser(1, invalidUsername, validPassword) // too short

	err := service.Update(updatedUser)
	if err == nil {
		t.Error("Update() expected error for invalid username, got nil")
	}
}

func TestUpdate_AvatarLanguageTheme(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	existingUser := mockRepo.users[0]

	avatar2 := enums.Avatar2
	updatedUser := createTestUser(existingUser.ID, existingUser.Username, existingUser.Password)
	updatedUser.Avatar = &avatar2
	updatedUser.Language = enums.French
	updatedUser.AppTheme = enums.Dark

	err := service.Update(updatedUser)
	if err != nil {
		t.Fatalf("Update() error = %v instead of nil", err)
	}

	// Verify changes
	user, _ := mockRepo.GetByID(1)
	if *user.Avatar != enums.Avatar2 {
		t.Fatalf("Update() failed to update avatar: %v instead of %v", *user.Avatar, enums.Avatar2)
	}
	if user.Language != enums.French {
		t.Fatalf("Update() failed to update language: %v instead of %v", user.Language, enums.French)
	}
	if user.AppTheme != enums.Dark {
		t.Fatalf("Update() failed to update theme: %v instead of %v", user.AppTheme, enums.Dark)
	}
}

func TestUpdateAdminRole_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	err := service.UpdateAdminRole(1, true)
	if err != nil {
		t.Fatalf("UpdateAdminRole() error = %v instead of nil", err)
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
	if err == nil {
		t.Error("UpdateAdminRole() expected error for non-existent user, got nil")
	}
}

func TestUpdatePassword_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, user.Password, validPassword)
	if err != nil {
		t.Fatalf("UpdatePassword() error = %v instead of nil", err)
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
	if err == nil {
		t.Error("UpdatePassword() expected error for empty old password, got nil")
	}

	err = service.UpdatePassword(user.ID, user.Password, "")
	if err == nil {
		t.Error("UpdatePassword() expected error for empty new password, got nil")
	}
}

func TestUpdatePassword_IncorrectOldPassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, validPassword, validPassword+"123")
	if err == nil {
		t.Error("UpdatePassword() expected error for incorrect old password, got nil")
	}
}

func TestUpdatePassword_SamePassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, user.Password, user.Password)
	if err != nil {
		t.Errorf("UpdatePassword() error = %v instead of nil for same password", err)
	}
}

func TestUpdatePassword_InvalidNewPassword(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(user.ID, user.Password, invalidPassword)
	if err == nil {
		t.Error("UpdatePassword() expected error for invalid new password, got nil")
	}
}

func TestUpdatePassword_UserNotFound(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	err := service.UpdatePassword(int64(invalidId), user.Password, validPassword)
	if err == nil {
		t.Error("UpdatePassword() expected error for non-existent user, got nil")
	}
}

/*** DELETE OPERATIONS TESTS ***/

func TestDelete_Success(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	id := mockRepo.users[0].ID

	err := service.Delete(id)
	if err != nil {
		t.Fatalf("Delete() error = %v instead of nil", err)
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
	if err == nil {
		t.Error("Delete() expected error for non-existent user, got nil")
	}

	// Verify it's the ErrUserNotFound error
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Delete() error should wrap ErrUserNotFound, got %v", err)
	}
}

func TestDelete_RepositoryError(t *testing.T) {
	mockRepo := setupTestData()
	mockRepo.deleteErr = errors.New("database error")

	service := &UserService{repo: mockRepo}

	err := service.Delete(1)
	if err == nil {
		t.Error("Delete() expected error from repository, got nil")
	}
}

/*** HELPER FUNCTIONS TESTS ***/

func TestUserExists_True(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	id := mockRepo.users[0].ID

	exists, err := service.userExists(id)
	if err != nil {
		t.Fatalf("userExists() error = %v instead of nil", err)
	}

	if !exists {
		t.Error("userExists() = false instead of true")
	}
}

func TestUserExists_False(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := &UserService{repo: mockRepo}

	exists, err := service.userExists(int64(invalidId))
	if err != nil {
		t.Fatalf("userExists() error = %v instead of nil", err)
	}

	if exists {
		t.Error("userExists() = true instead of false")
	}
}

func TestUsernameExists_True(t *testing.T) {
	mockRepo := setupTestData()
	service := &UserService{repo: mockRepo}

	user := mockRepo.users[0]

	exists, err := service.usernameExists(user.Username)
	if err != nil {
		t.Fatalf("usernameExists() error = %v instead of nil", err)
	}

	if !exists {
		t.Fatalf("usernameExists() = false instead of true")
	}
}

func TestUsernameExists_False(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := &UserService{repo: mockRepo}

	exists, err := service.usernameExists(invalidUsername)
	if err != nil {
		t.Fatalf("usernameExists() error = %v instead of nil", err)
	}

	if exists {
		t.Error("usernameExists() = true instead of false")
	}
}
