package controllers

import (
	"context"
	"sync"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra/metrics"
	"touchgift-job-manager/interface/gateways"
	"touchgift-job-manager/usecase"

	"github.com/pkg/errors"
)

type DeliveryStart interface {
	StartMonitoring(ctx context.Context, wg *sync.WaitGroup)
	Close()
}

type deliveryStart struct {
	logger               usecase.Logger
	monitor              *metrics.Monitor
	config               *config.DeliveryStart
	appTicker            AppTicker
	worker               deliveryStartWorker
	transaction          gateways.TransactionHandler
	deliveryStartUsecase usecase.DeliveryStart
	deliveryControlEvent usecase.DeliveryControlEvent
}

type deliveryStartWorker struct {
	wg *sync.WaitGroup
	q  chan *DeliveryStartCondition
}

var (
	metricDeliveryStartCampaignTotal       = "delivery_start_campaign_total"
	metricDeliveryStartCampaignTotalDesc   = "all start campaign count"
	metricDeliveryStartCampaignTotalLabels = []string{"status"}

	metricDeliveryStartCampaignDuration        = "delivery_start_campaign_duration_seconds"
	metricDeliveryStartCampaignDurationDesc    = "delivery start reserve processing time (seconds)"
	metricDeliveryStartCampaignDurationLabels  = []string{"kind"}
	metricDeliveryStartCampaignDurationBuckets = []float64{0.01, 0.025, 0.050, 0.075, 0.100, 0.300, 0.500}
)

func NewDeliveryStart(
	logger usecase.Logger,
	monitor *metrics.Monitor,
	config *config.DeliveryStart,
	appTicker AppTicker,
	transaction gateways.TransactionHandler,
	deliveryStartUsecase usecase.DeliveryStart,
	deliveryControlEvent usecase.DeliveryControlEvent,
) DeliveryStart {
	monitor.Metrics.AddCounter(metricDeliveryStartCampaignTotal, metricDeliveryStartCampaignTotalDesc, metricDeliveryStartCampaignTotalLabels)
	monitor.Metrics.AddHistogram(metricDeliveryStartCampaignDuration, metricDeliveryStartCampaignDurationDesc, metricDeliveryStartCampaignDurationLabels, metricDeliveryStartCampaignDurationBuckets)
	return &deliveryStart{
		logger:    logger,
		monitor:   monitor,
		config:    config,
		appTicker: appTicker,
		worker: deliveryStartWorker{
			wg: &sync.WaitGroup{},
			q:  make(chan *DeliveryStartCondition, config.NumberOfQueue),
		},
		transaction:          transaction,
		deliveryStartUsecase: deliveryStartUsecase,
		deliveryControlEvent: deliveryControlEvent,
	}
}

func (d *deliveryStart) StartMonitoring(ctx context.Context, wg *sync.WaitGroup) {
	d.logger.Info().Msg("Start monitoring delivery start")
	d.createWorker(ctx)
	d.deliveryStartUsecase.CreateWorker(ctx)
	wg.Add(1)
	ticker := d.appTicker.New(d.config.TaskInterval, time.Minute)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			baseTime := now.Truncate(time.Minute)
			// 配信開始処理
			go d.call(ctx, &DeliveryStartCondition{
				BaseTime: baseTime,
				// 10sはぎりぎりで配信開始するのを防ぐために追加している
				To:     baseTime.Add(d.config.TaskInterval).Add(10 * time.Second),
				Status: codes.StatusConfigured,
				r:      make(chan int, d.config.NumberOfQueue),
			})

			// 再起動等でwarmupになったままのものを処理する
			go d.call(ctx, &DeliveryStartCondition{
				BaseTime: baseTime,
				To:       baseTime,
				Status:   codes.StatusWarmup,
				r:        make(chan int, d.config.NumberOfQueue),
			})
		case <-ctx.Done():
			d.logger.Info().Msg("Close monitoring delivery start")
			wg.Done()
			return
		}
	}
}

func (d *deliveryStart) Close() {
	close(d.worker.q)
	d.worker.wg.Wait()
	d.deliveryStartUsecase.Close()
}

// 配信処理用のWorkerを作成
func (d *deliveryStart) createWorker(ctx context.Context) {
	for i := 0; i < d.config.NumberOfConcurrent; i++ {
		d.worker.wg.Add(1)
		// Workerが呼び出された時に実際に動く処理
		go d.execute(ctx)
	}
}

func (d *deliveryStart) call(ctx context.Context, condition *DeliveryStartCondition) {
	startTime := time.Now()
	defer func() {
		latency := time.Since(startTime)
		d.monitor.Metrics.GetHistogram(metricDeliveryStartCampaignDuration).WithLabelValues("close_start").Observe(latency.Seconds())
		close(condition.r)
	}()
	for {
		select {
		case <-ctx.Done():
			d.logger.Debug().Msg("Close call")
			return
		default:
			d.worker.q <- condition
			count := <-condition.r
			if count < d.config.TaskLimit {
				return
			}
		}
	}
}

// 実際の配信開始処理
func (d *deliveryStart) execute(ctx context.Context) {
	defer d.worker.wg.Done()
	wg := sync.WaitGroup{}
	for {
		select {
		case condition, ok := <-d.worker.q:
			if !ok {
				return
			}
			wg.Add(1)
			err := d.process(ctx, condition)
			if err != nil {
				d.logger.Error().Err(err).Time("baseTime", condition.BaseTime).Str("status", condition.Status).Msg("Failed to process")
			}
			wg.Done()
		case <-ctx.Done():
			wg.Wait()
			d.logger.Debug().Msg("Close execute")
			return
		}
	}
}

func (d *deliveryStart) process(ctx context.Context, condition *DeliveryStartCondition) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic. reason: %#v", r)
		}
	}()
	// 配信開始処理をする
	// 配信開始日の検索条件 (start_at < to and status = 'configured')
	// 例. baseTime: 2021/01/06 10:00 , TaskInterval: 1m の場合
	// start_at < 2021/01/06 10:01 and status = 'configured' のものを取得する
	baseTime := condition.BaseTime
	to := condition.To
	d.logger.Debug().Time("baseTime", baseTime).Str("status", condition.Status).Time("to", to).Msg("Condtion parameter")

	// 開始対象キャンペーンを取得
	campaigns, err := d.deliveryStartUsecase.GetCampaignToStart(ctx, to, condition.Status, d.config.TaskLimit)
	if err != nil {
		condition.r <- 0
		return errors.Wrap(err, "Failed to GetCampaignData")
	}

	condition.r <- len(campaigns)
	for i := range campaigns {
		campaign := (campaigns)[i]
		d.monitor.Metrics.GetCounter(metricDeliveryStartCampaignTotal).WithLabelValues(campaign.Status).Inc()
		switch campaign.Status {
		case codes.StatusConfigured:
			err := d.handleConfigured(ctx, baseTime, campaign)
			if err != nil {
				d.logger.Error().Err(err).Time("baseTime", baseTime).
					Int("campaign_id", campaign.ID).
					Msg("Failed to handle configured")
			}
		case codes.StatusWarmup:
			// 再起動等でwarmupになったままのものを処理する
			// 既に開始時間を過ぎているのですぐに開始する
			d.deliveryStartUsecase.ExecuteNow(campaign)
		}
	}
	return nil
}

func (d *deliveryStart) handleConfigured(ctx context.Context, baseTime time.Time, campaign *models.Campaign) (err error) {
	var tx gateways.Transaction
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic. reason: %#v", r)
		}
		if err != nil && tx != nil {
			if terr := tx.Rollback(); terr != nil {
				d.logger.Error().Time("baseTime", baseTime).Err(terr).
					Int("campaign_id", campaign.ID).Msg("Failed to rollback")
			}
		}
	}()
	tx, err = d.transaction.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to begin transaction")
	}
	// ステータスをwarmupに更新
	_, err = d.deliveryStartUsecase.UpdateStatus(ctx, tx, campaign, codes.StatusWarmup)
	if err != nil {
		return errors.Wrap(err, "Failed to update")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit")
	}
	// 配信制御イベントを発行する
	d.deliveryControlEvent.PublishCampaignEvent(ctx, campaign.ID, campaign.GroupID, campaign.OrgCode, codes.StatusConfigured, codes.StatusWarmup, "")
	// 取得した配信対象の開始時間を指定時間として実行する
	d.deliveryStartUsecase.Reserve(ctx, campaign.StartAt, campaign)
	return nil
}

// DeliveryStartCondition is struct
// 配信開始処理の条件
type DeliveryStartCondition struct {
	BaseTime time.Time
	To       time.Time
	Status   string
	r        chan int
}
