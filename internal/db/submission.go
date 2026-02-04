package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type SubmissionRepo struct {
	BaseRepo[model.Submission]
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
		IP:                 c.IP,
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
		BaseRepo: BaseRepo[model.Submission]{
			DB: tx,
		},
	}
}

func (s *SubmissionRepo) GetBloodTeamID(contestFlagID uint) ([]uint, model.RetVal) {
	var submissions []model.Submission
	teamIDL := make([]uint, 0)
	res := s.DB.Model(&model.Submission{}).Where("contest_flag_id = ?", contestFlagID).Order("id").Limit(3).Find(&submissions)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Submission: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Submission{}.ModelName(), "Error": res.Error.Error()}}
	}
	for _, submission := range submissions {
		if submission.TeamID != 0 {
			teamIDL = append(teamIDL, submission.TeamID)
		}
	}
	return teamIDL, model.SuccessRetVal()
}

type UserSolvedSubmission struct {
	SubmissionID            uint
	UserID                  uint
	TeamID                  uint
	ContestID               uint
	Score                   float64
	SubmissionTime          time.Time
	ContestFlagID           uint
	ContestFlagScore        float64
	ContestFlagCurrentScore float64
	ContestFlagMinScore     float64
	ContestFlagScoreType    uint
}

func (s *SubmissionRepo) GetUserSolvedSubmissions(userIDL ...uint) ([]UserSolvedSubmission, model.RetVal) {
	if len(userIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	var results []UserSolvedSubmission
	res := s.DB.Raw(`
		SELECT submissions.id, submissions.user_id, submissions.team_id, submissions.contest_id, submissions.contest_flag_id,
		submissions.score, submissions.created_at AS submission_time,
		contest_flags.score as contest_flag_score, contest_flags.current_score as contest_flag_current_flag,
		contest_flags.min_score as contest_flag_min_score, contest_flags.score_type as contest_flag_score_type
    	FROM submissions
		INNER JOIN contest_flags ON submissions.contest_flag_id = contest_flags.id AND contest_flags.deleted_at IS NULL
		WHERE submissions.user_id IN ? AND submissions.solved = true AND submissions.deleted_at IS NULL
	`, userIDL).Find(&results)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Submissions: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Submission{}.ModelName(), "Error": res.Error.Error()}}
	}
	return results, model.SuccessRetVal()
}
