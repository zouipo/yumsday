package dto

import "github.com/zouipo/yumsday/backend/internal/model/enum"

type ItemDto struct {
	ID                 int64           `json:"id"`
	Name               string          `json:"name"`
	Description        *string         `json:"description"`
	AverageMarketPrice *float64        `json:"average_market_price"`
	UnitType           enum.UnitType   `json:"unit_type"`
	GroupID            int64           `json:"group_id"`
	ItemCategory       ItemCategoryDto `json:"item_category"`
}
