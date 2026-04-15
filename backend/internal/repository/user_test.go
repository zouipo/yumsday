package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

var (
	yesterday = time.Now().AddDate(0, 0, -1)

	invalidId       = int64(-1)
	invalidUsername = "invalidUsername"

	expectedUsers = []model.User{
		{
			Username:           "testuser1",
			Password:           "$2a$12$q7Nm8q9c9g9unKbhjqcWS.Y7tQplxJvgTi8wjsWh7IOPE9ilUwNVm",
			AppAdmin:           false,
			CreatedAt:          yesterday,
			Avatar:             &enum.Avatar1,
			Language:           enum.English,
			AppTheme:           enum.Light,
			LastVisitedGroupID: new(int64(1)),
		},
		{
			Username:           "testuser2",
			Password:           "$2a$12$Z30jTp2WrTWT1jOcnZiXvOcIcqhFNyNnKt7yS7FcUUaIHdgVPy3k2",
			AppAdmin:           true,
			CreatedAt:          yesterday,
			Avatar:             &enum.Avatar2,
			Language:           enum.French,
			AppTheme:           enum.Dark,
			LastVisitedGroupID: new(int64(2)),
		},
		{
			Username:           "testuser3",
			Password:           "$2a$12$flHptXw2TVYQs3b74duKJO.AkxIoaFPctDSp0AtquuTc82xte4wwy",
			AppAdmin:           false,
			CreatedAt:          yesterday,
			Avatar:             &enum.Avatar3,
			Language:           enum.English,
			AppTheme:           enum.System,
			LastVisitedGroupID: new(int64(3)),
		},
		{
			Username:           "testuser4",
			Password:           "$2a$12$8dCvoylHH5QIRHlpurXJ3ORMqeGwRkfP3XzytQUVxuPjoIbzj9PWa",
			AppAdmin:           false,
			CreatedAt:          yesterday,
			Avatar:             nil,
			Language:           enum.English,
			AppTheme:           enum.Light,
			LastVisitedGroupID: nil,
		},
	}
)

func CompareSlicesUsers(actual, expected []model.User) error {
	if len(actual) != (len(expected) + 1) {
		return fmt.Errorf("expected %d users, got %d", len(expected)+1, len(actual))
	}

	actual = utils.SortSliceByFieldName(actual, "ID", false)
	expected = utils.SortSliceByFieldName(expected, "ID", false)

	// Start at 1 to skip the admin user created by the migration script
	for i := 1; i < len(actual); i++ {
		actualUser := actual[i]
		expectedUser := expected[i-1]

		if err := compareUsers(&actualUser, &expectedUser); err != nil {
			return err
		}
	}

	return nil
}

func compareUsers(actual, expected *model.User) error {
	if actual.Username != expected.Username {
		return fmt.Errorf("Username = %s instead of %s", actual.Username, expected.Username)
	}

	if actual.Password != expected.Password {
		return fmt.Errorf("Password = %s instead of %s", actual.Password, expected.Password)
	}

	if actual.AppAdmin != expected.AppAdmin {
		return fmt.Errorf("AppAdmin ='%v'instead of %v", actual.AppAdmin, expected.AppAdmin)
	}

	// allows a small time difference (±1 minute) to account for timing variations
	if !utils.TimesApproximatelyEqual(actual.CreatedAt, expected.CreatedAt, time.Minute) {
		return fmt.Errorf("CreatedAt ='%v'instead of around'%v'(±1min)", actual.CreatedAt, expected.CreatedAt)
	}

	// Compare Avatar pointers
	if (actual.Avatar == nil) != (expected.Avatar == nil) {
		return fmt.Errorf("Avatar nil mismatch: got'%v'instead of %v", actual.Avatar, expected.Avatar)
	} else if actual.Avatar != nil && expected.Avatar != nil && *actual.Avatar != *expected.Avatar {
		return fmt.Errorf("Avatar ='%v'instead of %v", *actual.Avatar, *expected.Avatar)
	}

	if actual.Language != expected.Language {
		return fmt.Errorf("Language = %s instead of %s", actual.Language, expected.Language)
	}

	if actual.AppTheme != expected.AppTheme {
		return fmt.Errorf("AppTheme = %s instead of %s", actual.AppTheme, expected.AppTheme)
	}

	// Compare LastVisitedGroupID pointers
	if (actual.LastVisitedGroupID == nil) != (expected.LastVisitedGroupID == nil) {
		return fmt.Errorf("LastVisitedGroupID nil mismatch: got'%v'instead of %v", actual.LastVisitedGroupID, expected.LastVisitedGroupID)
	} else if actual.LastVisitedGroupID != nil && expected.LastVisitedGroupID != nil && *actual.LastVisitedGroupID != *expected.LastVisitedGroupID {
		return fmt.Errorf("LastVisitedGroupID = %d instead of %d", *actual.LastVisitedGroupID, *expected.LastVisitedGroupID)
	}
	return nil
}

// setupUserTestDB initializes an in-memory SQLite database with test data for testing.
func setupUserTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Apply migrations using the migration package
	migrationsFS := os.DirFS("../../data/migrations")
	err = migration.Migrate(db, migrationsFS)
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	for groupID := 1; groupID <= 3; groupID++ {
		_, err = db.Exec(
			`INSERT INTO groups (id, name, image_url, created_at) VALUES (?, ?, ?, ?);`,
			groupID,
			"group-"+strconv.Itoa(groupID),
			nil,
			yesterday,
		)
		if err != nil {
			t.Fatalf("failed to insert test group '%d': %v", groupID, err)
		}
	}

	// Insert test users
	for i, user := range expectedUsers {
		res, err := db.Exec(
			`INSERT INTO users (username, password, app_admin, created_at, avatar, language, app_theme, last_visited_group_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
			user.Username,
			user.Password,
			user.AppAdmin,
			user.CreatedAt,
			user.Avatar,
			user.Language,
			user.AppTheme,
			user.LastVisitedGroupID,
		)
		if err != nil {
			t.Fatalf("failed to insert test user '%s'", user.Username)
		}

		expectedUsers[i].ID, err = res.LastInsertId()
		if err != nil {
			t.Fatalf("failed to get last insert ID for user '%s'", user.Username)
		}
	}

	return db
}

/*** TEST CONSTRUCTOR ***/

func TestNewUserRepository(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	if repo == nil {
		t.Fatal("expected non-nil UserRepository")
	}

	if repo.db == nil {
		t.Fatal("expected non-nil database connection")
	}
}

/*** READ OPERATIONS TESTS ***/

func TestGetAllUsers(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if err := CompareSlicesUsers(users, expectedUsers); err != nil {
		t.Error("GetAll() returned users with mismatched fields: " + err.Error())
	}
}

func TestGetByUserID(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name     string
		wantErr  error
		wantUser model.User
	}{
		{
			name:     "existing user 1",
			wantErr:  nil,
			wantUser: expectedUsers[0],
		},
		{
			name:     "existing user 2",
			wantErr:  nil,
			wantUser: expectedUsers[1],
		},
		{
			name:    "non-existing user",
			wantErr: customErrors.NewNotFoundError("User", fmt.Sprintf("%d", invalidId), sql.ErrNoRows),
			wantUser: model.User{
				ID: invalidId,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(tt.wantUser.ID)

			if tt.wantErr != nil {
				if !utils.CompareErrors(err, tt.wantErr) {
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
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name     string
		username string
		wantErr  error
		wantUser model.User
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
			wantErr:  customErrors.NewNotFoundError("User", fmt.Sprintf("%s", invalidUsername), sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByUsername(tt.username)

			if tt.wantErr != nil {
				if !utils.CompareErrors(err, tt.wantErr) {
					t.Errorf("GetByUsername() error = '%v' instead of '%v'", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByUsername() unexpected error = %v", err)
			}

			if err := compareUsers(user, &tt.wantUser); err != nil {
				t.Errorf("GetByUsername() returned user does not match expected user: %v", err.Error())
			}
		})
	}
}

/*** CREATE OPERATIONS TESTS ***/

func TestCreateUser(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	avatar1 := enum.Avatar1
	avatar2 := enum.Avatar2

	tests := []struct {
		name    string
		user    *model.User
		wantErr error
	}{
		{
			name: "create new user",
			user: &model.User{
				Username:  "newuser",
				Password:  "newpassword",
				AppAdmin:  false,
				Avatar:    &avatar1,
				Language:  enum.English,
				AppTheme:  enum.Light,
				CreatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "create admin user",
			user: &model.User{
				Username:  "newuser2",
				Password:  "newpassword",
				AppAdmin:  true,
				Avatar:    &avatar2,
				Language:  enum.French,
				AppTheme:  enum.Dark,
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
				if !utils.CompareErrors(err, tt.wantErr) {
					t.Errorf("Create() error = '%v' instead of '%v'", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() unexpected error = %v", err)
			}

			// Verify the user was actually created
			createdUser, err := repo.GetByID(id)
			if err != nil {
				t.Fatalf("failed to fetch created user: %v", err)
			}

			if err := compareUsers(createdUser, tt.user); err != nil {
				t.Errorf("Actual created user does not match expected user: %v", err.Error())
			}
		})
	}
}

/*** UPDATE OPERATIONS TESTS ***/

func TestUpdate(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	avatar1 := enum.Avatar1
	avatar2 := enum.Avatar2

	tests := []struct {
		name    string
		user    *model.User
		wantErr error
	}{
		{
			name: "duplicate username update",
			user: &model.User{
				ID:                 expectedUsers[1].ID,
				Username:           expectedUsers[0].Username,
				Password:           expectedUsers[1].Password + "_updated",
				AppAdmin:           true,
				CreatedAt:          expectedUsers[1].CreatedAt,
				Avatar:             &avatar2,
				Language:           enum.French,
				AppTheme:           enum.Dark,
				LastVisitedGroupID: expectedUsers[1].LastVisitedGroupID,
			},
			wantErr: customErrors.NewConflictError("User", "already exists", sqlite3.ErrConstraintUnique),
		},
		{
			name: "update non-existing user",
			user: &model.User{
				ID:       invalidId,
				Username: "nonexistent",
				Password: "password",
				AppAdmin: false,
				Avatar:   &avatar1,
				Language: enum.English,
				AppTheme: enum.Light,
			},
			wantErr: customErrors.NewNotFoundError("User", strconv.FormatInt(invalidId, 10), nil),
		},
		{
			name: "no field updated",
			user: &model.User{
				ID:                 expectedUsers[0].ID,
				Username:           expectedUsers[0].Username,
				Password:           expectedUsers[0].Password,
				AppAdmin:           expectedUsers[0].AppAdmin,
				CreatedAt:          expectedUsers[0].CreatedAt,
				Avatar:             expectedUsers[0].Avatar,
				Language:           expectedUsers[0].Language,
				AppTheme:           expectedUsers[0].AppTheme,
				LastVisitedGroupID: expectedUsers[0].LastVisitedGroupID,
			},
			wantErr: nil,
		},
		{
			name: "update existing user",
			user: &model.User{
				ID:                 expectedUsers[0].ID,
				Username:           expectedUsers[0].Username + "_updated",
				Password:           expectedUsers[0].Password + "_updated",
				AppAdmin:           true,
				CreatedAt:          expectedUsers[0].CreatedAt,
				Avatar:             &avatar2,
				Language:           enum.French,
				AppTheme:           enum.Dark,
				LastVisitedGroupID: expectedUsers[0].LastVisitedGroupID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.user)

			if tt.wantErr != nil {
				if !utils.CompareErrors(err, tt.wantErr) {
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
	db := setupUserTestDB(t)
	defer db.Close()

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
			wantErr: customErrors.NewNotFoundError("User", strconv.FormatInt(invalidId, 10), nil),
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
				if !utils.CompareErrors(err, tt.wantErr) {
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

func TestDeleteUser(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		id      int64
		wantErr error
	}{
		{
			name:    "delete non-existing user",
			id:      invalidId,
			wantErr: customErrors.NewNotFoundError("User", strconv.FormatInt(invalidId, 10), nil),
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
				if !utils.CompareErrors(err, tt.wantErr) {
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
			} else if appErr, ok := err.(customErrors.AppError); !ok || appErr.HTTPStatus() != 404 {
				t.Errorf("expected NotFound error, but GetByID() error = %v", err)
			}
		})
	}
}

func TestDeleteThenGetAll(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

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
