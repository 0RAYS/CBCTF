package repo

import (
	"CBCTF/internel/model"
	"crypto/md5"
	"fmt"
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
	References         model.UintList
}

func (c CreateCheatRepo) Convert2Model() model.Model {
	tmp := make([]uint, 5)
	if c.UserID != nil {
		tmp[0] = *c.UserID
	}
	if c.TeamID != nil {
		tmp[1] = *c.TeamID
	}
	if c.ContestID != nil {
		tmp[2] = *c.ContestID
	}
	if c.ContestChallengeID != nil {
		tmp[3] = *c.ContestChallengeID
	}
	if c.ContestFlagID != nil {
		tmp[4] = *c.ContestFlagID
	}
	hash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d%d%d%d%d%s%s%t", tmp[0], tmp[1], tmp[2], tmp[3], tmp[4], c.Reason, c.Type, c.Checked))))
	return model.Cheat{
		UserID:             c.UserID,
		TeamID:             c.TeamID,
		ContestID:          c.ContestID,
		ContestChallengeID: c.ContestChallengeID,
		ContestFlagID:      c.ContestFlagID,
		Reason:             c.Reason,
		Type:               c.Type,
		Checked:            c.Checked,
		References:         c.References,
		Hash:               hash,
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
