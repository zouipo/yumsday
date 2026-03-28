package model

import "time"

type Dish struct {
	ID       int64     `json:"id"`
	Portion  int       `json:"portion"`
	Bought   bool      `json:"bought"`
	Datetime time.Time `json:"datetime"`
	GroupID  int64     `json:"group_id"`
	Recipes  []Recipe  `json:"recipes"`
}
