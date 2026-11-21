package db

import (
	"database/sql"
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
	ContestID       *sql.Null[uint]
	Hidden          *bool
	Attempt         *int64
	Hints           *[]string
	Tags            *[]string
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
	if u.ContestID != nil {
		options["contest_id"] = *u.ContestID
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

// ChallengeFilter 描述题库聚合视图的过滤条件。
// ContestID 为 0 且 Unbound 为 false 时不过滤比赛，返回所有题目。
type ChallengeFilter struct {
	ContestID   uint
	Unbound     bool
	ContestName string
	Type        model.ChallengeType
	Category    string
	Name        string
	Description string
}

// apply 把过滤条件拼接到查询上。
func (f ChallengeFilter) apply(tx *gorm.DB) *gorm.DB {
	if f.Unbound {
		tx = tx.Where("challenges.contest_id IS NULL")
	} else if f.ContestID != 0 {
		tx = tx.Where("challenges.contest_id = ?", f.ContestID)
	}
	if f.ContestName != "" {
		tx = tx.Where(
			"challenges.contest_id IN (SELECT id FROM contests WHERE name ILIKE ? AND deleted_at IS NULL)",
			"%"+f.ContestName+"%",
		)
	}
	if f.Type != "" {
		tx = tx.Where("challenges.type = ?", f.Type)
	}
	if f.Category != "" {
		tx = tx.Where("challenges.category = ?", f.Category)
	}
	if f.Name != "" {
		tx = tx.Where("challenges.name ILIKE ?", "%"+f.Name+"%")
	}
	if f.Description != "" {
		tx = tx.Where("challenges.description ILIKE ?", "%"+f.Description+"%")
	}
	return tx
}

// ListChallenges 返回跨比赛聚合的题目列表。
func (c *ChallengeRepo) ListChallenges(limit, offset int, filter ChallengeFilter) ([]model.Challenge, int64, model.RetVal) {
	tx := filter.apply(c.DB.Model(&model.Challenge{}))
	target := fmt.Sprintf("contest_id=%d,contest_name=%s,type=%s,category=%s", filter.ContestID, filter.ContestName, filter.Type, filter.Category)

	var count int64
	if res := tx.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count Challenges: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}

	challenges := make([]model.Challenge, 0)
	if res := filter.apply(c.DB.Model(&model.Challenge{})).
		Order("challenges.id DESC").Limit(limit).Offset(offset).Find(&challenges); res.Error != nil {
		log.Logger.Warningf("Failed to list Challenges: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}
	return challenges, count, model.SuccessRetVal()
}

// ListCategories 返回过滤条件下的题目分类。
func (c *ChallengeRepo) ListCategories(filter ChallengeFilter) ([]string, model.RetVal) {
	categories := make([]string, 0)
	tx := filter.apply(c.DB.Model(&model.Challenge{}))
	target := fmt.Sprintf("contest_id=%d,type=%s", filter.ContestID, filter.Type)
	if res := tx.Distinct().Order("challenges.category ASC").Pluck("challenges.category", &categories); res.Error != nil {
		log.Logger.Warningf("Failed to get Categories: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}
	return categories, model.SuccessRetVal()
}

// ListForUnfinishedContests includes upcoming and running contests, regardless
// of visibility. Use Go duration arithmetic, matching Contest.IsOver exactly.
// Challenges without a contest are included so that their images are warm
// before they get copied into one of those contests.
func (c *ChallengeRepo) ListForUnfinishedContests(now time.Time) ([]model.Challenge, model.RetVal) {
	contests, ret := InitContestRepo(c.DB).ListUnfinished(now)
	if !ret.OK {
		return nil, ret
	}
	if len(contests) == 0 {
		return nil, model.SuccessRetVal()
	}
	contestIDs := make([]uint, 0, len(contests))
	for _, contest := range contests {
		contestIDs = append(contestIDs, contest.ID)
	}
	tx := c.DB.Model(&model.Challenge{}).
		Where("challenges.contest_id IS NULL OR challenges.contest_id IN ?", contestIDs)
	return InitChallengeRepo(tx).FindAll()
}

// ListUnsolvedID 返回队伍尚未解出的题目 id。
func (c *ChallengeRepo) ListUnsolvedID(teamID, contestID uint, category string, limit, offset int) ([]uint, int64, model.RetVal) {
	base := c.DB.Model(&model.Challenge{}).
		Where("challenges.contest_id = ?", contestID).
		Where("challenges.hidden = ?", false).
		Where("EXISTS (SELECT 1 FROM challenge_flags WHERE challenge_flags.challenge_id = challenges.id AND challenge_flags.deleted_at IS NULL AND challenge_flags.active = true)").
		Where(`
			EXISTS (
				SELECT 1
				FROM challenge_flags
				WHERE challenge_flags.challenge_id = challenges.id
					AND challenge_flags.deleted_at IS NULL
					AND challenge_flags.active = true
					AND NOT EXISTS (
						SELECT 1
						FROM team_flags
						WHERE team_flags.team_id = ?
							AND team_flags.challenge_flag_id = challenge_flags.id
							AND team_flags.solved = true
							AND team_flags.deleted_at IS NULL
					)
			)
		`, teamID)
	if category != "" {
		base = base.Where("challenges.category = ?", category)
	}
	conditions := map[string]any{"contest_id": contestID, "team_id": teamID}
	if category != "" {
		conditions["category"] = category
	}
	target := conditionTarget(conditions)

	var count int64
	if res := base.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to list Challenge: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}

	var ids []uint
	res := base.Select("challenges.id").
		Order("challenges.id ASC").
		Limit(limit).Offset(offset).
		Scan(&ids)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list Challenge: %s (target=%s)", res.Error, target)
		return nil, 0, model.FailRetVal(i18n.Model.Challenge.GetError, "Challenge", target, res.Error)
	}
	return ids, count, model.SuccessRetVal()
}

// DeleteByContestID 删除比赛下的全部题目。
func (c *ChallengeRepo) DeleteByContestID(contestIDL ...uint) model.RetVal {
	return c.DeleteByFieldID("contest_id", contestIDL...)
}

// Delete 删除题目及其解题记录。
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
	if ret = InitSubmissionRepo(c.DB).DeleteByChallengeID(challengeIDL...); !ret.OK {
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

// DeleteByID 删除题目及其解题记录。
func (c *ChallengeRepo) DeleteByID(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	var ret model.RetVal
	if ret = InitChallengeFlagRepo(c.DB).DeleteByChallengeID(idL...); !ret.OK {
		return ret
	}
	if ret = InitSubmissionRepo(c.DB).DeleteByChallengeID(idL...); !ret.OK {
		return ret
	}
	if ret = InitGeneratorRepo(c.DB).DeleteByChallengeID(idL...); !ret.OK {
		return ret
	}
	if ret = InitVictimRepo(c.DB).DeleteByChallengeID(idL...); !ret.OK {
		return ret
	}
	target := fmt.Sprintf("id=%v", idL)
	if res := c.DB.Model(&model.Challenge{}).Where("id IN ?", idL).Delete(&model.Challenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.Challenge.DeleteError, "Challenge", target, res.Error)
	}
	return model.SuccessRetVal()
}
