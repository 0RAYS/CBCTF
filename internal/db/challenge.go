package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type ChallengeRepo struct {
	BaseRepo[model.Challenge]
}

type UpdateChallengeOptions struct {
	Name            *string
	Description     *string
	Category        *string
	GeneratorImage  *string
	NetworkPolicies *model.NetworkPolicies
	Template        *model.ChallengeTemplate
}

func (u UpdateChallengeOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Description != nil {
		options["description"] = *u.Description
	}
	if u.Category != nil {
		options["category"] = *u.Category
	}
	if u.GeneratorImage != nil {
		options["generator_image"] = *u.GeneratorImage
	}
	if u.NetworkPolicies != nil {
		options["network_policies"] = *u.NetworkPolicies
	}
	if u.Template != nil {
		options["template"] = *u.Template
	}
	return options
}

func InitChallengeRepo(tx *gorm.DB) *ChallengeRepo {
	return &ChallengeRepo{
		DB: tx,
	}
}

func (c *ChallengeRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.Challenge, model.RetVal) {
	return c.GetByUniqueField("rand_id", randID, optionsL...)
}

// ListForUnfinishedContests includes upcoming and running contests, regardless
// of visibility. Use Go duration arithmetic, matching Contest.IsOver exactly.
func (c *ChallengeRepo) ListForUnfinishedContests(now time.Time) ([]model.Challenge, model.RetVal) {
	contests, ret := InitContestRepo(c.DB.Select("id", "start", "duration")).FindAll()
	if !ret.OK {
		return nil, ret
	}
	var contestIDs []uint
	for _, contest := range contests {
		if !now.After(contest.Start.Add(contest.Duration)) {
			contestIDs = append(contestIDs, contest.ID)
		}
	}
	if len(contestIDs) == 0 {
		return []model.Challenge{}, model.SuccessRetVal()
	}

	// IN selects each challenge once even when several eligible contests refer
	// to it. GORM applies soft-delete scopes to both references and challenges.
	references := c.DB.Model(&model.ContestChallenge{}).
		Select("challenge_id").
		Where("contest_id IN ?", contestIDs)
	return InitChallengeRepo(c.DB.Where("id IN (?)", references)).FindAll()
}

func (c *ChallengeRepo) ListCategories(t model.ChallengeType) ([]string, model.RetVal) {
	var categories = make([]string, 0)
	res := c.DB.Model(&model.Challenge{})
	target := ""
	if t != "" {
		res = res.Where("type = ?", t)
		target = fmt.Sprintf("type=%s", t)
	}
	res = res.Distinct().Order("category ASC").Pluck("category", &categories)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Categories: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}
	return categories, model.SuccessRetVal()
}

// ListChallengesAvailable 返回题库题目列表，同一题库题目允许在同场比赛多次添加。
func (c *ChallengeRepo) ListChallengesAvailable(
	limit, offset int,
	name, description, category string,
	t model.ChallengeType,
) ([]model.Challenge, int64, model.RetVal) {
	tx := c.DB.Model(&model.Challenge{})
	conditions := make(map[string]any)
	if name != "" {
		tx = tx.Where("challenges.name ILIKE ?", "%"+name+"%")
		conditions["name"] = name
	}
	if description != "" {
		tx = tx.Where("challenges.description ILIKE ?", "%"+description+"%")
		conditions["description"] = description
	}
	if category != "" {
		tx = tx.Where("challenges.category = ?", category)
		conditions["category"] = category
	}
	if t != "" {
		tx = tx.Where("challenges.type = ?", t)
		conditions["type"] = t
	}
	target := conditionTarget(conditions)

	var count int64
	if res := tx.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count Challenges: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}

	var challenges = make([]model.Challenge, 0)
	if res := tx.Order("challenges.id DESC").Limit(limit).Offset(offset).Find(&challenges); res.Error != nil {
		log.Logger.Warningf("Failed to list Challenges: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}

	return challenges, count, model.SuccessRetVal()
}

// Delete 删除题库题目。
// 比赛题目是独立副本，题库删除不会影响比赛题目及其解题记录。
func (c *ChallengeRepo) Delete(randIDL ...string) model.RetVal {
	if len(randIDL) == 0 {
		return model.SuccessRetVal()
	}
	var challengeIDL []uint
	target := fmt.Sprintf("rand_id=%v", randIDL)
	if res := c.DB.Model(&model.Challenge{}).Where("rand_id IN ?", randIDL).Pluck("id", &challengeIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get Challenges: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}
	var ret model.RetVal
	if ret = InitChallengeFlagRepo(c.DB).DeleteByChallengeID(challengeIDL...); !ret.OK {
		return ret
	}
	if ret = InitGeneratorRepo(c.DB).DeleteByChallengeID(challengeIDL...); !ret.OK {
		return ret
	}
	if ret = InitVictimRepo(c.DB).DeleteByChallengeID(challengeIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.Challenge{}).Where("rand_id IN ?", randIDL).Delete(&model.Challenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.Challenge.DeleteError, "Challenge", target, res.Error)
	}
	return model.SuccessRetVal()
}

// ListContestChallengeUsages 返回引用了指定题库题目的比赛题目实例。
// 用于在删除或修改题库题目前提示其影响范围。
func (c *ChallengeRepo) ListContestChallengeUsages(challengeID uint) ([]model.ContestChallenge, model.RetVal) {
	usages, _, ret := InitContestChallengeRepo(c.DB).List(-1, -1, GetOptions{
		Conditions: map[string]any{"challenge_id": challengeID},
		Preloads:   map[string]GetOptions{"Contest": {}},
	})
	return usages, ret
}
