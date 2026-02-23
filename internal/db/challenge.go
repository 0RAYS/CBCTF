package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type ChallengeRepo struct {
	BaseRepo[model.Challenge]
}

type CreateChallengeOptions struct {
	RandID          string
	Name            string
	Description     string
	Category        string
	Type            model.ChallengeType
	GeneratorImage  string
	Options         model.Options
	NetworkPolicies model.NetworkPolicies
}

func (c CreateChallengeOptions) Convert2Model() model.Model {
	return model.Challenge{
		RandID:          c.RandID,
		Name:            c.Name,
		Description:     c.Description,
		Category:        c.Category,
		Type:            c.Type,
		GeneratorImage:  c.GeneratorImage,
		Options:         c.Options,
		NetworkPolicies: c.NetworkPolicies,
	}
}

type UpdateChallengeOptions struct {
	Name            *string
	Description     *string
	Category        *string
	GeneratorImage  *string
	Options         *model.Options
	NetworkPolicies *model.NetworkPolicies
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
	if u.Options != nil {
		options["options"] = u.Options
	}
	if u.NetworkPolicies != nil {
		options["network_policies"] = *u.NetworkPolicies
	}
	return options
}

func InitChallengeRepo(tx *gorm.DB) *ChallengeRepo {
	return &ChallengeRepo{
		BaseRepo: BaseRepo[model.Challenge]{
			DB: tx,
		},
	}
}

func (c *ChallengeRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.Challenge, model.RetVal) {
	return c.GetByUniqueKey("rand_id", randID, optionsL...)
}

func (c *ChallengeRepo) ListCategories(t model.ChallengeType) ([]string, model.RetVal) {
	var categories = make([]string, 0)
	res := c.DB.Model(&model.Challenge{})
	if t != "" {
		res = res.Where("type = ?", t)
	}
	res = res.Select("distinct category").Find(&categories)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Categories: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Challenge{}.ModelName(), "Error": res.Error.Error()}}
	}
	return categories, model.SuccessRetVal()
}

func (c *ChallengeRepo) ListChallengesNotInContest(contestID uint, limit, offset int, category string, t model.ChallengeType) ([]model.Challenge, int64, model.RetVal) {
	base := c.DB.Table("challenges").
		Joins("LEFT JOIN contest_challenges ON challenges.id = contest_challenges.challenge_id AND contest_challenges.contest_id = ? AND contest_challenges.deleted_at IS NULL", contestID).
		Where("contest_challenges.id IS NULL AND challenges.deleted_at IS NULL")

	if category != "" {
		base = base.Where("challenges.category = ?", category)
	}
	if t != "" {
		base = base.Where("challenges.type = ?", t)
	}

	var count int64
	if res := base.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count Challenges not in contest: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Challenge{}.ModelName(), "Error": res.Error.Error()}}
	}

	var challenges = make([]model.Challenge, 0)
	if res := base.Select("challenges.*").Order("challenges.id DESC").Limit(limit).Offset(offset).Scan(&challenges); res.Error != nil {
		log.Logger.Warningf("Failed to list Challenges not in contest: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Challenge{}.ModelName(), "Error": res.Error.Error()}}
	}

	return challenges, count, model.SuccessRetVal()
}

func (c *ChallengeRepo) Delete(randIDL ...string) model.RetVal {
	challengeL, _, ret := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"rand_id": randIDL},
		Preloads: map[string]GetOptions{
			"Dockers":           {},
			"ChallengeFlags":    {},
			"ContestChallenges": {},
			"Submissions":       {},
		},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	dockerIDL, challengeFlagIDL, contestChallengeIDL, submissionIDL := make([]uint, 0), make([]uint, 0), make([]uint, 0), make([]uint, 0)
	for _, challenge := range challengeL {
		for _, docker := range challenge.Dockers {
			dockerIDL = append(dockerIDL, docker.ID)
		}
		for _, challengeFlag := range challenge.ChallengeFlags {
			challengeFlagIDL = append(challengeFlagIDL, challengeFlag.ID)
		}
		for _, contestChallenge := range challenge.ContestChallenges {
			contestChallengeIDL = append(contestChallengeIDL, contestChallenge.ID)
		}
		for _, submission := range challenge.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ret = InitDockerRepo(c.DB).Delete(dockerIDL...); !ret.OK {
		return ret
	}
	if ret = InitChallengeFlagRepo(c.DB).Delete(challengeFlagIDL...); !ret.OK {
		return ret
	}
	if ret = InitContestChallengeRepo(c.DB).Delete(contestChallengeIDL...); !ret.OK {
		return ret
	}
	if ret = InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.Challenge{}).Where("rand_id IN ?", randIDL).Delete(&model.Challenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": model.Challenge{}.ModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
