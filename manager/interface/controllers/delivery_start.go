package controllers

import (
	"context"
	"sync"
	"touchgift-job-manager/infra/metrics"
	"touchgift-job-manager/usecase"
)

type DeliveryStart interface {
	StartMonitoring(ctx context.Context, wg *sync.WaitGroup)
	Close()
}

type deliveryStart struct {
	logger  usecase.Logger
	monitor *metrics.Monitor
}
