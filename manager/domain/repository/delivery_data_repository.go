package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type DeliveryDataCampaign interface {
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

type DeliveryDataTouchPoint interface {
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

type DeliveryDataCreative interface {
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

type DeliveryDataContent interface {
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
