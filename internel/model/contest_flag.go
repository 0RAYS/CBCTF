package model

import (
	"CBCTF/internel/i18n"
	"time"
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
	Basic
}

func (c ContestFlag) GetModelName() string {
	return "ContestFlag"
}

func (c ContestFlag) GetID() uint {
	return c.ID
}

func (c ContestFlag) GetVersion() uint {
	return c.Version
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
