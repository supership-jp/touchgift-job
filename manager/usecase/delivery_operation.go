//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"
)

type DeliveryOperation interface {
	// キャンペーンのログを処理する
	Process(ctx context.Context, current time.Time, campaignLog *models.CampaignLog) error
}

type deliveryOperation struct {
	logger                 Logger
	monitor                *metrics.Monitor
	transaction            repository.TransactionHandler
	campaignRepository     repository.CampaignRepository
	campaignDataRepository repository.DeliveryDataCampaignRepository
	creative               Creative
	deliveryStart          DeliveryStart
	deliveryEnd            DeliveryEnd
	deliveryControlEvent   DeliveryControlEvent
}

func NewDeliveryOperation(
	logger Logger,
	monitor *metrics.Monitor,
	transaction repository.TransactionHandler,
	campaignRepository repository.CampaignRepository,
	campaignDataRepository repository.DeliveryDataCampaignRepository,
	creative Creative,
	deliveryStart DeliveryStart,
	deliveryEnd DeliveryEnd,
	deliveryControlEvent DeliveryControlEvent,
) DeliveryOperation {
	instance := deliveryOperation{
		logger:                 logger,
		monitor:                monitor,
		transaction:            transaction,
		campaignRepository:     campaignRepository,
		campaignDataRepository: campaignDataRepository,
		creative:               creative,
		deliveryStart:          deliveryStart,
		deliveryEnd:            deliveryEnd,
		deliveryControlEvent:   deliveryControlEvent,
	}
	// TODO:メトリクスの追加: どれだけデータが処理されたか
	// monitor.Metrics.AddCounter(metricDynamodbPutTotal, metricDynamodbPutTotalDesc, metricDynamodbPutTotalLabels)
	return &instance
}
func (d *deliveryOperation) Process(ctx context.Context, current time.Time, campaignLog *models.CampaignLog) (err error) {
	var tx repository.Transaction
	defer func() {
		if err != nil && tx != nil {
			if result := tx.Rollback(); result != nil {
				d.logger.Error().Err(result).Time("current", current).Msg("Failed to rollback")
			}
		}
	}()
	tx, err = d.transaction.Begin(ctx)
	if err != nil {
		return err
	}
	switch campaignLog.Event {
	case "insert", "update":
		campaign, beforeStatus, afterStatus, err := d.processCampaignLog(ctx, tx, current, campaignLog)
		if err != nil {
			return err
		}
		switch *afterStatus {
		case codes.StatusWarmup:
			if err := tx.Commit(); err != nil {
				d.logger.Error().Err(err).Time("current", current).Msg("Failed to commit")
				return err
			}
		default:
			if err := d.creative.Process(ctx, tx, current, &campaignLog.Creatives); err != nil {
				return err
			}
			if err := tx.Commit(); err != nil {
				d.logger.Error().Err(err).Time("current", current).Msg("Failed to commit")
				return err
			}
			// 配信制御イベントを発行する
			d.deliveryControlEvent.Publish(ctx, campaign.ID, campaign.OrgCode, *beforeStatus, *afterStatus, "")
		}
	case "delete":
		// キャンペーンの物理削除は配信後には起きないためdelivery_data削除はしない
		d.logger.Info().
			Time("current", current).
			Int("campaign_id", campaignLog.ID).
			Str("event", campaignLog.Event).
			Msg("Not delete delivery_data after delivery")
		// ログのCreative情報で必要な処理をする
		if err := d.creative.Process(ctx, tx, current, &campaignLog.Creatives); err != nil {
			return err
		}
		return codes.ErrDoNothing // rollbackさせるため
	default:
		if err := tx.Commit(); err != nil {
			d.logger.Error().Err(err).Time("current", current).Msg("Failed to commit")
			return err
		}
		// TODO: 配信制御イベントを発行する
		return nil

	}
	return nil
}

// キャンペーンログの処理
func (d *deliveryOperation) processCampaignLog(ctx context.Context,
	tx repository.Transaction, current time.Time, campaign *models.CampaignLog,
) (*models.Campaign, *string, *string, error) {
	// 該当する配信データを取得
	condition := repository.CampaignCondition{
		CampaignID: campaign.ID,
	}
	campaignData, err := d.campaignRepository.GetDeliveryToStart(ctx, tx, &condition)
	if err != nil && nil != codes.ErrNoData {
		return nil, nil, nil, err
	}
	if campaignData == nil {
		// 配信データが取得できない場合は何もしない
		d.logger.Error().
			Int("campaign_id", campaign.ID).
			Msg("No delivery data")
		return nil, nil, nil, codes.ErrDoNothing
	}
	// 配信データを同期する
	before, after, err := d.sync(ctx, tx, campaignData)
	return campaignData, &before, &after, err
}

// 配信データ同期処理
func (d *deliveryOperation) sync(
	ctx context.Context,
	tx repository.Transaction,
	campaign *models.Campaign,
) (string, string, error) {
	// 配信ステータスにより配信データを更新する
	switch campaign.Status {
	// 配信中 または 配信再開
	case codes.StatusStarted, codes.StatusResume:
		// 配信データを登録(更新)する
		_, err := d.deliveryStart.UpdateStatus(ctx, tx, campaign, codes.StatusStarted)
		if err != nil {
			return campaign.Status, "", err
		}
		err = d.deliveryStart.CreateDeliveryDatas(ctx, tx, campaign)
		if err != nil {
			return campaign.Status, "", err
		}
		return campaign.Status, codes.StatusStarted, nil
		// 配信一時停止
	case codes.StatusPause:
		err := d.deliveryEnd.Stop(ctx, tx, campaign, codes.StatusPaused)
		if err != nil {
			return campaign.Status, "", err
		}
		return campaign.Status, codes.StatusPaused, d.deliveryEnd.Delete(ctx, campaign)
	// 配信停止
	case codes.StatusStop:
		err := d.deliveryEnd.Stop(ctx, tx, campaign, codes.StatusStopped)
		if err != nil {
			return campaign.Status, "", err
		}
		return campaign.Status, codes.StatusStopped, d.deliveryEnd.Delete(ctx, campaign)
		// 配信終了済
	case codes.StatusEnded:
		// DynamoDBから削除(campaign.statusの更新はしない)
		return campaign.Status, codes.StatusEnded, d.deliveryEnd.Delete(ctx, campaign)
	case codes.StatusSuspend, codes.StatusConfigured:
		// 未配信のため何もしない
		return campaign.Status, "", codes.ErrDoNothing
	default:
		d.logger.Error().Interface("delivery_data", *campaign).Msg("Unknown campaign status")
		return campaign.Status, "", codes.ErrDoNothing
	}
}
