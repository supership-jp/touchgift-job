package models

type TouchPoint struct {
	GroupID int    `db:"group_id" json:"group_id"`
	StoreID string `db:"store_id" json:"store_id"`
	ID      string `db:"id" json:"id"`
}
