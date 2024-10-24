//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

// TODO: 予算消化追加時に実装する

// var (
// 	metricDeliveryControlProcess     = "delivery_control_usecase_process"
// 	metricDeliveryControlProcessDesc = "delivery control usecase processing metrics"

// 	// event: delivery_control_log.event
// 	metricDeliveryControlProcessLabels = []string{"event"}
// )

// // DeliveryControl is interface
// type DeliveryControl interface {
// 	// 配信制御ログを処理する
// 	Process(ctx context.Context, current time.Time, deliveryControlLog *models.DeliveryControlLog) error
// }

// type deliveryControl struct {
// 	logger              Logger
// 	monitor             *metrics.Monitor
// 	transaction         repository.TransactionHandler
// 	campaignRepository  repository.CampaignRepository
// 	deliveryControlEvent DeliveryControlEvent
// 	deliveryEnd         DeliveryEnd
// }

// // NewDeliveryControl is function
// func NewDeliveryControl(
// 	logger Logger,
// 	monitor *metrics.Monitor,
// 	transaction repository.TransactionHandler,
// 	campaignRepository repository.CampaignRepository,
// 	deliveryControlEvent DeliveryControlEvent,
// 	deliveryEnd DeliveryEnd,
// ) DeliveryControl {
// 	instance := deliveryControl{
// 		logger:              logger,
// 		monitor:             monitor,
// 		transaction:         transaction,
// 		campaignRepository:  campaignRepository,
// 		deliveryControlEvent: deliveryControlEvent,
// 		deliveryEnd:         deliveryEnd,
// 	}
// 	monitor.Metrics.AddCounter(metricDeliveryControlProcess,
// 		metricDeliveryControlProcessDesc,
// 		metricDeliveryControlProcessLabels)
// 	return &instance
// }

// //nolint:gocognit // [21]時間あるときに修正する
// func (d *deliveryControl) Process(ctx context.Context, current time.Time, deliveryControlLog *models.DeliveryControlLog) (err error) {
// 	var tx repository.Transaction
// 	defer func() {
// 		if err != nil && tx != nil {
// 			if rerr := tx.Rollback(); rerr != nil {
// 				d.logger.Error().Err(rerr).Time("current", current).Msg("Failed to rollback")
// 			}
// 		}
// 	}()
// 	d.monitor.Metrics.
// 		GetCounter(metricDeliveryControlProcess).
// 		WithLabelValues(deliveryControlLog.Event).
// 		Inc()
// 	switch deliveryControlLog.Event {
// 	// 予算を消化した場合
// 	case "expended":
// 		tx, err = d.transaction.Begin(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		// 配信制御ログを処理する
// 		campaign, status, err := d.stopDelivery(ctx, tx, deliveryControlLog, "ended", true)
// 		if err != nil {
// 			return err
// 		}
// 		if err = tx.Commit(); err != nil {
// 			return errors.Wrap(err, "Failed to commit")
// 		}
// 		// 配信制御イベントを発行する
// 		d.deliveryControlEvent.Publish(ctx, campaign.ID, campaign.OrgCode, campaign.Status, *status, codes.DetailExpended)
// 	// 予算が不足していた場合
// 	case "shortage":
// 		tx, err = d.transaction.Begin(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		// 配信制御ログを処理する
// 		campaign, status, err := d.stopDelivery(ctx, tx, deliveryControlLog, "paused", false)
// 		if err != nil {
// 			return err
// 		}
// 		if err = tx.Commit(); err != nil {
// 			return errors.Wrap(err, "Failed to commit")
// 		}
// 		// 配信制御イベントを発行する
// 		d.deliveryControlEvent.Publish(ctx, campaign.ID, campaign.OrgCode, campaign.Status, *status, codes.DetailShortage)
// 	default:
// 		d.logger.Info().Time("current", current).Interface("delivery_control_log", deliveryControlLog).Msg("Unknown event")
// 	}
// 	return nil
// }

// // DeliveryControlLogを処理する
// func (d *deliveryControl) stopDelivery(ctx context.Context,
// 	tx repository.Transaction, deliveryControlLog *models.DeliveryControlLog,
// 	afterStatus string, budgetExpended bool) (*models.Campaign, *string, error) {
// 	// 該当する更新日時を取得
// 	campaign, err := d.campaignRepository.GetCampaignToExpendedOrShortage(ctx, tx, deliveryControlLog.CampaignID, budgetExpended)
// 	if err != nil && err != codes.ErrNoData {
// 		return nil, nil, err
// 	}
// 	if err == codes.ErrNoData {
// 		// 更新日時が取得できない場合は何もしない(できない)
// 		return nil, nil, errors.Wrap(err, "Not found campaign")
// 	}
// 	err = d.deliveryEnd.Stop(ctx, tx, campaign, afterStatus)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return campaign, &afterStatus, nil
// }
