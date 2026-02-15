package enums

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestRecipeCategoryEnum_String(t *testing.T) {
	tests := []struct {
		name     string
		category RecipeCategoryEnum
		expected string
	}{
		{"Dessert", Dessert, "DESSERT"},
		{"Salad", Salad, "SALAD"},
		{"Main Course", MainCourse, "MAIN COURSE"},
		{"Soup", Soup, "SOUP"},
		{"Breakfast", Breakfast, "BREAKFAST"},
		{"Brunch", Brunch, "BRUNCH"},
		{"Starter", Starter, "STARTER"},
		{"Sauce", Sauce, "SAUCE"},
		{"Snack", Snack, "SNACK"},
		{"Beverage", Beverage, "BEVERAGE"},
		{"Vegan", Vegan, "VEGAN"},
		{"Vegetarian", Vegetarian, "VEGETARIAN"},
		{"Gluten Free", GlutenFree, "GLUTEN FREE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.category.String(); actual != tt.expected {
				t.Errorf("RecipeCategoryEnum.String() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestRecipeCategoryEnum_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expected  RecipeCategoryEnum
		expectErr bool
	}{
		{"Valid Dessert", `"DESSERT"`, Dessert, false},
		{"Valid Salad", `"SALAD"`, Salad, false},
		{"Valid Main Course", `"MAIN COURSE"`, MainCourse, false},
		{"Valid Soup", `"SOUP"`, Soup, false},
		{"Valid Vegan", `"VEGAN"`, Vegan, false},
		{"Valid Vegetarian", `"VEGETARIAN"`, Vegetarian, false},
		{"Valid Gluten Free", `"GLUTEN FREE"`, GlutenFree, false},
		{"Invalid value", `"INVALID"`, RecipeCategoryEnum{}, true},
		{"Invalid JSON", `invalid`, RecipeCategoryEnum{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var category RecipeCategoryEnum
			err := json.Unmarshal([]byte(tt.jsonData), &category)

			if tt.expectErr && err == nil {
				t.Error("Expected error but actual nil")
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

func TestRecipeCategoryEnum_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		category RecipeCategoryEnum
		expected string
	}{
		{"Dessert", Dessert, `"DESSERT"`},
		{"Salad", Salad, `"SALAD"`},
		{"Main Course", MainCourse, `"MAIN COURSE"`},
		{"Soup", Soup, `"SOUP"`},
		{"Vegan", Vegan, `"VEGAN"`},
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

func TestRecipeCategoryEnum_Scan(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  RecipeCategoryEnum
		expectErr bool
	}{
		{"Valid string", "DESSERT", RecipeCategoryEnum{value: "DESSERT"}, false},
		{"Another valid string", "SALAD", RecipeCategoryEnum{value: "SALAD"}, false},
		{"Nil value", nil, RecipeCategoryEnum{}, true},
		{"Invalid type", 123, RecipeCategoryEnum{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var category RecipeCategoryEnum
			err := category.Scan(tt.value)

			if tt.expectErr && err == nil {
				t.Error("Expected error but actual nil")
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

func TestRecipeCategoryEnum_Value(t *testing.T) {
	tests := []struct {
		name     string
		category RecipeCategoryEnum
		expected driver.Value
	}{
		{"Dessert", Dessert, "DESSERT"},
		{"Salad", Salad, "SALAD"},
		{"Main Course", MainCourse, "MAIN COURSE"},
		{"Soup", Soup, "SOUP"},
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
