package repository

import (
	"context"
	"time"
	"touchgift-job-manager/domain/models"
)

//　RDSへの接続操作

type UpdateCondition struct {
	CampaignID int
	Status     string
	UpdatedAt  time.Time
}

type CampaignToStartCondition struct {
	To     time.Time
	Status string
}

type CampaignDataToEndCondition struct {
	End    time.Time
	Status string
}

type CampaignCondition struct {
	CampaignID int
	Status     string
}

type CampaignRepository interface {
	// GetCampaignToStart 配信開始するキャンペーン情報を取得する
	GetCampaignToStart(ctx context.Context, tx Transaction, args *CampaignToStartCondition) ([]*models.Campaign, error)
	//	配信終了するキャンペーンを取得する
	GetCampaignToEnd(ctx context.Context, tx Transaction, campaign *CampaignDataToEndCondition) ([]*models.Campaign, error)
	//// UpdateStatus ステータス更新(status)更新
	UpdateStatus(ctx context.Context, tx Transaction, campaign *UpdateCondition) (time.Time, error)
}
