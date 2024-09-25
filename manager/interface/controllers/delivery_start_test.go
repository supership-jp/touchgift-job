package controllers

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra/metrics"
	mock_controllers "touchgift-job-manager/mock/controllers"
	mock_gateways "touchgift-job-manager/mock/gateways"
	mock_repository "touchgift-job-manager/mock/repository"
	mock_usecase "touchgift-job-manager/mock/usecase"
)

func TestDeliveryStart_Execute(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	createCampaign := func(id int, status string) *models.Campaign {
		return &models.Campaign{
			ID:        id,
			GroupID:   1,
			Status:    status,
			StartAt:   time.Time{},
			UpdatedAt: time.Time{},
			OrgCode:   "org1",
		}
	}

	t.Run("データ無しの場合正常に処理する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryStart
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryStart := NewDeliveryStart(
			logger,
			metrics.GetMonitor(),
			&configData,
			appTicker,
			transactionHandler,
			deliveryStartUsecase,
			deliveryControlEvent,
		)

		// mockの呼び出し定義(想定される呼び出し)
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}

		deliveryStartUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// configuredデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("configured"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return []models.Campaign{}, nil
			}).Times(1)
		// warmupデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("warmup"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return []models.Campaign{}, nil
			}).Times(1)
		deliveryStartUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryStart.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryStart.Close()
	})

	t.Run("configured,warmup両方ともデータ1件ありの場合正常にする", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		// テスト設定準備
		configData := config.Env.DeliveryStart
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryStart := NewDeliveryStart(
			logger,
			metrics.GetMonitor(),
			&configData,
			appTicker,
			transactionHandler,
			deliveryStartUsecase,
			deliveryControlEvent,
		)

		// mockの呼び出し定義(想定される呼び出し)
		campaigns := []*models.Campaign{
			createCampaign(1, "configured"),
		}
		campaignWarmups := []*models.Campaign{
			createCampaign(2, "warmup"),
		}
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}
		deliveryStartUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// configuredデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("configured"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return campaigns, nil
			}).Times(1)
		transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil).Times(1)
		deliveryStartUsecase.EXPECT().UpdateStatus(
			gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaigns[0]), codes.StatusWarmup).Return(1, nil).Times(1)
		tx.EXPECT().Commit().Return(nil).Times(1)
		deliveryControlEvent.EXPECT().PublishCampaignEvent(
			gomock.Eq(ctx), gomock.Eq(campaigns[0].ID), gomock.Eq(campaigns[0].GroupID), gomock.Eq(campaigns[0].OrgCode),
			gomock.Eq("configured"), gomock.Eq("warmup"), gomock.Eq(""),
		).Times(1)
		deliveryStartUsecase.EXPECT().Reserve(gomock.Eq(ctx), gomock.Eq(campaigns[0].StartAt), gomock.Eq(campaigns[0])).Return().Times(1)
		// warmupデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("warmup"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return campaignWarmups, nil
			}).Times(1)
		deliveryStartUsecase.EXPECT().ExecuteNow(campaignWarmups[0]).Return().Times(1)
		// 共通
		deliveryStartUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryStart.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryStart.Close()
	})

	t.Run("configuredのみデータ1件ありの場合正常に処理する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryStart
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryStart := NewDeliveryStart(
			logger,
			metrics.GetMonitor(),
			&configData,
			appTicker,
			transactionHandler,
			deliveryStartUsecase,
			deliveryControlEvent,
		)

		// mockの呼び出し定義(想定される呼び出し)
		campaigns := []*models.Campaign{
			createCampaign(1, "configured"),
		}
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}
		deliveryStartUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// configuredデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("configured"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return campaigns, nil
			}).Times(1)
		transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil).Times(1)
		deliveryStartUsecase.EXPECT().UpdateStatus(
			gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaigns[0]), gomock.Eq(codes.StatusWarmup)).Return(1, nil).Times(1)
		tx.EXPECT().Commit().Return(nil).Times(1)
		deliveryControlEvent.EXPECT().PublishCampaignEvent(
			gomock.Eq(ctx), gomock.Eq(campaigns[0].ID), gomock.Eq(campaigns[0].GroupID), gomock.Eq(campaigns[0].OrgCode),
			gomock.Eq("configured"), gomock.Eq("warmup"), gomock.Eq(""),
		).Times(1)
		deliveryStartUsecase.EXPECT().Reserve(gomock.Eq(ctx), gomock.Eq(campaigns[0].StartAt), campaigns[0]).Return().Times(1)
		// warmupデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("warmup"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return []*models.Campaign{}, nil
			}).Times(1)
		// 共通
		deliveryStartUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryStart.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryStart.Close()
	})

	t.Run("warmpupのみデータ1件ありの場合正常に処理", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryStart
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryStart := NewDeliveryStart(
			logger,
			metrics.GetMonitor(),
			&configData,
			appTicker,
			transactionHandler,
			deliveryStartUsecase,
			deliveryControlEvent,
		)

		// mockの呼び出し定義(想定される呼び出し)
		campaignWarmups := []*models.Campaign{
			createCampaign(2, "warmup"),
		}
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}
		deliveryStartUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// configuredデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("configured"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return []*models.Campaign{}, nil
			}).Times(1)
		// warmupデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("warmup"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return campaignWarmups, nil
			}).Times(1)
		deliveryStartUsecase.EXPECT().ExecuteNow(campaignWarmups[0]).Return().Times(1)
		// 共通
		deliveryStartUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryStart.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryStart.Close()
	})

	t.Run("configuredのみデータ1件ありで状態更新でエラーが起きた場合ロールバックする", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryStart
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryStart := NewDeliveryStart(
			logger,
			metrics.GetMonitor(),
			&configData,
			appTicker,
			transactionHandler,
			deliveryStartUsecase,
			deliveryControlEvent,
		)

		// mockの呼び出し定義(想定される呼び出し)
		campaigns := []*models.Campaign{
			createCampaign(1, "configured"),
		}
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}
		deliveryStartUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// configuredデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("configured"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return campaigns, nil
			}).Times(1)
		transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil).Times(1)
		deliveryStartUsecase.EXPECT().UpdateStatus(
			gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaigns[0]), gomock.Eq(codes.StatusWarmup)).Return(0, errors.New("Failed to update")).Times(1)
		tx.EXPECT().Rollback().Return(nil).Times(1)
		// warmupデータの処理
		deliveryStartUsecase.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Any(), gomock.Eq("warmup"), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return []*models.Campaign{}, nil
			}).Times(1)
		// 共通
		deliveryStartUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryStart.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryStart.Close()
	})

}
