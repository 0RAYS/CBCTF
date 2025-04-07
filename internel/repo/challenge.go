package repo

import (
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
	Flags     model.Strings
	Docker    model.Docker
	Dockers   model.Dockers
}

type UpdateChallengeOptions struct {
	Name      *string        `json:"name"`
	Desc      *string        `json:"desc"`
	Category  *string        `json:"category"`
	Type      *string        `json:"type"`
	Generator *string        `json:"generator"`
	Flags     *model.Strings `json:"flags"`
	Docker    *model.Docker  `json:"docker"`
	Dockers   *model.Dockers `json:"dockers"`
}

func InitChallengeRepo(tx *gorm.DB) *ChallengeRepo {
	return &ChallengeRepo{Repo: Repo[model.Challenge]{DB: tx, Model: "Challenge"}}
}

func (c *ChallengeRepo) getByUniqueKey(key string, value interface{}, preload bool, nestedL ...string) (model.Challenge, bool, string) {
	switch key {
	case "id":
		value = value.(string)
	default:
		return model.Challenge{}, false, "UnsupportedKey"
	}
	var challenge model.Challenge
	res := c.DB.Model(&model.Challenge{}).Where(key+" = ?", value)
	res = model.GetPreload(res, preload, nestedL...).Limit(1).Find(&challenge)
	if res.RowsAffected == 0 {
		return model.Challenge{}, false, "ChallengeNotFound"
	}
	return challenge, true, "Success"
}

func (c *ChallengeRepo) GetByID(id string, preload bool, nestedL ...string) (model.Challenge, bool, string) {
	return c.getByUniqueKey("id", id, preload, nestedL...)
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
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (c *ChallengeRepo) GetAll(limit, offset int, t, category string, preload bool, nestedL ...string) ([]model.Challenge, int64, bool, string) {
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
	res = model.GetPreload(res, preload, nestedL...)
	res = res.Limit(limit).Offset(offset).Find(&challenges)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Challenges: %s", res.Error)
		return challenges, count, false, "GetChallengeError"
	}
	return challenges, count, true, "Success"
}

func (c *ChallengeRepo) Update(id string, options UpdateChallengeOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Challenge: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		challenge, ok, msg := c.GetByID(id, false)
		if !ok {
			return ok, msg
		}
		data["version"] = challenge.Version + 1
		res := c.DB.Model(&model.Challenge{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, challenge.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Challenge: %s", res.Error)
			return false, "UpdateChallengeError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

func (c *ChallengeRepo) Delete(idL ...string) (bool, string) {
	res := c.DB.Model(&model.Challenge{}).Where("id IN ?", idL).Delete(&model.Challenge{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %s", res.Error)
		return false, "DeleteChallengeError"
	}
	return true, "Success"
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
		return categories, false, "GetChallengeError"
	}
	return categories, true, "Success"
}
