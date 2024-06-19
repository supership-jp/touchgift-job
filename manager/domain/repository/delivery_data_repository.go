package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type DeliveryDataRepository interface {
	// Get 取得する
	Get(ctx context.Context, campaignID int) (*models.Campaign, error)
	// Put 登録/更新する
	Put(ctx context.Context, campaign *models.Campaign) error
	// Delete 削除する
	Delete(ctx context.Context, campaignID int) error
	// PutTransact トランザクションでまとめて配信に必要なデータを登録/更新する
	PutTransact(ctx context.Context, campaignData *models.Campaign) error
	// DeleteTransact トランザクションでまとめて配信に必要なデータを削除する
	DeleteTransact(ctx context.Context, ID int) error
}
