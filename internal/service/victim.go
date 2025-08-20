package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func needVPC(dockers []model.Docker) bool {
	for _, docker := range dockers {
		for _, network := range docker.Networks {
			if network.CIDR != "" {
				return true
			}
		}
	}
	return false
}

func StartTeamVictim(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge) (model.Victim, bool, string) {
	var (
		challengeRepo = db.InitChallengeRepo(tx)
		victimRepo    = db.InitVictimRepo(tx)
		teamFlagRepo  = db.InitTeamFlagRepo(tx)
		podRepo       = db.InitPodRepo(tx)
		containerRepo = db.InitContainerRepo(tx)
	)
	challenge, ok, msg := challengeRepo.GetByID(contestChallenge.ChallengeID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Dockers": {Preloads: map[string]db.GetOptions{"ChallengeFlags": {}}}},
	})
	if !ok {
		return model.Victim{}, false, msg
	}
	if victim, ok, _ := victimRepo.HasAliveVictim(team.ID, contestChallenge.ID); ok {
		return victim, true, i18n.Success
	}
	vOptions := db.CreateVictimOptions{
		ContestChallengeID: contestChallenge.ID,
		TeamID:             team.ID,
		UserID:             user.ID,
		Start:              time.Now(),
		Duration:           time.Hour,
		NetworkPolicies:    challenge.NetworkPolicies,
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
				teamFlag, ok, msg := teamFlagRepo.Get(db.GetOptions{
					Selects:    []string{"id", "challenge_flag_id", "value"},
					Conditions: map[string]any{"team_id": team.ID, "challenge_flag_id": challengeFlag.ID},
					Preloads:   map[string]db.GetOptions{"ChallengeFlag": {Selects: []string{"id", "Name"}}},
				})
				if !ok {
					return model.Victim{}, false, msg
				}
				switch challengeFlag.InjectType {
				case model.EnvInjectType:
					envFlagL[teamFlag.ChallengeFlag.Name] = teamFlag.Value
				case model.VolumeInjectType:
					volumeFlagL[challengeFlag.Path] = teamFlag.Value
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
				teamFlag, ok, msg := teamFlagRepo.Get(db.GetOptions{
					Selects:    []string{"id", "challenge_flag_id", "value"},
					Conditions: map[string]any{"team_id": team.ID, "challenge_flag_id": challengeFlag.ID},
					Preloads:   map[string]db.GetOptions{"ChallengeFlag": {Selects: []string{"id", "Name"}}},
				})
				if !ok {
					return model.Victim{}, false, msg
				}
				switch challengeFlag.InjectType {
				case model.EnvInjectType:
					envFlagL[teamFlag.ChallengeFlag.Name] = teamFlag.Value
				case model.VolumeInjectType:
					volumeFlagL[challengeFlag.Path] = teamFlag.Value
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
	ipExposesMap, ok, msg := k8s.StartVictim(victim, victim.Pods)
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
		var frpc []string
		victim.ExposedEndpoints, frpc, ok, msg = k8s.CreateFrpc(victim)
		if !ok {
			return model.Victim{}, false, msg
		}
		for _, frpcPodName := range frpc {
			p, ok, msg := podRepo.Create(db.CreatePodOptions{
				VictimID: victim.ID,
				Name:     frpcPodName,
			})
			if !ok {
				return model.Victim{}, false, msg
			}
			victim.Pods = append(victim.Pods, p)
		}
	}
	if ok, msg = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
		VPC:              &victim.VPC,
		Endpoints:        &victim.Endpoints,
		ExposedEndpoints: &victim.ExposedEndpoints,
	}); !ok {
		return model.Victim{}, false, msg
	}
	prometheus.AddVictimContainerMetrics(contest, contestChallenge, 1)
	return victim, true, i18n.Success
}

// GetTeamVictimStatus contestChallenge 需要预加载 model.ContestChallenge
func GetTeamVictimStatus(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Down",
	}
	if contestChallenge.Type != model.PodsChallengeType {
		data["status"] = "NotDocker"
		return data
	}
	victims, _, ok, _ := db.InitVictimRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID, "team_id": team.ID},
		Preloads: map[string]db.GetOptions{
			"Pods": {},
		},
	})
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
	data["target"] = victims[0].RemoteAddr()
	data["status"] = "Running"
	data["remaining"] = victims[0].Remaining().Seconds()
	return data
}

func StopTeamVictim(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) (bool, string) {
	victimRepo := db.InitVictimRepo(tx)
	victims, _, ok, msg := victimRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID, "contest_challenge_id": contestChallenge.ID},
		Preloads: map[string]db.GetOptions{
			"Pods": {},
		},
	})
	if !ok {
		return false, msg
	}
	// 预期中, len(victims) == 1, 考虑意外情况
	victimIDL := make([]uint, 0)
	for _, victim := range victims {
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
		victimIDL = append(victimIDL, victim.ID)
		LoadTraffic(tx, victim)
	}
	ok, msg = victimRepo.Delete(victimIDL...)
	if ok {
		prometheus.SubVictimContainerMetrics(team.Contest, contestChallenge, 1)
	}
	return ok, msg
}
