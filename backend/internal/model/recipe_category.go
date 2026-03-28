package model

type RecipeCategory struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	GroupID int64    `json:"group_id"`
	Recipes []Recipe `json:"recipes"`
}
