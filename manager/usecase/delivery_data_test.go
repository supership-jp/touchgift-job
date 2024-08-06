package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"touchgift-job-manager/domain/models"

	mock_repository "touchgift-job-manager/mock/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// DeliveryDataのPutのテスト
func TestDeliveryData_Put(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("delivery_data4種を登録する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaign := createTestCampaign()
		creatives := []*models.Creative{
			{
				ID:  1,
				URL: "url1",
			},
		}
		content := models.DeliveryDataContent{
			CampaignID: 1,
		}
		touchPoints := []*models.DeliveryTouchPoint{
			{
				GroupID:      1,
				TouchPointID: "touchpoint1",
			},
		}

		gomock.InOrder(
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(campaign.CreateDeliveryDataCampaign())).Return(nil),
			touchPointDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(touchPoints[0])).Return(nil),
			creativeDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(creatives[0].CreateDeliveryDataCreative())).Return(nil),
			contentDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(&content)).Return(nil),
		)

		// テストを実行する
		deliveryDataUsecase := NewDeliveryData(logger, touchPointDataRepository, campaignDataRepository, contentDataRepository, creativeDataRepository)
		err := deliveryDataUsecase.Put(ctx, campaign, creatives, &content, touchPoints)
		assert.NoError(t, err)
	})

	t.Run("Creativesが空の場合もdelivery_dataを登録する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaign := createTestCampaign()
		creatives := []*models.Creative{}
		content := models.DeliveryDataContent{
			CampaignID: 1,
		}
		touchPoints := []*models.DeliveryTouchPoint{
			{
				GroupID:      1,
				TouchPointID: "touchpoint1",
			},
		}

		gomock.InOrder(
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(campaign.CreateDeliveryDataCampaign())).Return(nil),
			touchPointDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(touchPoints[0])).Return(nil),
			contentDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(&content)).Return(nil),
		)

		// テストを実行する
		deliveryDataUsecase := NewDeliveryData(logger, touchPointDataRepository, campaignDataRepository, contentDataRepository, creativeDataRepository)
		err := deliveryDataUsecase.Put(ctx, campaign, creatives, &content, touchPoints)
		assert.NoError(t, err)
	})

	t.Run("delivery_dataの登録処理でエラーが発生した場合、エラーを返して終了する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaign := createTestCampaign()
		creatives := []*models.Creative{}
		content := models.DeliveryDataContent{}
		touchPoints := []*models.DeliveryTouchPoint{}

		expectedError := errors.New("Failed to put")
		gomock.InOrder(
			campaignDataRepository.EXPECT().Put(gomock.Eq(ctx), gomock.Eq(campaign.CreateDeliveryDataCampaign())).Return(expectedError),
		)

		// テストを実行する
		deliveryDataUsecase := NewDeliveryData(logger, touchPointDataRepository, campaignDataRepository, contentDataRepository, creativeDataRepository)
		err := deliveryDataUsecase.Put(ctx, campaign, creatives, &content, touchPoints)
		assert.EqualError(t, err, expectedError.Error())
	})
}

// DeliveryDataのDeleteのテスト
func TestDeliveryData_Delete(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("CampaignIDを指定してdelivery_dataを削除できる", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaign := createTestCampaign()
		campaignID := campaign.CreateDeliveryDataCampaign().ID
		gomock.InOrder(
			campaignDataRepository.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(&campaignID)).Return(nil),
		)

		// テストを実行する
		deliveryDataUsecase := NewDeliveryData(logger, touchPointDataRepository, campaignDataRepository, contentDataRepository, creativeDataRepository)
		err := deliveryDataUsecase.Delete(ctx, campaignID)
		assert.NoError(t, err)
	})

	t.Run("delivery_dataの削除処理でエラーが発生した場合、エラーを返して終了する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		campaignDataRepository := mock_repository.NewMockDeliveryDataCampaignRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		touchPointDataRepository := mock_repository.NewMockDeliveryDataTouchPointRepository(ctrl)
		contentDataRepository := mock_repository.NewMockDeliveryDataContentRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		campaign := createTestCampaign()
		campaignID := campaign.CreateDeliveryDataCampaign().ID

		expectedError := errors.New("Failed delete delivery_data by Campaign_id.")
		gomock.InOrder(
			campaignDataRepository.EXPECT().Delete(gomock.Eq(ctx), gomock.Eq(&campaignID)).Return(expectedError),
		)

		// テストを実行する
		deliveryDataUsecase := NewDeliveryData(logger, touchPointDataRepository, campaignDataRepository, contentDataRepository, creativeDataRepository)
		err := deliveryDataUsecase.Delete(ctx, campaignID)
		assert.EqualError(t, err, expectedError.Error())
	})
}

func createTestCampaign() *models.Campaign {
	return &models.Campaign{
		ID:        1,
		Status:    "",
		StartAt:   time.Now(),
		EndAt:     sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: time.Now(),
	}
}
