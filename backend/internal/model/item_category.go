package model

type ItemCategory struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	GroupID int64  `json:"group_id"`
	Items   []Item `json:"items"`
}
