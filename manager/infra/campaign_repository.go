package infra

import (
	"context"
	"fmt"
	"time"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type CampaignRepository struct {
	logger     *Logger
	sqlHandler SQLHandler
}

func NewCampaignRepository(logger *Logger, sqlHandler SQLHandler) repository.CampaignRepository {
	campaignRepository := &CampaignRepository{
		logger:     logger,
		sqlHandler: sqlHandler,
	}
	return campaignRepository
}

// GetCampaignToStart 配信開始するキャンペーン情報を取得する
func (c *CampaignRepository) GetCampaignToStart(ctx context.Context, tx repository.Transaction, args *repository.CampaignToStartCondition) ([]*models.Campaign, error) {
	query := `SELECT
    c.id as id,
    sg.id as group_id,
    c.organization_code as org_code,
		c.daily_coupon_limit_per_user as daily_coupon_limit_per_user,
    c.status as status,
	c.updated_at as updated_at
FROM campaign c
INNER JOIN store_group sg ON c.store_group_id = sg.id
WHERE
    c.start_at <= :to AND
	c.status = :status`
	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var campaigns []*models.Campaign
	err = stmt.SelectContext(ctx, &campaigns, map[string]interface{}{
		"to":     args.To,
		"status": args.Status,
	})
	if err != nil {
		c.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}
	return campaigns, nil
}

// GetCampaignToEnd 配信が終了するキャンペーン情報を取得する
func (c *CampaignRepository) GetCampaignToEnd(ctx context.Context, tx repository.Transaction, args *repository.CampaignDataToEndCondition) ([]*models.Campaign, error) {
	query := `SELECT
    c.id as id,
    sg.id as group_id,
    c.organization_code as org_code,
    c.daily_coupon_limit_per_user as daily_coupon_limit_per_user,
    c.status as status,
	c.updated_at as updated_at
FROM campaign c
INNER JOIN store_group sg ON c.store_group_id = sg.id
WHERE
    c.end_at <= :end AND
	c.status = :status`
	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var campaigns []*models.Campaign
	err = stmt.SelectContext(ctx, &campaigns, map[string]interface{}{
		"end":    args.End,
		"status": args.Status,
	})
	if err != nil {
		c.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}
	return campaigns, nil
}

func (c *CampaignRepository) GetDeliveryToStart(ctx context.Context,
	tx repository.Transaction, args *repository.CampaignCondition,
) (*models.Campaign, error) {
	query := `SELECT
		c.id as id,
		sg.id as group_id,
		c.organization_code as org_code,
		c.daily_coupon_limit_per_user as daily_coupon_limit_per_user,
		c.status as status
	FROM campaign c
	INNER JOIN store_group sg ON c.store_group_id = sg.id
	WHERE
		c.id = :id`
	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// TODO: 修正する(PKでフィルタリングしているため配列ではなく構造体で取得する)
	campaigns := []models.Campaign{}
	err = stmt.SelectContext(ctx, &campaigns, map[string]interface{}{
		"id": args.CampaignID,
	})
	if err != nil {
		c.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}
	campaign := campaigns[0]
	return &campaign, nil
}

func (c *CampaignRepository) GetCreativeByCampaignID(ctx context.Context, tx repository.Transaction, args *repository.CreativeByCampaignIDCondition) ([]*models.Creative, error) {
	query := `
	SELECT
		creative.id as id,
		COALESCE(banner.height, video.height) AS height,
		COALESCE(banner.width, video.width) AS width,
		COALESCE(banner.img_url, video.video_url) AS url,
		CASE
			WHEN banner.id IS NOT NULL THEN 'banner'
			WHEN video.id IS NOT NULL THEN 'video'
		END AS type,
		CASE
			WHEN banner.id IS NOT NULL THEN banner.extension
			WHEN video.id IS NOT NULL THEN video.extension
		END AS extension,
		campaign_creative.skip_offset AS skip_offset,
		video.endcard_url AS end_card_url,
		video.endcard_width AS end_card_width,
		video.endcard_height AS end_card_height,
		video.endcard_extension AS end_card_extension
	FROM campaign_creative
			 INNER JOIN  campaign ON campaign_creative.campaign_id = campaign.id
			 INNER JOIN creative ON campaign_creative.creative_id = creative.id
			 LEFT JOIN banner ON creative.banner_id = banner.id
			 LEFT JOIN video ON creative.video_id = video.id
	WHERE campaign.id = :campaign_id
	LIMIT :limit
`
	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var creatives []*models.Creative
	err = stmt.SelectContext(ctx, &creatives, map[string]interface{}{
		"campaign_id": args.CampaignID,
		"limit":       args.Limit,
	})
	if err != nil {
		c.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}

	return creatives, nil
}

func (c *CampaignRepository) UpdateStatus(ctx context.Context, tx repository.Transaction, target *repository.UpdateCondition) (int, error) {
	query := `UPDATE campaign
	SET
	    status = ?,
	    updated_at = ?
	WHERE id = ?`

	// PrepareContextを使用してステートメントを準備します
	stmt, err := tx.(*Transaction).Tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// ExecContextを使用してクエリを実行し、結果を確認します
	result, err := stmt.ExecContext(ctx, target.Status, time.Now(), target.CampaignID)
	if err != nil {
		return 0, err
	}

	// RowsAffectedをチェックして、実際に更新が行われたかを確認します
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows updated")
	}

	// 成功した場合、更新されたcampaign_idを戻り値として返します。
	return target.CampaignID, nil
}
