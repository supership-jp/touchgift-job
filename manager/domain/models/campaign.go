package models

import (
	"database/sql"
	"strconv"
	"time"
)

// Campaign RDBから取得した配信開始・終了に必要なデータ
type Campaign struct {
	ID                      int          `db:"id" json:"id"`
	StartAt                 time.Time    `db:"start_at" json:"start_at"`
	EndAt                   sql.NullTime `db:"end_at" json:"end_at"`
	UpdatedAt               time.Time    `db:"updated_at"`
	GroupID                 int          `db:"group_id" json:"group_id"`
	OrgCode                 string       `db:"org_code" json:"org_code"`
	DailyCouponLimitPerUser int          `db:"daily_coupon_limit_per_user" json:"daily_coupon_limit_per_user"`
	Status                  string       `db:"status" json:"status"`
}

func (c *Campaign) CreateDeliveryDataCampaign(creative []*Creative) *DeliveryDataCampaign {
	creatives := make([]CampaignCreative, 0, len(creative))
	for _, c := range creative {
		creatives = append(creatives, *c.CreateCampaignCreative())
	}
	return &DeliveryDataCampaign{
		ID:         strconv.Itoa(c.ID),
		GroupID:    c.GroupID,
		OrgCode:    c.OrgCode,
		DailyLimit: c.DailyCouponLimitPerUser,
		Status:     c.Status,
		Creatives:  creatives,
	}
}

type CampaignCreative struct {
	ID         int `json:"id"`
	Rate       int `json:"rate"`
	SkipOffset int `json:"skip_offset"`
}
