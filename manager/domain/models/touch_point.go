package models

type TouchPoint struct {
	GroupID      int    `db:"group_id" json:"group_id"`
	TouchPointID string `db:"touch_point_id" json:"touch_point_id"`
}
