package infra

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type DeliveryDataRepository struct {
	logger     *Logger
	sqlHandler SQLHandler
}

func NewDeliveryRepository(logger *Logger, sqlHandler SQLHandler) *DeliveryDataRepository {
	return &DeliveryDataRepository{
		logger:     logger,
		sqlHandler: sqlHandler,
	}
}

func (d *DeliveryDataRepository) GetDeliveryData(ctx context.Context) ([]*models.Delivery, error) {
	// SQLクエリ
	query := `
    SELECT c.id as id, c.name as name, v.video_url as video, v.endcard_url as end_card
    FROM creative c
    JOIN creative_video cv ON cv.creative_id = c.id
    JOIN video v ON cv.video_id = v.id
    WHERE c.creative_type = ?
    `
	var deliveries []*models.Delivery                            // スライスの宣言
	err := d.sqlHandler.Select(ctx, &deliveries, query, "video") // &deliveries でスライスのアドレスを渡す
	if err != nil {
		d.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}

	return deliveries, nil
}
