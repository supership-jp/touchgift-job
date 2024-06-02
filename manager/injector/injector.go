package injector

import (
	"touchgift-job-manager/infra"
	"touchgift-job-manager/interface/controllers"
)

func InjectPingController(logger *infra.Logger) controllers.HTTPHandler {
	return controllers.NewPing(
		logger,
	)
}
