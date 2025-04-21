package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"io"
	"net/http"
	"time"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	HttpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Histogram of HTTP request body sizes.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10), // 100B ~ 51KB
		},
		[]string{"path", "method"},
	)

	HttpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Histogram of HTTP response body sizes.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"path", "method"},
	)

	InFlightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "in_flight_requests",
			Help: "Current number of in-flight requests being handled.",
		},
	)
)

func getRequestBodySize(ctx *gin.Context) int {
	if ctx.Request.Body == nil {
		return 0
	}

	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return 0
	}

	size := len(bodyBytes)

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return size
}

func init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(HttpRequestSize)
	prometheus.MustRegister(HttpResponseSize)
	prometheus.MustRegister(InFlightRequests)
	prometheus.MustRegister(collectors.NewGoCollector())
	prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}

func Prometheus(ctx *gin.Context) {
	start := time.Now()
	InFlightRequests.Inc()
	reqSize := getRequestBodySize(ctx)

	ctx.Next()

	duration := time.Since(start).Seconds()
	status := ctx.Writer.Status()
	path := ctx.FullPath()
	if path == "" {
		path = ctx.Request.URL.Path
	}

	InFlightRequests.Dec()
	HttpRequestsTotal.WithLabelValues(path, ctx.Request.Method, http.StatusText(status)).Inc()
	HttpRequestDuration.WithLabelValues(path, ctx.Request.Method).Observe(duration)
	HttpRequestSize.WithLabelValues(path, ctx.Request.Method).Observe(float64(reqSize))
	HttpResponseSize.WithLabelValues(path, ctx.Request.Method).Observe(float64(ctx.Writer.Size()))
}
