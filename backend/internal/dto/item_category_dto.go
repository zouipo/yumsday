package dto

type ItemCategoryDto struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	GroupID int64  `json:"group_id"`
}
