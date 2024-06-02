package infra

import (
	"github.com/gin-gonic/gin"
	"io"
	"touchgift-job-manager/manager/infra/accesslog"
	"touchgift-job-manager/manager/infra/metrics"
	"touchgift-job-manager/manager/infra/requestid"
)

func NewRouter(log *Logger) *gin.Engine {

	validatorSupport := NewValidatorSupport(log)
	errorSupport := NewErrorSupport(log, validatorSupport)

	router := gin.New()

	monitor := metrics.GetMonitor()
	router.Use(monitor.Middleware())
	router.Use(requestid.New())
	router.Use(accesslog.New())
	router.Use(errorSupport.middleware())
	router.Use(recovery(log))
	return router
}

func recovery(log *Logger) gin.HandlerFunc {
	return gin.RecoveryWithWriter(NewErrorWriter(log))
}

// NewErrorWriter is function
func NewErrorWriter(log *Logger) io.Writer {
	return &ErrorWriter{
		log: log,
	}
}

// ErrorWriter is struct
type ErrorWriter struct {
	log *Logger
}

func (e *ErrorWriter) Write(p []byte) (n int, err error) {
	e.log.Error().Msg(string(p))
	return len(p), nil
}
