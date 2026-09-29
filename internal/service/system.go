package service

import (
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/validation"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
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
	contests, count, ret := repo.List(-1, -1, db.GetOptions{Sort: []string{"start ASC"}})
	if ret.OK {
		contestIDs := make([]uint, 0, len(contests))
		for _, contest := range contests {
			contestIDs = append(contestIDs, contest.ID)
		}
		userCountMap, userRet := repo.CountUsersMap(contestIDs...)
		teamCountMap, teamRet := repo.CountTeamsMap(contestIDs...)
		limit := min(len(contests), 3)
		for i := range limit {
			contest := contests[i]
			var users, teams any
			if userRet.OK {
				users = userCountMap[contest.ID]
			}
			if teamRet.OK {
				teams = teamCountMap[contest.ID]
			}
			data["upcoming"] = append(data["upcoming"].([]gin.H), gin.H{
				"name":     contest.Name,
				"start":    contest.Start,
				"duration": int64(contest.Duration.Seconds()),
				"users":    users,
				"teams":    teams,
				"picture":  contest.Picture,
			})
		}
		data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "CTF Events", "value": count})
	}
	appendCount := func(label string, count int64, ret model.RetVal) {
		var value any
		if ret.OK {
			value = count
		}
		data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": label, "value": value})
	}
	count, ret = db.InitUserRepo(tx).Count()
	appendCount("Activate CTFers", count, ret)
	count, ret = db.InitChallengeRepo(tx).Count()
	appendCount("Challenges", count, ret)
	count, ret = db.InitSubmissionRepo(tx).Count()
	appendCount("Submissions", count, ret)
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
	unavailable := make([]string, 0)
	metrics, skipped, metricRet := redis.GetMetrics(tx.Statement.Context)
	ret["metrics"], ret["metrics_skipped"] = metrics, skipped
	if !metricRet.OK {
		unavailable = append(unavailable, "metrics")
	}
	for _, counter := range []struct {
		name string
		read func() (int64, model.RetVal)
	}{
		{"users", func() (int64, model.RetVal) { return db.InitUserRepo(tx).Count() }},
		{"contests", func() (int64, model.RetVal) { return db.InitContestRepo(tx).Count() }},
		{"ip", db.InitRequestRepo(tx).CountIP},
		{"challenges", func() (int64, model.RetVal) { return db.InitChallengeRepo(tx).Count() }},
		{"submissions", func() (int64, model.RetVal) { return db.InitSubmissionRepo(tx).Count(db.CountOptions{Deleted: true}) }},
		{"victims", func() (int64, model.RetVal) { return db.InitVictimRepo(tx).Count(db.CountOptions{Deleted: true}) }},
		{"requests", func() (int64, model.RetVal) { return db.InitRequestRepo(tx).Count(db.CountOptions{Deleted: true}) }},
	} {
		value, result := counter.read()
		if result.OK {
			ret[counter.name] = value
		} else {
			ret[counter.name] = nil
			unavailable = append(unavailable, counter.name)
		}
	}
	cache, err := redis.Count(tx.Statement.Context)
	if err == nil {
		ret["cache"] = cache
	} else {
		ret["cache"] = nil
		unavailable = append(unavailable, "cache")
	}
	ret["unavailable"] = unavailable
	return ret
}

func validateK8sSettings(form dto.UpdateSettingForm) model.RetVal {
	if name := form.K8SNamespace; name != nil {
		if errors := validation.IsDNS1123Label(*name); len(errors) > 0 {
			return model.RetVal{Msg: i18n.Response.BadRequest, Attr: map[string]any{"Error": strings.Join(errors, "; ")}}
		}
	}
	if name := form.K8SPriorityClassName; name != nil && *name != "" {
		if errors := validation.IsDNS1123Subdomain(*name); len(errors) > 0 {
			return model.RetVal{Msg: i18n.Response.BadRequest, Attr: map[string]any{"Error": strings.Join(errors, "; ")}}
		}
	}
	return model.SuccessRetVal()
}

var settingsUpdateMu sync.Mutex

func UpdateSystemSettings(tx *gorm.DB, form dto.UpdateSettingForm) model.RetVal {
	settingsUpdateMu.Lock()
	defer settingsUpdateMu.Unlock()
	if ret := validateK8sSettings(form); !ret.OK {
		return ret
	}
	if form.K8SFrpFrps != nil {
		for i, server := range *form.K8SFrpFrps {
			if server.Token != "" {
				continue
			}
			// 优先按 host:port 精确匹配, 保证删除或重排其他节点后 token 仍能归位.
			matched := false
			for _, existing := range config.Env.K8S.Frp.Frps {
				if existing.Host == server.Host && existing.Port == server.Port {
					(*form.K8SFrpFrps)[i].Token = existing.Token
					matched = true
					break
				}
			}
			// 编辑既有节点的 host/port 而未重新填写 token 时, 回退到按位置匹配, 避免 token 丢失.
			if !matched && i < len(config.Env.K8S.Frp.Frps) {
				(*form.K8SFrpFrps)[i].Token = config.Env.K8S.Frp.Frps[i].Token
			}
		}
	}
	kv := map[string]any{
		model.HostSettingKey: form.Host,

		model.AsyncQLogLevelSettingKey:       form.AsyncQLogLevel,
		model.AsyncQVictimConcurrencyKey:     form.AsyncQVictimConcurrency,
		model.AsyncQTrafficConcurrencyKey:    form.AsyncQTrafficConcurrency,
		model.AsyncQGeneratorConcurrencyKey:  form.AsyncQGeneratorConcurrency,
		model.AsyncQAttachmentConcurrencyKey: form.AsyncQAttachmentConcurrency,
		model.AsyncQEmailConcurrencyKey:      form.AsyncQEmailConcurrency,
		model.AsyncQWebhookConcurrencyKey:    form.AsyncQWebhookConcurrency,
		model.AsyncQImageConcurrencyKey:      form.AsyncQImageConcurrency,

		model.GinModeSettingKey:               form.GinMode,
		model.GinUploadPictureSettingKey:      form.GinUploadPicture,
		model.GinUploadChallengeSettingKey:    form.GinUploadChallenge,
		model.GinUploadWriteupSettingKey:      form.GinUploadWriteup,
		model.GinProxiesSettingKey:            form.GinProxies,
		model.GinRateLimitGlobalSettingKey:    form.GinRateLimitGlobal,
		model.GinRateLimitWhitelistSettingKey: form.GinRateLimitWhitelist,
		model.GinOriginsSettingKey:            form.GinOrigins,
		model.GinLogWhitelistSettingKey:       form.GinLogWhitelist,
		model.GinJWTSecretSettingKey:          form.GinJWTSecret,
		model.GinMetricsWhitelistSettingKey:   form.GinMetricsWhitelist,
		model.GinPProfWhitelistSettingKey:     form.GinPProfWhitelist,

		model.K8SNamespaceSettingKey:         form.K8SNamespace,
		model.K8SCaptureImageSettingKey:      form.K8SCaptureImage,
		model.K8SCaptureEnabledSettingKey:    form.K8SCaptureEnabled,
		model.K8SPriorityClassSettingKey:     form.K8SPriorityClassName,
		model.K8SWorkerImageSettingKey:       form.K8SWorkerImage,
		model.K8SGeneratorPoolSizeSettingKey: form.K8SGeneratorPoolSize,
		model.K8SFrpOnSettingKey:             form.K8SFrpOn,
		model.K8SFrpFrpcImageSettingKey:      form.K8SFrpFrpcImage,
		model.K8SFrpNginxImageSettingKey:     form.K8SFrpNginxImage,
		model.K8SFrpFrpsSettingKey:           form.K8SFrpFrps,

		model.CheatIPWhitelistSettingKey:         form.CheatIPWhitelist,
		model.WebhookWhitelistSettingKey:         form.WebhookWhitelist,
		model.RegistrationEnabledSettingKey:      form.RegistrationEnabled,
		model.RegistrationDefaultGroupSettingKey: form.RegistrationDefaultGroup,
	}
	var snapshot *config.Config
	ret := db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtextextended(?, 0))", "cbctf:settings").Error; err != nil {
			return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
		}
		repo := db.InitSettingRepo(tx)
		keys := make([]string, 0, len(kv))
		for key := range kv {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value := kv[key]
			if ret := repo.Update(key, db.UpdateSettingOptions{Value: &model.SettingValue{V: value}}); !ret.OK {
				return ret
			}
		}
		var ret model.RetVal
		snapshot, ret = repo.ReadSnapshot()
		return ret
	})
	if ret.OK {
		config.Env = snapshot
	}
	return ret
}

func GetPublicSystemConfig() map[string]any {
	return map[string]any{
		"registration_enabled": config.Env.Registration.Enabled,
	}
}
