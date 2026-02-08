package router

import (
	p "CBCTF/internal/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsRouter(router *gin.Engine) {
	// 不使用默认 registry, 防止重启时重复注册导致 panic
	var (
		registry                         = prometheus.NewRegistry()
		registerer prometheus.Registerer = registry
		gatherer   prometheus.Gatherer   = registry
	)

	// 注册HTTP基础指标
	registerer.MustRegister(p.HttpRequestsTotal)
	registerer.MustRegister(p.HttpRequestDuration)
	registerer.MustRegister(p.HttpRequestSize)
	registerer.MustRegister(p.HttpResponseSize)
	registerer.MustRegister(p.InFlightRequests)

	// 注册CTF业务指标
	registerer.MustRegister(p.FlagSubmissionTotal)
	registerer.MustRegister(p.ContestActiveTeams)
	registerer.MustRegister(p.ContestActiveUsers)
	registerer.MustRegister(p.VictimContainerTotal)
	registerer.MustRegister(p.UserRegistrationTotal)
	registerer.MustRegister(p.UserLoginTotal)
	registerer.MustRegister(p.FileUploadTotal)
	registerer.MustRegister(p.FileUploadSize)
	registerer.MustRegister(p.WebSocketConnections)
	registerer.MustRegister(p.EmailSentTotal)
	registerer.MustRegister(p.CacheHitRate)
	registerer.MustRegister(p.RateLimitHits)
	registerer.MustRegister(p.ErrorTotal)

	registerer.MustRegister(collectors.NewGoCollector())
	registerer.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	router.GET("/metrics", gin.WrapH(promhttp.InstrumentMetricHandler(
		registerer, promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}),
	)))
}
