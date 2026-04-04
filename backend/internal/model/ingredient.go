package model

type Ingredient struct {
	ID       int64    `json:"id"`
	Quantity *float64 `json:"quantity"`
	Item     Item     `json:"item"`
	RecipeID int64    `json:"recipe_id"`
	UnitID   int64    `json:"unit_id"`
}
