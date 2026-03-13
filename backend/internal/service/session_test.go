package service

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
)

var (
	cookieName = "session_id"
	expiration = 1 * time.Hour
)

type MockSessionRepository struct {
	mu         sync.RWMutex
	sessions   map[string]*model.Session
	getByIDErr error
	writeErr   error
	deleteErr  error
	cleanUpErr error
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions: make(map[string]*model.Session),
	}
}

/*** SESSION REPOSITORY IMPLEMENTATION ***/

func (m *MockSessionRepository) GetByID(id string) (*model.Session, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[id]
	if !exists {
		return nil, customErrors.NewEntityNotFoundError("Session", id, nil)
	}

	return session, nil
}

func (m *MockSessionRepository) Write(s *model.Session) error {
	if m.writeErr != nil {
		return m.writeErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[s.ID] = s
	return nil
}

func (m *MockSessionRepository) Delete(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, id)
	return nil
}

func (m *MockSessionRepository) CleanUp(exp time.Duration) int64 {
	if m.cleanUpErr != nil {
		return 0
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	removed := int64(0)
	cutoff := time.Now().Add(-exp).UTC()

	for id, session := range m.sessions {
		if session.LastActivity.Before(cutoff) {
			delete(m.sessions, id)
			removed++
		}
	}

	return removed
}

/*** HELPER FUNCTIONS ***/

func (m *MockSessionRepository) addSession(session *model.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[session.ID] = session
}

func (m *MockSessionRepository) hasSession(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.sessions[id]
	return exists
}

func (m *MockSessionRepository) sessionCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.sessions)
}

func (m *MockSessionRepository) getSession(id string) (*model.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, exists := m.sessions[id]
	return s, exists
}

func createTestSession(id string, lastActivity time.Time) *model.Session {
	return &model.Session{
		ID:           id,
		CreatedAt:    time.Now().UTC(),
		LastActivity: lastActivity,
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Agent",
		UserID:       1,
	}
}

/*** TEST CONSTRUCTOR ***/

func TestNewSessionService(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	expiration := 24 * time.Hour

	service := NewSessionService(mockRepo, cookieName, expiration)

	if service == nil {
		t.Fatal("NewSessionService() returned nil")
	}

	if service.repo == nil {
		t.Error("NewSessionService() repo is nil")
	}

	if service.cookieName != cookieName {
		t.Errorf("NewSessionService() cookieName = %s, want %s", service.cookieName, cookieName)
	}

	if service.expiration != expiration {
		t.Errorf("NewSessionService() expiration = %v, want %v", service.expiration, expiration)
	}
}

func TestNewSessionService_CleansUpExpiredSessions(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	testExpiration := 1 * time.Hour

	// expired sessions
	expiredSession1 := createTestSession("expired-1", time.Now().UTC().Add(-2*time.Hour))
	expiredSession2 := createTestSession("expired-2", time.Now().UTC().Add(-3*time.Hour))
	mockRepo.addSession(expiredSession1)
	mockRepo.addSession(expiredSession2)

	// valid sessions
	validSession1 := createTestSession("valid-1", time.Now().UTC().Add(-30*time.Minute))
	validSession2 := createTestSession("valid-2", time.Now().UTC())
	mockRepo.addSession(validSession1)
	mockRepo.addSession(validSession2)

	_ = NewSessionService(mockRepo, cookieName, testExpiration)

	// Pause to allow the cleanup goroutine to run
	time.Sleep(100 * time.Millisecond)

	// Expired sessions should be removed
	if mockRepo.hasSession("expired-1") {
		t.Error("NewSessionService() did not clean up expired-1 session")
	}
	if mockRepo.hasSession("expired-2") {
		t.Error("NewSessionService() did not clean up expired-2 session")
	}

	// Valid sessions should remain
	if !mockRepo.hasSession("valid-1") {
		t.Error("NewSessionService() incorrectly removed valid-1 session")
	}
	if !mockRepo.hasSession("valid-2") {
		t.Error("NewSessionService() incorrectly removed valid-2 session")
	}

	if mockRepo.sessionCount() != 2 {
		t.Errorf("NewSessionService() expected 2 sessions after cleanup, got %d", mockRepo.sessionCount())
	}
}

func TestNewSessionService_CleansUpError(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	mockRepo.cleanUpErr = customErrors.NewInternalServerError("Failed to clean up sessions", nil)
	testExpiration := 1 * time.Hour

	expiredSession1 := createTestSession("expired-1", time.Now().UTC().Add(-2*time.Hour))
	mockRepo.addSession(expiredSession1)
	validSession1 := createTestSession("valid-1", time.Now().UTC().Add(-30*time.Minute))
	mockRepo.addSession(validSession1)

	_ = NewSessionService(mockRepo, cookieName, testExpiration)

	time.Sleep(100 * time.Millisecond)

	if mockRepo.sessionCount() != 2 {
		t.Errorf("NewSessionService() expected 2 sessions after cleanup error, got %d", mockRepo.sessionCount())
	}
}

/*** GET TESTS ***/

func TestGetSession_NoCookie_ReturnsNewSession(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	session := service.GetSession(req)

	if session == nil {
		t.Fatal("GetSession() returned nil")
	}

	if session.ID == "" {
		t.Error("GetSession() returned session with empty ID")
	}

	if time.Since(session.CreatedAt) > 1*time.Second {
		t.Error("GetSession() session CreatedAt is not recent")
	}

	if time.Since(session.LastActivity) > 1*time.Second {
		t.Error("GetSession() session LastActivity is not recent")
	}
}

func TestGetSession_ValidCookie_ReturnsExistingSession(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	sessionID := "test-session-123"
	existingSession := createTestSession(sessionID, time.Now().UTC())
	mockRepo.addSession(existingSession)

	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: sessionID,
	})

	session := service.GetSession(req)

	if session == nil {
		t.Fatal("GetSession() returned nil")
	}

	if session.ID != sessionID {
		t.Errorf("GetSession() session ID = %s, want %s", session.ID, sessionID)
	}

	if session.UserID != existingSession.UserID {
		t.Errorf("GetSession() session UserID = %d, want %d", session.UserID, existingSession.UserID)
	}

	if session.LastActivity != existingSession.LastActivity {
		t.Error("GetSession() session LastActivity does not match existing session")
	}
}

func TestGetSession_SessionNotFound_ReturnsNewSession(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	nonExistingSessionID := "non-existent-session"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: nonExistingSessionID,
	})

	session := service.GetSession(req)

	if session == nil {
		t.Fatal("GetSession() returned nil")
	}

	if session.ID == nonExistingSessionID {
		t.Error("GetSession() should return new session when existing session not found")
	}

	if time.Since(session.CreatedAt) > 1*time.Second {
		t.Error("GetSession() new session CreatedAt is not recent")
	}
}

func TestGetSession_ExpiredSession_ReturnsNewSession(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	sessionID := "expired-session-123"

	// Create a session with LastActivity older than expiration
	expiredSession := createTestSession(sessionID, time.Now().UTC().Add(-2*time.Hour))
	mockRepo.addSession(expiredSession)

	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: sessionID,
	})

	session := service.GetSession(req)

	if session == nil {
		t.Fatal("GetSession() returned nil")
	}

	if session.ID == sessionID {
		t.Error("GetSession() should return new session when existing session is expired")
	}

	if time.Since(session.CreatedAt) > 1*time.Second {
		t.Error("GetSession() new session CreatedAt is not recent")
	}

	// Verify the expired session was deleted from the repository
	_, exists := mockRepo.getSession(sessionID)
	if exists {
		t.Error("GetSession() should delete expired session from repository")
	}
}

func TestGetSession_RepositoryError_ReturnsNewSession(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	mockRepo.getByIDErr = customErrors.NewInternalServerError("Failed to fetch session by ID", nil)

	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: "session-id",
	})

	session := service.GetSession(req)

	if session == nil {
		t.Fatal("GetSession() returned nil")
	}

	if session.ID == "session-id" {
		t.Error("GetSession() should return new session when repository error occurs")
	}

	// Should return a new session when repository error occurs
	if time.Since(session.CreatedAt) > 1*time.Second {
		t.Error("GetSession() new session CreatedAt is not recent")
	}
}

/*** SAVE TESTS ***/

func TestSave_Success(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	oldLastActivity := time.Now().UTC().Add(-10 * time.Minute)
	session := createTestSession("test-session", oldLastActivity)

	// Wait to ensure time difference in LastActivity
	time.Sleep(10 * time.Millisecond)

	service.Save(session)

	savedSession, exists := mockRepo.getSession(session.ID)
	if !exists {
		t.Fatal("Save() did not save session to repository")
	}

	if savedSession.ID != session.ID {
		t.Errorf("Save() session ID = %s, want %s", savedSession.ID, session.ID)
	}

	if !savedSession.LastActivity.After(oldLastActivity) {
		t.Error("Save() did not update LastActivity")
	}
	// Verify LastActivity was updated to recent time
	if time.Since(savedSession.LastActivity) > 1*time.Second {
		t.Error("Save() LastActivity is not recent")
	}
}

func TestSave_WithRepositoryError(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	mockRepo.writeErr = customErrors.NewInternalServerError("Failed to write in session", nil)

	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	session := createTestSession("test-session", time.Now().UTC())

	service.Save(session)

	// LastActivity should still be updated
	if time.Since(session.LastActivity) > 1*time.Second {
		t.Error("Save() did not update LastActivity even with repository error")
	}
}

/*** REMOVE TESTS ***/

func TestRemove_Success(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	sessionID := "test-session-123"
	session := createTestSession(sessionID, time.Now().UTC())
	mockRepo.addSession(session)

	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	if !mockRepo.hasSession(sessionID) {
		t.Fatal("Setup failed: session does not exist before removal")
	}

	service.Remove(session)

	if mockRepo.hasSession(sessionID) {
		t.Error("Remove() did not delete session from repository")
	}
}

func TestRemove_NonExistentSession(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	sessionID := "test-session-123"
	session := createTestSession(sessionID, time.Now().UTC())
	mockRepo.addSession(session)

	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	falseSession := createTestSession("non-existent", time.Now().UTC())

	service.Remove(falseSession)

	if mockRepo.sessionCount() != 1 {
		t.Error("Remove() should not have add or delete sessions")
	}
}

func TestRemove_Error(t *testing.T) {
	mockRepo := NewMockSessionRepository()
	mockRepo.deleteErr = customErrors.NewInternalServerError("Failed to delete session", nil)

	sessionID := "test-session-123"
	session := createTestSession(sessionID, time.Now().UTC())
	mockRepo.addSession(session)

	service := &SessionService{
		repo:       mockRepo,
		cookieName: cookieName,
		expiration: expiration,
	}

	service.Remove(session)

	if mockRepo.sessionCount() != 1 {
		t.Error("Remove() should not have add or delete sessions")
	}
}
