package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/xlsx"
)

// showHiddenTeams 隐藏队伍是否出现在排行榜中。管理员始终可以看到全部队伍。
func showHiddenTeams(ctx *gin.Context, form dto.LeaderboardForm) bool {
	return middleware.IsFullAccess(ctx) || form.ShowHidden
}

func GetTeamRanking(ctx *gin.Context) {
	var form dto.LeaderboardForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teams, count, ret := service.GetTeamRankingViews(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form.Limit, form.Offset, showHiddenTeams(ctx, form))
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

func GetLeaderboardMatrix(ctx *gin.Context) {
	var form dto.LeaderboardForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teams, count, ret := service.GetLeaderboardMatrixViews(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form.Limit, form.Offset, showHiddenTeams(ctx, form))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := resp.GetLeaderboardMatrixResp(teams)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"teams": data, "count": count}))
}

func GetRankTimeline(ctx *gin.Context) {
	var form dto.LeaderboardForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data, ret := service.GetRankTimelineViews(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), showHiddenTeams(ctx, form))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetRankTimelineResp(data)))
}

// GetUserRanking 全局用户总分榜，登录用户可见。
func GetUserRanking(ctx *gin.Context) {
	var form dto.LeaderboardForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	users, count, ret := service.GetUserRankingViews(db.DB.WithContext(ctx.Request.Context()), form.Limit, form.Offset, form.ShowHidden)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(users))
	for _, user := range users {
		data = append(data, resp.GetUserRankingResp(user))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"users": data, "count": count}))
}

// ExportLeaderboard 导出比赛排行榜为 xlsx，一行一支队伍，含每道题的得分与解题时间。
func ExportLeaderboard(ctx *gin.Context) {
	var form dto.LeaderboardForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	data, ret := service.GetLeaderboardExportView(db.DB.WithContext(ctx.Request.Context()), contest, showHiddenTeams(ctx, form))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	file, err := xlsx.LeaderboardWorkbook(data)
	if err != nil {
		log.Logger.Warningf("Failed to build leaderboard workbook: %s", err)
		resp.JSON(ctx, model.FailRetVal(i18n.Common.UnknownError, "Leaderboard", contest.Name, err))
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Logger.Warningf("Failed to close leaderboard workbook: %s", closeErr)
		}
	}()
	buffer, err := file.WriteToBuffer()
	if err != nil {
		log.Logger.Warningf("Failed to encode leaderboard workbook: %s", err)
		resp.JSON(ctx, model.FailRetVal(i18n.Common.UnknownError, "Leaderboard", contest.Name, err))
		return
	}
	filename := "contest-" + strconv.FormatUint(uint64(contest.ID), 10) + "-leaderboard.xlsx"
	ctx.Writer.Header().Set("File", "true")
	ctx.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}
