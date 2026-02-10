package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/errors"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/models/enums"
)

var (
	yesterday = time.Now().AddDate(0, 0, -1)

	nextId = int64(5)

	invalidId       = int64(-1)
	invalidUsername = "invalidUsername"

	expectedUsers = []models.User{
		{
			ID:               1,
			Username:         "testuser1",
			Password:         "password123",
			AppAdmin:         false,
			CreatedAt:        yesterday,
			Avatar:           func() *enums.Avatar { a := enums.Avatar1; return &a }(),
			Language:         enums.English,
			AppTheme:         enums.Light,
			LastVisitedGroup: func() *int64 { v := int64(1); return &v }(),
		},
		{
			ID:               2,
			Username:         "testuser2",
			Password:         "password456",
			AppAdmin:         true,
			CreatedAt:        yesterday,
			Avatar:           func() *enums.Avatar { a := enums.Avatar2; return &a }(),
			Language:         enums.French,
			AppTheme:         enums.Dark,
			LastVisitedGroup: func() *int64 { v := int64(1); return &v }(),
		},
		{
			ID:               3,
			Username:         "testuser3",
			Password:         "password789",
			AppAdmin:         false,
			CreatedAt:        yesterday,
			Avatar:           func() *enums.Avatar { a := enums.Avatar3; return &a }(),
			Language:         enums.English,
			AppTheme:         enums.System,
			LastVisitedGroup: func() *int64 { v := int64(2); return &v }(),
		},
		{
			ID:               4,
			Username:         "testuser4",
			Password:         "password000",
			AppAdmin:         false,
			CreatedAt:        yesterday,
			Avatar:           nil,
			Language:         enums.English,
			AppTheme:         enums.Light,
			LastVisitedGroup: nil,
		},
	}
)

func compareListUsers(actual, expected []models.User) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("user list length mismatch: actual is %d instead of %d", len(actual), len(expected))
	}

	// Sort both lists by ID to ensure consistent comparison
	sortList(actual)
	sortList(expected)

	for i := range expected {
		actualUser := actual[i]
		expectedUser := expected[i]

		if err := compareUsers(&actualUser, &expectedUser); err != nil {
			return err
		}
	}

	return nil
}

func compareUsers(actual, expected *models.User) error {
	if actual.ID != expected.ID {
		return fmt.Errorf("ID = %d instead of %d", actual.ID, expected.ID)
	}

	if actual.Username != expected.Username {
		return fmt.Errorf("Username = %s instead of %s", actual.Username, expected.Username)
	}

	if actual.Password != expected.Password {
		return fmt.Errorf("Password = %s instead of %s", actual.Password, expected.Password)
	}

	if actual.AppAdmin != expected.AppAdmin {
		return fmt.Errorf("AppAdmin = %v instead of %v", actual.AppAdmin, expected.AppAdmin)
	}

	// allows a small time difference (±1 minute) to account for timing variations
	minTime := expected.CreatedAt.Add(-1 * time.Minute)
	maxTime := expected.CreatedAt.Add(1 * time.Minute)
	if actual.CreatedAt.Before(minTime) || actual.CreatedAt.After(maxTime) {
		return fmt.Errorf("CreatedAt = %v instead of around %v (±1min)", actual.CreatedAt, expected.CreatedAt)
	}

	// Compare Avatar pointers
	if (actual.Avatar == nil) != (expected.Avatar == nil) {
		return fmt.Errorf("Avatar nil mismatch: got %v instead of %v", actual.Avatar, expected.Avatar)
	} else if actual.Avatar != nil && expected.Avatar != nil && *actual.Avatar != *expected.Avatar {
		return fmt.Errorf("Avatar = %v instead of %v", *actual.Avatar, *expected.Avatar)
	}

	if actual.Language != expected.Language {
		return fmt.Errorf("Language = %s instead of %s", actual.Language, expected.Language)
	}

	if actual.AppTheme != expected.AppTheme {
		return fmt.Errorf("AppTheme = %s instead of %s", actual.AppTheme, expected.AppTheme)
	}

	// Compare LastVisitedGroup pointers
	if (actual.LastVisitedGroup == nil) != (expected.LastVisitedGroup == nil) {
		return fmt.Errorf("LastVisitedGroup nil mismatch: got %v instead of %v", actual.LastVisitedGroup, expected.LastVisitedGroup)
	} else if actual.LastVisitedGroup != nil && expected.LastVisitedGroup != nil && *actual.LastVisitedGroup != *expected.LastVisitedGroup {
		return fmt.Errorf("LastVisitedGroup = %d instead of %d", *actual.LastVisitedGroup, *expected.LastVisitedGroup)
	}
	return nil
}

func sortList(users []models.User) {
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
}

func resetNextId() {
	nextId = int64(5)
}

// compareErrors compares two errors to check if they are equivalent AppErrors.
// It compares the Message, StatusCode, and underlying Err fields.
func compareErrors(actual, expected error) bool {
	if actual == nil && expected == nil {
		return true
	}

	if (actual == nil) != (expected == nil) {
		return false
	}

	// Cast both to *AppError
	actualAppErr, actualIsAppErr := actual.(*customErrors.AppError)
	expectedAppErr, expectedIsAppErr := expected.(*customErrors.AppError)

	if actualIsAppErr && expectedIsAppErr {
		if actualAppErr.Code != expectedAppErr.Code && actualAppErr.Message != expectedAppErr.Message && actualAppErr.StatusCode != expectedAppErr.StatusCode {
			return false
		}

		// Cast both into sqlite3.Error to compare their ExtendedCode if possible
		actualSQLErr, actualIsSQLErr := actualAppErr.Err.(sqlite3.Error)
		expectedSQLErr, expectedIsSQLErr := expectedAppErr.Err.(sqlite3.Error)

		if actualIsSQLErr && expectedIsSQLErr {
			return actualSQLErr.ExtendedCode == expectedSQLErr.ExtendedCode
		}

		// If actual is sqlite3.Error but expected is an error code constant (ErrNoExtended or ErrNo),
		// compare the actual error's ExtendedCode with the expected constant
		if actualIsSQLErr {
			if errNoExt, ok := expectedAppErr.Err.(sqlite3.ErrNoExtended); ok {
				return actualSQLErr.ExtendedCode == errNoExt
			}
			if errNo, ok := expectedAppErr.Err.(sqlite3.ErrNo); ok {
				return actualSQLErr.ExtendedCode == sqlite3.ErrNoExtended(errNo)
			}
		}

		// non-SQLite errors
		return actualAppErr.Err == expectedAppErr.Err
	}

	// If not AppErrors, compare their error messages
	return actual.Error() == expected.Error()
}

// setupTestDB initializes an in-memory SQLite database with test data for testing.
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

/*** TEST CONSTRUCTOR ***/

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

/*** READ OPERATIONS TESTS ***/

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if err := compareListUsers(users, expectedUsers); err != nil {
		t.Error("GetAll() returned users with mismatched fields: " + err.Error())
	}
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name     string
		id       int64
		wantErr  error
		wantUser models.User
	}{
		{
			name:     "existing user 1",
			id:       expectedUsers[0].ID,
			wantErr:  nil,
			wantUser: expectedUsers[0],
		},
		{
			name:     "existing user 2",
			id:       expectedUsers[1].ID,
			wantErr:  nil,
			wantUser: expectedUsers[1],
		},
		{
			name:    "non-existing user",
			id:      invalidId,
			wantErr: customErrors.NewEntityNotFoundError("User", fmt.Sprintf("%d", invalidId), sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(tt.id)

			if tt.wantErr != nil {
				if !compareErrors(err, tt.wantErr) {
					t.Errorf("GetByID() error = '%v' instead of '%v'", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if err := compareUsers(user, &tt.wantUser); err != nil {
				t.Errorf("GetByID() returned user does not match expected user: %v", err.Error())
			}
		})
	}
}

func TestGetByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name     string
		username string
		wantErr  error
		wantUser models.User
	}{
		{
			name:     "existing user 1",
			username: expectedUsers[0].Username,
			wantErr:  nil,
			wantUser: expectedUsers[0],
		},
		{
			name:     "existing user 2",
			username: expectedUsers[1].Username,
			wantErr:  nil,
			wantUser: expectedUsers[1],
		},
		{
			name:     "non-existing user",
			username: invalidUsername,
			wantErr:  customErrors.NewEntityNotFoundError("User", fmt.Sprintf("%s", invalidUsername), sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByUsername(tt.username)

			if tt.wantErr != nil {
				if !compareErrors(err, tt.wantErr) {
					t.Errorf("GetByID() error = '%v' instead of '%v'", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if err := compareUsers(user, &tt.wantUser); err != nil {
				t.Errorf("GetByID() returned user does not match expected user: %v", err.Error())
			}
		})
	}
}

/*** CREATE OPERATIONS TESTS ***/

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	avatar1 := enums.Avatar1
	avatar2 := enums.Avatar2

	tests := []struct {
		name    string
		user    *models.User
		wantErr error
	}{
		{
			name: "create new user",
			user: &models.User{
				Username:  "newuser",
				Password:  "newpassword",
				AppAdmin:  false,
				Avatar:    &avatar1,
				Language:  enums.English,
				AppTheme:  enums.Light,
				CreatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "create admin user",
			user: &models.User{
				Username:  "newuser2",
				Password:  "newpassword",
				AppAdmin:  true,
				Avatar:    &avatar2,
				Language:  enums.French,
				AppTheme:  enums.Dark,
				CreatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name:    "create duplicate username",
			user:    &expectedUsers[0],
			wantErr: customErrors.NewConflictError("User", "already exists", sqlite3.ErrConstraintUnique),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(tt.user)

			if tt.wantErr != nil {
				if !compareErrors(err, tt.wantErr) {
					t.Errorf("Create() error = '%v' instead of '%v'", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() unexpected error = %v", err)
			}

			if id != nextId {
				t.Errorf("Create() returned invalid ID = %d", id)
			}
			nextId += 1

			// Verify the user was actually created
			createdUser, err := repo.GetByID(id)
			if err != nil {
				t.Fatalf("failed to fetch created user: %v", err)
			}

			tt.user.ID = id

			if err := compareUsers(createdUser, tt.user); err != nil {
				t.Errorf("Actual created user does not match expected user: %v", err.Error())
			}
		})
	}

	resetNextId()
}

/*** UPDATE OPERATIONS TESTS ***/

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	avatar1 := enums.Avatar1
	avatar2 := enums.Avatar2

	tests := []struct {
		name    string
		user    *models.User
		wantErr error
	}{
		{
			name: "duplicate username update",
			user: &models.User{
				ID:               expectedUsers[1].ID,
				Username:         expectedUsers[0].Username,
				Password:         expectedUsers[1].Password + "_updated",
				AppAdmin:         true,
				CreatedAt:        expectedUsers[1].CreatedAt,
				Avatar:           &avatar2,
				Language:         enums.French,
				AppTheme:         enums.Dark,
				LastVisitedGroup: expectedUsers[1].LastVisitedGroup,
			},
			wantErr: customErrors.NewConflictError("User", "already exists", sqlite3.ErrConstraintUnique),
		},
		{
			name: "update non-existing user",
			user: &models.User{
				ID:       invalidId,
				Username: "nonexistent",
				Password: "password",
				AppAdmin: false,
				Avatar:   &avatar1,
				Language: enums.English,
				AppTheme: enums.Light,
			},
			wantErr: customErrors.NewEntityNotFoundError("User", strconv.FormatInt(invalidId, 10), nil),
		},
		{
			name: "no field updated",
			user: &models.User{
				ID:               expectedUsers[0].ID,
				Username:         expectedUsers[0].Username,
				Password:         expectedUsers[0].Password,
				AppAdmin:         expectedUsers[0].AppAdmin,
				CreatedAt:        expectedUsers[0].CreatedAt,
				Avatar:           expectedUsers[0].Avatar,
				Language:         expectedUsers[0].Language,
				AppTheme:         expectedUsers[0].AppTheme,
				LastVisitedGroup: expectedUsers[0].LastVisitedGroup,
			},
			wantErr: nil,
		},
		{
			name: "update existing user",
			user: &models.User{
				ID:               expectedUsers[0].ID,
				Username:         expectedUsers[0].Username + "_updated",
				Password:         expectedUsers[0].Password + "_updated",
				AppAdmin:         true,
				CreatedAt:        expectedUsers[0].CreatedAt,
				Avatar:           &avatar2,
				Language:         enums.French,
				AppTheme:         enums.Dark,
				LastVisitedGroup: expectedUsers[0].LastVisitedGroup,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.user)

			if tt.wantErr != nil {
				if !compareErrors(err, tt.wantErr) {
					t.Errorf("Update() error = '%v' instead of '%v'", err, tt.wantErr)
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

			if err := compareUsers(updatedUser, tt.user); err != nil {
				t.Errorf("Actual updated user does not match expected user: %v", err.Error())
			}
		})
	}
}

func TestUpdateAdminRole(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		userID  int64
		role    bool
		wantErr error
	}{
		{
			name:    "update admin role for non-existing user",
			userID:  invalidId,
			role:    true,
			wantErr: customErrors.NewEntityNotFoundError("User", strconv.FormatInt(invalidId, 10), nil),
		},
		{
			name:    "set admin role for existing user",
			userID:  expectedUsers[0].ID,
			role:    true,
			wantErr: nil,
		},
		{
			name:    "clear admin role for existing user",
			userID:  expectedUsers[1].ID,
			role:    false,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateAdminRole(tt.userID, tt.role)

			if tt.wantErr != nil {
				if !compareErrors(err, tt.wantErr) {
					t.Errorf("UpdateAdminRole() error = '%v' instead of '%v'", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("UpdateAdminRole() unexpected error = %v", err)
			}

			// Verify the user's admin role was actually updated
			updatedUser, err := repo.GetByID(tt.userID)
			if err != nil {
				t.Fatalf("failed to fetch updated user: %v", err)
			}

			if updatedUser.AppAdmin != tt.role {
				t.Errorf("user admin role not updated: got %v, want %v", updatedUser.AppAdmin, tt.role)
			}
		})
	}
}

/*** DELETE OPERATIONS TESTS ***/

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
			name:    "delete non-existing user",
			id:      invalidId,
			wantErr: customErrors.NewEntityNotFoundError("User", strconv.FormatInt(invalidId, 10), nil),
		},
		{
			name:    "delete existing user",
			id:      expectedUsers[0].ID,
			wantErr: nil,
		},
		{
			name:    "delete another existing user",
			id:      expectedUsers[1].ID,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)

			if tt.wantErr != nil {
				if !compareErrors(err, tt.wantErr) {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Delete() unexpected error = %v", err)
			}

			// Verify the user was actually deleted
			_, err = repo.GetByID(tt.id)
			if err == nil {
				t.Errorf("expected user to be deleted, but GetByID() returned no error")
			} else if appErr, ok := err.(*customErrors.AppError); !ok || appErr.StatusCode != 404 {
				t.Errorf("expected NotFound error, but GetByID() error = %v", err)
			}
		})
	}
}

func TestDeleteThenGetAll(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	initialCount := len(users)

	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	users, err = repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() after delete error = %v", err)
	}

	if len(users) != initialCount-1 {
		t.Errorf("after delete, GetAll() returned %d users, expected %d", len(users), initialCount-1)
	}

	for _, user := range users {
		if user.ID == 1 {
			t.Error("deleted user still appears in GetAll() results")
		}
	}
}
