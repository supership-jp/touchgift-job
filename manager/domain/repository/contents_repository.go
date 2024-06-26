package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type ContentsByCampaignIDCondition struct {
	CampaignID int
}

type GenerateContentsCondition struct {
	Coupons    []*models.Coupon
	GimmickURL *string
}

type ContentsRepository interface {
	// GetContentsByCampaignID キャンペーンIDからコンテンツデータを取得する
	GetCouponsByCampaignID(ctx context.Context, tx Transaction, args *ContentsByCampaignIDCondition) ([]*models.Coupon, error)
	GetGimmickURLByCampaignID(ctx context.Context, tx Transaction, campaignID *ContentsByCampaignIDCondition) (*string, error)
}

// ContentsHelper gimmick_urlとクーポン一覧からコンテンツを作成する
type ContentsHelper interface {
	GenerateContents(ctx context.Context, args *GenerateContentsCondition) ([]*models.Contents, error)
}
