package service

import (
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// GetRemoteStatus model.Usage 需要预加载
func GetRemoteStatus(tx *gorm.DB, usage model.Usage) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Running",
	}
	if usage.Challenge.Type != model.DockerChallenge && usage.Challenge.Type != model.DockersChallenge {
		data["status"] = "NotDocker"
		return data
	}
	repo := db.InitContainerRepo(tx)
	var minTime float64
	for _, container := range usage.Containers {
		_, ok, _ := repo.GetByID(container.ID, false, 0)
		if !ok {
			data["status"] = "Down"
			continue
		}
		if minTime == 0 || minTime > container.Remaining().Seconds() {
			minTime = container.Remaining().Seconds()
		}
		data["target"] = append(data["target"].([]string), container.RemoteAddr()...)
	}
	if len(data["target"].([]string)) > 0 && data["status"] == "Down" {
		data["status"] = "PartDown"
	}
	data["remaining"] = minTime
	return data
}

func StartContainer(tx *gorm.DB, user model.User, team model.Team, usage model.Usage) (bool, string) {
	answerRepo := db.InitAnswerRepo(tx)
	containerRepo := db.InitContainerRepo(tx)
	containers, ok, _ := containerRepo.GetBy2ID(team.ID, usage.ID, false, 0)
	if ok {
		return true, "Success"
	}
	switch usage.Challenge.Type {
	case model.DockerChallenge:
		options := db.CreateContainerOptions{
			UsageID:           usage.ID,
			TeamID:            team.ID,
			UserID:            user.ID,
			Exposes:           usage.Docker.Ports,
			Start:             time.Now(),
			Duration:          time.Hour,
			Image:             usage.Docker.Image,
			PodName:           fmt.Sprintf("victim-%s-%d-pod", usage.ChallengeID, team.ID),
			ContainerName:     fmt.Sprintf("victim-%s-%d", usage.ChallengeID, team.ID),
			ServiceName:       fmt.Sprintf("victim-%s-%d-svc", usage.ChallengeID, team.ID),
			NetworkPolicyName: fmt.Sprintf("victim-%s-%d-net", usage.ChallengeID, team.ID),
			NetworkPolicies:   usage.Docker.NetworkPolicies,
		}
		for _, flagID := range usage.Docker.FlagsID {
			answer, ok, msg := answerRepo.GetBy2ID(team.ID, flagID, false, 0)
			if !ok {
				return false, msg
			}
			options.Flags = append(options.Flags, answer.Value)
		}
		container, ok, msg := containerRepo.Create(options)
		if !ok {
			return false, msg
		}
		containers = append(containers, container)
	case model.DockersChallenge:
		for i, docker := range usage.Dockers {
			options := db.CreateContainerOptions{
				UsageID:           usage.ID,
				TeamID:            team.ID,
				UserID:            user.ID,
				Exposes:           docker.Ports,
				Start:             time.Now(),
				Duration:          time.Hour,
				Image:             docker.Image,
				PodName:           fmt.Sprintf("victim-%s-%d-pod-%d", usage.ChallengeID, team.ID, i),
				ContainerName:     fmt.Sprintf("victim-%s-%d-%d", usage.ChallengeID, team.ID, i),
				ServiceName:       fmt.Sprintf("victim-%s-%d-svc-%d", usage.ChallengeID, team.ID, i),
				NetworkPolicyName: fmt.Sprintf("victim-%s-%d-net-%d", usage.ChallengeID, team.ID, i),
				NetworkPolicies:   docker.NetworkPolicies,
			}
			for _, flagID := range docker.FlagsID {
				answer, ok, msg := answerRepo.GetBy2ID(team.ID, flagID, false, 0)
				if !ok {
					return false, msg
				}
				options.Flags = append(options.Flags, answer.Value)
			}
			container, ok, msg := containerRepo.Create(options)
			if !ok {
				return false, msg
			}
			containers = append(containers, container)
		}
	default:
		return false, "InvalidChallengeType"
	}
	for _, container := range containers {
		ip, ok, msg := k8s.StartContainer(container)
		if !ok {
			return false, msg
		}
		ok, msg = containerRepo.Update(container.ID, db.UpdateContainerOptions{
			IP: &ip,
		})
		if !ok {
			return false, msg
		}
	}
	return true, "Success"
}

func StopContainer(tx *gorm.DB, container model.Container) (bool, string) {
	ok, msg := k8s.StopContainer(container)
	if !ok {
		return false, "DeleteContainerError"
	}
	repo := db.InitContainerRepo(tx)
	duration := time.Now().Sub(container.Start)
	ok, msg = repo.Update(container.ID, db.UpdateContainerOptions{
		Duration: &duration,
	})
	if !ok {
		return false, msg
	}
	ok, msg = repo.Delete(container.ID)
	if !ok {
		return false, msg
	}
	return true, "Success"
}
