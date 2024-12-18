//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
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
	Limit  int
}

type CampaignDataToEndCondition struct {
	End    time.Time
	Status []string
	Limit  int
}

type CampaignCondition struct {
	CampaignID int
	Status     string
}

type CampaignRepository interface {
	// GetCampaignToStart 配信開始するキャンペーン情報を取得する
	GetCampaignToStart(ctx context.Context, args *CampaignToStartCondition) ([]*models.Campaign, error)
	// GetCampaignToEnd 配信が終了するキャンペーン情報を取得する
	GetCampaignToEnd(ctx context.Context, args *CampaignDataToEndCondition) ([]*models.Campaign, error)
	// UpdateStatus キャンペーン情報のステータス更新(status)更新
	UpdateStatus(ctx context.Context, tx Transaction, campaign *UpdateCondition) (int, error)
	// 配信操作するキャンペーン情報を取得する
	GetDeliveryToStart(ctx context.Context, tx Transaction, args *CampaignCondition) (*models.Campaign, error)
	// キャンペーンに紐づくクリエイティブの配信レートやスキップオフセットを取得する
	GetCampaignCreative(ctx context.Context, tx Transaction, args *CampaignCondition) ([]*models.CampaignCreative, error)
	// groupIDに紐づく配信中のキャンペーン数を取得する
	GetDeliveryCampaignCountByGroupID(ctx context.Context, groupID int) (int, error)
}
