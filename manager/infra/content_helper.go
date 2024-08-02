package infra

import (
	"context"
	"errors"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type ContentsHelper struct {
	logger *Logger
}

func NewContentsHelper(logger *Logger) *ContentsHelper {
	return &ContentsHelper{
		logger: logger,
	}
}

func (c *ContentsHelper) GenerateContents(ctx context.Context, args *repository.GenerateContentCondition) ([]*models.Content, error) {
	if args == nil {
		c.logger.Error().Msg("GenerateContentsCondition is nil")
		return nil, errors.New("generate contents condition must not be nil")
	}

	var contents []*models.Content

	// Prepare content data
	content := &models.Content{
		Coupons: make([]models.Coupon, len(args.Coupons)),
		Gimmicks: []models.Gimmick{
			{
				URL: "",
			},
		},
	}

	// Copy coupons data
	for i, coupon := range args.Coupons {
		content.Coupons[i] = models.Coupon{
			ID:       coupon.ID,
			Name:     coupon.Name,
			Code:     coupon.Code,
			ImageURL: coupon.ImageURL,
			Rate:     coupon.Rate,
		}
	}

	// Set gimmick URL if available
	if args.GimmickURL != nil {
		content.Gimmicks[0].URL = *args.GimmickURL
	}

	contents = append(contents, content)

	return contents, nil
}
