package view

import (
	"time"

	"CBCTF/internal/model"
)

type VictimStatusView struct {
	Targets   []string
	Duration  float64
	Remaining float64
	Status    string
}

type ChallengeRuntimeView struct {
	Challenge model.Challenge
	Attempts  int64
	Init      bool
	Solved    bool
	Remote    VictimStatusView
	FileName  string
}

type ChallengeStatusView struct {
	Attempts int64
	Init     bool
	Solved   bool
	Remote   VictimStatusView
	FileName string
}

type FlagSolverView struct {
	UserID   uint
	UserName string
	TeamID   uint
	TeamName string
	Score    float64
	SolvedAt time.Time
}
