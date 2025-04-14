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
	"time"
)

// StartVictim model.Usage 需要预加载 model.Challenge
func StartVictim(tx *gorm.DB, user model.User, team model.Team, usage model.Usage) (model.Victim, bool, string) {
	var (
		victim model.Victim
		ok     bool
		msg    string
		err    error
	)
	victimRepo := db.InitVictimRepo(tx)
	victims, ok, _ := victimRepo.GetBy2ID(team.ID, usage.ID, false)
	if ok {
		return victims[0], true, "Success"
	}
	n := team.ID
	block, err := utils.GetIPBlock(n, config.Env.K8S.IPPool.CIDR, config.Env.K8S.IPPool.BlockSize)
	if err != nil {
		return model.Victim{}, false, "GetIPBlockError"
	}
	if len(block) == 0 {
		return model.Victim{}, false, "EmptyIPBlock"
	}
	ipBlock := fmt.Sprintf("%s-%d", block[0], len(block))
	dns := make(map[string]string)
	podRepo := db.InitPodRepo(tx)
	containerRepo := db.InitContainerRepo(tx)
	answerRepo := db.InitAnswerRepo(tx)
	if usage.Challenge.Type == model.PodsChallenge {
		if len(block) < len(usage.Dockers) {
			return model.Victim{}, false, "EmptyIPBlock"
		}
		vOptions := db.CreateVictimOptions{
			UsageID:  usage.ID,
			TeamID:   team.ID,
			UserID:   user.ID,
			IPBlock:  ipBlock,
			Start:    time.Now(),
			Duration: time.Hour,
		}
		victim, ok, msg = victimRepo.Create(vOptions)
		if !ok {
			return model.Victim{}, false, msg
		}
		podGroup := make(map[uint][]model.Docker)
		for _, docker := range usage.Dockers {
			podGroup[docker.PodGroup] = append(podGroup[docker.PodGroup], docker)
		}
		for i, dockers := range podGroup {
			for _, docker := range dockers {
				dns[docker.Hostname] = block[i]
			}
		}
		for i, dockers := range podGroup {
			exposes := make([]int32, 0)
			policies := make([]model.NetworkPolicy, 0)
			for _, docker := range dockers {
				exposes = append(exposes, docker.Ports...)
				for x, policy := range docker.NetworkPolicies {
					for y, target := range policy.To {
						if target.Hostname != "" {
							ip, ok := dns[target.Hostname]
							if !ok {
								return model.Victim{}, false, "InvalidNetworkPolicy"
							}
							dockers[i].NetworkPolicies[x].To[y].CIDR = fmt.Sprintf("%s/32", ip)
							dockers[i].NetworkPolicies[x].To[y].Except = nil
						}
					}
					for y, target := range policy.From {
						if target.Hostname != "" {
							ip, ok := dns[target.Hostname]
							if !ok {
								return model.Victim{}, false, "InvalidNetworkPolicy"
							}
							dockers[i].NetworkPolicies[x].From[y].CIDR = fmt.Sprintf("%s/32", ip)
							dockers[i].NetworkPolicies[x].From[y].Except = nil
						}
					}
				}
				policies = append(policies, docker.NetworkPolicies...)
			}

			pOptions := db.CreatePodOptions{
				VictimID:          victim.ID,
				Name:              fmt.Sprintf("victim-%s-%d-pod-%d", usage.ChallengeID, team.ID, i),
				PodIP:             block[i],
				ServiceName:       fmt.Sprintf("victim-%s-%d-svc-%d", usage.ChallengeID, team.ID, i),
				NetworkPolicyName: fmt.Sprintf("victim-%s-%d-net-%d", usage.ChallengeID, team.ID, i),
				ExposePorts:       exposes,
				NetworkPolicies:   policies,
			}
			pod, ok, msg := podRepo.Create(pOptions)
			if !ok {
				return model.Victim{}, false, msg
			}
			for _, docker := range dockers {
				cOptions := db.CreateContainerOptions{
					PodID:    pod.ID,
					Name:     fmt.Sprintf("victim-%s-%d-%d", usage.ChallengeID, team.ID, i),
					Image:    docker.Image,
					Hostname: docker.Hostname,
					Exposes:  docker.Ports,
				}
				for _, flagID := range docker.FlagIDL {
					answer, ok, msg := answerRepo.GetBy2ID(team.ID, flagID)
					if !ok {
						return model.Victim{}, false, msg
					}
					cOptions.Flags = append(cOptions.Flags, answer.Value)
				}
				container, ok, msg := containerRepo.Create(cOptions)
				if !ok {
					return model.Victim{}, false, msg
				}
				pod.Containers = append(pod.Containers, container)
			}
			victim.Pods = append(victim.Pods, pod)
		}
	} else {
		return model.Victim{}, false, "InvalidChallengeType"
	}
	ipL, ok, msg := k8s.StartVictim(victim, dns)
	if !ok {
		return model.Victim{}, false, msg
	}
	for i, pod := range victim.Pods {
		ip := ipL[pod.Name]
		ok, msg := podRepo.Update(pod.ID, db.UpdatePodOptions{
			ExposeIP: &ip,
		})
		if !ok {
			return model.Victim{}, false, msg
		}
		victim.Pods[i].ExposeIP = ip
	}
	return victim, true, "Success"
}

// GetVictimStatus model.Usage 需要预加载 model.Challenge
func GetVictimStatus(tx *gorm.DB, team model.Team, usage model.Usage) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Down",
	}
	if usage.Challenge.Type != model.PodsChallenge {
		data["status"] = "NotDocker"
		return data
	}
	repo := db.InitVictimRepo(tx)
	victims, ok, _ := repo.GetBy2ID(team.ID, usage.ID, false, "Pods")
	if !ok {
		return data
	}
	if len(victims) == 0 {
		return data
	}
	if len(victims) > 1 {
		data["status"] = "Error"
		return data
	}
	for _, pod := range victims[0].Pods {
		data["target"] = append(data["target"].([]string), pod.RemoteAddr()...)
	}
	data["status"] = "Running"
	data["remaining"] = victims[0].Remaining().Seconds()
	return data
}

func StopVictim(tx *gorm.DB, team model.Team, usage model.Usage) (bool, string) {
	victimRepo := db.InitVictimRepo(tx)
	victims, ok, msg := victimRepo.GetBy2ID(team.ID, usage.ID, false, "Pods")
	if !ok {
		return false, msg
	}
	// 预期中, len(victims) == 1, 考虑意外情况
	for _, victim := range victims {
		ok, msg = k8s.StopVictim(victim)
		if !ok {
			return false, msg
		}
		utils.RemoveIPBlock(victim.IPBlock)
		duration := time.Now().Sub(victim.Start)
		if ok, msg = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
			Duration: &duration,
		}); !ok {
			return false, msg
		}
		if ok, msg = victimRepo.Delete(victim.ID); !ok {
			return false, msg
		}
		go LoadTraffic(tx, victim)
	}
	return true, "Success"
}
