package infra

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"touchgift-job-manager/config"
)

func TestNewDynamoDBHandlerWithLocalStack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックの作成
	mockLogger := GetLogger()
	mockRegion := NewRegion(mockLogger)

	// DynamoDBハンドラの作成
	handler := NewDynamoDBHandler(mockLogger, mockRegion)

	// テストの検証
	assert.NotNil(t, handler, "handlerがnilではありません。")
	assert.NotNil(t, handler.Svc, "Dynamoクライアントが正常に動作しています")
	assert.Equal(t, config.Env.DynamoDB.EndPoint, handler.Svc.Endpoint, "エンドポイントが正しく設定されています")
}
