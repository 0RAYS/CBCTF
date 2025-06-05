package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type ChallengeRepo struct {
	Basic[model.Challenge]
}

type CreateChallengeOptions struct {
	RandID         string
	Name           string
	Desc           string
	Category       string
	Type           string
	GeneratorImage string
}

func (c CreateChallengeOptions) Convert2Model() model.Model {
	return model.Challenge{
		RandID:         c.RandID,
		Name:           c.Name,
		Desc:           c.Desc,
		Category:       c.Category,
		Type:           c.Type,
		GeneratorImage: c.GeneratorImage,
	}
}

type UpdateChallengeOptions struct {
	Name           *string
	Desc           *string
	Category       *string
	Type           *string
	GeneratorImage *string
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
	if u.Type != nil {
		options["type"] = *u.Type
	}
	if u.GeneratorImage != nil {
		options["generator_image"] = *u.GeneratorImage
	}
	return options
}

func InitChallengeRepo(tx *gorm.DB) *ChallengeRepo {
	return &ChallengeRepo{
		Basic: Basic[model.Challenge]{
			DB: tx,
		},
	}
}

func (c *ChallengeRepo) GetByRandID(randID string, preloadL ...string) (model.Challenge, bool, string) {
	return c.getUniqueByKey("rand_id", randID, preloadL...)
}

func (c *ChallengeRepo) Delete(randIDL ...string) (bool, string) {
	dockerGroupIDL, challengeFlagIDL, contestChallengeIDL := make([]uint, 0), make([]uint, 0), make([]uint, 0)
	for _, randID := range randIDL {
		challenge, ok, msg := c.GetByRandID(randID, "DockerGroups", "ChallengeFlags", "ContestChallenges")
		if !ok {
			return false, msg
		}
		for _, dockerGroup := range challenge.DockerGroups {
			dockerGroupIDL = append(dockerGroupIDL, dockerGroup.ID)
		}
		for _, challengeFlag := range challenge.ChallengeFlags {
			challengeFlagIDL = append(challengeFlagIDL, challengeFlag.ID)
		}
		for _, contestChallenge := range challenge.ContestChallenges {
			contestChallengeIDL = append(contestChallengeIDL, contestChallenge.ID)
		}
	}
	if ok, msg := InitDockerGroupRepo(c.DB).Delete(dockerGroupIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitChallengeFlagRepo(c.DB).Delete(challengeFlagIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitContestChallengeRepo(c.DB).Delete(contestChallengeIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.Challenge{}).Where("rand_id IN ?", randIDL).Delete(&model.Challenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %v", res.Error)
		return false, model.Challenge{}.DeleteErrorString()
	}
	return true, i18n.Success
}
