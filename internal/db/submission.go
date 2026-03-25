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

type BloodRankRow struct {
	ContestFlagID uint
	TeamID        uint
	BloodRank     int
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
	rankMap, ret := s.GetBloodRankMap(contestFlagID)
	if !ret.OK {
		return nil, ret
	}
	teamRanks := rankMap[contestFlagID]
	if len(teamRanks) == 0 {
		return nil, model.SuccessRetVal()
	}
	orderedTeamIDs := make([]uint, 3)
	for teamID, rank := range teamRanks {
		if rank <= 0 || rank > len(orderedTeamIDs) {
			continue
		}
		orderedTeamIDs[rank-1] = teamID
	}
	teamIDL := make([]uint, 0, len(orderedTeamIDs))
	for _, teamID := range orderedTeamIDs {
		if teamID != 0 {
			teamIDL = append(teamIDL, teamID)
		}
	}
	return teamIDL, model.SuccessRetVal()
}

func (s *SubmissionRepo) GetBloodRankMap(contestFlagIDL ...uint) (map[uint]map[uint]int, model.RetVal) {
	rankMap := make(map[uint]map[uint]int)
	if len(contestFlagIDL) == 0 {
		return rankMap, model.SuccessRetVal()
	}

	firstSolves := s.DB.Table("submissions").
		Select("contest_flag_id, team_id, MIN(created_at) AS first_solved_at, MIN(id) AS first_submission_id").
		Where("contest_flag_id = ANY(?) AND solved = true AND team_id <> 0 AND deleted_at IS NULL", contestFlagIDL).
		Group("contest_flag_id, team_id")

	ranked := s.DB.Table("(?) AS first_solves", firstSolves).
		Select("contest_flag_id, team_id, ROW_NUMBER() OVER (PARTITION BY contest_flag_id ORDER BY first_solved_at ASC, first_submission_id ASC) AS blood_rank")

	var rows []BloodRankRow
	res := s.DB.Table("(?) AS ranked", ranked).
		Select("contest_flag_id, team_id, blood_rank").
		Where("blood_rank <= 3").
		Scan(&rows)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get blood rank: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Submission.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}

	for _, row := range rows {
		if rankMap[row.ContestFlagID] == nil {
			rankMap[row.ContestFlagID] = make(map[uint]int)
		}
		rankMap[row.ContestFlagID][row.TeamID] = row.BloodRank
	}
	return rankMap, model.SuccessRetVal()
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
		Joins("INNER JOIN users ON submissions.user_id = users.id AND users.deleted_at IS NULL").
		Joins("INNER JOIN teams ON submissions.team_id = teams.id AND teams.deleted_at IS NULL").
		Where("submissions.contest_flag_id = ? AND submissions.solved = true AND submissions.deleted_at IS NULL", contestFlagID).
		Order("submissions.created_at ASC, submissions.id ASC").
		Scan(&rows)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get flag solvers: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Submission.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return rows, model.SuccessRetVal()
}

func (s *SubmissionRepo) ListSolvedByTeamID(teamIDL ...uint) ([]model.Submission, model.RetVal) {
	if len(teamIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	submissions := make([]model.Submission, 0)
	res := s.DB.Model(&model.Submission{}).
		Where("team_id = ANY(?) AND solved = true", teamIDL).
		Order("team_id ASC, created_at ASC, id ASC").
		Find(&submissions)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list solved submissions: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Submission.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return submissions, model.SuccessRetVal()
}
