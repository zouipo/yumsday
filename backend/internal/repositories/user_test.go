package repositories

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/models/enums"
)

var (
	yesterday = time.Now().AddDate(0, 0, -1)
	minTime   = yesterday.Add(-1 * time.Minute)
	maxTime   = yesterday.Add(1 * time.Minute)
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Apply migrations using the migration package
	migrationsFS := os.DirFS("../../data/migrations")
	err = migration.Migrate(db, migrationsFS)
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	// Read and execute the test data
	testDataSQL, err := os.ReadFile("../../data/test/test_data.sql")
	if err != nil {
		t.Fatalf("failed to read test data file: %v", err)
	}

	_, err = db.Exec(string(testDataSQL))
	if err != nil {
		t.Fatalf("failed to execute test data: %v", err)
	}

	return db
}

// teardownTestDB closes the database connection.
func teardownTestDB(db *sql.DB) {
	db.Close()
}

func TestNewUserRepository(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	if repo == nil {
		t.Fatal("expected non-nil UserRepository")
	}

	if repo.db == nil {
		t.Fatal("expected non-nil database connection")
	}
}

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	expectedCount := 4
	if len(users) != expectedCount {
		t.Errorf("GetAll() returned %d users, expected %d", len(users), expectedCount)
	}

	if len(users) > 0 {
		user := users[0]
		if user.ID != 1 {
			t.Errorf("expected user ID 1, got %d", user.ID)
		}
		if user.Username != "testuser1" {
			t.Errorf("expected username 'testuser1', got '%s'", user.Username)
		}
		if user.Password != "password123" {
			t.Errorf("expected password 'password123', got '%s'", user.Password)
		}
		if user.AppAdmin {
			t.Errorf("expected testuser1 to not be admin")
		}
		// Check that CreatedAt is approximately yesterday (within ±1 minute)
		if user.CreatedAt.Before(minTime) || user.CreatedAt.After(maxTime) {
			t.Errorf("expected creation date around %v (±1min), got '%v'", yesterday, user.CreatedAt)
		}
		if user.Avatar == nil || *user.Avatar != enums.Avatar1 {
			t.Errorf("expected avatar 'AVATAR_1', got '%v'", user.Avatar)
		}
		if user.Language != enums.English {
			t.Errorf("expected language 'EN', got '%s'", user.Language)
		}
		if user.AppTheme != enums.Light {
			t.Errorf("expected app theme 'LIGHT', got '%s'", user.AppTheme)
		}
	}
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name       string
		id         int64
		wantErr    error
		wantUser   string
		wantAdmin  bool
		wantAvatar *enums.Avatar
		wantLang   enums.Language
		wantTheme  enums.AppTheme
	}{
		{
			name:       "existing user 1",
			id:         1,
			wantErr:    nil,
			wantUser:   "testuser1",
			wantAdmin:  false,
			wantAvatar: func() *enums.Avatar { a := enums.Avatar1; return &a }(),
			wantLang:   enums.English,
			wantTheme:  enums.Light,
		},
		{
			name:       "existing user 2",
			id:         2,
			wantErr:    nil,
			wantUser:   "testuser2",
			wantAdmin:  true,
			wantAvatar: func() *enums.Avatar { a := enums.Avatar2; return &a }(),
			wantLang:   enums.French,
			wantTheme:  enums.Dark,
		},
		{
			name:    "non-existing user",
			id:      999,
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(tt.id)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if user.ID != tt.id {
				t.Errorf("GetByID() ID = %d, want %d", user.ID, tt.id)
			}
			if user.Username != tt.wantUser {
				t.Errorf("GetByID() Username = %s, want %s", user.Username, tt.wantUser)
			}
			if user.Password == "" {
				t.Errorf("GetByID() Password should not be empty")
			}
			if user.AppAdmin != tt.wantAdmin {
				t.Errorf("GetByID() AppAdmin = %v, want %v", user.AppAdmin, tt.wantAdmin)
			}
			yesterday := time.Now().AddDate(0, 0, -1)
			minTime := yesterday.Add(-1 * time.Minute)
			maxTime := yesterday.Add(1 * time.Minute)
			if user.CreatedAt.Before(minTime) || user.CreatedAt.After(maxTime) {
				t.Errorf("GetByID() CreatedAt = %v, expected around %v (±1min)", user.CreatedAt, yesterday)
			}
			if tt.wantAvatar != nil && (user.Avatar == nil || *user.Avatar != *tt.wantAvatar) {
				t.Errorf("GetByID() Avatar = %v, want %v", user.Avatar, tt.wantAvatar)
			}
			if user.Language != tt.wantLang {
				t.Errorf("GetByID() Language = %s, want %s", user.Language, tt.wantLang)
			}
			if user.AppTheme != tt.wantTheme {
				t.Errorf("GetByID() AppTheme = %s, want %s", user.AppTheme, tt.wantTheme)
			}
		})
	}
}

func TestGetByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name      string
		username  string
		wantErr   error
		wantID    int64
		wantAdmin bool
	}{
		{
			name:      "existing user testuser1",
			username:  "testuser1",
			wantErr:   nil,
			wantID:    1,
			wantAdmin: false,
		},
		{
			name:      "existing user testuser2",
			username:  "testuser2",
			wantErr:   nil,
			wantID:    2,
			wantAdmin: true,
		},
		{
			name:     "non-existing user",
			username: "nonexistent",
			wantErr:  sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByUsername(tt.username)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByUsername() unexpected error = %v", err)
			}

			if user.ID != tt.wantID {
				t.Errorf("GetByUsername() ID = %d, want %d", user.ID, tt.wantID)
			}
			if user.Username != tt.username {
				t.Errorf("GetByUsername() Username = %s, want %s", user.Username, tt.username)
			}
			if user.Password == "" {
				t.Errorf("GetByUsername() Password should not be empty")
			}
			if user.AppAdmin != tt.wantAdmin {
				t.Errorf("GetByUsername() AppAdmin = %v, want %v", user.AppAdmin, tt.wantAdmin)
			}
			yesterday := time.Now().AddDate(0, 0, -1)
			minTime := yesterday.Add(-1 * time.Minute)
			maxTime := yesterday.Add(1 * time.Minute)
			if user.CreatedAt.Before(minTime) || user.CreatedAt.After(maxTime) {
				t.Errorf("GetByUsername() CreatedAt = %v, expected around %v (±1min)", user.CreatedAt, yesterday)
			}
			if user.Avatar == nil {
				t.Errorf("GetByUsername() Avatar should not be nil for test users 1 and 2")
			}
			if user.Language == "" {
				t.Errorf("GetByUsername() Language should not be empty")
			}
			if user.AppTheme == "" {
				t.Errorf("GetByUsername() AppTheme should not be empty")
			}
		})
	}
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "create new user",
			user: &models.User{
				Username:  "newuser",
				Password:  "newpass123",
				AppAdmin:  false,
				CreatedAt: time.Now(),
				Avatar:    func() *enums.Avatar { a := enums.Avatar1; return &a }(),
				Language:  enums.English,
				AppTheme:  enums.Light,
			},
			wantErr: false,
		},
		{
			name: "create admin user",
			user: &models.User{
				Username:  "adminuser",
				Password:  "adminpass456",
				AppAdmin:  true,
				CreatedAt: time.Now(),
				Avatar:    func() *enums.Avatar { a := enums.Avatar3; return &a }(),
				Language:  enums.French,
				AppTheme:  enums.Dark,
			},
			wantErr: false,
		},
		{
			name: "create duplicate username",
			user: &models.User{
				Username:  "testuser1", // Already exists
				Password:  "password",
				AppAdmin:  false,
				CreatedAt: time.Now(),
				Avatar:    func() *enums.Avatar { a := enums.Avatar1; return &a }(),
				Language:  enums.English,
				AppTheme:  enums.Light,
			},
			wantErr: true, // Should fail due to unique constraint
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(tt.user)

			if tt.wantErr {
				if err == nil {
					t.Error("Create() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() unexpected error = %v", err)
			}

			if id <= 0 {
				t.Errorf("Create() returned invalid ID = %d", id)
			}

			// Verify the user was actually created
			createdUser, err := repo.GetByID(id)
			if err != nil {
				t.Fatalf("failed to fetch created user: %v", err)
			}

			if createdUser.Username != tt.user.Username {
				t.Errorf("created user Username = %s, want %s", createdUser.Username, tt.user.Username)
			}
			if createdUser.Password != tt.user.Password {
				t.Errorf("created user Password = %s, want %s", createdUser.Password, tt.user.Password)
			}
			if createdUser.AppAdmin != tt.user.AppAdmin {
				t.Errorf("created user AppAdmin = %v, want %v", createdUser.AppAdmin, tt.user.AppAdmin)
			}
			if tt.user.Avatar != nil && (createdUser.Avatar == nil || *createdUser.Avatar != *tt.user.Avatar) {
				t.Errorf("created user Avatar = %v, want %v", createdUser.Avatar, tt.user.Avatar)
			}
			if createdUser.Language != tt.user.Language {
				t.Errorf("created user Language = %s, want %s", createdUser.Language, tt.user.Language)
			}
			if createdUser.AppTheme != tt.user.AppTheme {
				t.Errorf("created user AppTheme = %s, want %s", createdUser.AppTheme, tt.user.AppTheme)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		wantErr error
	}{
		{
			name: "update existing user",
			user: &models.User{
				ID:       1,
				Username: "updateduser1",
				Password: "newpassword",
				AppAdmin: true, // Changed from false
				Avatar:   func() *enums.Avatar { a := enums.Avatar2; return &a }(),
				Language: enums.French, // Changed from English
				AppTheme: enums.Dark,   // Changed from Light
			},
			wantErr: nil,
		},
		{
			name: "update non-existing user",
			user: &models.User{
				ID:       999,
				Username: "nonexistent",
				Password: "password",
				AppAdmin: false,
				Avatar:   func() *enums.Avatar { a := enums.Avatar1; return &a }(),
				Language: enums.English,
				AppTheme: enums.Light,
			},
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.user)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Update() unexpected error = %v", err)
			}

			// Verify the user was actually updated
			updatedUser, err := repo.GetByID(tt.user.ID)
			if err != nil {
				t.Fatalf("failed to fetch updated user: %v", err)
			}

			if updatedUser.Username != tt.user.Username {
				t.Errorf("updated user Username = %s, want %s", updatedUser.Username, tt.user.Username)
			}
			if updatedUser.Password != tt.user.Password {
				t.Errorf("updated user Password = %s, want %s", updatedUser.Password, tt.user.Password)
			}
			if updatedUser.AppAdmin != tt.user.AppAdmin {
				t.Errorf("updated user AppAdmin = %v, want %v", updatedUser.AppAdmin, tt.user.AppAdmin)
			}
			if tt.user.Avatar != nil && (updatedUser.Avatar == nil || *updatedUser.Avatar != *tt.user.Avatar) {
				t.Errorf("updated user Avatar = %v, want %v", updatedUser.Avatar, tt.user.Avatar)
			}
			if updatedUser.Language != tt.user.Language {
				t.Errorf("updated user Language = %s, want %s", updatedUser.Language, tt.user.Language)
			}
			if updatedUser.AppTheme != tt.user.AppTheme {
				t.Errorf("updated user AppTheme = %s, want %s", updatedUser.AppTheme, tt.user.AppTheme)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		id      int64
		wantErr error
	}{
		{
			name:    "delete existing user",
			id:      1,
			wantErr: nil,
		},
		{
			name:    "delete another existing user",
			id:      2,
			wantErr: nil,
		},
		{
			name:    "delete non-existing user",
			id:      999,
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Delete() unexpected error = %v", err)
			}

			// Verify the user was actually deleted
			_, err = repo.GetByID(tt.id)
			if err != sql.ErrNoRows {
				t.Errorf("expected user to be deleted, but GetByID() error = %v", err)
			}
		})
	}
}

func TestDeleteThenGetAll(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	// Get initial count
	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	initialCount := len(users)

	// Delete one user
	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify count decreased
	users, err = repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() after delete error = %v", err)
	}

	if len(users) != initialCount-1 {
		t.Errorf("after delete, GetAll() returned %d users, expected %d", len(users), initialCount-1)
	}

	// Verify the deleted user is not in the list
	for _, user := range users {
		if user.ID == 1 {
			t.Error("deleted user still appears in GetAll() results")
		}
	}
}
