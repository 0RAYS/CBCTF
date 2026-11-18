package view

import "CBCTF/internal/model"

type ChallengeFlagView struct {
	ID    uint
	Value string
}

type ChallengeView struct {
	Unavailable   []string
	Challenge     model.Challenge
	Flags         []ChallengeFlagView
	DockerCompose string
	FileName      string
}

type SimpleChallengeView struct {
	Challenge model.Challenge
}

// ChallengeUsageView 描述一个引用了题库题目的比赛题目实例。
type ChallengeUsageView struct {
	// ContestID 比赛 id
	ContestID uint
	// ContestName 比赛名称
	ContestName string
	// ContestChallenge 比赛题目实例
	ContestChallenge model.ContestChallenge
	// Solvers 已解出该题目的队伍数
	Solvers int64
}
