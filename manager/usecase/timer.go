//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"time"
)

// Timer is interface
type Timer interface {
	ExecuteAtTime(ctx context.Context, specifiedTime time.Time, process func())
}

type timer struct {
	logger Logger
}

// NewTimer is function
func NewTimer(
	logger Logger,
) Timer {
	return &timer{
		logger: logger,
	}
}

// 指定時間に実行する
func (d *timer) ExecuteAtTime(ctx context.Context, specifiedTime time.Time, process func()) {
	duration := time.Until(specifiedTime)
	timer := time.NewTimer(duration)
	go func() {
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			d.logger.Debug().Msg("End timer")
			return
		case <-timer.C:
			// 実行する
			process()
			return
		}
	}()
}
