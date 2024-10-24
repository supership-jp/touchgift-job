package infra

import (
	"context"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra/metrics"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

//TODO: 監視対象メトリクスを考える

type SNSHandler interface {
	Publish(ctx context.Context, message string, messageAttributes map[string]string, topicArn string) (*string, error)
}

type snsHandler struct {
	logger  *Logger
	svc     *sns.SNS
	monitor *metrics.Monitor
}

func NewSNSHandler(
	logger *Logger,
	region Region,
	monitor *metrics.Monitor) SNSHandler {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig().WithEndpoint(config.Env.SNS.EndPoint),
		SharedConfigState: session.SharedConfigEnable,
	}))
	if config.Env.RegionFromEC2Metadata {
		awsSession.Config.Region = region.Get()
	}

	return &snsHandler{
		logger:  logger,
		svc:     sns.New(awsSession),
		monitor: monitor,
	}
}

func (s *snsHandler) Publish(ctx context.Context, message string, messageAttributes map[string]string, topicArn string) (*string, error) {
	attributes := make(map[string]*sns.MessageAttributeValue, len(messageAttributes))
	for k, v := range messageAttributes {
		attributes[k] = &sns.MessageAttributeValue{
			StringValue: aws.String(v),
			DataType:    aws.String("String"),
		}
	}
	output, err := s.svc.PublishWithContext(ctx, &sns.PublishInput{
		TopicArn:          aws.String(topicArn),
		Message:           aws.String(message),
		MessageAttributes: attributes,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("message", message).Msg("Failed to publish message.")
		return nil, err
	}
	s.logger.Debug().Str("message_id", aws.StringValue(output.MessageId)).Msg("Publish message.")
	return output.MessageId, nil
}
