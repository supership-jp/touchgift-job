package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

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
	Get(ctx context.Context, id *string) (*models.TouchPoint, error)
	//	登録/更新する
	Put(ctx context.Context, updateData *models.TouchPoint) error
	// まとめて登録更新する
	PutAll(ctx context.Context, updateData *[]models.TouchPoint) error
	// 削除する
	Delete(ctx context.Context, campaignID *string) error
	// まとめて削除する
	DeleteAll(ctx context.Context, deleteDatas *[]models.TouchPoint) error
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
}

type DeliveryDataContentRepository interface {
	// 取得する
	Get(ctx context.Context, id *string) (*models.DeliveryDataContent, error)
	//	登録/更新する
	Put(ctx context.Context, updateData *models.DeliveryDataContent) error
	// まとめて登録更新する
	PutAll(ctx context.Context, updateData *[]models.DeliveryDataContent) error
	// 削除する
	Delete(ctx context.Context, campaignID *string) error
	// まとめて削除する
	DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryDataContent) error
}
