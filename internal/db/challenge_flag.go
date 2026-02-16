package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"

	"gorm.io/gorm"
)

type ChallengeFlagRepo struct {
	BaseRepo[model.ChallengeFlag]
}

type CreateChallengeFlagOptions struct {
	ChallengeID uint
	DockerID    sql.Null[uint]
	Name        string
	Value       string
	InjectType  model.FlagInjectType
	Path        string
}

func (c CreateChallengeFlagOptions) Convert2Model() model.Model {
	return model.ChallengeFlag{
		ChallengeID: c.ChallengeID,
		DockerID:    c.DockerID,
		Name:        c.Name,
		Value:       c.Value,
		InjectType:  c.InjectType,
		Path:        c.Path,
	}
}

type UpdateChallengeFlagOptions struct {
	Value      *string
	InjectType *model.FlagInjectType
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
		BaseRepo: BaseRepo[model.ChallengeFlag]{
			DB: tx,
		},
	}
}

func (c *ChallengeFlagRepo) Delete(idL ...uint) model.RetVal {
	challengeFlagL, _, ret := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Preloads: map[string]GetOptions{
			"ContestFlags": {},
			"TeamFlags":    {},
		},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	contestFlagIDL, teamFlagIDL := make([]uint, 0), make([]uint, 0)
	for _, challengeFlag := range challengeFlagL {
		for _, contestFlag := range challengeFlag.ContestFlags {
			contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
		}
		for _, teamFlag := range challengeFlag.TeamFlags {
			teamFlagIDL = append(teamFlagIDL, teamFlag.ID)
		}
	}
	if ret = InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ret.OK {
		return ret
	}
	if ret = InitTeamFlagRepo(c.DB).Delete(teamFlagIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.ChallengeFlag{}).Where("id IN ?", idL).Delete(&model.ChallengeFlag{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ChallengeFlag: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]interface{}{"Model": model.ChallengeFlag{}.ModelName(), "Error": res.Error}}
	}
	return model.SuccessRetVal()
}
