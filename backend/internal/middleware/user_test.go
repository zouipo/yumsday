package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/model"
)

type mockUserService struct {
	getByIDCalls int
	getByIDErr   error
}

func (m *mockUserService) GetAll() ([]model.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetByID(id int64) (*model.User, error) {
	m.getByIDCalls++
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return &model.User{ID: id, Username: "test"}, nil
}

func (m *mockUserService) GetByUsername(username string) (*model.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserService) Create(user *model.User) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockUserService) Update(user *model.User) error {
	return errors.New("not implemented")
}

func (m *mockUserService) UpdateAdminRole(userID int64, role bool) error {
	return errors.New("not implemented")
}

func (m *mockUserService) UpdatePassword(userID int64, oldPassword string, newPassword string) error {
	return errors.New("not implemented")
}

func (m *mockUserService) Delete(id int64) error {
	return errors.New("not implemented")
}

// Test unauthenticated
func TestUserInjector_unauthenticated_nonLogin(t *testing.T) {
	mockService := &mockUserService{}
	mw := UserInjector(mockService)

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	session := model.NewSession()
	session.UserID = 0

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	mw(next).ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d instead of %d", http.StatusUnauthorized, w.Code)
	}

	if handlerCalled {
		t.Fatal("expected handler not to be called")
	}

	if mockService.getByIDCalls != 0 {
		t.Fatalf("expected GetByID not to be called, got %d", mockService.getByIDCalls)
	}
}

func TestUserInjector_unauthenticated_login(t *testing.T) {
	mockService := &mockUserService{}
	mw := UserInjector(mockService)

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	session := model.NewSession()
	session.UserID = 0

	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	mw(next).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d instead of %d", http.StatusOK, w.Code)
	}

	if !handlerCalled {
		t.Fatal("expected handler to be called")
	}

	if mockService.getByIDCalls != 0 {
		t.Fatalf("expected GetByID not to be called, got %d", mockService.getByIDCalls)
	}
}

// Tests authenticated
func TestUserInjector_authenticated_login(t *testing.T) {
	mockService := &mockUserService{}
	mw := UserInjector(mockService)

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	session := model.NewSession()
	session.UserID = 1

	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	mw(next).ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Fatalf("expected status %d instead of %d", http.StatusFound, w.Code)
	}

	if got := w.Header().Get("Location"); got != "/" {
		t.Fatalf("expected redirect location %q instead of %q", "/", got)
	}

	if handlerCalled {
		t.Fatal("expected handler not to be called")
	}

	if mockService.getByIDCalls != 0 {
		t.Fatalf("expected GetByID not to be called, got %d", mockService.getByIDCalls)
	}
}

func TestUserInjector_authenticated_nonLogin(t *testing.T) {
	mockService := &mockUserService{}
	mw := UserInjector(mockService)

	var userInCtx *model.User
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInCtx = r.Context().Value(ctx.UserCtxKey{}).(*model.User)
		w.WriteHeader(http.StatusOK)
	})

	session := model.NewSession()
	session.UserID = 1

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	mw(next).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d instead of %d", http.StatusOK, w.Code)
	}

	if mockService.getByIDCalls != 1 {
		t.Fatalf("expected GetByID to be called once, got %d", mockService.getByIDCalls)
	}

	if userInCtx == nil {
		t.Fatal("expected user in context, got nil")
	}

	if userInCtx.ID != 1 {
		t.Fatalf("expected user ID 1 instead of %d", userInCtx.ID)
	}
}

func TestUserInjector_authenticated_getByIDError(t *testing.T) {
	getByIDErr := errors.New("db error")
	mockService := &mockUserService{getByIDErr: getByIDErr}
	mw := UserInjector(mockService)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	session := model.NewSession()
	session.UserID = 1

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r = r.WithContext(context.WithValue(r.Context(), ctx.SessionCtxKey{}, session))
	w := httptest.NewRecorder()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic, got none")
		}
		if recovered != getByIDErr {
			t.Fatalf("expected panic value %v instead of %v", getByIDErr, recovered)
		}
	}()

	mw(next).ServeHTTP(w, r)
}
