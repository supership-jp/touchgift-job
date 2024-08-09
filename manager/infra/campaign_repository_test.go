package infra

import (
	"context"
	"fmt"
	"testing"
	"time"
	"touchgift-job-manager/domain/repository"
	mock_infra "touchgift-job-manager/mock/infra"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCampaignRepository_GetCampaignToStart(t *testing.T) {
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
		_, err = createStartCampaignData(ctx, t, tx, time.Now().Local().Add(time.Duration(-1)*time.Hour).Format("2006-01-02 15:04:05"))
		assert.NoError(t, err)

		repo := NewCampaignRepository(logger, sqlHandler)
		campaigns, _ := repo.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
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

		//	start_atが現在時刻から1時間後のテストデータを登録する
		id, err := createStartCampaignData(ctx, t, tx, time.Now().Local().Add(time.Duration(1)*time.Hour).Format("2006-01-02 15:04:05"))
		assert.NoError(t, err)
		assert.NotNil(t, id)

		repo := NewCampaignRepository(logger, sqlHandler)
		campaigns, _ := repo.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
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
			Status: "configured",
		})
		assert.Len(t, campaigns, 0)
	})
}
func TestCampaignRepository_GetCampaignToEnd(t *testing.T) {
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
		campaignRepository := NewCampaignRepository(logger, sqlHandler)
		actuals, err := campaignRepository.GetCampaignToEnd(ctx, &repository.CampaignDataToEndCondition{
			End:    time.Now(),
			Status: []string{"configured"},
		})
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(actuals))
		}
	})
	// TODO: テストが通らない
	t.Run("配信中(status=started)かつ終了日を過ぎているキャンペーンを1件返す", func(t *testing.T) {
		// mockを使用する準備
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

		_sqlHandler := mock_infra.NewMockSQLHandler(ctrl)
		_sqlHandler.EXPECT().PrepareContext(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context, query string) (*sqlx.Stmt, error) {
			return tx.(*Transaction).Tx.PreparexContext(ctx, query)
		}).Times(1)
		_sqlHandler.EXPECT().In(gomock.Any(), gomock.Any()).DoAndReturn(func(query string, arg interface{}) (*string, []interface{}, error) {
			return sqlHandler.In(query, arg)
		}).Times(1)
		//	end_atが現在時刻から1時間前のテストデータを登録する
		_, err = createEndedCampaignData(ctx, t, tx, time.Now().Local().Add(time.Duration(-1)*time.Hour).Format("2006-01-02 15:04:05"))
		assert.NoError(t, err)

		campaignRepository := NewCampaignRepository(logger, _sqlHandler)
		actuals, err := campaignRepository.GetCampaignToEnd(ctx, &repository.CampaignDataToEndCondition{
			End:    time.Now().Local(),
			Status: []string{"started", "warmup"},
		})
		assert.NoError(t, err)
		fmt.Println(actuals)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(actuals))
		}
	})
	t.Run("配信中(status=started)かつ終了日を過ぎていない場合はキャンペーンを返却しない", func(t *testing.T) {
		// mockを使用する準備
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		tx, err := sqlHandler.Begin(ctx)

		if !assert.NoError(t, err) {
			return
		}

		// ロールバックする
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()
		sqlHandler := mock_infra.NewMockSQLHandler(ctrl)

		id, err := createEndedCampaignData(ctx, t, tx, time.Now().Local().Add(time.Duration(1)*time.Hour).Format("2006-01-02 15:04:05"))
		assert.NoError(t, err)
		assert.NotNil(t, id)

		repo := NewCampaignRepository(logger, sqlHandler)
		campaigns, _ := repo.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
			Status: "configured",
		})

		assert.Len(t, campaigns, 0)
	})
}

func TestCampaignRepository_UpdateStatus(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()

	t.Run("正常にステータスがupdateされる", func(t *testing.T) {
		// mockを使用する準備
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
		_, err = createStartCampaignData(ctx, t, tx, time.Now().Local().Add(time.Duration(-1)*time.Hour).Format("2006-01-02 15:04:05"))
		assert.NoError(t, err)

		// リポジトリの取得
		campaignRepository := NewCampaignRepository(logger, sqlHandler)

		campaigns, err := campaignRepository.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
			Status: "configured",
		})
		assert.NoError(t, err)

		campaign_id := campaigns[0].ID

		updatedID, err := campaignRepository.UpdateStatus(ctx, tx, &repository.UpdateCondition{
			CampaignID: campaign_id,
			Status:     "started",
			UpdatedAt:  time.Now(),
		})
		if err != nil {
			assert.NoError(t, err)
		}
		assert.Equal(t, campaigns[0].ID, updatedID)

		campaigns, err = campaignRepository.GetCampaignToStart(ctx, tx, &repository.CampaignToStartCondition{
			To:     time.Now(),
			Status: "started",
		})
		assert.NoError(t, err)
		// startedのキャンペーンが取得できるか
		assert.Len(t, campaigns, 1)
	})
}

// TODO: 時間を現在時刻を基準に
func createStartCampaignData(ctx context.Context, t testing.TB, tx repository.Transaction, startAt string) (*int, error) {
	rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)
	// 店舗情報('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
	rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
	//　store_group情報登録('グループA', 'ORG001', 1)
	store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

	// キャンペーン情報登録
	id, err := rdbUtil.InsertCampaign(
		"ORG001",     // organizationCode
		"configured", // status
		"Project X",  // name
		startAt,      // startAt
		time.Now().Local().Add(time.Duration(1)*time.Hour).Format("2006-01-02 15:04:05"), // endAt
		1,              // lastUpdatedBy
		store_group_id, // storeGroupId
	)
	return &id, err
}

// 終了したキャンペーン作成
func createEndedCampaignData(ctx context.Context, t testing.TB, tx repository.Transaction, endAt string) (*int, error) {
	rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)
	// 店舗情報('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
	rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
	//　store_group情報登録('グループA', 'ORG001', 1)
	store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

	// キャンペーン情報登録
	id, err := rdbUtil.InsertCampaign(
		"ORG001",              // organizationCode
		"started",             // status
		"Project X",           // name
		"2006-06-01 18:41:11", // startAt
		endAt,                 // endAt
		1,                     // lastUpdatedBy
		store_group_id,        // storeGroupId
	)
	if err != nil {
		fmt.Println(err)
	}
	return &id, err
}
