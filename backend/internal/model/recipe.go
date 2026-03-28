package model

import "time"

type Recipe struct {
	ID                 int64            `json:"id"`
	Name               string           `json:"name"`
	Description        *string          `json:"description"`
	ImageURL           *string          `json:"image_url"`
	OriginalLink       *string          `json:"original_link"`
	PreparationTimeMin *int             `json:"preparation_time_min"`
	CookingTimeMin     *int             `json:"cooking_time_min"`
	Servings           *int             `json:"servings"`
	Instructions       *string          `json:"instructions"`
	CreatedAt          time.Time        `json:"created_at"`
	Public             bool             `json:"public"`
	Comment            *string          `json:"comment"`
	GroupID            int64            `json:"group_id"`
	Categories         []RecipeCategory `json:"recipe_categories"`
	Ingredients        []Ingredient     `json:"ingredients"`
	Dishes             []Dish           `json:"dishes"`
}
