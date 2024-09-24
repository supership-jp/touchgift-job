package infra

/*
import (
	"context"
	"errors"
	"testing"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestContentsHelper_GenerateContents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := GetLogger()

	// 正常な入力データの準備
	coupons := []*models.Coupon{
		{
			ID:       1,
			Name:     "Summer Sale",
			Code:     "SUMMER2024",
			ImageURL: "https://example.com/summer-sale.jpg",
			Rate:     "100",
		},
		{
			ID:       2,
			Name:     "Winter Sale",
			Code:     "WINTER2024",
			ImageURL: "https://example.com/winter-sale.jpg",
			Rate:     "200",
		},
		{
			ID:       3,
			Name:     "Spring Sale",
			Code:     "SPRING2024",
			ImageURL: "https://example.com/spring-sale.jpg",
			Rate:     "300",
		},
	}

	gimmickURL := "https://example.com/gimmick.jpg"

	contentsHelper := NewContentsHelper(logger)

	t.Run("正常な条件でのコンテンツ生成（複数クーポン）", func(t *testing.T) {
		args := &repository.GenerateContentCondition{
			Coupons:    coupons,
			GimmickURL: &gimmickURL,
		}
		result, err := contentsHelper.GenerateContents(context.Background(), args)
		assert.NoError(t, err)
		assert.Len(t, result, 1)

		// 確認: 返されたコンテンツにクーポンが3件含まれているか
		assert.Equal(t, 3, len(result[0].Coupons))

		// クーポンの詳細確認
		assert.Equal(t, 1, result[0].Coupons[0].ID)
		assert.Equal(t, "Summer Sale", result[0].Coupons[0].Name)
		assert.Equal(t, "SUMMER2024", result[0].Coupons[0].Code)
		assert.Equal(t, "https://example.com/summer-sale.jpg", result[0].Coupons[0].ImageURL)
		assert.Equal(t, "100", result[0].Coupons[0].Rate)

		assert.Equal(t, 2, result[0].Coupons[1].ID)
		assert.Equal(t, "Winter Sale", result[0].Coupons[1].Name)
		assert.Equal(t, "WINTER2024", result[0].Coupons[1].Code)
		assert.Equal(t, "https://example.com/winter-sale.jpg", result[0].Coupons[1].ImageURL)
		assert.Equal(t, "200", result[0].Coupons[1].Rate)

		assert.Equal(t, 3, result[0].Coupons[2].ID)
		assert.Equal(t, "Spring Sale", result[0].Coupons[2].Name)
		assert.Equal(t, "SPRING2024", result[0].Coupons[2].Code)
		assert.Equal(t, "https://example.com/spring-sale.jpg", result[0].Coupons[2].ImageURL)
		assert.Equal(t, "300", result[0].Coupons[2].Rate)

		// ギミックURLの確認
		assert.Equal(t, "https://example.com/gimmick.jpg", result[0].Gimmicks[0].URL)
	})

	t.Run("引数がnilの場合のエラーハンドリング", func(t *testing.T) {
		result, err := contentsHelper.GenerateContents(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("generate contents condition must not be nil"), err)
	})

	t.Run("クーポンリストが空の場合", func(t *testing.T) {
		args := &repository.GenerateContentCondition{
			Coupons:    []*models.Coupon{},
			GimmickURL: &gimmickURL,
		}
		result, err := contentsHelper.GenerateContents(context.Background(), args)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Len(t, result[0].Coupons, 0) // Ensure no coupons are present
	})

	t.Run("Gimmick URLがnilの場合", func(t *testing.T) {
		args := &repository.GenerateContentCondition{
			Coupons:    coupons,
			GimmickURL: nil,
		}
		result, err := contentsHelper.GenerateContents(context.Background(), args)
		assert.NoError(t, err)
		assert.Equal(t, "", result[0].Gimmicks[0].URL) // Gimmick URL should be empty string
	})

	t.Run("クーポン情報に不正なデータが含まれている場合", func(t *testing.T) {
		invalidCoupons := []*models.Coupon{
			{
				ID:       -1, // Invalid ID
				Name:     "",
				Code:     "",
				ImageURL: "invalid-url", // Potentially invalid URL
				Rate:     "NaN",         // Non-numeric rate
			},
		}
		args := &repository.GenerateContentCondition{
			Coupons:    invalidCoupons,
			GimmickURL: &gimmickURL,
		}
		result, err := contentsHelper.GenerateContents(context.Background(), args)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, -1, result[0].Coupons[0].ID)
		assert.Equal(t, "invalid-url", result[0].Coupons[0].ImageURL)
		assert.Equal(t, "NaN", result[0].Coupons[0].Rate)
	})
}
*/
