//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package notification

import "context"

type NotificationHandler interface {
	Publish(ctx context.Context, message string, messageAttributes map[string]string, topicArn string) (*string, error)
}
