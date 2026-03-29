package repository

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

var (
	invalidItemId   = int64(-1)
	invalidNameId   = "invalidName"
	ErrItemNotFound = "Item not found"

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
			Name:               "Milk",
			Description:        func() *string { v := "Whole milk, 1 liter"; return &v }(),
			AverageMarketPrice: nil,
			UnitType:           enum.Volume,
			GroupID:            1,
			ItemCategory: model.ItemCategory{
				ID:   1,
				Name: "Baking",
			},
		},
		{
			ID:                 3,
			Name:               "Potatoes",
			Description:        func() *string { v := "Fresh potatoes, 1 kg"; return &v }(),
			AverageMarketPrice: func() *float64 { v := 3.5; return &v }(),
			UnitType:           enum.Weight,
			GroupID:            1,
			ItemCategory: model.ItemCategory{
				ID:   1,
				Name: "Vegetables",
			},
		},
	}
)

func compareListItems(actual, expected []model.Item) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d items, got %d", len(expected), len(actual))
	}

	sortItemsByID(actual)
	sortItemsByID(expected)

	for i := range actual {
		if err := compareItems(&actual[i], &expected[i]); err != nil {
			return fmt.Errorf("item at index %d does not match: %v", i, err)
		}
	}

	return nil
}

func compareItems(actual, expected *model.Item) error {
	if actual.ID != expected.ID {
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

func sortItemsByID(items []model.Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
}

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

	// Insert test items
	for i, item := range expectedItems {
		res, err := db.Exec(
			`INSERT INTO items (name, description, average_market_price, unit_type, item_category_id, group_id)
			VALUES (?, ?, ?, ?, ?, ?);`,
			item.Name,
			item.Description,
			item.AverageMarketPrice,
			item.UnitType,
			item.ItemCategory.ID,
			item.GroupID,
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
