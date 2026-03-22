package service

import (
	"errors"
	"net/http"
	"testing"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	username = "testuser"
	userID   = int64(1)
)

type MockSessionService struct {
	savedSessions   []*model.Session
	removedSessions []*model.Session
	removeErr       error
}

func (m *MockSessionService) GetSession(_ *http.Request) *model.Session {
	return model.NewSession()
}

func (m *MockSessionService) CookieName() string {
	return "session_id"
}

func (m *MockSessionService) Expiration() time.Duration {
	return time.Hour
}

func (m *MockSessionService) Save(session *model.Session) error {
	m.savedSessions = append(m.savedSessions, session)
	return nil
}

func (m *MockSessionService) Remove(session *model.Session) error {
	if m.removeErr != nil {
		return m.removeErr
	}
	m.removedSessions = append(m.removedSessions, session)
	return nil
}

type MockUserService struct {
	user             *model.User
	getByUsernameErr error
}

func (m *MockUserService) GetAll() ([]model.User, error) {
	return nil, nil
}

func (m *MockUserService) GetByID(_ int64) (*model.User, error) {
	return nil, nil
}

func (m *MockUserService) GetByUsername(_ string) (*model.User, error) {
	if m.getByUsernameErr != nil {
		return nil, m.getByUsernameErr
	}
	return m.user, nil
}

func (m *MockUserService) Create(_ *model.User) (int64, error) {
	return 0, nil
}

func (m *MockUserService) Update(_ *model.User) error {
	return nil
}

func (m *MockUserService) UpdateAdminRole(_ int64, _ bool) error {
	return nil
}

func (m *MockUserService) UpdatePassword(_ int64, _, _ string) error {
	return nil
}

func (m *MockUserService) Delete(_ int64) error {
	return nil
}

func createAuthTestUser(t *testing.T, id int64, username, password string) *model.User {
	t.Helper()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword() error = %v, want nil", err)
	}

	avatar := enum.Avatar1
	return &model.User{
		ID:        id,
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		Avatar:    &avatar,
		Language:  enum.English,
		AppTheme:  enum.Light,
	}
}

func TestNewAuthService(t *testing.T) {
	mockSessionService := &MockSessionService{}
	mockUserService := &MockUserService{}

	service := NewAuthService(mockSessionService, mockUserService)

	if service == nil {
		t.Fatal("NewAuthService() returned nil")
	}

	if service.sessionService == nil {
		t.Error("NewAuthService() sessionService is nil")
	}

	if service.userService == nil {
		t.Error("NewAuthService() userService is nil")
	}
}

func TestAuthenticate_Success(t *testing.T) {
	testUser := createAuthTestUser(t, userID, username, ValidPassword)
	mockUserService := &MockUserService{user: testUser}
	mockSessionService := &MockSessionService{}
	service := NewAuthService(mockSessionService, mockUserService)

	session := model.NewSession()
	err := service.Authenticate(session, username, ValidPassword)

	if err != nil {
		t.Fatalf("Authenticate() error = %v, want nil", err)
	}

	if session.UserID != testUser.ID {
		t.Errorf("Authenticate() session UserID = %d, want %d", session.UserID, testUser.ID)
	}

	if len(mockSessionService.savedSessions) != 1 {
		t.Fatalf("Authenticate() save calls = %d, want 1", len(mockSessionService.savedSessions))
	}

	if mockSessionService.savedSessions[0] != session {
		t.Error("Authenticate() saved session pointer does not match input session")
	}
}

func TestAuthenticate_UserServiceError(t *testing.T) {
	expectedErr := customErrors.NewNotFoundError("User", username, nil)
	mockUserService := &MockUserService{getByUsernameErr: expectedErr}
	mockSessionService := &MockSessionService{}
	service := NewAuthService(mockSessionService, mockUserService)

	session := model.NewSession()
	err := service.Authenticate(session, username, "irrelevant")

	if !utils.CompareErrors(err, expectedErr) {
		t.Errorf("Authenticate() error = %v, want %v", err, expectedErr)
	}

	if len(mockSessionService.savedSessions) != 0 {
		t.Error("Authenticate() should not save session when user retrieval fails")
	}
}

func TestAuthenticate_WrongPassword(t *testing.T) {
	mockUserService := &MockUserService{
		user: createAuthTestUser(t, userID, username, ValidPassword),
	}
	mockSessionService := &MockSessionService{}
	service := NewAuthService(mockSessionService, mockUserService)

	session := model.NewSession()
	err := service.Authenticate(session, username, InvalidPassword)

	expectedErr := customErrors.NewUnauthorizedError("invalid credentials", bcrypt.ErrMismatchedHashAndPassword)
	if !utils.CompareErrors(err, expectedErr) {
		t.Errorf("Authenticate() error = %v, want %v", err, expectedErr)
	}

	if len(mockSessionService.savedSessions) != 0 {
		t.Error("Authenticate() should not save session when credentials are invalid")
	}

	if session.UserID != 0 {
		t.Errorf("Authenticate() session UserID = %d, want 0", session.UserID)
	}
}

func TestAuthenticate_InvalidPasswordHash_ReturnsInternalServerError(t *testing.T) {
	avatar := enum.Avatar1
	mockUserService := &MockUserService{
		user: &model.User{
			ID:        userID,
			Username:  username,
			Password:  "not-a-bcrypt-hash",
			CreatedAt: time.Now().UTC(),
			Avatar:    &avatar,
			Language:  enum.English,
			AppTheme:  enum.Light,
		},
	}
	mockSessionService := &MockSessionService{}
	service := NewAuthService(mockSessionService, mockUserService)

	session := model.NewSession()
	err := service.Authenticate(session, username, "any-password")

	if err == nil {
		t.Fatal("Authenticate() error = nil, want non-nil")
	}

	appErr, ok := errors.AsType[customErrors.AppError](err)
	if !ok {
		t.Fatalf("Authenticate() error type = %T, want *AppError", err)
	}

	if appErr.HTTPStatus() != http.StatusInternalServerError {
		t.Errorf("Authenticate() status = %d, want %d", appErr.HTTPStatus(), http.StatusInternalServerError)
	}

	if appErr.Error() != "an error occurred while checking credentials" {
		t.Errorf("Authenticate() message = %q, want %q", appErr.Error(), "an error occurred while checking credentials")
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		t.Error("Authenticate() expected non-mismatch bcrypt error for invalid hash")
	}

	if len(mockSessionService.savedSessions) != 0 {
		t.Error("Authenticate() should not save session when password hash is invalid")
	}
}

func TestLogout_RemovesSession(t *testing.T) {
	mockUserService := &MockUserService{}
	mockSessionService := &MockSessionService{}
	service := NewAuthService(mockSessionService, mockUserService)

	session := model.NewSession()
	err := service.Logout(session)

	if err != nil {
		t.Fatalf("Logout() error = %v, want nil", err)
	}

	if len(mockSessionService.removedSessions) != 1 {
		t.Fatalf("Logout() remove calls = %d, want 1", len(mockSessionService.removedSessions))
	}

	if mockSessionService.removedSessions[0] != session {
		t.Error("Logout() removed session pointer does not match input session")
	}
}

func TestLogout_RepositoryError(t *testing.T) {
	mockUserService := &MockUserService{}
	mockSessionService := &MockSessionService{
		removeErr: customErrors.NewInternalError("failed to remove session", nil),
	}
	service := NewAuthService(mockSessionService, mockUserService)

	session := model.NewSession()
	err := service.Logout(session)

	if err == nil {
		t.Fatal("Logout() error = nil, want non-nil")
	}

	if !utils.CompareErrors(err, mockSessionService.removeErr) {
		t.Errorf("Logout() error = %v, want %v", err, mockSessionService.removeErr)
	}

	if len(mockSessionService.removedSessions) != 0 {
		t.Fatalf("Logout() shouldn't remove session on repository error")
	}
}
