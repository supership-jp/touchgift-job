//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"encoding/json"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/notification"

	"github.com/rs/xid"
)

type DeliveryControlEvent interface {
	PublishCampaignEvent(ctx context.Context, CampaignID int, organization string, before string, after string, detail string)
}

type deliveryControlEvent struct {
	logger              Logger
	notificationHandler notification.NotificationHandler
}

func NewDeliveryControlEvent(
	logger Logger,
	notificationHandler notification.NotificationHandler,
) DeliveryControlEvent {
	instance := deliveryControlEvent{
		logger:              logger,
		notificationHandler: notificationHandler,
	}
	return &instance
}

// SNSへPublishを行う
// CampaignID, org_code, cacheOperation(サーバー上のキャッシュ操作), before(更新前のCampaign.status), after(更新後のCampaign.status)
func (d *deliveryControlEvent) PublishCampaignEvent(ctx context.Context,
	CampaignID int, organization string, before string, after string, detail string) {
	deliveryControl := d.createDeliveryControlLog(CampaignID, organization, before, after, detail)

	message, err := json.Marshal(deliveryControl)
	if err != nil {
		d.failedToPublishLog(deliveryControl, err)
	}
	messageAttributes := map[string]string{
		"event":           deliveryControl.Event,
		"cache_operation": deliveryControl.CacheOperation,
	}
	messageID, err := d.notificationHandler.Publish(ctx, string(message), messageAttributes)
	if err != nil {
		d.failedToPublishLog(deliveryControl, err)
	} else {
		d.logger.Info().
			Str("message_id", *messageID).
			Str("trace_id", deliveryControl.TraceID).
			Str("trace_time", deliveryControl.Time).
			Int("version", deliveryControl.Version).
			Str("event", deliveryControl.Event).
			Str("event_detail", deliveryControl.EventDetail).
			Str("cache_operation", deliveryControl.CacheOperation).
			Str("source", deliveryControl.Source).
			Str("org_code", deliveryControl.OrgCode).
			Int("campaign_id", deliveryControl.CampaignID).
			Msg("Publish delivery control event")
	}
}

func (d *deliveryControlEvent) failedToPublishLog(deliveryControl *models.DeliveryControlLog, err error) {
	// このログが出た場合はcloudwatch logsのmetric alarmでアラートを通知する
	d.logger.Error().Err(err).
		Str("trace_id", deliveryControl.TraceID).
		Str("trace_time", deliveryControl.Time).
		Int("version", deliveryControl.Version).
		Str("event", deliveryControl.Event).
		Str("cache_operation", deliveryControl.CacheOperation).
		Str("source", deliveryControl.Source).
		Str("organization", deliveryControl.OrgCode).
		Int("campaign_id", deliveryControl.CampaignID).
		Msg("Failed to publish sns event")
}

// delivery_controlログに整形
func (d *deliveryControlEvent) createDeliveryControlLog(campaignID int,
	organization string, before string, after string,
	eventDetail string) *models.DeliveryControlLog {

	event, operation := d.deliveryEvent(before, after)
	current := time.Now().Format(time.RFC3339Nano)
	return &models.DeliveryControlLog{
		TraceID:        d.createTraceID(),
		Time:           current,
		Version:        config.Env.Version,
		Event:          event,
		EventDetail:    eventDetail,
		CacheOperation: operation,
		OrgCode:        organization,
		Source:         "touchgift-job-manager",
		CampaignID:     campaignID,
	}
}

func (d *deliveryControlEvent) deliveryEvent(before string, after string) (string, string) {
	var event string
	var operation string
	d.logger.Info().
		Str("status_before_update", before).
		Str("status_after_update", after).Msg("Check delivery control event")
	switch {
	case before == codes.StatusConfigured && after == codes.StatusWarmup:
		event = codes.StatusWarmup
		operation = "NONE"
	case before == codes.StatusWarmup && after == codes.StatusStarted:
		event = codes.StatusStart
		operation = "PUT"
	case before == codes.StatusResume && after == codes.StatusStarted:
		event = codes.StatusResume
		operation = "PUT"
	case before == codes.StatusStarted && after == codes.StatusStarted:
		event = "update"
		operation = "PUT"
	case before == codes.StatusStop && after == codes.StatusStopped:
		event = codes.StatusStop
		operation = "DELETE"
	case after == codes.StatusPaused:
		event = codes.StatusPause
		operation = "DELETE"
	case after == codes.StatusEnded:
		event = codes.StatusEnd
		operation = "DELETE"
	default:
		d.logger.Warn().
			Str("status_before_update", before).
			Str("status_after_update", after).Msg("Unknown status")
	}
	return event, operation
}

func (d *deliveryControlEvent) createTraceID() string {
	result := xid.New().String()
	return result
}
