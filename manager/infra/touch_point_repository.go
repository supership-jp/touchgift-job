package infra

import (
	"context"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type TouchPointRepository struct {
	logger     *Logger
	sqlHandler SQLHandler
}

func NewTouchPointRepository(logger *Logger, sqlHandler SQLHandler) *TouchPointRepository {
	return &TouchPointRepository{
		logger:     logger,
		sqlHandler: sqlHandler,
	}
}

func (t *TouchPointRepository) GetTouchPointByGroupID(ctx context.Context,
	tx repository.Transaction, args *repository.TouchPointByGroupIDCondition) ([]*models.TouchPoint, error) {

	query := `SELECT
		c.store_group_id AS "group_id",
		tp.id AS "touch_point_id"
	FROM campaign c
	JOIN store_map sm on c.store_group_id = sm.store_group_id
	JOIN store s on sm.store_id = s.id
	JOIN touch_point tp on tp.store_id = s.id
	WHERE c.store_group_id IN (:group_id)
	LIMIT :limit`
	stmt, err := tx.(*Transaction).Tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var touchPoints []*models.TouchPoint
	err = stmt.SelectContext(ctx, &touchPoints, map[string]interface{}{
		"group_id": args.GroupID,
		"limit":    args.Limit,
	})

	if err != nil {
		t.logger.Error().Msgf("Error getting deliveries: %v", err)
		return nil, err
	}
	return touchPoints, nil
}