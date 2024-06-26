package infra

import (
	"context"
	"github.com/jmoiron/sqlx"
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

func (d *CampaignRepository) UpdateStatus(ctx context.Context, tx repository.Transaction, target *repository.UpdateCondition) ([]int, error) {
	query, args, err := sqlx.In(
		`UPDATE campaign
         SET status = ?, updated_at = NOW()
         WHERE campaign.id IN (?)`,
		target.Status, target.CampaignIDs)
	if err != nil {
		return nil, err
	}

	// クエリをデータベースのプレースホルダ形式に再バインドする
	query = tx.(*Transaction).Tx.Rebind(query)

	// UPDATE文を実行
	_, err = tx.(*Transaction).Tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return target.CampaignIDs, nil
}
