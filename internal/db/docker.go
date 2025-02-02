package db

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"time"
)

func CreateDocker(ctx context.Context, flag model.Flag, creatorID uint) (model.Docker, bool, string) {
	var (
		docker model.Docker
		ok     bool
		port   int
	)
	if docker, ok, _ = GetDockerBy3ID(ctx, flag.ContestID, flag.TeamID, flag.ChallengeID); ok {
		return docker, ok, "Success"
	}
	challenge, ok, msg := GetChallengeByID(ctx, flag.ChallengeID)
	if !ok || challenge.Type != model.Container {
		return model.Docker{}, false, msg
	}

	docker = model.InitDocker(flag, challenge, creatorID)
	res := DB.WithContext(ctx).Model(model.Docker{}).Create(&docker)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Docker: %s", res.Error)
		return model.Docker{}, false, "CreateDockerError"
	}
	log.Logger.Debugf("Starting container for team %d challenge %s", flag.TeamID, flag.ChallengeID)
	port, ok, msg = k8s.StartContainer(challenge, flag, docker)
	if !ok {
		log.Logger.Warningf("Failed to start container for challenge %s: %s", flag.ChallengeID, msg)
		_, _ = DeleteDocker(ctx, docker.ID)
		return model.Docker{}, false, msg
	}
	UpdateDocker(ctx, docker.ID, map[string]interface{}{"port": port})
	docker.Port = int32(port)
	return docker, true, "Success"
}

func GetDockerByID(ctx context.Context, id uint) (model.Docker, bool, string) {
	var docker model.Docker
	res := DB.WithContext(ctx).Model(model.Docker{}).Where("id = ?", id).Find(&docker).Limit(1)
	if res.RowsAffected != 1 {
		return model.Docker{}, false, "DockerNotFound"
	}
	return docker, true, "Success"
}

func GetDockerBy3ID(ctx context.Context, contestID, teamID uint, challengeID string) (model.Docker, bool, string) {
	var docker model.Docker
	res := DB.WithContext(ctx).Model(model.Docker{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&docker).Limit(1)
	if res.RowsAffected != 1 {
		return model.Docker{}, false, "DockerNotFound"
	}
	if docker.Start.Add(docker.Duration).Before(time.Now()) {
		_, _ = DeleteDocker(ctx, docker.ID)
		return model.Docker{}, false, "DockerNotFound"
	}
	return docker, true, "Success"
}

func DeleteDocker(ctx context.Context, id uint) (bool, string) {
	docker, ok, msg := GetDockerByID(ctx, id)
	if !ok {
		return false, msg
	}
	go func(d model.Docker) {
		log.Logger.Debugf("Stopping container for team %d challenge %s", d.TeamID, d.ChallengeID)
		ok, msg = k8s.StopContainer(d)
		if !ok {
			log.Logger.Warningf("Failed to stop container for challenge %s: %s", d.ChallengeID, msg)
		}
	}(docker)
	res := DB.WithContext(ctx).Model(model.Docker{}).
		Where("id = ?", id).Delete(&model.Docker{})
	if res.Error != nil {
		return false, "DeleteDockerError"
	}
	return true, "Success"
}

func UpdateDocker(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	res := DB.WithContext(ctx).Model(model.Docker{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Docker: %v", res.Error)
		return false, "UpdateDockerError"
	}
	return true, "Success"
}
