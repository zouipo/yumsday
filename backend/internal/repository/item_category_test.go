package repository

import (
	"database/sql"
	"fmt"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	invalidICID      = int64(-1)
	invalidICName    = "INVALID CATEGORY NAME"
	invalidICGroupID = int64(-1)
)

var testItemCategories = []model.ItemCategory{
	{ID: 1, Name: "GRAINS AND PASTA", GroupID: 1},
	{ID: 2, Name: "BAKED GOODS", GroupID: 1},
	{ID: 3, Name: "SPICES AND CONDIMENTS", GroupID: 1},
	{ID: 4, Name: "DAIRY", GroupID: 1},
	{ID: 5, Name: "MEAT", GroupID: 2},
	{ID: 6, Name: "VEGETABLES", GroupID: 2},
	{ID: 7, Name: "SNACKS", GroupID: 2},
	{ID: 8, Name: "CANNED GOODS", GroupID: 2},
	{ID: 9, Name: "BEVERAGE", GroupID: 2},
}

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
	db := utils.SetUpTestDB(t)
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
	db := utils.SetUpTestDB(t)
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
			categoryID: testItemCategories[0].ID,
			expected:   &testItemCategories[0],
			expectErr:  nil,
		},
		{
			name:       "Get item category by invalid ID",
			categoryID: invalidICID,
			expected:   nil,
			expectErr:  customErrors.NewNotFoundError("ItemCategory", "item_categories.id", sql.ErrNoRows),
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

func TestGetItemCategoryByNameAndGroupID(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewItemCategoryRepository(db)

	tests := []struct {
		name      string
		icName    string
		groupID   int64
		expected  *model.ItemCategory
		expectErr error
	}{
		{
			name:      "Valid name and group ID",
			icName:    testItemCategories[0].Name,
			groupID:   testItemCategories[0].GroupID,
			expected:  &testItemCategories[0],
			expectErr: nil,
		},
		{
			name:      "Invalid name and valid group ID",
			icName:    invalidICName,
			groupID:   testItemCategories[0].GroupID,
			expected:  nil,
			expectErr: customErrors.NewNotFoundError("ItemCategory", "item_categories.name,item_categories.group_id", nil),
		},
		{
			name:      "Valid name and invalid group ID",
			icName:    testItemCategories[0].Name,
			groupID:   invalidICGroupID,
			expected:  nil,
			expectErr: customErrors.NewNotFoundError("ItemCategory", "item_categories.name,item_categories.group_id", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := repo.GetByNameAndGroupID(tt.icName, tt.groupID)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByNameAndGroupID() unexpected error = %v", err)
			}

			if err := compareItemCategory(category, tt.expected); err != nil {
				t.Errorf("GetByNameAndGroupID() item category does not match expected: %v", err)
			}
		})
	}
}
