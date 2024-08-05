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

	// 2分後に始まり、5分後に終わるキャンペーン(セグメントなし)
	storeID := rdbUtil.InsertStore("org_1", "store1", "store_name1", "0000000", "00", "住所")
	groupID := rdbUtil.InsertStoreGroup("group1", "org_1", 1)
	gimmickID := rdbUtil.InsertGimmick("gimmick1", "http://example.com", "org_1", "2", "code1", "xid", 1)
	startAt, endAt := time.Now().Add(2*time.Minute).Local(), sql.NullTime{Time: time.Now().Add(5 * time.Minute).Local(), Valid: true}
	campaignID, _ := rdbUtil.InsertCampaign("org_1", "configured", "campaign1", startAt.Format("2006-01-02 15:04:05"), endAt.Time.Format("2006-01-02 15:04:05"), 1, groupID)
	couponID := rdbUtil.InsertCoupon("coupon1", "org_1", "2", "code1", "http://example.com", "xid", 1)
	rdbUtil.InsertCampaignCoupon(campaignID, couponID, 100)
	rdbUtil.InsertCampaignGimmick(campaignID, gimmickID)
	videoID, _ := rdbUtil.InsertVideo("video_url", "endcard_url", "video_xid", "endcard_xid", 100, 100, "mp4", 100, 100, "png", 1, 100, "http://example.com")
	creativeID, _ := rdbUtil.InsertCreative("org_1", "2", "creative1", "http://example.com", 1, videoID)
	rdbUtil.InsertCampaignCreative(1, creativeID, 100, 1)
	rdbUtil.InsertTouchPoint("org_1", "point_unique_id", "print_management_id", storeID, "QR", "name", 1)

	/*
		// 5分後に始まり、終了日なしキャンペーン。(セグメントあり)
		StoreID2 := rdbUtil.InsertStore("adversier2", "code2", "Store_origin2")
		campaignID2 := rdbUtil.InsertOrder("origin2", "code2", StoreID2)
		bannerID2 := rdbUtil.InsertBanner(1, 1, "image_url2")
		contentID2 := rdbUtil.InsertContent("title2")
		itemCategoryIds2, _ := json.Marshal([]string{"1012"})
		creativeID2 := rdbUtil.InsertCreative("click_url2", bannerID2, "code2", contentID2, StoreID2, itemCategoryIds2)
		startAt2, endAt2 := time.Now().Add(5*time.Minute).Local(), sql.NullTime{Time: time.Time{}, Valid: false}
		campaignID2, _ := rdbUtil.InsertScheduleWithStatus(campaignID2, "campaign2", startAt2, endAt2, "configured", "campaign_origin2", 100, "cpc", 1000, false)
		rdbUtil.InsertScheduleCreative(campaignID2, creativeID2, 1)
		locationID2 := rdbUtil.InsertLocation("location2", "org1")
		rdbUtil.InsertScheduleLocation(campaignID2, locationID2, 1)
		segmentDataID2 := rdbUtil.InsertSegmentData("description1", "type1")
		sectionID2 := rdbUtil.InsertSection("section", "description2", "org", "service", segmentDataID2)
		segmentID2 := rdbUtil.InsertSegment("segment1", "description2", "key1", "code1", sectionID2)
		segmentKeyID2 := rdbUtil.InsertSegmentKey(segmentID2, 1, "description")
		rdbUtil.InsertScheduleSegment(campaignID2, segmentID2, segmentKeyID2)

		// 3分後に始まり10分後に終わるキャンペーン(item_category_idsは[])
		StoreID3 := rdbUtil.InsertStore("adversier3", "code3", "Store_origin3")
		campaignID3 := rdbUtil.InsertOrder("origin3", "code3", StoreID3)
		bannerID3 := rdbUtil.InsertBanner(1, 1, "image_url3")
		contentID3 := rdbUtil.InsertContent("title3")
		itemCategoryIds3, _ := json.Marshal([]string{})
		creativeID3 := rdbUtil.InsertCreative("click_url3", bannerID3, "code3", contentID3, StoreID3, itemCategoryIds3)
		startAt3, endAt3 := time.Now().Add(5*time.Minute).Local(), sql.NullTime{Time: time.Time{}, Valid: false}
		campaignID3, _ := rdbUtil.InsertScheduleWithStatus(campaignID3, "campaign3", startAt3, endAt3, "configured", "campaign_origin3", 100, "cpc", 1000, false)
		rdbUtil.InsertScheduleCreative(campaignID3, creativeID3, 1)
		locationID3 := rdbUtil.InsertLocation("location3", "org1")
		rdbUtil.InsertScheduleLocation(campaignID3, locationID3, 1)
	*/
}
