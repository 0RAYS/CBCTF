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
	_, ok, _ := c.GetByName(name)
	return !ok
}

func (c *ContestRepo) getByUniqueKey(key string, value interface{}, preloadL ...string) (model.Contest, bool, string) {
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
	res = preload(res, preloadL...).Limit(1).Find(&contest)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Contest: %s", res.Error)
		return model.Contest{}, false, "GetContestError"
	}
	if res.RowsAffected == 0 {
		return model.Contest{}, false, "ContestNotFound"
	}
	return contest, true, "Success"
}
func (c *ContestRepo) GetByName(name string, preloadL ...string) (model.Contest, bool, string) {
	return c.getByUniqueKey("name", name, preloadL...)
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

func (c *ContestRepo) GetAll(limit, offset int, hidden bool, preloadL ...string) ([]model.Contest, int64, bool, string) {
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
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&contests)
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
		contest, ok, msg := c.GetByID(id)
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

func (c *ContestRepo) Delete(idL ...uint) (bool, string) {
	teamIDL, noticeIDL, usageIDL, flagIDL, submissionIDL := make([]uint, 0), make([]uint, 0), make([]uint, 0), make([]uint, 0), make([]uint, 0)
	for _, id := range idL {
		contest, ok, msg := c.GetByID(id, "Teams", "Notices", "Usages", "Flags", "Submissions")
		if !ok {
			return ok, msg
		}
		for _, team := range contest.Teams {
			teamIDL = append(teamIDL, team.ID)
		}
		for _, notice := range contest.Notices {
			noticeIDL = append(noticeIDL, notice.ID)
		}
		for _, usage := range contest.Usages {
			usageIDL = append(usageIDL, usage.ID)
		}
		for _, flag := range contest.Flags {
			flagIDL = append(flagIDL, flag.ID)
		}
		for _, submission := range contest.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ok, msg := InitTeamRepo(c.DB).Delete(teamIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitNoticeRepo(c.DB).Delete(noticeIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitUsageRepo(c.DB).Delete(usageIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitFlagRepo(c.DB).Delete(flagIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.Contest{}).Where("id IN ?", idL).Delete(&model.Challenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Contest: %v", res.Error)
		return false, "DeleteContestError"
	}
	return true, "Success"
}
