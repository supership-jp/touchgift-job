package infra

import (
	"context"
	"testing"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra/metrics"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSNSHandler_NewSNSHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := GetLogger()
	region := NewRegion(logger)
	monitor := metrics.GetMonitor()
	t.Run("正常にsnsとの通信が確立される", func(t *testing.T) {
		handler := NewSNSHandler(logger, region, monitor)
		assert.NotNil(t, handler)
	})
}

func TestSNSHandler_Publish(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := GetLogger()
	region := NewRegion(logger)
	monitor := metrics.GetMonitor()
	t.Run("snsにメッセージが正常に送れる", func(t *testing.T) {
		handler := NewSNSHandler(logger, region, monitor)
		message := "Hello, this is a test message"
		messageAttributes := map[string]string{
			"Key1": "Value1",
			"Key2": "Value2",
		}
		messageID, err := handler.Publish(ctx, message, messageAttributes, config.Env.SNS.ControlLogTopicArn)
		assert.NoError(t, err)
		assert.NotEmpty(t, messageID)
	})
}
