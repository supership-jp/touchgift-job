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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	metricDynamodbPutTotal       = "dynamodb_put_total"
	metricDynamodbPutTotalDesc   = "put total count from dynamodb"
	metricDynamodbPutTotalLabels = []string{"table_name", "kind"}

	metricDynamodbUpdateTotal       = "dynamodb_update_total"
	metricDynamodbUpdateTotalDesc   = "update total count from dynamodb"
	metricDynamodbUpdateTotalLabels = []string{"table_name", "kind"}

	metricDynamodbDeleteTotal       = "dynamodb_delete_total"
	metricDynamodbDeleteTotalDesc   = "get total count from dynamodb"
	metricDynamodbDeleteTotalLabels = []string{"table_name", "kind"}
)

// DeliveryDataCreativeRepository is struvt
type DeliveryDataCreativeRepository struct {
	logger          *Logger
	dynamoDBHandler *DynamoDBHandler
	tableName       *string
	monitor         *metrics.Monitor
}

// NewDeliveryDataCreativeRepository is function
func NewDeliveryDataCreativeRepository(handler *DynamoDBHandler, logger *Logger, monitor *metrics.Monitor) repository.DeliveryDataCreativeRepository {
	tableName := config.Env.DynamoDB.CreativeTableName
	if len(config.Env.DynamoDB.TableNamePrefix) > 0 {
		// CIやローカル用
		tableName = config.Env.DynamoDB.TableNamePrefix + tableName
	}

	monitor.Metrics.AddCounter(metricDynamodbPutTotal, metricDynamodbPutTotalDesc, metricDynamodbPutTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbDeleteTotal, metricDynamodbDeleteTotalDesc, metricDynamodbDeleteTotalLabels)
	monitor.Metrics.AddCounter(metricDynamodbUpdateTotal, metricDynamodbUpdateTotalDesc, metricDynamodbUpdateTotalLabels)

	DeliveryDataCreativeRepository := DeliveryDataCreativeRepository{
		logger:          logger,
		dynamoDBHandler: handler,
		tableName:       &tableName,
		monitor:         monitor,
	}
	return &DeliveryDataCreativeRepository
}

// Get is function
func (r *DeliveryDataCreativeRepository) Get(ctx context.Context, id *string) (*models.DeliveryDataCreative, error) {
	result, err := r.dynamoDBHandler.Svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
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
	item := models.DeliveryDataCreative{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Put is function
func (r *DeliveryDataCreativeRepository) Put(ctx context.Context, updateData *models.DeliveryDataCreative) error {
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
func (r *DeliveryDataCreativeRepository) PutAll(ctx context.Context, updateDatas *[]models.DeliveryDataCreative) error {
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
func (r *DeliveryDataCreativeRepository) Delete(ctx context.Context, id *string) error {
	defer func() {
		r.monitor.Metrics.GetCounter(metricDynamodbDeleteTotal).WithLabelValues(*r.tableName, "success").Inc()
	}()
	_, err := r.dynamoDBHandler.Svc.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
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
func (r *DeliveryDataCreativeRepository) DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryDataCreative) error {
	for i := range *deleteDatas {
		deleteData := (*deleteDatas)[i]
		ID := strconv.Itoa(deleteData.ID)
		err := r.Delete(ctx, &ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateTTL is function
func (r *DeliveryDataCreativeRepository) UpdateTTL(ctx context.Context, id string, ttl int64) error {
	defer func() {
		r.monitor.Metrics.GetCounter(metricDynamodbUpdateTotal).WithLabelValues(*r.tableName, "success").Inc()
	}()

	_, err := r.dynamoDBHandler.Svc.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#ttl": aws.String("ttl"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":ttl": {
				N: aws.String(strconv.FormatInt(ttl, 10)),
			},
		},
		UpdateExpression:    aws.String("SET #ttl = :ttl"),
		ConditionExpression: aws.String("attribute_exists(#ttl)"),
		ReturnValues:        aws.String("NONE"),
	})
	if err != nil {
		if awserr, ok := err.(awserr.RequestFailure); ok && awserr.Code() == "ConditionalCheckFailedException" {
			r.monitor.Metrics.GetCounter(metricDynamodbUpdateTotal).WithLabelValues(*r.tableName, "failed").Inc()
			r.logger.Error().Err(err).Msg("Condition mismatch.")
			return codes.ErrConditionFailed
		}
		r.monitor.Metrics.GetCounter(metricDynamodbUpdateTotal).WithLabelValues(*r.tableName, "error").Inc()
		r.logger.Error().Err(err).Msg("Failed to connect.")
		return err
	}
	return nil
}
