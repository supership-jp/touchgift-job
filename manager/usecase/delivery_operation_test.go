package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"

	mock_repository "touchgift-job-manager/mock/repository"
	mock_usecase "touchgift-job-manager/mock/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// キャンペーンログのイベントがinsert,update
// Campaignのsyncのテスト (配信期間中)
func TestDeliveryOperation_Process_Sync_StoreDelivery_DuringDelivery(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("ステータスがconfiguredかつ配信期間中の場合何もせず終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)
		current := time.Now()

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaignLog := createTestDeliveryOperationLog("insert")
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			"configured")
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, campaignLog)
		assert.EqualError(t, err, codes.ErrDoNothing.Error())
	})

	t.Run("ステータスがresumeかつ配信期間中の場合、delivery_data, creative_dataを登録してステータスをstartedに更新する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)
		current := time.Now()

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaignLog := createTestDeliveryOperationLog("insert")
		before := "resume"
		after := codes.StatusStarted
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			before)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryStartUsecase.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(codes.StatusStarted)).Return(1, nil),
			deliveryStartUsecase.EXPECT().CreateDeliveryDatas(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&campaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEvent.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(campaign.ID), gomock.Eq(campaign.GroupID), gomock.Eq(campaign.OrgCode), gomock.Eq(campaign.Status),
				gomock.Eq(after), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, campaignLog)
		assert.NoError(t, err)
	})

	t.Run("ステータスがstartedかつ配信期間中の場合、delivery_data, creative_data, budget_dataを登録してステータスの更新はstartedのまま", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)
		current := time.Now()

		// mockの処理を定義
		// テスト対象のexecuteは、 locationDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaignLog := createTestDeliveryOperationLog("insert")
		before := codes.StatusStarted
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			before)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryStartUsecase.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(codes.StatusStarted)).Return(1, nil),
			deliveryStartUsecase.EXPECT().CreateDeliveryDatas(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&campaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEvent.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(campaign.ID), gomock.Eq(campaign.GroupID), gomock.Eq(campaign.OrgCode), gomock.Eq(campaign.Status),
				gomock.Eq(codes.StatusStarted), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, campaignLog)
		assert.NoError(t, err)
	})

	t.Run("ステータスがstartedかつ配信期間中かつキャンペーンに紐づくロケーションが無い場合、delivery_data, creative_dataを登録してステータスの更新はstartedのまま", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaignLog := createTestDeliveryOperationLog("insert")
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			codes.StatusStarted)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}
		current := time.Now()

		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryStartUsecase.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(codes.StatusStarted)).Return(1, nil),
			deliveryStartUsecase.EXPECT().CreateDeliveryDatas(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&campaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEvent.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(campaign.ID), gomock.Eq(campaign.GroupID), gomock.Eq(campaign.OrgCode), gomock.Eq(campaign.Status),
				gomock.Eq(codes.StatusStarted), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, campaignLog)
		assert.NoError(t, err)
	})

	t.Run("ステータスがstartedかつ配信期間中かつキャンペーンに紐づくクリエイティブが無い場合、delivery_dataを登録してステータスの更新はstartedのまま", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		CampaignLog := createTestDeliveryOperationLog("insert")
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			codes.StatusStarted)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}
		current := time.Now()

		// どう呼ばれるかを定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryStartUsecase.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(codes.StatusStarted)).Return(1, nil),
			deliveryStartUsecase.EXPECT().CreateDeliveryDatas(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&CampaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEvent.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(campaign.ID), gomock.Eq(campaign.GroupID), gomock.Eq(campaign.OrgCode), gomock.Eq(campaign.Status),
				gomock.Eq(codes.StatusStarted), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.NoError(t, err)
	})
}

// Campaignのsyncのテスト(配信停止系)
func TestDeliveryOperation_Process_Sync_DeleteDelivery(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("ステータスがpauseの場合、delivery_data, budget_dataを削除してステータスはpausedに更新する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)
		current := time.Now()

		// mockの処理を定義
		// テスト対象のexecuteは、 locationDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		CampaignLog := createTestDeliveryOperationLog("insert")
		before := "pause"
		after := "paused"
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			before)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryEndUsecase.EXPECT().Stop(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(after)).Return(nil),
			deliveryEndUsecase.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&CampaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEvent.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(campaign.ID), gomock.Eq(campaign.GroupID), gomock.Eq(campaign.OrgCode), gomock.Eq(campaign.Status),
				gomock.Eq(after), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.NoError(t, err)
	})

	t.Run("ステータスがstopの場合、delivery_data, budget_dataを削除してステータスはstoppedに更新する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)
		current := time.Now()

		// mockの処理を定義
		// テスト対象のexecuteは、 locationDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		CampaignLog := createTestDeliveryOperationLog("insert")
		before := "stop"
		after := "stopped"
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			before)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryEndUsecase.EXPECT().Stop(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(after)).Return(nil),
			deliveryEndUsecase.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&CampaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEvent.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(campaign.ID), gomock.Eq(campaign.GroupID), gomock.Eq(campaign.OrgCode), gomock.Eq(campaign.Status),
				gomock.Eq(after), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.NoError(t, err)
	})

}

func TestDeliveryOperation_Process_Sync_EndedDelivery(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("ステータスがendedの場合、delivery_dataを削除する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)
		current := time.Now()

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		CampaignLog := createTestDeliveryOperationLog("update")
		status := "ended"
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			status)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryEndUsecase.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&CampaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(nil),
			deliveryControlEvent.EXPECT().PublishCampaignEvent(
				gomock.Eq(ctx), gomock.Eq(campaign.ID), gomock.Eq(campaign.GroupID), gomock.Eq(campaign.OrgCode), gomock.Eq(campaign.Status),
				gomock.Eq(status), gomock.Eq(""),
			),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.NoError(t, err)
	})
}

// CampaignのProcessのテスト (キャンペーンログのイベントがdelete)
func TestDeliveryOperation_Process_Delete(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("キャンペーンログのイベントがdeleteの場合、delivery_data,delivery_budget_dataの削除はせずcreativeを処理する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		current := time.Now()
		CampaignLog := createTestDeliveryOperationLog("delete")
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&CampaignLog.Creatives)).Return(nil),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.EqualError(t, err, codes.ErrDoNothing.Error())
	})
}

// CampaignのProcessの異常系のテスト
func TestDeliveryOperation_Process_Error(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("トランザクションが開始できなかった場合、エラーを返して終了する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 locationDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		current := time.Now()
		CampaignLog := createTestDeliveryOperationLog("delete")
		expectedErr := errors.New("begin is error")
		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, expectedErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("処理の途中でエラーが発生した場合、エラーを返してRollbackする", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		current := time.Now()
		CampaignLog := createTestDeliveryOperationLog("insert")
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			codes.StatusStarted)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}
		expectedErr := errors.New("creative is error")

		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryStartUsecase.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(codes.StatusStarted)).Return(1, nil),
			deliveryStartUsecase.EXPECT().CreateDeliveryDatas(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&CampaignLog.Creatives)).Return(expectedErr),
			tx.EXPECT().Rollback().Return(nil),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("commitでエラーが発生した場合、エラーを返して終了する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		current := time.Now()
		CampaignLog := createTestDeliveryOperationLog("insert")
		before := codes.StatusStarted
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			before)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		expectedErr := errors.New("commit is error")

		gomock.InOrder(
			transactionHandler.EXPECT().Begin(gomock.Eq(ctx)).Return(tx, nil),
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryStartUsecase.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(codes.StatusStarted)).Return(1, nil),
			deliveryStartUsecase.EXPECT().CreateDeliveryDatas(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign)).Return(nil),
			creativeUsecase.EXPECT().Process(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(current), gomock.Eq(&CampaignLog.Creatives)).Return(nil),
			tx.EXPECT().Commit().Return(expectedErr),
			tx.EXPECT().Rollback().Return(nil),
		)
		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)
		err := deliveryOperationUsecase.Process(ctx, current, CampaignLog)
		assert.EqualError(t, err, expectedErr.Error())
	})
}

// DeliveryOperationのprocessCampaignLogのテスト
func TestDeliveryOperation_processCampaignLog(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	// t.Parallel() メトリクス取得の重複エラーとなって、CIが失敗するため、逐次実行にする(TODO:恒久対応検討する)

	t.Run("Campaign_logを処理する", func(t *testing.T) {
		current := time.Now()
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		CampaignLog := createTestDeliveryOperationLog("insert")
		before := codes.StatusStarted
		campaignData := models.DeliveryDataCampaign{
			ID: "1",
		}
		campaign := createSyncTestDeliveryOperation(
			&campaignData,
			time.Now().Add(-5*time.Minute),
			sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true},
			time.Now().Add(1*time.Second),
			before)
		condition := repository.CampaignCondition{
			CampaignID: campaign.ID,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(campaign, nil),
			deliveryStartUsecase.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign), gomock.Eq(codes.StatusStarted)).Return(1, nil),
			deliveryStartUsecase.EXPECT().CreateDeliveryDatas(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(campaign)).Return(nil),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)

		// private methodのテストを行うためにcastする
		deliveryOperationInteractor := deliveryOperationUsecase.(*deliveryOperation)
		_, _, _, err := deliveryOperationInteractor.processCampaignLog(ctx, tx, current, CampaignLog)
		assert.NoError(t, err)
	})

	t.Run("Campaignが取得できなかった場合ErrDoNothingエラーを返して終了", func(t *testing.T) {
		current := time.Now()
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		transactionHandler := mock_repository.NewMockTransactionHandler(ctrl)
		deliveryStartUsecase := mock_usecase.NewMockDeliveryStart(ctrl)
		creativeUsecase := mock_usecase.NewMockCreative(ctrl)
		deliveryEndUsecase := mock_usecase.NewMockDeliveryEnd(ctrl)
		deliveryControlEvent := mock_usecase.NewMockDeliveryControlEvent(ctrl)

		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		CampaignLog := createTestDeliveryOperationLog("insert")
		condition := repository.CampaignCondition{
			CampaignID: 1,
			Status:     "",
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			campaignRepository.EXPECT().GetDeliveryToStart(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(nil, nil),
		)

		// テストを実行する
		deliveryOperationUsecase := NewDeliveryOperation(logger, metrics.GetMonitor(), transactionHandler, campaignRepository, campaignDataRepository, creativeUsecase, deliveryStartUsecase, deliveryEndUsecase, deliveryControlEvent)

		// private methodのテストを行うためにcastする
		deliveryOperationInteractor := deliveryOperationUsecase.(*deliveryOperation)
		_, _, _, err := deliveryOperationInteractor.processCampaignLog(ctx, tx, current, CampaignLog)
		assert.EqualError(t, err, codes.ErrDoNothing.Error())
	})
}

func createSyncTestDeliveryOperation(campaignData *models.DeliveryDataCampaign, startAt time.Time, endAt sql.NullTime, updatedAt time.Time, status string) *models.Campaign {
	id, _ := strconv.Atoi(campaignData.ID)
	return &models.Campaign{
		OrgCode:   "org1",
		ID:        id,
		GroupID:   1,
		Status:    status,
		StartAt:   startAt,
		EndAt:     endAt,
		UpdatedAt: updatedAt,
	}
}

func createTestDeliveryOperationLog(event string) *models.CampaignLog {
	return &models.CampaignLog{
		ID:    1,
		Event: event,
		Creatives: []models.CreativeLog{
			{
				ID:    1,
				Event: event,
			},
		},
	}
}
