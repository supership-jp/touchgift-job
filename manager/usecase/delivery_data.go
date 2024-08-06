//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type DeliveryData interface {
	// Put 配信データを登録する
	Put(ctx context.Context, campaign *models.Campaign, creatives []*models.Creative, content *models.DeliveryDataContent, touchPoints []*models.DeliveryTouchPoint) error
	// Delete 配信データを1件削除
	Delete(ctx context.Context, campaignID string) error
}

type deliveryData struct {
	logger                   Logger
	touchPointDataRepository repository.DeliveryDataTouchPointRepository
	campaignDataRepository   repository.DeliveryDataCampaignRepository
	contentDataRepository    repository.DeliveryDataContentRepository
	creativeDataRepository   repository.DeliveryDataCreativeRepository
}

func NewDeliveryData(
	logger Logger,
	touchPointDataRepository repository.DeliveryDataTouchPointRepository,
	campaignDataRepository repository.DeliveryDataCampaignRepository,
	contentDataRepository repository.DeliveryDataContentRepository,
	creativeDataRepository repository.DeliveryDataCreativeRepository,
) DeliveryData {
	return &deliveryData{
		logger:                   logger,
		touchPointDataRepository: touchPointDataRepository,
		campaignDataRepository:   campaignDataRepository,
		contentDataRepository:    contentDataRepository,
		creativeDataRepository:   creativeDataRepository,
	}
}

func (d *deliveryData) Put(ctx context.Context,
	campaign *models.Campaign, creatives []*models.Creative, content *models.DeliveryDataContent, touchPoints []*models.DeliveryTouchPoint,
) error {
	err := d.campaignDataRepository.Put(ctx, campaign.CreateDeliveryDataCampaign(creatives))
	if err != nil {
		return err
	}

	for _, tp := range touchPoints {
		err := d.touchPointDataRepository.Put(ctx, tp)
		if err != nil {
			return err
		}
	}

	for _, creative := range creatives {
		err := d.creativeDataRepository.Put(ctx, creative.CreateDeliveryDataCreative(campaign.ID))
		if err != nil {
			return err
		}
	}

	err = d.contentDataRepository.Put(ctx, content)
	if err != nil {
		return err
	}
	return nil
}

func (d *deliveryData) Delete(ctx context.Context, campaignID string) error {
	if err := d.campaignDataRepository.Delete(ctx, &campaignID); err != nil {
		return err
	}
	return nil
}
