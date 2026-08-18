package mapper

import (
	"reflect"
	"testing"

	"github.com/zouipo/yumsday/backend/internal/dto"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

/*** DATA ***/

var (
	ing = &model.Ingredient{
		ID:       1,
		Quantity: new(2.0),
	}

	item1Complete = &model.Item{
		ID:                 1,
		Name:               "Flour",
		Description:        new("All-purpose flour"),
		AverageMarketPrice: new(2.50),
		UnitType:           enum.Weight,
		ItemCategory:       *itemCategory1,
		GroupID:            1,
		Ingredients: []model.Ingredient{
			*ing,
		},
	}

	itemNoCategory = &model.Item{
		ID:                 3,
		Name:               "Chicken Breast",
		Description:        new("Boneless skinless chicken breast"),
		AverageMarketPrice: new(8.50),
		UnitType:           enum.Weight,
		GroupID:            2,
	}

	itemNoCategoryDto = &dto.ItemDto{
		ID:                 3,
		Name:               "Chicken Breast",
		Description:        new("Boneless skinless chicken breast"),
		AverageMarketPrice: new(8.50),
		UnitType:           enum.Weight,
		GroupID:            2,
	}

	itemEmptyFields = &model.Item{
		ID:   4,
		Name: "Empty item",
	}

	itemEmptyFieldsDto = &dto.ItemDto{
		ID:   4,
		Name: "Empty item",
	}
)

/*** TESTS ***/
func TestToItemDto(t *testing.T) {
	tests := []struct {
		name        string
		item        *model.Item
		expectedDto *dto.ItemDto
	}{
		{
			name:        "Item 1 ingredients",
			item:        item1Complete,
			expectedDto: itemDto1,
		},
		{
			name:        "Item 2 no ingredient",
			item:        item2,
			expectedDto: itemDto2,
		},
		{
			name:        "Item 3 no category",
			item:        itemNoCategory,
			expectedDto: itemNoCategoryDto,
		},
		{
			name:        "Item 4 empty fields",
			item:        itemNoCategory,
			expectedDto: itemNoCategoryDto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualDto := ToItemDto(tt.item)

			if !reflect.DeepEqual(actualDto, tt.expectedDto) {
				t.Errorf("Actual DTO %v mismatched the expected dto %v", actualDto, tt.expectedDto)
			}
		})
	}
}

func TestToItem(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dto.ItemDto
		expected *model.Item
	}{
		{
			name:     "Item 1",
			dto:      itemDto1,
			expected: item1,
		},
		{
			name:     "Item 2",
			dto:      itemDto2,
			expected: item2,
		},
		{
			name:     "Item 3",
			dto:      itemNoCategoryDto,
			expected: itemNoCategory,
		},
		{
			name:     "Item 2",
			dto:      itemEmptyFieldsDto,
			expected: itemEmptyFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToItem(tt.dto)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Actual item category %v mismatched the expected %v", actual, tt.expected)
			}
		})
	}
}
