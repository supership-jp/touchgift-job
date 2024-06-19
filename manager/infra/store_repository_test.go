package infra

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoreRepository_Get(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()

	t.Run("IDを指定して取得", func(t *testing.T) {
		ctx := context.Background()
		storeRepository := NewStoreRepository(logger, sqlHandler)
		store, err := storeRepository.Get(ctx, 1)
		if err != nil {
			return
		}
		assert.Equal(t, "東京本店", store.Name)
		assert.Equal(t, store.OrganizationCode, "ORG001")
	})
}

func TestStoreRepository_Select(t *testing.T) {
	logger := GetLogger()
	sqlHandler := NewSQLHandler(logger)
	defer sqlHandler.Close()

	t.Run("店舗取得", func(t *testing.T) {
		ctx := context.Background()
		storeRepository := NewStoreRepository(logger, sqlHandler)
		stores, err := storeRepository.Select(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, stores, 2)
	})
}
