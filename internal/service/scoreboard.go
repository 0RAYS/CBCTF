package service

import (
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
)

func UpdateTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int) ([]model.Team, int64, model.RetVal) {
	var snapshot []model.Team
	ret := db.WithTransactionDB(tx, func(work *gorm.DB) model.RetVal {
		repo := db.InitTeamRepo(work.Clauses(clause.Locking{Strength: "UPDATE"}))
		teams, ret := repo.FindAll(db.GetOptions{Conditions: map[string]any{"contest_id": contest.ID, "banned": false}, Sort: []string{"id"}})
		if !ret.OK {
			return ret
		}
		scores, ret := CalcTeamScores(work, contest.Blood, teams...)
		if !ret.OK {
			return ret
		}
		for _, team := range teams {
			if ret := db.InitTeamRepo(work).Update(team.ID, db.UpdateTeamOptions{Score: new(scores[team.ID].Total)}); !ret.OK {
				return ret
			}
		}
		var ret2 model.RetVal
		snapshot, ret2 = rankedTeams(work, contest.ID)
		if !ret2.OK {
			return ret2
		}
		for i := range snapshot {
			snapshot[i].Rank = i + 1
			if ret := db.InitTeamRepo(work).Update(snapshot[i].ID, db.UpdateTeamOptions{Rank: new(i + 1)}); !ret.OK {
				return ret
			}
		}
		snapshot, ret2 = rankedTeams(work, contest.ID)
		return ret2
	})
	if !ret.OK {
		return nil, 0, ret
	}
	if ret = redis.UpdateTeamRanking(tx.Statement.Context, contest.ID, snapshot); !ret.OK {
		return nil, 0, ret
	}
	start, end := utils.TidyPaginate(len(snapshot), limit, offset)
	return snapshot[start:end], int64(len(snapshot)), model.SuccessRetVal()
}

func rankedTeams(tx *gorm.DB, contestID uint) ([]model.Team, model.RetVal) {
	return db.InitTeamRepo(tx).FindAll(db.GetOptions{Conditions: map[string]any{"contest_id": contestID, "banned": false}, Sort: []string{"score DESC", "last ASC", "id ASC"}})
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
	teams, ret := redis.GetTeamRanking(tx.Statement.Context, contest.ID, int64(start), int64(end-1))
	if !ret.OK || len(teams) != end-start {
		// Cache failure is not a reason to mutate scores or recursively rebuild.
		var dbRet model.RetVal
		teams, _, dbRet = db.InitTeamRepo(tx).List(end-start, start, db.GetOptions{Conditions: map[string]any{"contest_id": contest.ID, "banned": false}, Sort: []string{"score DESC", "last ASC", "id ASC"}})
		if !dbRet.OK {
			return nil, count, dbRet
		}
	}
	for i := range teams {
		teams[i].Rank = start + i + 1
	}
	return teams, count, model.SuccessRetVal()
}

func UpdateUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, model.RetVal) {
	users, _, ret := db.InitUserRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"banned": false},
		Sort:       []string{"score DESC", "solved DESC", "id ASC"},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	if ret = redis.UpdateUserRanking(tx.Statement.Context, users); !ret.OK {
		return nil, 0, ret
	}
	start, end := utils.TidyPaginate(len(users), limit, offset)
	return users[start:end], int64(len(users)), model.SuccessRetVal()
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
	users, ret := redis.GetUserRanking(tx.Statement.Context, int64(start), int64(end-1))
	if !ret.OK || len(users) != end-start {
		var dbRet model.RetVal
		users, _, dbRet = db.InitUserRepo(tx).List(end-start, start, db.GetOptions{Conditions: map[string]any{"banned": false}, Sort: []string{"score DESC", "solved DESC", "id ASC"}})
		if !dbRet.OK {
			return nil, count, dbRet
		}
	}
	return users, count, model.SuccessRetVal()
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
			Solved:    view.BuildSolvedStateViews(solvedMap[team.ID], contestFlags),
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
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	challengeMap := make(map[string]model.ContestChallenge)
	for _, contestChallenge := range contestChallenges {
		if contestChallenge.Hidden {
			continue
		}
		challengeMap[contestChallenge.RandID] = contestChallenge
	}
	globalMap := make(map[string]int)
	for _, contestFlag := range contestFlags {
		if contestFlag.ContestChallenge.Hidden {
			continue
		}
		globalMap[contestFlag.ContestChallenge.RandID]++
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

// GetLeaderboardExportView 汇总排行榜导出所需的排名与分题解题明细。
// 非管理员看不到隐藏队伍与隐藏题目，名次沿用积分榜中的位置。
func GetLeaderboardExportView(tx *gorm.DB, contest model.Contest, admin bool) (view.LeaderboardExportView, model.RetVal) {
	teams, _, ret := GetTeamRanking(tx, contest, -1, 0)
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}
	scores, ret := CalcTeamScores(tx, contest.Blood, teams...)
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}

	visibleChallenges := make([]model.ContestChallenge, 0, len(contestChallenges))
	for _, contestChallenge := range contestChallenges {
		if contestChallenge.Hidden && !admin {
			continue
		}
		visibleChallenges = append(visibleChallenges, contestChallenge)
	}
	slices.SortFunc(visibleChallenges, func(a, b model.ContestChallenge) int {
		if order := strings.Compare(a.Category, b.Category); order != 0 {
			return order
		}
		return strings.Compare(a.Name, b.Name)
	})

	visibleTeams := make([]model.Team, 0, len(teams))
	for _, team := range teams {
		if team.Hidden && !admin {
			continue
		}
		visibleTeams = append(visibleTeams, team)
	}
	teamIDs := make([]uint, 0, len(visibleTeams))
	for _, team := range visibleTeams {
		teamIDs = append(teamIDs, team.ID)
	}
	memberCounts, ret := db.InitTeamRepo(tx).CountUsersMap(teamIDs...)
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}
	submissions, ret := db.InitSubmissionRepo(tx).ListSolvedByTeamID(teamIDs...)
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}
	solvedAt := make(map[uint]map[uint]time.Time, len(teamIDs))
	for _, submission := range submissions {
		if submission.CreatedAt.IsZero() {
			continue
		}
		times, ok := solvedAt[submission.TeamID]
		if !ok {
			times = make(map[uint]time.Time)
			solvedAt[submission.TeamID] = times
		}
		if earliest, ok := times[submission.ContestChallengeID]; !ok || submission.CreatedAt.Before(earliest) {
			times[submission.ContestChallengeID] = submission.CreatedAt
		}
	}

	exportView := view.LeaderboardExportView{
		ContestName: contest.Name,
		Challenges:  make([]view.LeaderboardChallengeView, 0, len(visibleChallenges)),
		Teams:       make([]view.LeaderboardTeamView, 0, len(visibleTeams)),
	}
	for _, contestChallenge := range visibleChallenges {
		exportView.Challenges = append(exportView.Challenges, view.LeaderboardChallengeView{
			ID:       contestChallenge.ID,
			Name:     contestChallenge.Name,
			Category: contestChallenge.Category,
		})
	}
	for _, team := range visibleTeams {
		teamScore := scores[team.ID]
		exportView.Teams = append(exportView.Teams, view.LeaderboardTeamView{
			Rank:        team.Rank,
			Name:        team.Name,
			Score:       team.Score,
			MemberCount: memberCounts[team.ID],
			SolvedCount: len(teamScore.ByChallenge),
			LastSolve:   team.Last,
			Scores:      teamScore.ByChallenge,
			SolvedAt:    solvedAt[team.ID],
		})
	}
	return exportView, model.SuccessRetVal()
}

func GetRankTimelineViews(tx *gorm.DB, contest model.Contest) ([]view.RankTimelineTeamView, model.RetVal) {
	teams, _, ret := GetTeamRanking(tx, contest, 10, 0)
	if !ret.OK {
		return nil, ret
	}
	// Redis returns current ranking order, but each cached team may carry an old Rank.
	// Assign positions before filtering so omitted zero-score teams do not shift ranks.
	for i := range teams {
		teams[i].Rank = i + 1
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
