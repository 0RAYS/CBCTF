package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func GetTeamRanking(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	teams, count, ret := service.GetTeamRanking(db.DB, contest, form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	repo := db.InitTeamRepo(db.DB)
	data := make([]gin.H, 0)
	isAdmin := middleware.IsAdmin(ctx)
	for _, team := range teams {
		if !middleware.IsAdmin(ctx) && team.Hidden {
			count--
			continue
		}
		solved, ret := db.InitContestFlagRepo(db.DB).GetTeamSolvedContestFlags(team.ID)
		if !ret.OK {
			count--
			continue
		}
		tmp := resp.GetTeamRankingResp(team, solved, contestFlags, isAdmin)
		tmp["users"] = repo.CountAssociation(team, "Users")
		data = append(data, tmp)
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"teams": data, "count": count}))
}

func GetScoreboard(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	teams, count, ret := service.GetTeamRanking(db.DB, contest, form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contestFlagRepo := db.InitContestFlagRepo(db.DB)
	contestFlags, _, ret := contestFlagRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {Preloads: map[string]db.GetOptions{"Challenge": {}}}},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contestChallengeRepo := db.InitContestChallengeRepo(db.DB)
	contestChallenges, _, ret := contestChallengeRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
		Preloads:   map[string]db.GetOptions{"Challenge": {}},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
	teamIDL := make([]uint, 0, len(teams))
	teamFlagRepo := db.InitTeamFlagRepo(db.DB)
	for _, team := range teams {
		if !middleware.IsAdmin(ctx) && team.Hidden {
			count--
			continue
		}
		teamMap[team.ID] = make(map[string]int)
		for challengeID := range globalMap {
			teamMap[team.ID][challengeID] = 0
		}
		teamIDL = append(teamIDL, team.ID)
	}
	teamFlags, ret := teamFlagRepo.GetTeamFlagsWithChallenge(teamIDL...)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	teamFlagsMap := make(map[uint][]db.TeamFlagWithChallenge)
	for _, teamFlag := range teamFlags {
		teamFlagsMap[teamFlag.TeamID] = append(teamFlagsMap[teamFlag.TeamID], teamFlag)
	}
	for teamID := range teamMap {
		for _, teamFlag := range teamFlagsMap[teamID] {
			if teamFlag.ContestChallengeHidden {
				continue
			}
			if teamFlag.Solved {
				teamMap[teamID][teamFlag.ChallengeRandID] += 1
			}
		}
	}
	data := resp.GetScoreboardResp(challengeMap, globalMap, teamMap, teams)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"teams": data, "count": count}))
}

func GetRankTimeline(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	teams, _, ret := service.GetTeamRanking(db.DB, contest, 10, 0)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	teams = slices.DeleteFunc(teams, func(team model.Team) bool {
		if team.Score == 0 {
			return true
		}
		return false
	})
	for i, team := range teams {
		submissions, _, ret := db.InitSubmissionRepo(db.DB).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"solved": true, "team_id": team.ID},
			Selects:    []string{"id", "score", "created_at"},
		})
		if !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		teams[i].Submissions = submissions
	}
	data := resp.GetRankTimelineResp(teams)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}
