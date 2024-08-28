//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/domain/models"
)

// CreativeCondition is struct
type CreativeCondition struct {
	ID                int
	Status            string
	JobProcessedState codes.JobProcessedState
}

type DeliveryDataCampaignRepository interface {
	// 取得する
	Get(ctx context.Context, id *string) (*models.DeliveryDataCampaign, error)
	//	登録/更新する
	Put(ctx context.Context, updateData *models.DeliveryDataCampaign) error
	// まとめて登録更新する
	PutAll(ctx context.Context, updateData *[]models.DeliveryDataCampaign) error
	// 削除する
	Delete(ctx context.Context, campaignID *string) error
	// まとめて削除する
	DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryDataCampaign) error
}

type DeliveryDataTouchPointRepository interface {
	// 取得する
	Get(ctx context.Context, groupID *string) (*models.DeliveryTouchPoint, error)
	//	登録/更新する
	Put(ctx context.Context, updateData *models.DeliveryTouchPoint) error
	// まとめて登録更新する
	PutAll(ctx context.Context, updateData *[]models.DeliveryTouchPoint) error
	// 削除する
	Delete(ctx context.Context, id *string, groupID *string) error
	// まとめて削除する
	DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryTouchPoint) error
	// TTLを更新する (更新対象がない場合 codes.ErrConditionFailed)
	// UpdateTTL(ctx context.Context, id string, ttl int64) error
}

type DeliveryDataCreativeRepository interface {
	// 取得する
	Get(ctx context.Context, id *string) (*models.DeliveryDataCreative, error)
	//	登録/更新する
	Put(ctx context.Context, updateData *models.DeliveryDataCreative) error
	// まとめて登録更新する
	PutAll(ctx context.Context, updateData *[]models.DeliveryDataCreative) error
	// 削除する
	Delete(ctx context.Context, campaignID *string) error
	// まとめて削除する
	DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryDataCreative) error
	// TTLを更新する (更新対象がない場合 codes.ErrConditionFailed)
	UpdateTTL(ctx context.Context, id string, ttl int64) error
}

type DeliveryDataContentRepository interface {
	// 取得する
	Get(ctx context.Context, campaignID *string) (*models.DeliveryDataContent, error)
	//	登録/更新する
	Put(ctx context.Context, updateData *models.DeliveryDataContent) error
	// まとめて登録更新する
	PutAll(ctx context.Context, updateData *[]models.DeliveryDataContent) error
	// 削除する
	Delete(ctx context.Context, campaignID *string) error
	// まとめて削除する
	DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryDataContent) error
}
