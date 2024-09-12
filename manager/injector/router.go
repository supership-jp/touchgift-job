package injector

import (
	"context"
	"sync"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra"
	"touchgift-job-manager/infra/metrics"

	"github.com/gin-gonic/gin"
)

type Initialize func() error

type Terminate func() error

func Route(router *gin.Engine) *gin.Engine {
	logger := infra.GetLogger()
	router.GET("/ping", func(c *gin.Context) {
		InjectPingController(logger).Handler(infra.NewContext(c))
	})
	return router
}

func AdminRoute(ctx context.Context, router *gin.Engine) (*gin.Engine, Initialize, Terminate) {
	logger := infra.GetLogger()
	monitor := metrics.GetMonitor()

	monitor.AddRoute(router, config.Env.Server.MetricsPath)

	deliveryOperationSync := InjectDeliveryOperationSyncController(logger)
	deliveryStart := InjectDeliveryStartController(logger)
	deliveryEnd := InjectDeliveryEndController(logger)
	// deliveryControlSync := InjectDeliveryControlSyncController(logger)

	var wg sync.WaitGroup
	initialize := func() error {
		deliveryOperationSync.Start(ctx, &wg)
		go deliveryStart.StartMonitoring(ctx, &wg)
		go deliveryEnd.StartMonitoring(ctx, &wg)
		return nil
	}
	terminate := func() error {
		wg.Wait()
		deliveryOperationSync.Close()
		deliveryStart.Close()
		deliveryEnd.Close()
		return nil
	}
	return router, initialize, terminate
}
