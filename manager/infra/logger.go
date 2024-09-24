//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package infra

import (
	"os"
	"time"
	"touchgift-job-manager/config"

	"github.com/rs/zerolog"
)

var log *Logger

func init() {
	zerolog.DisableSampling(true)
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// GetLogger is function
func GetLogger() *Logger {
	if log != nil {
		return log
	}
	level, err := zerolog.ParseLevel(config.Env.LogLevel)
	if err != nil {
		panic(err)
	} else {
		zerolog.SetGlobalLevel(level)
	}
	logger := zerolog.New(os.Stdout).With().Str("ev", "application").Timestamp().Caller().Logger()
	log = &Logger{
		delegate: &logger,
	}
	return log
}

// Logger is struct
type Logger struct {
	delegate *zerolog.Logger
}

// NewLogger is function
func NewLogger(l *zerolog.Logger) Logger {
	return Logger{
		delegate: l,
	}
}

// Fatal is function
func (logger Logger) Fatal() *zerolog.Event {
	return logger.delegate.Fatal()
}

// Error is function
func (logger Logger) Error() *zerolog.Event {
	return logger.delegate.Error()
}

// Warn is function
func (logger Logger) Warn() *zerolog.Event {
	return logger.delegate.Warn()
}

// Info is function
func (logger Logger) Info() *zerolog.Event {
	return logger.delegate.Info()
}

// Debug is function
func (logger Logger) Debug() *zerolog.Event {
	return logger.delegate.Debug()
}

// Fatalf is function
func (logger Logger) Fatalf(format string, v ...interface{}) {
	logger.delegate.Fatal().Msgf(format, v...)
}

// Errorf is function
func (logger Logger) Errorf(format string, v ...interface{}) {
	logger.delegate.Error().Msgf(format, v...)
}

// Warnf is function
func (logger Logger) Warnf(format string, v ...interface{}) {
	logger.delegate.Warn().Msgf(format, v...)
}

// Infof is function
func (logger Logger) Infof(format string, v ...interface{}) {
	logger.delegate.Info().Msgf(format, v...)
}

// Debugf is function
func (logger Logger) Debugf(format string, v ...interface{}) {
	logger.delegate.Debug().Msgf(format, v...)
}

// With is function
func (logger Logger) With() zerolog.Context {
	return logger.delegate.With()
}
