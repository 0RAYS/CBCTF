package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type CheatRepo struct {
	Basic[model.Cheat]
}

type CreateCheatRepo struct {
	UserID             *uint
	TeamID             *uint
	ContestID          *uint
	ContestChallengeID *uint
	ContestFlagID      *uint
	Reason             string
	Type               string
	Checked            bool
	Hash               string
	References         model.UintList
}

func (c CreateCheatRepo) Convert2Model() model.Model {
	return model.Cheat{
		UserID:             c.UserID,
		TeamID:             c.TeamID,
		ContestID:          c.ContestID,
		ContestChallengeID: c.ContestChallengeID,
		ContestFlagID:      c.ContestFlagID,
		Reason:             c.Reason,
		Type:               c.Type,
		Checked:            c.Checked,
		Hash:               c.Hash,
		References:         c.References,
	}
}

type UpdateCheatRepo struct {
	Reason     *string
	Type       *string
	Checked    *bool
	Hash       *string
	References *model.UintList
}

func (u UpdateCheatRepo) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Reason != nil {
		options["reason"] = *u.Reason
	}
	if u.Type != nil {
		options["type"] = *u.Type
	}
	if u.Checked != nil {
		options["checked"] = *u.Checked
	}
	if u.Hash != nil {
		options["hash"] = *u.Hash
	}
	if u.References != nil {
		options["references"] = u.References
	}
	return options
}

func InitCheatRepo(tx *gorm.DB) *CheatRepo {
	return &CheatRepo{
		Basic: Basic[model.Cheat]{
			DB: tx,
		},
	}
}
