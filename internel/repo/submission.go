package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type SubmissionRepo struct {
	Basic[model.Submission]
}

type CreateSubmissionOptions struct {
	ContestChallengeID uint
	ContestID          uint
	ChallengeID        uint
	TeamID             uint
	UserID             uint
	ContestFlagID      uint
	Value              string
	Solved             bool
	Score              float64
}

func (c CreateSubmissionOptions) Convert2Model() model.Model {
	return model.Submission{
		ContestChallengeID: c.ContestChallengeID,
		ContestID:          c.ContestID,
		ChallengeID:        c.ChallengeID,
		TeamID:             c.TeamID,
		UserID:             c.UserID,
		ContestFlagID:      c.ContestFlagID,
		Value:              c.Value,
		Solved:             c.Solved,
		Score:              c.Score,
	}
}

type UpdateSubmissionOptions struct {
	Solved *bool
	Score  *float64
}

func (u UpdateSubmissionOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Solved != nil {
		options["solved"] = *u.Solved
	}
	if u.Score != nil {
		options["score"] = *u.Score
	}
	return options
}

func InitSubmissionRepo(tx *gorm.DB) *SubmissionRepo {
	return &SubmissionRepo{
		Basic: Basic[model.Submission]{
			DB: tx,
		},
	}
}
