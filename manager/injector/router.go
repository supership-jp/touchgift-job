package injector

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync"
	"touchgift-job-manager/manager/config"
	"touchgift-job-manager/manager/infra"
	"touchgift-job-manager/manager/infra/metrics"
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

	monitor := metrics.GetMonitor()

	monitor.AddRoute(router, config.Env.Server.MetricsPath)

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
