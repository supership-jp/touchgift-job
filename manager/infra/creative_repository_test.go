package infra

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"touchgift-job-manager/domain/repository"
	mock_infra "touchgift-job-manager/mock/infra"
)

func TestCreativeRepository_GetCreativeByCampaignID(t *testing.T) {
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
		creativeRepository := NewCreativeRepository(logger, sqlHandler)
		actuals, err := creativeRepository.GetCreativeByCampaignID(ctx, tx, &repository.CreativeByCampaignIDCondition{
			CampaignID: 0,
			Limit:      1,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(actuals))
		}
	})
	t.Run("データ1を一件返却する", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		tx, err := sqlHandler.Begin(ctx)
		if !assert.NoError(t, err) {
			return
		}
		//	ロールバックする(テストデータは不要なので)
		defer func() {
			err := tx.Rollback()
			assert.NoError(t, err)
		}()
		// テストデータを作成する
		//	テストデータを登録する
		rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)
		// 店舗情報('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
		rdbUtil.InsertStore("ORG001", "S001", "東京本店", "100-0001", "13", "東京都千代田区丸の内1-1-1")
		//　store_group情報登録('グループA', 'ORG001', 1)
		store_group_id := rdbUtil.InsertStoreGroup("グループA", "ORG001", 1)

		gimmick_id := rdbUtil.InsertGimmick("ギミックA", "htttps://gimmck.jpg", "S001", "0", "xxx", 1)
		// キャンペーン情報登録
		campaign_id, _ := rdbUtil.InsertCampaign(
			"ORG001",              // organizationCode
			"configured",          // status
			"Project X",           // name
			"2024-06-01 18:41:11", // startAt
			"2024-06-29 18:41:11", // endAt
			1,                     // lastUpdatedBy
			store_group_id,        // storeGroupId
			gimmick_id,
		)

		// video
		video_id, err := rdbUtil.InsertVideo(
			"https://example.com/video.mp4",   // video_url
			"https://example.com/endcard.jpg", // endcard_url
			"video_xid",                       // video_xid
			"endcard_xid",                     // endcard_xid
			100,                               // height
			200,                               // width
			"mp4",                             // extension
			100,                               // endcard_height
			200,                               // endcard_width
			"jpg",                             // endcard_extension
			1,                                 // last_updated_by
			10,
			"https://example.com/endcard.jpg",
		)

		// creative
		creative_id, err := rdbUtil.InsertCreative(
			"ORG001",
			"0",
			"creative_name",
			"click_url",
			"video",
			1,
			video_id,
		)

		// campaign_creative
		rdbUtil.InsertCampaignCreative(
			campaign_id,
			creative_id,
			100,
			3,
		)

		sqlHandler := mock_infra.NewMockSQLHandler(ctrl)
		creativeRepository := NewCreativeRepository(logger, sqlHandler)
		actuals, err := creativeRepository.GetCreativeByCampaignID(ctx, tx, &repository.CreativeByCampaignIDCondition{
			CampaignID: campaign_id,
			Limit:      10,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(actuals))
		}
	})
}
