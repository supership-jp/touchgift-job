package infra

import (
	"context"
	"database/sql"
	"errors"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type ContentsRepository struct {
	logger     *Logger
	sqlHandler SQLHandler
}

func NewContentsRepository(logger *Logger, sqlHandler SQLHandler) *ContentsRepository {
	return &ContentsRepository{
		logger:     logger,
		sqlHandler: sqlHandler,
	}
}

func (c *ContentsRepository) GetCouponsByCampaignID(ctx context.Context, tx repository.Transaction, args *repository.ContentsByCampaignIDCondition) ([]*models.Coupon, error) {
	query := `SELECT
    coupon.id AS coupon_id,
    coupon.name AS coupon_name,
    coupon.code AS coupon_code,
    coupon.img_url AS coupon_image_url,
    campaign_coupon.delivery_rate AS coupon_rate
FROM campaign
JOIN campaign_coupon ON campaign.id = campaign_coupon.campaign_id
JOIN coupon ON campaign_coupon.coupon_id = coupon.id
WHERE campaign.id = :campaign_id`

	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var coupons []*models.Coupon
	err = stmt.SelectContext(ctx, &coupons, map[string]interface{}{
		"campaign_id": args.CampaignID,
		"limit":       args.Limit,
	})

	if err != nil {
		c.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}

	return coupons, nil
}

func (c *ContentsRepository) GetGimmickURLByCampaignID(ctx context.Context, tx repository.Transaction, args *repository.ContentsByCampaignIDCondition) (*string, error) {
	query := `SELECT
    gimmick.img_url AS gimmick_url
FROM campaign
JOIN gimmick ON campaign.gimmick_id = gimmick.id
WHERE campaign.id = :campaign_id`

	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var gimmickURL string
	err = stmt.QueryRowxContext(ctx, map[string]interface{}{
		"campaign_id": args.CampaignID,
	}).Scan(&gimmickURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows found, return nil without error
			return nil, nil
		}
		c.logger.Error().Msgf("Error getting gimmick URL: %v", err)
		return nil, err
	}

	return &gimmickURL, nil
}
