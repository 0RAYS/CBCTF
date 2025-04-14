package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"time"
)

type VictimRepo struct {
	Repo[model.Victim]
}

type CreateVictimOptions struct {
	UsageID  uint
	TeamID   uint
	UserID   uint
	IPBlock  string
	Start    time.Time
	Duration time.Duration
}

type UpdateVictimOptions struct {
	Duration *time.Duration `json:"duration"`
}

func InitVictimRepo(tx *gorm.DB) *VictimRepo {
	return &VictimRepo{Repo: Repo[model.Victim]{DB: tx, Model: "Victim"}}
}

func (v *VictimRepo) Count(deleted bool) (int64, bool, string) {
	var count int64
	res := v.DB.Model(&model.Victim{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Victims: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (v *VictimRepo) CountByTeam(teamID uint, deleted bool) (int64, bool, string) {
	var count int64
	res := v.DB.Model(&model.Victim{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("team_id = ?", teamID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Victims: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (v *VictimRepo) GetByTeam(teamID uint, limit, offset int, deleted bool, preloadL ...string) ([]model.Victim, int64, bool, string) {
	var (
		victims        = make([]model.Victim, 0)
		count, ok, msg = v.CountByTeam(teamID, deleted)
	)
	if !ok {
		return victims, count, false, msg
	}
	res := v.DB.Model(&model.Victim{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("team_id = ?", teamID)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&victims)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Victims: %s", res.Error)
		return victims, count, false, "GetVictimError"
	}
	return victims, count, true, "Success"
}

// GetBy2ID ok == true, 必存在一个 victim
func (v *VictimRepo) GetBy2ID(teamID, usageID uint, deleted bool, preloadL ...string) ([]model.Victim, bool, string) {
	victims := make([]model.Victim, 0)
	res := v.DB.Model(&model.Victim{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("team_id = ? AND usage_id = ?", teamID, usageID)
	res = preload(res, preloadL...).Find(&victims)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Victims: %s", res.Error)
		return victims, false, "GetVictimError"
	}
	if res.RowsAffected == 0 {
		return victims, false, "VictimNotFound"
	}
	return victims, true, "Success"
}

func (v *VictimRepo) GetAll(limit, offset int, deleted bool, preloadL ...string) ([]model.Victim, int64, bool, string) {
	var (
		victims        = make([]model.Victim, 0)
		count, ok, msg = v.Count(deleted)
	)
	if !ok {
		return victims, count, false, msg
	}
	res := v.DB.Model(&model.Victim{})
	if deleted {
		res = res.Unscoped()
	}
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&victims)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Victims: %s", res.Error)
		return victims, count, false, "GetVictimError"
	}
	return victims, count, true, "Success"
}

func (v *VictimRepo) Update(id uint, options UpdateVictimOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Victim: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		victim, ok, msg := v.GetByID(id)
		if !ok {
			return false, msg
		}
		data["version"] = victim.Version + 1
		res := v.DB.Model(&model.Victim{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, victim.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Victim: %s", res.Error)
			return false, "UpdateVictimError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}
