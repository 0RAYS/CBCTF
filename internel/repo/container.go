package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"time"
)

type ContainerRepo struct {
	Repo[model.Container]
}

type CreateContainerOptions struct {
	UsageID           uint
	TeamID            uint
	UserID            uint
	PodIP             string
	IPBlock           string
	Exposes           model.Exposes
	Start             time.Time
	Duration          time.Duration
	Image             string
	Hostname          string
	PodName           string
	ContainerName     string
	ServiceName       string
	NetworkPolicyName string
	NetworkPolicies   model.NetworkPolicies
	Flags             model.Strings
}

type UpdateContainerOptions struct {
	IP       *string        `json:"ip"`
	PodIP    *string        `json:"pod_ip"`
	Duration *time.Duration `json:"duration"`
}

func InitContainerRepo(tx *gorm.DB) *ContainerRepo {
	return &ContainerRepo{Repo: Repo[model.Container]{DB: tx, Model: "Container"}}
}

func (c *ContainerRepo) Count(teamID uint, deleted bool) (int64, bool, string) {
	var count int64
	res := c.DB.Model(&model.Container{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("team_id = ?", teamID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Containers: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (c *ContainerRepo) GetByTeam(teamID uint, limit, offset int, deleted bool, preloadL ...string) ([]model.Container, int64, bool, string) {
	var (
		containers     = make([]model.Container, 0)
		count, ok, msg = c.Count(teamID, deleted)
	)
	if !ok {
		return containers, count, false, msg
	}
	res := c.DB.Model(&model.Container{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("team_id = ?", teamID)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&containers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Containers: %s", res.Error)
		return containers, count, false, "GetContainerError"
	}
	return containers, count, true, "Success"
}

func (c *ContainerRepo) GetBy2ID(teamID uint, usageID uint, deleted bool, preloadL ...string) ([]model.Container, bool, string) {
	containers := make([]model.Container, 0)
	res := c.DB.Model(&model.Container{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("team_id = ? AND usage_id = ?", teamID, usageID)
	res = preload(res, preloadL...).Find(&containers)
	if res.RowsAffected == 0 {
		return containers, false, "ContainerNotFound"
	}
	return containers, true, "Success"
}

func (c *ContainerRepo) GetByName(key, value string, deleted bool, preloadL ...string) ([]model.Container, bool, string) {
	containers := make([]model.Container, 0)
	res := c.DB.Model(&model.Container{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where(key+" = ?", value)
	res = preload(res, preloadL...).Find(&containers)
	if res.RowsAffected == 0 {
		return containers, false, "ContainerNotFound"
	}
	return containers, true, "Success"
}

func (c *ContainerRepo) Update(id uint, options UpdateContainerOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Container: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		container, ok, msg := c.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = container.Version + 1
		res := c.DB.Model(&model.Container{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, container.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Container: %s", res.Error)
			return false, "UpdateContainerError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}
