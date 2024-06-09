package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type creativeRepository interface {
	GetDeliveryData(ctx context.Context) ([]*models.Delivery, error)
}
