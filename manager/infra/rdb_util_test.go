package infra

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"touchgift-job-manager/domain/repository"
)

type RDBUtil struct {
	ctx context.Context
	t   testing.TB
	tx  *sqlx.Tx
}

func NewTouchGiftRDBUtil(ctx context.Context, t testing.TB, tx repository.Transaction) *RDBUtil {
	return &RDBUtil{
		ctx: ctx,
		t:   t,
		tx:  tx.(*Transaction).Tx,
	}
}

func (r *RDBUtil) InsertTouchPoint(organization_code string, point_unique_id string, print_management_id string, store_id int, type_param string, name string, last_updated_by int) int {
	_, err := r.tx.ExecContext(r.ctx, `
				INSERT INTO touch_point (
				  organization_code,
				  point_unique_id,
				  print_management_id,
				  store_id,
				  type,
				  name,
				  last_updated_by        
				) VALUES (
				    ?, ?, ?, ?, ?, ?, ?
				)`, organization_code, point_unique_id, print_management_id, store_id, type_param, name, last_updated_by)
	if !assert.NoError(r.t, err) {
		r.t.Fatal(err)
	}

	var id int
	err = r.tx.QueryRowxContext(r.ctx, `SELECT LAST_INSERT_ID()`).Scan(&id)
	if !assert.NoError(r.t, err) {
		r.t.Fatal(err)
	}
	return id

}

//INSERT INTO `store` (organization_code, store_id, name, post_code, province_code, address)
//VALUES
//    ('ORG001', 'S001', '東京本店', '100-0001', '13', '東京都千代田区丸の内1-1-1'),
//    ('ORG002', 'S002', '大阪店', '530-0001', '27', '大阪府大阪市北区梅田1-2-2');

func (r *RDBUtil) InsertStore(organization_code string, store_id string, name string, post_code string, province_code string, address string) int {
	_, err := r.tx.ExecContext(r.ctx,
		`INSERT INTO store (
            organization_code,
            store_id,
            name,
            post_code,
            province_code,
            address
        ) VALUES (?, ?, ?, ?, ?, ?)`, // 変更点: プレースホルダーを `?` に
		organization_code,
		store_id,
		name,
		post_code,
		province_code,
		address,
	)
	if !assert.NoError(r.t, err) {
		r.t.Fatal(err)
	}

	var id int
	err = r.tx.QueryRowxContext(r.ctx, `SELECT LAST_INSERT_ID()`).Scan(&id)
	if !assert.NoError(r.t, err) {
		r.t.Fatal(err)
	}
	return id
}

//INSERT INTO `store_group` (name, organization_code, last_updated_by)
//VALUES
//('グループA', 'ORG001', 1);

func (r *RDBUtil) InsertStoreGroup(name string, organization_code string, last_updated_by int) int {
	_, err := r.tx.ExecContext(r.ctx,
		`INSERT INTO store_group (
			name,
		 	organization_code,
			last_updated_by
        ) VALUES (?, ?, ?)`, // 変更点: プレースホルダーを `?` に
		name,
		organization_code,
		last_updated_by,
	)
	if !assert.NoError(r.t, err) {
		r.t.Fatal(err)
	}

	var id int
	err = r.tx.QueryRowxContext(r.ctx, `SELECT LAST_INSERT_ID()`).Scan(&id)
	if !assert.NoError(r.t, err) {
		r.t.Fatal(err)
	}
	return id
}

func (r *RDBUtil) InsertStoreMap(store_group_id int, store_id int) (int, error) {
	query := `
		INSERT INTO store_map(
		     store_group_id,
			 store_id
		) VALUES (?, ?)
	`

	// ExecContext を使用して SQL を実行
	result, err := r.tx.ExecContext(r.ctx, query,
		store_group_id,
		store_id,
	)

	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	// LAST_INSERT_ID()を使って最後に挿入された行のIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	return int(id), nil // 成功時はIDとnilを返す

}

func (r *RDBUtil) InsertCampaign(organizationCode string, status string, name string, startAt string, endAt string, lastUpdatedBy int, storeGroupId int) (int, error) {
	query := `
        INSERT INTO campaign (
            organization_code,
            status,
            name,
            start_at,
            end_at,
            created_at,
            updated_at,
            last_updated_by,
            store_group_id
        ) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?)`

	// ExecContext を使用して SQL を実行
	result, err := r.tx.ExecContext(r.ctx, query,
		organizationCode,
		status,
		name,
		startAt,
		endAt,
		lastUpdatedBy,
		storeGroupId,
	)
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	// LAST_INSERT_ID()を使って最後に挿入された行のIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	return int(id), nil // 成功時はIDとnilを返す
}

func (r *RDBUtil) InsertVideo(video_url string, endcard_url string, video_xid string, endcard_xid string,
	height int, width int, extension string, endcard_height int, endcard_width int, endcard_extension string, last_updated_by int) (int, error) {
	query := `
        INSERT INTO video (
            video_url,
            endcard_url,
            video_xid,
            endcard_xid,
            height,
            width,
            extension,
            endcard_height,
            endcard_width,
		    endcard_extension,
			last_updated_by
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// ExecContext を使用して SQL を実行
	result, err := r.tx.ExecContext(r.ctx, query,
		video_url,
		endcard_url,
		video_xid,
		endcard_xid,
		height,
		width,
		extension,
		endcard_height,
		endcard_width,
		endcard_extension,
		last_updated_by,
	)
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	// LAST_INSERT_ID()を使って最後に挿入された行のIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	return int(id), nil // 成功時はIDとnilを返す
}

//
//INSERT INTO creative
//(organization_code, name, status, click_url, creative_type, last_updated_by, video_id)
//VALUES
//('ORG123', 'Example Video Creative', '1', 'https://example.com/click_here', 'video', 1, 2);

func (r *RDBUtil) InsertCreative(organizationCode string, status string, name string, click_url string, creative_type string, last_updated_by int, video_id int) (int, error) {
	query := `
        INSERT INTO creative (
            organization_code,
            status,
		    name, 
            click_url,
            creative_type,
            last_updated_by,
            video_id
        ) VALUES (?,?,?,?,?,?,?)`

	// ExecContext を使用して SQL を実行
	result, err := r.tx.ExecContext(r.ctx, query,
		organizationCode,
		status,
		name,
		click_url,
		creative_type,
		last_updated_by,
		video_id,
	)
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	// LAST_INSERT_ID()を使って最後に挿入された行のIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	return int(id), nil // 成功時はIDとnilを返す
}

//INSERT INTO campaign_creative
//(campaign_id, creative_id, delivery_rate)
//VALUES
//    (76, 2, 100);

func (r *RDBUtil) InsertCampaignCreative(campaign_id int, creative_id int, delivery_rate int, skip_offset int) (int, error) {
	query := `
        INSERT INTO campaign_creative (
			campaign_id,
		    creative_id,
		    delivery_rate,
		    skip_offset
        ) VALUES (?,?,?, ?)`

	// ExecContext を使用して SQL を実行
	result, err := r.tx.ExecContext(r.ctx, query,
		campaign_id,
		creative_id,
		delivery_rate,
		skip_offset,
	)
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	// LAST_INSERT_ID()を使って最後に挿入された行のIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err // エラーを呼び出し元に返す
	}

	return int(id), nil // 成功時はIDとnilを返す
}
