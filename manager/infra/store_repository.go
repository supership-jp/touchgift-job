package infra

import (
	"context"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
)

type StoreRepository struct {
	logger     *Logger
	sqlHandler SQLHandler
}

func NewStoreRepository(logger *Logger, sqlHandler SQLHandler) repository.StoreRepository {
	storeRepository := StoreRepository{
		logger:     logger,
		sqlHandler: sqlHandler,
	}
	return &storeRepository
}

func (s *StoreRepository) Get(ctx context.Context, id int) (*models.Store, error) {

	query := "SELECT id, organization_code, name FROM store WHERE id = ?"
	stmt, err := s.sqlHandler.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to close statement")
		}
	}()
	dest := models.Store{}
	err = stmt.GetContext(ctx, &dest, id)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func (s *StoreRepository) Select(ctx context.Context) ([]*models.Store, error) {
	// SQLクエリ
	query := "SELECT id, organization_code, name FROM store"

	// 結果を格納するためのスライス
	var dest []*models.Store

	// SelectContextを使用してデータベースからデータを取得
	err := s.sqlHandler.Select(ctx, &dest, query)
	if err != nil {
		return nil, err // エラーがあればそのまま返す
	}

	// エラーがなければデータを返す
	return dest, nil
}
