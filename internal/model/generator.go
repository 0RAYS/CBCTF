package model

import (
	"database/sql"
	"time"
)

const (
	WaitingGeneratorStatus     = WaitingStatus
	PendingGeneratorStatus     = PendingStatus
	RunningGeneratorStatus     = RunningStatus
	TerminatingGeneratorStatus = TerminatingStatus
	StoppedGeneratorStatus     = StoppedStatus
)

// Generator 是动态题目生成器实例。
// ChallengeID 题库题目溯源。
// ContestChallengeID 比赛题目 id，同一题库题目在多场比赛加入时各自持有独立副本。
// Generator 是动态题目生成器实例。
// ChallengeID 题库题目溯源。
// ContestChallengeID 比赛题目 id，同一题库题目在多场比赛加入时各自持有独立副本。
type Generator struct {
	WorkerToken        string           `json:"-"`
	Image              string           `json:"image"`
	ChallengeID        uint             `gorm:"index" json:"challenge_id"`
	ChallengeName      string           `json:"challenge_name"`
	Challenge          Challenge        `json:"-"`
	ContestChallengeID sql.Null[uint]   `gorm:"index" json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	ContestID          sql.Null[uint]   `gorm:"index" json:"contest_id"`
	Contest            Contest          `json:"-"`
	Name               string           `json:"pod_name"`
	Success            int64            `gorm:"default:0" json:"success"`
	SuccessLast        time.Time        `gorm:"default:null" json:"success_last"`
	Failure            int64            `gorm:"default:0" json:"failure"`
	FailureLast        time.Time        `gorm:"default:null" json:"failure_last"`
	Status             string           `json:"status"`
	BaseModel
}

// Target 返回生成器池的目标标识
func (g Generator) Target() GeneratorTarget {
	return GeneratorTarget{
		ContestID:          g.ContestID.V,
		ChallengeID:        g.ChallengeID,
		ContestChallengeID: g.ContestChallengeID.V,
	}
}
