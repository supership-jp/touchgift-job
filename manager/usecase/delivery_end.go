//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"strconv"
	"sync"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"

	"github.com/pkg/errors"
)

var (
	metricDeliveryEndDuration        = "adm_delivery_end_duration_seconds"
	metricDeliveryEndDurationDesc    = "adm delivery end processing time (seconds)"
	metricDeliveryEndDurationBuckets = []float64{0.025, 0.050, 0.100, 0.300, 0.500}
)

// DeliveryEnd is interface
type DeliveryEnd interface {
	// 終了対象キャンペーンを取得する
	GetDeliveryDataCampaigns(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error)
	// キャンペーンのステータスをTERMINATEにする
	Terminate(ctx context.Context, tx repository.Transaction, campaignID int, updatedAt time.Time) (int, error)
	// 配信終了処理を予約する
	Reserve(ctx context.Context, endAt time.Time, campaign *models.Campaign)
	// 配信開始処理を実行する(即時)
	ExecuteNow(campaign *models.Campaign)
	// 配信停止処理
	Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error
	// 配信データ削除
	Delete(ctx context.Context, campaign *models.Campaign) error
	// 終了する
	Close()
	// Workerを作成する
	CreateWorker(ctx context.Context)
}

type deliveryEnd struct {
	logger                   Logger
	monitor                  *metrics.Monitor
	config                   *config.DeliveryEnd
	configUsecase            *config.DeliveryEndUsecase
	worker                   deliveryEndWorker
	transaction              repository.TransactionHandler
	timer                    Timer
	deliveryControlEvent     DeliveryControlEvent
	campaignRepository       repository.CampaignRepository
	campaignDataRepository   repository.DeliveryDataCampaignRepository
	contentDataRepository    repository.DeliveryDataContentRepository
	touchPointDataRepository repository.DeliveryDataTouchPointRepository
	touchPointRepository     repository.TouchPointRepository
}

type deliveryEndWorker struct {
	wg *sync.WaitGroup
	q  chan *models.Campaign
}

// NewDeliveryEnd is function
func NewDeliveryEnd(
	logger Logger,
	monitor *metrics.Monitor,
	config *config.DeliveryEnd,
	configUsecase *config.DeliveryEndUsecase,
	transaction repository.TransactionHandler,
	timer Timer,
	deliveryControlEvent DeliveryControlEvent,
	campaignRepository repository.CampaignRepository,
	campaignDataRepository repository.DeliveryDataCampaignRepository,
	contentDataRepository repository.DeliveryDataContentRepository,
	touchPointDataRepository repository.DeliveryDataTouchPointRepository,
	touchPointRepository repository.TouchPointRepository,
) DeliveryEnd {
	instance := deliveryEnd{
		logger:        logger,
		monitor:       monitor,
		config:        config,
		configUsecase: configUsecase,
		worker: deliveryEndWorker{
			wg: &sync.WaitGroup{},
			q:  make(chan *models.Campaign, config.NumberOfQueue),
		},
		transaction:              transaction,
		timer:                    timer,
		deliveryControlEvent:     deliveryControlEvent,
		campaignRepository:       campaignRepository,
		campaignDataRepository:   campaignDataRepository,
		contentDataRepository:    contentDataRepository,
		touchPointDataRepository: touchPointDataRepository,
		touchPointRepository:     touchPointRepository,
	}
	monitor.Metrics.AddHistogram(metricDeliveryEndDuration,
		metricDeliveryEndDurationDesc,
		nil,
		metricDeliveryEndDurationBuckets)
	return &instance
}

// 配信終了処理用のWorkerを作成
func (d *deliveryEnd) CreateWorker(ctx context.Context) {
	for i := 0; i < d.configUsecase.NumberOfConcurrent; i++ {
		d.worker.wg.Add(1)
		// Workerが呼び出された時に実際に動く処理
		go d.execute(ctx)
	}
}

func (d *deliveryEnd) Close() {
	close(d.worker.q)
	d.worker.wg.Wait()
}

// 配信終了処理を指定時間に実行するように予約する
func (d *deliveryEnd) Reserve(ctx context.Context, endAt time.Time, campaign *models.Campaign) {
	d.timer.ExecuteAtTime(ctx, endAt, func() {
		d.ExecuteNow(campaign)
	})
}

// 配信終了処理を実行する(即時)
func (d *deliveryEnd) ExecuteNow(campaign *models.Campaign) {
	d.worker.q <- campaign // 実行する
}

// 終了対象キャンペーンを取得する
func (d *deliveryEnd) GetDeliveryDataCampaigns(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
	condition := repository.CampaignDataToEndCondition{
		End:    to,
		Status: status,
	}
	return d.campaignRepository.GetCampaignToEnd(ctx, &condition)
}

func (d *deliveryEnd) Terminate(ctx context.Context, tx repository.Transaction, campaignID int, updatedAt time.Time) (int, error) {
	condition := &repository.UpdateCondition{
		CampaignID: campaignID,
		Status:     codes.StatusTerminate,
		UpdatedAt:  updatedAt,
	}
	updatedCampaignID, err := d.campaignRepository.UpdateStatus(ctx, tx, condition)
	if err != nil {
		return 0, err
	}
	return updatedCampaignID, nil
}

// 配信停止処理 RDBのキャンペーンステータスを更新する
func (d *deliveryEnd) Stop(ctx context.Context, tx repository.Transaction, campaign *models.Campaign, status string) error {
	condition := &repository.UpdateCondition{
		CampaignID: campaign.ID,
		Status:     status,
		UpdatedAt:  time.Now(),
	}
	_, err := d.campaignRepository.UpdateStatus(ctx, tx, condition)
	if err != nil {
		return err
	}
	return nil
}

// DynamoDBから配信データを削除する
func (d *deliveryEnd) Delete(ctx context.Context, campaign *models.Campaign) error {
	campaignID := strconv.Itoa(campaign.ID)
	if err := d.campaignDataRepository.Delete(ctx, &campaignID); err != nil {
		return err
	}
	if err := d.contentDataRepository.Delete(ctx, &campaignID); err != nil {
		return err
	}
	// グループに紐づく配信中のキャンペーンがない場合はタッチポイントデータも削除する
	count, err := d.campaignRepository.GetDeliveryCampaignCountByGroupID(ctx, campaign.GroupID)
	if err != nil {
		return err
	}
	if count == 0 {
		touchPoints, err := d.touchPointRepository.GetTouchPointByGroupID(ctx, &repository.TouchPointByGroupIDCondition{
			GroupID: campaign.GroupID,
			Limit:   100000,
		})
		if err != nil {
			return err
		}
		for _, touchPoint := range touchPoints {
			groupIDStr := strconv.Itoa(touchPoint.GroupID)
			if err := d.touchPointDataRepository.Delete(ctx, &touchPoint.ID, &groupIDStr); err != nil {
				return err
			}
			d.deliveryControlEvent.PublishDeliveryEvent(ctx, touchPoint.ID, touchPoint.GroupID, campaign.ID, campaign.OrgCode, "DELETE")
		}
	}
	return nil
}

// 配信終了処理
func (d *deliveryEnd) execute(ctx context.Context) {
	defer d.worker.wg.Done()
	wg := sync.WaitGroup{}
	for {
		select {
		case reservedData, ok := <-d.worker.q:
			if !ok {
				return
			}
			wg.Add(1)
			startTime := time.Now()
			err := d.end(ctx, startTime, reservedData)
			if err != nil {
				d.logger.Error().Err(err).Time("baseTime", startTime).Int("id", reservedData.ID).Msg("Failed to end")
			} else {
				latency := time.Since(startTime)
				d.monitor.Metrics.
					GetHistogram(metricDeliveryEndDuration).
					WithLabelValues().Observe(latency.Seconds())
			}
			wg.Done()
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}

//nolint:gocognit // [21]時間あるときに修正する
func (d *deliveryEnd) end(ctx context.Context, startTime time.Time, reservedData *models.Campaign) (err error) {
	var tx repository.Transaction
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic. reason: %#v", r)
		}
		if err != nil && tx != nil {
			if terr := tx.Rollback(); terr != nil {
				d.logger.Error().Err(terr).Time("baseTime", startTime).Int("campaign_id", reservedData.ID).Msg("Failed to rollback")
			}
		}
	}()
	tx, err = d.transaction.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to start transaction")
	}
	d.logger.Debug().Int("campaign_id", reservedData.ID).Msg("Get campaign")

	condition := repository.CampaignCondition{
		CampaignID: reservedData.ID,
	}
	deliveryData, err := d.campaignRepository.GetDeliveryToStart(ctx, tx, &condition)
	if err != nil {
		return errors.Wrap(err, "Failed to get deliveryData")
	}

	// terminate以外はエラーを返す
	if deliveryData.Status != codes.StatusTerminate {
		return errors.Errorf("campaign other than terminate. status: %s", deliveryData.Status)
	}
	// 終了処理
	afterStatus, err := d.handleTerminateStatus(ctx, tx, deliveryData)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit")
	}
	// 配信制御イベントを発行する
	d.deliveryControlEvent.PublishCampaignEvent(ctx, deliveryData.ID, deliveryData.GroupID, deliveryData.OrgCode, deliveryData.Status, *afterStatus, "")
	return nil
}

func (d *deliveryEnd) handleTerminateStatus(
	ctx context.Context, tx repository.Transaction, deliveryData *models.Campaign) (status *string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic. reason: %#v", r)
		}
	}()
	// 終了処理
	afterStatus := codes.StatusEnded
	if err := d.Stop(ctx, tx, deliveryData, afterStatus); err != nil {
		return nil, errors.Wrap(err, "Failed to delete process")
	}
	if err := d.Delete(ctx, deliveryData); err != nil {
		return nil, errors.Wrap(err, "Failed to delete process")
	}
	return &afterStatus, nil
}
