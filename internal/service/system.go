package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetHomePageData(tx *gorm.DB) gin.H {
	data := gin.H{
		"upcoming":   []gin.H{},
		"stats":      []gin.H{},
		"scoreboard": []gin.H{},
	}
	if branding, ret := GetDefaultBranding(tx); ret.OK {
		data["branding"] = resp.GetBrandingResp(branding)
	}
	repo := db.InitContestRepo(tx)
	contests, count, ret := repo.List(-1, -1)
	if ret.OK {
		contestIDs := make([]uint, 0, len(contests))
		for _, contest := range contests {
			contestIDs = append(contestIDs, contest.ID)
		}
		userCountMap, _ := repo.CountUsersMap(contestIDs...)
		teamCountMap, _ := repo.CountTeamsMap(contestIDs...)
		limit := len(contests)
		if limit > 3 {
			limit = 3
		}
		for i := 0; i < limit; i++ {
			contest := contests[i]
			data["upcoming"] = append(data["upcoming"].([]gin.H), gin.H{
				"name":     contest.Name,
				"start":    contest.Start,
				"duration": int64(contest.Duration.Seconds()),
				"users":    userCountMap[contest.ID],
				"teams":    teamCountMap[contest.ID],
				"picture":  contest.Picture,
			})
		}
		data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "CTF Events", "value": count})
	}
	count, _ = db.InitUserRepo(tx).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Activate CTFers", "value": count})
	count, _ = db.InitChallengeRepo(tx).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Challenges", "value": count})
	count, _ = db.InitSubmissionRepo(tx).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Submissions", "value": count})
	users, _, _ := GetUserRanking(tx, 5, 0)
	for _, user := range users {
		data["scoreboard"] = append(data["scoreboard"].([]gin.H), gin.H{
			"name":   user.Name,
			"score":  user.Score,
			"solved": user.Solved,
		})
	}
	return data
}

func GetSystemStatus(tx *gorm.DB) map[string]any {
	ret := make(map[string]any)
	ret["metrics"] = redis.GetMetrics()
	ret["users"], _ = db.InitUserRepo(tx).Count()
	ret["contests"], _ = db.InitContestRepo(tx).Count()
	ret["ip"], _ = db.InitRequestRepo(tx).CountIP()
	ret["challenges"], _ = db.InitChallengeRepo(tx).Count()
	ret["submissions"], _ = db.InitSubmissionRepo(tx).Count(db.CountOptions{Deleted: true})
	ret["victims"], _ = db.InitVictimRepo(tx).Count(db.CountOptions{Deleted: true})
	ret["requests"], _ = db.InitRequestRepo(tx).Count(db.CountOptions{Deleted: true})
	ret["cache"] = redis.Count()
	return ret
}

func UpdateSystemSettings(tx *gorm.DB, form dto.UpdateSettingForm) model.RetVal {
	if form.K8SFrpFrps != nil {
		for i, server := range *form.K8SFrpFrps {
			if server.Token == "" {
				for _, existing := range config.Env.K8S.Frp.Frps {
					if existing.Host == server.Host && existing.Port == server.Port {
						(*form.K8SFrpFrps)[i].Token = existing.Token
						break
					}
				}
			}
		}
	}
	kv := map[string]any{
		model.HostSettingKey: form.Host,
		model.PathSettingKey: form.Path,

		model.AsyncQLogLevelSettingKey:       form.AsyncQLogLevel,
		model.AsyncQConcurrencySettingKey:    form.AsyncQConcurrency,
		model.AsyncQVictimConcurrencyKey:     form.AsyncQVictimConcurrency,
		model.AsyncQGeneratorConcurrencyKey:  form.AsyncQGeneratorConcurrency,
		model.AsyncQAttachmentConcurrencyKey: form.AsyncQAttachmentConcurrency,
		model.AsyncQEmailConcurrencyKey:      form.AsyncQEmailConcurrency,
		model.AsyncQWebhookConcurrencyKey:    form.AsyncQWebhookConcurrency,
		model.AsyncQImageConcurrencyKey:      form.AsyncQImageConcurrency,

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
		model.GinMetricsWhitelistSettingKey:   form.GinMetricsWhitelist,

		model.GormPostgresHostSettingKey:    form.GormPostgresHost,
		model.GormPostgresPortSettingKey:    form.GormPostgresPort,
		model.GormPostgresUserSettingKey:    form.GormPostgresUser,
		model.GormPostgresPwdSettingKey:     form.GormPostgresPwd,
		model.GormPostgresDBSettingKey:      form.GormPostgresDB,
		model.GormPostgresSSLModeSettingKey: form.GormPostgresSSLMode,
		model.GormPostgresMXOpenSettingKey:  form.GormPostgresMXOpen,
		model.GormPostgresMXIdleSettingKey:  form.GormPostgresMXIdle,
		model.GormLogLevelSettingKey:        form.GormLogLevel,

		model.RedisHostSettingKey: form.RedisHost,
		model.RedisPortSettingKey: form.RedisPort,
		model.RedisPwdSettingKey:  form.RedisPwd,

		model.K8SConfigSettingKey:        form.K8SConfig,
		model.K8SNamespaceSettingKey:     form.K8SNamespace,
		model.K8STCPDumpImageSettingKey:  form.K8STCPDumpImage,
		model.K8SFrpOnSettingKey:         form.K8SFrpOn,
		model.K8SFrpFrpcImageSettingKey:  form.K8SFrpFrpcImage,
		model.K8SFrpNginxImageSettingKey: form.K8SFrpNginxImage,
		model.K8SFrpFrpsSettingKey:       form.K8SFrpFrps,

		model.CheatIPWhitelistSettingKey: form.CheatIPWhitelist,
		model.WebhookWhitelistSettingKey: form.WebhookWhitelist,

		model.RegistrationEnabledSettingKey:      form.RegistrationEnabled,
		model.RegistrationDefaultGroupSettingKey: form.RegistrationDefaultGroup,

		model.GeoCityDBSettingKey: form.GeoCityDB,
	}
	repo := db.InitSettingRepo(tx)
	for key, value := range kv {
		if ret := repo.Update(key, db.UpdateSettingOptions{Value: &model.SettingValue{V: value}}); !ret.OK {
			return ret
		}
	}
	return repo.ReadSettings()
}

func GetPublicSystemConfig() map[string]any {
	return map[string]any{
		"registration_enabled": config.Env.Registration.Enabled,
	}
}
