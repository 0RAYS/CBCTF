package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsWhitelist(ctx *gin.Context) {
	clientIP := ctx.ClientIP()
	ip := net.ParseIP(clientIP)
	for _, entry := range config.Env.Gin.Metrics.Whitelist {
		if entry == clientIP {
			ctx.Next()
			return
		}
		_, cidr, err := net.ParseCIDR(entry)
		if err == nil && ip != nil && cidr.Contains(ip) {
			ctx.Next()
			return
		}
	}
	resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Forbidden})
}

func Prometheus(ctx *gin.Context) {
	start := time.Now()
	prometheus.InFlightRequests.Inc()

	ctx.Next()

	status := ctx.GetInt(resp.CTXStatusCodeKey)
	if status == 0 {
		status = ctx.Writer.Status()
	}
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
