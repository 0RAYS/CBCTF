package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetScoreboard(ctx *gin.Context) {
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
	showAll := middleware.GetRole(ctx) == "admin"
	contest := middleware.GetContest(ctx)
	teams, count, ok, msg := service.GetTeamRanking(db.DB.WithContext(ctx), contest.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestFlagRepo := db.InitContestFlagRepo(db.DB.WithContext(ctx))
	contestFlags, _, ok, msg := contestFlagRepo.ListWithConditions(-1, -1, db.GetOptions{
		{Key: "contest_id", Value: contest.ID, Op: "and"},
	}, false, "ContestChallenge", "ContestChallenge.Challenge")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestChallengeRepo := db.InitContestChallengeRepo(db.DB.WithContext(ctx))
	contestChallenges, _, ok, msg := contestChallengeRepo.ListWithConditions(-1, -1, db.GetOptions{
		{Key: "contest_id", Value: contest.ID, Op: "and"},
		{Key: "hidden", Value: false, Op: "and"},
	}, false, "Challenge")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	challengeMap := make(map[string]model.Challenge)
	for _, contestChallenge := range contestChallenges {
		if contestChallenge.Hidden {
			continue
		}
		challengeMap[contestChallenge.Challenge.RandID] = contestChallenge.Challenge
	}
	globalMap := make(map[string]int)
	for _, contestFlag := range contestFlags {
		if contestFlag.ContestChallenge.Hidden {
			continue
		}
		globalMap[contestFlag.ContestChallenge.Challenge.RandID] += 1
	}
	teamMap := make(map[uint]map[string]int)
	teamFlagRepo := db.InitTeamFlagRepo(db.DB.WithContext(ctx))
	for _, team := range teams {
		if !showAll && team.Hidden {
			count--
			continue
		}
		teamMap[team.ID] = make(map[string]int)
		for challengeID, _ := range globalMap {
			teamMap[team.ID][challengeID] = 0
		}
		teamFlags, _, ok, msg := teamFlagRepo.ListWithConditions(-1, -1, db.GetOptions{
			{Key: "team_id", Value: team.ID, Op: "and"},
		}, false, "ContestFlag", "ContestFlag.ContestChallenge", "ContestFlag.ContestChallenge.Challenge")
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		for _, teamFlag := range teamFlags {
			if teamFlag.ContestFlag.ContestChallenge.Hidden {
				continue
			}
			if teamFlag.Solved {
				teamMap[team.ID][teamFlag.ContestFlag.ContestChallenge.Challenge.RandID] += 1
			}
		}
	}
	data := resp.GetScoreboardResp(challengeMap, globalMap, teamMap, teams)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"teams": data, "count": count}})
}
