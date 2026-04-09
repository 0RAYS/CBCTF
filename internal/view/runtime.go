package view

import (
	"CBCTF/internal/model"
	"time"
)

type VictimStatusView struct {
	Targets   []string
	Duration  float64
	Remaining float64
	Status    string
}

type ContestChallengeView struct {
	ContestChallenge model.ContestChallenge
	Attempts         int64
	Init             bool
	Solved           bool
	Remote           VictimStatusView
	FileName         string
}

type ContestChallengeStatusView struct {
	Attempts int64
	Init     bool
	Solved   bool
	Remote   VictimStatusView
	FileName string
}

type ContestFlagSolverView struct {
	UserID   uint
	UserName string
	TeamID   uint
	TeamName string
	Score    float64
	SolvedAt time.Time
}
