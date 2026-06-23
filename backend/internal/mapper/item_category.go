package mapper

import (
	"github.com/zouipo/yumsday/backend/internal/dto"
	"github.com/zouipo/yumsday/backend/internal/model"
)

func ToItemCategoryDto(itemCategory *model.ItemCategory) *dto.ItemCategoryDto {
	return &dto.ItemCategoryDto{
		ID:      itemCategory.ID,
		Name:    itemCategory.Name,
		GroupID: itemCategory.GroupID,
	}
}

func ToItemCategory(dto *dto.ItemCategoryDto) *model.ItemCategory {
	return &model.ItemCategory{
		ID:      dto.ID,
		Name:    dto.Name,
		GroupID: dto.GroupID,
	}
}
