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
		return nil, model.RetVal{Msg: i18n.Model.Submission.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	for _, submission := range submissions {
		if submission.TeamID != 0 {
			teamIDL = append(teamIDL, submission.TeamID)
		}
	}
	return teamIDL, model.SuccessRetVal()
}

type FlagSolverRow struct {
	UserID   uint
	UserName string
	TeamID   uint
	TeamName string
	Score    float64
	SolvedAt time.Time
}

func (s *SubmissionRepo) ListFlagSolvers(contestFlagID uint) ([]FlagSolverRow, model.RetVal) {
	var rows []FlagSolverRow
	res := s.DB.Table("submissions").
		Select("submissions.user_id, users.name AS user_name, submissions.team_id, teams.name AS team_name, submissions.score, submissions.created_at AS solved_at").
		Joins("LEFT JOIN users ON submissions.user_id = users.id AND users.deleted_at IS NULL").
		Joins("LEFT JOIN teams ON submissions.team_id = teams.id AND teams.deleted_at IS NULL").
		Where("submissions.contest_flag_id = ? AND submissions.solved = true AND submissions.deleted_at IS NULL", contestFlagID).
		Order("submissions.id ASC").
		Scan(&rows)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get flag solvers: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Submission.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return rows, model.SuccessRetVal()
}
