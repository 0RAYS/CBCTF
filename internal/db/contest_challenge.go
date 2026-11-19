package db

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type ContestChallengeRepo struct {
	BaseRepo[model.ContestChallenge]
}

type UpdateContestChallengeOptions struct {
	Name            *string
	Description     *string
	Hidden          *bool
	Attempt         *int64
	Hints           *[]string
	Tags            *[]string
	Category        *string
	Type            *model.ChallengeType
	Template        *model.ChallengeTemplate
	GeneratorImage  *string
	NetworkPolicies *model.NetworkPolicies
}

func (u UpdateContestChallengeOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Description != nil {
		options["description"] = *u.Description
	}
	if u.Hidden != nil {
		options["hidden"] = *u.Hidden
	}
	if u.Attempt != nil {
		options["attempt"] = *u.Attempt
	}
	if u.Hints != nil {
		options["hints"] = *u.Hints
	}
	if u.Tags != nil {
		options["tags"] = *u.Tags
	}
	if u.Category != nil {
		options["category"] = *u.Category
	}
	if u.Type != nil {
		options["type"] = *u.Type
	}
	if u.Template != nil {
		options["template"] = *u.Template
	}
	if u.GeneratorImage != nil {
		options["generator_image"] = *u.GeneratorImage
	}
	if u.NetworkPolicies != nil {
		options["network_policies"] = *u.NetworkPolicies
	}
	return options
}

func InitContestChallengeRepo(tx *gorm.DB) *ContestChallengeRepo {
	return &ContestChallengeRepo{
		DB: tx,
	}
}

func (c *ContestChallengeRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.ContestChallenge, model.RetVal) {
	return c.GetByUniqueField("rand_id", randID, optionsL...)
}

func (c *ContestChallengeRepo) DeleteByContestID(contestIDL ...uint) model.RetVal {
	return c.DeleteByFieldID("contest_id", contestIDL...)
}

func (c *ContestChallengeRepo) DeleteByChallengeID(challengeIDL ...uint) model.RetVal {
	return c.DeleteByFieldID("challenge_id", challengeIDL...)
}

func (c *ContestChallengeRepo) ListCategories(contestID uint, t model.ChallengeType) ([]string, model.RetVal) {
	var categories []string
	tx := c.DB.Model(&model.ContestChallenge{}).Where("contest_id = ?", contestID)
	conditions := map[string]any{"contest_id": contestID}
	if t != "" {
		tx = tx.Where("type = ?", t)
		conditions["type"] = t
	}
	target := conditionTarget(conditions)
	if res := tx.Distinct().Order("category ASC").Pluck("category", &categories); res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge categories: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.ContestChallenge.GetError, "ContestChallenge", target, res.Error)
	}
	return categories, model.SuccessRetVal()
}

func (c *ContestChallengeRepo) ListUnsolvedID(teamID, contestID uint, category string, limit, offset int) ([]uint, int64, model.RetVal) {
	base := c.DB.Model(&model.ContestChallenge{}).
		Where("contest_challenges.contest_id = ?", contestID).
		Where("contest_challenges.hidden = ?", false).
		Where("EXISTS (SELECT 1 FROM contest_flags WHERE contest_flags.contest_challenge_id = contest_challenges.id AND contest_flags.deleted_at IS NULL AND contest_flags.active = true)").
		Where(`
			EXISTS (
				SELECT 1
				FROM contest_flags
				WHERE contest_flags.contest_challenge_id = contest_challenges.id
					AND contest_flags.deleted_at IS NULL
					AND contest_flags.active = true
					AND NOT EXISTS (
						SELECT 1
						FROM team_flags
						WHERE team_flags.team_id = ?
							AND team_flags.contest_flag_id = contest_flags.id
							AND team_flags.solved = true
							AND team_flags.deleted_at IS NULL
					)
			)
		`, teamID)
	if category != "" {
		base = base.Where("contest_challenges.category = ?", category)
	}
	conditions := map[string]any{"contest_id": contestID, "team_id": teamID}
	if category != "" {
		conditions["category"] = category
	}
	target := conditionTarget(conditions)

	var count int64
	if res := base.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.ContestChallenge.GetError, "ContestChallenge", target, res.Error)
	}

	var ids []uint
	res := base.Select("contest_challenges.id").
		Order("contest_challenges.id ASC").
		Limit(limit).Offset(offset).
		Scan(&ids)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.ContestChallenge.GetError, "ContestChallenge", target, res.Error)
	}
	return ids, count, model.SuccessRetVal()
}

func (c *ContestChallengeRepo) Delete(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	if ret := InitContestFlagRepo(c.DB).DeleteByContestChallengeID(idL...); !ret.OK {
		return ret
	}
	if ret := InitSubmissionRepo(c.DB).DeleteByContestChallengeID(idL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.ContestChallenge{}).Where("id IN ?", idL).Delete(&model.ContestChallenge{}); res.Error != nil {
		target := fmt.Sprintf("id=%v", idL)
		log.Logger.Warningf("Failed to delete ContestChallenge: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.ContestChallenge.DeleteError, "ContestChallenge", target, res.Error)
	}
	return model.SuccessRetVal()
}
