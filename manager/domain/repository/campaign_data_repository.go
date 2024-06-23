package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type CampaignDataRepository interface {
	// 取得する
	Get(ctx context.Context, id *string) (*models.Campaign, error)
	//	登録/更新する
	Put(ctx context.Context, updateData *models.Campaign) error
	// まとめて登録更新する
	PutAll(ctx context.Context, updateData *[]models.Campaign) error
	// 削除する
	Delete(ctx context.Context, campaignID *string) error
	// まとめて削除する
	DeleteAll(ctx context.Context, deleteDatas *[]models.Campaign) error
}
