package infra

import (
	"context"
	"strconv"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DeliveryContentRepository is struvt
type DeliveryContentRepository struct {
	logger          *Logger
	dynamoDBHandler *DynamoDBHandler
	tableName       *string
	monitor         *metrics.Monitor
}

// NewDeliveryContentRepository is function
func NewDeliveryDataContentRepository(handler *DynamoDBHandler, logger *Logger, monitor *metrics.Monitor) repository.DeliveryDataContentRepository {
	tableName := config.Env.DynamoDB.ContentTableName
	if len(config.Env.DynamoDB.TableNamePrefix) > 0 {
		// CIやローカル用
		tableName = config.Env.DynamoDB.TableNamePrefix + tableName
	}

	monitor.Metrics.AddCounter(metricDynamodbPutTotal, metricDynamodbPutTotalDesc, metricDynamodbPutTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbDeleteTotal, metricDynamodbDeleteTotalDesc, metricDynamodbDeleteTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbUpdateTotal, metricDynamodbUpdateTotalDesc, metricDynamodbUpdateTotalLabels)

	DeliveryContentRepository := DeliveryContentRepository{
		logger:          logger,
		dynamoDBHandler: handler,
		tableName:       &tableName,
		monitor:         monitor,
	}
	return &DeliveryContentRepository
}

// Get is function
func (r *DeliveryContentRepository) Get(ctx context.Context, id *string) (*models.DeliveryDataContent, error) {
	result, err := r.dynamoDBHandler.Svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"campaign_id": {
				N: id,
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, codes.ErrNoData
	}
	item := models.DeliveryDataContent{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Put is function
func (r *DeliveryContentRepository) Put(ctx context.Context, updateData *models.DeliveryDataContent) error {
	defer func() {
		r.monitor.Metrics.GetCounter(metricDynamodbPutTotal).WithLabelValues(*r.tableName, "success").Inc()
	}()

	item, err := dynamodbattribute.MarshalMap(updateData)
	if err != nil {
		r.monitor.Metrics.GetCounter(metricDynamodbPutTotal).WithLabelValues(*r.tableName, "marshal_error").Inc()
		return err
	}
	_, err = r.dynamoDBHandler.Svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName:    r.tableName,
		Item:         item,
		ReturnValues: aws.String("NONE"),
	})
	if err != nil {
		r.monitor.Metrics.GetCounter(metricDynamodbPutTotal).WithLabelValues(*r.tableName, "error").Inc()
		return err
	}
	return nil
}

// PutAll is function
func (r *DeliveryContentRepository) PutAll(ctx context.Context, updateDatas *[]models.DeliveryDataContent) error {
	for i := range *updateDatas {
		updateData := (*updateDatas)[i]
		err := r.Put(ctx, &updateData)
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete is function
func (r *DeliveryContentRepository) Delete(ctx context.Context, id *string) error {
	defer func() {
		r.monitor.Metrics.GetCounter(metricDynamodbDeleteTotal).WithLabelValues(*r.tableName, "success").Inc()
	}()
	_, err := r.dynamoDBHandler.Svc.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"campaign_id": {
				N: id,
			},
		},
	})
	if err != nil {
		r.monitor.Metrics.GetCounter(metricDynamodbDeleteTotal).WithLabelValues(*r.tableName, "error").Inc()
		return err
	}
	return nil
}

// DeleteAll is function
func (r *DeliveryContentRepository) DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryDataContent) error {
	for i := range *deleteDatas {
		deleteData := (*deleteDatas)[i]
		ID := deleteData.CampaignID
		IDStr := strconv.Itoa(ID)
		err := r.Delete(ctx, &IDStr)
		if err != nil {
			return err
		}
	}
	return nil
}
