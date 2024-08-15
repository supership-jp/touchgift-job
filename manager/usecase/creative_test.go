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

	mock_repository "touchgift-job-manager/mock/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// CreativeのProcessのテスト (event: deleteの場合)
func TestCreative_Process_Delete(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("クリエイティブログのイベントがdeleteかつRDBから対象のキャンペーンIDに紐づくcreativeが取得できた場合、何もせず終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 locationDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		creativeLogs := []models.CreativeLog{
			{ID: 1, Event: "delete"},
		}
		condition := repository.CreativeCondition{
			ID: creativeLogs[0].ID,
		}
		creatives := createTestCreatives(
			sql.NullTime{Time: time.Now(), Valid: true},
		)
		creativeDatas := make([]models.DeliveryDataCreative, len(creatives))
		for i := range creatives {
			creativeDatas[i] = *(creatives)[i].CreateDeliveryDataCreative()
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			creativeRepository.EXPECT().GetCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(creatives, nil),
		)

		// テストを実行する
		creative := NewCreative(logger, creativeDataRepository, creativeRepository)
		err := creative.Process(ctx, tx, time.Now(), &creativeLogs)
		assert.NoError(t, err)
	})

	t.Run("クリエイティブログのイベントがdeleteかつRDBから対象のキャンペーンIDに紐づくcreativeが取得できなかった場合、creativeのttlを1日後にして更新", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 locationDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		creativeLogs := []models.CreativeLog{
			{ID: 1, Event: "delete"},
		}
		creativeID := strconv.Itoa(creativeLogs[0].ID)
		creatives := []models.Creative{}
		condition := repository.CreativeCondition{
			ID: creativeLogs[0].ID,
		}

		// 何回呼ばれるか (Times)
		// を定義する
		gomock.InOrder(
			creativeRepository.EXPECT().GetCreative(gomock.Eq(ctx), gomock.Eq(tx), gomock.Eq(&condition)).Return(creatives, nil),
			creativeDataRepository.EXPECT().UpdateTTL(gomock.Eq(ctx), gomock.Eq(creativeID), gomock.Eq(time.Now().Add(24*time.Hour).Truncate(time.Millisecond).Unix())).Return(nil),
		)
		// テストを実行する
		creative := NewCreative(logger, creativeDataRepository, creativeRepository)
		err := creative.Process(ctx, tx, time.Now(), &creativeLogs)
		assert.NoError(t, err)
	})
}

// CreativeのProcessのテスト (イベントがdelete以外)
func TestCreative_Process_Etc(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("クリエイティブログのイベントがdelete以外の場合、campaign処理で行うため何もしない", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)
		tx := mock_repository.NewMockTransaction(ctrl)

		// mockの処理を定義
		// テスト対象のexecuteは、 locationDataRepository.Put を使っているのでその処理を定義する
		// 引数に渡ると想定される値
		ctx := context.Background()
		creativeLogs := []models.CreativeLog{
			{ID: 1, Event: "insert"},
		}

		// テストを実行する
		creative := NewCreative(logger, creativeDataRepository, creativeRepository)
		err := creative.Process(ctx, tx, time.Now(), &creativeLogs)
		assert.NoError(t, err)
	})
}

// CreativeのPutのテスト
func TestCreative_Put(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("creativeを登録する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		creatives := createTestCreatives(sql.NullTime{})
		creativeDatas := make([]models.DeliveryDataCreative, len(creatives))
		for i := range creatives {
			creativeDatas[i] = *(creatives)[i].CreateDeliveryDataCreative()
		}
		gomock.InOrder(
			creativeDataRepository.EXPECT().PutAll(gomock.Eq(ctx), gomock.Eq(&creativeDatas)).Return(nil),
		)

		// テストを実行する
		creative := NewCreative(logger, creativeDataRepository, creativeRepository)
		err := creative.Put(ctx, &creativeDatas)
		assert.NoError(t, err)
	})

	t.Run("creative登録処理でエラーが起きた場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		creatives := createTestCreatives(sql.NullTime{})
		creativeDatas := make([]models.DeliveryDataCreative, len(creatives))
		for i := range creatives {
			creativeDatas[i] = *(creatives)[i].CreateDeliveryDataCreative()
		}
		expectedError := errors.New("Failed to put")
		gomock.InOrder(
			creativeDataRepository.EXPECT().PutAll(gomock.Eq(ctx), gomock.Eq(&creativeDatas)).Return(expectedError),
		)

		// テストを実行する
		creative := NewCreative(logger, creativeDataRepository, creativeRepository)
		err := creative.Put(ctx, &creativeDatas)
		assert.EqualError(t, err, expectedError.Error())
	})
}

// CreativeのupdateTTLのテスト
func TestCreative_updateTTL(t *testing.T) {
	// テスト用のLoggerを作成
	logger := NewTestLogger(t)
	t.Parallel()

	t.Run("ttlを更新する", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		creativeLog := createTestCreativeLog()
		ttl := time.Now()
		gomock.InOrder(
			creativeDataRepository.EXPECT().UpdateTTL(gomock.Eq(ctx), gomock.Eq(strconv.Itoa(creativeLog.ID)), gomock.Eq(ttl.Unix())).Return(nil),
		)

		// テストを実行する
		creativeUsecase := NewCreative(logger, creativeDataRepository, creativeRepository)
		// private methodのテストを行うためにcastする
		creativeInteractor := creativeUsecase.(*creative)
		err := creativeInteractor.updateTTL(ctx, ttl, creativeLog)
		assert.NoError(t, err)
	})

	t.Run("ttlの更新処理でエラーが発生した場合、エラーを返して終了", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		creativeLog := createTestCreativeLog()
		ttl := time.Now()
		expectedError := errors.New("Failed to put")
		gomock.InOrder(
			creativeDataRepository.EXPECT().UpdateTTL(gomock.Eq(ctx), gomock.Eq(strconv.Itoa(creativeLog.ID)), gomock.Eq(ttl.Unix())).Return(expectedError),
		)

		// テストを実行する
		creativeUsecase := NewCreative(logger, creativeDataRepository, creativeRepository)
		// private methodのテストを行うためにcastする
		creativeInteractor := creativeUsecase.(*creative)
		err := creativeInteractor.updateTTL(ctx, ttl, creativeLog)
		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("ttlの更新処理でErrConditionFailedエラーが発生した場合、エラーは返さないで何もしない", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // 定義したmockの処理が想定どおり呼ばれているかチェックが行われる

		// 必要なmockを作成
		creativeRepository := mock_repository.NewMockCreativeRepository(ctrl)
		creativeDataRepository := mock_repository.NewMockDeliveryDataCreativeRepository(ctrl)

		// mockの処理を定義
		// 引数に渡ると想定される値
		ctx := context.Background()
		creativeLog := createTestCreativeLog()
		ttl := time.Now()

		gomock.InOrder(
			creativeDataRepository.EXPECT().UpdateTTL(gomock.Eq(ctx), gomock.Eq(strconv.Itoa(creativeLog.ID)), gomock.Eq(ttl.Unix())).Return(codes.ErrConditionFailed),
		)

		// テストを実行する
		creativeUsecase := NewCreative(logger, creativeDataRepository, creativeRepository)
		// private methodのテストを行うためにcastする
		creativeInteractor := creativeUsecase.(*creative)
		err := creativeInteractor.updateTTL(ctx, ttl, creativeLog)
		assert.NoError(t, err)
	})
}

func createTestCreatives(expirationDate sql.NullTime) []models.Creative {
	return []models.Creative{
		{
			ID:     1,
			Height: 1.0,
			Width:  1.0,
			URL:    "url1",
			// ExpirationDate: expirationDate,
		},
	}
}

func createTestCreativeLog() *models.CreativeLog {
	return &models.CreativeLog{
		ID:    1,
		Event: "insert",
	}
}
