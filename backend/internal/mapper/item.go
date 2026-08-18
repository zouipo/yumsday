package mapper

import (
	"github.com/zouipo/yumsday/backend/internal/dto"
	"github.com/zouipo/yumsday/backend/internal/model"
)

func ToItemDto(item *model.Item) *dto.ItemDto {
	return &dto.ItemDto{
		ID:                 item.ID,
		Name:               item.Name,
		Description:        item.Description,
		AverageMarketPrice: item.AverageMarketPrice,
		UnitType:           item.UnitType,
		GroupID:            item.GroupID,
		ItemCategory:       *ToItemCategoryDto(&item.ItemCategory),
	}
}

func ToItem(dto *dto.ItemDto) *model.Item {
	return &model.Item{
		ID:                 dto.ID,
		Name:               dto.Name,
		Description:        dto.Description,
		AverageMarketPrice: dto.AverageMarketPrice,
		UnitType:           dto.UnitType,
		GroupID:            dto.GroupID,
		ItemCategory:       *ToItemCategory(&dto.ItemCategory),
	}
}
