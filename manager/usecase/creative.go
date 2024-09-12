//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"strconv"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

// Creative is interface
type Creative interface {
	Process(ctx context.Context, tx repository.Transaction, current time.Time, creativeLogs *[]models.CreativeLog) error
	// クリエイティブを登録/更新する
	Put(ctx context.Context, creatives *[]models.DeliveryDataCreative) error
}

type creative struct {
	logger                 Logger
	creativeDataRepository repository.DeliveryDataCreativeRepository
	creativeRepository     repository.CreativeRepository
}

// NewCreative is function
func NewCreative(
	logger Logger,
	creativeDataRepository repository.DeliveryDataCreativeRepository,
	creativeRepository repository.CreativeRepository,
) Creative {
	instance := creative{
		logger:                 logger,
		creativeDataRepository: creativeDataRepository,
		creativeRepository:     creativeRepository,
	}
	return &instance
}

func (c *creative) Process(ctx context.Context, tx repository.Transaction, current time.Time, creativeLogs *[]models.CreativeLog) error {
	for i := range *creativeLogs {
		creativeLog := (*creativeLogs)[i]
		switch creativeLog.Event {
		case "insert", "update":
			//　Campaign処理で処理済のため何もしない
		case "delete":
			condition := repository.CreativeCondition{
				ID: creativeLog.ID,
			}
			creatives, err := c.creativeRepository.GetCreative(ctx, tx, &condition)
			if err != nil {
				return err
			}
			if len(creatives) == 0 {
				c.logger.Info().Time("current", current).Str("org_code", creativeLog.OrgCode).Int("creative_id", creativeLog.ID).Msg("Delete (change ttl)")
				// どのキャンペーンにも紐付かないデータの場合、有効期限(TTL)を1日後に更新
				ttl := time.Now().Add(24 * time.Hour).Truncate(time.Millisecond)
				if err := c.updateTTL(ctx, ttl, &creativeLog); err != nil {
					return err
				}
			}
			// 違うキャンペーンに紐づく場合は何もしない (そのキャンペーンの配信データの方で操作される)
			return nil
		default:
			c.logger.Error().Interface("creative_log", creativeLog).Msg("Unknown event")
		}
	}
	return nil
}

// TTLを更新する
func (c *creative) updateTTL(ctx context.Context, ttl time.Time, creativeLog *models.CreativeLog) error {
	creativeID := strconv.Itoa(creativeLog.ID)
	// DynamoDB更新
	if err := c.creativeDataRepository.UpdateTTL(ctx, creativeID, ttl.Unix()); err != nil && err != codes.ErrConditionFailed {
		return err
	}
	return nil
}

func (c *creative) Put(ctx context.Context, creatives *[]models.DeliveryDataCreative) error {
	if err := c.creativeDataRepository.PutAll(ctx, creatives); err != nil {
		return err
	}
	return nil
}
