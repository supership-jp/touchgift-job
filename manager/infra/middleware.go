package infra

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"unsafe"
)

type middleware struct{}

func pathAll(path *string, raw *string) *string {
	pathAll := make([]byte, 0, len(*path)+len(*raw)+1)
	pathAll = append(pathAll, *path...)
	pathAll = append(pathAll, "?"...)
	pathAll = append(pathAll, *raw...)
	return (*string)(unsafe.Pointer(&pathAll))
}
func (mw middleware) useCors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOriginFunc = func(origin string) bool { return true }
	return cors.New(config)
}
func (mw middleware) accessLog(logger *Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		appCtx := NewContext(ctx)
		appCtx.SetResponseHeader()
		defer func() {
			path := ctx.Request.URL.Path
			raw := ctx.Request.URL.RawQuery
			if raw != "" {
				path = *pathAll(&path, &raw)
			}
			dumplogger := logger.With().
				Str("ev", "acc").
				Int("status", ctx.Writer.Status()).
				Str("method", ctx.Request.Method).
				Str("path", path).
				Str("ip", ctx.ClientIP()).
				Dur("latency", time.Since(startTime)).
				Str("user-agent", ctx.Request.UserAgent()).
				Str("request_id", appCtx.RequestID())
			dumplogger.Logger()
		}()
		ctx.Next()
	}
}
