package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type TeamFlagRepo struct {
	BaseRepo[model.TeamFlag]
}

type CreateTeamFlagOptions struct {
	TeamID          uint
	ContestFlagID   uint
	ChallengeFlagID uint
	Value           string
	Solved          bool
}

func (c CreateTeamFlagOptions) Convert2Model() model.Model {
	return model.TeamFlag{
		TeamID:          c.TeamID,
		ContestFlagID:   c.ContestFlagID,
		ChallengeFlagID: c.ChallengeFlagID,
		Value:           c.Value,
		Solved:          c.Solved,
	}
}

type UpdateTeamFlagRepo struct {
	Value  *string
	Solved *bool
}

func (u UpdateTeamFlagRepo) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Value != nil {
		options["value"] = *u.Value
	}
	if u.Solved != nil {
		options["solved"] = *u.Solved
	}
	return options
}

func InitTeamFlagRepo(tx *gorm.DB) *TeamFlagRepo {
	return &TeamFlagRepo{
		BaseRepo: BaseRepo[model.TeamFlag]{
			DB: tx,
		},
	}
}

type TeamFlagWithChallenge struct {
	model.TeamFlag
	// ContestFlag 字段
	ContestFlagValue         string
	ContestFlagScore         float64
	ContestFlagCurrentScore  float64
	ContestFlagMinScore      float64
	ContestFlagScoreType     uint
	ContestChallengeName     string
	ContestChallengeCategory string
	ContestChallengeHidden   bool
	ChallengeID              uint
	ChallengeRandID          string
	ChallengeName            string
}

func (t *TeamFlagRepo) GetTeamFlagsWithChallenge(teamIDL ...uint) ([]TeamFlagWithChallenge, model.RetVal) {
	if len(teamIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	var results []TeamFlagWithChallenge
	res := t.DB.Table("team_flags").
		Select(`team_flags.*,
			contest_flags.value AS contest_flag_value,
			contest_flags.score AS contest_flag_score,
			contest_flags.current_score AS contest_flag_current_score,
			contest_flags.min_score AS contest_flag_min_score,
			contest_flags.score_type AS contest_flag_score_type,
			contest_challenges.name AS contest_challenge_name,
			contest_challenges.category AS contest_challenge_category,
			contest_challenges.hidden AS contest_challenge_hidden,
			challenges.id AS challenge_id,
			challenges.rand_id AS challenge_rand_id,
			challenges.name AS challenge_name`).
		Joins("INNER JOIN contest_flags ON team_flags.contest_flag_id = contest_flags.id AND contest_flags.deleted_at IS NULL").
		Joins("INNER JOIN contest_challenges ON contest_flags.contest_challenge_id = contest_challenges.id AND contest_challenges.deleted_at IS NULL").
		Joins("INNER JOIN challenges ON contest_challenges.challenge_id = challenges.id AND challenges.deleted_at IS NULL").
		Where("team_flags.team_id IN ? AND team_flags.deleted_at IS NULL", teamIDL).
		Scan(&results)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get TeamFlags: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.TeamFlag.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return results, model.SuccessRetVal()
}

func (t *TeamFlagRepo) Exists(teamID uint, contestFlagIDL ...uint) (bool, model.RetVal) {
	if len(contestFlagIDL) == 0 {
		return false, model.SuccessRetVal()
	}
	var exists bool
	res := t.DB.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM team_flags
			WHERE team_id = ? AND contest_flag_id IN ? AND deleted_at IS NULL
		)
	`, teamID, contestFlagIDL).Scan(&exists)
	if res.Error != nil {
		log.Logger.Warningf("Failed to check team flag existence: %s", res.Error)
		return false, model.RetVal{Msg: i18n.Model.TeamFlag.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return exists, model.SuccessRetVal()
}

func (t *TeamFlagRepo) CountGenerated(teamID uint, contestFlagIDL ...uint) (int64, model.RetVal) {
	if len(contestFlagIDL) == 0 {
		return 0, model.SuccessRetVal()
	}
	var count int64
	res := t.DB.Model(&model.TeamFlag{}).
		Where("team_id = ? AND contest_flag_id IN ?", teamID, contestFlagIDL).
		Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count generated team flags: %s", res.Error)
		return 0, model.RetVal{Msg: i18n.Model.TeamFlag.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}
