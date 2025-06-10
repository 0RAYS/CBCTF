package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
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
	contestFlagRepo := db.InitContestFlagRepo(db.DB.WithContext(ctx))
	contestFlagL, _, ok, msg := contestFlagRepo.ListWithConditions(-1, -1, db.GetOptions{
		{Key: "contest_id", Value: team.ContestID, Op: "and"},
	}, false, "ContestChallenge", "ContestChallenge.Challenge")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	solvedFlagL, _, _ := service.GetTeamSolvedFlags(db.DB.WithContext(ctx), team)
	data := resp.GetTeamResp(team)
	data["solved"] = resp.GetSolvedStateResp(solvedFlagL, contestFlagL)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func GetTeams(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	DB := db.DB.WithContext(ctx)
	contest := middleware.GetContest(ctx)
	teams, count, ok, msg := db.InitTeamRepo(DB).ListWithConditions(form.Limit, form.Offset, db.GetOptions{
		{Key: "contest_id", Value: contest.ID, Op: "and"},
	}, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, team := range teams {
		data = append(data, resp.GetTeamResp(team))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "teams": data}})
}

func GetTeamCaptcha(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": middleware.GetTeam(ctx).Captcha})
}

func GetTeammates(ctx *gin.Context) {
	all := middleware.GetRole(ctx) == "admin"
	data := make([]gin.H, 0)
	for _, user := range middleware.GetTeam(ctx).Users {
		data = append(data, resp.GetUserResp(*user, all))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
	return
}

func GetTeamRanking(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBindQuery(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	var teamsData []struct {
		Team   model.Team
		Solved []model.ContestFlag
	}
	contest := middleware.GetContest(ctx)
	showAll := middleware.GetRole(ctx) == "admin"
	teams, count, ok, msg := service.GetTeamRanking(db.DB.WithContext(ctx), contest.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	for _, team := range teams {
		if !showAll && team.Hidden {
			count--
			continue
		}
		solved, ok, _ := service.GetTeamSolvedFlags(db.DB.WithContext(ctx), team)
		if !ok {
			count--
			continue
		}
		teamsData = append(teamsData, struct {
			Team   model.Team
			Solved []model.ContestFlag
		}{Team: team, Solved: solved})
	}
	contestFlags, _, ok, msg := db.InitContestFlagRepo(db.DB.WithContext(ctx)).ListWithConditions(-1, -1, db.GetOptions{
		{Key: "contest_id", Value: contest.ID, Op: "and"},
	}, false, "ContestChallenge", "ContestChallenge.Challenge")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := resp.GetTeamRankingResp(teamsData, contestFlags, showAll)
	data["count"] = count
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
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
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		ok, msg = service.AdminUpdateTeam(tx, team, form)
	} else {
		var form f.UpdateTeamForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
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
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
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
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.WithContext(ctx).Begin()
	team, ok, msg := service.JoinTeam(tx, contest, user, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.Set("Team", team)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateTeam(ctx *gin.Context) {
	var form f.CreateTeamForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.WithContext(ctx).Begin()
	team, ok, msg := service.CreateTeam(tx, contest, user, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.Set("Team", team)
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
