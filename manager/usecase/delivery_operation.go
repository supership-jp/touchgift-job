//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"time"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"
)

type DeliveryOperation interface {
	// キャンペーンのログを処理する
	Process(ctx context.Context, current time.Time, campaignLog *models.CampaignLog) error
}

type deliveryOperation struct {
	logger      Logger
	monitor     *metrics.Monitor
	transaction repository.TransactionHandler
}

func NewDeliveryOperation(
	logger Logger,
	monitor *metrics.Monitor,
	transaction repository.TransactionHandler,
) DeliveryOperation {
	instance := deliveryOperation{
		logger:      logger,
		monitor:     monitor,
		transaction: transaction,
	}

	// TODO:メトリクスの追加: どれだけデータが処理されたか
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
		// TODO:キャンペーンログを処理する
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
func (d *deliveryOperation) processCampaignLog(
	ctx context.Context, campaign *models.CampaignLog) error {

	// 該当する配信データを取得

	// 配信データが取得できない場合は何もしない

	// 配信データを同期する

	return nil
}

// 配信データ同期処理
func (d *deliveryOperation) sync(
	ctx context.Context,
	tx repository.Transaction,
	// TODO: キャンペーン情報入れる
) (string, string, error) {
	// 配信ステータスにより配信データを更新する

	return "", "", nil
}
