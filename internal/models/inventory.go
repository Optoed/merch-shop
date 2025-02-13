package models

type InventoryItem struct {
	ID       int    `json:"id" db:"id"`
	UserID   int    `json:"user_id" db:"user_id"`
	ItemName string `json:"item_name" db:"item_name"`
	Count    int    `json:"count" db:"count"`
}
