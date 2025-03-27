package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTeam(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	flags, ok, msg := service.GetContestFlag(db.DB.WithContext(ctx), team.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	solved, _, _ := service.GetTeamSolved(db.DB.WithContext(ctx), team.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": resp.GetTeamResp(team, solved, flags)})
}

func GetTeamCaptcha(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetTeam(ctx).Captcha})
}

func GetTeammates(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetTeam(ctx).Users})
}

func UpdateTeam(ctx *gin.Context) {
	var (
		team = middleware.GetTeam(ctx)
		tx   = db.DB.WithContext(ctx).Begin()
		ok   bool
		msg  string
	)
	if middleware.GetRole(ctx) == "admin" {
		var form f.AdminUpdateTeamForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		ok, msg = service.AdminUpdateTeam(tx, team, form)
	} else {
		var form f.UpdateTeamForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		ok, msg = service.UpdateTeam(tx, team, form)
	}
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateCaptcha(ctx *gin.Context) {
	captcha := utils.UUID()
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateTeamCaptcha(tx, middleware.GetTeam(ctx), captcha)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": captcha})
}

func DeleteTeam(ctx *gin.Context) {
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.DeleteTeam(tx, middleware.GetTeam(ctx))
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func JoinTeam(ctx *gin.Context) {
	var form f.JoinTeamForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.JoinTeam(tx, contest, user, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateTeam(ctx *gin.Context) {
	var form f.CreateTeamForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.CreateTeam(tx, contest, user, form)
	if !ok {
		tx.Rollback()
		return
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetTeamRanking(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBindQuery(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	DB := db.DB.WithContext(ctx)
	var teamsData []struct {
		Team   model.Team
		Solved []model.Flag
	}
	flags, ok, msg := service.GetContestFlag(DB, contest.ID)
	teams, count, ok, msg := service.GetTeamRanking(DB, contest.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	for _, team := range teams {
		solved, ok, _ := service.GetTeamSolved(DB, team.ID)
		if !ok {
			count--
		}
		teamsData = append(teamsData, struct {
			Team   model.Team
			Solved []model.Flag
		}{Team: team, Solved: solved})
	}
	data := resp.GetTeamRankingResp(teamsData, flags)
	data["count"] = count
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}
