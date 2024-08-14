//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package controllers

import "time"

type AppTicker interface {
	New(interval time.Duration, unit time.Duration) *time.Ticker
}

type appTicker struct{}

func NewAppTicker() AppTicker {
	return &appTicker{}
}

func (a *appTicker) New(interval time.Duration, unit time.Duration) *time.Ticker {
	time.Sleep(time.Until(time.Now().Add(1 * unit).Truncate(unit)))
	return time.NewTicker(interval)
}
