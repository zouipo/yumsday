package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var invalidItemCategoryRepositoryID = int64(-1)

func compareItemCategory(actual, expected *model.ItemCategory) error {
	if actual.ID != expected.ID {
		return fmt.Errorf("expected ID %d, got %d", expected.ID, actual.ID)
	}
	if actual.Name != expected.Name {
		return fmt.Errorf("expected Name %s, got %s", expected.Name, actual.Name)
	}
	if actual.GroupID != expected.GroupID {
		return fmt.Errorf("expected GroupID %d, got %d", expected.GroupID, actual.GroupID)
	}

	return nil
}

func TestNewItemCategoryRepository(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	repo := NewItemCategoryRepository(db)

	if repo == nil {
		t.Fatal("expected non-nil repository, got nil")
	}

	if repo.db == nil {
		t.Fatal("expected non-nil database connection, got nil")
	}
}

func TestGetItemCategoryByID(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	repo := NewItemCategoryRepository(db)

	tests := []struct {
		name       string
		categoryID int64
		expected   *model.ItemCategory
		expectErr  error
	}{
		{
			name:       "Get item category by valid ID",
			categoryID: 1,
			expected:   &model.ItemCategory{ID: 1, Name: "Baking"},
			expectErr:  nil,
		},
		{
			name:       "Get item category by second valid ID",
			categoryID: 2,
			expected:   &model.ItemCategory{ID: 2, Name: "Vegetables"},
			expectErr:  nil,
		},
		{
			name:       "Get item category by invalid ID",
			categoryID: invalidItemCategoryRepositoryID,
			expected:   nil,
			expectErr:  customErrors.NewNotFoundError("Item category", strconv.FormatInt(invalidItemCategoryRepositoryID, 10), sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := repo.GetByID(tt.categoryID)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if err := compareItemCategory(category, tt.expected); err != nil {
				t.Errorf("GetByID() item category does not match expected: %v", err)
			}
		})
	}
}
