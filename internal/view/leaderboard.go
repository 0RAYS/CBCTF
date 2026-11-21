package view

import (
	"slices"
	"strings"
	"time"

	"CBCTF/internal/model"
)

type LeaderboardSolvedStateView struct {
	Category string
	Solved   int64
	All      int64
}

// BuildSolvedStateViews 按题目分类统计解题数量，两组 flag 都需预加载 Challenge
func BuildSolvedStateViews(solved []model.ChallengeFlag, all []model.ChallengeFlag) []LeaderboardSolvedStateView {
	categories := make(map[uint]string)
	for _, v := range all {
		categories[v.ChallengeID] = v.Challenge.Category
	}
	allCount := make(map[string]int64)
	for _, v := range all {
		allCount[v.Challenge.Category] += 1
	}
	solvedCount := make(map[string]int64)
	for _, flag := range solved {
		solvedCount[categories[flag.ChallengeID]] += 1
	}
	data := make([]LeaderboardSolvedStateView, 0)
	for category, total := range allCount {
		if _, ok := solvedCount[category]; !ok {
			solvedCount[category] = 0
		}
		data = append(data, LeaderboardSolvedStateView{
			Category: category,
			Solved:   solvedCount[category],
			All:      total,
		})
	}
	slices.SortFunc(data, func(a, b LeaderboardSolvedStateView) int {
		return strings.Compare(a.Category, b.Category)
	})
	return data
}

type TeamRankingView struct {
	Team      model.Team
	UserCount int64
	Solved    []LeaderboardSolvedStateView
}

type LeaderboardChallengeSolveView struct {
	ID       string
	Total    int
	Solved   int
	Name     string
	Category string
}

type LeaderboardMatrixTeamView struct {
	Team       model.Team
	UserCount  int64
	Challenges []LeaderboardChallengeSolveView
}

type RankTimelinePointView struct {
	Time  time.Time
	Score float64
}

type RankTimelineTeamView struct {
	ID       uint
	Name     string
	Picture  model.FileURL
	Rank     int
	Score    float64
	Timeline []RankTimelinePointView
}

// LeaderboardChallengeView 排行榜导出中的一道比赛题目。
type LeaderboardChallengeView struct {
	ID       uint
	Name     string
	Category string
}

// LeaderboardTeamView 排行榜导出中的一支队伍。
// Scores 与 SolvedAt 以比赛题目 ID 为键，未解出的题目不会出现在其中。
type LeaderboardTeamView struct {
	Rank        int
	Name        string
	Score       float64
	MemberCount int64
	SolvedCount int
	LastSolve   time.Time
	Scores      map[uint]float64
	SolvedAt    map[uint]time.Time
}

// LeaderboardExportView 排行榜导出的完整数据，一行对应一支队伍。
type LeaderboardExportView struct {
	ContestName string
	Challenges  []LeaderboardChallengeView
	Teams       []LeaderboardTeamView
}
