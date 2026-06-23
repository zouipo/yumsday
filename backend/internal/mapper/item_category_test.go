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
	ItemCategory1Complete = &model.ItemCategory{
		ID:      1,
		Name:    "GRAINS AND PASTA",
		GroupID: 1,
		Items: []model.Item{
			{
				ID:                 1,
				Name:               "Flour",
				Description:        new("All-purpose flour"),
				AverageMarketPrice: new(2.50),
				UnitType:           enum.Weight,
				GroupID:            1,
			},
			{
				ID:                 2,
				Name:               "Rice",
				Description:        new("White rice"),
				AverageMarketPrice: new(1.80),
				UnitType:           enum.Weight,
				GroupID:            1,
			},
		},
	}

	ItemCategory1Light = &model.ItemCategory{
		ID:      1,
		Name:    "GRAINS AND PASTA",
		GroupID: 1,
	}

	ItemCategory2 = &model.ItemCategory{
		ID:      2,
		Name:    "VEGETABLES",
		GroupID: 2,
	}

	ItemCategoryDto1 = &dto.ItemCategoryDto{
		ID:      1,
		Name:    "GRAINS AND PASTA",
		GroupID: 1,
	}

	ItemCategoryDto2 = &dto.ItemCategoryDto{
		ID:      2,
		Name:    "VEGETABLES",
		GroupID: 2,
	}
)

/*** TESTS ***/
func TestToItemCategoryDto(t *testing.T) {
	tests := []struct {
		name        string
		item        *model.ItemCategory
		expectedDto *dto.ItemCategoryDto
	}{
		{
			name:        "ItemCategory 1",
			item:        ItemCategory1Complete,
			expectedDto: ItemCategoryDto1,
		},
		{
			name:        "ItemCategory 2",
			item:        ItemCategory2,
			expectedDto: ItemCategoryDto2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualDto := ToItemCategoryDto(tt.item)

			if !reflect.DeepEqual(actualDto, tt.expectedDto) {
				t.Errorf("Actual DTO %v mismatched the expected dto %v", actualDto, tt.expectedDto)
			}
		})
	}
}

func TestToItemCategory(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dto.ItemCategoryDto
		expected *model.ItemCategory
	}{
		{
			name:     "ItemCategory 1",
			dto:      ItemCategoryDto1,
			expected: ItemCategory1Light,
		},
		{
			name:     "ItemCategory 2",
			dto:      ItemCategoryDto2,
			expected: ItemCategory2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToItemCategory(tt.dto)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Actual item category %v mismatched the expected %v", actual, tt.expected)
			}
		})
	}
}
