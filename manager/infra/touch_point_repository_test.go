package infra

import (
	"context"
	"testing"
	"touchgift-job-manager/domain/repository"
	mock_infra "touchgift-job-manager/mock/infra"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTouchPointRepository_GetTouchPointByGroupID(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()

	t.Run("対象のデータが存在しないため0件を返す", func(t *testing.T) {
		//	mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		//	 トランザクションを開始(トランザクション内でテストする)
		tx, err := sqlHandler.Begin(ctx)
		if !assert.NoError(t, err) {
			return
		}
		//	ロールバックする(テストデータは不要なので)
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()
		touchPointRepository := NewTouchPointRepository(logger, sqlHandler)
		actuals, err := touchPointRepository.GetTouchPointByGroupID(ctx, &repository.TouchPointByGroupIDCondition{
			GroupID: 1,
			Limit:   10,
		})
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(actuals))
		}

	})
	t.Run("1件だけ存在した場合", func(t *testing.T) {
		//	mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		// トランザクションを開始する(トランザクション内でテストする)
		tx, err := sqlHandler.Begin(ctx)
		if !assert.NoError(t, err) {
			return
		}

		//　ロールバックする(テストデータは不要なので)
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()

		//	テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)
		// 店舗情報('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
		store_id := rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録('グループA', 'ORG001', 1)
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

		// キャンペーン情報登録
		_, err = rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-01 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)
		if err != nil {
			if !assert.NoError(t, err) {
				return
			}
		}

		rdbUtil.InsertTouchPoint(
			"ORG001",
			"xxx",
			"yyy",
			store_id,
			"nfc",
			"ポイントA",
			1,
		)

		_, err = rdbUtil.InsertStoreMap(
			store_group_id,
			store_id,
		)

		if !assert.NoError(t, err) {
			return
		}

		_sqlHandler := mock_infra.NewMockSQLHandler(ctrl)
		_sqlHandler.EXPECT().PrepareContext(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context, query string) (*sqlx.Stmt, error) {
			return tx.(*Transaction).Tx.PreparexContext(ctx, query)
		}).Times(1)
		_sqlHandler.EXPECT().In(gomock.Any(), gomock.Any()).DoAndReturn(func(query string, arg interface{}) (*string, []interface{}, error) {
			return sqlHandler.In(query, arg)
		}).Times(1)
		touchPointRepository := NewTouchPointRepository(logger, _sqlHandler)
		actuals, err := touchPointRepository.GetTouchPointByGroupID(ctx, &repository.TouchPointByGroupIDCondition{
			GroupID: store_group_id,
			Limit:   10,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(actuals))
		}

	})
	t.Run("データが複数件存在する場合", func(t *testing.T) {
		//	mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		// トランザクションを開始する(トランザクション内でテストする)
		tx, err := sqlHandler.Begin(ctx)
		if !assert.NoError(t, err) {
			return
		}

		//　ロールバックする(テストデータは不要なので)
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()

		//	テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)

		// 店舗登録
		storeId := rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//	store_group情報登録
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

		//　キャンペーン情報登録
		_, err = rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-01 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)
		if !assert.NoError(t, err) {
			return
		}

		rdbUtil.InsertTouchPoint(
			"ORG001",
			"xxx01",
			"yyy01",
			storeId,
			"nfc",
			"ポイントA",
			1,
		)

		rdbUtil.InsertTouchPoint(
			"ORG001",
			"xxx02",
			"yyy02",
			storeId,
			"nfc",
			"ポイントB",
			1,
		)

		rdbUtil.InsertTouchPoint(
			"ORG001",
			"xxx03",
			"yyy03",
			storeId,
			"nfc",
			"ポイントC",
			1,
		)

		_, err = rdbUtil.InsertStoreMap(
			store_group_id,
			storeId,
		)
		if !assert.NoError(t, err) {
			return
		}

		_sqlHandler := mock_infra.NewMockSQLHandler(ctrl)
		_sqlHandler.EXPECT().PrepareContext(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context, query string) (*sqlx.Stmt, error) {
			return tx.(*Transaction).Tx.PreparexContext(ctx, query)
		}).Times(1)
		_sqlHandler.EXPECT().In(gomock.Any(), gomock.Any()).DoAndReturn(func(query string, arg interface{}) (*string, []interface{}, error) {
			return sqlHandler.In(query, arg)
		}).Times(1)
		touchPointRepository := NewTouchPointRepository(logger, _sqlHandler)
		actuals, err := touchPointRepository.GetTouchPointByGroupID(ctx, &repository.TouchPointByGroupIDCondition{
			GroupID: store_group_id,
			Limit:   10,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, 3, len(actuals))
		}

	})
}
