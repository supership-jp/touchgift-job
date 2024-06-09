package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type StoreRepository interface {
	Get(ctx context.Context, id int) (*models.Store, error)
	Select(ctx context.Context) ([]*models.Store, error)
}
