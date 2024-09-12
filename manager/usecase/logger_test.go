package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

type TestLogger struct {
	t        testing.TB
	delegate *zerolog.Logger
}

func NewTestLogger(t testing.TB) *TestLogger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	logger := TestLogger{t: t, delegate: &log}
	return &logger
}

func (i *TestLogger) Fatal() *zerolog.Event {
	return i.delegate.Fatal()
}
func (i *TestLogger) Error() *zerolog.Event {
	return i.delegate.Error()
}
func (i *TestLogger) Warn() *zerolog.Event {
	return i.delegate.Warn()
}
func (i *TestLogger) Info() *zerolog.Event {
	return i.delegate.Info()
}
func (i *TestLogger) Debug() *zerolog.Event {
	return i.delegate.Debug()
}
func (i *TestLogger) Fatalf(format string, v ...interface{}) {
	i.delegate.Fatal().Msgf(format, v...)
}
func (i *TestLogger) Errorf(format string, v ...interface{}) {
	i.delegate.Error().Msgf(format, v...)
}
func (i *TestLogger) Infof(format string, v ...interface{}) {
	i.delegate.Info().Msgf(format, v...)
}
func (i *TestLogger) Warnf(format string, v ...interface{}) {
	i.delegate.Warn().Msgf(format, v...)
}
func (i *TestLogger) Debugf(format string, v ...interface{}) {
	i.delegate.Debug().Msgf(format, v...)
}
func (i *TestLogger) With() zerolog.Context {
	return i.delegate.With()
}
