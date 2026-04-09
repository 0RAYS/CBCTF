package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetTeamFlags(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ReadFlagEventType)
	teamFlags, ret := service.ListTeamFlagViews(db.DB, middleware.GetTeam(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(teamFlags))
	for _, challenge := range teamFlags {
		flags := make([]gin.H, 0, len(challenge.Flags))
		for _, flag := range challenge.Flags {
			flags = append(flags, gin.H{
				"value":         flag.Value,
				"solved":        flag.Solved,
				"template":      flag.Template,
				"init_score":    flag.InitScore,
				"current_score": flag.CurrentScore,
				"decay":         flag.Decay,
				"min_score":     flag.MinScore,
				"solvers":       flag.Solvers,
			})
		}
		data = append(data, gin.H{
			"name":     challenge.Name,
			"type":     challenge.Type,
			"category": challenge.Category,
			"hidden":   challenge.Hidden,
			"flags":    flags,
		})
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func InitTeamFlag(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.InitChallengeEventType)
	ret := service.InitTeamChallenge(
		db.DB,
		middleware.GetSelf(ctx),
		middleware.GetTeam(ctx),
		middleware.GetContest(ctx),
		middleware.GetChallenge(ctx),
		middleware.GetContestChallenge(ctx),
	)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func ResetTeamFlag(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ResetChallengeEventType)
	ret := service.ResetTeamChallenge(
		db.DB,
		middleware.GetSelf(ctx),
		middleware.GetTeam(ctx),
		middleware.GetContest(ctx),
		middleware.GetChallenge(ctx),
		middleware.GetContestChallenge(ctx),
	)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
