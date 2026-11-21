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
// ChallengeID 题目 id，生成器池按题目划分。
type Generator struct {
	WorkerToken   string         `json:"-"`
	Image         string         `json:"image"`
	ChallengeID   uint           `gorm:"index" json:"challenge_id"`
	ChallengeName string         `json:"challenge_name"`
	Challenge     Challenge      `json:"-"`
	ContestID     sql.Null[uint] `gorm:"index" json:"contest_id"`
	Contest       Contest        `json:"-"`
	Name          string         `json:"pod_name"`
	Success       int64          `gorm:"default:0" json:"success"`
	SuccessLast   time.Time      `gorm:"default:null" json:"success_last"`
	Failure       int64          `gorm:"default:0" json:"failure"`
	FailureLast   time.Time      `gorm:"default:null" json:"failure_last"`
	Status        string         `json:"status"`
	BaseModel
}

// Target 返回生成器池的目标标识
func (g Generator) Target() GeneratorTarget {
	return GeneratorTarget{
		ChallengeID: g.ChallengeID,
	}
}
