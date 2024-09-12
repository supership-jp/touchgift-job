package controllers

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra"
	"touchgift-job-manager/infra/metrics"
	"touchgift-job-manager/interface/gateways"
	"touchgift-job-manager/usecase"
)

type DeliveryOperationSync interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	Close()
}

type deliveryOperationSync struct {
	logger                   usecase.Logger
	monitor                  *metrics.Monitor
	queueHandler             gateways.QueueHandler
	deliveryOperationUsecase usecase.DeliveryOperation
	wg                       *sync.WaitGroup
}

// 　TODO: メトリクスちゃんとやる
var (
	metricDeliveryOperationSyncTotal       = "delivery_operation_sync_total"
	metricDeliveryOperationSyncTotalDesc   = "all delivery operation sync count"
	metricDeliveryOperationSyncTotalLabels = []string{"target", "event"}

	metricDeliveryOperationSyncDuration        = "delivery_operation_sync_duration_seconds"
	metricDeliveryOperationSyncDurationDesc    = "delivery operation sync processing time (seconds)"
	metricDeliveryOperationSyncDurationLabels  = []string{"kind"}
	metricDeliveryOperationSyncDurationBuckets = []float64{0.01, 0.025, 0.050, 0.075, 0.100, 0.300, 0.500}
)

func NewDeliveryOperationSync(
	logger usecase.Logger,
	monitor *metrics.Monitor,
	queueHandler gateways.QueueHandler,
	deliveryOperationUsecase usecase.DeliveryOperation) DeliveryOperationSync {
	instance := deliveryOperationSync{
		logger:                   logger,
		monitor:                  monitor,
		queueHandler:             queueHandler,
		deliveryOperationUsecase: deliveryOperationUsecase,
		wg:                       &sync.WaitGroup{},
	}
	monitor.Metrics.AddCounter(
		metricDeliveryOperationSyncTotal, metricDeliveryOperationSyncTotalDesc,
		metricDeliveryOperationSyncTotalLabels)
	monitor.Metrics.AddHistogram(
		metricDeliveryOperationSyncDuration, metricDeliveryOperationSyncDurationDesc,
		metricDeliveryOperationSyncDurationLabels, metricDeliveryOperationSyncDurationBuckets)
	return &instance
}

func (d *deliveryOperationSync) Start(ctx context.Context, wg *sync.WaitGroup) {
	maxMessages := config.Env.SQS.MaxMessages
	ch := make(chan gateways.QueueMessage, maxMessages)
	go d.queueHandler.Poll(ctx, wg, ch, maxMessages)
	go func() {
		defer func() {
			close(ch)
			if r := recover(); r != nil {
				d.logger.Error().Msgf("goroutine unrecoverable detail: %#v", r)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				d.wg.Done()
				return
			default:
				for queueMessage := range ch {
					startTime := time.Now()
					d.wg.Add(1)
					d.process(ctx, startTime, queueMessage)
					d.wg.Done()
				}
			}
		}
	}()
}

func (d *deliveryOperationSync) process(ctx context.Context, startTime time.Time, queueMessage infra.QueueMessage) {
	defer func() {
		if r := recover(); r != nil {
			d.logger.Error().Msgf("Failed to process %#v", r)
		}
		endLatency := time.Since(startTime)
		d.monitor.Metrics.GetHistogram(metricDeliveryOperationSyncDuration).
			WithLabelValues("end_process").Observe(endLatency.Seconds())
	}()
	var deliveryOperationLog models.DeliveryOperationLog
	message := *queueMessage.Message()
	messageID := queueMessage.MessageID()
	decoder := json.NewDecoder(strings.NewReader(message))
	if err := decoder.Decode(&deliveryOperationLog); err != nil {
		d.logger.Error().Err(err).Str("body", message).Msg("Failed to parse message")
		d.queueHandler.UnprocessableMessage()
		d.queueHandler.DeleteMessage(ctx, queueMessage)
	} else {
		process := func() error {
			for i := range deliveryOperationLog.CampaignLogs {
				current := time.Now()
				campaign := deliveryOperationLog.CampaignLogs[i]

				// d.monitor.Metrics.GetCounter(metricDeliveryOperationSyncTotal).
				// 	WithLabelValues(campaign.Event).Inc()
				if err := d.deliveryOperationUsecase.Process(ctx, current, &campaign); err != nil {
					if err == codes.ErrDoNothing {
						return nil
					}
					return err
				}
			}
			return nil
		}
		err := process()
		latency := time.Since(startTime)
		if err != nil {
			d.logger.Error().Dur("latency", latency).Err(err).Str("message_id", *messageID).Str("body", message).Msg("Failed to process")
			d.queueHandler.UnprocessableMessage()
			// リランできるように SQS からは削除しない代わりに、ログ出力しておく
			d.queueHandler.OutputDeleteCliLog(queueMessage)
		} else {
			d.queueHandler.DeleteMessage(ctx, queueMessage)
		}
	}
}

func (d *deliveryOperationSync) Close() {
	d.wg.Wait()
}
