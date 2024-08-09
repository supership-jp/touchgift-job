package infra

import (
	"context"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/repository"
	"touchgift-job-manager/infra/metrics"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DeliveryTouchPointRepository is struvt
type DeliveryTouchPointRepository struct {
	logger          *Logger
	dynamoDBHandler *DynamoDBHandler
	tableName       *string
	monitor         *metrics.Monitor
}

// NewDeliveryTouchPointRepository is function
func NewDeliveryDataTouchPointRepository(handler *DynamoDBHandler, logger *Logger, monitor *metrics.Monitor) repository.DeliveryDataTouchPointRepository {
	tableName := config.Env.DynamoDB.TouchPointTableName
	if len(config.Env.DynamoDB.TableNamePrefix) > 0 {
		// CIやローカル用
		tableName = config.Env.DynamoDB.TableNamePrefix + tableName
	}

	monitor.Metrics.AddCounter(metricDynamodbPutTotal, metricDynamodbPutTotalDesc, metricDynamodbPutTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbDeleteTotal, metricDynamodbDeleteTotalDesc, metricDynamodbDeleteTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbUpdateTotal, metricDynamodbUpdateTotalDesc, metricDynamodbUpdateTotalLabels)

	DeliveryTouchPointRepository := DeliveryTouchPointRepository{
		logger:          logger,
		dynamoDBHandler: handler,
		tableName:       &tableName,
		monitor:         monitor,
	}
	return &DeliveryTouchPointRepository
}

// Get is function
func (r *DeliveryTouchPointRepository) Get(ctx context.Context, id *string) (*models.DeliveryTouchPoint, error) {
	result, err := r.dynamoDBHandler.Svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"group_id": {
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
	item := models.DeliveryTouchPoint{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Put is function
func (r *DeliveryTouchPointRepository) Put(ctx context.Context, updateData *models.DeliveryTouchPoint) error {
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
func (r *DeliveryTouchPointRepository) PutAll(ctx context.Context, updateDatas *[]models.DeliveryTouchPoint) error {
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
func (r *DeliveryTouchPointRepository) Delete(ctx context.Context, id *string) error {
	defer func() {
		r.monitor.Metrics.GetCounter(metricDynamodbDeleteTotal).WithLabelValues(*r.tableName, "success").Inc()
	}()
	_, err := r.dynamoDBHandler.Svc.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"group_id": {
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
func (r *DeliveryTouchPointRepository) DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryTouchPoint) error {
	for i := range *deleteDatas {
		deleteData := (*deleteDatas)[i]
		ID := &deleteData.TouchPointID
		err := r.Delete(ctx, ID)
		if err != nil {
			return err
		}
	}
	return nil
}
