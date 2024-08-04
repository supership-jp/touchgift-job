//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type ContentByCampaignIDCondition struct {
	CampaignID int
}

type GenerateContentCondition struct {
	Coupons    []*models.Coupon
	GimmickURL *string
}

type ContentRepository interface {
	// GetCouponsByCampaignID  キャンペーンIDからクーポンデータを取得する
	GetCouponsByCampaignID(ctx context.Context, tx Transaction, args *ContentByCampaignIDCondition) ([]*models.Coupon, error)
	// GetGimmickURLByCampaignID キャンペーンIDからギミックURLを取得する
	GetGimmicksByCampaignID(ctx context.Context, tx Transaction, campaignID *ContentByCampaignIDCondition) (*string, *string, error)
}

// ContentsHelper gimmick_urlとクーポン一覧からコンテンツを作成する
type ContentHelper interface {
	GenerateContent(ctx context.Context, args *GenerateContentCondition) ([]*models.Content, error)
}
