package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
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

func (c *ChallengeFlagRepo) Delete(idL ...uint) (bool, string) {
	contestFlagIDL := make([]uint, 0)
	for _, id := range idL {
		challengeFlag, ok, msg := c.GetByID(id, "ContestFlags")
		if !ok && msg != i18n.ChallengeFlagNotFound {
			return false, msg
		}
		for _, contestFlag := range challengeFlag.ContestFlags {
			contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
		}
	}
	if ok, msg := InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.ChallengeFlag{}).Where("id IN ?", idL).Delete(&model.ChallengeFlag{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ChallengeFlag: %s", res.Error)
		return false, model.ChallengeFlag{}.DeleteErrorString()
	}
	return true, i18n.Success
}
