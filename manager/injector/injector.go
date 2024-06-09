package injector

import (
	"touchgift-job-manager/config"
	"touchgift-job-manager/infra"
	"touchgift-job-manager/infra/metrics"
	"touchgift-job-manager/interface/controllers"
	"touchgift-job-manager/usecase"
)

func InjectPingController(logger *infra.Logger) controllers.HTTPHandler {
	return controllers.NewPing(
		logger,
	)
}

var region infra.Region

func InjectRegion(logger *infra.Logger) infra.Region {
	if region == nil {
		region = infra.NewRegion(logger)
	}
	return region
}

func InjectSQSHandler(logger *infra.Logger, queueURL string) infra.SQSHandler {
	return infra.NewSQSHandler(
		logger,
		InjectRegion(logger),
		&queueURL,
		&config.Env.SQS.VisibilityTimeoutSeconds,
		&config.Env.SQS.WaitTimeSeconds,
		metrics.GetMonitor(),
	)
}

var sqlHandler infra.SQLHandler

func InjectSQLHandler(logger *infra.Logger) infra.SQLHandler {
	if sqlHandler == nil {
		sqlHandler = infra.NewSQLHandler(logger)
	}
	return sqlHandler
}

var deliveryOperationUsecase usecase.DeliveryOperation

func InjectDeliveryOperationUsecase(logger *infra.Logger) usecase.DeliveryOperation {
	// TODO: 依存関係のあるusecaser, repository作ってからDI
	if deliveryOperationUsecase == nil {
		deliveryOperationUsecase = usecase.NewDeliveryOperation(
			logger,
			metrics.GetMonitor(),
			InjectSQLHandler(logger),
		)
	}
	return deliveryOperationUsecase
}

func InjectDeliveryOperationSyncController(logger *infra.Logger) controllers.DeliveryOperationSync {
	subLogger := logger.With().Str("type", "delivery_sync").Logger()
	return controllers.NewDeliveryOperationSync(
		infra.NewLogger(&subLogger),
		metrics.GetMonitor(),
		InjectSQSHandler(logger, config.Env.SQS.DeliveryOperationQueueURL),
		InjectDeliveryOperationUsecase(logger),
	)
}
