package model

import "github.com/zouipo/yumsday/backend/internal/model/enum"

type Unit struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Factor   float64       `json:"factor"`
	UnitType enum.UnitType `json:"unit_type"`
}
