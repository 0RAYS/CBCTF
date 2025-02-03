package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTeam(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetTeam(ctx)})
}

func GetTeamCaptcha(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetTeam(ctx).Captcha})
}

func GetTeams(ctx *gin.Context) {
	var form constants.GetModelsForm
	all := false
	if middleware.GetRole(ctx) == "admin" {
		all = true
	}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	teams, count, ok, msg := db.GetTeams(ctx, middleware.GetContest(ctx).ID, form.Limit, form.Offset, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "teams": teams}})
}

func JoinTeam(ctx *gin.Context) {
	var form constants.JoinTeamForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	team, ok, msg := db.GetTeamByName(ctx, form.Name, contest.ID, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if form.Captcha != team.Captcha {
		ctx.JSON(http.StatusOK, gin.H{"msg": "CaptchaError", "data": nil})
		return
	}
	_, msg = db.JoinTeam(ctx, middleware.GetSelfID(ctx), contest.ID, team.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateTeam(ctx *gin.Context) {
	var form constants.CreateTeamForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	if form.Captcha != contest.Captcha {
		ctx.JSON(http.StatusOK, gin.H{"msg": "CaptchaError", "data": nil})
		return
	}
	userID := middleware.GetSelfID(ctx)
	team, ok, msg := db.CreateTeam(ctx, form, userID, contest.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": team})
}

func LeaveTeam(ctx *gin.Context) {
	userID := middleware.GetSelfID(ctx)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	_, msg := db.LeaveTeam(ctx, userID, contest.ID, team.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateTeam(ctx *gin.Context) {
	var (
		team model.Team
		msg  string
		data map[string]interface{}
	)
	team = middleware.GetTeam(ctx)
	if middleware.GetRole(ctx) == "admin" {
		var form constants.AdminUpdateTeamForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		data = utils.Form2Map(form)
	} else if middleware.GetRole(ctx) == "user" {
		var form constants.UpdateTeamForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		if team.CaptainID != middleware.GetSelfID(ctx) {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
			return
		}
		data = utils.Form2Map(form)
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		return
	}
	if name, ok := data["name"]; ok && name.(string) != team.Name {
		if !db.IsUniqueTeamName(name.(string), middleware.GetContest(ctx).ID) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "TeamNameExists", "data": nil})
			return
		}
	}
	if captainID, ok := data["captain_id"]; ok && captainID.(uint) != team.CaptainID {
		if !db.IsMemberInTeam(team.ID, captainID.(uint)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "UserNotInTeam", "data": nil})
			return
		}
	}
	_, msg = db.UpdateTeam(ctx, team.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteTeam(ctx *gin.Context) {
	_, msg := db.DeleteTeam(ctx, middleware.GetTeam(ctx).ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func KickMember(ctx *gin.Context) {
	var form constants.KickMemberForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	if !db.IsMemberInTeam(team.ID, form.UserID) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNotInTeam", "data": nil})
		return
	}
	_, msg := db.LeaveTeam(ctx, form.UserID, middleware.GetContest(ctx).ID, team.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetTeammates(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetTeam(ctx).Users})
}
