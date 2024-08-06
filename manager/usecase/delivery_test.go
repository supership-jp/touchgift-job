package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"
	mock_repository "touchgift-job-manager/mock/repository"
	mock_usecase "touchgift-job-manager/mock/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryUpdateStatus(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("正常ケース(引数にて処理は変わらない)", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		status := "started"
		updatedAt := time.Now()
		condition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     status,
			UpdatedAt:  updatedAt,
		}
		// 更新後のcampaign.updated_at
		// expectedUpdatedAt := updatedAt.Add(1 * time.Minute)

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(1, nil),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		count, err := deliveryUsecase.UpdateStatus(ctx, tx, campaignID, status, updatedAt)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, count)
		}
	})
}

func TestDeliveryStartOrSync(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("正常に登録する(タッチポイント、コンテントあり)", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		groupID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		// テストデータの作成 クリエイティブ
		creativeCondition := repository.CreativeByCampaignIDCondition{
			CampaignID: campaignID,
		}
		creatives := []*models.Creative{
			{ID: 1, URL: ""},
		}
		// キャンペーン
		campaignData := models.Campaign{
			ID:        campaignID,
			GroupID:   1,
			StartAt:   updatedAt,
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)
		// コンテント
		contentCondtion := repository.ContentByCampaignIDCondition{
			CampaignID: campaignID,
		}
		gimmickURL := "https://example.com"
		gimmickCode := "gimmick_code"
		coupons := []*models.Coupon{
			{
				ID: 1,
			},
		}
		content := createDeliveryDataContent(campaignID, coupons[0].ID, gimmickURL, gimmickCode)
		// タッチポイント
		touchPointCondition := &repository.TouchPointByGroupIDCondition{
			GroupID: 1,
			Limit:   1,
		}
		touchPoints := []*models.TouchPoint{
			{
				GroupID: groupID,
			},
		}
		touchPointDatas := []*models.DeliveryTouchPoint{
			{
				GroupID: groupID,
			},
		}

		// status更新後のcampaign.updated_at
		// expectedUpdatedAt := updatedAt.Add(1 * time.Minute)

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondtion)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondtion)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(touchPointCondition)).Return(touchPoints, nil),
			deliveryDataUsecase.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(campaign), gomock.Eq(creatives), gomock.Eq(content), gomock.Eq(touchPointDatas)).Return(nil),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.NoError(t, err)
	})

	t.Run("正常に登録する(セグメントあり)", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		groupID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		// テストデータの作成 クリエイティブ
		creativeCondition := repository.CreativeByCampaignIDCondition{
			CampaignID: campaignID,
		}
		creatives := []*models.Creative{
			{ID: 1, URL: ""},
		}
		// キャンペーン
		campaignData := models.Campaign{
			ID:        campaignID,
			GroupID:   1,
			StartAt:   updatedAt,
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)
		// コンテント
		contentCondtion := repository.ContentByCampaignIDCondition{
			CampaignID: campaignID,
		}
		gimmickURL := "https://example.com"
		gimmickCode := "gimmick_code"
		coupons := []*models.Coupon{
			{
				ID: 1,
			},
		}
		content := createDeliveryDataContent(campaignID, coupons[0].ID, gimmickURL, gimmickCode)
		// タッチポイント
		touchPointCondition := &repository.TouchPointByGroupIDCondition{
			GroupID: 1,
			Limit:   1,
		}
		touchPoints := []*models.TouchPoint{
			{
				GroupID: groupID,
			},
		}
		touchPointDatas := []*models.DeliveryTouchPoint{
			{
				GroupID: groupID,
			},
		}

		// status更新後のcampaign.updated_at
		// expectedUpdatedAt := updatedAt.Add(1 * time.Minute)

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&creativeCondition)).Return(creatives, nil),
			contentRepository.EXPECT().GetGimmicksByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondtion)).Return(&gimmickURL, &gimmickCode, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&contentCondtion)).Return(coupons, nil),
			touchPointRepository.EXPECT().GetTouchPointByGroupID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(touchPointCondition)).Return(touchPoints, nil),
			deliveryDataUsecase.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(campaign), gomock.Eq(creatives), gomock.Eq(content), gomock.Eq(touchPointDatas)).Return(nil),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.NoError(t, err)
	})

	t.Run("campaignに紐づくlocationが無い場合配信データを登録する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		condition := repository.CampaignCondition{
			CampaignID: campaignID,
		}
		creatives := []models.Creative{
			{ID: 1, URL: ""},
		}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)

		// status更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, expectedUpdatedAt, nil),
			//creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(segments, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(creatives, nil),
			//deliveryDataUsecase.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(campaign), gomock.Eq(&creatives), gomock.Eq(&segments), gomock.Eq(expectedUpdatedAt)).Return(nil),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.NoError(t, err)
	})

	t.Run("creativeが無い場合配信データを登録する(creative_dataは登録しない)", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		condition := repository.CampaignCondition{
			CampaignID: campaignID,
		}
		creatives := []models.Creative{}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)

		// status更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, expectedUpdatedAt, nil),
			//creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(segments, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(creatives, nil),
			//deliveryDataUsecase.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(campaign), gomock.Eq(&creatives), gomock.Eq(&segments), gomock.Eq(expectedUpdatedAt)).Return(nil),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.NoError(t, err)
	})

	t.Run("campaignに紐づくSegmentを取得する際にエラーが発生した場合、処理を中断", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		condition := repository.CampaignCondition{
			CampaignID: campaignID,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)

		// status更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)
		expected := errors.New("Failed to get")
		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, expectedUpdatedAt, nil),
			creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(nil, expected),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.EqualError(t, err, expected.Error())
	})

	t.Run("campaignに紐づくCreativeを取得する際にエラーが発生した場合、処理を中断", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		condition := repository.CampaignCondition{
			CampaignID: campaignID,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)

		// status更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)
		expected := errors.New("Failed to get")
		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, expectedUpdatedAt, nil),
			//creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(segments, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(nil, expected),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.EqualError(t, err, expected.Error())
	})

	t.Run("配信状態の更新に失敗した場合DynamoDBに登録しない(Publishしない)", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()
		condition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)
		// 更新後のcampaign.updated_at
		expectedUpdatedAt := time.Time{}

		expected := errors.New("Failed to update")
		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(0, expectedUpdatedAt, expected),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.EqualError(t, err, expected.Error())
	})

	t.Run("CreativeのDynamoDB登録時にエラーが発生した場合、処理を中断", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()

		condition := repository.CampaignCondition{
			CampaignID: campaignID,
		}
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		creatives := []models.Creative{
			{ID: 1, URL: ""},
		}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)
		// 更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)
		expected := errors.New("Failed to put creative")

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, expectedUpdatedAt, nil),
			//creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(segments, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(creatives, nil),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.EqualError(t, err, expected.Error())
	})

	t.Run("配信データのDynamoDB登録時にエラーが発生した場合、処理を中断", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "warmup"
		after := "started"
		updatedAt := time.Now()

		condition := repository.CampaignCondition{
			CampaignID: campaignID,
		}
		updateCondition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		creatives := []models.Creative{
			{ID: 1, URL: ""},
		}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)

		// 更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)
		expected := errors.New("Failed to put delivery_data")

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&updateCondition)).Return(1, expectedUpdatedAt, nil),
			//creativeRepository.EXPECT().GetCreativeByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(segments, nil),
			contentRepository.EXPECT().GetCouponsByCampaignID(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(creatives, nil),
			//deliveryDataUsecase.EXPECT().Put(gomock.Eq(ctx),
			//gomock.Eq(campaign), gomock.Eq(&creatives), gomock.Eq(&segments),
			//gomock.Eq(expectedUpdatedAt)).Return(expected),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.StartOrSync(ctx, tx, campaign)
		assert.EqualError(t, err, expected.Error())
	})
}

func TestDeliveryStop(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("正常に停止する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "terminate"
		after := "ended"
		updatedAt := time.Now()
		condition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)
		// 更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)

		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(1, expectedUpdatedAt, nil),
			deliveryDataUsecase.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(campaignID)).Return(nil),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.Stop(ctx, tx, campaign, after)
		assert.NoError(t, err)
	})

	t.Run("配信状態の更新に失敗した場合DynamoDBから削除しない(Publishしない)", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "terminate"
		after := "ended"
		updatedAt := time.Now()
		condition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)
		// 更新後のcampaign.updated_at
		expectedUpdatedAt := time.Time{}

		expected := errors.New("Failed to update")
		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(0, expectedUpdatedAt, expected),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.Stop(ctx, tx, campaign, after)
		assert.EqualError(t, err, expected.Error())
	})

	t.Run("DynamoDBのデータ削除でエラーが発生した場合、Publishしない", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignRepository := mock_repository.NewMockCampaignRepository(ctrl)
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		contentRepository := mock_repository.NewMockContentRepository(ctrl)
		touchPointRepository := mock_repository.NewMockTouchPointRepository(ctrl)
		deliveryDataUsecase := mock_usecase.NewMockDeliveryData(ctrl)
		deliveryControlEventUsecase := mock_usecase.NewMockDeliveryControlEvent(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		ctx := context.Background()

		campaignID := 1
		before := "terminate"
		after := "ended"
		updatedAt := time.Now()
		condition := repository.UpdateCondition{
			CampaignID: campaignID,
			Status:     after,
			UpdatedAt:  updatedAt,
		}
		campaignData := models.Campaign{
			ID:        campaignID,
			StartAt:   time.Now(),
			UpdatedAt: updatedAt,
		}
		campaign := createSyncTestDeliveryOperation(
			campaignData.CreateDeliveryDataCampaign(),
			campaignData.StartAt,
			sql.NullTime{Time: campaignData.StartAt.Add(20 * time.Minute), Valid: true},
			campaignData.UpdatedAt,
			before)
		// 更新後のcampaign.updated_at
		expectedUpdatedAt := updatedAt.Add(1 * time.Minute)

		expected := errors.New("Failed to delete")
		gomock.InOrder(
			campaignRepository.EXPECT().UpdateStatus(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(1, expectedUpdatedAt, nil),
			deliveryDataUsecase.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(campaignID)).Return(expected),
		)

		deliveryUsecase := NewDelivery(logger, metrics.GetMonitor(), campaignRepository, creativeRepository, contentRepository, touchPointRepository, deliveryDataUsecase, deliveryControlEventUsecase)
		err := deliveryUsecase.Stop(ctx, tx, campaign, after)
		assert.EqualError(t, err, expected.Error())
	})
}

func createDeliveryDataContent(campaignID int, couponID int, gimmickURL string, gimmickCode string) *models.DeliveryDataContent {
	return &models.DeliveryDataContent{
		CampaignID: campaignID,
		Coupons: []models.DeliveryCouponData{
			{
				ID: couponID,
			},
		},
		Gimmicks: []models.Gimmick{
			{
				URL:  gimmickURL,
				Code: gimmickCode,
			},
		},
	}
}
