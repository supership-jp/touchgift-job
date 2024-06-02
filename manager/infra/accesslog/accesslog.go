package accesslog

import (
	"net/http"
	"os"
	"time"
	"touchgift-job-manager/manager/infra/requestid"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func pathAll(path *string, raw *string) *string {
	pathAll := make([]byte, 0, len(*path)+len(*raw)+1)
	pathAll = append(pathAll, *path...)
	pathAll = append(pathAll, "?"...)
	pathAll = append(pathAll, *raw...)
	return (*string)(unsafe.Pointer(&pathAll))
}

// New is gin.HandlerFunc
func New() gin.HandlerFunc {
	sublog := zerolog.New(os.Stdout).With().Str("ev", "acc").Timestamp().Logger()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = *pathAll(&path, &raw)
		}
		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		msg := "Request"
		if len(c.Errors) > 0 {
			msg = c.Errors.String()
		}

		dumplogger := sublog.With().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Dur("latency", latency). // Millisecond
			Str("user-agent", c.Request.UserAgent()).
			Str("request-id", requestid.Get(c)).
			Logger()

		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			{
				dumplogger.Warn().Msg(msg)
			}
		case c.Writer.Status() >= http.StatusInternalServerError:
			{
				dumplogger.Error().Msg(msg)
			}
		default:
			dumplogger.Info().Msg(msg)
		}
	}
}
