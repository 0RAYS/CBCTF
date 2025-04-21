package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/k8s"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetUsageStatus(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	usage := middleware.GetUsage(ctx)
	DB := db.DB.WithContext(ctx)
	data := gin.H{
		"attempts": service.CountAttempts(DB, team, usage),
		"init":     service.IsGenerated(DB, team, usage),
		"solved":   service.IsSolved(DB, team, usage),
		"remote":   service.GetVictimStatus(DB, team, usage),
		"file": func() string {
			if _, err := os.Stat(usage.Challenge.AttachmentPath(team.ID)); err != nil {
				return ""
			}
			return usage.Challenge.AttachmentPath(team.ID)
		}(),
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetUsage(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	data := resp.GetUsageResp(usage, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetUsages(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	var (
		all     = middleware.GetRole(ctx) == "admin"
		DB      = db.DB.WithContext(ctx)
		contest = middleware.GetContest(ctx)
	)
	usages, count, ok, msg := db.InitUsageRepo(DB).
		GetAll(contest.ID, form.Limit, form.Offset, all, "Challenge", "Flags")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, usage := range usages {
		tmp := resp.GetUsageResp(usage, all)
		if !all {
			team := middleware.GetTeam(ctx)
			tmp["attempts"] = service.CountAttempts(DB, team, usage)
			tmp["init"] = service.IsGenerated(DB, team, usage)
			tmp["solved"] = service.IsSolved(DB, team, usage)
			tmp["remote"] = service.GetVictimStatus(DB, team, usage)
			tmp["file"] = func() string {
				if _, err := os.Stat(usage.Challenge.AttachmentPath(team.ID)); err != nil {
					return ""
				}
				return usage.Challenge.AttachmentPath(team.ID)
			}()
		}
		data = append(data, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"challenges": data, "count": count}})
}

func AddUsage(ctx *gin.Context) {
	var form f.CreateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	DB := db.DB.WithContext(ctx)
	usages, failed, _, _ := service.CreateUsage(DB, middleware.GetContest(ctx), form)
	data := make([]gin.H, 0)
	for _, usage := range usages {
		data = append(data, resp.GetUsageResp(usage, true))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"usages": data, "failed": failed}})
}

func UpdateUsage(ctx *gin.Context) {
	var form f.UpdateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	usage := middleware.GetUsage(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateUsage(tx, usage, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GenerateTeamUsage(reset bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			team    = middleware.GetTeam(ctx)
			usage   = middleware.GetUsage(ctx)
			answers = make([]model.Answer, 0)
			ok      bool
			msg     string
		)
		if ok, err := redis.CheckChallengeInit(team.ID, usage.ChallengeID); ok || err != nil {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "TooQuick", "data": nil})
			return
		}
		_ = redis.RecordChallengeInit(team.ID, usage.ChallengeID)
		tx := db.DB.WithContext(ctx).Begin()
		if reset {
			answers, ok, msg = service.ResetAnswer(tx, team, usage)
		} else {
			answers, ok, msg = service.GenerateAnswer(tx, team, usage)
		}
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		switch usage.Challenge.Type {
		case model.DynamicChallenge:
			ok, msg = k8s.GenerateAttachment(usage, team, answers)
		case model.PodsChallenge:
			// 不考虑失败
			go service.StopVictim(db.DB.WithContext(ctx), team, usage)
			ok, msg = true, "Success"
		default:
			ok, msg = true, "Success"
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	}
}

func RemoveUsage(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.DeleteUsage(tx, usage)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
