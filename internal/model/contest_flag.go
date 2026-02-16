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
	ContestID          uint             `json:"contest_id"`
	ContestChallengeID uint             `json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	ChallengeFlagID    uint             `json:"challenge_flag_id"`
	Submissions        []Submission     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags          []TeamFlag       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Value              string           `json:"value"`
	Score              float64          `gorm:"default:1000" json:"score"`
	CurrentScore       float64          `gorm:"default:1000" json:"current_score"`
	Decay              float64          `gorm:"default:50" json:"decay"`
	MinScore           float64          `gorm:"default:100" json:"min_score"`
	ScoreType          ScoreType        `gorm:"default:0" json:"score_type"`
	Solvers            int64            `json:"solvers"`
	Last               time.Time        `gorm:"default:null" json:"last"`
	BaseModel
}

func (c ContestFlag) TableName() string {
	return "contest_flags"
}

func (c ContestFlag) ModelName() string {
	return "ContestFlag"
}

func (c ContestFlag) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c ContestFlag) UniqueFields() []string {
	return []string{"id"}
}

func (c ContestFlag) QueryFields() []string {
	return []string{}
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
