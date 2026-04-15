package repository

import (
	"database/sql"
	"reflect"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	invalidItemId     = int64(-1)
	invalidGroupId    = int64(-1)
	invalidCategoryId = int64(-1)
	invalidFieldSort  = "invalidField"
	invalidName       = "invalidName"
	validItemID       = int64(1)

	itemCategory1 = model.ItemCategory{ID: 1, Name: "GRAINS AND PASTA"}
	itemCategory2 = model.ItemCategory{ID: 2, Name: "BAKED GOODS"}
	itemCategory3 = model.ItemCategory{ID: 3, Name: "SPICES AND CONDIMENTS"}
	itemCategory4 = model.ItemCategory{ID: 4, Name: "DAIRY"}
	itemCategory5 = model.ItemCategory{ID: 5, Name: "MEAT"}
	itemCategory6 = model.ItemCategory{ID: 6, Name: "VEGETABLES"}
	itemCategory7 = model.ItemCategory{ID: 7, Name: "SNACKS"}
	itemCategory8 = model.ItemCategory{ID: 8, Name: "CANNED GOODS"}
	itemCategory9 = model.ItemCategory{ID: 9, Name: "BEVERAGE"}

	expectedItems = []model.Item{
		{
			ID:                 1,
			Name:               "Flour",
			Description:        new("All-purpose flour"),
			AverageMarketPrice: new(2.50),
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory:       itemCategory1,
		},
		{
			ID:                 2,
			Name:               "Sugar",
			Description:        new("White granulated sugar"),
			AverageMarketPrice: new(1.80),
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory:       itemCategory2,
		},
		{
			ID:                 3,
			Name:               "Salt",
			Description:        new("Table salt"),
			AverageMarketPrice: new(0.50),
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory:       itemCategory3,
		},
		{
			ID:                 4,
			Name:               "Eggs",
			Description:        new("Large eggs"),
			AverageMarketPrice: new(3.50),
			UnitType:           enum.Piece,
			GroupID:            1,
			ItemCategory:       itemCategory4,
		},
		{
			ID:                 5,
			Name:               "Milk",
			Description:        new("Whole milk"),
			AverageMarketPrice: new(2.20),
			UnitType:           enum.Volume,
			GroupID:            1,
			ItemCategory:       itemCategory4,
		},
		{
			ID:                 6,
			Name:               "Butter",
			Description:        new("Unsalted butter"),
			AverageMarketPrice: new(4.00),
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory:       itemCategory4,
		},
		{
			ID:                 7,
			Name:               "Chicken Breast",
			Description:        new("Boneless skinless chicken breast"),
			AverageMarketPrice: new(8.50),
			UnitType:           enum.Weight,
			GroupID:            2,
			ItemCategory:       itemCategory5,
		},
		{
			ID:                 8,
			Name:               "Tomatoes",
			Description:        new("Fresh tomatoes"),
			AverageMarketPrice: new(3.00),
			UnitType:           enum.Weight,
			GroupID:            2,
			ItemCategory:       itemCategory6,
		},
		{
			ID:                 9,
			Name:               "Onions",
			Description:        new("Yellow onions"),
			AverageMarketPrice: new(1.50),
			UnitType:           enum.Weight,
			GroupID:            2,
			ItemCategory:       itemCategory6,
		},
		{
			ID:                 10,
			Name:               "Garlic",
			Description:        new("Fresh garlic"),
			AverageMarketPrice: new(2.00),
			UnitType:           enum.Weight,
			GroupID:            2,
			ItemCategory:       itemCategory6,
		},
		{
			ID:                 11,
			Name:               "Water",
			Description:        nil,
			AverageMarketPrice: nil,
			UnitType:           enum.Volume,
			GroupID:            2,
			ItemCategory:       itemCategory9,
		},
		{
			ID:                 12,
			Name:               "Pepper",
			Description:        nil,
			AverageMarketPrice: new(1.20),
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory:       itemCategory3,
		},
		{
			ID:                 13,
			Name:               "Olive Oil",
			Description:        new("Extra virgin olive oil"),
			AverageMarketPrice: nil,
			UnitType:           enum.Volume,
			GroupID:            1,
			ItemCategory:       itemCategory3,
		},
		{
			ID:                 14,
			Name:               "Potato Chips",
			Description:        new("Salted potato chips"),
			AverageMarketPrice: new(2.99),
			UnitType:           enum.Bag,
			GroupID:            2,
			ItemCategory:       itemCategory7,
		},
		{
			ID:                 15,
			Name:               "Canned Beans",
			Description:        new("Black beans"),
			AverageMarketPrice: new(1.50),
			UnitType:           enum.Numeric,
			GroupID:            1,
			ItemCategory:       itemCategory8,
		},
	}
)

func itemsByGroupID(items []model.Item, groupID int64) []model.Item {
	filtered := make([]model.Item, 0)
	for _, item := range items {
		if item.GroupID == groupID {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func CompareSlicesItems(s1, s2 []model.Item) bool {
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

/*** TEST CONSTRUCTOR ***/
func TestNewItemRepository(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	if repo == nil {
		t.Fatal("expected non-nil repository, got nil")
	}

	if repo.db == nil {
		t.Fatal("expected non-nil database connection, got nil")
	}
}

/*** READ OPERATIONS TESTS ***/
func TestGetItemsByGroupID(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	result, err := db.Exec("INSERT INTO groups (name, image_url, created_at) VALUES (?, ?, ?)", "EmptyGroupForItems", nil, 0)
	if err != nil {
		t.Fatalf("failed to create empty group for test: %v", err)
	}

	emptyGroupID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get empty group ID: %v", err)
	}

	repo := NewItemRepository(db)

	tests := []struct {
		name       string
		groupID    int64
		sortBy     string
		descending bool
		expected   []model.Item
		expectErr  error
	}{
		{
			name:       "Valid group ID with sorting by name",
			groupID:    1,
			sortBy:     "items.name",
			descending: false,
			expected:   utils.SortSliceByFieldName(itemsByGroupID(expectedItems, 1), "Name", false),
			expectErr:  nil,
		},
		{
			name:       "Valid group ID with sorting by average market price",
			groupID:    1,
			sortBy:     "items.average_market_price",
			descending: false,
			expected:   utils.SortSliceByFieldName(itemsByGroupID(expectedItems, 1), "AverageMarketPrice", false),
			expectErr:  nil,
		},
		{
			name:       "Valid group ID with sorting by unit type",
			groupID:    1,
			sortBy:     "items.unit_type",
			descending: false,
			expected:   utils.SortSliceByFieldName(itemsByGroupID(expectedItems, 1), "UnitType", false),
			expectErr:  nil,
		},
		{
			name:       "Valid group ID with sorting by item category name",
			groupID:    1,
			sortBy:     "item_categories.name",
			descending: false,
			expected:   utils.SortSliceByFieldName(itemsByGroupID(expectedItems, 1), "ItemCategory.Name", false),
			expectErr:  nil,
		},
		{
			name:       "Valid group ID with sorting by invalid field",
			groupID:    1,
			sortBy:     invalidFieldSort,
			descending: false,
			expected:   nil,
			expectErr:  customErrors.NewInternalError("failed to fetch items", nil),
		},
		{
			name:       "Valid group ID with no items",
			groupID:    invalidGroupId,
			sortBy:     "items.name",
			descending: false,
			expected:   []model.Item{},
			expectErr:  nil,
		},
		{
			name:       "Valid group ID 2 with sorting by name",
			groupID:    2,
			sortBy:     "items.name",
			descending: false,
			expected:   utils.SortSliceByFieldName(itemsByGroupID(expectedItems, 2), "Name", false),
			expectErr:  nil,
		},
		{
			name:       "Valid group ID 3 with no items",
			groupID:    3,
			sortBy:     "items.name",
			descending: false,
			expected:   []model.Item{},
			expectErr:  nil,
		},
		{
			name:       "Existing group with no items returns empty list",
			groupID:    emptyGroupID,
			sortBy:     "items.name",
			descending: false,
			expected:   []model.Item{},
			expectErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := repo.GetByGroupID(tt.groupID, tt.sortBy, tt.descending)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByGroupID() unexpected error = %v", err)
			}

			if !CompareSlicesItems(items, tt.expected) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.expected, items)
			}
		})
	}
}

func TestGetItemById(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	tests := []struct {
		name      string
		id        int64
		expected  model.Item
		expectErr error
	}{
		{
			name:      "Get item by valid ID 1",
			id:        validItemID,
			expected:  utils.SortSliceByFieldName(expectedItems, "ID", false)[validItemID-1],
			expectErr: nil,
		},
		{
			name:      "Get item by valid ID 2",
			id:        validItemID + 1,
			expected:  utils.SortSliceByFieldName(expectedItems, "ID", false)[validItemID],
			expectErr: nil,
		},
		{
			name:      "Get item by invalid ID",
			id:        invalidItemId,
			expected:  model.Item{},
			expectErr: customErrors.NewNotFoundError("Item", "items.id", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := repo.GetByID(tt.id)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(item, &tt.expected) {
				t.Errorf("Items should be equal: %v, %v", item, tt.expected)
			}
		})
	}
}

func TestGetItemByName(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	tests := []struct {
		name      string
		itemName  string
		expected  []model.Item
		expectErr error
	}{
		{
			name:      "Get item by valid name",
			itemName:  utils.SortSliceByFieldName(expectedItems, "Name", false)[0].Name,
			expected:  []model.Item{utils.SortSliceByFieldName(expectedItems, "Name", false)[0]},
			expectErr: nil,
		},
		{
			name:      "Get item by invalid name returns empty slice",
			itemName:  invalidName,
			expected:  []model.Item{},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := repo.GetByName(tt.itemName)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByName() unexpected error = %v", err)
			}

			if !CompareSlicesItems(items, tt.expected) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.expected, items)
			}
		})
	}
}

/*** CREATE OPERATIONS TESTS ***/
func TestCreateItem(t *testing.T) {
	tests := []struct {
		name      string
		item      model.Item
		expectErr error
	}{
		{
			name:      "Create item with nil values",
			item:      expectedItems[0],
			expectErr: nil,
		},
		{
			name:      "Create item with no nil values",
			item:      expectedItems[1],
			expectErr: nil,
		},
		{
			name: "Create item with invalid group ID",
			item: model.Item{
				Name:    "Invalid Group Item",
				GroupID: invalidGroupId,
			},
			expectErr: customErrors.NewInternalError("failed to create item", nil),
		},
		{
			name: "Create item with invalid category ID",
			item: model.Item{
				Name:    "Invalid Category Item",
				GroupID: invalidGroupId,
				ItemCategory: model.ItemCategory{
					ID:   invalidCategoryId,
					Name: "Invalid Category",
				},
			},
			expectErr: customErrors.NewInternalError("failed to create item", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := utils.SetUpTestDB(t)
			defer db.Close()

			repo := NewItemRepository(db)

			id, err := repo.Create(&tt.item)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() unexpected error = %v", err)
			}

			if id == 0 {
				t.Errorf("expected non-zero item ID after creation, got %d", id)
			}

			lastId := utils.SortSliceByFieldName(expectedItems, "ID", false)[len(expectedItems)-1].ID
			if id <= lastId {
				t.Errorf("expected item ID to be sequential: expected > %d, got %d", lastId, id)
			}

			// Verify the item was created correctly
			createdItem, err := repo.GetByID(id)
			if err != nil {
				t.Fatalf("failed to retrieve created item: %v", err)
			}

			tt.item.ID = lastId + 1
			if !reflect.DeepEqual(createdItem, &tt.item) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.item, createdItem)
			}
		})
	}
}

/*** UPDATE OPERATIONS TESTS ***/
func TestUpdateItem(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	tests := []struct {
		name         string
		item         model.Item
		expectedItem model.Item
		expectErr    error
	}{
		{
			name: "Update existing item",
			item: model.Item{
				ID:                 4,
				Name:               "Free-range Eggs",
				Description:        new("Large free-range eggs"),
				AverageMarketPrice: new(4.0),
				UnitType:           enum.Piece,
				GroupID:            1,
				ItemCategory: model.ItemCategory{
					ID:   4,
					Name: "DAIRY",
				},
			},
			expectedItem: model.Item{
				ID:                 4,
				Name:               "Free-range Eggs",
				Description:        new("Large free-range eggs"),
				AverageMarketPrice: new(4.0),
				UnitType:           enum.Piece,
				GroupID:            1,
				ItemCategory: model.ItemCategory{
					ID:   4,
					Name: "DAIRY",
				},
			},
			expectErr: nil,
		},
		{
			name:         "No field updated",
			item:         expectedItems[0],
			expectedItem: expectedItems[0],
			expectErr:    nil,
		},
		{
			name: "Update non-existing item",
			item: model.Item{
				ID:      invalidItemId,
				Name:    "Non-existing Item",
				GroupID: 1,
				ItemCategory: model.ItemCategory{
					ID:   1,
					Name: "Baking",
				},
			},
			expectedItem: model.Item{},
			expectErr:    customErrors.NewNotFoundError("Item", "items.id", sql.ErrNoRows),
		},
		{
			name: "Update item without taking group ID into account",
			item: model.Item{
				ID:                 expectedItems[0].ID,
				Name:               "Super Flour",
				Description:        expectedItems[0].Description,
				AverageMarketPrice: expectedItems[0].AverageMarketPrice,
				UnitType:           expectedItems[0].UnitType,
				GroupID:            2,
				ItemCategory: model.ItemCategory{
					ID:   1,
					Name: "GRAINS AND PASTA",
				},
			},
			expectedItem: model.Item{
				ID:                 expectedItems[0].ID,
				Name:               "Super Flour",
				Description:        expectedItems[0].Description,
				AverageMarketPrice: expectedItems[0].AverageMarketPrice,
				UnitType:           expectedItems[0].UnitType,
				GroupID:            expectedItems[0].GroupID,
				ItemCategory: model.ItemCategory{
					ID:   1,
					Name: "GRAINS AND PASTA",
				},
			},
			expectErr: nil,
		},
		{
			name: "Update item with invalid category ID",
			item: model.Item{
				ID:      expectedItems[0].ID,
				Name:    "Valid name",
				GroupID: 1,
				ItemCategory: model.ItemCategory{
					ID:   invalidCategoryId,
					Name: "Invalid Category",
				},
			},
			expectedItem: model.Item{},
			expectErr:    customErrors.NewInternalError("failed to update item", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(&tt.item)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Update() unexpected error = %v", err)
			}

			// Verify the item was updated correctly
			updatedItem, err := repo.GetByID(tt.item.ID)
			if err != nil {
				t.Fatalf("failed to retrieve updated item: %v", err)
			}

			if !reflect.DeepEqual(updatedItem, &tt.expectedItem) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.expectedItem, updatedItem)
			}
		})
	}
}

/*** DELETE OPERATIONS TESTS ***/
func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		expectErr error
	}{
		{
			name: "Delete existing item",
			// Can be deleted because it is not referenced in any other tables (typically ingredients, groceries)
			id:        14,
			expectErr: nil,
		},
		{
			name:      "Delete non-existing item",
			id:        invalidItemId,
			expectErr: customErrors.NewNotFoundError("Item", "items.id", sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := utils.SetUpTestDB(t)
			defer db.Close()

			repo := NewItemRepository(db)

			err := repo.Delete(tt.id)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Delete() unexpected error = %v", err)
			}

			// Verify the item was deleted
			_, err = repo.GetByID(tt.id)
			if !utils.CompareErrors(err, customErrors.NewNotFoundError("Item", "items.id", sql.ErrNoRows)) {
				t.Errorf("expected item to be deleted, but it still exists")
			}
		})
	}
}
