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

// TouchPointDataRepository の Get のテスト
func TestTouchPointDataRepository_Get(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	t.Run("touchpoint_dataが空の場合はエラーを返す", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		ID := "100"
		actual, err := touchPointDataRepository.Get(ctx, &ID)
		if assert.Error(t, err) {
			assert.Nil(t, actual)
			assert.EqualError(t, err, codes.ErrNoData.Error())
		}
	})

	t.Run("touchpoint_dataを1件返す", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		ID := 1
		IDString := strconv.Itoa(ID)
		expected := models.DeliveryTouchPoint{
			GroupID:      ID,
			TouchPointID: "test1",
		}
		// データを用意
		if err := touchPointDataRepository.Put(ctx, &expected); !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := touchPointDataRepository.Delete(ctx, &IDString); err != nil {
				assert.NoError(t, err)
			}
		}()
		// テスト実行する
		actual, err := touchPointDataRepository.Get(ctx, &IDString)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected, *actual)
		}
	})
}

// TouchPointDataRepository の Put のテスト
func TestTouchPointDataRepository_Put(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)
	id := 1
	IDString := strconv.Itoa(id)

	createData := func() *models.DeliveryTouchPoint {
		return &models.DeliveryTouchPoint{
			GroupID:      id,
			TouchPointID: "test1",
		}
	}

	t.Run("touchpoint_dataを1件登録", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := touchPointDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := touchPointDataRepository.Delete(ctx, &IDString); err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := touchPointDataRepository.Get(ctx, &IDString)
		if assert.NoError(t, err) {
			assert.Exactly(t, *expected, *actual)
		}
	})

	t.Run("touchpoint_dataを1件更新", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := touchPointDataRepository.Put(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 更新用データ用意
		expected2 := *expected
		expected2.TouchPointID = "test2"
		// 更新する
		err = touchPointDataRepository.Put(ctx, &expected2)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			if err := touchPointDataRepository.Delete(ctx, &IDString); err != nil {
				assert.NoError(t, err)
			}
		}()
		actual, err := touchPointDataRepository.Get(ctx, &IDString)
		if assert.NoError(t, err) {
			assert.Exactly(t, expected2, *actual)
		}
	})
}

// TouchPointDataRepository の PutAll のテスト
func TestTouchPointDataRepository_PutAll(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	createData := func() *[]models.DeliveryTouchPoint {
		return &[]models.DeliveryTouchPoint{
			{
				GroupID:      1,
				TouchPointID: "test1",
			},
			{
				GroupID:      2,
				TouchPointID: "test2",
			},
		}
	}

	t.Run("touchpoint_dataを複数件登録", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		expected := createData()
		// 登録する
		err := touchPointDataRepository.PutAll(ctx, expected)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			for i := range *expected {
				data := (*expected)[i]
				IDString := strconv.Itoa(data.GroupID)
				if err := touchPointDataRepository.Delete(ctx, &IDString); err != nil {
					assert.NoError(t, err)
				}
			}
		}()
		for i := range *expected {
			data := (*expected)[i]
			IDString := strconv.Itoa(data.GroupID)
			actual, err := touchPointDataRepository.Get(ctx, &IDString)
			if assert.NoError(t, err) {
				assert.Exactly(t, data, *actual)
			}
		}
	})

	t.Run("touchPoint_dataを複数件更新", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		initData := createData()
		// 登録する
		err := touchPointDataRepository.PutAll(ctx, initData)
		if !assert.NoError(t, err) {
			return
		}
		// 更新用データ整理
		updateData := make([]models.DeliveryTouchPoint, len(*initData))
		for i := range *initData {
			data := (*initData)[i]
			data.GroupID = i
			updateData[i] = data
		}
		// 更新する
		err = touchPointDataRepository.PutAll(ctx, &updateData)
		if !assert.NoError(t, err) {
			return
		}
		// 用意したデータを削除
		defer func() {
			for i := range updateData {
				data := (updateData)[i]
				IDString := strconv.Itoa(data.GroupID)
				if err := touchPointDataRepository.Delete(ctx, &IDString); err != nil {
					assert.NoError(t, err)
				}
			}
		}()
		for i := range updateData {
			data := (updateData)[i]
			IDString := strconv.Itoa(data.GroupID)
			actual, err := touchPointDataRepository.Get(ctx, &IDString)
			if assert.NoError(t, err) {
				assert.Exactly(t, data, *actual)
			}
		}
	})
}

// TouchPointDataRepository の Delete のテスト
func TestTouchPointDataRepository_Delete(t *testing.T) {
	ctx := context.Background()
	logger := GetLogger()
	monitor := metrics.GetMonitor()
	region := NewRegion(logger)
	dynamodbHandler := NewDynamoDBHandler(logger, region)

	t.Run("touchpoint_dataを1件削除", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		ID := 1
		expected := models.DeliveryTouchPoint{
			GroupID:      ID,
			TouchPointID: "test1",
		}
		if err := touchPointDataRepository.Put(ctx, &expected); !assert.NoError(t, err) {
			return
		}
		IDString := strconv.Itoa(ID)
		err := touchPointDataRepository.Delete(ctx, &IDString)
		if assert.NoError(t, err) {
			actual, err := touchPointDataRepository.Get(ctx, &IDString)
			if assert.Error(t, err) {
				assert.Nil(t, actual)
				assert.EqualError(t, err, codes.ErrNoData.Error())
			}
		}
	})
	t.Run("touchpoint_dataの削除は対象がない場合エラーは返さない", func(t *testing.T) {
		touchPointDataRepository := NewDeliveryDataTouchPointRepository(dynamodbHandler, logger, monitor)
		ID := "100"
		err := touchPointDataRepository.Delete(ctx, &ID)
		assert.NoError(t, err)
	})
}
