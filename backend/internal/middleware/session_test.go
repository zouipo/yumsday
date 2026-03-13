package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	cookieName = "session_id"
	expiration = 1 * time.Hour
)

// mockSessionHandler captures whether it was called and the request it received.
type mockSessionHandler struct {
	called  bool
	request *http.Request
}

// ServeHTTP implements http.Handler for test assertions on middleware behavior.
func (m *mockSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.called = true
	m.request = r
	w.WriteHeader(http.StatusOK)
}

type mockSessionService struct {
	// mu protects access to sessionToReturn and saveCalled;
	// Without a mutex, one goroutine can write while another reads, causing a data race and unreliable behavior.
	mu              sync.Mutex
	sessionToReturn *model.Session
	cookieName      string
	expiration      time.Duration
	saveCalled      int
}

func newMockSessionService(session *model.Session) *mockSessionService {
	return &mockSessionService{
		sessionToReturn: session,
		cookieName:      "session_id",
		expiration:      expiration,
	}
}

func (m *mockSessionService) GetSession(_ *http.Request) *model.Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.sessionToReturn == nil {
		m.sessionToReturn = model.NewSession()
	}
	return m.sessionToReturn
}

func (m *mockSessionService) CookieName() string {
	return m.cookieName
}

func (m *mockSessionService) Expiration() time.Duration {
	return m.expiration
}

func (m *mockSessionService) Save(_ *model.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saveCalled++
}

func (m *mockSessionService) getSaveCalled() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saveCalled
}

/*** HELPERS ***/
func compareCookiesSession(t *testing.T, cookie *http.Cookie, expectedName string, expectedValue string) {
	// If an error is thrown, it will point to the test that called compareCookiesSession instead of this helper function.
	// It makes debug easier.
	t.Helper()

	if cookie == nil {
		t.Fatal("session cookie not set in response")
	}
	if cookie.Name != expectedName {
		t.Errorf("expected cookie Name %q, got %q", expectedName, cookie.Name)
	}
	if cookie.Value != expectedValue {
		t.Errorf("expected cookie Value %q, got %q", expectedValue, cookie.Value)
	}
	if cookie.Domain != "localhost" {
		t.Errorf("expected cookie Domain %q, got %q", "localhost", cookie.Domain)
	}
	if !cookie.HttpOnly {
		t.Error("expected cookie HttpOnly to be true")
	}
	if cookie.Path != "/" {
		t.Errorf("expected cookie Path %q, got %q", "/", cookie.Path)
	}
	if !cookie.Secure {
		t.Error("expected cookie Secure to be true")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected cookie SameSite Lax, got %v", cookie.SameSite)
	}
	if !utils.TimesApproximatelyEqual(cookie.Expires, time.Now().Add(expiration).UTC(), time.Minute) {
		t.Errorf("expected cookie Expires %v, got %v", time.Now().Add(expiration).UTC(), cookie.Expires)
	}
	expectedMaxAge := int(expiration.Seconds())
	if cookie.MaxAge != expectedMaxAge {
		t.Errorf("expected cookie MaxAge %d, got %d", expectedMaxAge, cookie.MaxAge)
	}
}

// TestSessionInjector_InjectsNewSessionIntoContext verifies that the SessionInjector middleware
// correctly injects a new session into the request context when no valid session cookie is present.
func TestSessionInjector_InjectsNewSessionIntoContext(t *testing.T) {
	// No session so GetService will return a new session
	svc := newMockSessionService(nil)
	next := &mockSessionHandler{}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder() // session recorder for capturing response headers

	// Three steps chained together:
	// 1. middleware := SessionInjector(svc)
	// 2. handler := middleware(next)
	// 3. handler.ServeHTTP(rr, req)
	SessionInjector(svc)(next).ServeHTTP(rr, req)

	if !next.called {
		t.Fatal("expected next handler to be called")
	}
	if next.request == nil {
		t.Fatal("expected request to be captured by next handler")
	}

	val := next.request.Context().Value(ctx.SessionCtxKey{})
	if val == nil {
		t.Fatal("expected session in context, got nil")
	}

	capturedSession := val.(*model.Session)
	if capturedSession == nil {
		t.Fatal("session was not injected into context")
	}
	if capturedSession.ID == "" {
		t.Error("injected session has empty ID")
	}
}

// TestSessionInjector_UsesExistingSession verifies that when a valid session cookie is present
// and the session exists in the repository, the existing session is reused.
func TestSessionInjector_UsesExistingSession(t *testing.T) {
	existingSessionID := "existing-session-id"
	existing := &model.Session{
		ID:           existingSessionID,
		CreatedAt:    time.Now().UTC(),
		LastActivity: time.Now().UTC(),
	}
	svc := newMockSessionService(existing)
	next := &mockSessionHandler{}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: existingSessionID})
	rr := httptest.NewRecorder()

	SessionInjector(svc)(next).ServeHTTP(rr, req)

	if !next.called {
		t.Fatal("expected next handler to be called")
	}
	if next.request == nil {
		t.Fatal("expected request to be captured by next handler")
	}

	val := next.request.Context().Value(ctx.SessionCtxKey{})
	if val == nil {
		t.Fatal("expected session in context, got nil")
	}

	capturedSession := val.(*model.Session)
	if capturedSession == nil {
		t.Fatal("session was not injected into context")
	}
	if capturedSession.ID != existingSessionID {
		t.Errorf("expected session ID %q, got %q", existingSessionID, capturedSession.ID)
	}
}

// TestSessionInjector_SetsCookie verifies that the middleware sets the session cookie
// with the correct name, security attributes, and expiry.
func TestSessionInjector_SetsCookie(t *testing.T) {
	svc := newMockSessionService(model.NewSession())
	next := &mockSessionHandler{}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	SessionInjector(svc)(next).ServeHTTP(rr, req)

	if !next.called {
		t.Fatal("expected next handler to be called")
	}

	var sessionCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == cookieName {
			sessionCookie = c
			break
		}
	}

	compareCookiesSession(t, sessionCookie, cookieName, svc.sessionToReturn.ID)
}

// TestSessionInjector_CookieValueMatchesSessionID verifies that the value of the session
// cookie matches the ID of the session injected into the context.
func TestSessionInjector_CookieValueMatchesSessionID(t *testing.T) {
	svc := newMockSessionService(model.NewSession())
	next := &mockSessionHandler{}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	SessionInjector(svc)(next).ServeHTTP(rr, req)

	if !next.called {
		t.Fatal("expected next handler to be called")
	}
	if next.request == nil {
		t.Fatal("expected request to be captured by next handler")
	}

	sessionFromCtx := next.request.Context().Value(ctx.SessionCtxKey{})
	if sessionFromCtx == nil {
		t.Fatal("expected session in context, got nil")
	}

	sessionID := sessionFromCtx.(*model.Session).ID

	var cookieValue string
	for _, c := range rr.Result().Cookies() {
		if c.Name == cookieName {
			cookieValue = c.Value
			compareCookiesSession(t, c, cookieName, sessionID)
			break
		}
	}

	if cookieValue != sessionID {
		t.Errorf("cookie value %q does not match session ID %q", cookieValue, sessionID)
	}
}

// TestSessionInjector_SavesSessionAfterHandler verifies that the session is saved
// after the handler returns when the request path is not /logout.
func TestSessionInjector_SavesSessionAfterHandler(t *testing.T) {
	svc := newMockSessionService(model.NewSession())
	next := &mockSessionHandler{}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	SessionInjector(svc)(next).ServeHTTP(rr, req)

	if !next.called {
		t.Fatal("expected next handler to be called")
	}

	// Allow the save goroutine to complete.
	time.Sleep(50 * time.Millisecond)

	if svc.getSaveCalled() == 0 {
		t.Error("expected session to be saved after handler, but Save was never called")
	}
}

// TestSessionInjector_DoesNotSaveSessionOnLogout verifies that the session is NOT saved
// after the handler returns when the request path is /logout.
func TestSessionInjector_DoesNotSaveSessionOnLogout(t *testing.T) {
	svc := newMockSessionService(model.NewSession())
	next := &mockSessionHandler{}

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rr := httptest.NewRecorder()

	SessionInjector(svc)(next).ServeHTTP(rr, req)

	if !next.called {
		t.Fatal("expected next handler to be called")
	}

	// Wait long enough so a goroutine would have run if it was going to.
	time.Sleep(50 * time.Millisecond)

	if svc.getSaveCalled() > 0 {
		t.Errorf("expected session NOT to be saved on /logout, but Save was called %d time(s)", svc.getSaveCalled())
	}
}
