package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type ContentsByCampaignIDCondition struct {
	CampaignID int
	Limit      int
}

type ContentsRepository interface {
	// GetContentsByCampaignID キャンペーンIDからコンテンツデータを取得する
	GetContentsByCampaignID(ctx context.Context, tx Transaction, args *ContentsByCampaignIDCondition) ([]*models.Contents, error)
}
