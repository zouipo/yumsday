package repository

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var invalidGroupRepositoryID = int64(-1)

func compareGroupMembers(actual, expected []model.GroupMember) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d group members, got %d", len(expected), len(actual))
	}

	sort.Slice(actual, func(i, j int) bool {
		return actual[i].UserID < actual[j].UserID
	})
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].UserID < expected[j].UserID
	})

	for i := range actual {
		if actual[i].UserID != expected[i].UserID {
			return fmt.Errorf("expected UserID %d, got %d", expected[i].UserID, actual[i].UserID)
		}
		if actual[i].Admin != expected[i].Admin {
			return fmt.Errorf("expected Admin %v, got %v", expected[i].Admin, actual[i].Admin)
		}
		if !utils.TimesApproximatelyEqual(actual[i].JoinedAt, expected[i].JoinedAt, time.Minute) {
			return fmt.Errorf("expected JoinedAt around %v, got %v", expected[i].JoinedAt, actual[i].JoinedAt)
		}
	}

	return nil
}

func compareGroup(actual, expected *model.Group) error {
	if actual.ID != expected.ID {
		return fmt.Errorf("expected ID %d, got %d", expected.ID, actual.ID)
	}
	if actual.Name != expected.Name {
		return fmt.Errorf("expected Name %s, got %s", expected.Name, actual.Name)
	}
	if (actual.ImageURL == nil) != (expected.ImageURL == nil) ||
		(actual.ImageURL != nil && *actual.ImageURL != *expected.ImageURL) {
		return fmt.Errorf("expected ImageURL %v, got %v", expected.ImageURL, actual.ImageURL)
	}
	if !utils.TimesApproximatelyEqual(actual.CreatedAt, expected.CreatedAt, time.Minute) {
		return fmt.Errorf("expected CreatedAt around %v, got %v", expected.CreatedAt, actual.CreatedAt)
	}
	if err := compareGroupMembers(actual.Members, expected.Members); err != nil {
		return err
	}

	return nil
}

func setupGroupTestDB(t *testing.T) (*sql.DB, *model.Group) {
	db := setUpTestDB(t)

	createdAt := time.Now().UTC()
	joinedAtOne := createdAt.Add(-2 * time.Hour)
	joinedAtTwo := createdAt.Add(-time.Hour)

	_, err := db.Exec(
		`INSERT INTO users (username, password, app_admin, created_at, language, app_theme)
		VALUES (?, ?, ?, ?, ?, ?);`,
		"group-member-1",
		"password-1",
		false,
		createdAt,
		"EN",
		"LIGHT",
	)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	_, err = db.Exec(`UPDATE groups SET image_url = ?, created_at = ? WHERE id = ?;`, nil, createdAt, 1)
	if err != nil {
		t.Fatalf("failed to update test group: %v", err)
	}

	_, err = db.Exec(
		`INSERT INTO group_members (user_id, group_id, admin, joined_at)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?);`,
		1, 1, true, joinedAtOne,
		2, 1, false, joinedAtTwo,
	)
	if err != nil {
		t.Fatalf("failed to insert test group members: %v", err)
	}

	return db, &model.Group{
		ID:        1,
		Name:      "Test Group",
		ImageURL:  nil,
		CreatedAt: createdAt,
		Members: []model.GroupMember{
			{UserID: 1, GroupId: 1, Admin: true, JoinedAt: joinedAtOne},
			{UserID: 2, GroupId: 1, Admin: false, JoinedAt: joinedAtTwo},
		},
	}
}

func TestNewGroupRepository(t *testing.T) {
	db, _ := setupGroupTestDB(t)
	defer db.Close()

	repo := NewGroupRepository(db)

	if repo == nil {
		t.Fatal("expected non-nil repository, got nil")
	}

	if repo.db == nil {
		t.Fatal("expected non-nil database connection, got nil")
	}
}

func TestGetGroupByID(t *testing.T) {
	db, expected := setupGroupTestDB(t)
	defer db.Close()

	repo := NewGroupRepository(db)

	tests := []struct {
		name      string
		groupID   int64
		expected  *model.Group
		expectErr error
	}{
		{
			name:      "Get group by valid ID",
			groupID:   1,
			expected:  expected,
			expectErr: nil,
		},
		{
			name:      "Get group by invalid ID",
			groupID:   invalidGroupRepositoryID,
			expected:  nil,
			expectErr: customErrors.NewNotFoundError("Group", strconv.FormatInt(invalidGroupRepositoryID, 10), sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group, err := repo.GetByID(tt.groupID)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if err := compareGroup(group, tt.expected); err != nil {
				t.Errorf("GetByID() group does not match expected: %v", err)
			}
		})
	}
}
