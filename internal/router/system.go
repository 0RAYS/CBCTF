package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/sys"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/net"
)

func HomePage(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(service.GetHomePageData(db.DB)))
}

func PublicSystemConfig(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(service.GetPublicSystemConfig()))
}

func SystemStatus(ctx *gin.Context) {
	ret := make(map[string]any)
	ret["metrics"] = redis.GetMetrics()

	ioStats, err := net.IOCounters(false)
	if err != nil || len(ioStats) == 0 {
		ret["io"] = 0
		ret["sent"] = 0
		ret["recv"] = 0
	} else {
		ret["io"] = ioStats[0].BytesSent + ioStats[0].BytesRecv
		ret["sent"] = ioStats[0].BytesSent
		ret["recv"] = ioStats[0].BytesRecv
	}

	for key, value := range service.GetSystemStatus(db.DB) {
		ret[key] = value
	}
	middleware.MU.Lock()
	if middleware.TotalRequests == 0 {
		ret["duration"] = 0
	} else {
		ret["duration"] = middleware.TotalDuration.Milliseconds() / int64(middleware.TotalRequests)
	}
	middleware.MU.Unlock()
	resp.JSON(ctx, model.SuccessRetVal(ret))
}

func GetLogs(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data, _ := redis.GetLogs(int64(form.Offset), int64(form.Offset+form.Limit))
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func SystemConfig(ctx *gin.Context) {
	maskedFrps := make([]gin.H, len(config.Env.K8S.Frp.Frps))
	for i, frps := range config.Env.K8S.Frp.Frps {
		maskedFrps[i] = gin.H{
			"host":    frps.Host,
			"port":    frps.Port,
			"allowed": frps.Allowed,
		}
	}
	data := gin.H{
		"host": config.Env.Host,
		"path": config.Env.Path,

		"log_level": config.Env.Log.Level,
		"log_save":  config.Env.Log.Save,

		"asyncq_log_level":              config.Env.AsyncQ.Log.Level,
		"asyncq_concurrency":            config.Env.AsyncQ.Concurrency,
		"asyncq_victim_concurrency":     config.Env.AsyncQ.Queues.Victim,
		"asyncq_generator_concurrency":  config.Env.AsyncQ.Queues.Generator,
		"asyncq_attachment_concurrency": config.Env.AsyncQ.Queues.Attachment,
		"asyncq_email_concurrency":      config.Env.AsyncQ.Queues.Email,
		"asyncq_webhook_concurrency":    config.Env.AsyncQ.Queues.Webhook,
		"asyncq_image_concurrency":      config.Env.AsyncQ.Queues.Image,

		"gin_mode":                config.Env.Gin.Mode,
		"gin_host":                config.Env.Gin.Host,
		"gin_port":                config.Env.Gin.Port,
		"gin_upload_max":          config.Env.Gin.Upload.Max,
		"gin_proxies":             config.Env.Gin.Proxies,
		"gin_ratelimit_global":    config.Env.Gin.RateLimit.Global,
		"gin_ratelimit_whitelist": config.Env.Gin.RateLimit.Whitelist,
		"gin_cors":                config.Env.Gin.CORS,
		"gin_log_whitelist":       config.Env.Gin.Log.Whitelist,
		"gin_metrics_whitelist":   config.Env.Gin.Metrics.Whitelist,

		"gorm_postgres_host":    config.Env.Gorm.Postgres.Host,
		"gorm_postgres_port":    config.Env.Gorm.Postgres.Port,
		"gorm_postgres_user":    config.Env.Gorm.Postgres.User,
		"gorm_postgres_pwd":     "******",
		"gorm_postgres_db":      config.Env.Gorm.Postgres.DB,
		"gorm_postgres_sslmode": config.Env.Gorm.Postgres.SSLMode,
		"gorm_postgres_mxopen":  config.Env.Gorm.Postgres.MaxOpenConns,
		"gorm_postgres_mxidle":  config.Env.Gorm.Postgres.MaxIdleConns,
		"gorm_log_level":        config.Env.Gorm.Log.Level,

		"redis_host": config.Env.Redis.Host,
		"redis_port": config.Env.Redis.Port,
		"redis_pwd":  "******",

		"k8s_config":    config.Env.K8S.Config,
		"k8s_namespace": config.Env.K8S.Namespace,
		"k8s_tcpdump":   config.Env.K8S.TCPDumpImage,

		"k8s_frp_on":    config.Env.K8S.Frp.On,
		"k8s_frp_frpc":  config.Env.K8S.Frp.FrpcImage,
		"k8s_frp_nginx": config.Env.K8S.Frp.NginxImage,
		"k8s_frp_frps":  maskedFrps,

		"cheat_ip_whitelist": config.Env.Cheat.IP.Whitelist,

		"webhook_whitelist": config.Env.Webhook.Whitelist,

		"registration_enabled":       config.Env.Registration.Enabled,
		"registration_default_group": config.Env.Registration.DefaultGroup,

		"geocity_db": config.Env.GeoCityDB,
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func UpdateSystem(ctx *gin.Context) {
	var form dto.UpdateSettingForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateSettingEventType)
	if ret := service.UpdateSystemSettings(db.DB, form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func RestartSystem(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.RestartSystemEventType)
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	go func(proc *os.Process) {
		_ = sys.Restart(proc)
	}(proc)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
