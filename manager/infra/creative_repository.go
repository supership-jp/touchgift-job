package infra

import (
	"context"
	"fmt"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type CreativeRepository struct {
	logger     *Logger
	sqlHandler SQLHandler
}

func NewCreativeRepository(logger *Logger, sqlHandler SQLHandler) *CreativeRepository {
	return &CreativeRepository{
		logger:     logger,
		sqlHandler: sqlHandler,
	}
}

func (c *CreativeRepository) GetCreativeByCampaignID(ctx context.Context, tx repository.Transaction, args *repository.CreativeByCampaignIDCondition) ([]*models.Creative, error) {

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
		IFNULL(video.endcard_url, '') AS end_card_url,
		video.endcard_width AS end_card_width,
		video.endcard_height AS end_card_height,
		IFNULL(video.endcard_extension, '') AS end_card_extension
	FROM campaign_creative
			 INNER JOIN  campaign ON campaign_creative.campaign_id = campaign.id
			 INNER JOIN creative ON campaign_creative.creative_id = creative.id
			 LEFT JOIN banner ON creative.banner_id = banner.id
			 LEFT JOIN video ON creative.video_id = video.id
	WHERE campaign.id = :campaign_id
	GROUP BY
	  creative.id, banner.id, video.id, video.endcard_url, video.endcard_width, video.endcard_height, video.endcard_extension
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

// repository.Transactionの実態は、infra.Transaction
func (c *CreativeRepository) GetCreative(ctx context.Context, tx repository.Transaction, args *repository.CreativeCondition) ([]models.Creative, error) {
	var data string

	query := fmt.Sprintf(`SELECT
	creative.id as id,
	creative.click_url as link,
	campaign_creative.delivery_rate as delivery_rate,
	%s
	banner.id as banner_id,
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
	max(campaign.end_at) as expiration_date
FROM campaign_creative
	INNER JOIN  campaign ON campaign_creative.campaign_id = campaign.id
	INNER JOIN creative ON campaign_creative.creative_id = creative.id
	LEFT JOIN banner ON creative.banner_id = banner.id
	LEFT JOIN video ON creative.video_id = video.id
WHERE creative.id = :creative_id
GROUP BY
	creative.id, banner.id, video.id, campaign_creative.delivery_rate
ORDER BY creative.id`, data)

	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			c.logger.Error().Err(err).Msg("Failed to close statement")
		}
	}()
	dest := []models.Creative{}
	err = stmt.SelectContext(ctx, &dest, map[string]interface{}{
		"creative_id": args.ID,
		"status":      args.Status,
	})
	return dest, err
}
