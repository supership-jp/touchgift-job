package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type CreativeByCampaignIDCondition struct {
	CampaignID int
	Limit      int
}

type CreativeRepository interface {
	// GetCreativeByCampaignID クリエイティブの取得をキャンペーンIDから行う
	GetCreativeByCampaignID(ctx context.Context, args *CreativeByCampaignIDCondition) ([]*models.Creative, error)
}
