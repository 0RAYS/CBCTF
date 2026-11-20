package view

import (
	"time"

	"CBCTF/internal/model"
)

type ScoreboardSolvedStateView struct {
	Category string
	Solved   int64
	All      int64
}

type TeamRankingView struct {
	Team      model.Team
	UserCount int64
	Solved    []ScoreboardSolvedStateView
}

type ScoreboardChallengeSolveView struct {
	ID       string
	Total    int
	Solved   int
	Name     string
	Category string
}

type ScoreboardTeamView struct {
	Team       model.Team
	UserCount  int64
	Challenges []ScoreboardChallengeSolveView
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
