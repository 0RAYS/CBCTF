package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"CBCTF/internal/service"
	"fmt"
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
	middleware.MU.Lock()
	if middleware.TotalRequests == 0 {
		ret["requests"] = 0
		ret["duration"] = 0
	} else {
		ret["requests"] = middleware.TotalRequests
		ret["duration"] = middleware.TotalDuration.Milliseconds() / int64(middleware.TotalRequests)
	}
	middleware.MU.Unlock()

	total, hit, miss := redis.Status()
	ret["cache"] = total
	ret["hit"] = hit
	if hit+miss == 0 {
		ret["rate"] = "0.00"
	} else {
		ret["rate"] = fmt.Sprintf("%.2f", float64(hit)/float64(hit+miss)*100)
	}
	prometheus.UpdateCacheMetrics("redis", hit, miss)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(ret))
}

func SystemConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, model.SuccessRetVal(config.Env))
}

func UpdateSystem(ctx *gin.Context) {
	var form dto.UpdateSettingForm
	if ret := form.Bind(ctx); !ret.OK {
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
	}
	repo := db.InitSettingRepo(db.DB)
	for key, value := range kv {
		if ret := repo.Update(key, db.UpdateSettingOptions{Value: &model.SettingValue{V: value}}); !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
	}
	// 读取数据库配置至内存
	if ret := repo.ReadSettings(); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	// 覆写文件
	if err := config.Save(); err != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func RestartSystem(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateSettingEventType)
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
