package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

var (
	metricRequestTotal     = "gin_request_total"
	metricRequestTotalDesc = "all the server received request num"

	metricURIRequestTotal       = "gin_request_uri_total"
	metricURIRequestTotalDesc   = "all the server received request num with every uri"
	metricURIRequestTotalLabels = []string{"uri", "method", "code"}

	metricResponseBody     = "gin_response_body_total"
	metricResponseBodyDesc = "the server send response body size, unit byte"

	metricRequestDuration        = "gin_request_duration_seconds"
	metricRequestDurationDesc    = "the time server took to handle the request (seconds)"
	metricRequestDurationLabels  = []string{"uri"}
	metricRequestDurationBuckets = []float64{0.025, 0.050, 0.100, 0.300, 0.500}
)

var monitor *Monitor

type Monitor struct {
	Metrics    *Metrics
	metricPath string
}

func (m *Monitor) Initialize() {
	m.Metrics.AddCounter(metricRequestTotal, metricRequestTotalDesc, nil)
	m.Metrics.AddCounter(metricURIRequestTotal, metricURIRequestTotalDesc, metricURIRequestTotalLabels)
	m.Metrics.AddCounter(metricResponseBody, metricResponseBodyDesc, nil)
	m.Metrics.AddHistogram(metricRequestDuration, metricRequestDurationDesc, metricRequestDurationLabels, metricRequestDurationBuckets)
}

func NewMonitor() *Monitor {
	monitor := Monitor{
		Metrics: NewMetrics(),
	}
	return &monitor
}

func GetMonitor() *Monitor {
	if monitor != nil {
		return monitor
	}
	monitor = NewMonitor()
	monitor.Initialize()
	return monitor
}

func (m *Monitor) AddRoute(router *gin.Engine, metricPath string) {
	m.metricPath = metricPath
	router.GET(metricPath, gin.WrapH(promhttp.Handler()))
}

func (m *Monitor) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.Request
		if request.URL.Path == m.metricPath {
			ctx.Next()
			return
		}
		startTime := time.Now()

		ctx.Next()

		latency := time.Since(startTime)

		writer := ctx.Writer
		m.Metrics.GetCounter(metricRequestTotal).WithLabelValues().Inc()
		m.Metrics.GetCounter(metricURIRequestTotal).WithLabelValues(ctx.FullPath(), request.Method, strconv.Itoa(writer.Status())).Inc()
		m.Metrics.GetHistogram(metricRequestDuration).WithLabelValues(ctx.FullPath()).Observe(latency.Seconds())
		if writer.Size() > 0 {
			m.Metrics.GetCounter(metricResponseBody).WithLabelValues().Add(float64(writer.Size()))
		}
	}
}
