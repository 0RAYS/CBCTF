package middleware

import (
	"CBCTF/internal/prometheus"
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func Prometheus(ctx *gin.Context) {
	start := time.Now()
	prometheus.InFlightRequests.Inc()
	reqSize := getRequestBodySize(ctx)

	ctx.Next()

	duration := time.Since(start).Seconds()
	status := ctx.Writer.Status()
	path := ctx.FullPath()
	if path == "" {
		path = ctx.Request.URL.Path
	}

	prometheus.InFlightRequests.Dec()
	prometheus.HttpRequestsTotal.WithLabelValues(path, ctx.Request.Method, http.StatusText(status)).Inc()
	prometheus.HttpRequestDuration.WithLabelValues(path, ctx.Request.Method).Observe(duration)
	prometheus.HttpRequestSize.WithLabelValues(path, ctx.Request.Method).Observe(float64(reqSize))
	prometheus.HttpResponseSize.WithLabelValues(path, ctx.Request.Method).Observe(float64(ctx.Writer.Size()))

	// 记录错误
	if status >= 400 {
		prometheus.ErrorTotal.WithLabelValues(http.StatusText(status), "http").Inc()
	}
}
