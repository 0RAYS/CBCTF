package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type ChallengeFlagRepo struct {
	Basic[model.ChallengeFlag]
}

type CreateChallengeFlagOptions struct {
	ChallengeID uint
	DockerID    *uint
	Value       string
	InjectType  string
	Path        string
}

func (c CreateChallengeFlagOptions) Convert2Model() model.Model {
	return model.ChallengeFlag{
		ChallengeID: c.ChallengeID,
		DockerID:    c.DockerID,
		Value:       c.Value,
		InjectType:  c.InjectType,
		Path:        c.Path,
	}
}

type UpdateChallengeFlagOptions struct {
	Value      *string
	InjectType *string
	Path       *string
}

func (u UpdateChallengeFlagOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Value != nil {
		options["value"] = *u.Value
	}
	if u.InjectType != nil {
		options["inject_type"] = *u.InjectType
	}
	if u.Path != nil {
		options["path"] = *u.Path
	}
	return options
}

func InitChallengeFlagRepo(tx *gorm.DB) *ChallengeFlagRepo {
	return &ChallengeFlagRepo{
		Basic: Basic[model.ChallengeFlag]{
			DB: tx,
		},
	}
}
