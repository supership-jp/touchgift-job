package injector

import (
	"touchgift-job-manager/manager/infra"
	"touchgift-job-manager/manager/interface/controllers"
)

func InjectPingController(logger *infra.Logger) controllers.HTTPHandler {
	return controllers.NewPing(
		logger,
	)
}
