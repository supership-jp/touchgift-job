package injector

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra"
	"touchgift-job-manager/infra/metrics"
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
	fmt.Println(&deliveryOperationSync)

	var wg sync.WaitGroup
	initialize := func() error {
		return nil
	}
	terminate := func() error {
		wg.Wait()
		return nil
	}
	return router, initialize, terminate
}
