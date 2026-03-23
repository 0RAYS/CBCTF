package model

import (
	"database/sql"
	"time"
)

const (
	WaitingGeneratorStatus = "waiting"
	PendingGeneratorStatus = "pending"
	RunningGeneratorStatus = "running"
	StoppedGeneratorStatus = "stopped"
)

// Generator
// BelongsTo Challenge
// BelongsTo Contest
type Generator struct {
	ChallengeID   uint           `json:"challenge_id"`
	ChallengeName string         `json:"challenge_name"`
	Challenge     Challenge      `json:"-"`
	ContestID     sql.Null[uint] `json:"contest_id"`
	Contest       Contest        `json:"-"`
	Name          string         `json:"pod_name"`
	Success       int64          `gorm:"default:0" json:"success"`
	SuccessLast   time.Time      `gorm:"default:null" json:"success_last"`
	Failure       int64          `gorm:"default:0" json:"failure"`
	FailureLast   time.Time      `gorm:"default:null" json:"failure_last"`
	Status        string         `json:"status"`
	BaseModel
}
