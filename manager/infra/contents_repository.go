package infra

import (
	"context"
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

func (c *ContentsRepository) GetContentsByCampaignID(ctx context.Context, tx repository.Transaction, args *repository.ContentsByCampaignIDCondition) ([]*models.Contents, error) {
	return nil, nil
}
