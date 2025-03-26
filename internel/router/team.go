package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	flags, ok, msg := service.GetContestFlag(DB, contest.ID)
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
