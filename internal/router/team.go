package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTeam(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": team})
}

func GetTeams(ctx *gin.Context) {
	var form GetModelsForm
	all := false
	if middleware.GetRole(ctx) == "admin" {
		all = true
	}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	teams, count, ok, msg := db.GetTeams(ctx, middleware.GetContestID(ctx), form.Limit, form.Offset, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "teams": teams}})
}

func JoinTeam(ctx *gin.Context) {
	var form JoinTeamForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contestID := middleware.GetContestID(ctx)
	team, ok, msg := db.GetTeamByName(ctx, form.Name, contestID, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if form.Captcha != team.Captcha {
		ctx.JSON(http.StatusOK, gin.H{"msg": "PasswordError", "data": nil})
		return
	}
	_, msg = db.JoinTeam(ctx, middleware.GetSelfID(ctx), contestID, team.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateTeam(ctx *gin.Context) {
	var form CreateTeamForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contestID := middleware.GetContestID(ctx)
	userID := middleware.GetSelfID(ctx)
	team, ok, msg := db.CreateTeam(ctx, form.Name, userID, contestID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": team})
}

func LeaveTeam(ctx *gin.Context) {
	userID := middleware.GetSelfID(ctx)
	contestID := middleware.GetContestID(ctx)
	team, ok, msg := db.GetTeamByUserID(ctx, userID, contestID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	_, msg = db.LeaveTeam(ctx, userID, contestID, team.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateTeam(ctx *gin.Context) {
	var (
		team model.Team
		ok   bool
		msg  string
		data map[string]interface{}
	)
	if middleware.GetRole(ctx) == "admin" {
		var form AdminUpdateTeamForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		team, ok, msg = db.GetTeamByID(ctx, middleware.GetTeamID(ctx))
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data = utils.Form2Map(form)
	} else if middleware.GetRole(ctx) == "user" {
		var form UpdateTeamForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		team, ok, msg = db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		if team.CaptainID != middleware.GetSelfID(ctx) {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
			return
		}
		data = utils.Form2Map(form)
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	if data["name"].(string) != team.Name {
		if !db.IsUniqueTeamName(data["name"].(string), middleware.GetContestID(ctx)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "TeamNameExists", "data": nil})
			return
		}
	}
	if data["captain_id"].(uint) != team.CaptainID {
		if !db.IsMemberInTeam(team.ID, data["captain_id"].(uint)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "UserNotInTeam", "data": nil})
			return
		}
	}
	_, msg = db.UpdateTeam(ctx, team.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteTeam(ctx *gin.Context) {
	teamID := middleware.GetTeamID(ctx)
	_, msg := db.DeleteTeam(ctx, teamID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func KickMember(ctx *gin.Context) {
	var form KickMemberForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	teamID := middleware.GetTeamID(ctx)
	if !db.IsMemberInTeam(teamID, form.UserID) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNotInTeam", "data": nil})
		return
	}
	_, msg := db.LeaveTeam(ctx, form.UserID, middleware.GetContestID(ctx), teamID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetTeamUsers(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByID(ctx, middleware.GetTeamID(ctx), true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": team.Users})
}
