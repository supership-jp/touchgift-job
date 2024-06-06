package models

type TouchPointData struct {
	GroupID      string `db:"" json:"group_id"`
	TouchPointID string `db:"" json:"touch_point_id"`
}
