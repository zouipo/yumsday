package model

type Grocery struct {
	ID             int64   `json:"id"`
	QuantityBought float64 `json:"quantity_bought"`
	UserQuantity   float64 `json:"user_quantity"`
	ItemID         int64   `json:"item_id"`
	UnitID         int64   `json:"unit_id"`
	GroupID        int64   `json:"group_id"`
}
