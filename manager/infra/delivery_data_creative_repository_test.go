package infra

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/infra/metrics"

	"github.com/stretchr/testify/assert"
)

// CreativeDataRepository の Get のテスト
func TestCreativeDataRepository_Get(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	t.Run("creative_dataが空の場合はエラーを返す", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		ID := "100"
		actual, err := creativeDataRepository.Get(ctx, &ID)
		if assert.Error(t, err) {
			assert.Nil(t, actual)
			assert.EqualError(t, err, codes.ErrNoData.Error())
		}
	})

	t.Run("creative_dataを1件返す", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		ID := 1
		IDString := strconv.Itoa(ID)
		expected := models.DeliveryDataCreative{
			CampaignID: ID,
			URL:        "id_get1_url",
			TTL:        time.Now().Unix(),
		}
		// データを用意
		if err := creativeDataRepository.Put(ctx, &expected); !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := creativeDataRepository.Delete(ctx, &IDString); err != nil {
				assert.NoError(t, err)
			}
		}()
		// テスト実行する
		actual, err := creativeDataRepository.Get(ctx, &IDString)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected, *actual)
		}
	})
}

// CreativeDataRepository の Put のテスト
func TestCreativeDataRepository_Put(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)
	campaignID := 1
	IDString := strconv.Itoa(campaignID)

	createData := func() *models.DeliveryDataCreative {
		return &models.DeliveryDataCreative{
			CampaignID: campaignID,
			URL:        "id_put1_url",
			TTL:        time.Now().Unix(),
		}
	}

	t.Run("creative_dataを1件登録", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := creativeDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := creativeDataRepository.Delete(ctx, &IDString); err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := creativeDataRepository.Get(ctx, &IDString)
		if assert.NoError(t, err) {
			assert.Exactly(t, *expected, *actual)
		}
	})

	t.Run("creative_dataを1件更新", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := creativeDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 更新用データ用意
		expected2 := *expected
		expected2.URL = "id_put1_url2"
		// 更新する
		err = creativeDataRepository.Put(ctx, &expected2)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := creativeDataRepository.Delete(ctx, &IDString); err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := creativeDataRepository.Get(ctx, &IDString)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected2, *actual)
		}
	})
}

// CreativeDataRepository の PutAll のテスト
func TestCreativeDataRepository_PutAll(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	createData := func() *[]models.DeliveryDataCreative {
		return &[]models.DeliveryDataCreative{
			{
				CampaignID: 1,
				URL:        "id_put1_url",
				TTL:        time.Now().Unix(),
			},
			{
				CampaignID: 2,
				URL:        "id_put2_url",
				TTL:        time.Now().Unix(),
			},
		}
	}

	t.Run("creative_dataを複数件登録", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := creativeDataRepository.PutAll(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			for i := range *expected {
				data := (*expected)[i]
				IDString := strconv.Itoa(data.CampaignID)
				if err := creativeDataRepository.Delete(ctx, &IDString); err != nil {
					assert.NoError(t, err)
				}
			}
		}()
		for i := range *expected {
			data := (*expected)[i]
			IDString := strconv.Itoa(data.CampaignID)
			actual, err := creativeDataRepository.Get(ctx, &IDString)
			if assert.NoError(t, err) {
				assert.Exactly(t, data, *actual)
			}
		}
	})

	t.Run("creative_dataを複数件更新", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		initData := createData()
		// 登録する
		err := creativeDataRepository.PutAll(ctx, initData)
		if !assert.NoError(t, err) {
			return
		}
		// 更新用データ整理
		updateData := make([]models.DeliveryDataCreative, len(*initData))
		for i := range *initData {
			data := (*initData)[i]
			data.URL = fmt.Sprintf("update_id%d_put_url", i)
			updateData[i] = data
		}
		// 更新する
		err = creativeDataRepository.PutAll(ctx, &updateData)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			for i := range updateData {
				data := (updateData)[i]
				IDString := strconv.Itoa(data.CampaignID)
				if err := creativeDataRepository.Delete(ctx, &IDString); err != nil {
					assert.NoError(t, err)
				}
			}
		}()
		for i := range updateData {
			data := (updateData)[i]
			IDString := strconv.Itoa(data.CampaignID)
			actual, err := creativeDataRepository.Get(ctx, &IDString)
			if assert.NoError(t, err) {
				assert.Exactly(t, data, *actual)
			}
		}
	})
}

// CreativeDataRepository の Delete のテスト
func TestCreativeDataRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	t.Run("creative_dataを1件削除", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		ID := 1
		expected := models.DeliveryDataCreative{
			CampaignID: ID,
			URL:        "id_delete1_url",
			TTL:        time.Now().Unix(),
		}
		if err := creativeDataRepository.Put(ctx, &expected); !assert.NoError(t, err) {
			return
		}
		IDString := strconv.Itoa(ID)
		err := creativeDataRepository.Delete(ctx, &IDString)
		if assert.NoError(t, err) {
			actual, err := creativeDataRepository.Get(ctx, &IDString)
			if assert.Error(t, err) {
				assert.Nil(t, actual)
				assert.EqualError(t, err, codes.ErrNoData.Error())
			}
		}
	})
	t.Run("creative_dataの削除は対象がない場合エラーは返さない", func(t *testing.T) {
		creativeDataRepository := NewDeliveryDataCreativeRepository(dynamodbHandler, logger, monitor)
		ID := "100"
		err := creativeDataRepository.Delete(ctx, &ID)
		assert.NoError(t, err)
	})
}
