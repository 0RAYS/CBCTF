package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"slices"

	"github.com/gin-gonic/gin"
)

func GetTeamRanking(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teams, count, ret := service.GetTeamRanking(db.DB, contest, form.Limit, form.Offset)
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

	solvedRows, ret := db.InitContestFlagRepo(db.DB).GetTeamsSolvedContestFlags(teamIDL...)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	solvedMap := make(map[uint][]model.ContestFlag, len(teamIDL))
	for _, row := range solvedRows {
		solvedMap[row.TeamID] = append(solvedMap[row.TeamID], row.ContestFlag)
	}

	data := make([]gin.H, 0)
	for _, team := range teams {
		if !middleware.IsFullAccess(ctx) && team.Hidden {
			count--
			continue
		}
		tmp := resp.GetTeamRankingResp(team, solvedMap[team.ID], contestFlags, middleware.IsFullAccess(ctx))
		tmp["users"] = userCountMap[team.ID]
		data = append(data, tmp)
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"teams": data, "count": count}))
}

func GetScoreboard(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	teams, count, ret := service.GetTeamRanking(db.DB, contest, form.Limit, form.Offset)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contestFlagRepo := db.InitContestFlagRepo(db.DB)
	contestFlags, _, ret := contestFlagRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {Preloads: map[string]db.GetOptions{"Challenge": {}}}},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contestChallengeRepo := db.InitContestChallengeRepo(db.DB)
	contestChallenges, _, ret := contestChallengeRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
		Preloads:   map[string]db.GetOptions{"Challenge": {}},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
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
		if !middleware.IsFullAccess(ctx) && team.Hidden {
			count--
			continue
		}
		teamMap[team.ID] = make(map[string]int)
		for challengeID := range globalMap {
			teamMap[team.ID][challengeID] = 0
		}
		teamIDL = append(teamIDL, team.ID)
	}

	userCountMap, ret := db.InitTeamRepo(db.DB).CountUsersMap(teamIDL...)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teamFlags, ret := teamFlagRepo.GetTeamFlagsWithChallenge(teamIDL...)
	if !ret.OK {
		resp.JSON(ctx, ret)
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
	data := resp.GetScoreboardResp(challengeMap, globalMap, teamMap, teams, userCountMap)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"teams": data, "count": count}))
}

func GetRankTimeline(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	teams, _, ret := service.GetTeamRanking(db.DB, contest, 10, 0)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	teams = slices.DeleteFunc(teams, func(team model.Team) bool {
		if team.Score == 0 {
			return true
		}
		return false
	})
	teamIDL := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDL = append(teamIDL, team.ID)
	}
	submissions, ret := db.InitSubmissionRepo(db.DB).ListSolvedByTeamID(teamIDL...)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	timelineMap := make(map[uint][]model.Submission, len(teamIDL))
	for _, submission := range submissions {
		timelineMap[submission.TeamID] = append(timelineMap[submission.TeamID], submission)
	}
	for i := range teams {
		teams[i].Submissions = timelineMap[teams[i].ID]
	}
	data := resp.GetRankTimelineResp(teams)
	resp.JSON(ctx, model.SuccessRetVal(data))
}
