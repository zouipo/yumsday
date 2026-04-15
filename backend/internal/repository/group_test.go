package repository

import (
	"database/sql"
	"fmt"
	"sort"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	invalidGroupRepositoryID = int64(-1)

	testGroups = newTestGroups()
)

func newTestGroups() []model.Group {
	epoch := time.Unix(0, 0)

	return []model.Group{
		{
			ID:        1,
			Name:      "Family",
			ImageURL:  new("/static/images/family.jpg"),
			CreatedAt: epoch,
			Members: []model.GroupMember{
				{UserID: 2, GroupId: 1, Admin: true, JoinedAt: epoch},
				{UserID: 3, GroupId: 1, Admin: true, JoinedAt: epoch},
				{UserID: 4, GroupId: 1, Admin: false, JoinedAt: epoch},
			},
		},
		{
			ID:        2,
			Name:      "Friends",
			ImageURL:  new("/static/images/friends.jpg"),
			CreatedAt: epoch,
			Members: []model.GroupMember{
				{UserID: 2, GroupId: 2, Admin: false, JoinedAt: epoch},
				{UserID: 4, GroupId: 2, Admin: true, JoinedAt: epoch},
			},
		},
		{
			ID:        3,
			Name:      "Work",
			ImageURL:  nil,
			CreatedAt: epoch,
			Members:   []model.GroupMember{},
		},
	}
}

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

func TestNewGroupRepository(t *testing.T) {
	db := utils.SetUpTestDB(t)
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
	db := utils.SetUpTestDB(t)
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
			expected:  &testGroups[0],
			expectErr: nil,
		},
		{
			name:      "Get group by invalid ID",
			groupID:   invalidGroupRepositoryID,
			expected:  nil,
			expectErr: customErrors.NewNotFoundError("groups", "id", sql.ErrNoRows),
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
