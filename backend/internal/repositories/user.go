package repositories

import (
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	"github.com/mattn/go-sqlite3"
	customErrors "github.com/zouipo/yumsday/backend/internal/errors"

	"github.com/zouipo/yumsday/backend/internal/models"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
	GetAll() ([]models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) (int64, error)
	Update(user *models.User) error
	UpdateAdminRole(userID int64, role bool) error
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
	users, err := r.fetchUsers()
	if err != nil {
		return nil, customErrors.NewInternalServerError("Failed to fetch users", err)
	}

	return users, nil
}

// GetByID fetches the user by ID.
// Returns an AppError if not found.
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user, err := r.fetchUser("id", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), err)
		}
		return nil, customErrors.NewInternalServerError("Failed to fetch user by ID", err)
	}

	return user, nil
}

// GetByUsername fetches the user that matches the provided username.
// Returns an AppError if not found.
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user, err := r.fetchUser("username", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewEntityNotFoundError("User", username, err)
		}
		return nil, err
	}

	return user, nil
}

// Create inserts a new user into the database and returns the inserted ID.
// Returns an AppError if creation fails.
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
		if sqlerr, ok := err.(sqlite3.Error); ok {
			if sqlerr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, customErrors.NewConflictError("User", "already exists", sqlerr)
			}
			return 0, customErrors.NewInternalServerError("Failed to create user", err)
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, customErrors.NewInternalServerError("Failed to retrieve created user", err)
	}

	return id, nil
}

// Update updates an existing user, except the the createdAt field.
// Returns an AppError if update fails.
func (r *UserRepository) Update(user *models.User) error {
	existingUser, err := r.GetByID(user.ID)
	slog.Debug("Existing user before update", "user", existingUser, "error", err)

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
		if sqlerr, ok := err.(sqlite3.Error); ok {
			if sqlerr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return customErrors.NewConflictError("User", "already exists", sqlerr)
			}
			return customErrors.NewInternalServerError("Failed to update user", err)
		}
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return customErrors.NewInternalServerError("Failed to retrieve updated user", err)
	}

	// If no row was updated (because the resource was not found), returns an AppError of type EntityNotFoundError
	if updatedRow == 0 {
		return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(user.ID, 10), err)
	}

	return nil
}

// UpdateAdminRole sets or clears the admin flag for the user with the given ID.
// Returns an AppError if update fails.
func (r *UserRepository) UpdateAdminRole(id int64, role bool) error {
	result, err := r.db.Exec(
		"UPDATE user SET app_admin = ? WHERE id = ?",
		role,
		id,
	)
	if err != nil {
		return customErrors.NewInternalServerError("Failed to update user admin role", err)
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return customErrors.NewInternalServerError("Failed to retrieve updated user", err)
	}

	if updatedRow == 0 {
		return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), err)
	}

	return nil
}

// Delete removes a user by its ID.
// Returns an AppError if deletion fails.
func (r *UserRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return customErrors.NewInternalServerError("Failed to delete user", err)
	}

	deletedRow, err := result.RowsAffected()
	if err != nil {
		return customErrors.NewInternalServerError("Failed to retrieve deleted user", err)
	}

	if deletedRow == 0 {
		return customErrors.NewEntityNotFoundError("User", strconv.FormatInt(id, 10), err)
	}

	return nil
}

/*** PRIVATE HELPER METHODS ***/
// fetchUsers returns a slice of all the users in the database.
func (r *UserRepository) fetchUsers() ([]models.User, error) {
	users := []models.User{}

	rows, err := r.db.Query("SELECT * FROM user")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User
		// Scan assign each row's column values to a field of the struct User through the pointers.
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.AppAdmin,
			&user.CreatedAt,
			&user.Avatar,
			&user.Language,
			&user.AppTheme,
			&user.LastVisitedGroup,
		)

		if err != nil {
			// Close rows before returning to prevent resource leaks.
			rows.Close()
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

// fetchUser executes the provided query and returns a single user.
func (r *UserRepository) fetchUser(column string, args ...any) (*models.User, error) {
	user := &models.User{}

	query := "SELECT * FROM user WHERE " + column + " = ?"

	row := r.db.QueryRow(query, args...)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.AppAdmin,
		&user.CreatedAt,
		&user.Avatar,
		&user.Language,
		&user.AppTheme,
		&user.LastVisitedGroup,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
