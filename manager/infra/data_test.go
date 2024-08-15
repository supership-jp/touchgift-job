//go:build createdata
// +build createdata

package infra

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestData(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	ctx := context.Background()
	tx, err := sqlHandler.Begin(ctx)
	if !assert.NoError(t, err) {
		t.Fatal("Failed to begin transaction")
	}

	defer func() {
		if err := tx.Commit(); err != nil {
			t.Error(err)
		}
		sqlHandler.Close()
	}()

	rdbUtil := NewTouchGiftRDBUtil(ctx, t, tx)

	// 1分後に始まり、5分後に終わるキャンペーン(クリエイティブあり)
	storeID := rdbUtil.InsertStore("org_1", "store1", "store_name1", "0000000", "00", "住所")
	groupID := rdbUtil.InsertStoreGroup("group1", "org_1", 1)
	gimmickID := rdbUtil.InsertGimmick("gimmick1", "http://example.com", "org_1", "2", "code1", "xid", 1)
	startAt, endAt := time.Now().Add(1*time.Minute).Local(), sql.NullTime{Time: time.Now().Add(5 * time.Minute).Local(), Valid: true}
	campaignID, err := rdbUtil.InsertCampaign("org_1", "configured", "campaign1", startAt.Format("2006-01-02 15:04:05"), endAt.Time.Format("2006-01-02 15:04:05"), 1, groupID)
	if !assert.NoError(t, err) {
		t.Fatal("Failed to insert campaign")
	}
	couponID := rdbUtil.InsertCoupon("coupon1", "org_1", "2", "code1", "http://example.com", "xid", 1)
	rdbUtil.InsertCampaignCoupon(campaignID, couponID, 100)
	rdbUtil.InsertCampaignGimmick(campaignID, gimmickID)
	videoID, _ := rdbUtil.InsertVideo("video_url", "endcard_url", "video_xid", "endcard_xid", 100, 100, "mp4", 100, 100, "png", 1, 100, "http://example.com")
	creativeID, _ := rdbUtil.InsertCreative("org_1", "2", "creative1", "http://example.com", 1, videoID)
	_, _ = rdbUtil.InsertCampaignCreative(1, creativeID, 100, 1)
	rdbUtil.InsertTouchPoint("org_1", "point_unique_id", "print_management_id", storeID, "QR", "name", 1)
	_, _ = rdbUtil.InsertStoreMap(groupID, storeID)
}
