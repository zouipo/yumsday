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
)

/*** TESTS ***/
func TestToItemCategoryDto(t *testing.T) {
	tests := []struct {
		name         string
		itemCategory *model.ItemCategory
		expectedDto  *dto.ItemCategoryDto
	}{
		{
			name:         "ItemCategory 1",
			itemCategory: ItemCategory1Complete,
			expectedDto:  itemCategoryDto1,
		},
		{
			name:         "ItemCategory 2",
			itemCategory: itemCategory2,
			expectedDto:  itemCategoryDto2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualDto := ToItemCategoryDto(tt.itemCategory)

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
			dto:      itemCategoryDto1,
			expected: itemCategory1,
		},
		{
			name:     "ItemCategory 2",
			dto:      itemCategoryDto2,
			expected: itemCategory2,
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
