package infra

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"touchgift-job-manager/domain/repository"
)

func TestGetCampaignDataToStart(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()
	t.Run("現在時刻が配信時間内のキャンペーンを取得", func(t *testing.T) {
		ctx := context.Background()

		condition := &repository.CampaignToStartCondition{
			To:     time.Now(),
			Status: "configured",
			Limit:  5,
		}

		campaignRepository := NewCampaignDataRepository(logger, sqlHandler)
		campaigns, err := campaignRepository.GetCampaignToStart(ctx, condition)
		assert.Nil(t, err)
		assert.Len(t, campaigns, 1)
	})
}
