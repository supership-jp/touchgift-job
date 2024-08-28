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

// DeliveryEnd is interface
type DeliveryEnd interface {
	StartMonitoring(ctx context.Context, wg *sync.WaitGroup)
	Close()
}

type deliveryEnd struct {
	logger             usecase.Logger
	monitor            *metrics.Monitor
	config             *config.DeliveryEnd
	appTicker          AppTicker
	worker             deliveryEndWorker
	transaction        gateways.TransactionHandler
	deliveryEndUsecase usecase.DeliveryEnd
}

type deliveryEndWorker struct {
	wg *sync.WaitGroup
	q  chan *DeliveryEndCondition
}

var (
	metricDeliveryEndCampaignTotal       = "delivery_end_campaign_total"
	metricDeliveryEndCampaignTotalDesc   = "all end campaign count"
	metricDeliveryEndCampaignTotalLabels = []string{"status"}

	metricDeliveryEndCampaignDuration        = "delivery_end_campaign_duration_seconds"
	metricDeliveryEndCampaignDurationDesc    = "delivery end reserve processing time (seconds)"
	metricDeliveryEndCampaignDurationLabels  = []string{"kind"}
	metricDeliveryEndCampaignDurationBuckets = []float64{0.01, 0.025, 0.050, 0.075, 0.100, 0.300, 0.500}
)

// NewDeliveryEnd is function
func NewDeliveryEnd(
	logger usecase.Logger,
	monitor *metrics.Monitor,
	config *config.DeliveryEnd,
	appTicker AppTicker,
	transaction gateways.TransactionHandler,
	deliveryEndUsecase usecase.DeliveryEnd,
) DeliveryEnd {
	instance := deliveryEnd{
		logger:    logger,
		monitor:   monitor,
		config:    config,
		appTicker: appTicker,
		worker: deliveryEndWorker{
			wg: &sync.WaitGroup{},
			q:  make(chan *DeliveryEndCondition, config.NumberOfQueue),
		},
		transaction:        transaction,
		deliveryEndUsecase: deliveryEndUsecase,
	}
	monitor.Metrics.AddCounter(metricDeliveryEndCampaignTotal, metricDeliveryEndCampaignTotalDesc, metricDeliveryEndCampaignTotalLabels)
	monitor.Metrics.AddHistogram(metricDeliveryEndCampaignDuration, metricDeliveryEndCampaignDurationDesc,
		metricDeliveryEndCampaignDurationLabels, metricDeliveryEndCampaignDurationBuckets)
	return &instance
}

// 配信終了の監視を始める
// 指定時間毎に開始対象があるかチェックして処理する
func (d *deliveryEnd) StartMonitoring(ctx context.Context, wg *sync.WaitGroup) {
	d.logger.Info().Msg("Start monitoring")
	d.createWorker(ctx)
	d.deliveryEndUsecase.CreateWorker(ctx)
	wg.Add(1)
	ticker := d.appTicker.New(d.config.TaskInterval, time.Minute)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			baseTime := now.Truncate(time.Minute)
			// 配信終了処理
			go d.call(ctx, &DeliveryEndCondition{
				BaseTime: baseTime,
				// 10sはぎりぎりで配信終了するのを防ぐために追加している
				To:     baseTime.Add(d.config.TaskInterval).Add(10 * time.Second),
				Status: []string{codes.StatusStarted, codes.StatusPaused},
				r:      make(chan int, d.config.NumberOfQueue),
			})
			// 再起動等でterminateになったままのものを処理する
			go d.call(ctx, &DeliveryEndCondition{
				BaseTime: baseTime,
				To:       baseTime,
				Status:   []string{codes.StatusTerminate},
				r:        make(chan int, d.config.NumberOfQueue),
			})
		case <-ctx.Done():
			d.logger.Info().Msg("Close monitoring")
			wg.Done()
			return
		}
	}
}

func (d *deliveryEnd) Close() {
	close(d.worker.q)
	d.worker.wg.Wait()
	d.deliveryEndUsecase.Close()
}

// 配信処理用のWorkerを作成
func (d *deliveryEnd) createWorker(ctx context.Context) {
	for i := 0; i < d.config.NumberOfConcurrent; i++ {
		d.worker.wg.Add(1)
		// Workerが呼び出された時に実際に動く処理
		go d.execute(ctx)
	}
}

// 配信終了タスクを呼び出す
func (d *deliveryEnd) call(ctx context.Context, condition *DeliveryEndCondition) {
	startTime := time.Now()
	defer func() {
		latency := time.Since(startTime)
		d.monitor.Metrics.GetHistogram(metricDeliveryEndCampaignDuration).WithLabelValues("execute_end").Observe(latency.Seconds())
		close(condition.r)
	}()
	for {
		select {
		case <-ctx.Done():
			d.logger.Debug().Msg("Close to call")
			return
		default:
			// 配信終了タスクを呼び出す
			d.worker.q <- condition
			// 配信終了タスクでの処理数を取得
			count := <-condition.r
			if count < d.config.TaskLimit {
				// 残りのデータがなくなったら終了
				return
			}
		}
	}
}

// 実際の配信終了処理
func (d *deliveryEnd) execute(ctx context.Context) {
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
				d.logger.Error().Err(err).Time("baseTime", condition.BaseTime).Strs("status", condition.Status).Msg("Failed to process")
			}
			wg.Done()
		case <-ctx.Done():
			wg.Wait()
			d.logger.Debug().Msg("Close execute")
			return
		}
	}
}

func (d *deliveryEnd) process(ctx context.Context, condition *DeliveryEndCondition) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic. reason: %#v", r)
		}
	}()
	// 配信終了処理をする
	// 配信終了日の検索条件 (end_at < to and status = ('started','paused'))
	// 例. baseTime: 2021/01/06 10:00 , TaskInterval: 1m の場合
	// end_at < 2021/01/06 10:01 and status in ('started','paused') のものを取得する
	baseTime := condition.BaseTime
	to := condition.To
	// 終了対象キャンペーンを取得
	campaigns, err := d.deliveryEndUsecase.GetDeliveryDataCampaigns(ctx, to, condition.Status, d.config.TaskLimit)
	if err != nil {
		condition.r <- 0
		return errors.Wrap(err, "Failed to GetDeliveryDataCampaign")
	}
	condition.r <- len(campaigns)
	for i := range campaigns {
		campaign := (campaigns)[i]
		d.monitor.Metrics.GetCounter(metricDeliveryEndCampaignTotal).WithLabelValues(campaign.Status).Inc()
		switch campaign.Status {
		case codes.StatusPaused, codes.StatusStarted:
			err := d.handlePausedOrStarted(ctx, baseTime, campaign)
			if err != nil {
				d.logger.Error().Err(err).Time("baseTime", baseTime).Int("campaign_id", campaign.ID).Msgf("Failed to handle %s", campaign.Status)
			}
		case codes.StatusTerminate:
			// 再起動等でterminateになったままのものを処理する
			// 既に終了時間を過ぎているのですぐに終了する
			d.deliveryEndUsecase.ExecuteNow(campaign)
		}
	}
	return nil
}

func (d *deliveryEnd) handlePausedOrStarted(ctx context.Context, baseTime time.Time, campaign *models.Campaign) (err error) {
	var tx gateways.Transaction
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic. reason: %#v", r)
		}
		if err != nil && tx != nil {
			if terr := tx.Rollback(); terr != nil {
				d.logger.Error().Time("baseTime", baseTime).Err(terr).Int("campaign_id", campaign.ID).Msg("Failed to rollback")
			}
		}
	}()
	tx, err = d.transaction.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to begin transaction")
	}
	// ステータスをterminateに更新
	_, err = d.deliveryEndUsecase.Terminate(ctx, tx, campaign.ID, campaign.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "Failed to update")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit")
	}
	// 取得した対象の配信終了時間を指定時間として実行する
	if campaign.EndAt.Valid {
		d.deliveryEndUsecase.Reserve(ctx, campaign.EndAt.Time, campaign)
	}
	return nil
}

// DeliveryEndCondition is struct
// 配信開始処理の条件
type DeliveryEndCondition struct {
	BaseTime time.Time
	To       time.Time
	Status   []string
	r        chan int
}
