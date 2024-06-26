package usecase

import (
	"context"
	"time"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"
)

type Delivery interface {
	// UpdateStatus 配信状態を更新する
	UpdateStatus(ctx context.Context, tx repository.Transaction, campaignId int, status string, updatedAt time.Time) (int64, time.Time, error)
	// DeliveryControlEvent 配信制御イベントを発行する
	DeliveryControlEvent(ctx context.Context, campaign *models.Campaign, beforeStatus string, afterStatus string, detail string, deliveryType string, priceType string)
	// StartOrSync データの同期or配信を開始する
	StartOrSync(ctx context.Context, tx repository.Transaction, campaign *models.Campaign) error
	// Stop 配信を停止する
	Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error
}

type delivery struct {
	logger             Logger
	monitor            *metrics.Monitor
	contentsRepository repository.ContentsRepository
}

func NewDelivery(
	logger Logger,
	monitor *metrics.Monitor,
) Delivery {
	return &delivery{
		logger:  logger,
		monitor: monitor,
	}
}

func (d *delivery) StartOrSync(ctx context.Context, tx repository.Transaction, campaign *models.Campaign) error {
	//campaignID := campaign.ID
	// TODO: campaign.statusの更新
	//_, _, err = d.UpdateStatus(ctx, tx, campaignID, "started", time.Now())
	//if err != nil {
	//	return err
	//}

	// TODO: IDの型の取り扱いを考える
	//condition := repository.ContentsByCampaignIDCondition{
	//	CampaignID: campaignID,
	//}

	// TODO: コンテンツをそれぞれキャンペーンから取得してメモリに展開

	// ギミックURLの取得
	//gimmickURL, err := d.contentsRepository.GetGimmickURLByCampaignID(ctx, tx, &condition)
	//if err != nil {
	//	return err
	//}
	//　クーポン一覧の取得
	//coupons, err := d.contentsRepository.GetCouponsByCampaignID(ctx, tx, &condition)

	//　TODO: 配信データ(キャンペーン, タッチポイント, コンテンツ)登録

	return nil
}

func (d *delivery) Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error {
	return nil
}

func (d *delivery) UpdateStatus(ctx context.Context, tx repository.Transaction, campaignId int, status string, updatedAt time.Time) (int64, time.Time, error) {
	//condition := repository.UpdateCondition{
	//	CampaignID: campaignId,
	//	Status:     status,
	//	UpdatedAt:  updatedAt,
	//}
	//count, executeUpdatedAt, err :=
	//if err != nil {
	//	return 0, time.Time{}, err
	//}
	return 0, time.Now(), nil
}

func (d *delivery) DeliveryControlEvent(ctx context.Context, campaign *models.Campaign, beforeStatus string, afterStatus string, detail string, deliveryType string, priceType string) {
	// メソッドの実装
}
