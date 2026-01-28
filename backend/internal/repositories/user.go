package repositories

import (
	"database/sql"

	"github.com/zouipo/yumsday/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Fetch all users
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	users, err := r.fetchUsers("SELECT * FROM user")
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Fetch a user by its ID
func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {
	user, err := r.fetchUser("SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Fetch a user by its username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user, err := r.fetchUser("SELECT * FROM user WHERE username = ?", username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Create a new user
func (r *UserRepository) CreateUser(user *models.User) (int64, error) {
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

// Update an existing user
func (r *UserRepository) UpdateUser(user *models.User) error {
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

// Delete a user by its ID
func (r *UserRepository) DeleteUser(id int64) error {
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

/*** PRIVATE METHODS ***/
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
	); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}

	return user, nil
}
