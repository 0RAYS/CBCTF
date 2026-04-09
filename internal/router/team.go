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
	view := service.GetTeamView(db.DB, team)
	solvedFlagL, contestFlagL, ret := service.GetTeamSolvedFlags(db.DB, middleware.GetContest(ctx), team)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := resp.GetTeamResp(view, middleware.IsFullAccess(ctx))
	data["solved"] = resp.GetSolvedStateResp(solvedFlagL, contestFlagL)
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetTeams(ctx *gin.Context) {
	var form dto.ListTeamForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teams, count, ret := service.ListTeams(db.DB, middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, team := range teams {
		data = append(data, resp.GetTeamResp(team, true))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "teams": data}))
}

func GetTeamCaptcha(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(middleware.GetTeam(ctx).Captcha))
}

func GetTeammates(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	users, ret := service.GetTeammates(db.DB, team, middleware.IsFullAccess(ctx))
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
	ret := service.UpdateTeamCaptcha(db.DB, team, captcha)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(captcha))
}

func DeleteTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteTeamEventType)
	ret := service.DeleteTeamWithTransaction(db.DB, middleware.GetTeam(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
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
	ret := service.KickMember(db.DB, middleware.GetContest(ctx), middleware.GetTeam(ctx), form.UserID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
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
	team, ret := service.JoinTeamWithTransaction(db.DB, middleware.GetContest(ctx), middleware.GetSelf(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
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
	team, ret := service.CreateTeamWithTransaction(db.DB, middleware.GetContest(ctx), middleware.GetSelf(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": team.ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func LeaveTeam(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.LeaveTeamEventType)
	ret := service.LeaveTeamWithTransaction(db.DB, middleware.GetContest(ctx), middleware.GetTeam(ctx), middleware.GetSelf(ctx).ID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Team": middleware.GetTeam(ctx).ID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
