package infra

import (
	"context"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type CampaignDataRepository struct {
	logger     *Logger
	sqlHandler SQLHandler
}

func NewCampaignDataRepository(logger *Logger, sqlHandler SQLHandler) repository.CampaignRepository {
	campaignRepository := &CampaignDataRepository{
		logger:     logger,
		sqlHandler: sqlHandler,
	}
	return campaignRepository
}

// GetCampaignDataToStart 配信開始時間になったらキャンペーン対象のキャンペーン情報を取得する
func (c *CampaignDataRepository) GetCampaignToStart(ctx context.Context, args *repository.CampaignToStartCondition) ([]*models.Campaign, error) {
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
	stmt, err := c.sqlHandler.PrepareNamedContext(ctx, query)
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

func (c *CampaignDataRepository) GetTouchPointByGroupID(ctx context.Context, args *repository.CampaignIDCondition) ([]*models.TouchPoint, error) {
	//query := fmt.Sprintf(`SELECT
	//id
	//FROM touch_point`, args)
	//fmt.Println(query)
	return nil, nil
}
