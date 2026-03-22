package model

type Ingredient struct {
	ID       int64    `json:"id"`
	Quantity *float64 `json:"quantity"`
	RecipeID int64    `json:"recipe_id"`
	ItemID   int64    `json:"item_id"`
	UnitID   *int64   `json:"unit_id"`
}
