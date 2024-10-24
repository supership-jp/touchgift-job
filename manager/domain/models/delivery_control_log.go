package models

type CampaignCacheLog struct {
	TraceID     string `json:"trace_id"`
	Time        string `json:"time"`
	Version     int    `json:"version"`
	Event       string `json:"event"`
	EventDetail string `json:"event_detail"`
	Action      string `json:"action"` // PUT or DELETE or NONE
	OrgCode     string `json:"org_code"`
	Source      string `json:"source"` // touchgift-delivery-manager
	ID          string `json:"id"`
	GroupID     string `json:"group_id"`
}

type CreativeCacheLog struct {
	ID               string   `json:"id"`
	Link             string   `json:"link"`
	URL              string   `json:"url"`
	Width            float32  `json:"width"`
	Height           float32  `json:"height"`
	Type             string   `json:"type"`
	Extension        string   `json:"extension"`
	Duration         *int     `json:"duration"`
	EndCardUrl       *string  `json:"endcard_url"`
	EndCardWidth     *float32 `json:"endcard_width"`
	EndCardHeight    *float32 `json:"endcard_height"`
	EndCardExtension *string  `json:"endcard_extension"`
	EndCardLink      *string  `json:"endcard_link"`
	Action           string   `json:"action"` // PUT or DELETE
}

type DeliveryCacheLog struct {
	ID         string `json:"id"`
	GroupID    int    `json:"group_id"`
	StoreID    string `json:"store_id"`
	Action     string `json:"action"` // PUT or DELETE
	OrgCode    string `json:"org_code"`
	CampaignID int    `json:"campaign_id"`
}
