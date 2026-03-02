package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
)

type SessionRepositoryInterface interface {
	// Define session-related data operations here
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// GetByID retrieves a session by its ID.
func (r *SessionRepository) GetByID(id int64) (*model.Session, error) {
	row := r.db.QueryRow("SELECT * FROM session WHERE id = ?", id)

	s := &model.Session{}

	err := row.Scan(
		&s.ID,
		&s.CreatedAt,
		&s.LastActivity,
		&s.IPAddress,
		&s.UserAgent,
		&s.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewEntityNotFoundError("Session", strconv.FormatInt(id, 10), err)
		}
		return nil, customErrors.NewInternalServerError("Failed to fetch session by ID", err)
	}

	return s, nil
}

// Write inserts a new session or updates an existing one based on the session ID.
func (r *SessionRepository) Write(s *model.Session) error {
	_, err := r.db.Exec(
		`INSERT INTO session (id, created_at, last_activity, ip_address, user_agent, user_id)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   user_id = excluded.user_id,
		   last_activity = excluded.last_activity,
		   ip_address = excluded.ip_address,
		   user_agent = excluded.user_agent`,
		s.ID, s.CreatedAt, s.LastActivity, s.IPAddress, s.UserAgent, s.UserID,
	)
	return customErrors.NewInternalServerError("Failed to write in session", err)
}

func (r *SessionRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM session WHERE id = ?", id)
	return customErrors.NewInternalServerError("Failed to delete session", err)
}

// CleanUp removes sessions that have been inactive for longer than the specified expiration duration.
// It returns the number of sessions that were removed.
func (r *SessionRepository) CleanUp(expiration time.Duration) int64 {
	result, err := r.db.Exec("DELETE FROM session WHERE last_activity < ?", time.Now().Add(-expiration).UTC())
	if err != nil {
		return 0
	}

	removedRows, _ := result.RowsAffected()
	return removedRows
}
