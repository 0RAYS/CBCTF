package db

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

// CreateDocker 创建 Docker, 并注入 flag
func CreateDocker(tx *gorm.DB, flag model.Flag, challenge model.Challenge, creatorID uint) (model.Docker, bool, string) {
	var (
		docker model.Docker
		ok     bool
		msg    string
		ip     string
		port   int32
	)
	if docker, ok, _ = GetDockerBy3ID(tx, flag.ContestID, flag.TeamID, flag.ChallengeID); ok {
		return docker, ok, "Success"
	}
	if challenge.Type != model.Container {
		return model.Docker{}, false, "InvalidChallengeType"
	}
	docker = model.InitDocker(flag, challenge, creatorID)
	log.Logger.Debugf("Starting container for team %d challenge %s", flag.TeamID, flag.ChallengeID)
	ip, port, ok, msg = k8s.StartContainer(challenge, flag, docker)
	if !ok {
		log.Logger.Warningf("Failed to start container for challenge %s: %s", flag.ChallengeID, msg)
		return model.Docker{}, false, msg
	}
	docker.IP = ip
	docker.Port = port
	res := tx.Model(model.Docker{}).Create(&docker)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Docker: %s", res.Error)
		return model.Docker{}, false, "CreateDockerError"
	}
	return docker, true, "Success"
}

// GetDockers 获取所有 Docker
func GetDockers(tx *gorm.DB, deleted bool) ([]model.Docker, bool, string) {
	var dockers []model.Docker
	res := tx.Model(model.Docker{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Find(&dockers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Dockers: %s", res.Error)
		return nil, false, "GetDockersError"
	}
	return dockers, true, "Success"
}

// GetDockerByID 根据 ID 获取 Docker
func GetDockerByID(tx *gorm.DB, id uint, deleted bool) (model.Docker, bool, string) {
	var docker model.Docker
	res := tx.Model(model.Docker{}).Where("id = ?", id)
	if deleted {
		res = res.Unscoped()
	}
	res = res.Find(&docker).Limit(1)
	if res.RowsAffected != 1 {
		return model.Docker{}, false, "DockerNotFound"
	}
	return docker, true, "Success"
}

func GetDockerByTeamID(tx *gorm.DB, teamID uint, limit, offset int, deleted bool) ([]model.Docker, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	res := tx.Model(model.Docker{}).Where("team_id = ?", teamID)
	if deleted {
		res = res.Unscoped()
	}
	var dockers []model.Docker
	var count int64
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count Dockers: %s", err)
		return nil, -1, false, "UnknownError"
	}
	res = res.Limit(limit).Offset(offset).Find(&dockers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Dockers: %s", res.Error)
		return nil, -1, false, "GetDockersError"
	}
	return dockers, count, true, "Success"
}

// GetDockerBy3ID 根据 contestID, teamID, challengeID 获取 Docker
func GetDockerBy3ID(tx *gorm.DB, contestID, teamID uint, challengeID string) (model.Docker, bool, string) {
	var docker model.Docker
	res := tx.Model(model.Docker{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&docker).Limit(1)
	if res.RowsAffected != 1 {
		return model.Docker{}, false, "DockerNotFound"
	}
	return docker, true, "Success"
}

// DeleteDocker 删除 Docker
func DeleteDocker(tx *gorm.DB, docker model.Docker) (bool, string) {
	log.Logger.Infof("Stopping container for team %d challenge %s", docker.TeamID, docker.ChallengeID)
	ok, msg := k8s.StopContainer(docker)
	if !ok {
		log.Logger.Warningf("Failed to stop container for challenge %s: %s", docker.ChallengeID, msg)
		return false, "DeleteDockerError"
	}
	res := tx.Model(model.Docker{}).
		Where("id = ?", docker.ID).Delete(&model.Docker{})
	if res.Error != nil {
		return false, "DeleteDockerError"
	}
	return true, "Success"
}

// UpdateDocker 更新 Docker, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateDocker(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	res := tx.Model(model.Docker{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Docker: %v", res.Error)
		return false, "UpdateDockerError"
	}
	return true, "Success"
}
