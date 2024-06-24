package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strings"
	"sync"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra/metrics"
)

var (
	metricSqsReceivedMessageTotal       = "sqs_received_message_total"
	metricSqsReceivedMessageTotalDesc   = "all received message count from sqs"
	metricSqsReceivedMessageTotalLabels = []string{"url"}

	metricSqsUnprocessableMessageTotal       = "sqs_unprocessable_message_total"
	metricSqsUnprocessableMessageTotalDesc   = "all received unprocessable message count from sqs"
	metricSqsUnprocessableMessageTotalLabels = []string{"url"}

	metricSqsDeletedMessageTotal       = "sqs_deleted_message_total"
	metricSqsDeletedMessageTotalDesc   = "all deleted message count from sqs"
	metricSqsDeletedMessageTotalLabels = []string{"url"}
)

type sqsHandler struct {
	logger                   *Logger
	svc                      *sqs.SQS
	queueURL                 *string
	visibilityTimeoutSeconds *int64
	waitTimeSeconds          *int64
	monitor                  *metrics.Monitor
}

type SQSHandler interface {
	Poll(ctx context.Context, wg *sync.WaitGroup, ch chan QueueMessage, sqsMaxMessages int64)
	UnprocessableMessage()
	OutputDeleteCliLog(message QueueMessage)
	DeleteMessage(ctx context.Context, message QueueMessage)
}

func NewSQSHandler(
	logger *Logger,
	region Region,
	queueURL *string,
	visibilityTimeoutSeconds *int64,
	waitTimeSeconds *int64,
	monitor *metrics.Monitor,
) SQSHandler {
	monitor.Metrics.AddCounter(metricSqsReceivedMessageTotal, metricSqsReceivedMessageTotalDesc, metricSqsReceivedMessageTotalLabels)
	monitor.Metrics.AddCounter(metricSqsUnprocessableMessageTotal, metricSqsUnprocessableMessageTotalDesc, metricSqsUnprocessableMessageTotalLabels)
	monitor.Metrics.AddCounter(metricSqsDeletedMessageTotal, metricSqsDeletedMessageTotalDesc, metricSqsDeletedMessageTotalLabels)

	sqsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig().WithEndpoint(config.Env.SQS.EndPoint),
		SharedConfigState: session.SharedConfigEnable,
	}))
	if config.Env.RegionFromEC2Metadata {
		sqsSession.Config.Region = region.Get()
	}
	return &sqsHandler{
		logger:                   logger,
		svc:                      sqs.New(sqsSession),
		queueURL:                 queueURL,
		visibilityTimeoutSeconds: visibilityTimeoutSeconds,
		waitTimeSeconds:          waitTimeSeconds,
		monitor:                  monitor,
	}
}

func (s *sqsHandler) Poll(ctx context.Context, wg *sync.WaitGroup, ch chan QueueMessage, sqsMaxMessages int64) {
	defer func() {
		wg.Done()
	}()
	wg.Add(1)
	s.logger.Info().Str("queue_url", *s.queueURL).Msg("Start sqs polling")
	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Str("queue_url", *s.queueURL).Msg("Stop sqs polling")
			return
		default:
			output, err := s.svc.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            s.queueURL,
				MaxNumberOfMessages: &sqsMaxMessages,
				VisibilityTimeout:   s.visibilityTimeoutSeconds,
				WaitTimeSeconds:     s.waitTimeSeconds,
			})
			if err != nil {
				s.logger.Error().Err(err).Str("queue_url", *s.queueURL).Msg("Failed to fetch sqs message")
			}
			for _, message := range output.Messages {
				var snsMessage SnsMessage
				decoder := json.NewDecoder(strings.NewReader(*message.Body))
				if err := decoder.Decode(&snsMessage); err != nil {
					s.logger.Error().Err(err).Str("queue_url", *s.queueURL).Str("body", *message.Body).Msg("Failed to parse sns message")
					s.deleteMessage(ctx, message.ReceiptHandle, message.MessageId)
				} else {
					ch <- NewMessage(message, &snsMessage)
				}
				s.monitor.Metrics.GetCounter(metricSqsReceivedMessageTotal).WithLabelValues(*s.queueURL).Inc()
			}
		}
	}
}

func (s *sqsHandler) UnprocessableMessage() {
	s.monitor.Metrics.GetCounter(metricSqsUnprocessableMessageTotal).WithLabelValues(*s.queueURL).Inc()
}

func (s *sqsHandler) OutputDeleteCliLog(message QueueMessage) {
	cli := fmt.Sprintf(`aws sqs delete-message --queue-url "%s" --receipt-handle "%s"`, *s.queueURL, *message.ReceiptHandle())
	s.logger.Warn().Str("cli", cli).Msg(`If you want to delete`)
}

func (s *sqsHandler) DeleteMessage(ctx context.Context, message QueueMessage) {
	s.deleteMessage(ctx, message.ReceiptHandle(), message.MessageID())
}

func (s *sqsHandler) deleteMessage(ctx context.Context, receiptHandle *string, messageID *string) {
	_, err := s.svc.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      s.queueURL,
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("message_id", *messageID).Str("receipt_handle", *receiptHandle).Msg("Failed to delete message.")
		return
	}
	s.monitor.Metrics.GetCounter(metricSqsDeletedMessageTotal).WithLabelValues(*s.queueURL).Inc()
}
