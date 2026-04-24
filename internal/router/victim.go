package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func StartVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	user := middleware.GetSelf(ctx)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	ret := service.StartVictim(db.DB, user.ID, team.ID, contest.ID, contestChallenge.ID, challenge.ID)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func ExtendVictimDuration(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ExtendVictimEventType)
	victim, ret := service.ExtendVictimDuration(db.DB, middleware.GetTeam(ctx), middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"target":    victim.RemoteAddr(),
		"duration":  victim.Duration.Seconds(),
		"remaining": victim.Remaining().Seconds(),
		"status":    "Running",
	}))
}

func StopVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	if ret := service.StopAliveVictim(db.DB, middleware.GetTeam(ctx), middleware.GetChallenge(ctx)); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func GetVictimHistories(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	victims, count, ret := service.ListVictimHistories(db.DB, middleware.GetTeam(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, victim := range victims {
		data = append(data, resp.GetVictimResp(victim))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"victims": data, "count": count}))
}

func GetVictims(ctx *gin.Context) {
	var form dto.GetVictimsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	victims, count, total, _ := service.GetVictims(db.DB, contest, form)
	data := make([]gin.H, 0)
	for _, victim := range victims {
		info := resp.GetVictimResp(victim)
		info["remote"] = victim.RemoteAddr()
		info["remaining"] = victim.Remaining().Seconds()
		info["team"] = victim.Team.Name
		info["user"] = victim.User.Name
		info["challenge"] = victim.ContestChallenge.Name
		data = append(data, info)
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"victims": data, "count": total, "running": count}))
}

func StartVictims(ctx *gin.Context) {
	var form dto.StartVictimsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	contest := middleware.GetContest(ctx)
	ret := service.StartVictims(db.DB, contest, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func StopVictims(ctx *gin.Context) {
	var form dto.StopVictimsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	go service.StopVictims(db.DB, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
