package models

import (
	"time"
)

// CampaignLog
// 同期対象配信データログのCampaign情報
type CampaignLog struct {
	ID        int    `json:"id"`
	OrgCode   string `json:"org_code"`
	Event     string `json:"event"`
	Creatives []CreativeLog
}

type CreativeLog struct {
	ID      int    `json:"id"`
	OrgCode string `json:"org_code"`
	Event   string `json:"event"`
}

// DeliveryOperationLog
// 同期対象配信データログ
type DeliveryOperationLog struct {
	Time         time.Time     `json:"time"`
	Type         string        `json:"type"`
	RequestID    string        `json:"request_id"`
	CampaignLogs []CampaignLog `json:"campaigns,omitempty"`
}
