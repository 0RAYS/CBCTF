package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type TeamFlagRepo struct {
	Basic[model.TeamFlag]
}

type CreateTeamFlagRepo struct {
	TeamID        uint
	ContestFlagID uint
	Value         string
	Solved        bool
}

func (c CreateTeamFlagRepo) Convert2Model() model.Model {
	return model.TeamFlag{
		TeamID:        c.TeamID,
		ContestFlagID: c.ContestFlagID,
		Value:         c.Value,
		Solved:        c.Solved,
	}
}

type UpdateTeamFlagRepo struct {
	Value  *string
	Solved *bool
}

func (u UpdateTeamFlagRepo) Convert2Map() map[string]any {
	m := make(map[string]any)
	if u.Value != nil {
		m["value"] = *u.Value
	}
	if u.Solved != nil {
		m["solved"] = *u.Solved
	}
	return m
}

func InitTeamFlagRepo(tx *gorm.DB) *TeamFlagRepo {
	return &TeamFlagRepo{
		Basic: Basic[model.TeamFlag]{
			DB: tx,
		},
	}
}
