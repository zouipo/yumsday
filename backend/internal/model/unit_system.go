package model

type UnitSystem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Units []Unit `json:"units"`
}
