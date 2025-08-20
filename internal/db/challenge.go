package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type ChallengeRepo struct {
	BasicRepo[model.Challenge]
}

type CreateChallengeOptions struct {
	RandID          string
	Name            string
	Desc            string
	Category        string
	Type            string
	GeneratorImage  string
	Options         model.Options
	NetworkPolicies model.NetworkPolicies
}

func (c CreateChallengeOptions) Convert2Model() model.Model {
	return model.Challenge{
		RandID:          c.RandID,
		Name:            c.Name,
		Desc:            c.Desc,
		Category:        c.Category,
		Type:            c.Type,
		GeneratorImage:  c.GeneratorImage,
		Options:         c.Options,
		NetworkPolicies: c.NetworkPolicies,
	}
}

type UpdateChallengeOptions struct {
	Name            *string
	Desc            *string
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
	if u.Desc != nil {
		options["desc"] = *u.Desc
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
		BasicRepo: BasicRepo[model.Challenge]{
			DB: tx,
		},
	}
}

func (c *ChallengeRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.Challenge, bool, string) {
	return c.GetByUniqueKey("rand_id", randID, optionsL...)
}

func (c *ChallengeRepo) ListCategories(t string) ([]string, bool, string) {
	var categories = make([]string, 0)
	res := c.DB.Model(&model.Challenge{})
	if t != "" {
		res = res.Where("type = ?", t)
	}
	res = res.Select("distinct category").Find(&categories)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Categories: %s", res.Error)
		return nil, false, i18n.GetChallengeError
	}
	return categories, true, i18n.Success
}

func (c *ChallengeRepo) Delete(randIDL ...string) (bool, string) {
	challengeL, _, ok, msg := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"rand_id": randIDL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"Dockers":           {Selects: []string{"id", "challenge_id"}},
			"ChallengeFlags":    {Selects: []string{"id", "challenge_id"}},
			"ContestChallenges": {Selects: []string{"id", "challenge_id"}},
			"Submissions":       {Selects: []string{"id", "challenge_id"}},
		},
	})
	if !ok && msg != i18n.ChallengeNotFound {
		return false, msg
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
	if ok, msg = InitDockerRepo(c.DB).Delete(dockerIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitChallengeFlagRepo(c.DB).Delete(challengeFlagIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitContestChallengeRepo(c.DB).Delete(contestChallengeIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.Challenge{}).Where("rand_id IN ?", randIDL).Delete(&model.Challenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %s", res.Error)
		return false, i18n.DeleteChallengeError
	}
	return true, i18n.Success
}
