package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"time"
)

type ContestRepo struct {
	Repo[model.Contest]
}

type CreateContestOptions struct {
	Name      string
	Desc      string
	Captcha   string
	Avatar    string
	Prefix    string
	Size      int
	Start     time.Time
	Duration  time.Duration
	Blood     bool
	Hidden    bool
	Rules     model.Strings
	Prizes    model.Prizes
	Timelines model.Timelines
}

type UpdateContestOptions struct {
	Name      *string          `json:"name"`
	Desc      *string          `json:"desc"`
	Captcha   *string          `json:"captcha"`
	Avatar    *string          `json:"avatar"`
	Prefix    *string          `json:"prefix"`
	Size      *int             `json:"size"`
	Start     *time.Time       `json:"start"`
	Duration  *time.Duration   `json:"duration"`
	Blood     *bool            `json:"blood"`
	Hidden    *bool            `json:"hidden"`
	Rules     *model.Strings   `json:"rules"`
	Prizes    *model.Prizes    `json:"prizes"`
	Timelines *model.Timelines `json:"timelines"`
}

func InitContestRepo(tx *gorm.DB) *ContestRepo {
	return &ContestRepo{Repo: Repo[model.Contest]{DB: tx}}
}

func (c *ContestRepo) IsUniqueName(name string) bool {
	_, ok, _ := c.GetByName(name, false, 0)
	return !ok
}

//func (c *ContestRepo) Create(options CreateContestOptions) (model.Contest, bool, string) {
//	contest, err := utils.S2S[model.Contest](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Contest: %s", err)
//		return model.Contest{}, false, "Options2ModelError"
//	}
//	res := c.DB.Model(&model.Contest{}).Create(&contest)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create Contest: %s", res.Error)
//		return model.Contest{}, false, "CreateContestError"
//	}
//	return contest, true, "Success"
//}

func (c *ContestRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Contest, bool, string) {
	switch key {
	case "name":
		value = value.(string)
	case "id":
		value = value.(uint)
	default:
		return model.Contest{}, false, "UnsupportedKey"
	}
	var contest model.Contest
	res := c.DB.Model(&model.Contest{}).Where(key+" = ?", value)
	res = model.GetPreload(res, c.Model, preload, depth).Limit(1).Find(&contest)
	if res.RowsAffected == 0 {
		return model.Contest{}, false, "ContestNotFound"
	}
	return contest, true, "Success"
}

//func (c *ContestRepo) GetByID(id uint, preload bool, depth int) (model.Contest, bool, string) {
//	return c.getByUniqueKey("id", id, preload, depth)
//}

func (c *ContestRepo) GetByName(name string, preload bool, depth int) (model.Contest, bool, string) {
	return c.getByUniqueKey("name", name, preload, depth)
}

func (c *ContestRepo) Count(hidden bool) (int64, bool, string) {
	var count int64
	res := c.DB.Model(&model.Contest{})
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	res = res.Count(&count)
	if res.Error != nil {
		log.Logger.Errorf("Failed to count Contests: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (c *ContestRepo) GetAll(limit, offset int, preload bool, depth int, hidden bool) ([]model.Contest, int64, bool, string) {
	var (
		contests       = make([]model.Contest, 0)
		count, ok, msg = c.Count(hidden)
	)
	if !ok {
		return contests, count, false, msg
	}
	res := c.DB.Model(&model.Contest{})
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	res = model.GetPreload(res, c.Model, preload, depth).Limit(limit).Offset(offset).Find(&contests)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get Contests: %s", res.Error)
		return contests, count, false, "GetContestsError"
	}
	return contests, count, true, "Success"
}

func (c *ContestRepo) Update(id uint, options UpdateContestOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Contest: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		contest, ok, msg := c.GetByID(id, false, 0)
		if !ok {
			return ok, msg
		}
		data["version"] = contest.Version + 1
		res := c.DB.Model(&model.Contest{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, contest.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Contest: %v", res.Error)
			return false, "UpdateContestError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

//func (c *ContestRepo) Delete(idL ...uint) (bool, string) {
//	res := c.DB.Model(&model.Contest{}).Where("id IN ?", idL).Delete(&model.Contest{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Contest: %s", res.Error)
//		return false, "DeleteContestError"
//	}
//	return true, "Success"
//}
