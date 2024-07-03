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

// TODO: メトリクスの監視項目を精査する

type CampaignDataRepository struct {
	logger          *Logger
	dynamoDBHandler *DynamoDBHandler
	tableName       *string
	monitor         *metrics.Monitor
}

func NewCampaignDataRepository(handler *DynamoDBHandler, logger *Logger, monitor *metrics.Monitor) repository.DeliveryDataCampaignRepository {
	tableName := config.Env.DynamoDB.CampaignTableName
	if len(config.Env.DynamoDB.TableNamePrefix) > 0 {
		// CIやローカル用
		tableName = config.Env.DynamoDB.TableNamePrefix + tableName
	}

	campaignDataRepository := CampaignDataRepository{
		logger:          logger,
		dynamoDBHandler: handler,
		tableName:       &tableName,
		monitor:         monitor,
	}

	return &campaignDataRepository
}

func (c *CampaignDataRepository) Get(ctx context.Context, id *string) (*models.DeliveryDataCampaign, error) {
	result, err := c.dynamoDBHandler.Svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: c.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: id,
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
	item := models.DeliveryDataCampaign{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Put キャンペーン配信データの登録/更新を行う
// TODO: メトリクス項目を考える(成功時、失敗時)
func (c *CampaignDataRepository) Put(ctx context.Context, updateData *models.DeliveryDataCampaign) error {

	item, err := dynamodbattribute.MarshalMap(updateData)
	if err != nil {
		return err
	}
	_, err = c.dynamoDBHandler.Svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName:    c.tableName,
		Item:         item,
		ReturnValues: aws.String("NONE"),
	})
	if err != nil {
		return err
	}
	return nil

}

func (c *CampaignDataRepository) PutAll(ctx context.Context, updateData *[]models.DeliveryDataCampaign) error {
	for i := range *updateData {
		updateData := (*updateData)[i]
		err := c.Put(ctx, &updateData)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (c *CampaignDataRepository) Delete(ctx context.Context, campaignID *string) error {
	_, err := c.dynamoDBHandler.Svc.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: c.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: campaignID,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *CampaignDataRepository) DeleteAll(ctx context.Context, deleteDatas *[]models.DeliveryDataCampaign) error {
	for i := range *deleteDatas {
		deleteData := (*deleteDatas)[i]
		ID := &deleteData.ID
		err := c.Delete(ctx, ID)
		if err != nil {
			return err
		}
	}
	return nil
}
