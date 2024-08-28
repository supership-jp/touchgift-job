package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"

	mock_repository "touchgift-job-manager/mock/repository"
	mock_usecase "touchgift-job-manager/mock/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// DeliveryStartのGetcampaignDataのテスト
func TestDeliveryStart_GetCampaignData(t *testing.T) {
	// 並列で実行する (mockを使っているため処理順番は気にせず実行できるように書ける)
	t.Parallel()

	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	t.Run("開始対象のキャンペーンがない場合、空のキャンペーンを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のGetcampaignDataは、campaignRepository.GetCampaignToStart を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		tx := mock_repository.NewMockTransaction(ctrl)
		to := time.Now()
		status := "configured"
		limit := 1
		condition := repository.CampaignToStartCondition{To: to, Status: status}
		// その際の戻り値
		expected := []*models.Campaign{}
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(expected, nil).Times(1),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &config.Env.DeliveryStart, &config.Env.DeliveryStartUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		actual, err := deliveryStart.GetCampaignToStart(ctx, to, status, limit)
		if assert.NoError(t, err) {
			assert.Equal(t, len(expected), len(actual))
			assert.Exactly(t, expected, actual)
		}
	})

	t.Run("開始対象のキャンペーンがある場合、キャンペーンを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のGetcampaignDataは、campaignRepository.GetCampaignToStart を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		tx := mock_repository.NewMockTransaction(ctrl)
		to := time.Now()
		limit := 1
		status := "configured"
		condition := repository.CampaignToStartCondition{To: to, Status: status}
		// その際の戻り値
		expected := []*models.Campaign{
			{
				ID: 1,
			},
		}
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(expected, nil).Times(1),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &config.Env.DeliveryStart, &config.Env.DeliveryStartUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		actual, err := deliveryStart.GetCampaignToStart(ctx, to, status, limit)
		if assert.NoError(t, err) {
			assert.Equal(t, len(expected), len(actual))
			assert.Exactly(t, expected, actual)
		}
	})

	t.Run("キャンペーン取得処理でエラー発生した場合、エラーを返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)

		// mockの処理を定義
		// テスト対象のGetcampaignDataは、campaignRepository.GetCampaignToStart を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		tx := mock_repository.NewMockTransaction(ctrl)
		to := time.Now()
		limit := 1
		status := "configured"
		condition := repository.CampaignToStartCondition{To: to, Status: status}
		// その際の戻り値
		expected := errors.New("Failed")
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetCampaignToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(nil, expected).Times(1),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &config.Env.DeliveryStart, &config.Env.DeliveryStartUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		actual, err := deliveryStart.GetCampaignToStart(ctx, to, status, limit)
		if assert.Error(t, err) {
			assert.Nil(t, actual)
			assert.EqualError(t, err, expected.Error())
		}
	})
}

// DeliveryStartのUpdateStatusのテスト
func TestDeliveryStart_UpdateStatus(t *testing.T) {
	// 並列で実行する (mockを使っているため処理順番は気にせず実行できるように書ける)
	t.Parallel()

	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	t.Run("更新対象のデータがない場合、データ更新件数の0が返却されること", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のUpdateStatusは、deliveryControlEventUsecase.UpdateStatus を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		ID := 1
		campaignData := &models.Campaign{
			ID:        ID,
			Status:    "configured",
			StartAt:   time.Time{},
			UpdatedAt: time.Time{},
			OrgCode:   "org1",
		}
		// 仮データ
		expected := 0
		// 何回呼ばれるか (Times)
		// を定義する
		condition := repository.UpdateCondition{
			CampaignID: campaignData.ID,
			Status:     codes.StatusWarmup,
			UpdatedAt:  campaignData.UpdatedAt,
		}
		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(expected, nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &config.Env.DeliveryStart, &config.Env.DeliveryStartUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		count, err := deliveryStart.UpdateStatus(ctx, tx, campaignData, codes.StatusWarmup)
		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("warmupに1件データを更新した場合、変更したデータ件数である1が返却されること", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のUpdateStatusは、deliveryControlEventUsecase.UpdateStatus を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		ID := 1
		status := codes.StatusWarmup
		before := "configured"
		organization := "org1"
		expected := 1
		campaignData := &models.Campaign{
			ID:        ID,
			Status:    before,
			StartAt:   time.Time{},
			UpdatedAt: time.Time{},
			OrgCode:   organization,
		}
		// 何回呼ばれるか
		// を定義する
		condition := repository.UpdateCondition{
			CampaignID: campaignData.ID,
			Status:     status,
			UpdatedAt:  campaignData.UpdatedAt,
		}
		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(expected, nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &config.Env.DeliveryStart, &config.Env.DeliveryStartUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)

		count, err := deliveryStart.UpdateStatus(ctx, tx, campaignData, codes.StatusWarmup)
		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("warmup更新処理中にエラーの場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のUpdateStatusは、deliveryControlEventUsecase.UpdateStatus を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		ID := 1
		status := codes.StatusWarmup
		campaignData := &models.Campaign{
			ID:        ID,
			Status:    "configured",
			StartAt:   time.Time{},
			UpdatedAt: time.Time{},
			OrgCode:   "org1",
		}
		condition := repository.UpdateCondition{
			CampaignID: campaignData.ID,
			Status:     status,
			UpdatedAt:  campaignData.UpdatedAt,
		}

		// その際の戻り値
		expected := errors.New("Failed")
		// 何回呼ばれるか (Times)
		// を定義する
		campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
			gomock.Eq(tx), gomock.Eq(&condition)).Return(0, expected).Times(1)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &config.Env.DeliveryStart, &config.Env.DeliveryStartUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		count, err := deliveryStart.UpdateStatus(ctx, tx, campaignData, codes.StatusWarmup)
		assert.EqualError(t, err, "Failed to update status. status: warmup: Failed")
		assert.Equal(t, 0, count)
	})

}

// DeliveryStartのExecuteのテスト (warmup以外)
func TestDeliveryStart_Execute_NotWARMUP(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)

	t.Run("warmup以外のステータスの場合ログ出力して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		campaignData := []*models.Campaign{
			{
				ID:        1,
				StartAt:   time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		deliveryData := createStartTestCampaign(
			campaignData[0], campaignData[0].StartAt.Add(5*time.Minute), sql.NullTime{}, codes.StatusStarted, campaignData[0].UpdatedAt.Add(1*time.Second))
		condition := repository.CampaignCondition{CampaignID: campaignData[0].ID, Status: codes.StatusWarmup}
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			tx.EXPECT().Rollback().Return(nil),
		)
		// テスト用の設定
		configS := config.Env.DeliveryStart
		configUsecase := config.Env.DeliveryStartUsecase
		// workerを1にする
		configUsecase.NumberOfConcurrent = 1
		configUsecase.NumberOfQueue = 1

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), campaignData[0]) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})
}

// DeliveryStartのExecuteのテスト (配信開始)
func TestDeliveryStart_Execute_DeliveryStart(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	// テスト用データ
	campaignData := models.Campaign{ID: 1, GroupID: 1, StartAt: time.Now(), UpdatedAt: time.Now()}
	deliveryData := createStartTestCampaign(
		&campaignData, campaignData.StartAt, sql.NullTime{}, codes.StatusWarmup, campaignData.UpdatedAt.Add(1*time.Second),
	)
	creatives := []*models.Creative{{ID: 1}}
	cc := []*models.CampaignCreative{{ID: creatives[0].ID}}
	coupons := []*models.Coupon{{ID: 1}}
	gimmickURL := "https://example.com"
	gimmickCode := "gimmick_code"
	touchPoints := []*models.TouchPoint{{GroupID: 1}}
	contentData := &models.DeliveryDataContent{
		CampaignID: campaignData.ID,
		Coupons:    []models.DeliveryCouponData{{ID: 1}},
		Gimmicks:   []models.Gimmick{{URL: gimmickURL, Code: gimmickCode}},
	}
	// DBから取得するデータの条件
	contentCondition := repository.ContentByCampaignIDCondition{CampaignID: campaignData.ID}
	updateCondition := repository.UpdateCondition{
		CampaignID: campaignData.ID,
		Status:     codes.StatusStarted,
		UpdatedAt:  deliveryData[0].UpdatedAt,
	}
	condition := repository.CampaignCondition{CampaignID: campaignData.ID, Status: codes.StatusWarmup}
	creativeCondition := repository.CreativeByCampaignIDCondition{
		CampaignID: campaignData.ID,
		Limit:      100,
	}
	// テスト用の設定
	configS := config.Env.DeliveryStart
	configUsecase := config.Env.DeliveryStartUsecase
	// workerを1にする
	configUsecase.NumberOfConcurrent = 1
	configUsecase.NumberOfQueue = 1
	t.Run("配信開始時間のキャンペーンがwarmupの場合、startedに更新してdelivery_data,creative_dataを登録して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(&repository.TouchPointByGroupIDCondition{GroupID: 1, Limit: 1})).Return(touchPoints, nil),
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(deliveryData[0].CreateDeliveryDataCampaign(cc))).Return(nil),
			touchPointDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(&models.DeliveryTouchPoint{GroupID: 1})).Return(nil),
			creativeDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(creatives[0].CreateDeliveryDataCreative())).Return(nil),
			contentDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(contentData)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEventUsecase.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(deliveryData[0].ID), gomock.Eq(deliveryData[0].OrgCode), gomock.Eq(deliveryData[0].Status),
				gomock.Eq(codes.StatusStarted), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンのstatus変更でエラーが起きた場合、エラーを返してロールバックする", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる
		dbErr := errors.New("db error")

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(0, dbErr).Times(1),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでクリエイティブ取得でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(nil, dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでギミック取得でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		errRes := ""
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&errRes, &errRes, dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでクーポン取得でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(nil, dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでタッチポイント取得でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(&repository.TouchPointByGroupIDCondition{GroupID: 1, Limit: 1})).Return(nil, dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでキャンペーン登録でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(&repository.TouchPointByGroupIDCondition{GroupID: 1, Limit: 1})).Return(touchPoints, nil),
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(deliveryData[0].CreateDeliveryDataCampaign(cc))).Return(dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでタッチポイント登録でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(&repository.TouchPointByGroupIDCondition{GroupID: 1, Limit: 1})).Return(touchPoints, nil),
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(deliveryData[0].CreateDeliveryDataCampaign(cc))).Return(nil),
			touchPointDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(&models.DeliveryTouchPoint{GroupID: 1})).Return(dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでクリエイティブ登録でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(&repository.TouchPointByGroupIDCondition{GroupID: 1, Limit: 1})).Return(touchPoints, nil),
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(deliveryData[0].CreateDeliveryDataCampaign(cc))).Return(nil),
			touchPointDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(&models.DeliveryTouchPoint{GroupID: 1})).Return(nil),
			creativeDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(creatives[0].CreateDeliveryDataCreative())).Return(dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})

	t.Run("配信開始時間のキャンペーンがwarmupでContent登録でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)
		timer := NewTimer(logger)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 deliveryDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		octx := context.Background()
		ctx, cancel := context.WithCancel(octx)
		dbErr := errors.New("db error")
		// どう呼ばれるか (呼び出し順も考慮)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(deliveryData[0], nil),
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx),
				gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil).Times(1),
			campaignRepository.EXPECT().GetCampaignCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&repository.CampaignCondition{CampaignID: campaignData.ID})).Return(cc, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondition)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(&repository.TouchPointByGroupIDCondition{GroupID: 1, Limit: 1})).Return(touchPoints, nil),
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(deliveryData[0].CreateDeliveryDataCampaign(cc))).Return(nil),
			touchPointDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(&models.DeliveryTouchPoint{GroupID: 1})).Return(nil),
			creativeDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(creatives[0].CreateDeliveryDataCreative())).Return(nil),
			contentDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(contentData)).Return(dbErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryStart := NewDeliveryStart(
			logger, metrics.GetMonitor(), &configS, &configUsecase, transactionHandler, timer,
			deliveryControlEventUsecase, campaignRepository, creativeRepository, contentRepository, touchPointRepository,
			campaignDataRepository, contentDataRepository, creativeDataRepository, touchPointDataRepository)
		// Workerを使って実行するので作成
		deliveryStart.CreateWorker(ctx)

		deliveryStart.Reserve(ctx, time.Now(), &campaignData) // 即時実行させる

		time.Sleep(100 * time.Millisecond) // 非同期で処理が実行されるので待つ
		// Workerを終了させる
		cancel()
		deliveryStart.Close()
	})
}

//nolint:unparam // `endAt` always receives `sql.NullTime{}` となっているが今後変わる可能性があるため
func createStartTestCampaign(campaignData *models.Campaign, startAt time.Time, endAt sql.NullTime, status string, updatedAt time.Time) []*models.Campaign {
	return []*models.Campaign{
		{
			ID:        campaignData.ID,
			GroupID:   campaignData.GroupID,
			Status:    status,
			StartAt:   startAt,
			EndAt:     endAt,
			UpdatedAt: updatedAt,
		},
	}
}
