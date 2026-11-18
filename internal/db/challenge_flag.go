package db

import (
	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type ChallengeFlagRepo struct {
	BaseRepo[model.ChallengeFlag]
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
		DB: tx,
	}
}

func (c *ChallengeFlagRepo) DeleteByChallengeID(challengeIDL ...uint) model.RetVal {
	return c.DeleteByFieldID("challenge_id", challengeIDL...)
}

// Delete 只删除题库侧的 flag。
// 比赛题目自加入时即持有 flag 副本，题库 flag 仅用于溯源，
// 因此删除题库 flag 不再级联删除比赛侧的 contest_flags / team_flags / submissions。
func (c *ChallengeFlagRepo) Delete(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	if res := c.DB.Model(&model.ChallengeFlag{}).Where("id IN ?", idL).Delete(&model.ChallengeFlag{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ChallengeFlag: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.ChallengeFlag.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
