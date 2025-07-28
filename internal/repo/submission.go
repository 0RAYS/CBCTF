package repo

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

type SubmissionRepo struct {
	BasicRepo[model.Submission]
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
	IP                 string
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
		BasicRepo: BasicRepo[model.Submission]{
			DB: tx,
		},
	}
}

func (s *SubmissionRepo) GetBloodTeam(contestFlagID uint) ([]uint, bool, string) {
	var submissions []model.Submission
	teamIDL := make([]uint, 0)
	res := s.DB.Model(&model.Submission{}).Where("contest_flag_id = ?", contestFlagID).Order("id").Limit(3).Find(&submissions)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Submission: %s", res.Error)
		return teamIDL, false, i18n.DeleteSubmissionError
	}
	for _, submission := range submissions {
		if submission.TeamID != 0 {
			teamIDL = append(teamIDL, submission.TeamID)
		}
	}
	return teamIDL, true, i18n.Success
}
