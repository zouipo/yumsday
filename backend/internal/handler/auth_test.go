package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/dto"
	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
)

var (
	username      = "username"
	password      = "password1234"
	wrongPassword = "wrong-password"
)

type mockAuthService struct {
	authErr      error
	logoutErr    error
	authCalls    int
	logoutCalls  int
	lastSession  *model.Session
	lastUsername string
	lastPassword string
}

func (m *mockAuthService) Authenticate(session *model.Session, username, password string) error {
	m.authCalls++
	m.lastSession = session
	m.lastUsername = username
	m.lastPassword = password
	return m.authErr
}

func (m *mockAuthService) Logout(session *model.Session) error {
	m.logoutCalls++
	m.lastSession = session
	return m.logoutErr
}

func TestNewAuthHandler(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}

	if handler.s != mockService {
		t.Error("handler service does not match provided service")
	}
}

/*** TESTS PostLogin ***/

func TestPostLogin_Success(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)
	session := model.NewSession()

	loginReq := dto.LoginDto{
		Username: username,
		Password: password,
	}
	body, _ := json.Marshal(loginReq)

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	handler.postLogin(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d instead of %d", http.StatusNoContent, w.Code)
	}

	if mockService.authCalls != 1 {
		t.Fatalf("expected auth calls 1 instead of %d", mockService.authCalls)
	}

	if mockService.lastSession != session {
		t.Error("expected same session pointer passed to service")
	}

	if mockService.lastUsername != username {
		t.Errorf("expected username %q instead of %q", username, mockService.lastUsername)
	}

	if mockService.lastPassword != password {
		t.Errorf("expected password %q instead of %q", password, mockService.lastPassword)
	}
}

func TestPostLogin_MissingCredentials(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	loginReq := dto.LoginDto{
		Username: "",
		Password: "",
	}
	body, _ := json.Marshal(loginReq)

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.postLogin(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d instead of %d", http.StatusBadRequest, w.Code)
	}

	if !strings.Contains(w.Body.String(), "missing username or password") {
		t.Errorf("expected error message containing %q instead of %q", "missing username or password", w.Body.String())
	}

	if mockService.authCalls != 0 {
		t.Errorf("expected auth calls 0 instead of %d", mockService.authCalls)
	}
}

func TestPostLogin_AppError(t *testing.T) {
	mockService := &mockAuthService{
		authErr: customErrors.NewUnauthorizedError("invalid credentials", nil),
	}
	handler := NewAuthHandler(mockService)
	session := model.NewSession()

	loginReq := dto.LoginDto{
		Username: username,
		Password: wrongPassword,
	}
	body, _ := json.Marshal(loginReq)

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	handler.postLogin(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d instead of %d", http.StatusUnauthorized, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Unauthorized") {
		t.Errorf("expected unauthorized error in response instead of %q", w.Body.String())
	}
}

func TestPostLogin_GenericError(t *testing.T) {
	mockService := &mockAuthService{
		authErr: customErrors.NewInternalError("an error occurred while checking credentials", nil),
	}
	handler := NewAuthHandler(mockService)
	session := model.NewSession()

	loginReq := dto.LoginDto{
		Username: username,
		Password: password,
	}
	body, _ := json.Marshal(loginReq)

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	handler.postLogin(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d instead of %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "an error occurred while checking credentials") {
		t.Errorf("expected error message containing %q instead of %q", "an error occurred while checking credentials", w.Body.String())
	}
}

/*** TESTS PostLogout ***/

func TestPostLogout_Success(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)
	session := model.NewSession()

	r := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	handler.postLogout(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d instead of %d", http.StatusNoContent, w.Code)
	}

	if mockService.logoutCalls != 1 {
		t.Fatalf("expected logout calls 1 instead of %d", mockService.logoutCalls)
	}

	if mockService.lastSession != session {
		t.Error("expected same session pointer passed to service")
	}
}

func TestPostLogout_Error(t *testing.T) {
	mockService := &mockAuthService{logoutErr: customErrors.NewInternalError("Failed to remove session", nil)}
	handler := NewAuthHandler(mockService)
	session := model.NewSession()

	r := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	handler.postLogout(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d instead of %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Failed to remove session") {
		t.Errorf("expected error message containing %q instead of %q", "Failed to remove session", w.Body.String())
	}
}
