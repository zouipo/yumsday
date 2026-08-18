package mapper

import (
	"github.com/zouipo/yumsday/backend/internal/dto"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

var (
	itemCategory1 = &model.ItemCategory{
		ID:      1,
		Name:    "GRAINS AND PASTA",
		GroupID: 1,
	}

	itemCategory2 = &model.ItemCategory{
		ID:      2,
		Name:    "VEGETABLES",
		GroupID: 2,
	}

	itemCategoryDto1 = &dto.ItemCategoryDto{
		ID:      1,
		Name:    "GRAINS AND PASTA",
		GroupID: 1,
	}

	itemCategoryDto2 = &dto.ItemCategoryDto{
		ID:      2,
		Name:    "VEGETABLES",
		GroupID: 2,
	}

	item1 = &model.Item{
		ID:                 1,
		Name:               "Flour",
		Description:        new("All-purpose flour"),
		AverageMarketPrice: new(2.50),
		UnitType:           enum.Weight,
		ItemCategory:       *itemCategory1,
		GroupID:            1,
	}

	item2 = &model.Item{
		ID:                 2,
		Name:               "Onions",
		Description:        new("Yellow onions"),
		AverageMarketPrice: new(1.50),
		UnitType:           enum.Weight,
		ItemCategory:       *itemCategory2,
		GroupID:            2,
	}

	itemDto1 = &dto.ItemDto{
		ID:                 1,
		Name:               "Flour",
		Description:        new("All-purpose flour"),
		AverageMarketPrice: new(2.50),
		UnitType:           enum.Weight,
		ItemCategory:       *itemCategoryDto1,
		GroupID:            1,
	}

	itemDto2 = &dto.ItemDto{
		ID:                 2,
		Name:               "Onions",
		Description:        new("Yellow onions"),
		AverageMarketPrice: new(1.50),
		UnitType:           enum.Weight,
		ItemCategory:       *itemCategoryDto2,
		GroupID:            2,
	}
)
