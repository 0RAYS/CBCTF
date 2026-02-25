package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	p "CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsRouter(router *gin.Engine) {
	// 不使用默认 registry，防止重启时重复注册导致 panic
	var (
		registry                         = prometheus.NewRegistry()
		registerer prometheus.Registerer = registry
		gatherer   prometheus.Gatherer   = registry
	)

	// 注册 DB 驱动的自定义 Collector
	registerer.MustRegister(p.NewCTFCollector())

	// 注册 HTTP 基础指标
	registerer.MustRegister(p.HttpRequestsTotal)
	registerer.MustRegister(p.HttpRequestDuration)
	registerer.MustRegister(p.HttpResponseSize)
	registerer.MustRegister(p.InFlightRequests)

	// 注册 CTF 业务事件指标
	registerer.MustRegister(p.FlagSubmissionsTotal)
	registerer.MustRegister(p.BloodTotal)
	registerer.MustRegister(p.UserRegistrationTotal)
	registerer.MustRegister(p.UserLoginTotal)
	registerer.MustRegister(p.FileUploadTotal)
	registerer.MustRegister(p.FileUploadSize)
	registerer.MustRegister(p.WebSocketConnections)
	registerer.MustRegister(p.EmailSentTotal)
	registerer.MustRegister(p.RateLimitHits)
	registerer.MustRegister(p.CheatDetectionsTotal)

	// 注册 Cron Job 指标
	registerer.MustRegister(p.CronJobDuration)
	registerer.MustRegister(p.CronJobRunsTotal)

	// 注册异步任务指标
	registerer.MustRegister(p.TaskEnqueuedTotal)
	registerer.MustRegister(p.TaskProcessedTotal)
	registerer.MustRegister(p.TaskProcessingDuration)

	// 注册 Go 运行时指标
	registerer.MustRegister(collectors.NewGoCollector())
	registerer.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	router.GET("/metrics", func(ctx *gin.Context) {
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
	}, gin.WrapH(promhttp.InstrumentMetricHandler(
		registerer, promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}),
	)))
}
