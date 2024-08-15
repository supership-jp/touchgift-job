package injector

import (
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/notification"
	"touchgift-job-manager/domain/repository"
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

var timer usecase.Timer

func InjectTimer(logger *infra.Logger) usecase.Timer {
	if timer == nil {
		timer = usecase.NewTimer(
			logger,
		)
	}
	return timer
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

var notifactionHandler notification.NotificationHandler

func InjectSNSHandler(logger *infra.Logger, topicArn string) notification.NotificationHandler {
	if notifactionHandler == nil {
		notifactionHandler = infra.NewSNSHandler(
			logger,
			InjectRegion(logger),
			topicArn,
			metrics.GetMonitor(),
		)
	}
	return notifactionHandler
}

var creativeUsecase usecase.Creative

func InjectCreativeUsecase(logger *infra.Logger) usecase.Creative {
	if creativeUsecase == nil {
		creativeUsecase = usecase.NewCreative(
			logger,
			InjectCreativeDataRepository(logger),
			InjectCreativeRepository(logger),
		)
	}
	return creativeUsecase
}

var deliveryStartUsecase usecase.DeliveryStart

func InjectDeliveryStartUsecase(logger *infra.Logger) usecase.DeliveryStart {
	if deliveryStartUsecase == nil {
		deliveryStartUsecase = usecase.NewDeliveryStart(
			logger,
			metrics.GetMonitor(),
			&config.Env.DeliveryStart,
			&config.Env.DeliveryStartUsecase,
			InjectSQLHandler(logger),
			InjectTimer(logger),
			InjectDeliveryControlEventUsecase(logger),
			InjectCampaignRepository(logger),
			InjectCreativeRepository(logger),
			InjectContentRepository(logger),
			InjectTouchPointRepository(logger),
			InjectCampaignDataRepository(logger),
			InjectContentDataRepository(logger),
			InjectCreativeDataRepository(logger),
			InjectTouchPointDataRepository(logger),
		)
	}
	return deliveryStartUsecase
}

var deliveryEndUsecase usecase.DeliveryEnd

func InjectDeliveryEndUsecase(logger *infra.Logger) usecase.DeliveryEnd {
	if deliveryEndUsecase == nil {
		deliveryEndUsecase = usecase.NewDeliveryEnd(
			logger,
			metrics.GetMonitor(),
			&config.Env.DeliveryEnd,
			&config.Env.DeliveryEndUsecase,
			InjectSQLHandler(logger),
			InjectTimer(logger),
			InjectDeliveryControlEventUsecase(logger),
			InjectCampaignRepository(logger),
			InjectCampaignDataRepository(logger),
			InjectContentDataRepository(logger),
			InjectTouchPointDataRepository(logger),
		)
	}
	return deliveryEndUsecase
}

var deliveryOperationUsecase usecase.DeliveryOperation

func InjectDeliveryOperationUsecase(logger *infra.Logger) usecase.DeliveryOperation {
	// TODO: 依存関係のあるusecaser, repository作ってからDI
	if deliveryOperationUsecase == nil {
		deliveryOperationUsecase = usecase.NewDeliveryOperation(
			logger,
			metrics.GetMonitor(),
			InjectSQLHandler(logger),
			InjectCampaignRepository(logger),
			InjectCampaignDataRepository(logger),
			InjectCreativeUsecase(logger),
			InjectDeliveryStartUsecase(logger),
			InjectDeliveryEndUsecase(logger),
			InjectDeliveryControlEventUsecase(logger),
		)
	}
	return deliveryOperationUsecase
}

var deliveryControlEventUsecase usecase.DeliveryControlEvent

func InjectDeliveryControlEventUsecase(logger *infra.Logger) usecase.DeliveryControlEvent {
	if deliveryControlEventUsecase == nil {
		deliveryControlEventUsecase = usecase.NewDeliveryControlEvent(
			logger,
			InjectSNSHandler(logger, config.Env.SNS.ControlLogTopicArn),
		)
	}
	return deliveryControlEventUsecase
}

var appTicker controllers.AppTicker

func InjectAppTicker() controllers.AppTicker {
	if appTicker == nil {
		appTicker = controllers.NewAppTicker()
	}
	return appTicker
}

var deliveryOperationSyncController controllers.DeliveryOperationSync

func InjectDeliveryOperationSyncController(logger *infra.Logger) controllers.DeliveryOperationSync {
	subLogger := logger.With().Str("type", "delivery_sync").Logger()
	if deliveryOperationSyncController == nil {
		deliveryOperationSyncController = controllers.NewDeliveryOperationSync(
			infra.NewLogger(&subLogger),
			metrics.GetMonitor(),
			InjectSQSHandler(logger, config.Env.SQS.DeliveryOperationQueueURL),
			InjectDeliveryOperationUsecase(logger),
		)
	}
	return deliveryOperationSyncController
}

var deliveryStartController controllers.DeliveryStart

func InjectDeliveryStartController(logger *infra.Logger) controllers.DeliveryStart {
	subLogger := logger.With().Str("type", "delivery_start").Logger()
	if deliveryStartController == nil {
		return controllers.NewDeliveryStart(
			infra.NewLogger(&subLogger),
			metrics.GetMonitor(),
			&config.Env.DeliveryStart,
			InjectAppTicker(),
			InjectSQLHandler(logger),
			InjectDeliveryStartUsecase(logger),
			InjectDeliveryControlEventUsecase(logger),
		)
	}
	return deliveryStartController
}

var deliveryEndController controllers.DeliveryEnd

func InjectDeliveryEndController(logger *infra.Logger) controllers.DeliveryEnd {
	subLogger := logger.With().Str("type", "delivery_end").Logger()
	if deliveryEndController == nil {
		return controllers.NewDeliveryEnd(
			infra.NewLogger(&subLogger),
			metrics.GetMonitor(),
			&config.Env.DeliveryEnd,
			InjectAppTicker(),
			InjectSQLHandler(logger),
			InjectDeliveryEndUsecase(logger),
		)
	}
	return deliveryEndController
}

// func InjectDeliveryControlSyncController(logger *infra.Logger) controllers.DeliveryControlSync {
// 	subLogger := logger.With().Str("type", "delivery_control_event").Logger()
// 	return controllers.NewDeliveryControlSync(
// 		infra.NewLogger(&subLogger),
// 		InjectSQSHandler(logger, config.Env.SQS.DeliveryControlQueueURL),
// 		InjectDeliveryControlEventUsecase(logger),
// 	)
// }

var campaignRepository repository.CampaignRepository

func InjectCampaignRepository(logger *infra.Logger) repository.CampaignRepository {
	if campaignRepository == nil {
		campaignRepository = infra.NewCampaignRepository(
			logger,
			InjectSQLHandler(logger),
		)
	}
	return campaignRepository
}

var creativeRepository repository.CreativeRepository

func InjectCreativeRepository(logger *infra.Logger) repository.CreativeRepository {
	if creativeRepository == nil {
		creativeRepository = infra.NewCreativeRepository(
			logger,
			InjectSQLHandler(logger),
		)
	}
	return creativeRepository
}

var contentRepository repository.ContentRepository

func InjectContentRepository(logger *infra.Logger) repository.ContentRepository {
	if contentRepository == nil {
		contentRepository = infra.NewContentRepository(
			logger,
			InjectSQLHandler(logger),
		)
	}
	return contentRepository
}

var touchPointRepository repository.TouchPointRepository

func InjectTouchPointRepository(logger *infra.Logger) repository.TouchPointRepository {
	if touchPointRepository == nil {
		touchPointRepository = infra.NewTouchPointRepository(
			logger,
			InjectSQLHandler(logger),
		)
	}
	return touchPointRepository
}

var campaignDataRepository repository.DeliveryDataCampaignRepository

func InjectCampaignDataRepository(logger *infra.Logger) repository.DeliveryDataCampaignRepository {
	if campaignDataRepository == nil {
		campaignDataRepository = infra.NewCampaignDataRepository(
			infra.NewDynamoDBHandler(logger, InjectRegion(logger)),
			logger,
			metrics.GetMonitor(),
		)
	}
	return campaignDataRepository
}

var creativeDataRepository repository.DeliveryDataCreativeRepository

func InjectCreativeDataRepository(logger *infra.Logger) repository.DeliveryDataCreativeRepository {
	if creativeDataRepository == nil {
		creativeDataRepository = infra.NewDeliveryDataCreativeRepository(
			infra.NewDynamoDBHandler(logger, InjectRegion(logger)),
			logger,
			metrics.GetMonitor(),
		)
	}
	return creativeDataRepository
}

var touchPointDataRepository repository.DeliveryDataTouchPointRepository

func InjectTouchPointDataRepository(logger *infra.Logger) repository.DeliveryDataTouchPointRepository {
	if touchPointDataRepository == nil {
		touchPointDataRepository = infra.NewDeliveryDataTouchPointRepository(
			infra.NewDynamoDBHandler(logger, InjectRegion(logger)),
			logger,
			metrics.GetMonitor(),
		)
	}
	return touchPointDataRepository
}

var contentDataRepository repository.DeliveryDataContentRepository

func InjectContentDataRepository(logger *infra.Logger) repository.DeliveryDataContentRepository {
	if contentDataRepository == nil {
		contentDataRepository = infra.NewDeliveryDataContentRepository(
			infra.NewDynamoDBHandler(logger, InjectRegion(logger)),
			logger,
			metrics.GetMonitor(),
		)
	}
	return contentDataRepository
}
