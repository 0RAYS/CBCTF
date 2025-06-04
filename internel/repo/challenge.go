package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type ChallengeRepo struct {
	Repo[model.Challenge]
}

type CreateChallengeOptions struct {
	ID        string
	Name      string
	Desc      string
	Category  string
	Type      string
	Generator string
	Flags     model.StringList
	Dockers   model.Dockers
}

type UpdateChallengeOptions struct {
	Name      *string           `json:"name"`
	Desc      *string           `json:"desc"`
	Category  *string           `json:"category"`
	Type      *string           `json:"type"`
	Generator *string           `json:"generator"`
	Flags     *model.StringList `json:"flags"`
	Dockers   *model.Dockers    `json:"dockers"`
}

func InitChallengeRepo(tx *gorm.DB) *ChallengeRepo {
	return &ChallengeRepo{
		Repo: Repo[model.Challenge]{
			DB: tx, Model: "Challenge",
			CreateError:   i18n.CreateChallengeError,
			DeleteError:   i18n.DeleteChallengeError,
			GetError:      i18n.GetChallengeError,
			NotFoundError: i18n.ChallengeNotFound,
		},
	}
}

func (c *ChallengeRepo) getByUniqueKey(key string, value any, preloadL ...string) (model.Challenge, bool, string) {
	switch key {
	case "id":
		value = value.(string)
	default:
		return model.Challenge{}, false, i18n.UnsupportedKey
	}
	var challenge model.Challenge
	res := c.DB.Model(&model.Challenge{}).Where(key+" = ?", value)
	res = preload(res, preloadL...).Limit(1).Find(&challenge)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Challenge: %s", res.Error)
		return model.Challenge{}, false, i18n.GetChallengeError
	}
	if res.RowsAffected == 0 {
		return model.Challenge{}, false, i18n.ChallengeNotFound
	}
	return challenge, true, i18n.Success
}

func (c *ChallengeRepo) GetByID(id string, preloadL ...string) (model.Challenge, bool, string) {
	return c.getByUniqueKey("id", id, preloadL...)
}

func (c *ChallengeRepo) Count(t, category string) (int64, bool, string) {
	var count int64
	res := c.DB.Model(&model.Challenge{})
	if t != "" && category != "" {
		res = res.Where("type = ? AND category = ?", t, category)
	} else if !(t == "" && category == "") {
		res = res.Where("type = ? OR category = ?", t, category)
	}
	res = res.Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Challenges: %s", res.Error)
		return 0, false, i18n.CountModelError
	}
	return count, true, i18n.Success
}

func (c *ChallengeRepo) GetAll(limit, offset int, t, category string, preloadL ...string) ([]model.Challenge, int64, bool, string) {
	var (
		challenges     = make([]model.Challenge, 0)
		count, ok, msg = c.Count(t, category)
	)
	if !ok {
		return challenges, count, false, msg
	}
	res := c.DB.Model(&model.Challenge{})
	if t != "" && category != "" {
		res = res.Where("type = ? AND category = ?", t, category)
	} else if !(t == "" && category == "") {
		res = res.Where("type = ? OR category = ?", t, category)
	}
	res = preload(res, preloadL...)
	res = res.Order("created_at DESC").Limit(limit).Offset(offset).Find(&challenges)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Challenges: %s", res.Error)
		return challenges, count, false, i18n.GetChallengeError
	}
	return challenges, count, true, i18n.Success
}

func (c *ChallengeRepo) Update(id string, options UpdateChallengeOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Challenge: too many times failed due to optimistic lock")
			return false, i18n.DeadLock
		}
		challenge, ok, msg := c.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = challenge.Version + 1
		res := c.DB.Model(&model.Challenge{}).Where("id = ? AND version = ?", id, challenge.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Challenge: %s", res.Error)
			return false, i18n.UpdateChallengeError
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}

func (c *ChallengeRepo) GetCategories(t string) ([]string, bool, string) {
	var categories = make([]string, 0)
	res := c.DB.Model(&model.Challenge{})
	if t != "" {
		res = res.Where("type = ?", t)
	}
	res = res.Select("distinct category").Find(&categories)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Categories: %s", res.Error)
		return categories, false, i18n.GetChallengeError
	}
	return categories, true, i18n.Success
}

func (c *ChallengeRepo) Delete(idL ...string) (bool, string) {
	usageIDL, submissionIDL := make([]uint, 0), make([]uint, 0)
	for _, id := range idL {
		challenge, ok, msg := c.GetByID(id, "Usages", "Submissions")
		if !ok {
			return false, msg
		}
		for _, usage := range challenge.Usages {
			usageIDL = append(usageIDL, usage.ID)
		}
		for _, submission := range challenge.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ok, msg := InitUsageRepo(c.DB).Delete(usageIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	res := c.DB.Model(&model.Challenge{}).Where("id IN ?", idL).Delete(&model.Challenge{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %s", res.Error)
		return false, i18n.DeleteChallengeError
	}
	return true, i18n.Success
}
