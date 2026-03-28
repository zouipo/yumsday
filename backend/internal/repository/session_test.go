package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	now            = time.Now().UTC()
	oneHourAgo     = now.Add(-1 * time.Hour)
	twoDaysAgo     = now.Add(-48 * time.Hour)
	invalidSession = "invalid-session-id"

	// UserID will be set during test setup after inserting test users into the database
	expectedSessions = []model.Session{
		{
			ID:           "session-id-1",
			CreatedAt:    twoDaysAgo,
			LastActivity: oneHourAgo,
			IPAddress:    "192.168.1.1",
			UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			UserID:       0,
		},
		{
			ID:           "session-id-2",
			CreatedAt:    oneHourAgo,
			LastActivity: now,
			IPAddress:    "192.168.1.2",
			UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			UserID:       0,
		},
		{
			ID:           "session-id-3",
			CreatedAt:    twoDaysAgo,
			LastActivity: twoDaysAgo,
			IPAddress:    "192.168.1.3",
			UserAgent:    "Mozilla/5.0 (X11; Linux x86_64)",
			UserID:       0,
		},
	}
)

// compareSessions compares two Session objects to check if they are equivalent.
func compareSessions(actual, expected *model.Session) error {
	if actual.ID != expected.ID {
		return fmt.Errorf("ID = %s instead of %s", actual.ID, expected.ID)
	}

	if !utils.TimesApproximatelyEqual(actual.CreatedAt, expected.CreatedAt, time.Minute) {
		return fmt.Errorf("CreatedAt = %v instead of around %v (±1min)", actual.CreatedAt, expected.CreatedAt)
	}

	if !utils.TimesApproximatelyEqual(actual.LastActivity, expected.LastActivity, time.Minute) {
		return fmt.Errorf("LastActivity = %v instead of around %v (±1min)", actual.LastActivity, expected.LastActivity)
	}

	if actual.IPAddress != expected.IPAddress {
		return fmt.Errorf("IPAddress = %s instead of %s", actual.IPAddress, expected.IPAddress)
	}

	if actual.UserAgent != expected.UserAgent {
		return fmt.Errorf("UserAgent = %s instead of %s", actual.UserAgent, expected.UserAgent)
	}

	if actual.UserID != expected.UserID {
		return fmt.Errorf("UserID = %d instead of %d", actual.UserID, expected.UserID)
	}

	return nil
}

// setupTestDB initializes an in-memory SQLite database with test data for testing.
func setupSessionTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	migrationsFS := os.DirFS("../../data/migrations")
	err = migration.Migrate(db, migrationsFS)
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	testUsers := []struct {
		username string
		password string
	}{
		{"testuser1", "validpassword123"},
		{"testuser2", "validpassword456"},
	}

	var userIDs []int64
	for _, user := range testUsers {
		res, err := db.Exec(
			`INSERT INTO users (username, password, app_admin, created_at, language, app_theme)
			VALUES (?, ?, ?, ?, ?, ?);`,
			user.username,
			user.password,
			false,
			now,
			"en",
			"light",
		)
		if err != nil {
			t.Fatalf("failed to insert test user '%s': %v", user.username, err)
		}

		userID, err := res.LastInsertId()
		if err != nil {
			t.Fatalf("failed to get last insert ID for user '%s': %v", user.username, err)
		}
		// Get user ID from the database for the foreign key reference in sessions
		userIDs = append(userIDs, userID)
	}

	expectedSessions[0].UserID = userIDs[0]
	expectedSessions[1].UserID = userIDs[1]
	expectedSessions[2].UserID = userIDs[0]

	for _, session := range expectedSessions {
		_, err := db.Exec(
			`INSERT INTO sessions (id, created_at, last_activity, ip_address, user_agent, user_id)
			VALUES (?, ?, ?, ?, ?, ?);`,
			session.ID,
			session.CreatedAt,
			session.LastActivity,
			session.IPAddress,
			session.UserAgent,
			session.UserID,
		)
		if err != nil {
			t.Fatalf("failed to insert test session '%s': %v", session.ID, err)
		}
	}

	return db
}

// teardownSessionTestDB closes the database connection.
func teardownSessionTestDB(db *sql.DB) {
	db.Close()
}

/*** TEST CONSTRUCTOR ***/

func TestNewSessionRepository(t *testing.T) {
	db := setupSessionTestDB(t)
	defer teardownSessionTestDB(db)

	repo := NewSessionRepository(db)

	if repo == nil {
		t.Fatal("expected non-nil SessionRepository")
	}

	if repo.db == nil {
		t.Fatal("expected non-nil database connection")
	}
}

/*** READ OPERATIONS TESTS ***/

func TestGetBySessionID(t *testing.T) {
	db := setupSessionTestDB(t)
	defer teardownSessionTestDB(db)

	repo := NewSessionRepository(db)

	tests := []struct {
		name      string
		sessionID string
		wantErr   error
		expected  *model.Session
	}{
		{
			name:      "existing session 1",
			sessionID: expectedSessions[0].ID,
			wantErr:   nil,
			expected:  &expectedSessions[0],
		},
		{
			name:      "non-existing session",
			sessionID: invalidSession,
			wantErr:   customErrors.NewNotFoundError("Session", invalidSession, sql.ErrNoRows),
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := repo.GetByID(tt.sessionID)

			if tt.wantErr != nil {
				if !utils.CompareErrors(err, tt.wantErr) {
					t.Errorf("GetByID() error = '%v' instead of '%v'", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if err := compareSessions(session, tt.expected); err != nil {
				t.Errorf("GetByID() returned session does not match expected session: %v", err.Error())
			}
		})
	}
}

/*** WRITE OPERATIONS TESTS ***/

func TestWriteSession(t *testing.T) {
	db := setupSessionTestDB(t)
	defer teardownSessionTestDB(db)

	repo := NewSessionRepository(db)

	tests := []struct {
		name    string
		session *model.Session
	}{
		{
			name: "insert new session",
			session: &model.Session{
				ID:           "new-session-id",
				CreatedAt:    now,
				LastActivity: now,
				IPAddress:    "10.0.0.1",
				UserAgent:    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)",
				UserID:       expectedSessions[0].UserID,
			},
		},
		{
			name: "update existing session",
			session: &model.Session{
				ID:           expectedSessions[0].ID,
				CreatedAt:    expectedSessions[0].CreatedAt,
				LastActivity: now,
				IPAddress:    "10.0.0.2",
				UserAgent:    "Updated User Agent",
				UserID:       expectedSessions[1].UserID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Write(tt.session)

			if err != nil {
				t.Fatalf("Write() unexpected error = %v", err)
			}

			// Verify the session was actually written
			writtenSession, err := repo.GetByID(tt.session.ID)
			if err != nil {
				t.Fatalf("failed to fetch written session: %v", err)
			}

			if err := compareSessions(writtenSession, tt.session); err != nil {
				t.Errorf("Actual written session does not match expected session: %v", err.Error())
			}
		})
	}
}

/*** DELETE OPERATIONS TESTS ***/

func TestDeleteSession(t *testing.T) {
	db := setupSessionTestDB(t)
	defer teardownSessionTestDB(db)

	repo := NewSessionRepository(db)

	tests := []struct {
		name      string
		sessionID string
	}{
		{
			name:      "delete existing session",
			sessionID: expectedSessions[0].ID,
		},
		{
			name:      "delete non-existing session",
			sessionID: invalidSession,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.sessionID)

			if err != nil {
				t.Fatalf("Delete() unexpected error = %v", err)
			}

			// Verify the session was actually deleted
			_, err = repo.GetByID(tt.sessionID)
			if err == nil {
				t.Error("session still exists after deletion")
			}
		})
	}
}

/*** CLEANUP OPERATIONS TESTS ***/

func TestCleanUp(t *testing.T) {
	tests := []struct {
		name            string
		expiration      time.Duration
		expectedRemoved int64
	}{
		{
			name:            "cleanup sessions older than 30 minutes",
			expiration:      30 * time.Minute,
			expectedRemoved: 2,
		},
		{
			name:            "cleanup sessions older than 3 hours",
			expiration:      3 * time.Hour,
			expectedRemoved: 1,
		},
		{
			name:            "cleanup sessions older than 7 days",
			expiration:      7 * 24 * time.Hour,
			expectedRemoved: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupSessionTestDB(t)
			defer teardownSessionTestDB(testDB)
			testRepo := NewSessionRepository(testDB)

			actualRemoved := testRepo.CleanUp(tt.expiration)

			if actualRemoved != tt.expectedRemoved {
				t.Errorf("CleanUp() removed %d sessions, expected %d", actualRemoved, tt.expectedRemoved)
			}

			// Verify the correct sessions were removed
			cutoffTime := now.Add(-tt.expiration)
			for _, session := range expectedSessions {
				retrievedSession, err := testRepo.GetByID(session.ID)
				notFoundErr := customErrors.NewNotFoundError("Session", session.ID, sql.ErrNoRows)
				if err != nil && !utils.CompareErrors(err, notFoundErr) {
					t.Errorf("Unexpected error while retrieving user session: %v", err)
				}

				if session.LastActivity.Before(cutoffTime) {
					if err == nil {
						t.Errorf("session %s should have been deleted but still exists", session.ID)
					}
				} else {
					if err != nil {
						t.Errorf("session %s should still exist but was deleted", session.ID)
					} else if err := compareSessions(retrievedSession, &session); err != nil {
						t.Errorf("session %s was modified: %v", session.ID, err)
					}
				}
			}
		})
	}
}
