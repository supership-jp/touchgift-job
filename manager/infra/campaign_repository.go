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
    c.organization_code as org_id,
    c.name as name,
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
    c.organization_code as org_id,
    c.name as name,
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
