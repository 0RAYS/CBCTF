package middleware

import (
	"CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Prometheus(ctx *gin.Context) {
	start := time.Now()
	prometheus.InFlightRequests.Inc()

	ctx.Next()

	status := ctx.GetInt(resp.CTXStatusCodeKey)
	path := ctx.FullPath()
	if path == "" {
		path = ctx.Request.URL.Path
	}

	prometheus.InFlightRequests.Dec()
	prometheus.HttpRequestsTotal.WithLabelValues(
		ctx.Request.Method, path, strconv.Itoa(status),
	).Inc()
	prometheus.HttpRequestDuration.WithLabelValues(
		ctx.Request.Method, path,
	).Observe(time.Since(start).Seconds())
	if size := ctx.Writer.Size(); size > 0 {
		prometheus.HttpResponseSize.WithLabelValues(
			ctx.Request.Method, path,
		).Observe(float64(size))
	}
}
