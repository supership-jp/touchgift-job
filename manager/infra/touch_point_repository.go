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
	args *repository.TouchPointByGroupIDCondition) ([]*models.TouchPoint, error) {

	query := `SELECT
		c.store_group_id AS "group_id",
		tp.point_unique_id AS "id",
		tp.store_id AS "store_id",
	FROM campaign c
	JOIN store_map sm on c.store_group_id = sm.store_group_id
	JOIN store s on sm.store_id = s.id
	JOIN touch_point tp on tp.store_id = s.id
	WHERE c.store_group_id IN (:group_id)
	LIMIT :limit`
	params := map[string]interface{}{
		"group_id": args.GroupID,
		"limit":    args.Limit,
	}
	_query, _params, err := t.sqlHandler.In(query, params)
	if err != nil {
		return nil, err
	}
	stmt, err := t.sqlHandler.PrepareContext(ctx, *_query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			t.logger.Error().Err(err).Msg("Failed to close statement")
		}
	}()
	var touchPoints []*models.TouchPoint
	err = stmt.SelectContext(ctx, &touchPoints, _params...)
	return touchPoints, err
}
