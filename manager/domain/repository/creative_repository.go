//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
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
	GetCreativeByCampaignID(ctx context.Context, tx Transaction, args *CreativeByCampaignIDCondition) ([]*models.Creative, error)
	GetCreative(ctx context.Context, tx Transaction, args *CreativeCondition) ([]models.Creative, error)
}
