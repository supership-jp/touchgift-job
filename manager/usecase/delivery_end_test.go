package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"
	"time"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"

	mock_repository "touchgift-job-manager/mock/repository"
	mock_usecase "touchgift-job-manager/mock/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// DeliveryEndのGetDeliveryDatasのテスト
func TestDeliveryEnd_GetCampaign(t *testing.T) {
	// 並列で実行する (mockを使っているため処理順番は気にせず実行できるように書ける)
	t.Parallel()

	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	t.Run("終了対象のキャンペーンがない場合空のキャンペーンを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のGetcampaignは、campaignRepository.GetCampaignToEnd を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		to := time.Now()
		limit := 1
		status := []string{"started"}
		condition := repository.CampaignDataToEndCondition{End: to, Status: status}
		// その際の戻り値
		expected := []*models.Campaign{}
		gomock.InOrder(
			campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(expected, nil).Times(1),
		)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &config.Env.DeliveryEnd, &config.Env.DeliveryEndUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		actual, err := deliveryEnd.GetDeliveryDataCampaigns(ctx, to, status, limit)
		if assert.NoError(t, err) {
			assert.Equal(t, len(expected), len(actual))
			assert.Exactly(t, expected, actual)
		}
	})

	t.Run("終了対象のキャンペーンがある場合そのキャンペーンを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のGetcampaignは、campaignRepository.GetCampaignToEnd を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		to := time.Now()
		status := []string{"started", "paused"}
		limit := 1
		condition := repository.CampaignDataToEndCondition{End: to, Status: status}
		// その際の戻り値
		expected := []*models.Campaign{
			{
				ID: 1,
			},
		}
		gomock.InOrder(
			campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(expected, nil).Times(1),
		)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &config.Env.DeliveryEnd, &config.Env.DeliveryEndUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		actual, err := deliveryEnd.GetDeliveryDataCampaigns(ctx, to, status, limit)
		if assert.NoError(t, err) {
			assert.Equal(t, len(expected), len(actual))
			assert.Exactly(t, expected, actual)
		}
	})

	t.Run("キャンペーン取得処理でエラーが発生した場合エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のGetcampaignは、campaignRepository.GetCampaignToEnd を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		to := time.Now()
		status := []string{"started"}
		limit := 1
		condition := repository.CampaignDataToEndCondition{End: to, Status: status}
		// その際の戻り値
		expected := errors.New("Failed")
		campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(nil, expected).Times(1)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &config.Env.DeliveryEnd, &config.Env.DeliveryEndUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		actual, err := deliveryEnd.GetDeliveryDataCampaigns(ctx, to, status, limit)
		if assert.Error(t, err) {
			assert.Nil(t, actual)
			assert.EqualError(t, err, expected.Error())
		}
	})
}

// DeliveryEndのTerminateのテスト
func TestDeliveryEnd_Terminate(t *testing.T) {
	// 並列で実行する (mockを使っているため処理順番は気にせず実行できるように書ける)
	t.Parallel()

	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	t.Run("terminate更新対象のデータがない場合、エラーは返さず何もしないで終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のTerminateは、deliveryControlUsecase.UpdateStatus を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		ID := 1
		updatedAt := time.Time{}
		status := "terminate"
		// 仮データ
		expected := 0
		// 何回呼ばれるか (Times)
		// を定義する
		updateCondition := &repository.UpdateCondition{
			CampaignID: ID,
			Status:     status,
			UpdatedAt:  updatedAt,
		}
		campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
			gomock.Eq(tx), updateCondition).Return(expected, nil).Times(1)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &config.Env.DeliveryEnd, &config.Env.DeliveryEndUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		actual, err := deliveryEnd.Terminate(ctx, tx, ID, updatedAt)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected, actual)
		}
	})

	t.Run("terminate更新対象のデータがある場合、statusを更新して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のTerminateは、deliveryControlUsecase.UpdateStatus を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		ID := 1
		status := "terminate"
		updatedAt := time.Time{}
		expected := 1
		// 何回呼ばれるか (Times)
		// を定義する
		updateCondition := &repository.UpdateCondition{
			CampaignID: ID,
			Status:     status,
			UpdatedAt:  updatedAt,
		}
		campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
			gomock.Eq(tx), updateCondition).Return(expected, nil).Times(1)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &config.Env.DeliveryEnd, &config.Env.DeliveryEndUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)

		actual, err := deliveryEnd.Terminate(ctx, tx, ID, updatedAt)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected, actual)
		}
	})

	t.Run("terminate更新処理でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のTerminateは、deliveryControlUsecase.UpdateStatus を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		ID := 1
		status := "terminate"
		updatedAt := time.Time{}

		// その際の戻り値
		expected := errors.New("Failed")
		// 何回呼ばれるか (Times)
		// を定義する
		updateCondition := &repository.UpdateCondition{
			CampaignID: ID,
			Status:     status,
			UpdatedAt:  updatedAt,
		}
		campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
			gomock.Eq(tx), gomock.Eq(updateCondition)).Return(0, expected).Times(1)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &config.Env.DeliveryEnd, &config.Env.DeliveryEndUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		_, err := deliveryEnd.Terminate(ctx, tx, ID, updatedAt)
		if assert.Error(t, err) {
			assert.EqualError(t, err, expected.Error())
		}
	})
}

// DeliveryEndのExecuteのテスト (terminate以外)
func TestDeliveryEnd_Execute_NotTERMINATE(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	// テスト用の設定
	configE := config.Env.DeliveryEnd
	configUsecase := config.Env.DeliveryEndUsecase
	// workerを1にする
	configUsecase.NumberOfConcurrent = 1
	configUsecase.NumberOfQueue = 1

	t.Run("terminate以外のステータスの場合、何もしないでログ出力して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		endAt := time.Now()
		campaign := models.Campaign{
			ID:        1,
			EndAt:     sql.NullTime{Time: endAt, Valid: true},
			UpdatedAt: time.Now(),
		}
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
		}
		// terminate以外の配信ステータス
		deliveryData := createEndTestCampaign(&campaign,
			campaign.EndAt.Time.Add(-10*time.Minute), sql.NullTime{Time: campaign.EndAt.Time.Add(5 * time.Minute)},
			"started", campaign.UpdatedAt.Add(1*time.Second))
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData, nil),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryEnd.CreateWorker(ctx)

		deliveryEnd.Reserve(ctx, time.Now(), &campaign) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryEnd.Close()
	})
}

// DeliveryEndのExecuteのテスト(配信終了)
func TestDeliveryEnd_Execute_DeliveryEnd(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	// テスト用の設定
	configE := config.Env.DeliveryEnd
	configUsecase := config.Env.DeliveryEndUsecase
	// workerを1にする
	configUsecase.NumberOfConcurrent = 1
	configUsecase.NumberOfQueue = 1
	// TODO: updateConditionの検証
	// status := "ended"
	// updateCondition := &repository.UpdateCondition{
	// 	CampaignID: deliveryData.ID,
	// 	Status:     status,
	// 	UpdatedAt:  time.Time{},
	// }
	t.Run("配信終了キャンペーンがterminateの場合はendedに更新して配信データを削除して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		// terminateにする前のデータ
		campaign := models.Campaign{ID: 1, GroupID: 2, StartAt: time.Now(), EndAt: sql.NullTime{Time: time.Now().Add(10 * time.Minute), Valid: true}, UpdatedAt: time.Now()}
		// terminateにした後のデータ
		deliveryData := createEndTestCampaign(&campaign, campaign.StartAt, campaign.EndAt, "terminate", campaign.UpdatedAt.Add(1*time.Second))
		status := "ended"
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
		}
		id := strconv.Itoa(deliveryData.ID)
		groupID := strconv.Itoa(deliveryData.GroupID)
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData, nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Any()).Return(1, nil),
			campaignDataRepository.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(&id)).Return(nil),
			contentDataRepository.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(&id)).Return(nil),
			campaignRepository.EXPECT().GetDeliveryCampaignCountByGroupID(gomock.Eq(ctx), gomock.Eq(deliveryData.GroupID)).Return(0, nil),
			touchPointDataRepository.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(&groupID)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlUsecase.EXPECT().Publish(
				gomock.Eq(ctx), gomock.Eq(deliveryData.ID), gomock.Eq(deliveryData.OrgCode), gomock.Eq(deliveryData.Status), gomock.Eq(status), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryEnd.CreateWorker(ctx)

		deliveryEnd.Reserve(ctx, time.Now(), &campaign) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryEnd.Close()
	})

	t.Run("groupIDに紐づく配信中のキャンペーンが存在した場合タッチポイントは削除しない", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		// terminateにする前のデータ
		campaign := models.Campaign{ID: 1, GroupID: 2, StartAt: time.Now(), EndAt: sql.NullTime{Time: time.Now().Add(10 * time.Minute), Valid: true}, UpdatedAt: time.Now()}
		// terminateにした後のデータ
		deliveryData := createEndTestCampaign(&campaign, campaign.StartAt, campaign.EndAt, "terminate", campaign.UpdatedAt.Add(1*time.Second))
		status := "ended"
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
		}
		id := strconv.Itoa(deliveryData.ID)
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData, nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Any()).Return(1, nil),
			campaignDataRepository.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(&id)).Return(nil),
			contentDataRepository.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(&id)).Return(nil),
			campaignRepository.EXPECT().GetDeliveryCampaignCountByGroupID(gomock.Eq(ctx), gomock.Eq(deliveryData.GroupID)).Return(1, nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlUsecase.EXPECT().Publish(
				gomock.Eq(ctx), gomock.Eq(deliveryData.ID), gomock.Eq(deliveryData.OrgCode), gomock.Eq(deliveryData.Status), gomock.Eq(status), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryEnd.CreateWorker(ctx)

		deliveryEnd.Reserve(ctx, time.Now(), &campaign) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryEnd.Close()
	})

	t.Run("配信終了キャンペーンのstatus変更でエラーが発生した場合、エラーを返してロールバックする", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)
		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		campaign := models.Campaign{
			ID:        1,
			EndAt:     sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: time.Now(),
		}
		deliveryData := createEndTestCampaign(&campaign, campaign.EndAt.Time.Add(-10*time.Minute), campaign.EndAt, "terminate", campaign.UpdatedAt.Add(1*time.Second))
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
		}
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData, nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Any()).Return(0, errors.New("Failed to update")),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テスト用の設定
		configE := config.Env.DeliveryEnd
		configUsecase := config.Env.DeliveryEndUsecase
		// workerを1にする
		configUsecase.NumberOfConcurrent = 1
		configUsecase.NumberOfQueue = 1

		// テストを実行する
		deliveryEnd := NewDeliveryEnd(
			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
			deliveryControlUsecase, campaignRepository, campaignDataRepository, contentDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryEnd.CreateWorker(ctx)

		deliveryEnd.Reserve(ctx, time.Now(), &campaign) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryEnd.Close()
	})
}

// // DeliveryEndのExecuteのテスト (直近終了予定)
// func TestDeliveryEnd_Execute_DeliveryMostRecentlyEnd(t *testing.T) {
// 	// テスト用のLoggerを作成
// 	logger := NewTestLogger(t)
// 	t.Run("直近終了予定のキャンペーンがterminateの場合はendedに更新してdelivery_dataを削除して終了", func(t *testing.T) {

// 		// mockを使用する準備
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

// 		// 必要なmockを作成
// 		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
// 		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
// 		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
// 		tx := mock_repository.NewMockTransaction(ctrl)
// 		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
// 		timer := mock_usecase.NewMockTimer(ctrl)

// 		// mockの処理を定義
// 		// テスト対象のexecuteは、 campaignDataRepository.Delete を使っているのでその処理を定義する
// 		// 引数に渡ると想定される値
// 		octx := context.Background()
// 		ctx, cancel := context.WithCancel(octx)
// 		campaign := models.Campaign{
// 			ID:        1,
// 			EndAt:     sql.NullTime{Time: time.Now(), Valid: true},
// 			UpdatedAt: time.Now(),
// 		}
// 		// config.Env.DeliveryEnd.TaskInterval が 1分のため、59*time.Second
// 		deliveryData := createEndTestCampaign(&campaign,
// 			campaign.EndAt.Time.Add(-5*time.Minute), sql.NullTime{Time: campaign.EndAt.Time.Add(59 * time.Second), Valid: true},
// 			"terminate", campaign.UpdatedAt.Add(1*time.Second))

// 		// テスト用の設定
// 		configE := config.Env.DeliveryEnd
// 		configUsecase := config.Env.DeliveryEndUsecase
// 		// workerを1にする
// 		configUsecase.NumberOfConcurrent = 1
// 		configUsecase.NumberOfQueue = 1
// 		deliveryEnd := NewDeliveryEnd(
// 			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
// 			deliveryControlUsecase, campaignRepository, campaignDataRepository)

// 		current := time.Now()
// 		status := "ended"
// 		condition := repository.CampaignCondition{
// 			CampaignID: campaign.ID,
// 		}
// 		updateCondition := &repository.UpdateCondition{
// 			CampaignID: deliveryData.ID,
// 			Status:     status,
// 			UpdatedAt:  time.Time{},
// 		}
// 		// 何回呼ばれるか (Times)
// 		// を定義する
// 		gomock.InOrder(
// 			timer.EXPECT().ExecuteAtTime(gomock.Eq(ctx), gomock.Eq(current), gomock.Any()).Do(func(ctx context.Context, specifiedTime time.Time, process func()) {
// 				process()
// 			}),
// 			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
// 			campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(deliveryData, nil),
// 			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(updateCondition)).Return(nil),
// 			tx.EXPECT().Commit().Return(nil),
// 			deliveryControlUsecase.EXPECT().Publish(
// 				gomock.Eq(ctx), gomock.Eq(deliveryData.ID), gomock.Eq(deliveryData.OrgCode), gomock.Eq(deliveryData.Status), gomock.Eq(status), gomock.Eq(""),
// 			),
// 		)

// 		// テストを実行する
// 		// Workerを使って実行するので作成
// 		deliveryEnd.CreateWorker(ctx)

// 		deliveryEnd.Reserve(ctx, current, &campaign) // 即時実行させる

// 		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
// 		// Workerを終了させる
// 		cancel()
// 		deliveryEnd.Close()
// 	})

// 	t.Run("直近終了予定のキャンペーンがterminateからendedの更新に失敗した場合はロールバックして終了", func(t *testing.T) {

// 		// mockを使用する準備
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

// 		// 必要なmockを作成
// 		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
// 		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
// 		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
// 		tx := mock_repository.NewMockTransaction(ctrl)
// 		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
// 		timer := mock_usecase.NewMockTimer(ctrl)

// 		// mockの処理を定義
// 		// テスト対象のexecuteは、 campaignDataRepository.Delete を使っているのでその処理を定義する
// 		// 引数に渡ると想定される値
// 		octx := context.Background()
// 		ctx, cancel := context.WithCancel(octx)
// 		campaign := models.Campaign{
// 			ID:        1,
// 			EndAt:     sql.NullTime{Time: time.Now(), Valid: true},
// 			UpdatedAt: time.Now(),
// 		}
// 		// config.Env.DeliveryEnd.TaskInterval が 1分のため、59*time.Second
// 		deliveryData := createEndTestCampaign(&campaign,
// 			campaign.EndAt.Time.Add(-5*time.Minute), sql.NullTime{Time: campaign.EndAt.Time.Add(59 * time.Second), Valid: true},
// 			"terminate", campaign.UpdatedAt.Add(1*time.Second))

// 		// テスト用の設定
// 		configE := config.Env.DeliveryEnd
// 		configUsecase := config.Env.DeliveryEndUsecase
// 		// workerを1にする
// 		configUsecase.NumberOfConcurrent = 1
// 		configUsecase.NumberOfQueue = 1
// 		deliveryEnd := NewDeliveryEnd(
// 			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
// 			deliveryControlUsecase, campaignRepository, campaignDataRepository)

// 		current := time.Now()
// 		status := "ended"
// 		condition := repository.CampaignCondition{
// 			CampaignID: campaign.ID,
// 		}
// 		updateCondition := &repository.UpdateCondition{
// 			CampaignID: deliveryData.ID,
// 			Status:     status,
// 			UpdatedAt:  time.Time{},
// 		}
// 		// 何回呼ばれるか (Times)
// 		// を定義する
// 		gomock.InOrder(
// 			timer.EXPECT().ExecuteAtTime(gomock.Eq(ctx), gomock.Eq(current), gomock.Any()).Do(func(ctx context.Context, specifiedTime time.Time, process func()) {
// 				process()
// 			}),
// 			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
// 			campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(deliveryData, nil),
// 			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(updateCondition)).Return(errors.New("Failed to update")),
// 			tx.EXPECT().Rollback().Return(nil),
// 		)

// 		// テストを実行する
// 		// Workerを使って実行するので作成
// 		deliveryEnd.CreateWorker(ctx)

// 		deliveryEnd.Reserve(ctx, current, &campaign) // 即時実行させる

// 		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
// 		// Workerを終了させる
// 		cancel()
// 		deliveryEnd.Close()
// 	})

// 	t.Run("直近終了予定のキャンペーンの削除に失敗した場合はロールバックして終了", func(t *testing.T) {

// 		// mockを使用する準備
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

// 		// 必要なmockを作成
// 		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
// 		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
// 		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
// 		tx := mock_repository.NewMockTransaction(ctrl)
// 		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
// 		timer := mock_usecase.NewMockTimer(ctrl)

// 		// mockの処理を定義
// 		// テスト対象のexecuteは、 campaignDataRepository.Delete を使っているのでその処理を定義する
// 		// 引数に渡ると想定される値
// 		octx := context.Background()
// 		ctx, cancel := context.WithCancel(octx)
// 		campaign := models.Campaign{
// 			ID:        1,
// 			EndAt:     sql.NullTime{Time: time.Now(), Valid: true},
// 			UpdatedAt: time.Now(),
// 		}
// 		// config.Env.DeliveryEnd.TaskInterval が 1分のため、59*time.Second
// 		deliveryData := createEndTestCampaign(&campaign,
// 			campaign.EndAt.Time.Add(-5*time.Minute), sql.NullTime{Time: campaign.EndAt.Time.Add(59 * time.Second), Valid: true},
// 			"terminate", campaign.UpdatedAt.Add(1*time.Second))

// 		// テスト用の設定
// 		configE := config.Env.DeliveryEnd
// 		configUsecase := config.Env.DeliveryEndUsecase
// 		// workerを1にする
// 		configUsecase.NumberOfConcurrent = 1
// 		configUsecase.NumberOfQueue = 1
// 		deliveryEnd := NewDeliveryEnd(
// 			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
// 			deliveryControlUsecase, campaignRepository, campaignDataRepository)

// 		current := time.Now()
// 		status := "ended"
// 		condition := repository.CampaignCondition{
// 			CampaignID: campaign.ID,
// 		}
// 		updateCondition := &repository.UpdateCondition{
// 			CampaignID: deliveryData.ID,
// 			Status:     status,
// 			UpdatedAt:  time.Time{},
// 		}
// 		// 何回呼ばれるか (Times)
// 		// を定義する
// 		gomock.InOrder(
// 			timer.EXPECT().ExecuteAtTime(gomock.Eq(ctx), gomock.Eq(current), gomock.Any()).Do(func(ctx context.Context, specifiedTime time.Time, process func()) {
// 				process()
// 			}),
// 			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
// 			campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(deliveryData, nil),
// 			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(updateCondition)).Return(errors.New("Failed to delete")),
// 			tx.EXPECT().Rollback().Return(nil),
// 		)

// 		// テストを実行する
// 		// Workerを使って実行するので作成
// 		deliveryEnd.CreateWorker(ctx)

// 		deliveryEnd.Reserve(ctx, current, &campaign) // 即時実行させる

// 		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
// 		// Workerを終了させる
// 		cancel()
// 		deliveryEnd.Close()
// 	})

// }

// // DeliveryEndのExecuteのテスト (今後終了予定)
// func TestDeliveryEnd_Execute_DeliveryEndFuture(t *testing.T) {
// 	// テスト用のLoggerを作成
// 	logger := NewTestLogger(t)
// 	t.Run("今後終了予定のキャンペーンがterminateの場合はendedに更新して終了", func(t *testing.T) {
// 		// mockを使用する準備
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

// 		// 必要なmockを作成
// 		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
// 		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
// 		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
// 		tx := mock_repository.NewMockTransaction(ctrl)
// 		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
// 		timer := NewTimer(logger)

// 		// mockの処理を定義
// 		// テスト対象のexecuteは、 campaignDataRepository.Delete を使っているのでその処理を定義する
// 		// 引数に渡ると想定される値
// 		octx := context.Background()
// 		ctx, cancel := context.WithCancel(octx)
// 		campaign := models.Campaign{
// 			ID:        1,
// 			EndAt:     sql.NullTime{Time: time.Now(), Valid: true},
// 			UpdatedAt: time.Now(),
// 		}
// 		// config.Env.DeliveryEnd.TaskInterval が 1分のため、5*time.Minute
// 		deliveryData := createEndTestCampaign(&campaign,
// 			campaign.EndAt.Time.Add(-5*time.Minute), sql.NullTime{Time: campaign.EndAt.Time.Add(5 * time.Minute), Valid: true},
// 			"terminate", campaign.UpdatedAt.Add(1*time.Second))
// 		status := "ended"
// 		condition := repository.CampaignCondition{
// 			CampaignID: campaign.ID,
// 		}
// 		updateCondition := &repository.UpdateCondition{
// 			CampaignID: deliveryData.ID,
// 			Status:     status,
// 			UpdatedAt:  time.Time{},
// 		}
// 		// 何回呼ばれるか (Times)
// 		// を定義する
// 		gomock.InOrder(
// 			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
// 			campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(deliveryData, nil),
// 			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(updateCondition)).Return(nil),
// 			tx.EXPECT().Commit().Return(nil),
// 			deliveryControlUsecase.EXPECT().Publish(
// 				gomock.Eq(ctx), gomock.Eq(deliveryData.ID), gomock.Eq(deliveryData.OrgCode), gomock.Eq(deliveryData.Status), gomock.Eq(status), gomock.Eq(""),
// 			),
// 		)

// 		// テスト用の設定
// 		configE := config.Env.DeliveryEnd
// 		configUsecase := config.Env.DeliveryEndUsecase
// 		// workerを1にする
// 		configUsecase.NumberOfConcurrent = 1
// 		configUsecase.NumberOfQueue = 1

// 		// テストを実行する
// 		deliveryEnd := NewDeliveryEnd(
// 			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
// 			deliveryControlUsecase, campaignRepository, campaignDataRepository)
// 		// Workerを使って実行するので作成
// 		deliveryEnd.CreateWorker(ctx)

// 		deliveryEnd.Reserve(ctx, time.Now(), &campaign) // 即時実行させる

// 		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
// 		// Workerを終了させる
// 		cancel()
// 		deliveryEnd.Close()
// 	})

// 	t.Run("今後終了予定のキャンペーンがありterminateからendedの更新に失敗した場合はロールバックして終了", func(t *testing.T) {
// 		// mockを使用する準備
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

// 		// 必要なmockを作成
// 		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
// 		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
// 		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
// 		tx := mock_repository.NewMockTransaction(ctrl)
// 		deliveryControlUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
// 		timer := NewTimer(logger)

// 		// mockの処理を定義
// 		// テスト対象のexecuteは、 campaignDataRepository.Delete を使っているのでその処理を定義する
// 		// 引数に渡ると想定される値
// 		octx := context.Background()
// 		ctx, cancel := context.WithCancel(octx)
// 		campaign := models.Campaign{
// 			ID:        1,
// 			EndAt:     sql.NullTime{Time: time.Now(), Valid: true},
// 			UpdatedAt: time.Now(),
// 		}
// 		// config.Env.DeliveryEnd.TaskInterval が 1分のため、5*time.Minute
// 		deliveryData := createEndTestCampaign(&campaign,
// 			campaign.EndAt.Time.Add(-5*time.Minute), sql.NullTime{Time: campaign.EndAt.Time.Add(5 * time.Minute), Valid: true},
// 			"terminate", campaign.UpdatedAt.Add(1*time.Second))
// 		status := "ended"
// 		condition := repository.CampaignCondition{
// 			CampaignID: campaign.ID,
// 		}
// 		updateCondition := &repository.UpdateCondition{
// 			CampaignID: deliveryData.ID,
// 			Status:     status,
// 			UpdatedAt:  time.Time{},
// 		}
// 		// 何回呼ばれるか (Times)
// 		// を定義する
// 		gomock.InOrder(
// 			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
// 			campaignRepository.EXPECT().GetCampaignToEnd(gomock.Eq(ctx), gomock.Eq(&condition)).Return(deliveryData, nil),
// 			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(updateCondition)).Return(errors.New("Failed to update")),
// 			tx.EXPECT().Rollback().Return(nil),
// 		)

// 		// テスト用の設定
// 		configE := config.Env.DeliveryEnd
// 		configUsecase := config.Env.DeliveryEndUsecase
// 		// workerを1にする
// 		configUsecase.NumberOfConcurrent = 1
// 		configUsecase.NumberOfQueue = 1

// 		// テストを実行する
// 		deliveryEnd := NewDeliveryEnd(
// 			logger, metrics.GetMonitor(), &configE, &configUsecase, transactionHandler, timer,
// 			deliveryControlUsecase, campaignRepository, campaignDataRepository)
// 		// Workerを使って実行するので作成
// 		deliveryEnd.CreateWorker(ctx)

// 		deliveryEnd.Reserve(ctx, time.Now(), &campaign) // 即時実行させる

// 		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
// 		// Workerを終了させる
// 		cancel()
// 		deliveryEnd.Close()
// 	})
// }

func createEndTestCampaign(campaign *models.Campaign, startAt time.Time, endAt sql.NullTime, status string, updatedAt time.Time) *models.Campaign {
	return &models.Campaign{
		ID:        campaign.ID,
		Status:    status,
		StartAt:   startAt,
		EndAt:     endAt,
		UpdatedAt: updatedAt,
		OrgCode:   "org1",
	}
}
