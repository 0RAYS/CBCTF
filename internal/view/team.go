package view

import "CBCTF/internal/model"

type TeamView struct {
	Unavailable []string
	Team        model.Team
	UserCount   int64
	// Adjustment 是管理员对该队伍做出的分数变动总和，队伍显示分数已包含它。
	Adjustment float64
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
