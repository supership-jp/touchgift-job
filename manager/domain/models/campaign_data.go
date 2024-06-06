package models

import (
	"database/sql"
	"time"
)

// CampaignData TODO: diも入れる
type CampaignData struct {
	ID        int          `db:"id" json:"id"`
	GroupID   string       `db:"group_id" json:"group_id"`
	OrgID     int          `db:"org_id" json:"org_id"`
	Status    string       `db:"status" json:"status"`
	StartAt   time.Time    `db:"start_at" json:"start_at"`
	EndAt     sql.NullTime `db:"end_at" json:"end_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
}
