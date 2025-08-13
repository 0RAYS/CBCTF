package model

import (
	"CBCTF/internal/i18n"
	"math"
	"time"
)

const (
	StaticScore      uint = 0
	LinearScore      uint = 1
	LogarithmicScore uint = 2

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
	Contest            Contest          `json:"-"`
	ContestChallengeID uint             `json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	ChallengeFlagID    uint             `json:"challenge_flag_id"`
	ChallengeFlag      ChallengeFlag    `json:"-"`
	Submissions        []Submission     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags          []TeamFlag       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Value              string           `json:"value"`
	Score              float64          `gorm:"default:1000" json:"score"`
	CurrentScore       float64          `gorm:"default:1000" json:"current_score"`
	Decay              float64          `gorm:"default:50" json:"decay"`
	MinScore           float64          `gorm:"default:100" json:"min_score"`
	ScoreType          uint             `gorm:"default:0" json:"score_type"`
	Solvers            int64            `json:"solvers"`
	Last               time.Time        `gorm:"default:null" json:"last"`
	BasicModel
}

func (c ContestFlag) GetModelName() string {
	return "ContestFlag"
}

func (c ContestFlag) GetVersion() uint {
	return c.Version
}

func (c ContestFlag) GetBasicModel() BasicModel {
	return c.BasicModel
}

func (c ContestFlag) CreateErrorString() string {
	return i18n.CreateContestFlagError
}

func (c ContestFlag) DeleteErrorString() string {
	return i18n.DeleteContestFlagError
}

func (c ContestFlag) GetErrorString() string {
	return i18n.GetContestFlagError
}

func (c ContestFlag) NotFoundErrorString() string {
	return i18n.ContestFlagNotFound
}

func (c ContestFlag) UpdateErrorString() string {
	return i18n.UpdateContestFlagError
}

func (c ContestFlag) GetUniqueKey() []string {
	return []string{"id"}
}

func (c ContestFlag) CalcScore(solvers int64) float64 {
	var calc float64 = 0
	switch c.ScoreType {
	case StaticScore:
		calc = c.Score
	case LinearScore:
		calc = c.Score - float64(solvers)*c.Decay
	case LogarithmicScore:
		calc = (c.MinScore-c.Score)/(c.Decay*c.Decay)*float64(solvers*solvers) + c.Score
	default:
		calc = c.Score
	}
	if calc < c.MinScore {
		calc = c.MinScore
	}
	return math.Trunc(calc*100) / 100
}
