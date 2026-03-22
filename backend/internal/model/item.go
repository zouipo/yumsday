package model

import "github.com/zouipo/yumsday/backend/internal/model/enum"

type Item struct {
	ID                 int64              `json:"id"`
	Name               string             `json:"name"`
	Description        *string            `json:"description"`
	AverageMarketPrice *float64           `json:"average_market_price"`
	ItemCategory       *enum.ItemCategory `json:"item_category"`
	UnitType           enum.UnitType      `json:"unit_type"`
}
