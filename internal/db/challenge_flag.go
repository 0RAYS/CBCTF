package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type ChallengeFlagRepo struct {
	BaseRepo[model.ChallengeFlag]
}

type CreateChallengeFlagOptions struct {
	ChallengeID uint
	Value       string
	Binding     model.FlagBinding
}

func (c CreateChallengeFlagOptions) Convert2Model() model.Model {
	return model.ChallengeFlag{
		ChallengeID: c.ChallengeID,
		Value:       c.Value,
		Binding:     c.Binding,
	}
}

type UpdateChallengeFlagOptions struct {
	Value   *string
	Binding *model.FlagBinding
}

func (u UpdateChallengeFlagOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Value != nil {
		options["value"] = *u.Value
	}
	if u.Binding != nil {
		options["binding"] = *u.Binding
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
	var contestFlagIDL []uint
	if res := c.DB.Model(&model.ContestFlag{}).Where("challenge_flag_id IN ?", idL).Pluck("id", &contestFlagIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get ContestFlags for challenge flags %v: %s", idL, res.Error)
		return model.RetVal{Msg: i18n.Model.ContestFlag.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if ret := InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Where("challenge_flag_id IN ?", idL).Delete(&model.TeamFlag{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete TeamFlags for challenge flags %v: %s", idL, res.Error)
		return model.RetVal{Msg: i18n.Model.TeamFlag.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if res := c.DB.Model(&model.ChallengeFlag{}).Where("id IN ?", idL).Delete(&model.ChallengeFlag{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ChallengeFlag: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.ChallengeFlag.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
