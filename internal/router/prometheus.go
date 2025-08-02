package router

import (
	"CBCTF/internal/log"
	p "CBCTF/internal/prometheus"
	"errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsRouter(router *gin.Engine) {
	pprof.Register(router)

	// 注册HTTP基础指标
	prometheus.MustRegister(p.HttpRequestsTotal)
	prometheus.MustRegister(p.HttpRequestDuration)
	prometheus.MustRegister(p.HttpRequestSize)
	prometheus.MustRegister(p.HttpResponseSize)
	prometheus.MustRegister(p.InFlightRequests)

	// 注册CTF业务指标
	prometheus.MustRegister(p.FlagSubmissionTotal)
	prometheus.MustRegister(p.ContestActiveTeams)
	prometheus.MustRegister(p.ContestActiveUsers)
	prometheus.MustRegister(p.VictimContainerTotal)
	prometheus.MustRegister(p.UserRegistrationTotal)
	prometheus.MustRegister(p.UserLoginTotal)
	prometheus.MustRegister(p.FileUploadTotal)
	prometheus.MustRegister(p.FileUploadSize)
	prometheus.MustRegister(p.WebSocketConnections)
	prometheus.MustRegister(p.EmailSentTotal)
	prometheus.MustRegister(p.CacheHitRate)
	prometheus.MustRegister(p.RateLimitHits)
	prometheus.MustRegister(p.ErrorTotal)
	var alreadyRegisteredError prometheus.AlreadyRegisteredError
	if err := prometheus.Register(collectors.NewGoCollector()); err != nil {
		if !errors.As(err, &alreadyRegisteredError) {
			log.Logger.Warningf("failed to register GoCollector: %v", err)
		}
	}
	if err := prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		if !errors.As(err, &alreadyRegisteredError) {
			log.Logger.Warningf("failed to register ProcessCollector: %v", err)
		}
	}
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
