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
	// GetCouponsByCampaignID  キャンペーンIDからクーポンデータを取得する
	GetCouponsByCampaignID(ctx context.Context, tx Transaction, args *ContentsByCampaignIDCondition) ([]*models.Coupon, error)
	// GetGimmickURLByCampaignID キャンペーンIDからギミックURLを取得する
	GetGimmickURLByCampaignID(ctx context.Context, tx Transaction, campaignID *ContentsByCampaignIDCondition) (*string, error)
}

// ContentsHelper gimmick_urlとクーポン一覧からコンテンツを作成する
type ContentsHelper interface {
	GenerateContents(ctx context.Context, args *GenerateContentsCondition) ([]*models.Contents, error)
}
