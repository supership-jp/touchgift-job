package infra

import (
	"context"
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

// GetCampaignDataToStart 配信開始時間になったらキャンペーン対象のキャンペーン情報を取得する
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
	c.status = :status
LIMIT :limit`
	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var campaigns []*models.Campaign
	err = stmt.SelectContext(ctx, &campaigns, map[string]interface{}{
		"to":     args.To,
		"limit":  args.Limit,
		"status": args.Status,
	})
	if err != nil {
		c.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}
	return campaigns, nil
}
