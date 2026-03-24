package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"

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
		resp.JSON(ctx, ret)
		return
	}
	solvedFlagL, _ := db.InitContestFlagRepo(db.DB).GetTeamSolvedContestFlags(team.ID)
	userCount, _ := db.InitTeamRepo(db.DB).CountUsers(team.ID)
	data := resp.GetTeamResp(team, middleware.IsFullAccess(ctx))
	data["users"] = userCount
	data["solved"] = resp.GetSolvedStateResp(solvedFlagL, contestFlagL)
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetTeams(ctx *gin.Context) {
	var form dto.ListTeamForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID},
		Search:     make(map[string]string),
	}
	if form.Name != "" {
		options.Search["name"] = form.Name
	}
	if form.Description != "" {
		options.Search["description"] = form.Description
	}
	teams, count, ret := db.InitTeamRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teamIDL := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDL = append(teamIDL, team.ID)
	}
	userCountMap, ret := db.InitTeamRepo(db.DB).CountUsersMap(teamIDL...)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, team := range teams {
		item := resp.GetTeamResp(team, true)
		item["users"] = userCountMap[team.ID]
		data = append(data, item)
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "teams": data}))
}

func GetTeamCaptcha(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(middleware.GetTeam(ctx).Captcha))
}

func GetTeammates(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	users, ret := db.InitUserRepo(db.DB).GetByTeamID(team.ID, -1, -1)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, middleware.IsFullAccess(ctx)))
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func UpdateTeam(ctx *gin.Context) {
	var (
		team = middleware.GetTeam(ctx)
		ret  model.RetVal
	)
	if middleware.IsFullAccess(ctx) {
		var form dto.AdminUpdateTeamForm
		if ret = dto.Bind(ctx, &form); !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
		ret = service.AdminUpdateTeam(db.DB, team, form)
	} else {
		var form dto.UpdateTeamForm
		if ret = dto.Bind(ctx, &form); !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
		ret = service.UpdateTeam(db.DB, team, form)
	}
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func UpdateCaptcha(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateTeamEventType)
	captcha := utils.UUID()
	team := middleware.GetTeam(ctx)
	ret := db.InitTeamRepo(db.DB).Update(team.ID, db.UpdateTeamOptions{Captcha: &captcha})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(captcha))
}

func DeleteTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteTeamEventType)
	team := middleware.GetTeam(ctx)
	tx := db.DB.Begin()
	ret := db.InitTeamRepo(tx).Delete(team.ID)
	if !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func KickMember(ctx *gin.Context) {
	var form dto.KickMemberForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.KickMemberEventType)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	tx := db.DB.Begin()
	if ret := service.LeaveTeam(tx, contest, team, form.UserID); !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Operator": middleware.GetSelf(ctx).ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func JoinTeam(ctx *gin.Context) {
	var form dto.JoinTeamForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.JoinTeamEventType)
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx)
	tx := db.DB.Begin()
	team, ret := service.JoinTeam(tx, contest, user, form)
	if !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func CreateTeam(ctx *gin.Context) {
	var form dto.CreateTeamForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateTeamEventType)
	contest := middleware.GetContest(ctx)
	user := middleware.GetSelf(ctx)
	tx := db.DB.Begin()
	team, ret := service.CreateTeam(tx, contest, user, form)
	if !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	go service.CreateTeamFlags(db.DB, team, contest)
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func LeaveTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.LeaveTeamEventType)
	user := middleware.GetSelf(ctx)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	tx := db.DB.Begin()
	if ret := service.LeaveTeam(tx, contest, team, user.ID); !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
