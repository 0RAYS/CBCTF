package middleware

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
)

// MetricsPath 是 Prometheus 抓取端点。它自身的请求不计入 HTTP 指标，
// 否则抓取频率和采集耗时会被算进流量与延迟分位数里。
const MetricsPath = "/metrics"

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
	resp.AbortJSON(ctx, model.FailRetVal(i18n.Response.Forbidden, "Prometheus", "", nil))
}

func Prometheus(ctx *gin.Context) {
	if ctx.FullPath() == MetricsPath {
		ctx.Next()
		return
	}

	start := time.Now()
	prometheus.InFlightRequests.Inc()
	defer prometheus.InFlightRequests.Dec()

	ctx.Next()

	code := ctx.GetInt(resp.CTXStatusCodeKey)
	if code == 0 {
		code = ctx.Writer.Status()
	}
	path := ctx.FullPath()
	if path == "" {
		path = "__not_found__"
	}
	statusClass := strconv.Itoa(code/100) + "xx"

	prometheus.HttpRequestsTotal.WithLabelValues(ctx.Request.Method, path, strconv.Itoa(code), statusClass).Inc()
	prometheus.HttpRequestDuration.WithLabelValues(ctx.Request.Method, path, statusClass).Observe(time.Since(start).Seconds())
	prometheus.HttpResponseSize.WithLabelValues(ctx.Request.Method, path, statusClass).Observe(float64(ctx.Writer.Size()))
	if size := ctx.Request.ContentLength; size > 0 {
		prometheus.HttpRequestSize.WithLabelValues(ctx.Request.Method, path).Observe(float64(size))
	}
}
