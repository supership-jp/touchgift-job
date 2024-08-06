//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"strconv"
	"time"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"
)

type Delivery interface {
	// UpdateStatus 配信状態を更新する
	UpdateStatus(ctx context.Context, tx repository.Transaction, campaignId int, status string, updatedAt time.Time) (int, error)
	// DeliveryControlEvent 配信制御イベントを発行する
	DeliveryControlEvent(ctx context.Context, campaign *models.Campaign, beforeStatus string, afterStatus string, detail string)
	// StartOrSync データの同期or配信を開始する
	StartOrSync(ctx context.Context, tx repository.Transaction, campaign *models.Campaign) error
	// Stop 配信を停止する
	Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error
}

type delivery struct {
	logger               Logger
	monitor              *metrics.Monitor
	campaignRepository   repository.CampaignRepository
	creativeRepository   repository.CreativeRepository
	contentsRepository   repository.ContentRepository
	touchPointRepository repository.TouchPointRepository
	deliveryData         DeliveryData
	deliveryControlEvent DeliveryControlEvent
}

func NewDelivery(
	logger Logger,
	monitor *metrics.Monitor,
	campaignRepository repository.CampaignRepository,
	creativeRepository repository.CreativeRepository,
	contentsRepository repository.ContentRepository,
	touchPointRepository repository.TouchPointRepository,
	deliveryData DeliveryData,
	deliveryControlEvent DeliveryControlEvent,
) Delivery {
	return &delivery{
		logger:               logger,
		monitor:              monitor,
		campaignRepository:   campaignRepository,
		creativeRepository:   creativeRepository,
		contentsRepository:   contentsRepository,
		touchPointRepository: touchPointRepository,
		deliveryData:         deliveryData,
		deliveryControlEvent: deliveryControlEvent,
	}
}

// 配信開始または同期処理
func (d *delivery) StartOrSync(ctx context.Context, tx repository.Transaction, campaign *models.Campaign) error {
	campaignID := campaign.ID
	_, err := d.UpdateStatus(ctx, tx, campaignID, "started", campaign.StartAt)
	if err != nil {
		return err
	}

	// TODO: IDの型の取り扱いを考える
	condition := repository.ContentByCampaignIDCondition{
		CampaignID: campaignID,
	}
	// クリエイティブの取得
	creativeCondition := repository.CreativeByCampaignIDCondition{
		CampaignID: campaignID,
	}
	creatives, err := d.creativeRepository.GetCreativeByCampaignID(ctx, tx, &creativeCondition)
	if err != nil {
		return err
	}

	// TODO: コンテンツをそれぞれキャンペーンから取得してメモリに展開
	// ギミックURLの取得
	gimmickURL, gimmickCode, err := d.contentsRepository.GetGimmicksByCampaignID(ctx, tx, &condition)
	if err != nil {
		return err
	}
	//　クーポン一覧の取得
	coupons, err := d.contentsRepository.GetCouponsByCampaignID(ctx, tx, &condition)
	if err != nil {
		return err
	}
	deliveryCouponDatas := make([]models.DeliveryCouponData, 0, len(coupons))
	for _, coupon := range coupons {
		deliveryCouponData := coupon.CreateDeliveryCouponData()
		deliveryCouponDatas = append(deliveryCouponDatas, *deliveryCouponData)
	}
	// タッチポイントの取得
	touchPointCondition := &repository.TouchPointByGroupIDCondition{
		GroupID: campaign.GroupID,
		Limit:   1,
	}
	touchPoints, err := d.touchPointRepository.GetTouchPointByGroupID(ctx, tx, touchPointCondition)
	if err != nil {
		return err
	}
	// content作成
	content := &models.DeliveryDataContent{
		CampaignID: campaign.ID,
		Coupons:    deliveryCouponDatas,
		Gimmicks: []models.Gimmick{
			{
				URL:  *gimmickURL,
				Code: *gimmickCode,
			},
		},
	}
	// touchPoint作成
	touchPointDatas := make([]*models.DeliveryTouchPoint, 0, len(touchPoints))
	for _, touchPoint := range touchPoints {
		touchPointData := models.DeliveryTouchPoint{
			TouchPointID: touchPoint.TouchPointID,
			GroupID:      touchPoint.GroupID,
		}
		touchPointDatas = append(touchPointDatas, &touchPointData)
	}

	err = d.deliveryData.Put(ctx, campaign, creatives, content, touchPointDatas)
	if err != nil {
		return err
	}
	return nil
}

// 配信停止処理
func (d *delivery) Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error {
	_, err := d.UpdateStatus(ctx, tx, campaign.ID, status, time.Now())
	if err != nil {
		return err
	}
	err = d.deliveryData.Delete(ctx, strconv.Itoa(campaign.ID))
	if err != nil {
		return err
	}
	return nil
}

func (d *delivery) UpdateStatus(ctx context.Context, tx repository.Transaction, campaignId int, status string, updatedAt time.Time) (int, error) {
	condition := &repository.UpdateCondition{
		CampaignID: campaignId,
		Status:     status,
		UpdatedAt:  updatedAt,
	}
	updatedCampaignId, err := d.campaignRepository.UpdateStatus(ctx, tx, condition)
	if err != nil {
		return 0, err
	}
	return updatedCampaignId, nil
}

func (d *delivery) DeliveryControlEvent(ctx context.Context,
	campaign *models.Campaign, beforeStatus string, afterStatus string, detail string) {
	d.deliveryControlEvent.Publish(ctx, campaign.ID, campaign.OrgCode, beforeStatus, afterStatus, detail)
}
