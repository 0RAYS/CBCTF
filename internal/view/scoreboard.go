package view

import (
	"CBCTF/internal/model"
	"time"
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
