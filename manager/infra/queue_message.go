//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE

package infra

import "github.com/aws/aws-sdk-go/service/sqs"

// SnsMessage
// SNS -> SQSのため、SNSメッセージの形式で受け取る
type SnsMessage struct {
	Type             string `json:"Type"`
	MessageID        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Subject          string `json:"Subject"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion int    `json:"SignatureVersion,string"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

// QueueMessage is interface
type QueueMessage interface {
	Message() *string
	MessageID() *string
	ReceiptHandle() *string
}

type queueMessage struct {
	sqsMessage *sqs.Message
	snsMessage *SnsMessage
}

func NewMessage(sqsMessage *sqs.Message, snsMessage *SnsMessage) QueueMessage {
	return &queueMessage{sqsMessage, snsMessage}
}

func (q *queueMessage) Message() *string {
	return &q.snsMessage.Message
}

func (q *queueMessage) MessageID() *string {
	return q.sqsMessage.MessageId
}

func (q *queueMessage) ReceiptHandle() *string {
	return q.sqsMessage.ReceiptHandle
}
