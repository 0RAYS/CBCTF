package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GenTestAttachment(tx *gorm.DB, challenge model.Challenge) (bool, string) {
	challengeFlags, _, ok, msg := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ok {
		return false, msg
	}
	ok, msg = k8s.GenTestAttachment(challenge, challengeFlags)
	if !ok {
		return false, msg
	}
	return true, msg
}

func GetTestVictimStatus(tx *gorm.DB, challenge model.Challenge) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Down",
	}
	if challenge.Type != model.PodsChallengeType {
		data["status"] = "NotDocker"
		return data
	}
	victim, ok, _ := db.InitVictimRepo(tx).HasAliveTestVictim(challenge.ID)
	if !ok {
		return data
	}
	data["target"] = victim.RemoteAddr()
	data["status"] = "Running"
	data["remaining"] = victim.Remaining().Seconds()
	return data
}

// StartTestVictim model.Challenge
func StartTestVictim(tx *gorm.DB, challenge model.Challenge) (model.Victim, bool, string) {
	var (
		challengeRepo = db.InitChallengeRepo(tx)
		victimRepo    = db.InitVictimRepo(tx)
		podRepo       = db.InitPodRepo(tx)
		containerRepo = db.InitContainerRepo(tx)
	)
	challenge, ok, msg := challengeRepo.GetByID(challenge.ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Dockers": {Preloads: map[string]db.GetOptions{"ChallengeFlags": {}}}},
	})
	if victim, ok, _ := victimRepo.HasAliveTestVictim(challenge.ID); ok {
		return victim, true, i18n.Success
	}
	vOptions := db.CreateVictimOptions{
		ChallengeID:     challenge.ID,
		Start:           time.Now(),
		Duration:        time.Hour,
		NetworkPolicies: challenge.NetworkPolicies,
	}
	var victim model.Victim
	if needVPC(challenge.Dockers) {
		pOptionsL := make(map[uint]db.CreatePodOptions)
		cOptionsL := make(map[uint]db.CreateContainerOptions)
		vpc := model.VPC{
			Name:    fmt.Sprintf("vpc-%s", utils.RandStr(20)),
			Subnets: make([]*model.Subnet, 0),
		}
		subnets := make(map[string]*model.Subnet)
		// DNat 去重
		networkDockerExposeDNat := make([]string, 0)
		// SNat 去重
		networkExternalSNat := make([]string, 0)
		for _, docker := range challenge.Dockers {
			for _, network := range docker.Networks {
				subnet, ok := subnets[network.Name]
				if !ok {
					subnet = &model.Subnet{
						DefName:   network.Name,
						Name:      fmt.Sprintf("net-%s", utils.RandStr(20)),
						CIDRBlock: network.CIDR,
						Gateway:   network.Gateway,
						//ExcludeIps:   []string{network.Gateway, network.IP},
						NetAttachDef: &model.NetAttachDef{Name: fmt.Sprintf("nad-%s", utils.RandStr(20))},
					}
					vpc.Subnets = append(vpc.Subnets, subnet)
					subnets[network.Name] = subnet
				}
				if network.External || len(docker.Exposes) > 0 {
					eip := &model.EIP{
						Name: fmt.Sprintf("eip-%s", utils.RandStr(20)),
					}
					if network.External {
						if !slices.Contains(networkExternalSNat, network.Name) {
							eip.SNats = []*model.SNat{{Name: fmt.Sprintf("snat-%s", utils.RandStr(20))}}
							networkExternalSNat = append(networkExternalSNat, network.Name)
						}
					}
					for _, expose := range docker.Exposes {
						key := fmt.Sprintf("%s-%d-%s", docker.Name, expose.Port, expose.Protocol)
						if !slices.Contains(networkDockerExposeDNat, key) {
							eip.DNats = append(eip.DNats, &model.DNat{
								Name:         fmt.Sprintf("dnat-%s", utils.RandStr(20)),
								ExternalPort: strconv.Itoa(int(expose.Port)),
								InternalIP:   network.IP,
								InternalPort: strconv.Itoa(int(expose.Port)),
								Protocol:     expose.Protocol,
							})
							networkDockerExposeDNat = append(networkDockerExposeDNat, key)
						}
					}
					if len(eip.SNats) > 0 || len(eip.DNats) > 0 {
						lanIP, err := utils.GetLastIP(subnet.CIDRBlock)
						if err != nil {
							return model.Victim{}, false, i18n.GetIPError
						}
						subnet.NatGateway = &model.NatGateway{
							Name:  fmt.Sprintf("nat-%s", utils.RandStr(20)),
							LanIP: lanIP,
						}
						subnet.NatGateway.EIPs = append(subnet.NatGateway.EIPs, eip)
					}
				}
			}
			pOptionsL[docker.ID] = db.CreatePodOptions{
				Name:     fmt.Sprintf("pod-%s", utils.RandStr(20)),
				PodPorts: docker.Exposes,
				Networks: docker.Networks,
			}
			envFlagL := make(model.StringMap)
			volumeFlagL := make(model.StringMap)
			for _, challengeFlag := range docker.ChallengeFlags {
				switch challengeFlag.InjectType {
				case model.EnvInjectType:
					envFlagL[challengeFlag.Name] = challengeFlag.Value
				case model.VolumeInjectType:
					volumeFlagL[challengeFlag.Path] = challengeFlag.Value
				default:
					return model.Victim{}, false, i18n.InvalidChallengeFlagInjectType
				}
			}
			cOptionsL[docker.ID] = db.CreateContainerOptions{
				Name:        fmt.Sprintf("ctn-%s", utils.RandStr(20)),
				Image:       docker.Image,
				CPU:         docker.CPU,
				Memory:      docker.Memory,
				WorkingDir:  docker.WorkingDir,
				Command:     docker.Command,
				Environment: docker.Environment,
				EnvFlags:    envFlagL,
				VolumeFlags: volumeFlagL,
				Exposes:     docker.Exposes,
			}
		}
		vOptions.VPC = vpc
		victim, ok, msg = victimRepo.Create(vOptions)
		if !ok {
			return model.Victim{}, false, msg
		}
		for _, docker := range challenge.Dockers {
			pOptions, ok := pOptionsL[docker.ID]
			if !ok {
				return model.Victim{}, false, i18n.UnknownError
			}
			pOptions.VictimID = victim.ID
			pod, ok, msg := podRepo.Create(pOptions)
			if !ok {
				return model.Victim{}, false, msg
			}
			cOptions, ok := cOptionsL[docker.ID]
			if !ok {
				return model.Victim{}, false, i18n.UnknownError
			}
			cOptions.PodID = pod.ID
			container, ok, msg := containerRepo.Create(cOptions)
			if !ok {
				return model.Victim{}, false, msg
			}
			pod.Containers = append(pod.Containers, container)
			victim.Pods = append(victim.Pods, pod)
		}
	} else {
		victim, ok, msg = victimRepo.Create(vOptions)
		if !ok {
			return model.Victim{}, false, msg
		}
		pOptions := db.CreatePodOptions{
			VictimID: victim.ID,
			Name:     fmt.Sprintf("pod-%s", utils.RandStr(20)),
			PodPorts: make(model.Exposes, 0),
		}
		cOptionsL := make([]db.CreateContainerOptions, 0)
		tmp := make([]string, 0)
		for _, docker := range challenge.Dockers {
			envFlagL := make(model.StringMap)
			volumeFlagL := make(model.StringMap)
			for _, challengeFlag := range docker.ChallengeFlags {
				switch challengeFlag.InjectType {
				case model.EnvInjectType:
					envFlagL[challengeFlag.Name] = challengeFlag.Value
				case model.VolumeInjectType:
					volumeFlagL[challengeFlag.Path] = challengeFlag.Value
				default:
					return model.Victim{}, false, i18n.InvalidChallengeFlagInjectType
				}
			}
			cOptionsL = append(cOptionsL, db.CreateContainerOptions{
				Name:        fmt.Sprintf("ctn-%s", utils.RandStr(20)),
				Image:       docker.Image,
				CPU:         docker.CPU,
				Memory:      docker.Memory,
				WorkingDir:  docker.WorkingDir,
				Command:     docker.Command,
				Environment: docker.Environment,
				EnvFlags:    envFlagL,
				VolumeFlags: volumeFlagL,
				Exposes:     docker.Exposes,
			})
			for _, p := range docker.Exposes {
				t := fmt.Sprintf("%d/%s", p.Port, p.Protocol)
				if !slices.Contains(tmp, t) {
					pOptions.PodPorts = append(pOptions.PodPorts, p)
					tmp = append(tmp, t)
				}
			}
		}
		pod, ok, msg := podRepo.Create(pOptions)
		if !ok {
			return model.Victim{}, false, msg
		}
		for _, cOptions := range cOptionsL {
			cOptions.PodID = pod.ID
			container, ok, msg := containerRepo.Create(cOptions)
			if !ok {
				return model.Victim{}, false, msg
			}
			pod.Containers = append(pod.Containers, container)
		}
		victim.Pods = append(victim.Pods, pod)
	}
	ipExposesMap, ok, msg := k8s.StartVictim(victim)
	if !ok {
		return model.Victim{}, false, msg
	}
	for ip, exposes := range ipExposesMap {
		for _, expose := range exposes {
			victim.Endpoints = append(victim.Endpoints, model.Endpoint{
				IP:       ip,
				Port:     expose.Port,
				Protocol: expose.Protocol,
			})
		}
	}
	victim.ExposedEndpoints = victim.Endpoints
	if config.Env.K8S.Frpc.On {
		var frpc string
		victim.ExposedEndpoints, frpc, ok, msg = k8s.CreateFrpc(victim)
		if !ok {
			return model.Victim{}, false, msg
		}
		p, ok, msg := podRepo.Create(db.CreatePodOptions{
			VictimID: victim.ID,
			Name:     frpc,
		})
		if !ok {
			return model.Victim{}, false, msg
		}
		victim.Pods = append(victim.Pods, p)
	}
	if ok, msg = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
		VPC:              &victim.VPC,
		Endpoints:        &victim.Endpoints,
		ExposedEndpoints: &victim.ExposedEndpoints,
	}); !ok {
		return model.Victim{}, false, msg
	}
	return victim, true, i18n.Success
}

func StopTestVictim(tx *gorm.DB, challenge model.Challenge) (bool, string) {
	victimRepo := db.InitVictimRepo(tx)
	victim, ok, msg := victimRepo.HasAliveTestVictim(challenge.ID)
	if !ok {
		return false, msg
	}
	ok, msg = k8s.StopVictim(victim)
	if !ok {
		return false, msg
	}
	duration := time.Now().Sub(victim.Start)
	if ok, msg = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
		Duration: &duration,
	}); !ok {
		return false, msg
	}
	LoadTraffic(tx, victim)
	return victimRepo.Delete(victim.ID)
}
