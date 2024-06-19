package models

import "time"

type Campaign struct {
	ID        int       `db:"id" json:"id"`
	GroupID   string    `db:"group_id" json:"group_id"`
	OrgID     string    `db:"org_id" json:"org_id"`
	Name      string    `db:"name" json:"name"`
	Status    string    `db:"status" json:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}
