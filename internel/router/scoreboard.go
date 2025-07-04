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

func GetScoreboard(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	teams, count, ok, msg := service.GetTeamRanking(db.DB.WithContext(ctx), contest.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestFlagRepo := db.InitContestFlagRepo(db.DB.WithContext(ctx))
	contestFlags, _, ok, msg := contestFlagRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads: map[string]db.GetOptions{
			"ContestChallenge": {
				Preloads: map[string]db.GetOptions{
					"Challenge": {},
				},
			},
		},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestChallengeRepo := db.InitContestChallengeRepo(db.DB.WithContext(ctx))
	contestChallenges, _, ok, msg := contestChallengeRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
		Preloads:   map[string]db.GetOptions{"Challenge": {}},
	})
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
		if !middleware.IsAdmin(ctx) && team.Hidden {
			count--
			continue
		}
		teamMap[team.ID] = make(map[string]int)
		for challengeID, _ := range globalMap {
			teamMap[team.ID][challengeID] = 0
		}
		teamFlags, _, ok, msg := teamFlagRepo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID},
			Preloads: map[string]db.GetOptions{
				"ContestFlag": {
					Preloads: map[string]db.GetOptions{
						"ContestChallenge": {
							Preloads: map[string]db.GetOptions{
								"Challenge": {},
							},
						},
					},
				},
			},
		})
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
