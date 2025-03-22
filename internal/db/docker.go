package db

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

// CreateDocker 创建 Docker, 并注入 flag
func CreateDocker(tx *gorm.DB, flag model.Flag, usage model.Usage, creatorID uint) (model.Container, bool, string) {
	var (
		docker model.Container
		ok     bool
		msg    string
		ip     string
		port   int32
	)
	if docker, ok, _ = GetDockerBy3ID(tx, flag.ContestID, flag.TeamID, flag.ChallengeID); ok {
		return docker, ok, "Success"
	}
	if usage.Type != model.Docker {
		return model.Container{}, false, "InvalidChallengeType"
	}
	docker = model.InitDocker(flag, usage, creatorID)
	log.Logger.Debugf("Starting container for team %d challenge %s", flag.TeamID, flag.ChallengeID)
	ip, port, ok, msg = k8s.StartContainer(usage, flag, docker)
	if !ok {
		go k8s.StopContainer(docker)
		log.Logger.Warningf("Failed to start container for challenge %s: %s", flag.ChallengeID, msg)
		return model.Container{}, false, msg
	}
	docker.IP = ip
	docker.Port = port
	res := tx.Model(&model.Container{}).Create(&docker)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Docker: %s", res.Error)
		return model.Container{}, false, "CreateDockerError"
	}
	return docker, true, "Success"
}

// GetDockers 获取所有 Docker, deleted 为 true 时获取已删除的 Docker
func GetDockers(tx *gorm.DB, deleted bool) ([]model.Container, bool, string) {
	var dockers []model.Container
	res := tx.Model(&model.Container{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Find(&dockers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Dockers: %s", res.Error)
		return make([]model.Container, 0), false, "GetDockersError"
	}
	return dockers, true, "Success"
}

// GetDockerByID 根据 ID 获取 Docker
func GetDockerByID(tx *gorm.DB, id uint, deleted bool) (model.Container, bool, string) {
	var docker model.Container
	res := tx.Model(&model.Container{}).Where("id = ?", id)
	if deleted {
		res = res.Unscoped()
	}
	res = res.Find(&docker).Limit(1)
	if res.RowsAffected != 1 {
		return model.Container{}, false, "DockerNotFound"
	}
	return docker, true, "Success"
}

// GetDockerByPodName 根据 podName 获取 Docker
func GetDockerByPodName(tx *gorm.DB, podName string) (model.Container, bool, string) {
	var docker model.Container
	res := tx.Model(&model.Container{}).Where("pod = ?", podName).Find(&docker).Limit(1)
	if res.RowsAffected != 1 {
		return model.Container{}, false, "DockerNotFound"
	}
	return docker, true, "Success"
}

// GetDockerByTeamID 根据 teamID 获取 Docker
func GetDockerByTeamID(tx *gorm.DB, teamID uint, limit, offset int, deleted bool) ([]model.Container, int64, bool, string) {
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
	var dockers []model.Container
	var count int64
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count Dockers: %s", err)
		return make([]model.Container, 0), -1, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	res = res.Limit(limit).Offset(offset).Find(&dockers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Dockers: %s", res.Error)
		return make([]model.Container, 0), -1, false, "GetDockersError"
	}
	return dockers, count, true, "Success"
}

// GetDockerBy3ID 根据 contestID, teamID, challengeID 获取 Docker
func GetDockerBy3ID(tx *gorm.DB, contestID, teamID uint, challengeID string) (model.Container, bool, string) {
	var docker model.Container
	res := tx.Model(&model.Container{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&docker).Limit(1)
	if res.RowsAffected != 1 {
		return model.Container{}, false, "DockerNotFound"
	}
	return docker, true, "Success"
}

// DeleteDocker 删除 Docker
func DeleteDocker(tx *gorm.DB, docker model.Container) (bool, string) {
	log.Logger.Infof("Stopping container for team %d challenge %s", docker.TeamID, docker.ChallengeID)
	ok, msg := k8s.StopContainer(docker)
	if !ok {
		log.Logger.Warningf("Failed to stop container for challenge %s: %s", docker.ChallengeID, msg)
		return false, "DeleteDockerError"
	}
	if ok, msg = UpdateDocker(tx, docker.ID, map[string]interface{}{"duration": time.Now().Sub(docker.Start)}); !ok {
		return false, msg
	}
	res := tx.Model(&model.Container{}).
		Where("id = ?", docker.ID).Delete(&model.Container{})
	if res.Error != nil {
		return false, "DeleteDockerError"
	}
	return true, "Success"
}

// UpdateDocker 更新 Docker, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateDocker(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	var count int
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed too many times to update user due to optimistic lock")
			return false, "FailedTooManyTimes"
		}
		var docker model.Container
		res := tx.Model(&model.Container{}).Where("id = ?", id).Find(&docker).Limit(1)
		if res.RowsAffected != 1 {
			return false, "DockerNotFound"
		}
		res = tx.Model(&docker).Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Docker: %v", res.Error)
			return false, "UpdateDockerError"
		}
		if res.RowsAffected == 0 {
			log.Logger.Debug("Failed to update docker due to optimistic lock")
			continue
		}
		break
	}
	return true, "Success"
}
