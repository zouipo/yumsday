package model

type GroceryList struct {
	ID             int64   `json:"id"`
	QuantityBought float64 `json:"quantity_bought"`
	UserQuantity   float64 `json:"user_quantity"`
	ItemID         int64   `json:"item_id"`
	UnitID         int64   `json:"unit_id"`
	UserGroupID    int64   `json:"user_group_id"`
}
