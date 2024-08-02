package infra

import (
	"context"
	"testing"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra/metrics"

	"github.com/stretchr/testify/assert"
)

func TestCampaignDataRepository_Get(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	t.Run("campaign_dataが空の場合はエラーを返す", func(t *testing.T) {
		campaignDataRepository := NewCampaignDataRepository(dynamodbHandler, logger, monitor)
		ID := "none"
		actual, err := campaignDataRepository.Get(ctx, &ID)
		if assert.Error(t, err) {
			assert.Nil(t, actual)
			assert.EqualError(t, err, codes.ErrNoData.Error())
		}
	})
	t.Run("campaign_dataを1件返す", func(t *testing.T) {
		campaignDataRepository := NewCampaignDataRepository(dynamodbHandler, logger, monitor)
		ID := "id_get1"
		expected := models.DeliveryDataCampaign{
			ID:      ID,
			GroupID: 1,
			OrgCode: "ORG1",
			Status:  "warmup",
		}
		// データを用意
		if err := campaignDataRepository.Put(ctx, &expected); !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := campaignDataRepository.Delete(ctx, &ID); err != nil {
				assert.NoError(t, err)
			}
		}()
		// テスト実行する
		actual, err := campaignDataRepository.Get(ctx, &ID)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected, *actual)
		}
	})

}
func TestCampaignDataRepository_Put(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.NewMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	createData := func() *models.DeliveryDataCampaign {
		return &models.DeliveryDataCampaign{
			ID:      "1",
			GroupID: 1,
			OrgCode: "ORG1",
			Status:  "warmup",
		}
	}

	t.Run("campaign_dataを1件登録できる", func(t *testing.T) {
		campaignDataRepository := NewCampaignDataRepository(dynamodbHandler, logger, monitor)
		expected := createData()

		// 登録する
		err := campaignDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}

		defer func() {
			err := campaignDataRepository.Delete(ctx, &expected.ID)
			if err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := campaignDataRepository.Get(ctx, &expected.ID)
		if assert.NoError(t, err) {
			assert.Exactly(t, *expected, *actual)
		}
	})
	t.Run("campaign_dataを1更新できる", func(t *testing.T) {
		campaignDataRepository := NewCampaignDataRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := campaignDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 更新用データ用意
		expected2 := *expected
		// 更新する
		err = campaignDataRepository.Put(ctx, &expected2)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := campaignDataRepository.Delete(ctx, &expected2.ID); err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := campaignDataRepository.Get(ctx, &expected2.ID)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected2, *actual)
		}
	})
}
func TestCampaignDataRepository_PutAll(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.NewMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	createData := func(id string, name string) *models.DeliveryDataCampaign {
		return &models.DeliveryDataCampaign{
			ID:      id,
			GroupID: 1,
			OrgCode: "ORG1",
			Status:  "active",
		}
	}

	initialCampaigns := []models.DeliveryDataCampaign{
		*createData("1", "キャンペーン初期A"),
		*createData("2", "キャンペーン初期B"),
		*createData("3", "キャンペーン初期C"),
	}

	updatedCampaigns := []models.DeliveryDataCampaign{
		*createData("1", "キャンペーン更新A"),
		*createData("2", "キャンペーン更新B"),
		*createData("3", "キャンペーン更新C"),
	}

	campaignDataRepository := NewCampaignDataRepository(dynamodbHandler, logger, monitor)

	t.Run("複数のcampaign_dataを登録できる", func(t *testing.T) {
		err := campaignDataRepository.PutAll(ctx, &initialCampaigns)
		if !assert.NoError(t, err) {
			return
		}

		defer func() {
			for _, campaign := range initialCampaigns {
				if err := campaignDataRepository.Delete(ctx, &campaign.ID); err != nil {
					assert.NoError(t, err)
				}
			}
		}()

		// 登録した各キャンペーンを検証
		for _, expected := range initialCampaigns {
			actual, err := campaignDataRepository.Get(ctx, &expected.ID)
			if assert.NoError(t, err) {
				assert.Exactly(t, expected, *actual)
			}
		}
	})

	t.Run("複数のcampaign_dataを更新できる", func(t *testing.T) {
		// データの更新を行う
		err := campaignDataRepository.PutAll(ctx, &updatedCampaigns)
		if !assert.NoError(t, err) {
			return
		}

		// 更新したデータを検証
		for _, expected := range updatedCampaigns {
			actual, err := campaignDataRepository.Get(ctx, &expected.ID)
			if assert.NoError(t, err) {
				assert.Exactly(t, expected, *actual)
			}
		}
	})
}
func TestCampaignDataRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.NewMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	createData := func(id string, name string) *models.DeliveryDataCampaign {
		return &models.DeliveryDataCampaign{
			ID:      id,
			GroupID: 1,
			OrgCode: "ORG1",
			Status:  "active",
		}
	}

	testCampaign := createData("100", "テストキャンペーン")

	campaignDataRepository := NewCampaignDataRepository(dynamodbHandler, logger, monitor)
	// Ensure the campaign is created before we try to delete it
	_ = campaignDataRepository.Put(ctx, testCampaign)

	t.Run("campaign_dataを1件削除できる", func(t *testing.T) {
		err := campaignDataRepository.Delete(ctx, &testCampaign.ID)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		// Confirm deletion
		_, err = campaignDataRepository.Get(ctx, &testCampaign.ID)
		assert.Error(t, err) // Expect an error as the data should no longer exist
	})

	t.Run("複数のcampaign_dataを削除できる", func(t *testing.T) {
		// Setup - Creating multiple campaigns to delete
		multipleCampaigns := []models.DeliveryDataCampaign{
			*createData("101", "キャンペーンX"),
			*createData("102", "キャンペーンY"),
			*createData("103", "キャンペーンZ"),
		}
		for _, cmp := range multipleCampaigns {
			_ = campaignDataRepository.Put(ctx, &cmp)
		}

		err := campaignDataRepository.DeleteAll(ctx, &multipleCampaigns)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		// Confirm all deletions
		for _, cmp := range multipleCampaigns {
			_, err := campaignDataRepository.Get(ctx, &cmp.ID)
			assert.Error(t, err) // Expect an error as the data should no longer exist
		}
	})
}
