package infra

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDeliveryData(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()
	t.Run("配信データを取得", func(t *testing.T) {
		ctx := context.Background()
		deliveryRepository := NewDeliveryRepository(logger, sqlHandler)
		deliveryData, err := deliveryRepository.GetDeliveryData(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, deliveryData, 1)
		// 配列の各要素のプロパティを検証
		assert.Equal(t, 1, deliveryData[0].ID)
		assert.Equal(t, "キャンペーン名", deliveryData[0].Name)
		assert.Equal(t, "https://example.com/video.mp4", deliveryData[0].Video)
		assert.Equal(t, "https://example.com/endcard.png", deliveryData[0].EndCard)
	})
}
