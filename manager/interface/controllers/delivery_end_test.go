package controllers

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra/metrics"
	mock_controllers "touchgift-job-manager/mock/controllers"
	mock_gateways "touchgift-job-manager/mock/gateways"
	mock_repository "touchgift-job-manager/mock/repository"
	mock_usecase "touchgift-job-manager/mock/usecase"
)

func TestDeliveryEnd_Execute(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	createCampaign := func(id int, status string) *models.Campaign {
		return &models.Campaign{
			ID:        id,
			Status:    status,
			EndAt:     sql.NullTime{Time: time.Time{}, Valid: true},
			UpdatedAt: time.Time{},
		}
	}

	t.Run("データ無しの場合正常に処理する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryEnd
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryEnd := NewDeliveryEnd(
			logger,
			metrics.GetMonitor(),
			&configData,
			appTicker,
			transactionHandler,
			deliveryEndUsecase,
		)

		// mockの呼び出し定義(想定される呼び出し)
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}

		deliveryEndUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// started,pausedデータの処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"started", "paused"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return []models.Campaign{}, nil
			}).Times(1)
		// terminateデータの処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"terminate"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return []models.Campaign{}, nil
			}).Times(1)
		deliveryEndUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryEnd.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryEnd.Close()
	})

	t.Run("started,terminate両方ともデータ1件ありで正常に処理する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryEnd
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryEnd := NewDeliveryEnd(logger, metrics.GetMonitor(), &configData, appTicker, transactionHandler, deliveryEndUsecase)

		// mockの呼び出し定義(想定される呼び出し)
		campaigns := []*models.Campaign{createCampaign(1, "started")}
		campaignTerminates := []*models.Campaign{createCampaign(2, "terminate")}
		// basetime := time.Now().Truncate(time.Minute).Add(configData.TaskInterval)
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}

		deliveryEndUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// started,pausedの場合の処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"started", "paused"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return campaigns, nil
			}).Times(1)
		transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil).Times(1)
		deliveryEndUsecase.EXPECT().Terminate(
			gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaigns[0].ID), gomock.Eq(campaigns[0].UpdatedAt)).Return(1, nil).Times(1)
		tx.EXPECT().Commit().Return(nil).Times(1)
		deliveryEndUsecase.EXPECT().Reserve(gomock.Eq(ctx), gomock.Eq(campaigns[0].EndAt.Time), campaigns[0]).Return().Times(1)
		// terminateの場合の処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"terminate"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return campaignTerminates, nil
			}).Times(1)
		deliveryEndUsecase.EXPECT().ExecuteNow(campaignTerminates[0]).Return().Times(1)
		// 共通処理
		deliveryEndUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryEnd.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryEnd.Close()
	})

	t.Run("startedのみデータ1件ありで正常に処理する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryEnd
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryEnd := NewDeliveryEnd(logger, metrics.GetMonitor(), &configData, appTicker, transactionHandler, deliveryEndUsecase)

		// mockの呼び出し定義(想定される呼び出し)
		campaigns := []*models.Campaign{createCampaign(1, "started")}
		// basetime := time.Now().Truncate(time.Minute).Add(configData.TaskInterval)
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}

		deliveryEndUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// started,pausedの場合の処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"started", "paused"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return campaigns, nil
			}).Times(1)
		transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil).Times(1)
		deliveryEndUsecase.EXPECT().Terminate(
			gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaigns[0].ID), gomock.Eq(campaigns[0].UpdatedAt)).Return(1, nil).Times(1)
		tx.EXPECT().Commit().Return(nil).Times(1)
		deliveryEndUsecase.EXPECT().Reserve(gomock.Eq(ctx), gomock.Eq(campaigns[0].EndAt.Time), campaigns[0]).Return().Times(1)
		// terminateの場合の処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"terminate"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return []*models.Campaign{}, nil
			}).Times(1)
		// 共通処理
		deliveryEndUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryEnd.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryEnd.Close()
	})

	t.Run("terminateのみデータ1件ありで正常に処理する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryEnd
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryEnd := NewDeliveryEnd(logger, metrics.GetMonitor(), &configData, appTicker, transactionHandler, deliveryEndUsecase)

		// mockの呼び出し定義(想定される呼び出し)
		campaignTerminates := []*models.Campaign{createCampaign(2, "terminate")}
		// basetime := time.Now().Truncate(time.Minute).Add(configData.TaskInterval)
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}

		deliveryEndUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// started,pausedの場合の処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"started", "paused"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return []*models.Campaign{}, nil
			}).Times(1)
		// terminateの場合の処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"terminate"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return campaignTerminates, nil
			}).Times(1)
		deliveryEndUsecase.EXPECT().ExecuteNow(campaignTerminates[0]).Return().Times(1)
		// 共通処理
		deliveryEndUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryEnd.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryEnd.Close()
	})

	t.Run("startedデータ1件ありで状態更新でエラー起きた場合ロールバックして終了する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// mockの準備
		transactionHandler := mock_gateways.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		appTicker := mock_controllers.NewMockAppTicker(ctrl)

		// テスト設定準備
		configData := config.Env.DeliveryEnd
		configData.NumberOfConcurrent = 1
		testExecuteInterval := 1 * time.Second

		// テスト対象準備
		pctx := context.Background()
		ctx, cancel := context.WithCancel(pctx)
		deliveryEnd := NewDeliveryEnd(
			logger,
			metrics.GetMonitor(),
			&configData,
			appTicker,
			transactionHandler,
			deliveryEndUsecase,
		)

		// mockの呼び出し定義(想定される呼び出し)
		campaigns := []*models.Campaign{
			createCampaign(1, "started"),
		}
		// basetime := time.Now().Truncate(time.Minute).Add(configData.TaskInterval)
		expectedTo := func() time.Time {
			return time.Now().
				Truncate(time.Minute).
				Add(configData.TaskInterval).Add(10 * time.Second)
		}

		deliveryEndUsecase.EXPECT().CreateWorker(gomock.Eq(ctx)).Return().Times(1)
		appTicker.EXPECT().New(gomock.Eq(configData.TaskInterval), time.Minute).DoAndReturn(func(interval time.Duration, unit time.Duration) *time.Ticker {
			return NewAppTicker().New(testExecuteInterval, time.Second)
		}).Times(1)
		// started,pausedの処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"started", "paused"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, expectedTo(), to, 1*time.Second)
				return campaigns, nil
			}).Times(1)
		transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil).Times(1)
		deliveryEndUsecase.EXPECT().Terminate(
			gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaigns[0].ID), gomock.Eq(campaigns[0].UpdatedAt)).Return(0, errors.New("Failed to update")).Times(1)
		tx.EXPECT().Rollback().Return(nil).Times(1)
		// terminateの処理
		deliveryEndUsecase.EXPECT().GetDeliveryDataCampaigns(gomock.Eq(ctx), gomock.Any(), gomock.Eq([]string{"terminate"}), gomock.Eq(configData.TaskLimit)).
			DoAndReturn(func(ctx context.Context, to time.Time, status []string, limit int) ([]*models.Campaign, error) {
				assert.WithinDuration(t, time.Now().Truncate(time.Minute), to, 1*time.Second)
				return []*models.Campaign{}, nil
			}).Times(1)
		deliveryEndUsecase.EXPECT().Close().Return().Times(1)

		// 実行時間の調整
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Truncate(time.Second).Add(-50 * time.Millisecond)))
		// テストを実行する
		var wg sync.WaitGroup
		go deliveryEnd.StartMonitoring(ctx, &wg)

		// 非同期で処理が実行されるので待つ
		time.Sleep(time.Until(time.Now().Add(testExecuteInterval).Add(100 * time.Millisecond)))
		// 終了させる
		cancel()
		wg.Wait()
		deliveryEnd.Close()
	})
}
