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
	Exposes           model.Exposes
	Start             time.Time
	Duration          time.Duration
	Image             string
	PodName           string
	ContainerName     string
	ServiceName       string
	NetworkPolicyName string
	NetworkPolicies   model.NetworkPolicies
	Flags             model.Strings
}

type UpdateContainerOptions struct {
	IP       *string        `json:"ip"`
	Duration *time.Duration `json:"duration"`
}

func InitContainerRepo(tx *gorm.DB) *ContainerRepo {
	return &ContainerRepo{Repo: Repo[model.Container]{DB: tx, Model: "Container"}}
}

func (c *ContainerRepo) IsUniqueContainer(usageID, teamID uint) bool {
	res := c.DB.Model(&model.Container{}).Where("usage_id = ? AND team_id = ?", usageID, teamID).Limit(1).Find(&model.Container{})
	return res.RowsAffected == 0
}

//func (c *ContainerRepo) Create(options CreateContainerOptions) (model.Container, bool, string) {
//	container, err := utils.S2S[model.Container](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Container: %s", err)
//		return model.Container{}, false, "Options2ModelError"
//	}
//	res := c.DB.Model(&model.Container{}).Create(&container)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create Container: %s", res.Error)
//		return model.Container{}, false, "CreateContainerError"
//	}
//	return container, true, "Success"
//}

//func (c *ContainerRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Container, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Container{}, false, "UnsupportedKey"
//	}
//	var container model.Container
//	res := c.DB.Model(&model.Container{}).Where(key+" = ?", key)
//	res = model.GetPreload(res, model.Container{}, preload, depth).Limit(1).Find(&container)
//	if res.RowsAffected == 0 {
//		return model.Container{}, false, "ContainerNotFound"
//	}
//	return container, true, "Success"
//}

//func (c *ContainerRepo) GetByID(id uint, preload bool, depth int) (model.Container, bool, string) {
//	return c.getByUniqueKey("id", id, preload, depth)
//}

func (c *ContainerRepo) Count(teamID uint) (int64, bool, string) {
	var count int64
	res := c.DB.Model(&model.Container{}).Where("team_id = ?", teamID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Containers: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (c *ContainerRepo) GetAll(teamID uint, limit, offset int, preload bool, depth int) ([]model.Container, int64, bool, string) {
	var (
		containers     = make([]model.Container, 0)
		count, ok, msg = c.Count(teamID)
	)
	if !ok {
		return containers, count, false, msg
	}
	res := c.DB.Model(&model.Container{}).Where("team_id = ?", teamID)
	res = model.GetPreload(res, c.Model, preload, depth).Limit(limit).Offset(offset).Find(&containers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Containers: %s", res.Error)
		return containers, count, false, "GetContainerError"
	}
	return containers, count, true, "Success"
}

func (c *ContainerRepo) GetBy2ID(teamID uint, usageID uint, preload bool, depth int) ([]model.Container, bool, string) {
	containers := make([]model.Container, 0)
	res := c.DB.Model(&model.Container{}).Where("team_id = ? AND usage_id = ?", teamID, usageID)
	res = model.GetPreload(res, c.Model, preload, depth).Limit(1).Find(&containers)
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
		container, ok, msg := c.GetByID(id, false, 0)
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

//func (c *ContainerRepo) Delete(idL ...uint) (bool, string) {
//	res := c.DB.Model(&model.Container{}).Where("id IN ?", idL).Delete(&model.Container{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Container: %s", res.Error)
//		return false, "DeleteContainerError"
//	}
//	return true, "Success"
//}
