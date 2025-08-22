package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"

	"gorm.io/gorm"
)

type ChallengeFlagRepo struct {
	BasicRepo[model.ChallengeFlag]
}

type CreateChallengeFlagOptions struct {
	ChallengeID uint
	DockerID    sql.Null[uint]
	Name        string
	Value       string
	InjectType  string
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
		BasicRepo: BasicRepo[model.ChallengeFlag]{
			DB: tx,
		},
	}
}

func (c *ChallengeFlagRepo) Delete(idL ...uint) (bool, string) {
	challengeFlagL, _, ok, msg := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"ContestFlags": {Selects: []string{"id", "challenge_flag_id"}},
			"TeamFlags":    {Selects: []string{"id", "challenge_flag_id"}},
		},
	})
	if !ok && msg != i18n.ChallengeFlagNotFound {
		return false, msg
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
	if ok, msg = InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitTeamFlagRepo(c.DB).Delete(teamFlagIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.ChallengeFlag{}).Where("id IN ?", idL).Delete(&model.ChallengeFlag{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ChallengeFlag: %s", res.Error)
		return false, i18n.DeleteChallengeFlagError
	}
	return true, i18n.Success
}
