package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mattn/go-sqlite3"
	customErrors "github.com/zouipo/yumsday/backend/internal/errors"

	"github.com/zouipo/yumsday/backend/internal/dtos"
	"github.com/zouipo/yumsday/backend/internal/mappers"
	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/models/enums"
)

var (
	testUser1 = createTestUser(1, "user1", "password123")
	testUser2 = createTestUser(2, "user2", "password456")
	testUser3 = createTestUser(3, "user3", "password789")

	notFoundErr = "No row found"
	conflictErr = "Conflict with User: already exists"

	validUsername = "validuser"
	validPassword = "ValidPass123"

	invalidId       = -1
	invalidUsername = "_"
	invalidPassword = "a"
)

// MockUserService is a mock implementation of UserService for testing handlers
type MockUserService struct {
	users            []models.User
	nextID           int64
	getAllErr        error
	getByIDErr       error
	getByUsernameErr error
	createErr        error
	updateErr        error
	deleteErr        error
	updateRoleErr    error
	updatePassErr    error
}

// NewMockUserService creates a new mock service with some test data
func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:  make([]models.User, 0),
		nextID: 1,
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
	return nil, customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), errors.New(notFoundErr))
}

func (m *MockUserService) GetByUsername(username string) (*models.User, error) {
	if m.getByUsernameErr != nil {
		return nil, m.getByUsernameErr
	}

	for i := range m.users {
		if m.users[i].Username == username {
			return &m.users[i], nil
		}
	}
	return nil, customErrors.NewEntityNotFoundError("User", username, errors.New(notFoundErr))
}

func (m *MockUserService) Create(user *models.User) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}
	user.ID = m.nextID
	m.nextID++

	m.users = append(m.users, *user)
	return user.ID, nil
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
	return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(user.ID, 10), errors.New(notFoundErr))
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
	return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), errors.New(notFoundErr))
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
	return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), errors.New(notFoundErr))
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
	return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), errors.New(notFoundErr))
}

/*** HELPER FUNCTIONS ***/

func (m *MockUserService) addUser(user *models.User) {
	user.ID = m.nextID
	m.nextID++
	m.users = append(m.users, *user)
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

// setupTestData creates a fresh mock repository with predefined test users for test independence.
// It is run at the start of each test to ensure a consistent state and avoid test interference.
func setupTestData() *MockUserService {
	mockService := NewMockUserService()

	// users
	mockService.addUser(testUser1)
	mockService.addUser(testUser2)
	mockService.addUser(testUser3)

	return mockService
}

// compareUsers compares the fields of a user to a userDTO;
// returns an error if any field does not match.
func compareUserToUserDto(userDto *dtos.UserDto, user *models.User) error {
	if userDto.ID != user.ID {
		return fmt.Errorf("ID = %d instead of %d", userDto.ID, user.ID)
	}
	if userDto.Username != user.Username {
		return fmt.Errorf("username = %s instead of %s", userDto.Username, user.Username)
	}
	if userDto.AppAdmin != user.AppAdmin {
		return fmt.Errorf("appAdmin ='%v'instead of %v", userDto.AppAdmin, user.AppAdmin)
	}
	// Verify both dates are within the last 2 minutes from now
	now := time.Now()
	threshold := now.Add(-2 * time.Minute)
	if userDto.CreatedAt.Before(threshold) || userDto.CreatedAt.After(now) {
		return fmt.Errorf("createdAt ='%v'is not within the last 2 minutes (threshold: %v, now: %v)", userDto.CreatedAt, threshold, now)
	}
	if user.CreatedAt.Before(threshold) || user.CreatedAt.After(now) {
		return fmt.Errorf("expected createdAt ='%v'is not within the last 2 minutes (threshold: %v, now: %v)", user.CreatedAt, threshold, now)
	}
	// Check Avatar with nil handling
	if (userDto.Avatar == nil) != (user.Avatar == nil) {
		return fmt.Errorf("avatar ='%v'instead of %v", userDto.Avatar, user.Avatar)
	}
	if userDto.Avatar != nil && user.Avatar != nil && *userDto.Avatar != *user.Avatar {
		return fmt.Errorf("avatar ='%v'instead of %v", *userDto.Avatar, *user.Avatar)
	}
	if userDto.Language != user.Language {
		return fmt.Errorf("language ='%v'instead of %v", userDto.Language, user.Language)
	}
	if userDto.AppTheme != user.AppTheme {
		return fmt.Errorf("appTheme ='%v'instead of %v", userDto.AppTheme, user.AppTheme)
	}
	return nil
}

// compareUserToNewUserDto compares the fields of a NewUserDto to a User;
// returns an error if any field does not match.
func compareUserToNewUserDto(user *models.User, newUserDto *dtos.NewUserDto) error {
	if user.Username != newUserDto.Username {
		return fmt.Errorf("username = %s instead of %s", user.Username, newUserDto.Username)
	}
	if user.AppAdmin != newUserDto.AppAdmin {
		return fmt.Errorf("appAdmin ='%v'instead of %v", user.AppAdmin, newUserDto.AppAdmin)
	}
	if user.Password != newUserDto.Password {
		return fmt.Errorf("password = %s instead of %s", user.Password, newUserDto.Password)
	}
	// Check Avatar with nil handling
	if (user.Avatar == nil) != (newUserDto.Avatar == nil) {
		return fmt.Errorf("avatar ='%v'instead of %v", user.Avatar, newUserDto.Avatar)
	}
	if user.Avatar != nil && newUserDto.Avatar != nil && *user.Avatar != *newUserDto.Avatar {
		return fmt.Errorf("avatar ='%v'instead of %v", *user.Avatar, *newUserDto.Avatar)
	}
	if user.Language != newUserDto.Language {
		return fmt.Errorf("language ='%v'instead of %v", user.Language, newUserDto.Language)
	}
	if user.AppTheme != newUserDto.AppTheme {
		return fmt.Errorf("appTheme ='%v'instead of %v", user.AppTheme, newUserDto.AppTheme)
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

func TestGetUsersAll_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	// Simulates a request to GET /user without query parameters to get all users
	r := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d instead of %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	var users []dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&users)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != len(mockService.users) {
		t.Errorf("expected %d users instead of %d", len(mockService.users), len(users))
	}
}

func TestGetUsersAll_RepoError(t *testing.T) {
	mockService := NewMockUserService()
	errMessage := errors.New("Failed to fetch users")
	mockService.getAllErr = customErrors.NewInternalServerError("Failed to fetch users", errMessage)

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	statusCode := mockService.getAllErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), errMessage.Error()) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage.Error(), w.Body.String())
	}
}

func TestGetUsersByUsername_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	user := mockService.users[0]

	r := httptest.NewRequest(http.MethodGet, "/user?username="+user.Username, nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d instead of %d", http.StatusOK, w.Code)
	}

	var users []dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&users)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("expected one user instead of %d", len(users))
	}

	if len(users) == 1 {
		if err := compareUserToUserDto(&users[0], &user); err != nil {
			t.Error("GetByUsername() returned user with mismatched fields: " + err.Error())
		}
	}
}

func TestGetUsersByUsername_NotFound(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users = []models.User{}

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user?username="+invalidUsername, nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d instead of %d", http.StatusNotFound, w.Code)
	}

	if !strings.Contains(w.Body.String(), notFoundErr) {
		t.Errorf("expected error message containing '%s' instead of '%s'", notFoundErr, w.Body.String())
	}
}

func TestGetUsersByUsername_EmptyUsername(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users = []models.User{}

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user?username=", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d instead of %d", http.StatusNotFound, w.Code)
	}

	if !strings.Contains(w.Body.String(), notFoundErr) {
		t.Errorf("expected error message containing '%s' instead of '%s'", notFoundErr, w.Body.String())
	}
}

// TestGetUsers_MultipleQueryParams tests the getUsers handler with multiple username query parameters
func TestGetUsers_MultipleQueryParams(t *testing.T) {
	mockService := setupTestData()
	handler := NewUserHandler(mockService)

	user1 := mockService.users[0]
	user2 := mockService.users[1]

	// Multiple username parameters
	r := httptest.NewRequest(http.MethodGet, "/user?username="+user1.Username+"&username="+user2.Username, nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d instead of %d", http.StatusBadRequest, w.Code)
	}

	expectedError := "Missing or invalid query parameters"
	if !strings.Contains(w.Body.String(), expectedError) {
		t.Errorf("expected error message containing '%s' instead of '%s'", expectedError, w.Body.String())
	}
}

// TestGetUsers_InvalidQueryParams tests the getUsers handler with invalid query parameters
func TestGetUsers_InvalidQueryParams(t *testing.T) {
	mockService := setupTestData()
	handler := NewUserHandler(mockService)

	// Multiple username parameters
	r := httptest.NewRequest(http.MethodGet, "/user?random=ok", nil)
	w := httptest.NewRecorder()

	handler.getUsers(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d instead of %d", http.StatusBadRequest, w.Code)
	}

	// A améliorer avec une erreur custom
	expectedError := "Missing or invalid query parameters"
	if !strings.Contains(w.Body.String(), expectedError) {
		t.Errorf("expected error message containing '%s' instead of '%s'", expectedError, w.Body.String())
	}
}

func TestGetUserByID_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	expected := mockService.users[0]

	r := httptest.NewRequest(http.MethodGet, "/user/"+strconv.FormatInt(expected.ID, 10), nil)
	// Add the ID to the context as the middleware would do
	ctx := context.WithValue(r.Context(), "id", expected.ID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.getUserByID(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d instead of %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	var actual dtos.UserDto
	err := json.NewDecoder(w.Body).Decode(&actual)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if err := compareUserToUserDto(&actual, &expected); err != nil {
		t.Error("GetByUsername() returned user with mismatched fields: " + err.Error())
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users = []models.User{}

	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/user/"+strconv.FormatInt(int64(invalidId), 10), nil)
	ctx := context.WithValue(r.Context(), "id", int64(invalidId))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.getUserByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d instead of %d", http.StatusNotFound, w.Code)
	}

	if !strings.Contains(w.Body.String(), notFoundErr) {
		t.Errorf("expected error message containing '%s' instead of '%s'", notFoundErr, w.Body.String())
	}
}

/*** CREATE OPERATIONS TESTS ***/

func TestCreateUser_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: validUsername,
		Password: validPassword,
		AppAdmin: false,
		Avatar:   &avatar,
		Language: enums.English,
		AppTheme: enums.Light,
	}

	usersNb := len(mockService.users)

	// Convert the newUser DTO to JSON
	body, _ := json.Marshal(newUser)
	// Create a POST request with the JSON body
	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d instead of %d", http.StatusCreated, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	var result map[string]int
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["id"] != int(mockService.nextID-1) {
		t.Errorf("expected id %d instead of %d", int(mockService.nextID-1), result["id"])
	}

	if usersNb+1 != len(mockService.users) {
		t.Errorf("expected %d users instead of %d", usersNb+1, len(mockService.users))
	}

	user, err := mockService.GetByID((int64)(result["id"]))
	if err != nil {
		t.Fatalf("failed to retrieve created user: %v", err)
	}

	if err := compareUserToNewUserDto(user, &newUser); err != nil {
		t.Error("Created user has mismatched fields: " + err.Error())
	}
}

func TestCreateUser_Success_AvatarNil(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	newUser := dtos.NewUserDto{
		Username: validUsername,
		Password: validPassword,
		AppAdmin: false,
		Avatar:   nil,
		Language: enums.English,
		AppTheme: enums.Light,
	}

	usersNb := len(mockService.users)

	// Convert the newUser DTO to JSON
	body, _ := json.Marshal(newUser)
	// Create a POST request with the JSON body
	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d instead of %d", http.StatusCreated, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	var result map[string]int
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["id"] != int(mockService.nextID-1) {
		t.Errorf("expected id %d instead of %d", int(mockService.nextID-1), result["id"])
	}

	if usersNb+1 != len(mockService.users) {
		t.Errorf("expected %d users instead of %d", usersNb+1, len(mockService.users))
	}

	user, err := mockService.GetByID((int64)(result["id"]))
	if err != nil {
		t.Fatalf("failed to retrieve created user: %v", err)
	}

	if err := compareUserToNewUserDto(user, &newUser); err != nil {
		t.Error("Created user has mismatched fields: " + err.Error())
	}
}

func TestCreateUser_InvalidBody(t *testing.T) {
	mockService := setupTestData()
	handler := NewUserHandler(mockService)

	usersNb := len(mockService.users)

	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d instead of %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "invalid character") {
		t.Errorf("expected error message containing JSON decode error, got: %s", w.Body.String())
	}

	if usersNb != len(mockService.users) {
		t.Errorf("expected %d users instead of %d", usersNb, len(mockService.users))
	}
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockService := setupTestData()
	mockService.createErr = customErrors.NewValidationError("username", "invalid username format", nil)

	handler := NewUserHandler(mockService)

	usersNb := len(mockService.users)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: invalidUsername,
		Password: invalidPassword,
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

	statusCode := mockService.createErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	errMessage := "Validation error on field 'username': invalid username format"
	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	if usersNb != len(mockService.users) {
		t.Errorf("expected %d users instead of %d", usersNb, len(mockService.users))
	}
}

func TestCreateUser_ConflictError(t *testing.T) {
	mockService := setupTestData()
	mockService.createErr = customErrors.NewConflictError("User", "already exists", sqlite3.ErrConstraintUnique)

	handler := NewUserHandler(mockService)

	usersNb := len(mockService.users)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: mockService.users[0].Username,
		Password: validPassword,
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

	statusCode := mockService.createErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), conflictErr) {
		t.Errorf("expected error message containing '%s' instead of '%s'", conflictErr, w.Body.String())
	}

	if usersNb != len(mockService.users) {
		t.Errorf("expected %d users instead of %d", usersNb, len(mockService.users))
	}
}

func TestCreateUser_RepoError(t *testing.T) {
	mockService := setupTestData()
	errMessage := "Failed to create user"
	mockService.createErr = customErrors.NewInternalServerError(errMessage, nil)

	handler := NewUserHandler(mockService)

	usersNb := len(mockService.users)

	avatar := enums.Avatar1
	newUser := dtos.NewUserDto{
		Username: mockService.users[0].Username,
		Password: validPassword,
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

	statusCode := mockService.createErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	if usersNb != len(mockService.users) {
		t.Errorf("expected %d users instead of %d", usersNb, len(mockService.users))
	}
}

/*** UPDATE OPERATIONS TESTS ***/

func TestUpdateUser_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	avatar := enums.Avatar2
	user := mappers.ToUserDtoNoPassword(&mockService.users[0])
	user.Username = validUsername
	user.Avatar = &avatar
	user.Language = enums.French
	user.AppTheme = enums.Dark

	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d instead of %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if err := compareUserToUserDto(user, actual); err != nil {
		t.Error("Updated user has mismatched fields: " + err.Error())
	}
}

func TestUpdateUser_InvalidBody(t *testing.T) {
	mockService := setupTestData()
	handler := NewUserHandler(mockService)

	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d instead of %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "invalid character") {
		t.Errorf("expected error message containing JSON decode error, got: %s", w.Body.String())
	}
}

func TestUpdateUser_ConflictError(t *testing.T) {
	mockService := setupTestData()
	mockService.updateErr = customErrors.NewConflictError("User", "already exists", sqlite3.ErrConstraintUnique)

	handler := NewUserHandler(mockService)

	user := mappers.ToUserDtoNoPassword(&mockService.users[0])
	user.Username = mockService.users[1].Username

	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	statusCode := mockService.updateErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), conflictErr) {
		t.Errorf("expected error message containing '%s' instead of '%s'", conflictErr, w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if err := compareUserToUserDto(mappers.ToUserDtoNoPassword(&mockService.users[0]), actual); err != nil {
		t.Error("Updated user has mismatched fields: " + err.Error())
	}
}

func TestUpdateUser_ValidationError(t *testing.T) {
	mockService := setupTestData()
	mockService.updateErr = customErrors.NewValidationError("username", "Invalid username format", nil)

	handler := NewUserHandler(mockService)

	user := mappers.ToUserDtoNoPassword(&mockService.users[0])
	user.Username = invalidUsername

	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	statusCode := mockService.updateErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	errMessage := "Validation error on field 'username': Invalid username format"
	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if err := compareUserToUserDto(mappers.ToUserDtoNoPassword(&mockService.users[0]), actual); err != nil {
		t.Error("updated user has mismatched fields: " + err.Error())
	}
}

func TestUpdateUser_RepoError(t *testing.T) {
	mockService := setupTestData()
	errMessage := "Failed to update user"
	mockService.updateErr = customErrors.NewInternalServerError(errMessage, nil)

	handler := NewUserHandler(mockService)

	avatar := enums.Avatar2
	user := mappers.ToUserDtoNoPassword(&mockService.users[0])
	user.Username = validUsername
	user.Avatar = &avatar
	user.Language = enums.French
	user.AppTheme = enums.System

	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.updateUser(w, r)

	statusCode := mockService.updateErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if err := compareUserToUserDto(mappers.ToUserDtoNoPassword(&mockService.users[0]), actual); err != nil {
		t.Error("Updated user has mismatched fields: " + err.Error())
	}
}

func TestUpdateUserAdminRole_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	user := mockService.users[0]
	adminRole := !user.AppAdmin

	rolePayload := dtos.AdminRolePayload{AppAdmin: adminRole}
	body, _ := json.Marshal(rolePayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/"+strconv.FormatInt(user.ID, 10)+"/role", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserAdminRole(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d instead of %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if actual.AppAdmin != adminRole {
		t.Errorf("expected appAdmin ='%v'instead of %v", adminRole, actual.AppAdmin)
	}
}

func TestUpdateUserAdminRole_InvalidBody(t *testing.T) {
	mockService := setupTestData()
	handler := NewUserHandler(mockService)

	user := mockService.users[0]

	r := httptest.NewRequest(http.MethodPatch, "/user/"+strconv.FormatInt(user.ID, 10)+"/role", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserAdminRole(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d instead of %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "invalid character") {
		t.Errorf("expected error message containing JSON decode error, got: %s", w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if actual.AppAdmin != user.AppAdmin {
		t.Errorf("expected appAdmin ='%v'instead of %v", user.AppAdmin, actual.AppAdmin)
	}
}

func TestUpdateUserAdminRole_RepoError(t *testing.T) {
	mockService := setupTestData()
	errMessage := "Failed to update user admin role"
	mockService.updateRoleErr = customErrors.NewInternalServerError(errMessage, nil)

	handler := NewUserHandler(mockService)

	user := mockService.users[0]
	adminRole := !user.AppAdmin

	rolePayload := dtos.AdminRolePayload{AppAdmin: adminRole}
	body, _ := json.Marshal(rolePayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/"+strconv.FormatInt(user.ID, 10)+"/role", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserAdminRole(w, r)

	statusCode := mockService.updateRoleErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if actual.AppAdmin != user.AppAdmin {
		t.Errorf("expected appAdmin ='%v'instead of %v", user.AppAdmin, actual.AppAdmin)
	}
}

func TestUpdateUserPassword_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	user := mockService.users[0]

	passwordPayload := dtos.PasswordPayload{
		OldPassword: user.Password,
		NewPassword: validPassword,
	}
	body, _ := json.Marshal(passwordPayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/"+strconv.FormatInt(user.ID, 10)+"/password", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserPassword(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d instead of %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if actual.Password != validPassword {
		t.Errorf("expected password ='%v'instead of %v", validPassword, actual.Password)
	}
}

func TestUpdateUserPassword_InvalidBody(t *testing.T) {
	mockService := setupTestData()
	handler := NewUserHandler(mockService)

	user := mockService.users[0]

	r := httptest.NewRequest(http.MethodPatch, "/user/"+strconv.FormatInt(user.ID, 10)+"/password", bytes.NewReader([]byte("invalid json")))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserPassword(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d instead of %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "invalid character") {
		t.Errorf("expected error message containing JSON decode error, got: %s", w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if actual.Password != user.Password {
		t.Errorf("expected password ='%v'instead of %v", user.Password, actual.Password)
	}
}

func TestUpdateUserPassword_ValidationError(t *testing.T) {
	mockService := setupTestData()
	mockService.updatePassErr = customErrors.NewValidationError("password", "Invalid password length", nil)

	handler := NewUserHandler(mockService)

	user := mockService.users[0]

	passwordPayload := dtos.PasswordPayload{
		OldPassword: user.Password,
		NewPassword: invalidPassword,
	}
	body, _ := json.Marshal(passwordPayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/"+strconv.FormatInt(user.ID, 10)+"/password", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserPassword(w, r)

	statusCode := mockService.updatePassErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	errMessage := "Validation error on field 'password': Invalid password length"
	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if actual.Password != user.Password {
		t.Errorf("expected password ='%v'instead of %v", user.Password, actual.Password)
	}
}

func TestUpdateUserPassword_RepoError(t *testing.T) {
	mockService := setupTestData()
	errMessage := "Failed to update user"
	mockService.updatePassErr = customErrors.NewInternalServerError(errMessage, nil)

	handler := NewUserHandler(mockService)

	user := mockService.users[0]

	passwordPayload := dtos.PasswordPayload{
		OldPassword: user.Password,
		NewPassword: validPassword,
	}
	body, _ := json.Marshal(passwordPayload)

	r := httptest.NewRequest(http.MethodPatch, "/user/"+strconv.FormatInt(user.ID, 10)+"/password", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.updateUserPassword(w, r)

	statusCode := mockService.updatePassErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	actual, err := mockService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if actual.Password != user.Password {
		t.Errorf("expected password ='%v'instead of %v", user.Password, actual.Password)
	}
}

/*** DELETE OPERATIONS TESTS ***/

func TestDeleteUser_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	usersNb := len(mockService.users)
	user := mockService.users[0]

	r := httptest.NewRequest(http.MethodDelete, "/user/"+strconv.FormatInt(user.ID, 10), nil)
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.deleteUser(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d instead of %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json instead of %s", contentType)
	}

	if len(mockService.users) != usersNb-1 {
		t.Errorf("expected %d users after deletion instead of %d", usersNb-1, len(mockService.users))
	}

	if _, err := mockService.GetByID(user.ID); err == nil {
		t.Errorf("expected error when retrieving deleted user, but got none")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)

	usersNb := len(mockService.users)

	r := httptest.NewRequest(http.MethodDelete, "/user/"+strconv.FormatInt(int64(invalidId), 10), nil)
	ctx := context.WithValue(r.Context(), "id", int64(invalidId))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.deleteUser(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d instead of %d", http.StatusNotFound, w.Code)
	}

	if len(mockService.users) != usersNb {
		t.Errorf("expected %d users after failed deletion instead of %d", usersNb, len(mockService.users))
	}
}

func TestDeleteUser_RepoError(t *testing.T) {
	mockService := setupTestData()
	errMessage := "Failed to delete user"
	mockService.deleteErr = customErrors.NewInternalServerError(errMessage, nil)

	handler := NewUserHandler(mockService)

	usersNb := len(mockService.users)
	user := mockService.users[0]

	r := httptest.NewRequest(http.MethodDelete, "/user/"+strconv.FormatInt(user.ID, 10), nil)
	ctx := context.WithValue(r.Context(), "id", int64(user.ID))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.deleteUser(w, r)

	statusCode := mockService.deleteErr.(*customErrors.AppError).StatusCode
	if w.Code != statusCode {
		t.Errorf("expected status %d instead of %d", statusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), errMessage) {
		t.Errorf("expected error message containing '%s' instead of '%s'", errMessage, w.Body.String())
	}

	if len(mockService.users) != usersNb {
		t.Errorf("expected %d users after failed deletion instead of %d", usersNb, len(mockService.users))
	}

	if _, err := mockService.GetByID(user.ID); err != nil {
		t.Errorf("expected user to still exist after failed deletion, but got error: %v", err)
	}
}

// TestRegisterRoutes tests the RegisterRoutes method
func TestRegisterRoutes_Success(t *testing.T) {
	mockService := setupTestData()

	handler := NewUserHandler(mockService)
	mux := http.NewServeMux()

	handler.RegisterRoutes(mux, "/test/api/user")

	// Test that routes are registered - test GET /api/user
	r := httptest.NewRequest(http.MethodGet, "/test/api/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d for GET /test/api/user instead of %d", http.StatusOK, w.Code)
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
	r = httptest.NewRequest(http.MethodPost, "/test/api/user", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d for POST /test/api/user instead of %d", http.StatusCreated, w.Code)
	}
}
