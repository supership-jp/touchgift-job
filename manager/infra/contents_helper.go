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

func (c *ContentsHelper) GenerateContents(ctx context.Context, args *repository.GenerateContentsCondition) ([]*models.Contents, error) {
	if args == nil {
		c.logger.Error().Msg("GenerateContentsCondition is nil")
		return nil, errors.New("generate contents condition must not be nil")
	}

	var contents []*models.Contents

	// Prepare content data
	content := &models.Contents{
		Coupons:    make([]models.Coupon, len(args.Coupons)),
		GimmickURL: "",
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
		content.GimmickURL = *args.GimmickURL
	}

	contents = append(contents, content)

	return contents, nil
}
