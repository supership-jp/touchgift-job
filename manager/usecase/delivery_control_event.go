//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"touchgift-job-manager/codes"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/models"
	"touchgift-job-manager/domain/notification"

	"github.com/rs/xid"
)

type DeliveryControlEvent interface {
	PublishCampaignEvent(ctx context.Context, CampaignID int, groupID int, organization string, before string, after string, detail string)
	PublishCreativeEvent(ctx context.Context, creative *models.DeliveryDataCreative, organization string, action string)
	PublishDeliveryEvent(ctx context.Context, id string, groupID int, storeID string, campaignID int, organization string, action string)
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

// サーバーのCampaignキャッシュ更新のためSNSへPublishを行う
// CampaignID, org_code, cacheOperation(サーバー上のキャッシュ操作), before(更新前のCampaign.status), after(更新後のCampaign.status)
func (d *deliveryControlEvent) PublishCampaignEvent(ctx context.Context,
	CampaignID int, groupID int, organization string, before string, after string, detail string) {
	deliveryControl := d.createCampaignCacheLog(CampaignID, groupID, organization, before, after, detail)

	message, err := json.Marshal(deliveryControl)
	if err != nil {
		d.failedToPublishLog(deliveryControl, err)
	}
	messageAttributes := map[string]string{
		"event":  deliveryControl.Event,
		"action": deliveryControl.Action,
	}
	messageID, err := d.notificationHandler.Publish(ctx, string(message), messageAttributes, config.Env.SNS.ControlLogTopicArn)
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
			Str("action", deliveryControl.Action).
			Str("source", deliveryControl.Source).
			Str("org_code", deliveryControl.OrgCode).
			Str("campaign_id", deliveryControl.ID).
			Str("group_id", deliveryControl.GroupID).
			Msg("Publish campaign cache event")
	}
}

// サーバーのCreativeキャッシュ更新のためSNSへPublishを行う
func (d *deliveryControlEvent) PublishCreativeEvent(ctx context.Context,
	creative *models.DeliveryDataCreative, organization string, action string) {
	deliveryControl := d.createCreativeEventLog(creative, organization, action)

	message, err := json.Marshal(deliveryControl)
	if err != nil {
		d.logger.Error().Err(err).
			Str("creative_id", deliveryControl.ID).
			Str("action", deliveryControl.Action).
			Msg("Failed to marshal json")
	}
	messageAttributes := map[string]string{
		"action": deliveryControl.Action,
	}
	messageID, err := d.notificationHandler.Publish(ctx, string(message), messageAttributes, config.Env.SNS.CreativeCacheTopicArn)
	if err != nil {
		d.logger.Error().Err(err).
			Str("creative_id", deliveryControl.ID).
			Str("action", deliveryControl.Action).
			Msg("Failed to creative cache sns publish")
	} else {
		d.logger.Info().
			Str("message_id", *messageID).
			Str("action", deliveryControl.Action).
			Str("creative_id", deliveryControl.ID).
			Msg("Publish creative cache event")
	}
}

// サーバーのTouchpointキャッシュ更新のためSNSへPublishを行う
func (d *deliveryControlEvent) PublishDeliveryEvent(ctx context.Context,
	id string, groupID int, storeID string, campaignID int, organization string, operation string) {
	deliveryControl := d.createDeliveryEventLog(id, groupID, storeID, organization, campaignID, operation)

	message, err := json.Marshal(deliveryControl)
	if err != nil {
		d.logger.Error().Err(err).
			Str("touchpoint_id", deliveryControl.ID).
			Str("action", deliveryControl.Action).
			Msg("Failed to marshal json")
	}
	messageAttributes := map[string]string{
		"action": deliveryControl.Action,
	}
	messageID, err := d.notificationHandler.Publish(ctx, string(message), messageAttributes, config.Env.SNS.DeliveryCacheTopicArn)
	if err != nil {
		d.logger.Error().Err(err).
			Str("touchpoint_id", deliveryControl.ID).
			Str("action", deliveryControl.Action).
			Msg("Failed to delivery cache sns publish")
	} else {
		d.logger.Info().
			Str("message_id", *messageID).
			Str("action", deliveryControl.Action).
			Str("org_code", deliveryControl.OrgCode).
			Int("campaign_id", deliveryControl.CampaignID).
			Int("group_id", deliveryControl.GroupID).
			Str("store_id", deliveryControl.StoreID).
			Msg("Publish delivery control event")
	}
}

func (d *deliveryControlEvent) failedToPublishLog(deliveryControl *models.CampaignCacheLog, err error) {
	// このログが出た場合はcloudwatch logsのmetric alarmでアラートを通知する
	d.logger.Error().Err(err).
		Str("trace_id", deliveryControl.TraceID).
		Str("trace_time", deliveryControl.Time).
		Int("version", deliveryControl.Version).
		Str("event", deliveryControl.Event).
		Str("action", deliveryControl.Action).
		Str("source", deliveryControl.Source).
		Str("organization", deliveryControl.OrgCode).
		Str("campaign_id", deliveryControl.ID).
		Str("group_id", deliveryControl.GroupID).
		Msg("Failed to publish sns event")
}

// delivery_controlログに整形
func (d *deliveryControlEvent) createCampaignCacheLog(campaignID int,
	groupID int, organization string, before string, after string,
	eventDetail string) *models.CampaignCacheLog {

	event, operation := d.deliveryEvent(before, after)
	current := time.Now().Format(time.RFC3339Nano)
	return &models.CampaignCacheLog{
		TraceID:     d.createTraceID(),
		Time:        current,
		Version:     config.Env.Version,
		Event:       event,
		EventDetail: eventDetail,
		Action:      operation,
		OrgCode:     organization,
		Source:      "touchgift-job-manager",
		ID:          strconv.Itoa(campaignID),
		GroupID:     strconv.Itoa(groupID),
	}
}

func (d *deliveryControlEvent) createCreativeEventLog(creative *models.DeliveryDataCreative,
	organization string, operation string) *models.CreativeCacheLog {
	return &models.CreativeCacheLog{
		ID:               creative.ID,
		Link:             creative.Link,
		URL:              creative.URL,
		Width:            creative.Width,
		Height:           creative.Height,
		Type:             creative.Type,
		Extension:        creative.Extension,
		Duration:         creative.Duration,
		EndCardUrl:       creative.EndCardURL,
		EndCardWidth:     creative.EndCardWidth,
		EndCardHeight:    creative.EndCardHeight,
		EndCardExtension: creative.EndCardExtension,
		EndCardLink:      creative.EndCardLink,
		Action:           operation,
	}
}

func (d *deliveryControlEvent) createDeliveryEventLog(id string, groupID int, storeID string,
	organization string, campaignID int, operation string) *models.DeliveryCacheLog {

	return &models.DeliveryCacheLog{
		Action:     operation,
		OrgCode:    organization,
		ID:         id,
		StoreID:    storeID,
		GroupID:    groupID,
		CampaignID: campaignID,
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
