package model

import (
	"math"
	"time"
)

type ScoreType uint

const (
	StaticScoreType      ScoreType = 0
	LinearScoreType      ScoreType = 1
	LogarithmicScoreType ScoreType = 2

	FirstBloodRate  float64 = 0.05
	SecondBloodRate float64 = 0.03
	ThirdBloodRate  float64 = 0.01
)

// ContestFlag
// BelongsTo Contest
// BelongsTo ContestChallenge
// BelongsTo ChallengeFlag
// HasMany Submission
// HasMany TeamFlag
type ContestFlag struct {
	ContestID          uint             `gorm:"index" json:"contest_id"`
	Contest            Contest          `json:"-"`
	ContestChallengeID uint             `gorm:"index;index:idx_contest_flags_challenge_active,priority:1" json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	ChallengeFlagID    uint             `gorm:"index" json:"challenge_flag_id"`
	ChallengeFlag      ChallengeFlag    `json:"-"`
	Submissions        []Submission     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags          []TeamFlag       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Value              string           `json:"value"`
	// Binding 为题库 flag 注入绑定的副本
	Binding FlagBinding `gorm:"serializer:json;type:jsonb;default:'{}'" json:"binding"`
	// Generation 为代次。flag 模板变更时旧代次归档（Active=false），新代次继续递增。
	Generation int `gorm:"default:1;index:idx_contest_flags_challenge_active,priority:3" json:"generation"`
	// Active 表示该 flag 是否为当前生效代次，只有生效代次参与发题与计分
	Active       bool      `gorm:"default:true;index:idx_contest_flags_challenge_active,priority:2" json:"active"`
	Score        float64   `gorm:"default:1000" json:"score"`
	CurrentScore float64   `gorm:"default:1000" json:"current_score"`
	Decay        float64   `gorm:"default:50" json:"decay"`
	MinScore     float64   `gorm:"default:100" json:"min_score"`
	ScoreType    ScoreType `gorm:"default:0" json:"score_type"`
	Solvers      int64     `json:"solvers"`
	Last         time.Time `gorm:"default:null" json:"last"`
	BaseModel
}

func (c ContestFlag) CalcScore(solvers int64) float64 {
	if solvers < 0 {
		solvers = 0
	}
	var calc float64 = 0
	switch c.ScoreType {
	case StaticScoreType:
		calc = c.Score
	case LinearScoreType:
		calc = c.Score - float64(solvers)*c.Decay
	case LogarithmicScoreType:
		if c.Decay > 0 {
			k := 5.0 / c.Decay
			calc = (c.Score-c.MinScore)*math.Exp(-k*float64(solvers)) + c.MinScore
		} else {
			calc = c.MinScore
		}
	default:
		calc = c.Score
	}
	if calc < c.MinScore {
		calc = c.MinScore
	}
	return math.Trunc(calc*100) / 100
}
