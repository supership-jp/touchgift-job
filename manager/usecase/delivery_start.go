//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"fmt"
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
	metricDeliveryStartDuration        = "touchgift_delivery_start_duration_seconds"
	metricDeliveryStartDurationDesc    = "touchgift delivery start processing time (seconds)"
	metricDeliveryStartDurationBuckets = []float64{0.025, 0.050, 0.100, 0.300, 0.500}
)

// DeliveryStart is interface
type DeliveryStart interface {
	// 開始対象キャンペーンを取得する
	GetCampaignToStart(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error)
	// キャンペーンのステータスを変更する
	UpdateStatus(ctx context.Context, tx repository.Transaction, Campaign *models.Campaign, status string) (int, error)
	// 配信開始処理を予約する
	Reserve(ctx context.Context, startAt time.Time, Campaign *models.Campaign)
	// 配信開始処理を実行する(即時)
	ExecuteNow(schduleData *models.Campaign)
	// 終了する
	Close()
	// Workerを作成する
	CreateWorker(ctx context.Context)
	// 配信データを作成する
	CreateDeliveryDatas(ctx context.Context, tx repository.Transaction, campaign *models.Campaign) error
}

type deliveryStart struct {
	logger                   Logger
	monitor                  *metrics.Monitor
	config                   *config.DeliveryStart
	configUsecase            *config.DeliveryStartUsecase
	worker                   deliveryStartWorker
	transaction              repository.TransactionHandler
	timer                    Timer
	deliveryControlEvent     DeliveryControlEvent
	campaignRepository       repository.CampaignRepository
	creativeRepository       repository.CreativeRepository
	contentRepository        repository.ContentRepository
	touchPointRepository     repository.TouchPointRepository
	campaignDataRepository   repository.DeliveryDataCampaignRepository
	contentDataRepository    repository.DeliveryDataContentRepository
	creativeDataRepository   repository.DeliveryDataCreativeRepository
	touchPointDataRepository repository.DeliveryDataTouchPointRepository
}

type deliveryStartWorker struct {
	wg *sync.WaitGroup
	q  chan *models.Campaign
}

// NewDeliveryStart is function
func NewDeliveryStart(
	logger Logger,
	monitor *metrics.Monitor,
	config *config.DeliveryStart,
	configUsecase *config.DeliveryStartUsecase,
	transaction repository.TransactionHandler,
	timer Timer,
	deliveryControlEvent DeliveryControlEvent,
	campaignRepository repository.CampaignRepository,
	creativeRepository repository.CreativeRepository,
	contentRepository repository.ContentRepository,
	touchPointRepository repository.TouchPointRepository,
	campaignDataRepository repository.DeliveryDataCampaignRepository,
	contentDataRepository repository.DeliveryDataContentRepository,
	creativeDataRepository repository.DeliveryDataCreativeRepository,
	touchPointDataRepository repository.DeliveryDataTouchPointRepository,
) DeliveryStart {
	instance := deliveryStart{
		logger:        logger,
		monitor:       monitor,
		config:        config,
		configUsecase: configUsecase,
		worker: deliveryStartWorker{
			wg: &sync.WaitGroup{},
			q:  make(chan *models.Campaign, config.NumberOfQueue),
		},
		transaction:              transaction,
		timer:                    timer,
		deliveryControlEvent:     deliveryControlEvent,
		campaignRepository:       campaignRepository,
		creativeRepository:       creativeRepository,
		contentRepository:        contentRepository,
		touchPointRepository:     touchPointRepository,
		campaignDataRepository:   campaignDataRepository,
		contentDataRepository:    contentDataRepository,
		creativeDataRepository:   creativeDataRepository,
		touchPointDataRepository: touchPointDataRepository,
	}
	monitor.Metrics.AddHistogram(metricDeliveryStartDuration, metricDeliveryStartDurationDesc, nil, metricDeliveryStartDurationBuckets)
	return &instance
}

// 配信開始処理用のWorkerを作成
func (d *deliveryStart) CreateWorker(ctx context.Context) {
	for i := 0; i < d.configUsecase.NumberOfConcurrent; i++ {
		d.worker.wg.Add(1)
		// Workerが呼び出された時に実際に動く処理
		go d.execute(ctx)
	}
}

func (d *deliveryStart) Close() {
	close(d.worker.q)
	d.worker.wg.Wait()
}

// 配信開始処理を指定時間に実行するように予約する
func (d *deliveryStart) Reserve(ctx context.Context, startAt time.Time, Campaign *models.Campaign) {
	// サーバキャッシュ用に150ms早く動かす
	d.timer.ExecuteAtTime(ctx, startAt.Add(-150*time.Millisecond), func() {
		d.ExecuteNow(Campaign)
	})
}

// 配信開始処理を実行する(即時)
func (d *deliveryStart) ExecuteNow(campaign *models.Campaign) {
	d.worker.q <- campaign // 実行する
}

// 開始対象キャンペーンを取得する
func (d *deliveryStart) GetCampaignToStart(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
	d.logger.Debug().Msg("DeliveryStart GetCampaignToStart")
	condition := repository.CampaignToStartCondition{
		To:     to,
		Status: status,
	}
	tx, err := d.transaction.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to start transaction")
	}
	return d.campaignRepository.GetCampaignToStart(ctx, tx, &condition)
}

func (d *deliveryStart) UpdateStatus(ctx context.Context, tx repository.Transaction, Campaign *models.Campaign, status string) (int, error) {
	condition := repository.UpdateCondition{
		CampaignID: Campaign.ID,
		Status:     status,
		UpdatedAt:  Campaign.UpdatedAt,
	}
	count, err := d.campaignRepository.UpdateStatus(ctx, tx, &condition)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Failed to update status. status: %s", status))
	}
	return count, nil
}

// 配信開始処理
func (d *deliveryStart) execute(ctx context.Context) {
	defer func() {
		d.logger.Debug().Msg("End execute")
		d.worker.wg.Done()
	}()
	wg := sync.WaitGroup{}
	for {
		select {
		case reservedData, ok := <-d.worker.q:
			if !ok {
				return
			}
			wg.Add(1)
			startTime := time.Now()
			err := d.start(ctx, startTime, reservedData)
			if err != nil {
				d.logger.Error().Err(err).Time("baseTime", startTime).Int("id", reservedData.ID).Msg("Failed to start")
			} else {
				latency := time.Since(startTime)
				d.monitor.Metrics.
					GetHistogram(metricDeliveryStartDuration).
					WithLabelValues().Observe(latency.Seconds())
			}
			wg.Done()
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}

//nolint:gocognit // [23]時間あるときに修正する
func (d *deliveryStart) start(
	ctx context.Context, startTime time.Time, reservedData *models.Campaign,
) (err error) {
	var tx repository.Transaction
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic. reason: %#v", r)
		}
		if err != nil && tx != nil {
			if terr := tx.Rollback(); terr != nil {
				d.logger.Error().Err(terr).Time("baseTime", startTime).Int("id", reservedData.ID).
					Msg("Failed to rollback")
			}
		}
	}()
	tx, err = d.transaction.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to start transaction")
	}

	condition := repository.CampaignCondition{
		CampaignID: reservedData.ID,
		Status:     codes.StatusWarmup,
	}
	d.logger.Debug().Int("id", reservedData.ID).Msg("Get Campaign")
	startCampaign, err := d.campaignRepository.GetDeliveryToStart(ctx, tx, &condition)
	if err != nil {
		return errors.Wrap(err, "Failed to get startcampaign")
	}

	// warmup以外のキャンペーンは処理しないためエラーを返す
	if startCampaign.Status != codes.StatusWarmup {
		return errors.Errorf("Campaign other than warmup. status: %s", startCampaign.Status)
	}
	// 配信データ作成処理
	_, err = d.UpdateStatus(ctx, tx, startCampaign, codes.StatusStarted)
	if err != nil {
		return err
	}
	err = d.CreateDeliveryDatas(ctx, tx, startCampaign)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit")
	}
	// 配信制御イベントを発行する
	d.deliveryControlEvent.Publish(
		ctx, startCampaign.ID, startCampaign.OrgCode, startCampaign.Status, codes.StatusStarted, "")
	return nil
}

func (d *deliveryStart) CreateDeliveryDatas(ctx context.Context, tx repository.Transaction, campaign *models.Campaign) error {
	creatives, content, touchPointDatas, err := d.getDataFromRDB(ctx, tx, campaign)
	if err != nil {
		return err
	}
	err = d.createDeliveryDatas(ctx, campaign, creatives, content, touchPointDatas)
	if err != nil {
		return err
	}
	return nil
}

// 配信開始時にRDBからデータを取得する処理
func (d *deliveryStart) getDataFromRDB(ctx context.Context, tx repository.Transaction, campaign *models.Campaign) (
	[]*models.Creative, *models.DeliveryDataContent, []*models.DeliveryTouchPoint, error,
) {
	// TODO: IDの型の取り扱いを考える
	condition := repository.ContentByCampaignIDCondition{
		CampaignID: campaign.ID,
	}
	// クリエイティブの取得
	creativeCondition := repository.CreativeByCampaignIDCondition{
		CampaignID: campaign.ID,
	}
	creatives, err := d.creativeRepository.GetCreativeByCampaignID(ctx, tx, &creativeCondition)
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: コンテンツをそれぞれキャンペーンから取得してメモリに展開
	// ギミックURLの取得
	gimmickURL, gimmickCode, err := d.contentRepository.GetGimmicksByCampaignID(ctx, tx, &condition)
	if err != nil {
		return nil, nil, nil, err
	}
	//　クーポン一覧の取得
	coupons, err := d.contentRepository.GetCouponsByCampaignID(ctx, tx, &condition)
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
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
	return creatives, content, touchPointDatas, nil
}

func (d *deliveryStart) createDeliveryDatas(ctx context.Context,
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
