package models

import "strconv"

// Dynamoに入れるデータ構造体はここに定義していく

type DeliveryDataCampaign struct {
	ID         string `json:"id"`
	GroupID    int    `json:"group_id"`
	OrgCode    string `json:"org_code"`
	DailyLimit int    `json:"daily_limit"`
	Status     string `json:"status"`
}

func (d *DeliveryDataCampaign) CreateCampaign() *Campaign {
	ID, _ := strconv.Atoi(d.ID)
	return &Campaign{
		ID:                      ID,
		GroupID:                 d.GroupID,
		OrgCode:                 d.OrgCode,
		DailyCouponLimitPerUser: d.DailyLimit,
		Status:                  d.Status,
	}
}

type DeliveryTouchPoint struct {
	TouchPointID string `json:"touch_point_id"`
	GroupID      int    `json:"group_id"`
}

// DeliveryDataCreative dyonamo用に整形するための構造体(クリエイティブ用)
type DeliveryDataCreative struct {
	CampaignID       string   `json:"campaign_id"` // これはキャンペーンIDです
	Link             *string  `json:"link,omitempty"`
	URL              string   `json:"url"`
	TTL              int64    `json:"ttl"`
	Width            float32  `json:"width"`
	Height           float32  `json:"height"`
	Type             string   `json:"type"`
	Extension        string   `json:"extension"`
	Duration         *int     `json:"duration"`
	SkipOffset       *int     `json:"skip_offset"`
	EndCardURL       *string  `json:"end_card_url,omitempty"`
	EndCardWidth     *float32 `json:"end_card_width,omitempty"`
	EndCardHeight    *float32 `json:"end_card_height,omitempty"`
	EndCardExtension *string  `json:"end_card_extension,omitempty"`
	EndCardLink      *string  `json:"end_card_link,omitempty"`
}

type DeliveryDataContent struct {
	CampaignID int                  `json:"campaign_id"`
	Coupons    []DeliveryCouponData `json:"coupons"`
	Gimmicks   []Gimmick            `json:"gimmicks"`
}

type DeliveryCouponData struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	ImageURL string `json:"image_url"`
	Rate     int    `json:"rate"`
}
