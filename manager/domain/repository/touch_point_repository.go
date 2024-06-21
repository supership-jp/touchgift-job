package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type TouchPointByGroupIDCondition struct {
	GroupID int
	Limit   int
}

type TouchPointRepository interface {
	GetTouchPointByGroupID(ctx context.Context, tx Transaction, args *TouchPointByGroupIDCondition) ([]*models.TouchPoint, error)
}
