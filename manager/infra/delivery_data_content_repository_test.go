package infra

import (
	"context"
	"strconv"
	"testing"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra/metrics"

	"github.com/stretchr/testify/assert"
)

// ContentDataRepository の Get のテスト
func TestContentDataRepository_Get(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	t.Run("content_dataが空の場合はエラーを返す", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		ID := "100"
		actual, err := contentDataRepository.Get(ctx, &ID)
		if assert.Error(t, err) {
			assert.Nil(t, actual)
			assert.EqualError(t, err, codes.ErrNoData.Error())
		}
	})

	t.Run("content_dataを1件返す", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		ID := "1"
		expected := models.DeliveryDataContent{
			CampaignID: ID,
			Coupons:    []models.DeliveryCouponData{{ID: 1}},
			Gimmicks:   []models.Gimmick{{URL: "URL1"}},
		}
		// データを用意
		if err := contentDataRepository.Put(ctx, &expected); !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := contentDataRepository.Delete(ctx, &ID); err != nil {
				assert.NoError(t, err)
			}
		}()
		// テスト実行する
		actual, err := contentDataRepository.Get(ctx, &ID)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected, *actual)
		}
	})
}

// ContentDataRepository の Put のテスト
func TestContentDataRepository_Put(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)
	campaignID := "1"

	createData := func() *models.DeliveryDataContent {
		return &models.DeliveryDataContent{
			CampaignID: campaignID,
			Coupons:    []models.DeliveryCouponData{{ID: 1}},
			Gimmicks:   []models.Gimmick{{URL: "URL1"}},
		}
	}

	t.Run("content_dataを1件登録", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := contentDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := contentDataRepository.Delete(ctx, &campaignID); err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := contentDataRepository.Get(ctx, &campaignID)
		if assert.NoError(t, err) {
			assert.Exactly(t, *expected, *actual)
		}
	})

	t.Run("content_dataを1件更新", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := contentDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 更新用データ用意
		expected2 := *expected
		expected2.Coupons = []models.DeliveryCouponData{{ID: 2}}
		// 更新する
		err = contentDataRepository.Put(ctx, &expected2)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := contentDataRepository.Delete(ctx, &campaignID); err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := contentDataRepository.Get(ctx, &campaignID)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected2, *actual)
		}
	})
}

// ContentDataRepository の PutAll のテスト
func TestContentDataRepository_PutAll(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	createData := func() *[]models.DeliveryDataContent {
		return &[]models.DeliveryDataContent{
			{
				CampaignID: "1",
				Coupons:    []models.DeliveryCouponData{{ID: 1}},
				Gimmicks:   []models.Gimmick{{URL: "URL1"}},
			},
			{
				CampaignID: "2",
				Coupons:    []models.DeliveryCouponData{{ID: 2}},
				Gimmicks:   []models.Gimmick{{URL: "URL2"}},
			},
		}
	}

	t.Run("content_dataを複数件登録", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := contentDataRepository.PutAll(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			for i := range *expected {
				data := (*expected)[i]
				if err := contentDataRepository.Delete(ctx, &data.CampaignID); err != nil {
					assert.NoError(t, err)
				}
			}
		}()
		for i := range *expected {
			data := (*expected)[i]
			actual, err := contentDataRepository.Get(ctx, &data.CampaignID)
			if assert.NoError(t, err) {
				assert.Exactly(t, data, *actual)
			}
		}
	})

	t.Run("content_dataを複数件更新", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		initData := createData()
		// 登録する
		err := contentDataRepository.PutAll(ctx, initData)
		if !assert.NoError(t, err) {
			return
		}
		// 更新用データ整理
		updateData := make([]models.DeliveryDataContent, len(*initData))
		for i := range *initData {
			data := (*initData)[i]
			data.CampaignID = strconv.Itoa(i)
			updateData[i] = data
		}
		// 更新する
		err = contentDataRepository.PutAll(ctx, &updateData)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			for i := range updateData {
				data := (updateData)[i]
				if err := contentDataRepository.Delete(ctx, &data.CampaignID); err != nil {
					assert.NoError(t, err)
				}
			}
		}()
		for i := range updateData {
			data := (updateData)[i]
			actual, err := contentDataRepository.Get(ctx, &data.CampaignID)
			if assert.NoError(t, err) {
				assert.Exactly(t, data, *actual)
			}
		}
	})
}

// ContentDataRepository の Delete のテスト
func TestContentDataRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	t.Run("content_dataを1件削除", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		ID := "1"
		expected := models.DeliveryDataContent{
			CampaignID: ID,
			Coupons:    []models.DeliveryCouponData{{ID: 1}},
			Gimmicks:   []models.Gimmick{{URL: "URL1"}},
		}
		if err := contentDataRepository.Put(ctx, &expected); !assert.NoError(t, err) {
			return
		}
		err := contentDataRepository.Delete(ctx, &ID)
		if assert.NoError(t, err) {
			actual, err := contentDataRepository.Get(ctx, &ID)
			if assert.Error(t, err) {
				assert.Nil(t, actual)
				assert.EqualError(t, err, codes.ErrNoData.Error())
			}
		}
	})
	t.Run("content_dataの削除は対象がない場合エラーは返さない", func(t *testing.T) {
		contentDataRepository := NewDeliveryDataContentRepository(dynamodbHandler, logger, monitor)
		ID := "100"
		err := contentDataRepository.Delete(ctx, &ID)
		assert.NoError(t, err)
	})
}
