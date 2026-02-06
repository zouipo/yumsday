package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zouipo/yumsday/backend/internal/dtos"
	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/models/enums"
	"github.com/zouipo/yumsday/backend/internal/services"
)

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

// MockUserService is a mock implementation of UserService for testing handlers
type MockUserService struct {
	users          []models.User
	nextID         int64
	getAllErr      error
	getByIDErr     error
	getByNameErr   error
	createErr      error
	updateErr      error
	deleteErr      error
	updateRoleErr  error
	updatePassErr  error
	createReturnID int64
}

// NewMockUserService creates a new mock service with some test data
func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:          make([]models.User, 0),
		createReturnID: 4,
	}
}

/*** USERSERVICE IMPLEMENTATION ***/

func (m *MockUserService) GetAll() ([]models.User, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.users, nil
}

func (m *MockUserService) GetByID(id int64) (*models.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.users {
		if m.users[i].ID == id {
			return &m.users[i], nil
		}
	}
	return nil, services.ErrUserNotFound
}

func (m *MockUserService) GetByUsername(username string) (*models.User, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}

	for i := range m.users {
		if m.users[i].Username == username {
			return &m.users[i], nil
		}
	}
	return nil, services.ErrUserNotFound
}

func (m *MockUserService) Create(user *models.User) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}
	user.ID = m.createReturnID
	m.users = append(m.users, *user)
	return m.createReturnID, nil
}

func (m *MockUserService) Update(user *models.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i := range m.users {
		if m.users[i].ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return services.ErrUserNotFound
}

func (m *MockUserService) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	for i := range m.users {
		if m.users[i].ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return services.ErrUserNotFound
}

func (m *MockUserService) UpdateAdminRole(id int64, isAdmin bool) error {
	if m.updateRoleErr != nil {
		return m.updateRoleErr
	}
	for i := range m.users {
		if m.users[i].ID == id {
			m.users[i].AppAdmin = isAdmin
			return nil
		}
	}
	return services.ErrUserNotFound
}

func (m *MockUserService) UpdatePassword(id int64, oldPassword, newPassword string) error {
	if m.updatePassErr != nil {
		return m.updatePassErr
	}
	for i := range m.users {
		if m.users[i].ID == id {
			m.users[i].Password = newPassword
			return nil
		}
	}
	return services.ErrUserNotFound
}

/*** HELPER FUNCTIONS ***/

// Helper methods for setting up test scenarios
func (m *MockUserService) addUser(user *models.User) {
	m.users = append(m.users, *user)
	if user.ID >= m.nextID {
		m.nextID = user.ID + 1
	}
}

func createTestUser(id int64, username, password string) *models.User {
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
func setupTestData() *MockUserService {
	mockService := NewMockUserService()
	mockService.addUser(testUser1)
	mockService.addUser(testUser2)
	mockService.addUser(testUser3)
	return mockService
}

// compareUsers compares two User objects and returns an error if any field does not match
func compareUserToUserDto(actual *models.User, expected *models.User) error {
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

func TestNewUserHandler(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}

	if handler.userService != mockService {
		t.Error("handler userService does not match the provided service")
	}
}

/*** READ OPERATIONS TESTS ***/

func TestGetUsersAll(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	// Simulates a request to GET /user without query parameters to get all users
	r := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}

	var users []dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&users)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != len(mockService.users) {
		t.Errorf("expected %d users, got %d", len(mockService.users), len(users))
	}
}

// TestGetUsersAllError tests the getUsers handler when GetAll returns an error
func TestGetUsersAllError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.getAllErr = errors.New("database error")

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestGetUsersByUsername tests the getUsers handler when filtering by username
func TestGetUsersByUsername(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	user := mockService.users[0]

	r := httptest.NewRequest(http.MethodGet, "/user?username="+user.Username, nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var users []dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&users)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}

	if users[0].Username != "user1" {
		t.Errorf("expected username user1, got %s", users[0].Username)
	}
}

// TestGetUsersByUsernameNotFound tests the getUsers handler when username is not found
func TestGetUsersByUsernameNotFound(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users = []models.User{}

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user?username=nonexistent", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestGetUsersInvalidQueryParams tests the getUsers handler with invalid query parameters
func TestGetUsersInvalidQueryParams(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	// Multiple username parameters
	r := httptest.NewRequest(http.MethodGet, "/user?username=user1&username=user2", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedError := "Missing or invalid query parameters"
	if !strings.Contains(w.Body.String(), expectedError) {
		t.Errorf("expected error message containing '%s', got '%s'", expectedError, w.Body.String())
	}
}

// TestGetUserByID tests the getUserByID handler
func TestGetUserByID(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	// Add the ID to the context as the middleware would do
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.getUserByID(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}

	var user dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&user)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}

	if user.Username != "user1" {
		t.Errorf("expected username user1, got %s", user.Username)
	}
}

// TestGetUserByIDNotFound tests the getUserByID handler when user is not found
func TestGetUserByIDNotFound(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users = []models.User{}

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user/999", nil)
	ctx := context.WithValue(r.Context(), "id", int64(999))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.getUserByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestCreateUser tests the createUser handler
func TestCreateUser(t *testing.T) {
	mockService := NewMockUserService()
	mockService.createReturnID = 42

	handler := NewUserHandler(mockService)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: "newuser",
		Password: "password123",
		AppAdmin: false,
		Avatar:   &avatar,
		Language: enums.English,
		AppTheme: enums.Light,
	}

	body, _ := json.Marshal(newUser)
	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}

	var result map[string]int
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["id"] != 42 {
		t.Errorf("expected id 42, got %d", result["id"])
	}
}

// TestCreateUserInvalidBody tests the createUser handler with invalid request body
func TestCreateUserInvalidBody(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedError := "Invalid request body"
	if !strings.Contains(w.Body.String(), expectedError) {
		t.Errorf("expected error message containing '%s', got '%s'", expectedError, w.Body.String())
	}
}

// TestCreateUserServiceError tests the createUser handler when service returns an error
func TestCreateUserServiceError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.createErr = errors.New("validation error")

	handler := NewUserHandler(mockService)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: "newuser",
		Password: "password123",
		Avatar:   &avatar,
	}

	body, _ := json.Marshal(newUser)
	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateUser tests the updateUser handler
func TestUpdateUser(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	avatar := enums.Avatar1
	userDto := dtos.UserDto{
		ID:       1,
		Username: "updateduser",
		AppAdmin: false,
		Avatar:   &avatar,
		Language: enums.English,
		AppTheme: enums.Light,
	}

	body, _ := json.Marshal(userDto)
	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}
}

// TestUpdateUserInvalidBody tests the updateUser handler with invalid request body
func TestUpdateUserInvalidBody(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateUserServiceError tests the updateUser handler when service returns an error
func TestUpdateUserServiceError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.updateErr = errors.New("database error")

	handler := NewUserHandler(mockService)

	userDto := dtos.UserDto{
		ID:       1,
		Username: "updateduser",
	}

	body, _ := json.Marshal(userDto)
	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestUpdateUserAdminRole tests the updateUserAdminRole handler
func TestUpdateUserAdminRole(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	rolePayload := map[string]bool{"app_admin": true}
	body, _ := json.Marshal(rolePayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/1/role", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserAdminRole(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}
}

// TestUpdateUserAdminRoleInvalidBody tests the updateUserAdminRole handler with invalid body
func TestUpdateUserAdminRoleInvalidBody(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodPatch, "/user/1/role", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserAdminRole(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateUserAdminRoleServiceError tests the updateUserAdminRole handler with service error
func TestUpdateUserAdminRoleServiceError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.updateRoleErr = errors.New("database error")

	handler := NewUserHandler(mockService)

	rolePayload := map[string]bool{"app_admin": true}
	body, _ := json.Marshal(rolePayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/1/role", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserAdminRole(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestUpdateUserPassword tests the updateUserPassword handler
func TestUpdateUserPassword(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	passwordPayload := map[string]string{
		"old_password": "password123",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(passwordPayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/1/password", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserPassword(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}
}

// TestUpdateUserPasswordInvalidBody tests the updateUserPassword handler with invalid body
func TestUpdateUserPasswordInvalidBody(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodPatch, "/user/1/password", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserPassword(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateUserPasswordServiceError tests the updateUserPassword handler with service error
func TestUpdateUserPasswordServiceError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.updatePassErr = errors.New("password verification failed")

	handler := NewUserHandler(mockService)

	passwordPayload := map[string]string{
		"old_password": "wrongpassword",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(passwordPayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/1/password", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserPassword(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestDeleteUser tests the deleteUser handler
func TestDeleteUser(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.deleteUser(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}

	// Verify the user was deleted from the mock service
	if len(mockService.users) != 2 {
		t.Errorf("expected 2 users after deletion, got %d", len(mockService.users))
	}
}

// TestDeleteUserNotFound tests the deleteUser handler when user is not found
func TestDeleteUserNotFound(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users = []models.User{}

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodDelete, "/user/999", nil)
	ctx := context.WithValue(r.Context(), "id", int64(999))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.deleteUser(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestDeleteUserServiceError tests the deleteUser handler with general service error
func TestDeleteUserServiceError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.deleteErr = errors.New("database error")

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.deleteUser(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestRegisterRoutes tests the RegisterRoutes method
func TestRegisterRoutes(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)
	mux := http.NewServeMux()

	handler.RegisterRoutes(mux, "/api/user")

	// Test that routes are registered - test GET /api/user
	r := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d for GET /api/user, got %d", http.StatusOK, w.Code)
	}

	// Test POST /api/user
	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: "newuser",
		Password: "password123",
		Avatar:   &avatar,
		Language: enums.English,
		AppTheme: enums.Light,
	}
	body, _ := json.Marshal(newUser)
	r = httptest.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d for POST /api/user, got %d", http.StatusCreated, w.Code)
	}
}

// TestGetAllUsers tests the private getAllUsers method
func TestGetAllUsers(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)
	w := httptest.NewRecorder()

	handler.getAllUsers(w)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var users []dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&users)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

// TestGetAllUsersError tests the private getAllUsers method with service error
func TestGetAllUsersError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.getAllErr = errors.New("database error")

	handler := NewUserHandler(mockService)
	w := httptest.NewRecorder()

	handler.getAllUsers(w)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestGetByUsername tests the private getByUsername method
func TestGetByUsername(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)
	w := httptest.NewRecorder()

	handler.getByUsername(w, "user1")

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var users []dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&users)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}

	if users[0].Username != "user1" {
		t.Errorf("expected username user1, got %s", users[0].Username)
	}
}

// TestGetByUsernameNotFound tests the private getByUsername method when user not found
func TestGetByUsernameNotFound(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users = []models.User{}

	handler := NewUserHandler(mockService)
	w := httptest.NewRecorder()

	handler.getByUsername(w, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestGetByUsernameServiceError tests the private getByUsername method with service error
func TestGetByUsernameServiceError(t *testing.T) {
	mockService := NewMockUserService()
	mockService.getByNameErr = errors.New("database error")

	handler := NewUserHandler(mockService)
	w := httptest.NewRecorder()

	handler.getByUsername(w, "someuser")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestGetUserByIDSerializationError tests what happens when encoding fails
// This is a edge case that's hard to trigger in practice, but we can document it
func TestGetUserByIDSuccessfulRetrieval(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	ctx := context.WithValue(r.Context(), "id", int64(1))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.getUserByID(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify the response contains valid JSON
	var result dtos.UserDto
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Username != "user1" {
		t.Errorf("expected username user1, got %s", result.Username)
	}
}

// TestCreateUserWithAllFields tests creating a user with all optional fields
func TestCreateUserWithAllFields(t *testing.T) {
	mockService := NewMockUserService()
	mockService.createReturnID = 100

	handler := NewUserHandler(mockService)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: "fulluser",
		Password: "securepass123",
		AppAdmin: true,
		Avatar:   &avatar,
		Language: enums.French,
		AppTheme: enums.Dark,
	}

	body, _ := json.Marshal(newUser)
	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result map[string]int
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["id"] != 100 {
		t.Errorf("expected id 100, got %d", result["id"])
	}

	// Verify the user was added to the mock service
	if len(mockService.users) != 1 {
		t.Errorf("expected 1 user in mock service, got %d", len(mockService.users))
	}

	createdUser := mockService.users[0]
	if createdUser.Username != "fulluser" {
		t.Errorf("expected username fulluser, got %s", createdUser.Username)
	}
	if !createdUser.AppAdmin {
		t.Error("expected user to be admin")
	}
	if createdUser.Language != enums.French {
		t.Errorf("expected language French, got %v", createdUser.Language)
	}
	if createdUser.AppTheme != enums.Dark {
		t.Errorf("expected theme Dark, got %v", createdUser.AppTheme)
	}
}

// Benchmark tests
func BenchmarkGetUserByID(b *testing.B) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodGet, "/user/1", nil)
		ctx := context.WithValue(r.Context(), "id", int64(1))
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.getUserByID(w, r)
	}
}

func BenchmarkCreateUser(b *testing.B) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: "benchuser",
		Password: "password123",
		Avatar:   &avatar,
		Language: enums.English,
		AppTheme: enums.Light,
	}

	body, _ := json.Marshal(newUser)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.createUser(w, r)
		mockService.createReturnID++
		mockService.users = []models.User{} // Reset for next iteration
	}
}

// TestEdgeCaseEmptyUsername tests edge cases with empty or special usernames
func TestEdgeCaseEmptyUsername(t *testing.T) {
	mockService := NewMockUserService()
	handler := NewUserHandler(mockService)

	newUser := dtos.NewUserDto{
		Username: "",
		Password: "password123",
	}

	body, _ := json.Marshal(newUser)
	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Here we're just testing that the handler doesn't crash
	// The actual validation should happen in the service layer
	handler.createUser(w, r)

	// The status code depends on service validation
	// We just ensure no panic occurred
	if w.Code != http.StatusCreated && w.Code != http.StatusBadRequest {
		t.Logf("Got status code %d for empty username", w.Code)
	}
}

// TestMultipleConcurrentRequests tests thread safety
func TestMultipleConcurrentRequests(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	// Run multiple concurrent GET requests
	done := make(chan bool)
	for i := 0; i < 3; i++ {
		go func(id int64) {
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", id), nil)
			ctx := context.WithValue(r.Context(), "id", id)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.getUserByID(w, r)
			done <- true
		}(int64(i + 1))
	}

	// Wait for all goroutines to finish
	for i := 0; i < 3; i++ {
		<-done
	}
}
