package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"log"
	"sync"
	"testing"
	"time"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra/metrics"
	"touchgift-job-manager/interface/gateways"
	mock_gateways "touchgift-job-manager/mock/gateways"
	mock_infra "touchgift-job-manager/mock/infra"
	mock_usecase "touchgift-job-manager/mock/usecase"
)

func TestDeliveryOperationSync_Start(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Run("Campaignsログがない場合何もしない", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queueHandler := mock_gateways.NewMockQueueHandler(ctrl)
		deliveryOperationUsecase := mock_usecase.NewMockDeliveryOperation(ctrl)
		queueMessage := mock_infra.NewMockQueueMessage(ctrl)

		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		wg := sync.WaitGroup{}

		gomock.InOrder(
			queueHandler.EXPECT().Poll(gomock.Eq(ctx), gomock.Eq(&wg), gomock.Any(), gomock.Eq(config.Env.SQS.MaxMessages)).Do(
				func(ctx context.Context, wg *sync.WaitGroup, ch chan gateways.QueueMessage, maxMessages int64) {
					ch <- queueMessage
				}),
			queueMessage.EXPECT().Message().DoAndReturn(func() *string {
				json := `{
					"time": "2021-10-01T10:00:00.000Z",
					"type": "delivery_operation",
					"campaigns":[]
				}`
				return &json
			}),
			queueMessage.EXPECT().MessageID().DoAndReturn(func() *string {
				messageID := "messageID1"
				return &messageID
			}),
			queueHandler.EXPECT().DeleteMessage(gomock.Eq(ctx), gomock.Eq(queueMessage)),
		)

		// テスト実行
		deliveryOperationSync := NewDeliveryOperationSync(logger, metrics.GetMonitor(), queueHandler, deliveryOperationUsecase)
		deliveryOperationSync.Start(ctx, &wg)
		time.Sleep(50 * time.Millisecond)

		// テスト完了待ち
		cancel()
		deliveryOperationSync.Close()
	})
	t.Run("SQSからのメッセージがパースエラーの場合、", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queueHandler := mock_gateways.NewMockQueueHandler(ctrl)
		deliveryOperationUsecase := mock_usecase.NewMockDeliveryOperation(ctrl)
		queueMessage := mock_infra.NewMockQueueMessage(ctrl)

		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		wg := sync.WaitGroup{}
		jsonText := `{
			"time": "2021-10-01T10:00:00.000Z",
			"type": "delivery_operation",
		}`
		gomock.InOrder(
			queueHandler.EXPECT().Poll(gomock.Eq(ctx), gomock.Eq(&wg), gomock.Any(), gomock.Eq(config.Env.SQS.MaxMessages)).Do(
				func(ctx context.Context, wg *sync.WaitGroup, ch chan gateways.QueueMessage, maxMessages int64) {
					ch <- queueMessage
				}),
			queueMessage.EXPECT().Message().DoAndReturn(func() *string {
				return &jsonText
			}),
			queueMessage.EXPECT().MessageID().DoAndReturn(func() *string {
				messageID := "messageID1"
				return &messageID
			}),
			queueHandler.EXPECT().UnprocessableMessage(),
			queueHandler.EXPECT().DeleteMessage(gomock.Eq(ctx), gomock.Eq(queueMessage)),
		)

		// テスト実行
		deliveryOperationSync := NewDeliveryOperationSync(logger, metrics.GetMonitor(), queueHandler, deliveryOperationUsecase)
		deliveryOperationSync.Start(ctx, &wg)
		time.Sleep(50 * time.Millisecond)

		// テスト完了待ち
		cancel()
		deliveryOperationSync.Close()
	})
	t.Run("campaignsログがある場合、通常の処理を行う", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queueHandler := mock_gateways.NewMockQueueHandler(ctrl)
		deliveryOperationUsecase := mock_usecase.NewMockDeliveryOperation(ctrl)
		queueMessage := mock_infra.NewMockQueueMessage(ctrl)

		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		wg := sync.WaitGroup{}
		jsonText := `{
			"time": "2021-10-01T10:00:00.000Z",
			"type": "delivery_operation",
			"target": "touch-gift",
			"campaigns":[{
				"id": 1, "origin_id": "origin_id1", "event": "insert", "Budget": 100 
			}]
		}`
		var deliveryOperationLog models.DeliveryOperationLog
		if err := json.Unmarshal([]byte(jsonText), &deliveryOperationLog); err != nil {
			log.Fatal(err)
		}
		gomock.InOrder(
			queueHandler.EXPECT().Poll(gomock.Eq(ctx), gomock.Eq(&wg), gomock.Any(), gomock.Eq(config.Env.SQS.MaxMessages)).Do(
				func(ctx context.Context, wg *sync.WaitGroup, ch chan gateways.QueueMessage, maxMessages int64) {
					ch <- queueMessage
				}),
			queueMessage.EXPECT().Message().DoAndReturn(func() *string {
				return &jsonText
			}),
			queueMessage.EXPECT().MessageID().DoAndReturn(func() *string {
				messageID := "messageID1"
				return &messageID
			}),
			deliveryOperationUsecase.EXPECT().Process(gomock.Eq(ctx),
				gomock.Any(), gomock.Eq(&deliveryOperationLog.CampaignLogs[0])).Return(nil),
			queueHandler.EXPECT().DeleteMessage(gomock.Eq(ctx), gomock.Eq(queueMessage)),
		)

		// テスト実行
		deliveryOperationSync := NewDeliveryOperationSync(logger, metrics.GetMonitor(), queueHandler, deliveryOperationUsecase)
		deliveryOperationSync.Start(ctx, &wg)
		time.Sleep(50 * time.Millisecond)

		// テスト完了待ち
		cancel()
		deliveryOperationSync.Close()

	})
	t.Run("campaignsログの処理でエラーが起きた場合、エラーを返して終了する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queueHandler := mock_gateways.NewMockQueueHandler(ctrl)
		deliveryOperationUsecase := mock_usecase.NewMockDeliveryOperation(ctrl)
		queueMessage := mock_infra.NewMockQueueMessage(ctrl)

		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		wg := sync.WaitGroup{}
		jsonText := `{
			"time": "2021-10-01T10:00:00.000Z",
			"type": "delivery_operation",
			"target": "touch-gift",
			"campaigns":[{
				"id": 1, "origin_id": "origin_id1", "event": "insert", "Budget": 100 
			}]
		}`
		var deliveryOperationLog models.DeliveryOperationLog
		if err := json.Unmarshal([]byte(jsonText), &deliveryOperationLog); err != nil {
			log.Fatal(err)
		}
		gomock.InOrder(
			queueHandler.EXPECT().Poll(gomock.Eq(ctx), gomock.Eq(&wg), gomock.Any(), gomock.Eq(config.Env.SQS.MaxMessages)).Do(
				func(ctx context.Context, wg *sync.WaitGroup, ch chan gateways.QueueMessage, maxMessages int64) {
					ch <- queueMessage
				}),
			queueMessage.EXPECT().Message().DoAndReturn(func() *string {
				return &jsonText
			}),
			queueMessage.EXPECT().MessageID().DoAndReturn(func() *string {
				messageID := "messageID1"
				return &messageID
			}),
			deliveryOperationUsecase.EXPECT().Process(gomock.Eq(ctx),
				gomock.Any(), gomock.Eq(&deliveryOperationLog.CampaignLogs[0])).Return(errors.New("Failed to process")),
			queueHandler.EXPECT().UnprocessableMessage(),
			queueHandler.EXPECT().OutputDeleteCliLog(gomock.Eq(queueMessage)),
		)

		// テスト実行
		deliveryOperationSync := NewDeliveryOperationSync(logger, metrics.GetMonitor(), queueHandler, deliveryOperationUsecase)
		deliveryOperationSync.Start(ctx, &wg)
		time.Sleep(50 * time.Millisecond)

		// テスト完了待ち
		cancel()
		deliveryOperationSync.Close()
	})
}
