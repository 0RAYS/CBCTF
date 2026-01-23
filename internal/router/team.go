package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
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
	contestFlagL, _, ret := contestFlagRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	solvedFlagL, _ := service.GetTeamSolvedFlags(db.DB, team)
	data := resp.GetTeamResp(team)
	data["solved"] = resp.GetSolvedStateResp(solvedFlagL, contestFlagL)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func GetTeams(ctx *gin.Context) {
	var form f.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	teams, count, ret := db.InitTeamRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"Users": {Selects: []string{"id"}}},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, team := range teams {
		data = append(data, resp.GetTeamResp(team))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "teams": data}))
}

func GetTeamCaptcha(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, model.SuccessRetVal(middleware.GetTeam(ctx).Captcha))
}

func GetTeammates(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	userIDL, ret := db.GetUserIDByTeamID(db.DB, team.ID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	users, _, ret := db.InitUserRepo(db.DB).List(-1, -1, db.GetOptions{Conditions: map[string]any{"id": userIDL}})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, middleware.IsAdmin(ctx)))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
	return
}

func UpdateTeam(ctx *gin.Context) {
	var (
		team = middleware.GetTeam(ctx)
		ret  model.RetVal
	)
	if middleware.IsAdmin(ctx) {
		var form f.AdminUpdateTeamForm
		if ret = form.Bind(ctx); !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
		ret = service.AdminUpdateTeam(db.DB, team, form)
	} else {
		var form f.UpdateTeamForm
		if ret = form.Bind(ctx); !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
		ret = service.UpdateTeam(db.DB, team, form)
	}
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func UpdateCaptcha(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
	captcha := utils.UUID()
	team := middleware.GetTeam(ctx)
	ret := db.InitTeamRepo(db.DB).Update(team.ID, db.UpdateTeamOptions{Captcha: &captcha})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(captcha))
}

func DeleteTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteTeamEventType)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	tx := db.DB.Begin()
	ret := db.InitTeamRepo(tx).Delete(team.ID)
	if !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	tx.Commit()
	prometheus.SubContestActiveTeamsMetrics(contest, 1)
	prometheus.SubContestActiveUsersMetrics(contest, int(db.InitTeamRepo(db.DB).CountAssociation(team, "Users")))
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}

func KickMember(ctx *gin.Context) {
	var form f.KickMemberForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.KickMemberEventType)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	tx := db.DB.Begin()
	if ret := service.LeaveTeam(tx, contest, team, form.UserID); !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Operator": middleware.GetSelfID(ctx)})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func JoinTeam(ctx *gin.Context) {
	var form f.JoinTeamForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.JoinTeamEventType)
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.Begin()
	team, ret := service.JoinTeam(tx, contest, user, form)
	if !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func CreateTeam(ctx *gin.Context) {
	var form f.CreateTeamForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateTeamEventType)
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.Begin()
	team, ret := service.CreateTeam(tx, contest, user, form)
	if !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	tx.Commit()
	go service.CreateTeamFlags(db.DB, team, contest)
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func LeaveTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.LeaveTeamEventType)
	user := middleware.GetSelf(ctx).(model.User)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	tx := db.DB.Begin()
	if ret := service.LeaveTeam(tx, contest, team, user.ID); !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
