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
		scores, ret := CalcTeamScores(work, contest.ID, contest.Blood, teams...)
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
	return db.InitTeamRepo(tx).FindAll(db.GetOptions{Conditions: map[string]any{"contest_id": contestID, "banned": false}, Sort: teamRankingSort})
}

// 排行榜排序规则：积分降序、最后得分时间升序、ID 升序。
var (
	teamRankingSort = []string{"score DESC", "last ASC", "id ASC"}
	userRankingSort = []string{"score DESC", "solved DESC", "id ASC"}
)

// teamRankingConditions 队伍排行榜查询条件，隐藏队伍是否参与排名由 showHidden 决定。
func teamRankingConditions(contestID uint, showHidden bool) map[string]any {
	conditions := map[string]any{"contest_id": contestID, "banned": false}
	if !showHidden {
		conditions["hidden"] = false
	}
	return conditions
}

// userRankingConditions 用户排行榜查询条件，隐藏用户是否参与排名由 showHidden 决定。
func userRankingConditions(showHidden bool) map[string]any {
	conditions := map[string]any{"banned": false}
	if !showHidden {
		conditions["hidden"] = false
	}
	return conditions
}

// teamRankingPage 返回队伍排行榜当前页。缓存快照与数据库一致时使用快照，否则回落到数据库。
func teamRankingPage(tx *gorm.DB, contestID uint, showHidden bool, count int64, start, end int) ([]model.Team, model.RetVal) {
	if cached, cacheRet := redis.GetTeamRanking(tx.Statement.Context, contestID); cacheRet.OK {
		visible := cached
		if !showHidden {
			visible = slices.DeleteFunc(cached, func(team model.Team) bool { return team.Hidden })
		}
		if int64(len(visible)) == count {
			return visible[start:end], model.SuccessRetVal()
		}
	}
	teams, _, ret := db.InitTeamRepo(tx).List(end-start, start, db.GetOptions{
		Conditions: teamRankingConditions(contestID, showHidden),
		Sort:       teamRankingSort,
	})
	return teams, ret
}

// GetTeamRanking 返回比赛队伍排名。名次按参与排名的队伍连续编号，隐藏队伍被剔除时不留空缺。
func GetTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int, showHidden bool) ([]model.Team, int64, model.RetVal) {
	conditions := teamRankingConditions(contest.ID, showHidden)
	count, ret := db.InitTeamRepo(tx).Count(db.CountOptions{Conditions: conditions})
	if !ret.OK {
		return nil, 0, ret
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, model.SuccessRetVal()
	}
	teams, ret := teamRankingPage(tx, contest.ID, showHidden, count, start, end)
	if !ret.OK {
		return nil, count, ret
	}
	for i := range teams {
		teams[i].Rank = start + i + 1
	}
	return teams, count, model.SuccessRetVal()
}

func UpdateUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, model.RetVal) {
	users, _, ret := db.InitUserRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"banned": false},
		Sort:       userRankingSort,
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

func GetUserRanking(tx *gorm.DB, limit, offset int, showHidden bool) ([]model.User, int64, model.RetVal) {
	conditions := userRankingConditions(showHidden)
	count, ret := db.InitUserRepo(tx).Count(db.CountOptions{Conditions: conditions})
	if !ret.OK {
		return nil, 0, ret
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, model.SuccessRetVal()
	}
	users, ret := userRankingPage(tx, showHidden, count, start, end)
	if !ret.OK {
		return nil, count, ret
	}
	return users, count, model.SuccessRetVal()
}

// userRankingPage 返回用户排行榜当前页。缓存快照与数据库一致时使用快照，否则回落到数据库。
func userRankingPage(tx *gorm.DB, showHidden bool, count int64, start, end int) ([]model.User, model.RetVal) {
	if cached, cacheRet := redis.GetUserRanking(tx.Statement.Context); cacheRet.OK {
		visible := cached
		if !showHidden {
			visible = slices.DeleteFunc(cached, func(user model.User) bool { return user.Hidden })
		}
		if int64(len(visible)) == count {
			return visible[start:end], model.SuccessRetVal()
		}
	}
	users, _, ret := db.InitUserRepo(tx).List(end-start, start, db.GetOptions{
		Conditions: userRankingConditions(showHidden),
		Sort:       userRankingSort,
	})
	return users, ret
}

// GetUserRankingViews 返回用户总分榜视图，名次按参与排名的用户连续编号。
func GetUserRankingViews(tx *gorm.DB, limit, offset int, showHidden bool) ([]view.UserRankingView, int64, model.RetVal) {
	users, count, ret := GetUserRanking(tx, limit, offset, showHidden)
	if !ret.OK {
		return nil, 0, ret
	}
	start := max(offset, 0)
	views := make([]view.UserRankingView, 0, len(users))
	for i, user := range users {
		views = append(views, view.UserRankingView{User: user, Rank: start + i + 1})
	}
	return views, count, model.SuccessRetVal()
}

func GetTeamRankingViews(tx *gorm.DB, contest model.Contest, limit, offset int, showHidden bool) ([]view.TeamRankingView, int64, model.RetVal) {
	flags, ret := db.InitChallengeFlagRepo(tx).ListForContest(contest.ID)
	if !ret.OK {
		return nil, 0, ret
	}
	teams, count, ret := GetTeamRanking(tx, contest, limit, offset, showHidden)
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
	solvedRows, ret := db.InitChallengeFlagRepo(tx).GetTeamsSolvedChallengeFlags(teamIDs...)
	if !ret.OK {
		return nil, 0, ret
	}
	solvedMap := make(map[uint][]model.ChallengeFlag, len(teamIDs))
	for _, row := range solvedRows {
		solvedMap[row.TeamID] = append(solvedMap[row.TeamID], row.ChallengeFlag)
	}
	views := make([]view.TeamRankingView, 0, len(teams))
	for _, team := range teams {
		views = append(views, view.TeamRankingView{
			Team:      team,
			UserCount: userCountMap[team.ID],
			Solved:    view.BuildSolvedStateViews(solvedMap[team.ID], flags),
		})
	}
	return views, count, model.SuccessRetVal()
}

func GetLeaderboardMatrixViews(tx *gorm.DB, contest model.Contest, limit, offset int, showHidden bool) ([]view.LeaderboardMatrixTeamView, int64, model.RetVal) {
	teams, count, ret := GetTeamRanking(tx, contest, limit, offset, showHidden)
	if !ret.OK {
		return nil, 0, ret
	}
	flags, ret := db.InitChallengeFlagRepo(tx).ListForContest(contest.ID)
	if !ret.OK {
		return nil, 0, ret
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	challengeMap := make(map[string]model.Challenge)
	for _, challenge := range challenges {
		if challenge.Hidden {
			continue
		}
		challengeMap[challenge.RandID] = challenge
	}
	globalMap := make(map[string]int)
	for _, flag := range flags {
		if flag.Challenge.Hidden {
			continue
		}
		globalMap[flag.Challenge.RandID]++
	}
	teamMap := make(map[uint]map[string]int)
	teamIDs := make([]uint, 0, len(teams))
	for _, team := range teams {
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
			if teamFlag.ChallengeHidden {
				continue
			}
			if teamFlag.Solved {
				teamMap[teamID][teamFlag.ChallengeRandID]++
			}
		}
	}
	views := make([]view.LeaderboardMatrixTeamView, 0, len(teamIDs))
	for _, team := range teams {
		challenges := make([]view.LeaderboardChallengeSolveView, 0, len(teamMap[team.ID]))
		for challengeRandID, solvedCount := range teamMap[team.ID] {
			challenge := challengeMap[challengeRandID]
			challenges = append(challenges, view.LeaderboardChallengeSolveView{
				ID:       challengeRandID,
				Total:    globalMap[challengeRandID],
				Solved:   solvedCount,
				Name:     challenge.Name,
				Category: challenge.Category,
			})
		}
		views = append(views, view.LeaderboardMatrixTeamView{
			Team:       team,
			UserCount:  userCountMap[team.ID],
			Challenges: challenges,
		})
	}
	return views, count, model.SuccessRetVal()
}

// GetLeaderboardExportView 汇总排行榜导出所需的排名与分题解题明细。
// 导出是管理员视角的产物，隐藏题目全部导出；隐藏队伍是否出现在导出中由 showHidden 决定。
func GetLeaderboardExportView(tx *gorm.DB, contest model.Contest, showHidden bool) (view.LeaderboardExportView, model.RetVal) {
	teams, _, ret := GetTeamRanking(tx, contest, -1, 0, showHidden)
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}
	scores, ret := CalcTeamScores(tx, contest.ID, contest.Blood, teams...)
	if !ret.OK {
		return view.LeaderboardExportView{}, ret
	}

	slices.SortFunc(challenges, func(a, b model.Challenge) int {
		if order := strings.Compare(a.Category, b.Category); order != 0 {
			return order
		}
		return strings.Compare(a.Name, b.Name)
	})

	visibleTeams := teams
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
		if earliest, ok := times[submission.ChallengeID]; !ok || submission.CreatedAt.Before(earliest) {
			times[submission.ChallengeID] = submission.CreatedAt
		}
	}

	exportView := view.LeaderboardExportView{
		ContestName: contest.Name,
		Challenges:  make([]view.LeaderboardChallengeView, 0, len(challenges)),
		Teams:       make([]view.LeaderboardTeamView, 0, len(visibleTeams)),
	}
	for _, challenge := range challenges {
		exportView.Challenges = append(exportView.Challenges, view.LeaderboardChallengeView{
			ID:       challenge.ID,
			Name:     challenge.Name,
			Category: challenge.Category,
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

func GetRankTimelineViews(tx *gorm.DB, contest model.Contest, showHidden bool) ([]view.RankTimelineTeamView, model.RetVal) {
	teams, _, ret := GetTeamRanking(tx, contest, 10, 0, showHidden)
	if !ret.OK {
		return nil, ret
	}
	// 名次在剔除零分队伍前确定，被省略的队伍不会挤占其他队伍的名次。
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
