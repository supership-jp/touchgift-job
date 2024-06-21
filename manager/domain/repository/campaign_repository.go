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

type CampaignRepository interface {
	// GetCampaignDataToStart 開始するキャンペーン情報を取得する
	GetCampaignToStart(ctx context.Context, tx Transaction, args *CampaignToStartCondition) ([]*models.Campaign, error)
}
