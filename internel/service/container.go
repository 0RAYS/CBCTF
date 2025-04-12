package service

import (
	"CBCTF/internel/config"
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"sync"
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
		_, ok, _ := repo.GetByID(container.ID)
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

func StartContainer(tx *gorm.DB, user model.User, team model.Team, usage model.Usage) ([]model.Container, bool, string) {
	answerRepo := db.InitAnswerRepo(tx)
	containerRepo := db.InitContainerRepo(tx)
	containers, ok, _ := containerRepo.GetBy2ID(team.ID, usage.ID, false)
	if ok {
		return containers, true, "Success"
	}
	n := team.ID
	block, err := utils.GetIPBlock(n, config.Env.K8S.IPPool.CIDR, config.Env.K8S.IPPool.BlockSize)
	if err != nil {
		return containers, false, "GetIPBlockError"
	}
	if len(block) == 0 {
		return containers, false, "EmptyIPBlock"
	}
	ipBlock := fmt.Sprintf("%s-%d", block[0], len(block))
	dns := make(map[string]string)
	switch usage.Challenge.Type {
	case model.DockerChallenge:
		dns[usage.Docker.Hostname] = block[rand.Intn(len(block))]
		options := db.CreateContainerOptions{
			UsageID:           usage.ID,
			TeamID:            team.ID,
			UserID:            user.ID,
			PodIP:             dns[usage.Docker.Hostname],
			IPBlock:           ipBlock,
			Exposes:           usage.Docker.Ports,
			Start:             time.Now(),
			Duration:          time.Hour,
			Image:             usage.Docker.Image,
			Hostname:          usage.Docker.Hostname,
			PodName:           fmt.Sprintf("victim-%s-%d-pod", usage.ChallengeID, team.ID),
			ContainerName:     fmt.Sprintf("victim-%s-%d", usage.ChallengeID, team.ID),
			ServiceName:       fmt.Sprintf("victim-%s-%d-svc", usage.ChallengeID, team.ID),
			NetworkPolicyName: fmt.Sprintf("victim-%s-%d-net", usage.ChallengeID, team.ID),
			NetworkPolicies:   usage.Docker.NetworkPolicies,
		}
		for _, flagID := range usage.Docker.FlagIDL {
			answer, ok, msg := answerRepo.GetBy2ID(team.ID, flagID)
			if !ok {
				return containers, false, msg
			}
			options.Flags = append(options.Flags, answer.Value)
		}
		container, ok, msg := containerRepo.Create(options)
		if !ok {
			return containers, false, msg
		}
		containers = append(containers, container)
	case model.DockersChallenge:
		if len(block) < len(usage.Dockers) {
			return containers, false, "EmptyIPBlock"
		}
		for i, docker := range usage.Dockers {
			dns[docker.Hostname] = block[i]
		}
		for i, docker := range usage.Dockers {
			for x, policy := range docker.NetworkPolicies {
				for y, target := range policy.To {
					if target.Hostname != "" {
						ip, ok := dns[target.Hostname]
						if !ok {
							return containers, false, "InvalidNetworkPolicy"
						}
						docker.NetworkPolicies[x].To[y].CIDR = fmt.Sprintf("%s/32", ip)
						docker.NetworkPolicies[x].To[y].Except = nil
					}
				}
				for y, target := range policy.From {
					if target.Hostname != "" {
						ip, ok := dns[target.Hostname]
						if !ok {
							return containers, false, "InvalidNetworkPolicy"
						}
						docker.NetworkPolicies[x].From[y].CIDR = fmt.Sprintf("%s/32", ip)
						docker.NetworkPolicies[x].From[y].Except = nil
					}
				}
			}
			options := db.CreateContainerOptions{
				UsageID:           usage.ID,
				TeamID:            team.ID,
				UserID:            user.ID,
				PodIP:             dns[docker.Hostname],
				IPBlock:           ipBlock,
				Exposes:           docker.Ports,
				Start:             time.Now(),
				Duration:          time.Hour,
				Image:             docker.Image,
				Hostname:          docker.Hostname,
				PodName:           fmt.Sprintf("victim-%s-%d-pod-%d", usage.ChallengeID, team.ID, i),
				ContainerName:     fmt.Sprintf("victim-%s-%d-%d", usage.ChallengeID, team.ID, i),
				ServiceName:       fmt.Sprintf("victim-%s-%d-svc-%d", usage.ChallengeID, team.ID, i),
				NetworkPolicyName: fmt.Sprintf("victim-%s-%d-net-%d", usage.ChallengeID, team.ID, i),
				NetworkPolicies:   docker.NetworkPolicies,
			}
			for _, flagID := range docker.FlagIDL {
				answer, ok, msg := answerRepo.GetBy2ID(team.ID, flagID)
				if !ok {
					return containers, false, msg
				}
				options.Flags = append(options.Flags, answer.Value)
			}
			container, ok, msg := containerRepo.Create(options)
			if !ok {
				return containers, false, msg
			}
			containers = append(containers, container)
		}
	default:
		return containers, false, "InvalidChallengeType"
	}
	type result struct {
		C   model.Container
		OK  bool
		Msg string
	}
	var wg sync.WaitGroup
	resultCh := make(chan result, len(containers))
	for _, container := range containers {
		wg.Add(1)
		go func(container model.Container) {
			defer wg.Done()
			pod, ip, ok, msg := k8s.StartContainer(container, dns)
			if !ok {
				resultCh <- result{C: container, OK: false, Msg: msg}
				return
			}
			ok, msg = containerRepo.Update(container.ID, db.UpdateContainerOptions{
				IP:    &ip,
				PodIP: &pod.Status.PodIP,
			})
			if !ok {
				resultCh <- result{C: container, OK: false, Msg: msg}
				return
			}
			container.IP = ip
			container.PodIP = pod.Status.PodIP
			resultCh <- result{C: container, OK: true, Msg: "Success"}
		}(container)
	}
	wg.Wait()
	close(resultCh)
	containers = make([]model.Container, 0)
	for res := range resultCh {
		if !res.OK {
			return containers, false, res.Msg
		}
		containers = append(containers, res.C)
	}
	return containers, true, "Success"
}

func StopContainer(tx *gorm.DB, team model.Team, usage model.Usage) (bool, string) {
	repo := db.InitContainerRepo(tx)
	containers, ok, msg := repo.GetBy2ID(team.ID, usage.ID, false)
	if !ok {
		return false, msg
	}
	type result struct {
		OK  bool
		Msg string
	}
	var wg sync.WaitGroup
	resultCh := make(chan result, len(containers))
	for _, container := range containers {
		wg.Add(1)
		go func(container model.Container) {
			defer wg.Done()
			ok, msg = k8s.StopContainer(container)
			if !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			utils.RemoveIPBlock(container.IPBlock)
			duration := time.Now().Sub(container.Start)
			ok, msg = repo.Update(container.ID, db.UpdateContainerOptions{
				Duration: &duration,
			})
			if !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			ok, msg = repo.Delete(container.ID)
			if !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			resultCh <- result{OK: true, Msg: "Success"}
		}(container)
	}
	wg.Wait()
	close(resultCh)
	for res := range resultCh {
		if !res.OK {
			return false, res.Msg
		}
	}
	return true, "Success"
}
