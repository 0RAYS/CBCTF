package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeam(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestFlagRepo := db.InitContestFlagRepo(db.DB)
	contestFlagL, _, ok, msg := contestFlagRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	solvedFlagL, _, _ := service.GetTeamSolvedFlags(db.DB, team)
	data := resp.GetTeamResp(team)
	data["solved"] = resp.GetSolvedStateResp(solvedFlagL, contestFlagL)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func GetTeams(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	teams, count, ok, msg := db.InitTeamRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"Users": {Selects: []string{"id"}}},
	})
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
	team := middleware.GetTeam(ctx)
	userIDL, ok, msg := db.GetUserIDByTeamID(db.DB, team.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	users, _, ok, msg := db.InitUserRepo(db.DB).List(-1, -1, db.GetOptions{Conditions: map[string]any{"id": userIDL}})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, middleware.IsAdmin(ctx)))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
	return
}

func UpdateTeam(ctx *gin.Context) {
	var (
		team = middleware.GetTeam(ctx)
		ok   bool
		msg  string
	)
	if middleware.IsAdmin(ctx) {
		var form f.AdminUpdateTeamForm
		if ok, msg = form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
		ok, msg = service.AdminUpdateTeam(db.DB, team, form)
	} else {
		var form f.UpdateTeamForm
		if ok, msg = form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
		ok, msg = service.UpdateTeam(db.DB, team, form)
	}
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateCaptcha(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
	captcha := utils.UUID()
	team := middleware.GetTeam(ctx)
	ok, msg := db.InitTeamRepo(db.DB).Update(team.ID, db.UpdateTeamOptions{Captcha: &captcha})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": captcha})
}

func DeleteTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteTeamEventType)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	tx := db.DB.Begin()
	ok, msg := db.InitTeamRepo(tx).Delete(team.ID)
	if ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	prometheus.SubContestActiveTeamsMetrics(contest, 1)
	prometheus.SubContestActiveUsersMetrics(contest, int(db.InitTeamRepo(db.DB).CountAssociation(team, "Users")))
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func KickMember(ctx *gin.Context) {
	var form f.KickMemberForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.KickMemberEventType)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	tx := db.DB.Begin()
	ok, msg := service.LeaveTeam(tx, contest, team, form.UserID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Operator": middleware.GetSelfID(ctx)})
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func JoinTeam(ctx *gin.Context) {
	var form f.JoinTeamForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.JoinTeamEventType)
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.Begin()
	team, ok, msg := service.JoinTeam(tx, contest, user, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateTeam(ctx *gin.Context) {
	var form f.CreateTeamForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateTeamEventType)
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.Begin()
	team, ok, msg := service.CreateTeam(tx, contest, user, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	go service.CreateTeamFlags(db.DB, team, contest)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func LeaveTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.LeaveTeamEventType)
	user := middleware.GetSelf(ctx).(model.User)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	tx := db.DB.Begin()
	ok, msg := service.LeaveTeam(tx, contest, team, user.ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
