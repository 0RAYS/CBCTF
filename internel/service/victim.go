package service

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func StartVictim(tx *gorm.DB, user model.User, team model.Team, contestChallenge model.ContestChallenge) (model.Victim, bool, string) {
	challenge, ok, msg := db.InitChallengeRepo(tx).
		GetByID(contestChallenge.ChallengeID, "DockerGroups", "DockerGroups.Dockers", "DockerGroups.Dockers.ChallengeFlags")
	if !ok {
		return model.Victim{}, false, msg
	}
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	victimRepo := db.InitVictimRepo(tx)
	podRepo := db.InitPodRepo(tx)
	containerRepo := db.InitContainerRepo(tx)
	if victim, ok, _ := victimRepo.HasAliveVictim(team.ID, contestChallenge.ID); ok {
		return victim, true, i18n.Success
	}
	ipBlock, err := utils.GetIPBlock(team.ID, config.Env.K8S.IPPool.CIDR, config.Env.K8S.IPPool.BlockSize)
	if err != nil {
		return model.Victim{}, false, i18n.GetIPBlockError
	}
	if len(ipBlock) == 0 || len(challenge.DockerGroups) > len(ipBlock) {
		return model.Victim{}, false, i18n.EmptyIPBlock
	}
	dns := make(model.StringMap)
	podPorts := make(map[uint]model.Ports)
	for i, dockerGroup := range challenge.DockerGroups {
		for _, docker := range dockerGroup.Dockers {
			if _, ok = dns[docker.Name]; ok {
				return model.Victim{}, false, i18n.DuplicateHostname
			}
			dns[docker.Name] = ipBlock[i]
			for _, port := range docker.Expose {
				p, _ := strconv.ParseInt(port, 10, 32)
				if !utils.In(int32(p), podPorts) {
					podPorts[dockerGroup.ID] = append(podPorts[dockerGroup.ID], int32(p))
				}
			}
		}
	}
	vOptions := db.CreateVictimOptions{
		ContestChallengeID: contestChallenge.ID,
		TeamID:             team.ID,
		UserID:             user.ID,
		IPBlock:            fmt.Sprintf("%s-%d", ipBlock[0], len(ipBlock)),
		Start:              time.Now(),
		Duration:           time.Hour,
		HostAlias:          dns,
	}
	victim, ok, msg := victimRepo.Create(vOptions)
	if !ok {
		return model.Victim{}, false, msg
	}
	for i, dockerGroup := range challenge.DockerGroups {
		pOptions := db.CreatePodOptions{
			VictimID:        victim.ID,
			Name:            victim.GenPodName(challenge.RandID),
			PodIP:           ipBlock[i],
			PodPorts:        podPorts[dockerGroup.ID],
			NetworkPolicies: dockerGroup.NetworkPolicies,
		}
		pod, ok, msg := podRepo.Create(pOptions)
		if !ok {
			return model.Victim{}, false, msg
		}
		for _, docker := range dockerGroup.Dockers {
			envFlagL := make(model.StringList, 0)
			volumeFlagL := make(model.StringMap)
			for _, challengeFlag := range docker.ChallengeFlags {
				teamFlag, ok, msg := teamFlagRepo.GetWithConditions(db.GetOptions{
					{Key: "team_id", Value: team.ID, Op: "and"},
					{Key: "challenge_flag_id", Value: challengeFlag.ID, Op: "and"},
				}, false)
				if !ok {
					return model.Victim{}, false, msg
				}
				switch challengeFlag.InjectType {
				case model.EnvInjectType:
					envFlagL = append(envFlagL, teamFlag.Value)
				case model.VolumeInjectType:
					volumeFlagL[challengeFlag.Path] = teamFlag.Value
				default:
					return model.Victim{}, false, i18n.InvalidChallengeFlagInjectType
				}
			}
			cOptions := db.CreateContainerOptions{
				PodID:       pod.ID,
				Name:        fmt.Sprintf("%s-%s", pod.Name, strings.ToLower(utils.RandStr(5))),
				Image:       docker.Image,
				Hostname:    docker.Name,
				WorkingDir:  docker.WorkingDir,
				Command:     docker.Command,
				Environment: docker.Environment,
				EnvFlags:    envFlagL,
				VolumeFlags: volumeFlagL,
				Exposes:     docker.Expose,
			}
			container, ok, msg := containerRepo.Create(cOptions)
			if !ok {
				return model.Victim{}, false, msg
			}
			pod.Containers = append(pod.Containers, container)
		}
		victim.Pods = append(victim.Pods, pod)
	}
	targets, ok, msg := k8s.StartVictim(victim)
	if !ok {
		go k8s.StopVictim(victim)
		return model.Victim{}, false, msg
	}
	for i, pod := range victim.Pods {
		target := targets[pod.Name]
		ip := target["ip"].(string)
		ports := model.Ports(target["ports"].([]int32))
		ok, msg = podRepo.Update(pod.ID, db.UpdatePodOptions{
			ExposedIP:    &ip,
			ExposedPorts: &ports,
		})
		if !ok {
			return model.Victim{}, false, msg
		}
		victim.Pods[i].ExposedIP = ip
		victim.Pods[i].ExposedPorts = ports
	}
	return victim, true, i18n.Success
}

// GetVictimStatus usage 需要预加载 model.ContestChallenge
func GetVictimStatus(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Down",
	}
	if contestChallenge.Challenge.Type != model.PodsChallengeType {
		data["status"] = "NotDocker"
		return data
	}
	victims, _, ok, _ := db.InitVictimRepo(tx).ListWithConditions(-1, -1, db.GetOptions{
		{Key: "team_id", Value: team.ID, Op: "and"},
		{Key: "contest_challenge_id", Value: contestChallenge.ID, Op: "and"},
	}, false, "Pods")
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

func StopVictim(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) (bool, string) {
	victimRepo := db.InitVictimRepo(tx)
	victims, _, ok, msg := victimRepo.ListWithConditions(-1, -1, db.GetOptions{
		{Key: "team_id", Value: team.ID, Op: "and"},
		{Key: "contest_challenge_id", Value: contestChallenge.ID, Op: "and"},
	}, false, "Pods")
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
	}
	return true, i18n.Success
}
