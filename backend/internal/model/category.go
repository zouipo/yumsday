package model

type Category struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	UserGroupID int64    `json:"user_group_id"`
	Recipes     []Recipe `json:"recipes"`
}
