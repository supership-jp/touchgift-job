package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	mock_notification "touchgift-job-manager/mock/notification"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryControlEvent_Publish(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	t.Run("配信制御イベントの通知ができること", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		current := time.Now()
		cacheOperation := "NONE"
		deliveryControlEvent := models.DeliveryControlLog{
			TraceID:        "traceid1",
			Time:           current.Format(time.RFC3339Nano),
			Version:        1,
			CacheOperation: cacheOperation,
			Event:          "warmup",
			Source:         "touchgift-job-manager",
			OrgCode:        "org1",
			CampaignID:     1,
		}
		_, err := json.Marshal(&deliveryControlEvent)
		if err != nil {
			t.Fatal("failed to marshal json")
		}
		messageAttributes := map[string]string{
			"event":           deliveryControlEvent.Event,
			"cache_operation": deliveryControlEvent.CacheOperation,
		}

		messageID := "message_id1"
		gomock.InOrder(
			notificationHandler.EXPECT().Publish(ctx, gomock.Any(), messageAttributes).Return(&messageID, nil),
		)

		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		deliveryControlEventUsecase.Publish(ctx, 1, "org1", "configured", "warmup", "")
	})

	t.Run("通知に失敗した場合は専用のエラーログを出力する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		current := time.Now()
		cacheOperation := "PUT"
		deliveryControlEvent := models.DeliveryControlLog{
			TraceID:        "traceid1",
			Time:           current.Format(time.RFC3339Nano),
			Version:        1,
			CacheOperation: cacheOperation,
			Event:          "start",
			Source:         "touchgift-job-manager",
			OrgCode:        "org1",
			CampaignID:     1,
		}
		_, err := json.Marshal(&deliveryControlEvent)
		if err != nil {
			t.Fatal("failed to marshal json")
		}
		messageAttributes := map[string]string{
			"event":           deliveryControlEvent.Event,
			"cache_operation": deliveryControlEvent.CacheOperation,
		}
		errUnexpected := errors.New("unexpected error")
		gomock.InOrder(
			notificationHandler.EXPECT().Publish(ctx, gomock.Any(), messageAttributes).Return(nil, errUnexpected),
		)

		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		deliveryControlEventUsecase.Publish(ctx, 1, "org1", "warmup", "started", "")
	})
}

// DeliveryControlEventのcreateDeliveryControlLogのテスト
func TestDeliveryControlEvent_createDeliveryControlLog(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("引数を渡してdelivery_control_logを作成できる", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		CampaignID := 1
		organization := "org"
		cache := "NONE"
		before := "configured"
		after := "warmup"
		expectedEvent := "warmup"
		expectedEventDetail := "shortage"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual := deliveryControlEventInteractor.createDeliveryControlLog(CampaignID, organization, before, after, codes.DetailShortage)
		assert.NotEmpty(t, actual.TraceID)
		assert.Equal(t, CampaignID, actual.CampaignID)
		assert.Equal(t, organization, actual.OrgCode)
		assert.Equal(t, expectedEvent, actual.Event)
		assert.Equal(t, cache, actual.CacheOperation)
		assert.Equal(t, expectedEventDetail, actual.EventDetail)
		assert.Equal(t, config.Env.Version, actual.Version)
		assert.NotEmpty(t, actual.Time)
	})

}

// DeliveryControlEventのdeliveryEventのテスト
func TestDeliveryControlEvent_deliveryEvent(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("campaignのstatus遷移がconfigured->warmupの場合、配信制御イベントはwarmupを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		expected := "warmup"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual, operation := deliveryControlEventInteractor.deliveryEvent("configured", "warmup")
		assert.Exactly(t, expected, actual)
		assert.Exactly(t, "NONE", operation)
	})

	t.Run("campaignのstatus遷移がwarmup->startedの場合、配信制御イベントはstartを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		expected := "start"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual, operation := deliveryControlEventInteractor.deliveryEvent("warmup", "started")
		assert.Exactly(t, expected, actual)
		assert.Exactly(t, "PUT", operation)
	})

	t.Run("campaignのstatus遷移がresume->startedの場合、配信制御イベントはresumeを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		expected := "resume"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual, operation := deliveryControlEventInteractor.deliveryEvent("resume", "started")
		assert.Exactly(t, expected, actual)
		assert.Exactly(t, "PUT", operation)
	})

	t.Run("campaignのstatus遷移がstarted->startedの場合、配信制御イベントはupdateを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		expected := "update"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual, operation := deliveryControlEventInteractor.deliveryEvent("started", "started")
		assert.Exactly(t, expected, actual)
		assert.Exactly(t, "PUT", operation)
	})

	t.Run("campaignのstatus遷移がpausedの場合、配信制御イベントはpauseを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		expected := "pause"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual, operation := deliveryControlEventInteractor.deliveryEvent("pause", "paused")
		assert.Exactly(t, expected, actual)
		assert.Exactly(t, "DELETE", operation)
	})

	t.Run("campaignのstatus遷移がstop->stoppedの場合、配信制御イベントはstopを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		expected := "stop"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual, operation := deliveryControlEventInteractor.deliveryEvent("stop", "stopped")
		assert.Exactly(t, expected, actual)
		assert.Exactly(t, "DELETE", operation)
	})

	t.Run("campaignのstatus遷移がterminate->endedの場合、配信制御イベントはendを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		notificationHandler := mock_notification.NewMockNotificationHandler(ctrl)

		expected := "end"
		// テストを実行する
		deliveryControlEventUsecase := NewDeliveryControlEvent(logger, notificationHandler)
		// private methodのテストを行うためにcastする
		deliveryControlEventInteractor := deliveryControlEventUsecase.(*deliveryControlEvent)
		actual, operation := deliveryControlEventInteractor.deliveryEvent("terminate", "ended")
		assert.Exactly(t, expected, actual)
		assert.Exactly(t, "DELETE", operation)
	})
}
