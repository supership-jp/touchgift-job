package controllers

// TODO: 予算消化追加時に実装する

// import (
// 	"context"
// 	"encoding/json"
// 	"strings"
// 	"sync"
// 	"time"
// 	"touchgift-job-manager/config"
// 	"touchgift-job-manager/domain/models"
// 	"touchgift-job-manager/infra"
// 	"touchgift-job-manager/infra/metrics"
// 	"touchgift-job-manager/interface/gateways"
// 	"touchgift-job-manager/usecase"
// )

// // DeliveryControl is interface
// type DeliveryControlSync interface {
// 	Start(ctx context.Context, wg *sync.WaitGroup)
// 	Close()
// }

// type deliveryControlSync struct {
// 	logger                 usecase.Logger
// 	monitor                *metrics.Monitor
// 	queueHandler           gateways.QueueHandler
// 	deliveryControlUsecase usecase.DeliveryControl
// 	wg                     *sync.WaitGroup
// }

// var (
// 	metricDeliveryControlSyncTotal       = "delivery_control_sync_total"
// 	metricDeliveryControlSyncTotalDesc   = "all delivery control sync count"
// 	metricDeliveryControlSyncTotalLabels = []string{"event"}

// 	metricDeliveryControlSyncDuration        = "delivery_control_sync_duration_seconds"
// 	metricDeliveryControlSyncDurationDesc    = "delivery control processing time (seconds)"
// 	metricDeliveryControlSyncDurationLabels  = []string{"kind"}
// 	metricDeliveryControlSyncDurationBuckets = []float64{0.01, 0.025, 0.050, 0.075, 0.100, 0.300, 0.500}
// )

// // NewDeliveryControl is function
// func NewDeliveryControlSync(
// 	logger usecase.Logger,
// 	monitor *metrics.Monitor,
// 	queueHandler gateways.QueueHandler,
// 	deliveryControlUsecase usecase.DeliveryControl) DeliveryControlSync {
// 	instance := deliveryControlSync{
// 		logger:                 logger,
// 		monitor:                monitor,
// 		queueHandler:           queueHandler,
// 		deliveryControlUsecase: deliveryControlUsecase,
// 		wg:                     &sync.WaitGroup{},
// 	}
// 	monitor.Metrics.AddCounter(
// 		metricDeliveryControlSyncTotal, metricDeliveryControlSyncTotalDesc,
// 		metricDeliveryControlSyncTotalLabels)
// 	monitor.Metrics.AddHistogram(
// 		metricDeliveryControlSyncDuration, metricDeliveryControlSyncDurationDesc,
// 		metricDeliveryControlSyncDurationLabels, metricDeliveryControlSyncDurationBuckets)
// 	return &instance
// }

// func (d *deliveryControlSync) Start(ctx context.Context, wg *sync.WaitGroup) {
// 	maxMessages := config.Env.SQS.MaxMessages
// 	ch := make(chan gateways.QueueMessage, maxMessages)
// 	go d.queueHandler.Poll(ctx, wg, ch, maxMessages)
// 	go func() {
// 		defer func() {
// 			close(ch)
// 			// メインの処理でrecoverを実行しているため基本的にはここには到達しない想定
// 			// このログが出た場合は予期せぬ形でgoroutineが停止している可能性があるためアプリの再起動が必要
// 			if r := recover(); r != nil {
// 				d.logger.Error().Msgf("goroutine unrecoverable detail: %#v", r)
// 			}
// 		}()
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				d.wg.Wait()
// 				return
// 			default:
// 				for queueMessage := range ch {
// 					d.wg.Add(1)
// 					d.process(ctx, queueMessage)
// 					d.wg.Done()
// 				}
// 			}
// 		}
// 	}()
// }

// func (d *deliveryControlSync) process(ctx context.Context, queueMessage infra.QueueMessage) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			d.logger.Error().Msgf("Failed to process. %#v", r)
// 		}
// 	}()
// 	startTime := time.Now()
// 	var deliveryControlLog models.DeliveryControlLog
// 	message := *queueMessage.Message()
// 	messageID := queueMessage.MessageID()
// 	decoder := json.NewDecoder(strings.NewReader(message))
// 	if err := decoder.Decode(&deliveryControlLog); err != nil {
// 		// このログが出た場合はcloudwatch logsのmetric alarmでアラートを通知する
// 		d.logger.Error().Err(err).Str("message_id", *messageID).Str("body", message).Msg("Failed to parse message")
// 		d.queueHandler.UnprocessableMessage()
// 		d.queueHandler.DeleteMessage(ctx, queueMessage)
// 	} else {
// 		d.monitor.Metrics.
// 			GetCounter(metricDeliveryControlSyncTotal).
// 			WithLabelValues(deliveryControlLog.Event).Inc()
// 		err := d.deliveryControlUsecase.Process(ctx, startTime, &deliveryControlLog)
// 		latency := time.Since(startTime)
// 		if err != nil {
// 			d.logger.Error().
// 				Str("trace_id", deliveryControlLog.TraceID).
// 				Str("trace_time", deliveryControlLog.Time).
// 				Int("campaign_id", deliveryControlLog.CampaignID).
// 				Dur("latency", latency).
// 				Err(err).Str("body", message).Msg("Failed to process")
// 			d.queueHandler.UnprocessableMessage()
// 			// リランできるように SQS からは削除しない代わりに、ログ出力しておく
// 			d.queueHandler.OutputDeleteCliLog(queueMessage)
// 		} else {
// 			d.monitor.Metrics.
// 				GetHistogram(metricDeliveryControlSyncDuration).
// 				WithLabelValues("end_process").Observe(latency.Seconds())
// 			d.queueHandler.DeleteMessage(ctx, queueMessage)
// 		}
// 	}
// }

// func (d *deliveryControlSync) Close() {
// 	d.wg.Wait()
// }
