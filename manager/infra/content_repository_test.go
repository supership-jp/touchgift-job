package infra

import (
	"context"
	"testing"
	"touchgift-job-manager/domain/repository"
	mock_infra "touchgift-job-manager/mock/infra"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestContentsRepository_GetGimmickURLByCampaignID(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()

	t.Run("対象のデータが存在しないためNilを返す", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		tx, err := sqlHandler.Begin(ctx)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()

		sqlHandler := mock_infra.NewMockSQLHandler(ctrl)
		contentsRepository := NewContentRepository(logger, sqlHandler)
		actuals, _, err := contentsRepository.GetGimmicksByCampaignID(ctx, tx, &repository.ContentByCampaignIDCondition{
			CampaignID: 0,
		})

		if assert.NoError(t, err) {
			assert.Nil(t, actuals)
		}

	})
	t.Run("Campaignに紐づいたGimmickが存在する時URLを返却する", func(t *testing.T) {
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

		// テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)

		// 店舗情報登録
		rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

		id, _ := rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-01 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)

		contentsRepository := NewContentRepository(logger, sqlHandler)

		gimmickURL, _, _ := contentsRepository.GetGimmicksByCampaignID(ctx, tx, &repository.ContentByCampaignIDCondition{
			CampaignID: id,
		})

		assert.NotNil(t, gimmickURL)
		assert.Equal(t, "https://gimmck.jpg", *gimmickURL)
	})
}

func TestContentsRepository_GetCouponsByCampaignID(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()

	t.Run("対象のデータが存在しないため0件を返す", func(t *testing.T) {
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
		contentsRepository := NewContentRepository(logger, sqlHandler)
		actuals, err := contentsRepository.GetCouponsByCampaignID(ctx, tx, &repository.ContentByCampaignIDCondition{
			CampaignID: 0,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(actuals))
		}
	})
	t.Run("キャンペーンに紐づくクーポンデータが1件存在する場合クーポンデータを取得できる", func(t *testing.T) {
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

		// テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)

		// 店舗情報登録
		rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

		couponID := rdbUtil.InsertCoupon(
			"Summer Sale",                         // name
			"ORG123",                              // organizationCode
			"2",                                   // status
			"SUMMER2024",                          // code
			"https://example.com/summer-sale.jpg", // imgUrl
			"XID1234",                             // xid
			1,
		)

		campaignID, _ := rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-01 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)

		rdbUtil.InsertCampaignCoupon(campaignID, couponID, 100)

		contentsRepository := NewContentRepository(logger, sqlHandler)
		actuals, err := contentsRepository.GetCouponsByCampaignID(ctx, tx, &repository.ContentByCampaignIDCondition{
			CampaignID: campaignID,
		})

		if assert.NoError(t, err) && assert.Equal(t, 1, len(actuals)) {
			assert.Equal(t, "Summer Sale", actuals[0].Name)
			assert.Equal(t, "SUMMER2024", actuals[0].Code)
			assert.Equal(t, "https://example.com/summer-sale.jpg", actuals[0].ImageURL)
			assert.Equal(t, "100", actuals[0].Rate)
		}

	})
	t.Run("キャンペーンに複数のクーポンデータが存在する場合、全てのクーポンデータを取得できる", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		// トランザクションを開始
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

		// テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)

		// 店舗情報登録
		rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)
		// 複数のクーポンを登録
		couponIDs := []int{
			rdbUtil.InsertCoupon("Summer Sale", "ORG123", "2", "SUMMER2024", "https://example.com/summer-sale.jpg", "XID1234", 1),
			rdbUtil.InsertCoupon("Winter Sale", "ORG123", "2", "WINTER2024", "https://example.com/winter-sale.jpg", "XID5678", 1),
		}

		campaignID, _ := rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-01 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
		)

		for _, couponID := range couponIDs {
			rdbUtil.InsertCampaignCoupon(campaignID, couponID, 50)
		}

		contentsRepository := NewContentRepository(logger, sqlHandler)
		actuals, err := contentsRepository.GetCouponsByCampaignID(ctx, tx, &repository.ContentByCampaignIDCondition{
			CampaignID: campaignID,
		})

		if assert.NoError(t, err) && assert.Equal(t, 2, len(actuals)) {
			assert.Equal(t, "Summer Sale", actuals[0].Name)
			assert.Equal(t, "SUMMER2024", actuals[0].Code)
			assert.Equal(t, "https://example.com/summer-sale.jpg", actuals[0].ImageURL)
			assert.Equal(t, "Winter Sale", actuals[1].Name)
			assert.Equal(t, "WINTER2024", actuals[1].Code)
			assert.Equal(t, "https://example.com/winter-sale.jpg", actuals[1].ImageURL)
			assert.Equal(t, "50", actuals[0].Rate)
			assert.Equal(t, "50", actuals[1].Rate)
		}
	})

}
