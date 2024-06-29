package infra

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra/metrics"
)

//TODO: 監視対象メトリクスを考える

type SNSHandler interface {
	Publish(ctx context.Context, message string, messageAttributes map[string]string) (*string, error)
}

type snsHandler struct {
	logger   *Logger
	svc      *sns.SNS
	topicArn string
	monitor  *metrics.Monitor
}

func NewSNSHandler(
	logger *Logger,
	region Region,
	topicArn string,
	monitor *metrics.Monitor) SNSHandler {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig().WithEndpoint(config.Env.SNS.EndPoint),
		SharedConfigState: session.SharedConfigEnable,
	}))
	if config.Env.RegionFromEC2Metadata {
		awsSession.Config.Region = region.Get()
	}

	return &snsHandler{
		logger:   logger,
		svc:      sns.New(awsSession),
		topicArn: topicArn,
		monitor:  monitor,
	}
}

func (s *snsHandler) Publish(ctx context.Context, message string, messageAttributes map[string]string) (*string, error) {
	attributes := make(map[string]*sns.MessageAttributeValue, len(messageAttributes))
	for k, v := range messageAttributes {
		attributes[k] = &sns.MessageAttributeValue{
			StringValue: aws.String(v),
			DataType:    aws.String("String"),
		}
	}
	output, err := s.svc.PublishWithContext(ctx, &sns.PublishInput{
		TopicArn:          aws.String(s.topicArn),
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
