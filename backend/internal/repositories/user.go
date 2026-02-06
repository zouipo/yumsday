package repositories

import (
	"database/sql"

	"github.com/zouipo/yumsday/backend/internal/models"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
	GetAll() ([]models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) (int64, error)
	Update(user *models.User) error
	Delete(id int64) error
}

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a new UserRepository using the provided database.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetAll fetches all users from the database.
func (r *UserRepository) GetAll() ([]models.User, error) {
	users, err := r.fetchUsers("SELECT * FROM user")
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID fetches the user by ID; returns sql.ErrNoRows if not found.
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user, err := r.fetchUser("SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByUsername fetches the user that matches the provided username; returns sql.ErrNoRows if not found.
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user, err := r.fetchUser("SELECT * FROM user WHERE username = ?", username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Create inserts a new user into the database and returns the inserted ID.
func (r *UserRepository) Create(user *models.User) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO user (username, password, app_admin, created_at, avatar, language, app_theme) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Username,
		user.Password,
		user.AppAdmin,
		user.CreatedAt,
		user.Avatar,
		user.Language,
		user.AppTheme,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update updates an existing user; returns sql.ErrNoRows if no row was affected.
func (r *UserRepository) Update(user *models.User) error {
	result, err := r.db.Exec(
		"UPDATE user SET username = ?, password = ?, app_admin = ?, avatar = ?, language = ?, app_theme = ? WHERE id = ?",
		user.Username,
		user.Password,
		user.AppAdmin,
		user.Avatar,
		user.Language,
		user.AppTheme,
		user.ID,
	)
	if err != nil {
		return err
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if updatedRow == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete removes a user by ID; returns sql.ErrNoRows if no row was deleted.
func (r *UserRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return err
	}

	deletedRow, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if deletedRow == 0 {
		return sql.ErrNoRows
	}

	return nil
}

/*** PRIVATE HELPER METHODS ***/
// fetchUsers executes the provided query and returns a slice of users.
func (r *UserRepository) fetchUsers(query string, args ...any) ([]models.User, error) {
	users := []models.User{}

	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	// Not mandatory
	defer rows.Close()

	for rows.Next() {
		var user models.User
		// // Scan assign each row's column values to a field of the struct User through the pointers.
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.AppAdmin,
			&user.CreatedAt,
			&user.Avatar,
			&user.Language,
			&user.AppTheme,
			&user.LastVisitedGroup,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// rows.Next stops either because it encounters an error (so row.Err() != nil),
	// either because it has reach the end of the rows.
	// rows.Err() automatically closes rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// fetchUser executes the provided query and returns a single user or sql.ErrNoRows.
func (r *UserRepository) fetchUser(query string, args ...any) (*models.User, error) {
	user := &models.User{}

	row := r.db.QueryRow(query, args...)

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.AppAdmin,
		&user.CreatedAt,
		&user.Avatar,
		&user.Language,
		&user.AppTheme,
		&user.LastVisitedGroup,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}

	return user, nil
}
