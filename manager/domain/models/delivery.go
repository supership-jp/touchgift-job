package models

// Dynamoに入れるデータ構造体はここに定義していく

type DeliveryDataCampaign struct {
	ID      string `json:"id"`
	GroupID int    `json:"group_id"`
	OrgID   string `json:"org_id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

type DeliveryTouchPoint struct {
	TouchPointID string `json:"touch_point_id"`
	GroupID      int    `json:"group_id"`
}

// DeliveryDataCreative dyonamo用に整形するための構造体(クリエイティブ用)
type DeliveryDataCreative struct {
	ID               string   `json:"id"` // これはキャンペーンIDです
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
	CampaignID string               `json:"campaign_id"`
	Coupons    []DeliveryCouponData `json:"coupons"`
	GimmickURL *string              `json:"gimmick_url,omitempty"`
}

type DeliveryCouponData struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	ImageURL string  `json:"image_url"`
	Rate     float64 `json:"rate"`
}
