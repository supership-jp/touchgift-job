package infra

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"touchgift-job-manager/domain/repository"
	mock_infra "touchgift-job-manager/mock/infra"
)

func TestCampaignDataRepository_GetCampaignToStart(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()

	t.Run("空データのためデータ無し0件返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		// トランザクションを開始(トランザクション内でテストする)
		tx, err := sqlHandler.Begin(ctx)
		if !assert.NoError(t, err) {
			return
		}
		// ロールバックする(テストデータは不要なので)
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()
		sqlHandler := mock_infra.NewMockSQLHandler(ctrl)
		campaignRepository := NewCampaignRepository(logger, sqlHandler)
		actuals, err := campaignRepository.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
			Limit:  1,
			Status: "configured",
		})
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(actuals))
		}
	})
	t.Run("配信前(status=configured)かつ期間範囲内かつ条件にStatus=configuredを指定した場合、データは1件返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		// トランザクションを開始する(トランザクション内でテストする)
		tx, err := sqlHandler.Begin(ctx)

		if !assert.NoError(t, err) {
			return
		}

		//	ロールバックする(テストデータは不要なので)
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()
		sqlHandler := mock_infra.NewMockSQLHandler(ctrl)

		//	テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)
		// 店舗情報('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
		rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録('グループA', 'ORG001', 1)
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

		// キャンペーン情報登録
		id, _ := rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-01 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)
		assert.NotNil(t, id)

		repo := NewCampaignRepository(logger, sqlHandler)
		campaigns, _ := repo.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
			Limit:  1,
			Status: "configured",
		})

		//tx.Commit()
		assert.Len(t, campaigns, 1)
	})
	t.Run("期間範囲外のためデータは0件を返す", func(t *testing.T) {
		//	mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		//	トランザクションを開始する
		tx, err := sqlHandler.Begin(ctx)

		if !assert.NoError(t, err) {
			return
		}

		//	ロールバックする
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()
		sqlHandler := mock_infra.NewMockSQLHandler(ctrl)

		//	テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)
		// 店舗情報('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
		rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録('グループA', 'ORG001', 1)
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

		// キャンペーン情報登録
		id, _ := rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-28 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)

		assert.NotNil(t, id)

		repo := NewCampaignRepository(logger, sqlHandler)
		campaigns, _ := repo.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
			Limit:  1,
			Status: "configured",
		})

		assert.Len(t, campaigns, 0)
	})
	t.Run("配信前(status=configured)以外で条件にStatus=configuredを指定した場合、データは0件返す", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		// トランザクションを開始する
		tx, err := sqlHandler.Begin(ctx)

		if !assert.NoError(t, err) {
			return
		}

		//	ロールバックする
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()
		sqlHandler := mock_infra.NewMockSQLHandler(ctrl)

		//	テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)
		// 店舗情報('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
		rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録('グループA', 'ORG001', 1)
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)
		// キャンペーン情報登録
		id, _ := rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"start",               // status
			"Project X",           // name
			"2024-06-28 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)

		assert.NotNil(t, id)
		repo := NewCampaignRepository(logger, sqlHandler)
		campaigns, _ := repo.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
			Limit:  1,
			Status: "configured",
		})
		assert.Len(t, campaigns, 0)
	})
}
