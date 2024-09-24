package models

import "strconv"

// Dynamoに入れるデータ構造体はここに定義していく

type DeliveryDataCampaign struct {
	ID         string              `json:"id"`
	GroupID    string              `json:"group_id"`
	OrgCode    string              `json:"org_code"`
	DailyLimit int                 `json:"daily_limit"`
	Creatives  []*CampaignCreative `json:"creatives,omitempty"`
	Status     string              `json:"status"`
}

func (d *DeliveryDataCampaign) CreateCampaign() *Campaign {
	ID, _ := strconv.Atoi(d.ID)
	groupID, _ := strconv.Atoi(d.GroupID)
	return &Campaign{
		ID:                      ID,
		GroupID:                 groupID,
		OrgCode:                 d.OrgCode,
		DailyCouponLimitPerUser: d.DailyLimit,
		Status:                  d.Status,
	}
}

type DeliveryTouchPoint struct {
	GroupID int    `json:"group_id"`
	ID      string `json:"id"`
}

// DeliveryDataCreative dynamo用に整形するための構造体(クリエイティブ用)
type DeliveryDataCreative struct {
	ID               string   `json:"id"`
	Link             string   `json:"link,omitempty"`
	URL              string   `json:"url"`
	TTL              int64    `json:"ttl"`
	Width            float32  `json:"width"`
	Height           float32  `json:"height"`
	Type             string   `json:"type"`
	Extension        string   `json:"extension"`
	Duration         *int     `json:"duration,omitempty"`
	SkipOffset       *int     `json:"skip_offset,omitempty"`
	EndCardURL       *string  `json:"end_card_url,omitempty"`
	EndCardWidth     *float32 `json:"end_card_width,omitempty"`
	EndCardHeight    *float32 `json:"end_card_height,omitempty"`
	EndCardExtension *string  `json:"end_card_extension,omitempty"`
	EndCardLink      *string  `json:"end_card_link,omitempty"`
}

type DeliveryDataContent struct {
	CampaignID string               `json:"campaign_id"`
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
