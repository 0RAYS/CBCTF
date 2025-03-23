package db

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

// CreateContainer 创建 Container, 并注入 flag
func CreateContainer(tx *gorm.DB, flag model.Flag, usage model.Usage, creatorID uint) (model.Container, bool, string) {
	var (
		container model.Container
		ok        bool
		msg       string
		ip        string
		port      int32
	)
	if container, ok, _ = GetContainerBy3ID(tx, flag.ContestID, flag.TeamID, flag.ChallengeID); ok {
		return container, ok, "Success"
	}
	if usage.Type != model.Docker {
		return model.Container{}, false, "InvalidChallengeType"
	}
	container = model.InitContainer(flag, usage, creatorID)
	log.Logger.Debugf("Starting container for team %d challenge %s", flag.TeamID, flag.ChallengeID)
	ip, port, ok, msg = k8s.StartContainer(usage, flag, container)
	if !ok {
		go k8s.StopContainer(container)
		log.Logger.Warningf("Failed to start container for challenge %s: %s", flag.ChallengeID, msg)
		return model.Container{}, false, msg
	}
	container.IP = ip
	container.Port = port
	res := tx.Model(&model.Container{}).Create(&container)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Container: %s", res.Error)
		return model.Container{}, false, "CreateContainerError"
	}
	return container, true, "Success"
}

// GetContainers 获取所有 Container, deleted 为 true 时获取已删除的 Container
func GetContainers(tx *gorm.DB, deleted bool) ([]model.Container, bool, string) {
	var containers []model.Container
	res := tx.Model(&model.Container{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Find(&containers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Containers: %s", res.Error)
		return make([]model.Container, 0), false, "GetContainersError"
	}
	return containers, true, "Success"
}

// GetContainerByID 根据 ID 获取 Container
func GetContainerByID(tx *gorm.DB, id uint, deleted bool) (model.Container, bool, string) {
	var container model.Container
	res := tx.Model(&model.Container{}).Where("id = ?", id)
	if deleted {
		res = res.Unscoped()
	}
	res = res.Find(&container).Limit(1)
	if res.RowsAffected != 1 {
		return model.Container{}, false, "ContainerNotFound"
	}
	return container, true, "Success"
}

// GetContainerByPodName 根据 podName 获取 Container
func GetContainerByPodName(tx *gorm.DB, podName string) (model.Container, bool, string) {
	var container model.Container
	res := tx.Model(&model.Container{}).Where("pod = ?", podName).Find(&container).Limit(1)
	if res.RowsAffected != 1 {
		return model.Container{}, false, "ContainerNotFound"
	}
	return container, true, "Success"
}

// GetContainerByTeamID 根据 teamID 获取 Container
func GetContainerByTeamID(tx *gorm.DB, teamID uint, limit, offset int, deleted bool) ([]model.Container, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	res := tx.Model(&model.Container{}).Where("team_id = ?", teamID)
	if deleted {
		res = res.Unscoped()
	}
	var containers []model.Container
	var count int64
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count Containers: %s", err)
		return make([]model.Container, 0), -1, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	res = res.Limit(limit).Offset(offset).Find(&containers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Containers: %s", res.Error)
		return make([]model.Container, 0), -1, false, "GetContainersError"
	}
	return containers, count, true, "Success"
}

// GetContainerBy3ID 根据 contestID, teamID, challengeID 获取 Container
func GetContainerBy3ID(tx *gorm.DB, contestID, teamID uint, challengeID string) (model.Container, bool, string) {
	var container model.Container
	res := tx.Model(&model.Container{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&container).Limit(1)
	if res.RowsAffected != 1 {
		return model.Container{}, false, "ContainerNotFound"
	}
	return container, true, "Success"
}

// DeleteContainer 删除 Container
func DeleteContainer(tx *gorm.DB, container model.Container) (bool, string) {
	log.Logger.Infof("Stopping container for team %d challenge %s", container.TeamID, container.ChallengeID)
	ok, msg := k8s.StopContainer(container)
	if !ok {
		log.Logger.Warningf("Failed to stop container for challenge %s: %s", container.ChallengeID, msg)
		return false, "DeleteContainerError"
	}
	if ok, msg = UpdateContainer(tx, container.ID, map[string]interface{}{"duration": time.Now().Sub(container.Start)}); !ok {
		return false, msg
	}
	res := tx.Model(&model.Container{}).
		Where("id = ?", container.ID).Delete(&model.Container{})
	if res.Error != nil {
		return false, "DeleteContainerError"
	}
	return true, "Success"
}

// UpdateContainer 更新 Container, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateContainer(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	var count int
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed too many times to update user due to optimistic lock")
			return false, "FailedTooManyTimes"
		}
		var container model.Container
		res := tx.Model(&model.Container{}).Where("id = ?", id).Find(&container).Limit(1)
		if res.RowsAffected != 1 {
			return false, "ContainerNotFound"
		}
		res = tx.Model(&container).Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Container: %v", res.Error)
			return false, "UpdateContainerError"
		}
		if res.RowsAffected == 0 {
			log.Logger.Debug("Failed to update Container due to optimistic lock")
			continue
		}
		break
	}
	return true, "Success"
}
