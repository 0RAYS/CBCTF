package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/service"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/net"
)

func HomePage(ctx *gin.Context) {
	data := gin.H{
		"upcoming":   []gin.H{},
		"stats":      []gin.H{},
		"scoreboard": []gin.H{},
	}
	repo := db.InitContestRepo(db.DB)
	contests, count, ret := repo.List(-1, -1)
	if ret.OK {
		for i := 0; i < func() int {
			if len(contests) > 3 {
				return 3
			}
			return len(contests)
		}(); i++ {
			contest := contests[i]
			info := gin.H{
				"name":     contest.Name,
				"start":    contest.Start,
				"duration": contest.Duration.Seconds(),
				"users":    repo.CountAssociation(contest, "Users"),
				"teams":    repo.CountAssociation(contest, "Teams"),
				"picture":  contest.Picture,
			}
			data["upcoming"] = append(data["upcoming"].([]gin.H), info)
		}
	}
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "CTF Events", "value": count})
	count, _ = db.InitUserRepo(db.DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Activate CTFers", "value": count})
	count, _ = db.InitChallengeRepo(db.DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Challenges", "value": count})
	count, _ = db.InitSubmissionRepo(db.DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Submissions", "value": count})
	users, _, _ := service.GetUserRanking(db.DB, 5, 0)
	for _, user := range users {
		data["scoreboard"] = append(data["scoreboard"].([]gin.H), gin.H{
			"name":   user.Name,
			"score":  user.Score,
			"solved": user.Solved,
		})
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
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

	ret["users"], _ = db.InitUserRepo(db.DB).Count()
	ret["contests"], _ = db.InitContestRepo(db.DB).Count()
	ret["ip"], _ = db.InitRequestRepo(db.DB).CountIP()
	ret["challenges"], _ = db.InitChallengeRepo(db.DB).Count()
	ret["submissions"], _ = db.InitSubmissionRepo(db.DB).Count()
	ret["victims"], _ = db.InitVictimRepo(db.DB).Count()
	ret["requests"], _ = db.InitRequestRepo(db.DB).Count()
	middleware.MU.Lock()
	if middleware.TotalRequests == 0 {
		ret["duration"] = 0
	} else {
		ret["duration"] = middleware.TotalDuration.Milliseconds() / int64(middleware.TotalRequests)
	}
	middleware.MU.Unlock()

	ret["cache"] = redis.Count()
	ctx.JSON(http.StatusOK, model.SuccessRetVal(ret))
}

func GetLogs(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data, _ := redis.GetLogs(int64(form.Offset), int64(form.Offset+form.Limit))
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func SystemConfig(ctx *gin.Context) {
	data := gin.H{
		"host": config.Env.Host,
		"path": config.Env.Path,

		"log_level": config.Env.Log.Level,
		"log_save":  config.Env.Log.Save,

		"asyncq_log_level":   config.Env.AsyncQ.Log.Level,
		"asyncq_concurrency": config.Env.AsyncQ.Concurrency,

		"gin_mode":                config.Env.Gin.Mode,
		"gin_host":                config.Env.Gin.Host,
		"gin_port":                config.Env.Gin.Port,
		"gin_upload_max":          config.Env.Gin.Upload.Max,
		"gin_proxies":             config.Env.Gin.Proxies,
		"gin_ratelimit_global":    config.Env.Gin.RateLimit.Global,
		"gin_ratelimit_whitelist": config.Env.Gin.RateLimit.Whitelist,
		"gin_cors":                config.Env.Gin.CORS,
		"gin_log_whitelist":       config.Env.Gin.Log.Whitelist,
		"gin_jwt_secret":          config.Env.Gin.JWT.Secret,
		"gin_jwt_static":          config.Env.Gin.JWT.Static,

		"gorm_mysql_host":   config.Env.Gorm.MySQL.Host,
		"gorm_mysql_port":   config.Env.Gorm.MySQL.Port,
		"gorm_mysql_user":   config.Env.Gorm.MySQL.User,
		"gorm_mysql_pwd":    "******",
		"gorm_mysql_db":     config.Env.Gorm.MySQL.DB,
		"gorm_mysql_mxopen": config.Env.Gorm.MySQL.MaxOpenConns,
		"gorm_mysql_mxidle": config.Env.Gorm.MySQL.MaxIdleConns,
		"gorm_log_level":    config.Env.Gorm.Log.Level,

		"redis_host": config.Env.Redis.Host,
		"redis_port": config.Env.Redis.Port,
		"redis_pwd":  "******",

		"k8s_config":                       config.Env.K8S.Config,
		"k8s_namespace":                    config.Env.K8S.Namespace,
		"k8s_external_network_cidr":        config.Env.K8S.ExternalNetwork.CIDR,
		"k8s_external_network_gateway":     config.Env.K8S.ExternalNetwork.Gateway,
		"k8s_external_network_interface":   config.Env.K8S.ExternalNetwork.Interface,
		"k8s_external_network_exclude_ips": config.Env.K8S.ExternalNetwork.ExcludeIPs,
		"k8s_tcpdump":                      config.Env.K8S.TCPDumpImage,

		"k8s_frp_on":    config.Env.K8S.Frp.On,
		"k8s_frp_frpc":  config.Env.K8S.Frp.FrpcImage,
		"k8s_frp_nginx": config.Env.K8S.Frp.NginxImage,
		"k8s_frp_frps":  config.Env.K8S.Frp.Frps,

		"k8s_generator_worker": config.Env.K8S.GeneratorWorker,

		"nfs_server":  config.Env.NFS.Server,
		"nfs_path":    config.Env.NFS.Path,
		"nfs_storage": config.Env.NFS.Storage,

		"cheat_ip_whitelist": config.Env.Cheat.IP.Whitelist,

		"webhook_blacklist": config.Env.Webhook.Blacklist,

		"registration_enabled":       config.Env.Registration.Enabled,
		"registration_default_group": config.Env.Registration.DefaultGroup,

		"geocity_db": config.Env.GeoCityDB,
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func UpdateSystem(ctx *gin.Context) {
	var form dto.UpdateSettingForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateSettingEventType)
	kv := map[string]any{
		model.HostSettingKey: form.Host,
		model.PathSettingKey: form.Path,

		model.LogLevelSettingKey: form.LogLevel,
		model.LogSaveSettingKey:  form.LogSave,

		model.AsyncQLogLevelSettingKey:    form.AsyncQLogLevel,
		model.AsyncQConcurrencySettingKey: form.AsyncQConcurrency,

		model.GinModeSettingKey:               form.GinMode,
		model.GinHostSettingKey:               form.GinHost,
		model.GinPortSettingKey:               form.GinPort,
		model.GinUploadMaxSettingKey:          form.GinUploadMax,
		model.GinProxiesSettingKey:            form.GinProxies,
		model.GinRateLimitGlobalSettingKey:    form.GinRateLimitGlobal,
		model.GinRateLimitWhitelistSettingKey: form.GinRateLimitWhitelist,
		model.GinCORSSettingKey:               form.GinCORS,
		model.GinLogWhitelistSettingKey:       form.GinLogWhitelist,
		model.GinJWTSecretSettingKey:          form.GinJWTSecret,
		model.GinJWTStaticSettingKey:          form.GinJWTStatic,

		model.GormMySQLHostSettingKey:   form.GormMySQLHost,
		model.GormMySQLPortSettingKey:   form.GormMySQLPort,
		model.GormMySQLUserSettingKey:   form.GormMySQLUser,
		model.GormMySQLPwdSettingKey:    form.GormMySQLPwd,
		model.GormMySQLDBSettingKey:     form.GormMySQLDB,
		model.GormMySQLMXOpenSettingKey: form.GormMySQLMXOpen,
		model.GormMySQLMXIdleSettingKey: form.GormMySQLMXIdle,
		model.GormLogLevelSettingKey:    form.GormLogLevel,

		model.RedisHostSettingKey: form.RedisHost,
		model.RedisPortSettingKey: form.RedisPort,
		model.RedisPwdSettingKey:  form.RedisPwd,

		model.K8SConfigSettingKey:                    form.K8SConfig,
		model.K8SNamespaceSettingKey:                 form.K8SNamespace,
		model.K8SExternalNetworkCIDRSettingKey:       form.K8SExternalNetworkCIDR,
		model.K8SExternalNetworkGatewaySettingKey:    form.K8SExternalNetworkGateway,
		model.K8SExternalNetworkInterfaceSettingKey:  form.K8SExternalNetworkInterface,
		model.K8SExternalNetworkExcludeIPsSettingKey: form.K8SExternalNetworkExcludeIPs,
		model.K8STCPDumpImageSettingKey:              form.K8STCPDumpImage,

		model.K8SFrpOnSettingKey:         form.K8SFrpOn,
		model.K8SFrpFrpcImageSettingKey:  form.K8SFrpFrpcImage,
		model.K8SFrpNginxImageSettingKey: form.K8SFrpNginxImage,
		model.K8SFrpFrpsSettingKey:       form.K8SFrpFrps,

		model.K8SGeneratorWorkerSettingKey: form.K8SGeneratorWorker,

		model.NFSServerSettingKey:  form.NFSServer,
		model.NFSPathSettingKey:    form.NFSPath,
		model.NFSStorageSettingKey: form.NFSStorage,

		model.CheatIPWhitelistSettingKey: form.CheatIPWhitelist,

		model.WebhookBlacklistSettingKey: form.WebhookBlacklist,

		model.RegistrationEnabledSettingKey:      form.RegistrationEnabled,
		model.RegistrationDefaultGroupSettingKey: form.RegistrationDefaultGroup,

		model.GeoCityDBSettingKey: form.GeoCityDB,
	}
	repo := db.InitSettingRepo(db.DB)
	for key, value := range kv {
		if ret := repo.Update(key, db.UpdateSettingOptions{Value: &model.SettingValue{V: value}}); !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
	}
	// 读取数据库配置至内存并覆写配置文件
	if ret := repo.ReadSettings(); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func RestartSystem(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.RestartSystemEventType)
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	go func(proc *os.Process) {
		_ = proc.Signal(syscall.SIGUSR1)
	}(proc)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
