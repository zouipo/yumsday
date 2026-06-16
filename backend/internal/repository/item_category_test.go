package repository

import (
	"database/sql"
	"reflect"
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
	{ID: 8, Name: "CANNED GOODS", GroupID: 1},
	{ID: 9, Name: "BEVERAGE", GroupID: 2},
}

func compareSlicesItemCategories(s1, s2 []model.ItemCategory) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if !reflect.DeepEqual(s1[i], s2[i]) {
			return false
		}
	}

	return true
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
			expectErr:  customErrors.NewNotFoundError("item_categories", "id", sql.ErrNoRows),
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

			if !reflect.DeepEqual(category, tt.expected) {
				t.Errorf("item categories should be equal: %v", err)
			}
		})
	}
}

func TestGetItemCategoryByNameAndGroupID(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewItemCategoryRepository(db)

	tests := []struct {
		name       string
		icName     string
		groupID    int64
		descending bool
		expected   []model.ItemCategory
	}{
		{
			name:       "Valid name and group ID",
			icName:     testItemCategories[0].Name,
			groupID:    testItemCategories[0].GroupID,
			descending: false,
			expected: []model.ItemCategory{
				testItemCategories[0],
			},
		},
		{
			name:       "Valid partial name upper-case and group ID, ascending",
			icName:     "GOOD",
			groupID:    testItemCategories[0].GroupID,
			descending: false,
			expected: []model.ItemCategory{
				testItemCategories[1],
				testItemCategories[7],
			},
		},
		{
			name:       "Valid partial name lower-case and group ID, descending",
			icName:     "good",
			groupID:    testItemCategories[0].GroupID,
			descending: true,
			expected: []model.ItemCategory{
				testItemCategories[7],
				testItemCategories[1],
			},
		},
		{
			name:     "Invalid name and valid group ID",
			icName:   invalidICName,
			groupID:  testItemCategories[0].GroupID,
			expected: []model.ItemCategory{},
		},
		{
			name:     "Valid name and invalid group ID",
			icName:   testItemCategories[0].Name,
			groupID:  invalidICGroupID,
			expected: []model.ItemCategory{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.GetByNameAndGroupID(tt.icName, tt.groupID, tt.descending)

			if err != nil {
				t.Fatalf("GetByNameAndGroupID() unexpected error = %v", err)
			}

			if !compareSlicesItemCategories(actual, tt.expected) {
				t.Errorf("item categories should be equal: expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
