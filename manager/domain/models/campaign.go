package models

import (
	"database/sql"
	"time"
)

// Campaign RDBから取得した配信開始・終了に必要なデータ
type Campaign struct {
	ID        int          `db:"id" json:"id"`
	StartAt   time.Time    `db:"start_at" json:"start_at"`
	EndAt     sql.NullTime `db:"end_at" json:"end_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	GroupID   int          `db:"group_id" json:"group_id"`
	OrgID     string       `db:"org_id" json:"org_id"`
	Name      string       `db:"name" json:"name"`
	Status    string       `db:"status" json:"status"`
}
