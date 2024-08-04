package models

type DeliveryControlLog struct {
	TraceID        string `json:"trace_id"`
	Time           string `json:"time"`
	Version        int    `json:"version"`
	Event          string `json:"event"`
	EventDetail    string `json:"event_detail"`
	CacheOperation string `json:"cache_operation"` // PUT or DELETE or NONE
	OrgCode        string `json:"org_code"`
	Source         string `json:"source"` // touchgift-delivery-manager
	CampaignID     int    `json:"campaign_id"`
	DeliveryType   string `json:"delivery_type"` // touchgift
	PriceType      string `json:"price_type"`
}
