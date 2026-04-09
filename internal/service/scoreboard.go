package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func UpdateTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int) ([]model.Team, int64, model.RetVal) {
	var (
		repo          = db.InitTeamRepo(tx)
		teams, _, ret = repo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
		})
	)
	if !ret.OK {
		return nil, 0, ret
	}
	scoreMap, ret := CalcTeamScores(tx, contest.Blood, teams...)
	if !ret.OK {
		return nil, 0, ret
	}
	for i := range teams {
		teams[i].Score = scoreMap[teams[i].ID]
	}
	if ret = redis.UpdateTeamRanking(contest.ID, teams); !ret.OK {
		return nil, 0, ret
	}
	teams, count, ret := GetTeamRanking(tx, contest, limit, offset)
	if !ret.OK {
		return nil, 0, ret
	}
	for i, team := range teams {
		teams[i].Rank = i + 1
		repo.Update(team.ID, db.UpdateTeamOptions{Score: &team.Score, Rank: new(i + 1)})
	}
	return teams, count, model.SuccessRetVal()
}

func GetTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int) ([]model.Team, int64, model.RetVal) {
	count, ret := db.InitTeamRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, model.SuccessRetVal()
	}
	teams, ret := redis.GetTeamRanking(contest.ID, int64(start), int64(end-1))
	if !ret.OK || (end-start > 0 && len(teams) == 0 && count > 0) {
		return UpdateTeamRanking(tx, contest, limit, offset)
	}
	return teams, count, model.SuccessRetVal()
}

func UpdateUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, model.RetVal) {
	users, _, ret := db.InitUserRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"banned": false},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	if ret = redis.UpdateUserRanking(users); !ret.OK {
		return nil, 0, ret
	}
	return GetUserRanking(tx, limit, offset)
}

func GetUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, model.RetVal) {
	count, ret := db.InitUserRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"banned": false},
	})
	if !ret.OK {
		return nil, count, ret
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, model.SuccessRetVal()
	}
	users, ret := redis.GetUserRanking(int64(start), int64(end-1))
	if !ret.OK || (end-start > 0 && len(users) == 0 && count > 0) {
		return UpdateUserRanking(tx, limit, offset)
	}
	return users, count, model.SuccessRetVal()
}

func buildSolvedStateViews(solved []model.ContestFlag, all []model.ContestFlag) []view.ScoreboardSolvedStateView {
	categories := make(map[uint]string)
	for _, v := range all {
		categories[v.ContestChallengeID] = v.ContestChallenge.Category
	}
	allCount := make(map[string]int64)
	for _, v := range all {
		allCount[v.ContestChallenge.Category] += 1
	}
	solvedCount := make(map[string]int64)
	for _, flag := range solved {
		solvedCount[categories[flag.ContestChallengeID]] += 1
	}
	data := make([]view.ScoreboardSolvedStateView, 0)
	for category, total := range allCount {
		if _, ok := solvedCount[category]; !ok {
			solvedCount[category] = 0
		}
		data = append(data, view.ScoreboardSolvedStateView{
			Category: category,
			Solved:   solvedCount[category],
			All:      total,
		})
	}
	slices.SortFunc(data, func(a, b view.ScoreboardSolvedStateView) int {
		return strings.Compare(a.Category, b.Category)
	})
	return data
}

func GetTeamRankingViews(tx *gorm.DB, contest model.Contest, limit, offset int, admin bool) ([]view.TeamRankingView, int64, model.RetVal) {
	contestFlags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	teams, count, ret := GetTeamRanking(tx, contest, limit, offset)
	if !ret.OK {
		return nil, 0, ret
	}
	teamIDs := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDs = append(teamIDs, team.ID)
	}
	userCountMap, ret := db.InitTeamRepo(tx).CountUsersMap(teamIDs...)
	if !ret.OK {
		return nil, 0, ret
	}
	solvedRows, ret := db.InitContestFlagRepo(tx).GetTeamsSolvedContestFlags(teamIDs...)
	if !ret.OK {
		return nil, 0, ret
	}
	solvedMap := make(map[uint][]model.ContestFlag, len(teamIDs))
	for _, row := range solvedRows {
		solvedMap[row.TeamID] = append(solvedMap[row.TeamID], row.ContestFlag)
	}
	views := make([]view.TeamRankingView, 0, len(teams))
	for _, team := range teams {
		if !admin && team.Hidden {
			count--
			continue
		}
		views = append(views, view.TeamRankingView{
			Team:      team,
			UserCount: userCountMap[team.ID],
			Solved:    buildSolvedStateViews(solvedMap[team.ID], contestFlags),
		})
	}
	return views, count, model.SuccessRetVal()
}

func GetScoreboardViews(tx *gorm.DB, contest model.Contest, limit, offset int, admin bool) ([]view.ScoreboardTeamView, int64, model.RetVal) {
	teams, count, ret := GetTeamRanking(tx, contest, limit, offset)
	if !ret.OK {
		return nil, 0, ret
	}
	contestFlags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {Preloads: map[string]db.GetOptions{"Challenge": {}}}},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
		Preloads:   map[string]db.GetOptions{"Challenge": {}},
	})
	if !ret.OK {
		return nil, 0, ret
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
		globalMap[contestFlag.ContestChallenge.Challenge.RandID]++
	}
	teamMap := make(map[uint]map[string]int)
	teamIDs := make([]uint, 0, len(teams))
	for _, team := range teams {
		if !admin && team.Hidden {
			count--
			continue
		}
		teamMap[team.ID] = make(map[string]int)
		for challengeID := range globalMap {
			teamMap[team.ID][challengeID] = 0
		}
		teamIDs = append(teamIDs, team.ID)
	}
	userCountMap, ret := db.InitTeamRepo(tx).CountUsersMap(teamIDs...)
	if !ret.OK {
		return nil, 0, ret
	}
	teamFlags, ret := db.InitTeamFlagRepo(tx).GetTeamFlagsWithChallenge(teamIDs...)
	if !ret.OK {
		return nil, 0, ret
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
				teamMap[teamID][teamFlag.ChallengeRandID]++
			}
		}
	}
	views := make([]view.ScoreboardTeamView, 0, len(teamIDs))
	for _, team := range teams {
		if !admin && team.Hidden {
			continue
		}
		challenges := make([]view.ScoreboardChallengeSolveView, 0, len(teamMap[team.ID]))
		for challengeRandID, solvedCount := range teamMap[team.ID] {
			challenge := challengeMap[challengeRandID]
			challenges = append(challenges, view.ScoreboardChallengeSolveView{
				ID:       challengeRandID,
				Total:    globalMap[challengeRandID],
				Solved:   solvedCount,
				Name:     challenge.Name,
				Category: challenge.Category,
			})
		}
		views = append(views, view.ScoreboardTeamView{
			Team:       team,
			UserCount:  userCountMap[team.ID],
			Challenges: challenges,
		})
	}
	return views, count, model.SuccessRetVal()
}

func GetRankTimelineViews(tx *gorm.DB, contest model.Contest) ([]view.RankTimelineTeamView, model.RetVal) {
	teams, _, ret := GetTeamRanking(tx, contest, 10, 0)
	if !ret.OK {
		return nil, ret
	}
	teams = slices.DeleteFunc(teams, func(team model.Team) bool {
		return team.Score == 0
	})
	teamIDs := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDs = append(teamIDs, team.ID)
	}
	submissions, ret := db.InitSubmissionRepo(tx).ListSolvedByTeamID(teamIDs...)
	if !ret.OK {
		return nil, ret
	}
	timelineMap := make(map[uint][]view.RankTimelinePointView, len(teamIDs))
	for _, submission := range submissions {
		timelineMap[submission.TeamID] = append(timelineMap[submission.TeamID], view.RankTimelinePointView{
			Time:  submission.CreatedAt,
			Score: submission.Score,
		})
	}
	views := make([]view.RankTimelineTeamView, 0, len(teams))
	for _, team := range teams {
		views = append(views, view.RankTimelineTeamView{
			ID:       team.ID,
			Name:     team.Name,
			Picture:  team.Picture,
			Rank:     team.Rank,
			Score:    team.Score,
			Timeline: timelineMap[team.ID],
		})
	}
	return views, model.SuccessRetVal()
}
