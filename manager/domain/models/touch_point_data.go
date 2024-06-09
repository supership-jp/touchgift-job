package models

type TouchPointData struct {
	GroupID      string `db:"group_id" json:"group_id"`
	TouchPointID string `db:"touch_point_id" json:"touch_point_id"`
}
