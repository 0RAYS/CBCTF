package view

import "CBCTF/internal/model"

type TeamView struct {
	Team      model.Team
	UserCount int64
}

type TeamFlagInfoView struct {
	Value        string
	Solved       bool
	Template     string
	InitScore    float64
	CurrentScore float64
	Decay        float64
	MinScore     float64
	Solvers      int64
}

type TeamFlagChallengeView struct {
	Name     string
	Type     model.ChallengeType
	Category string
	Hidden   bool
	Flags    []TeamFlagInfoView
}
