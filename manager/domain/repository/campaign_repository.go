package repository

import (
	"context"
	"time"
	"touchgift-job-manager/domain/models"
)

type UpdateCondition struct {
	CampaignID int
	Status     string
	UpdatedAt  time.Time
}

type CampaignToStartCondition struct {
	To     time.Time
	Status string
	Limit  int
}

type CampaignIDCondition struct {
	GroupID int
}

type CampaignRepository interface {
	// GetCampaignDataToStart 開始するキャンペーン情報を取得する
	GetCampaignToStart(ctx context.Context, args *CampaignToStartCondition) ([]*models.Campaign, error)
	// GetTouchPointDataByTouchPointID touch_pointIDからtouch_pointテーブルを取得する
	GetTouchPointByGroupID(ctx context.Context, args *CampaignIDCondition) ([]*models.TouchPoint, error)
}
