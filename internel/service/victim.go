package service

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/k8s"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"slices"
	"strings"
	"time"
)

func NeedVPC(dockers []model.Docker) bool {
	for _, docker := range dockers {
		for _, network := range docker.Networks {
			if network.CIDR != "" {
				return true
			}
		}
	}
	return false
}

// StartTeamVictim Todo
func StartTeamVictim(tx *gorm.DB, user model.User, team model.Team, contestChallenge model.ContestChallenge) (model.Victim, bool, string) {
	challenge, ok, msg := db.InitChallengeRepo(tx).
		GetByID(contestChallenge.ChallengeID, db.GetOptions{
			Preloads: map[string]db.GetOptions{
				"Dockers": {
					Preloads: map[string]db.GetOptions{
						"ChallengeFlags": {},
					},
				},
			},
		})
	if !ok {
		return model.Victim{}, false, msg
	}
	victimRepo := db.InitVictimRepo(tx)
	if victim, ok, _ := victimRepo.HasAliveVictim(team.ID, contestChallenge.ID); ok {
		return victim, true, i18n.Success
	}
	podRepo := db.InitPodRepo(tx)
	containerRepo := db.InitContainerRepo(tx)
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	vOptions := db.CreateVictimOptions{
		ContestChallengeID: contestChallenge.ID,
		TeamID:             team.ID,
		UserID:             user.ID,
		Start:              time.Now(),
		Duration:           time.Hour,
		NetworkPolicies:    challenge.NetworkPolicies,
	}
	var victim model.Victim
	if NeedVPC(challenge.Dockers) {
		vOptions.VPC = fmt.Sprintf("vpc-%s", utils.RandStr(10))

		subnetsName := make([]string, 0)
		gatewayName := make([]string, 0)
		subnets := make(model.Subnets, 0)
		netAttachDefs := make(model.StringMap)
		gateways := make(model.Gateways, 0)
		eips := make(model.EIPs, 0)
		dnats := make(model.DNats, 0)
		snats := make(model.SNats, 0)
		pOptionsMap := make(map[uint]db.CreatePodOptions)
		cOptionsMap := make(map[uint]db.CreateContainerOptions)
		for _, docker := range challenge.Dockers {
			pOptions := db.CreatePodOptions{
				Name:     fmt.Sprintf("pod-%s", utils.RandStr(10)),
				PodPorts: docker.Exposes,
				//IPs:      make(model.IPs, 0),
			}
			envFlagL := make(model.StringList, 0)
			volumeFlagL := make(model.StringMap)
			for _, challengeFlag := range docker.ChallengeFlags {
				teamFlag, ok, msg := teamFlagRepo.Get(db.GetOptions{
					Conditions: map[string]any{
						"team_id":           team.ID,
						"challenge_flag_id": challengeFlag.ID,
					},
				})
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
				Name:        fmt.Sprintf("ctn-%s", utils.RandStr(10)),
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
			cOptionsMap[docker.ID] = cOptions
			for _, network := range docker.Networks {
				name := strings.ReplaceAll(strings.ReplaceAll(network.CIDR, ".", "_"), "/", "_")

				subnet := model.Subnet{
					Name:     fmt.Sprintf("net-%s", utils.RandStr(10)),
					CIDR:     network.CIDR,
					Gateway:  network.Gateway,
					External: network.External,
				}
				netAttachDefs[subnet.Name] = fmt.Sprintf("nad-%s", utils.RandStr(5))
				if !slices.Contains(subnetsName, name) {
					subnets = append(subnets, subnet)
				}

				lanIP, err := utils.GetLastIP(subnet.CIDR)
				if err != nil {
					log.Logger.Warningf("Failed to get lan IP for subnet %s: %v", subnet.Name, err)
					return model.Victim{}, false, msg
				}
				gateway := model.Gateway{
					Name:   fmt.Sprintf("gw-%s", utils.RandStr(10)),
					VPC:    vOptions.VPC,
					Subnet: subnet.Name,
					LanIP:  lanIP,
				}
				if !slices.Contains(gatewayName, name) {
					gateways = append(gateways, gateway)
				}

				eip := model.EIP{
					Name:    fmt.Sprintf("eip-%s", utils.RandStr(10)),
					Gateway: gateway.Name,
				}
				eips = append(eips, eip)

				snat := model.SNat{
					Name:         fmt.Sprintf("snat-%s", utils.RandStr(10)),
					EIP:          eip.Name,
					InternalCIDR: subnet.CIDR,
				}
				snats = append(snats, snat)

				for _, e := range docker.Exposes {
					dnats = append(dnats, model.DNat{
						Name:         fmt.Sprintf("dnat-%s", utils.RandStr(10)),
						EIP:          eip.Name,
						ExternalPort: e.Port,
						InternalIP:   network.IP,
						InternalPort: e.Port,
						Protocol:     e.Protocol,
					})
				}
				//pOptions.IPs = append(pOptions.IPs, model.IP{
				//	Name:    fmt.Sprintf("%s.%s.%s.%s.ovn", pOptions.Name, k8s.GlobalNamespace, subnet.Name, k8s.GlobalNamespace),
				//	Subnet:  subnet.Name,
				//	PodName: pOptions.Name,
				//	IP:      network.IP,
				//})
			}
			pOptionsMap[docker.ID] = pOptions
		}
		vOptions.Subnets = subnets
		vOptions.NetAttachDefs = netAttachDefs
		vOptions.Gateways = gateways
		vOptions.EIPs = eips
		vOptions.SNats = snats

		victim, ok, msg = victimRepo.Create(vOptions)
		if !ok {
			return model.Victim{}, false, msg
		}
		for id, pOptions := range pOptionsMap {
			pOptions.VictimID = victim.ID
			pod, ok, msg := podRepo.Create(pOptions)
			if !ok {
				return model.Victim{}, false, msg
			}
			cOptions, ok := cOptionsMap[id]
			if !ok {
				log.Logger.Warningf("Failed to create container options for pod %s.", pod.Name)
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
			Name:     fmt.Sprintf("pod-%s", utils.RandStr(10)),
			PodPorts: make(model.Exposes, 0),
		}
		cOptionsL := make([]db.CreateContainerOptions, 0)
		tmp := make([]string, 0)
		for _, docker := range challenge.Dockers {
			envFlagL := make(model.StringList, 0)
			volumeFlagL := make(model.StringMap)
			for _, challengeFlag := range docker.ChallengeFlags {
				teamFlag, ok, msg := teamFlagRepo.Get(db.GetOptions{
					Conditions: map[string]any{
						"team_id":           team.ID,
						"challenge_flag_id": challengeFlag.ID,
					},
				})
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
			cOptionsL = append(cOptionsL, db.CreateContainerOptions{
				Name:        fmt.Sprintf("ctn-%s", utils.RandStr(10)),
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
	ok, msg = k8s.StartVictim(victim)
	if !ok {
		//go k8s.StopVictim(victim)
		return model.Victim{}, false, msg
	}
	//for i, pod := range victim.Pods {
	//	target := targets[pod.Name]
	//	ip := target["ip"].(string)
	//	ports := model.Ports(target["ports"].([]int32))
	//	ok, msg = podRepo.Update(pod.ID, db.UpdatePodOptions{
	//		ExposedIP:    &ip,
	//		ExposedPorts: &ports,
	//	})
	//	if !ok {
	//		return model.Victim{}, false, msg
	//	}
	//	victim.Pods[i].ExposedIP = ip
	//	victim.Pods[i].ExposedPorts = ports
	//}
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
		Conditions: map[string]any{
			"contest_challenge_id": contestChallenge.ID,
			"team_id":              team.ID,
		},
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
	for _, pod := range victims[0].Pods {
		data["target"] = append(data["target"].([]string), pod.RemoteAddr()...)
	}
	data["status"] = "Running"
	data["remaining"] = victims[0].Remaining().Seconds()
	return data
}

func StopTeamVictim(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) (bool, string) {
	victimRepo := db.InitVictimRepo(tx)
	victims, _, ok, msg := victimRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{
			"team_id":              team.ID,
			"contest_challenge_id": contestChallenge.ID,
		},
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
	return victimRepo.Delete(victimIDL...)
}
