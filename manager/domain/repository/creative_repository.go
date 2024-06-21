package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type CreativeByCampaignIDCondition struct {
	CampaignID int
	Limit      int
}

type CreativeRepository interface {
	GetCreativeByCampaignID(ctx context.Context, args *CreativeByCampaignIDCondition) ([]*models.Creative, error)
}
