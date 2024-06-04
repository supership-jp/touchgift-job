//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"time"
	"touchgift-job-manager/domain/models"
)

type DeliveryOperation interface {
	// キャンペーンのログを処理する
	Process(ctx context.Context, current time.Time, campaignLog *models.CampaignLog) error
}

// TODO: 一旦mock的に作成
type deliveryOperation struct{}

func NewDeliveryOperation() DeliveryOperation { return nil }
func (op *deliveryOperation) Process(ctx context.Context, current time.Time, campaignLog *models.CampaignLog) error {
	return nil
}
