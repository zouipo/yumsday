package model

import "time"

type Group struct {
	ID         int64            `json:"id"`
	Name       string           `json:"name"`
	ImageURL   *string          `json:"image_url"`
	CreatedAt  time.Time        `json:"created_at"`
	Groceries  []Grocery        `json:"groceries"`
	Categories []RecipeCategory `json:"categories"`
	Recipes    []Recipe         `json:"recipes"`
	Items      []Item           `json:"items"`
	Dishes     []Dish           `json:"dishes"`
	Users      []struct {
		UserID   int64     `json:"user_id"`
		Admin    bool      `json:"admin"`
		JoinedAt time.Time `json:"joined_at"`
	} `json:"users"`
}
