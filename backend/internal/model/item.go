package model

import "github.com/zouipo/yumsday/backend/internal/model/enum"

type Item struct {
	ID                 int64         `json:"id"`
	Name               string        `json:"name"`
	Description        *string       `json:"description"`
	AverageMarketPrice *float64      `json:"average_market_price"`
	UnitType           enum.UnitType `json:"unit_type"`
	ItemCategoryID     int64         `json:"item_category_id"`
	GroupID            int64         `json:"group_id"`
	Ingredients        []Ingredient  `json:"ingredients"`
}
