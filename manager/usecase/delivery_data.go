//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import "context"

type DeliveryData interface {
	// Put 配信データを登録する
	Put(ctx context.Context) error
	// Delete 配信データを1件削除
	Delete(ctx context.Context, campaignID int) error
}

type deliveryData struct {
	logger Logger
}

func (d *deliveryData) Put(ctx context.Context) error {

	// TODO:[Dynamo]TouchPointデータの作成

	// TODO:[Dynamo]Campaignデータの作成

	// TODO:[Dynamo]Contentsデータの作成

	return nil
}

func (d *deliveryData) Delete(ctx context.Context, campaignID int) error {
	return nil
}
