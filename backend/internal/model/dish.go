package model

import "time"

type Dish struct {
	ID          int64     `json:"id"`
	Portion     int       `json:"portion"`
	Bought      bool      `json:"bought"`
	Datetime    time.Time `json:"datetime"`
	UserGroupID int64     `json:"user_group_id"`
	RecipeID    int64     `json:"recipe_id"`
}
