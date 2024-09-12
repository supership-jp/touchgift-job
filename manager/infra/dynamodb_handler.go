package infra

import (
	"touchgift-job-manager/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBHandler is struct
type DynamoDBHandler struct {
	Svc    *dynamodb.DynamoDB
	logger *Logger
	region Region
}

// NewDynamoDBHandler is function
func NewDynamoDBHandler(logger *Logger, region Region) *DynamoDBHandler {
	dynamoSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig().WithEndpoint(config.Env.DynamoDB.EndPoint),
		SharedConfigState: session.SharedConfigEnable,
	}))
	if config.Env.RegionFromEC2Metadata {
		dynamoSession.Config.Region = region.Get()
	}
	handler := DynamoDBHandler{
		logger: logger,
		Svc:    dynamodb.New(dynamoSession),
		region: region,
	}
	return &handler
}
