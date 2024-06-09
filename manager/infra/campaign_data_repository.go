package infra

import (
	"context"
	"fmt"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"
)

var (
	// Dynamoに対して行われたPUT操作の総数
	metricDynamodbPutTotal = "dynamodb_put_total"
	// 複数のアイテムへのPUT操作の総数
	metricDynamodbPutTotalDesc = "put total count from dynamodb"
	// DynamoDB操作の種類とテーブル名を識別するラベル
	metricDynamodbPutTotalLabels = []string{"table_name", "kind"}

	// トランザクションを伴うPUT操作の総数
	metricDynamodbPutTransactTotal = "dynamodb_put_transact_total"
	// 複数のアイテムや操作を一つの処理で行うPUT操作の総数
	metricDynamodbPutTransactTotalDesc = "put transact total count from dynamodb"
	// トランザクションPUT操作が行われたDynamoDBのテーブル名と操作の種類を識別するラベル
	metricDynamodbPutTransactTotalLabels = []string{"table_name", "kind"}

	// トランザクションを伴うDELETE操作の総数
	metricDynamodbDeleteTransactTotal = "dynamodb_delete_transact_total"
	// 複数のアイテムの削除や他の操作を一つの処理で行うDELETE操作の総数
	metricDynamodbDeleteTransactTotalDesc = "delete transact total count from dynamodb"
	// トランザクションDELETE操作が行われたDynamoDBのテーブル名と操作の種類を識別するラベル
	metricDynamodbDeleteTransactTotalLabels = []string{"table_name", "kind"}
)

type CampaignDataRepository struct {
	logger          *Logger
	dynamoDBHandler *DynamoDBHandler
	tableName       *string
	monitor         *metrics.Monitor
}

func NewDeliveryDataRepository(handler *DynamoDBHandler, logger *Logger, monitor *metrics.Monitor) repository.DeliveryDataRepository {
	tableName := config.Env.DynamoDB.DeliveryTableName
	if len(config.Env.DynamoDB.TableNamePrefix) > 0 {
		// CIやローカル用
		tableName = config.Env.DynamoDB.TableNamePrefix + tableName
	}

	monitor.Metrics.AddCounter(metricDynamodbPutTotal, metricDynamodbPutTotalDesc, metricDynamodbPutTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbPutTransactTotal, metricDynamodbPutTransactTotalDesc, metricDynamodbPutTransactTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbDeleteTransactTotal, metricDynamodbDeleteTransactTotalDesc, metricDynamodbDeleteTransactTotalLabels)

	deliveryDataRepository := CampaignDataRepository{
		logger:          logger,
		dynamoDBHandler: handler,
		tableName:       &tableName,
		monitor:         monitor,
	}
	fmt.Print(deliveryDataRepository)
	return nil
	//return &deliveryDataRepository
}

func Get(ctx context.Context, campaignID int) (*models.CampaignData, error) {

}
