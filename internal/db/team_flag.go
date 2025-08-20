package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamFlagRepo struct {
	BasicRepo[model.TeamFlag]
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

type DiffUpdateTeamFlag struct{}

func (d DiffUpdateTeamFlag) Convert2Expr() map[string]clause.Expr {
	return nil
}

func InitTeamFlagRepo(tx *gorm.DB) *TeamFlagRepo {
	return &TeamFlagRepo{
		BasicRepo: BasicRepo[model.TeamFlag]{
			DB: tx,
		},
	}
}
