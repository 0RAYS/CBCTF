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
	DB := db.DB.WithContext(ctx)
	all := middleware.GetRole(ctx) == "admin"
	flags, _, ok, msg := db.InitFlagRepo(DB).GetByKeyID("contest_id", team.ContestID, -1, -1, true, 3)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	solved, _, _ := service.GetTeamSolved(db.DB.WithContext(ctx), team.ID)
	data := resp.GetTeamResp(team, all)
	data["solved"] = resp.GetSolvedStateResp(solved, flags)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetTeams(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	DB := db.DB.WithContext(ctx)
	contest := middleware.GetContest(ctx)
	teams, count, ok, msg := db.InitTeamRepo(DB).GetAll(contest.ID, form.Limit, form.Offset, false, 0, true, true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, team := range teams {
		data = append(data, resp.GetTeamResp(team, true))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "teams": data}})
}

func GetTeamCaptcha(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetTeam(ctx).Captcha})
}

func GetTeammates(ctx *gin.Context) {
	all := middleware.GetRole(ctx) == "admin"
	data := make([]gin.H, 0)
	for _, user := range middleware.GetTeam(ctx).Users {
		data = append(data, resp.GetUserResp(*user, all))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
	return
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

func KickMember(ctx *gin.Context) {
	var form f.KickMemberForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.LeaveTeam(tx, contest, team, form.UserID)
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

func LeaveTeam(ctx *gin.Context) {
	user := middleware.GetSelf(ctx).(model.User)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.LeaveTeam(tx, contest, team, user.ID)
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
	flags, _, ok, msg := db.InitFlagRepo(DB).GetByKeyID("contest_id", contest.ID, -1, -1, true, 3)
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
