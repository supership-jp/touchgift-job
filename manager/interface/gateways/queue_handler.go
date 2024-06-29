//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package gateways

import (
	"context"
	"sync"
	"touchgift-job-manager/infra"
)

type QueueMessage = infra.QueueMessage

type QueueHandler interface {
	Poll(ctx context.Context, wg *sync.WaitGroup, ch chan infra.QueueMessage, maxMessages int64)
	UnprocessableMessage()
	OutputDeleteCliLog(message infra.QueueMessage)
	DeleteMessage(ctx context.Context, message infra.QueueMessage)
}
