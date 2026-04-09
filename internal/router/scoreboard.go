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

func GetTeamRanking(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teams, count, ret := service.GetTeamRankingViews(db.DB, middleware.GetContest(ctx), form.Limit, form.Offset, middleware.IsFullAccess(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(teams))
	for _, team := range teams {
		item := resp.GetTeamRankingResp(team, middleware.IsFullAccess(ctx))
		item["users"] = team.UserCount
		data = append(data, item)
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"teams": data, "count": count}))
}

func GetScoreboard(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teams, count, ret := service.GetScoreboardViews(db.DB, middleware.GetContest(ctx), form.Limit, form.Offset, middleware.IsFullAccess(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := resp.GetScoreboardResp(teams)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"teams": data, "count": count}))
}

func GetRankTimeline(ctx *gin.Context) {
	data, ret := service.GetRankTimelineViews(db.DB, middleware.GetContest(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetRankTimelineResp(data)))
}
