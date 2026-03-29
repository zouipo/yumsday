package repository

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/migration"
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

	expectedItems = []model.Item{
		{
			ID:                 1,
			Name:               "Flour",
			Description:        nil,
			AverageMarketPrice: nil,
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory: model.ItemCategory{
				ID:   1,
				Name: "Baking",
			},
		},
		{
			ID:                 2,
			Name:               "Potatoes",
			Description:        utils.Ptr("Fresh potatoes, 1 kg"),
			AverageMarketPrice: utils.Ptr(3.5),
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory: model.ItemCategory{
				ID:   2,
				Name: "Vegetables",
			},
		},
		{
			ID:                 3,
			Name:               "Milk",
			Description:        utils.Ptr("Whole milk, 1 liter"),
			AverageMarketPrice: utils.Ptr(1.5),
			UnitType:           enum.Volume,
			GroupID:            1,
			ItemCategory: model.ItemCategory{
				ID:   1,
				Name: "Baking",
			},
		},
		{
			ID:                 4,
			Name:               "Apple",
			Description:        utils.Ptr("Good for making pies, compotes, juice, etc."),
			AverageMarketPrice: utils.Ptr(4.2),
			UnitType:           enum.Numeric,
			GroupID:            1,
			ItemCategory: model.ItemCategory{
				ID:   3,
				Name: "Fruits",
			},
		},
	}
)

// compareListItems compares two slices of Item objects and returns an error if they do not match.
func compareListItems(actual, expected []model.Item) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d items, got %d", len(expected), len(actual))
	}

	for i := range actual {
		if err := compareItems(&actual[i], &expected[i], true); err != nil {
			return fmt.Errorf("item at index %d does not match: %v", i, err)
		}
	}

	return nil
}

// compareItems compares two Item objects and returns an error if they do not match.
func compareItems(actual, expected *model.Item, compareId bool) error {
	if compareId && (actual.ID != expected.ID) {
		return fmt.Errorf("expected ID %d, got %d", expected.ID, actual.ID)
	}
	if actual.Name != expected.Name {
		return fmt.Errorf("expected Name %s, got %s", expected.Name, actual.Name)
	}
	if (actual.Description == nil) != (expected.Description == nil) ||
		(actual.Description != nil && *actual.Description != *expected.Description) {
		return fmt.Errorf("expected Description %v, got %v", expected.Description, actual.Description)
	}
	if (actual.AverageMarketPrice == nil) != (expected.AverageMarketPrice == nil) ||
		(actual.AverageMarketPrice != nil && *actual.AverageMarketPrice != *expected.AverageMarketPrice) {
		return fmt.Errorf("expected AverageMarketPrice %v, got %v", expected.AverageMarketPrice, actual.AverageMarketPrice)
	}
	if actual.UnitType != expected.UnitType {
		return fmt.Errorf("expected UnitType %s, got %s", expected.UnitType, actual.UnitType)
	}
	if actual.GroupID != expected.GroupID {
		return fmt.Errorf("expected GroupID %d, got %d", expected.GroupID, actual.GroupID)
	}
	if actual.ItemCategory.ID != expected.ItemCategory.ID {
		return fmt.Errorf("expected ItemCategory ID %d, got %d", expected.ItemCategory.ID, actual.ItemCategory.ID)
	}
	if actual.ItemCategory.Name != expected.ItemCategory.Name {
		return fmt.Errorf("expected ItemCategory Name %s, got %s", expected.ItemCategory.Name, actual.ItemCategory.Name)
	}
	return nil
}

// sortItemsByField sorts a slice of Item objects by the specified field and returns the sorted slice.
func sortItemsByField(items []model.Item, sortBy string) []model.Item {
	sorted := append([]model.Item(nil), items...)

	sort.Slice(sorted, func(i, j int) bool {
		switch sortBy {
		case "name":
			return sorted[i].Name < sorted[j].Name
		case "average_market_price":
			if sorted[i].AverageMarketPrice == nil {
				return true
			}
			if sorted[j].AverageMarketPrice == nil {
				return false
			}
			return *sorted[i].AverageMarketPrice < *sorted[j].AverageMarketPrice
		case "unit_type":
			return sorted[i].UnitType.String() < sorted[j].UnitType.String()
		case "item_categories.name":
			return sorted[i].ItemCategory.Name < sorted[j].ItemCategory.Name
		default:
			return sorted[i].ID < sorted[j].ID
		}
	})

	return sorted
}

// setUpTestDB initializes an in-memory SQLite database, applies migrations, and inserts test data for items.
func setUpTestDB(t *testing.T) *sql.DB {
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

	// Insert group and item categories for FK constraints on items
	_, err = db.Exec(`
		INSERT INTO groups (id, name, created_at)
		VALUES (1, 'Test Group', unixepoch());
	`)
	if err != nil {
		t.Fatalf("failed to insert test group: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO item_categories (id, name, group_id)
		VALUES
			(1, 'Baking', 1),
			(2, 'Vegetables', 1),
			(3, 'Fruits', 1);
	`)
	if err != nil {
		t.Fatalf("failed to insert test item categories: %v", err)
	}

	// Insert test items
	for i, item := range expectedItems {
		res, err := db.Exec(
			`INSERT INTO items (name, description, average_market_price, unit_type, group_id, item_category_id)
			VALUES (?, ?, ?, ?, ?, ?);`,
			item.Name,
			item.Description,
			item.AverageMarketPrice,
			item.UnitType,
			item.GroupID,
			item.ItemCategory.ID,
		)
		if err != nil {
			t.Fatalf("failed to insert test item '%s'", item.Name)
		}

		expectedItems[i].ID, err = res.LastInsertId()
		if err != nil {
			t.Fatalf("failed to get last insert ID for item '%s'", item.Name)
		}
	}

	return db
}

/*** TEST CONSTRUCTOR ***/
func TestNewItemRepository(t *testing.T) {
	db := setUpTestDB(t)
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
func TestGetAllItemsByGroupID(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	tests := []struct {
		name      string
		groupID   int64
		sortBy    string
		expected  []model.Item
		expectErr error
	}{
		{
			name:      "Valid group ID with sorting by name",
			groupID:   1,
			sortBy:    "i.name",
			expected:  sortItemsByField(expectedItems, "name"),
			expectErr: nil,
		},
		{
			name:      "Valid group ID with sorting by average market price",
			groupID:   1,
			sortBy:    "i.average_market_price",
			expected:  sortItemsByField(expectedItems, "average_market_price"),
			expectErr: nil,
		},
		{
			name:      "Valid group ID with sorting by unit type",
			groupID:   1,
			sortBy:    "i.unit_type",
			expected:  sortItemsByField(expectedItems, "unit_type"),
			expectErr: nil,
		},
		{
			name:      "Valid group ID with sorting by item category name",
			groupID:   1,
			sortBy:    "ic.name",
			expected:  sortItemsByField(expectedItems, "item_categories.name"),
			expectErr: nil,
		},
		{
			name:      "Valid group ID with sorting by invalid field",
			groupID:   1,
			sortBy:    invalidFieldSort,
			expected:  nil,
			expectErr: customErrors.NewInternalError("Failed to fetch items", nil),
		},
		{
			name:      "Invalid group ID",
			groupID:   invalidGroupId,
			sortBy:    "i.name",
			expected:  []model.Item{},
			expectErr: nil,
		},
		{
			name:      "Invalid group ID",
			groupID:   2,
			sortBy:    "i.name",
			expected:  []model.Item{},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := repo.GetAllByGroupID(tt.groupID, tt.sortBy)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetAllByGroupID() unexpected error = %v", err)
			}

			if err := compareListItems(items, tt.expected); err != nil {
				t.Errorf("GetAllByGroupID() items do not match expected: %v", err.Error())
			}
		})
	}
}

func TestGetItemById(t *testing.T) {
	db := setUpTestDB(t)
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
			expected:  sortItemsByField(expectedItems, "id")[validItemID-1],
			expectErr: nil,
		},
		{
			name:      "Get item by valid ID 2",
			id:        validItemID + 1,
			expected:  sortItemsByField(expectedItems, "id")[validItemID],
			expectErr: nil,
		},
		{
			name:      "Get item by invalid ID",
			id:        invalidItemId,
			expected:  model.Item{},
			expectErr: customErrors.NewNotFoundError("Item", strconv.FormatInt(invalidItemId, 10), nil),
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

			if err := compareItems(item, &tt.expected, true); err != nil {
				t.Errorf("GetByID() items do not match expected: %v", err.Error())
			}
		})
	}
}

func TestGetItemByName(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	tests := []struct {
		name      string
		itemName  string
		expected  model.Item
		expectErr error
	}{
		{
			name:      "Get item by valid name",
			itemName:  sortItemsByField(expectedItems, "name")[0].Name,
			expected:  sortItemsByField(expectedItems, "name")[0],
			expectErr: nil,
		},
		{
			name:      "Get item by invalid name",
			itemName:  invalidName,
			expected:  model.Item{},
			expectErr: customErrors.NewNotFoundError("Item", invalidName, nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := repo.GetByName(tt.itemName)

			if tt.expectErr != nil {
				if !utils.CompareErrors(err, tt.expectErr) {
					t.Errorf("expected error '%v', got '%v'", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByName() unexpected error = %v", err)
			}

			if err := compareItems(item, &tt.expected, true); err != nil {
				t.Errorf("GetByName() items do not match expected: %v", err.Error())
			}
		})
	}
}

/*** CREATE OPERATIONS TESTS ***/
func TestCreateItem(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

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
			expectErr: customErrors.NewInternalError("Failed to create item", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			lastId := sortItemsByField(expectedItems, "id")[len(expectedItems)-1].ID
			if id > lastId {
				t.Errorf("expected item ID to be sequential: expected <= %d, got %d", lastId, id)
			}

			// Verify the item was created correctly
			createdItem, err := repo.GetByID(id)
			if err != nil {
				t.Fatalf("failed to retrieve created item: %v", err)
			}

			if err := compareItems(createdItem, &tt.item, false); err != nil {
				t.Errorf("created item does not match expected: %v", err.Error())
			}
		})
	}
}

/*** UPDATE OPERATIONS TESTS ***/
func TestUpdateItem(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	tests := []struct {
		name      string
		item      model.Item
		expectErr error
	}{
		{
			name: "Update existing item",
			item: model.Item{
				ID:                 4,
				Name:               "Carot",
				Description:        utils.Ptr("Good for making pies, soups, juice, etc."),
				AverageMarketPrice: utils.Ptr(2.5),
				UnitType:           enum.Numeric,
				GroupID:            1,
				ItemCategory: model.ItemCategory{
					ID:   2,
					Name: "Vegetables",
				},
			},
			expectErr: nil,
		},
		{
			name:      "No field updated",
			item:      expectedItems[0],
			expectErr: nil,
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
			expectErr: customErrors.NewNotFoundError("Item", strconv.FormatInt(invalidItemId, 10), sql.ErrNoRows),
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

			if err := compareItems(updatedItem, &tt.item, true); err != nil {
				t.Errorf("updated item does not match expected: %v", err.Error())
			}
		})
	}
}

/*** DELETE OPERATIONS TESTS ***/
func TestDeleteItem(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	repo := NewItemRepository(db)

	tests := []struct {
		name      string
		id        int64
		expectErr error
	}{
		{
			name:      "Delete existing item",
			id:        validItemID,
			expectErr: nil,
		},
		{
			name:      "Delete non-existing item",
			id:        invalidItemId,
			expectErr: customErrors.NewNotFoundError("Item", strconv.FormatInt(invalidItemId, 10), sql.ErrNoRows),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if !utils.CompareErrors(err, customErrors.NewNotFoundError("Item", strconv.FormatInt(tt.id, 10), sql.ErrNoRows)) {
				t.Errorf("expected item to be deleted, but it still exists")
			}
		})
	}
}
