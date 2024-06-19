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
	logger  Logger
	monitor *metrics.Monitor
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
	// TODO: campaign.statusの更新
	_, _, err := d.UpdateStatus(ctx, tx, campaign.ID, "started", time.Now())
	if err != nil {
		return err
	}

	// TODO: コンテンツをそれぞれキャンペーンから取得してメモリに展開

	//　TODO: 配信データ(キャンペーン, タッチポイント, コンテンツ)登録

	return nil
}

func (d *delivery) Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error {
	return nil
}

func (d *delivery) UpdateStatus(ctx context.Context, tx repository.Transaction, campaignId int, status string, updatedAt time.Time) (int64, time.Time, error) {
	return 0, time.Now(), nil
}

func (d *delivery) DeliveryControlEvent(ctx context.Context, campaign *models.Campaign, beforeStatus string, afterStatus string, detail string, deliveryType string, priceType string) {
	// メソッドの実装
}
