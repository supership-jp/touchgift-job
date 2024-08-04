//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	"touchgift-job-manager/domain/models"
)

type TouchPointByGroupIDCondition struct {
	GroupID int
	Limit   int
}

type TouchPointRepository interface {
	// GetTouchPointByGroupID グループIDからタッチポイントデータを取得する
	GetTouchPointByGroupID(ctx context.Context, tx Transaction, args *TouchPointByGroupIDCondition) ([]*models.TouchPoint, error)
}
