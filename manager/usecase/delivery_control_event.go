//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"encoding/json"
	"time"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/notification"

	"github.com/rs/xid"
)

type DeliveryControlEvent interface {
	Publish(ctx context.Context, CampaignID int, service string, organization string, before string, after string, detail string, deliveryType string, priceType string)
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
// CampaignID, service, organization, cacheOperation(サーバー上のキャッシュ操作), before(更新前のCampaign.status), after(更新後のCampaign.status)
func (d *deliveryControlEvent) Publish(ctx context.Context,
	CampaignID int, service string, organization string, before string, after string, detail string, deliveryType string, priceType string) {
	deliveryControl := d.createDeliveryControlLog(CampaignID, service, organization, before, after, detail, deliveryType, priceType)

	message, err := json.Marshal(deliveryControl)
	if err != nil {
		d.failedToPublishLog(deliveryControl, deliveryType, priceType, err)
	}
	messageAttributes := map[string]string{
		"event":           deliveryControl.Event,
		"cache_operation": deliveryControl.CacheOperation,
		"delivery_type":   deliveryType,
		"price_type":      priceType,
	}
	messageID, err := d.notificationHandler.Publish(ctx, string(message), messageAttributes)
	if err != nil {
		d.failedToPublishLog(deliveryControl, deliveryType, priceType, err)
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
			Str("organization", deliveryControl.Organization).
			Str("service", deliveryControl.Service).
			Int("campaign_id", deliveryControl.CampaignID).
			Str("delivery_type", deliveryType).
			Str("price_type", priceType).
			Msg("Publish delivery control event")
	}
}

func (d *deliveryControlEvent) failedToPublishLog(deliveryControl *models.DeliveryControlLog, deliveryType string, priceType string, err error) {
	// このログが出た場合はcloudwatch logsのmetric alarmでアラートを通知する
	d.logger.Error().Err(err).
		Str("trace_id", deliveryControl.TraceID).
		Str("trace_time", deliveryControl.Time).
		Int("version", deliveryControl.Version).
		Str("event", deliveryControl.Event).
		Str("cache_operation", deliveryControl.CacheOperation).
		Str("source", deliveryControl.Source).
		Str("organization", deliveryControl.Organization).
		Str("service", deliveryControl.Service).
		Int("campaign_id", deliveryControl.CampaignID).
		Str("delivery_type", deliveryType).
		Str("price_type", priceType).
		Msg("Failed to publish sns event")
}

// delivery_controlログに整形
func (d *deliveryControlEvent) createDeliveryControlLog(campaignID int,
	service string, organization string, before string, after string,
	eventDetail string, deliveryType string, priceType string) *models.DeliveryControlLog {

	event, operation := d.deliveryEvent(before, after)
	current := time.Now().Format(time.RFC3339Nano)
	return &models.DeliveryControlLog{
		TraceID:        d.createTraceID(),
		Time:           current,
		Version:        config.Env.Version,
		Event:          event,
		EventDetail:    eventDetail,
		CacheOperation: operation,
		Organization:   organization,
		Service:        service,
		Source:         "touchgift-job-manager",
		CampaignID:     campaignID,
		DeliveryType:   deliveryType,
		PriceType:      priceType,
	}
}

func (d *deliveryControlEvent) deliveryEvent(before string, after string) (string, string) {
	var event string
	var operation string
	d.logger.Info().
		Str("status_before_update", before).
		Str("status_after_update", after).Msg("Check delivery control event")
	switch {
	case before == "configured" && after == "warmup":
		event = "warmup"
		operation = "NONE"
	case before == "warmup" && after == "started":
		event = "start"
		operation = "PUT"
	case before == "resume" && after == "started":
		event = "resume"
		operation = "PUT"
	case before == "started" && after == "started":
		event = "update"
		operation = "PUT"
	case before == "stop" && after == "stopped":
		event = "stop"
		operation = "DELETE"
	case after == "paused":
		event = "pause"
		operation = "DELETE"
	case after == "ended":
		event = "end"
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
