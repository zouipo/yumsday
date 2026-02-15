package enums

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestItemCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		category ItemCategory
		expected string
	}{
		{"Fruits", Fruits, "FRUITS"},
		{"Vegetables", Vegetables, "VEGETABLES"},
		{"Meat", Meat, "MEAT"},
		{"Seafood", Seafood, "SEAFOOD"},
		{"Dairy", Dairy, "DAIRY"},
		{"Starch", Starch, "STARCH"},
		{"Beverages", Beverages, "BEVERAGES"},
		{"Snacks", Snacks, "SNACKS"},
		{"Condiments", Condiments, "CONDIMENTS"},
		{"Bakery", Bakery, "BAKERY"},
		{"Baked Goods", BakedGoods, "BAKED GOODS"},
		{"Canned Goods", CannedGoods, "CANNED GOODS"},
		{"Frozen Foods", FrozenFoods, "FROZEN FOODS"},
		{"Personal Care", PersonalCare, "PERSONAL CARE"},
		{"Household Supplies", HouseholdSupplies, "HOUSEHOLD SUPPLIES"},
		{"Pet Care", PetCare, "PET CARE"},
		{"Baby Items", BabyItems, "BABY ITEMS"},
		{"Others", Others, "OTHERS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.category.String(); actual != tt.expected {
				t.Errorf("ItemCategory.String() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestItemCategory_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expected  ItemCategory
		expectErr bool
	}{
		{"Valid Fruits", `"FRUITS"`, Fruits, false},
		{"Valid Vegetables", `"VEGETABLES"`, Vegetables, false},
		{"Valid Meat", `"MEAT"`, Meat, false},
		{"Valid Dairy", `"DAIRY"`, Dairy, false},
		{"Valid Baked Goods", `"BAKED GOODS"`, BakedGoods, false},
		{"Valid Others", `"OTHERS"`, Others, false},
		{"Invalid value", `"INVALID"`, ItemCategory{}, true},
		{"Invalid JSON", `invalid`, ItemCategory{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var category ItemCategory
			err := json.Unmarshal([]byte(tt.jsonData), &category)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && category != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, expected %v", category, tt.expected)
			}
		})
	}
}

func TestItemCategory_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		category ItemCategory
		expected string
	}{
		{"Fruits", Fruits, `"FRUITS"`},
		{"Vegetables", Vegetables, `"VEGETABLES"`},
		{"Meat", Meat, `"MEAT"`},
		{"Baked Goods", BakedGoods, `"BAKED GOODS"`},
		{"Others", Others, `"OTHERS"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.category)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %v, expected %v", string(data), tt.expected)
			}
		})
	}
}

func TestItemCategory_Scan(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  ItemCategory
		expectErr bool
	}{
		{"Valid string", "FRUITS", ItemCategory{value: "FRUITS"}, false},
		{"Another valid string", "MEAT", ItemCategory{value: "MEAT"}, false},
		{"Nil value", nil, ItemCategory{}, true},
		{"Invalid type", 123, ItemCategory{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var category ItemCategory
			err := category.Scan(tt.value)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && category.value != tt.expected.value {
				t.Errorf("Scan() = %v, expected %v", category.value, tt.expected.value)
			}
		})
	}
}

func TestItemCategory_Value(t *testing.T) {
	tests := []struct {
		name     string
		category ItemCategory
		expected driver.Value
	}{
		{"Fruits", Fruits, "FRUITS"},
		{"Vegetables", Vegetables, "VEGETABLES"},
		{"Meat", Meat, "MEAT"},
		{"Others", Others, "OTHERS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.category.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if val != tt.expected {
				t.Errorf("Value() = %v, expected %v", val, tt.expected)
			}
		})
	}
}
