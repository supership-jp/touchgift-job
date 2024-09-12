//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import "github.com/rs/zerolog"

// Logger is interface
type Logger interface {
	Fatal() *zerolog.Event
	Error() *zerolog.Event
	Warn() *zerolog.Event
	Info() *zerolog.Event
	Debug() *zerolog.Event
	Fatalf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	With() zerolog.Context
}
