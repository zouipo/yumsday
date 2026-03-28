package repository

import (
	"database/sql"
	"errors"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
)

type SessionRepositoryInterface interface {
	GetByID(id string) (*model.Session, error)
	Write(s *model.Session) error
	Delete(id string) error
	CleanUp(expiration time.Duration) int64
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
func (r *SessionRepository) GetByID(id string) (*model.Session, error) {
	row := r.db.QueryRow("SELECT * FROM sessions WHERE id = ?", id)

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
			return nil, customErrors.NewNotFoundError("Session", id, err)
		}
		return nil, customErrors.NewInternalError("Failed to fetch session by ID", err)
	}

	return s, nil
}

// Write inserts a new session or updates an existing one based on the session ID.
func (r *SessionRepository) Write(s *model.Session) error {
	_, err := r.db.Exec(
		`INSERT INTO sessions (id, created_at, last_activity, ip_address, user_agent, user_id)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   user_id = excluded.user_id,
		   last_activity = excluded.last_activity,
		   ip_address = excluded.ip_address,
		   user_agent = excluded.user_agent`,
		s.ID, s.CreatedAt, s.LastActivity, s.IPAddress, s.UserAgent, s.UserID,
	)
	if err != nil {
		return customErrors.NewInternalError("Failed to write in session", err)
	}
	return nil
}

// Delete removes a session by its ID.
//
// NOTE: It does not return an error if the session doesn't exist, since SQLite's DELETE doesn't error on non-existent rows.
func (r *SessionRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	if err != nil {
		return customErrors.NewInternalError("Failed to delete session", err)
	}
	return nil
}

// CleanUp removes sessions that have been inactive for longer than the specified expiration duration.
// It returns the number of sessions that were removed.
func (r *SessionRepository) CleanUp(expiration time.Duration) int64 {
	result, err := r.db.Exec("DELETE FROM sessions WHERE last_activity < ?", time.Now().Add(-expiration).UTC())
	if err != nil {
		return 0
	}

	removedRows, err := result.RowsAffected()
	if err != nil {
		return 0
	}
	return removedRows
}
