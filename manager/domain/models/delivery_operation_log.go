package models

import (
	"time"
)

// CampaignLog
// 同期対象配信データログのCampaign情報
type CampaignLog struct {
	ID               int    `json:"id"`
	Budget           int    `json:"budget"`
	OrganizationCode string `json:"organization_code"`
	OriginID         string `json:"origin_id"`
	Event            string `json:"event"`
	Service          string `json:"service"`
	AdvertiserID     int    `json:"advertiser_id"`
}

// DeliveryOperationLog
// 同期対象配信データログ
type DeliveryOperationLog struct {
	Time         time.Time     `json:"time"`
	Type         string        `json:"type"`
	RequestID    string        `json:"request_id"`
	CampaignLogs []CampaignLog `json:"campaigns,omitempty"`
}
