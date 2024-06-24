package infra

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra/metrics"
)

func TestSQSHandler_NewSQSHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := GetLogger()
	region := NewRegion(logger)
	monitor := metrics.GetMonitor()

	t.Run("正常にsqsとの通信が確立される", func(t *testing.T) {
		handler := NewSQSHandler(
			logger, region, &config.Env.SQS.DeliveryControlQueueURL, &config.Env.SQS.VisibilityTimeoutSeconds, &config.Env.SQS.WaitTimeSeconds, monitor)
		assert.NotNil(t, handler)
	})

}
