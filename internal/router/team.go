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
	teams, count, ok, msg := db.GetTeams(db.DB.WithContext(ctx), middleware.GetContest(ctx).ID, form.Limit, form.Offset, all)
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
	team, ok, msg := db.GetTeamByName(db.DB.WithContext(ctx), form.Name, contest.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if form.Captcha != team.Captcha {
		ctx.JSON(http.StatusOK, gin.H{"msg": "CaptchaError", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg = db.JoinTeam(tx, middleware.GetSelf(ctx).(model.User), team, contest)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
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
	tx := db.DB.WithContext(ctx).Begin()
	team, ok, msg := db.CreateTeam(tx, form, middleware.GetSelf(ctx).(model.User), contest)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": team})
}

func LeaveTeam(ctx *gin.Context) {
	user := middleware.GetSelf(ctx).(model.User)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.LeaveTeam(tx, user, team, contest)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
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
		if !db.IsUniqueTeamName(db.DB.WithContext(ctx), name.(string), middleware.GetContest(ctx).ID) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "TeamNameExists", "data": nil})
			return
		}
	}
	if captainID, ok := data["captain_id"]; ok && captainID.(uint) != team.CaptainID {
		if !db.IsMemberInTeam(db.DB.WithContext(ctx), team.ID, captainID.(uint)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "UserNotInTeam", "data": nil})
			return
		}
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.UpdateTeam(tx, team.ID, data)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteTeam(ctx *gin.Context) {
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.DeleteTeam(tx, middleware.GetTeam(ctx))
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func KickMember(ctx *gin.Context) {
	var form constants.KickMemberForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	if !db.IsMemberInTeam(db.DB.WithContext(ctx), team.ID, form.UserID) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNotInTeam", "data": nil})
		return
	}
	user, ok, msg := db.GetUserByID(db.DB.WithContext(ctx), form.UserID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg = db.LeaveTeam(tx, user, team, middleware.GetContest(ctx))
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetTeammates(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetTeam(ctx).Users})
}
