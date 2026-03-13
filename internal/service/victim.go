package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"slices"
	"strings"
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

func StartVictim(tx *gorm.DB, userID, teamID, contestID uint, contestChallengeID, challengeID uint) model.RetVal {
	var (
		challengeRepo = db.InitChallengeRepo(tx)
		victimRepo    = db.InitVictimRepo(tx)
		teamFlagRepo  = db.InitTeamFlagRepo(tx)
		podRepo       = db.InitPodRepo(tx)
		containerRepo = db.InitContainerRepo(tx)
	)
	if _, ret := victimRepo.HasAliveVictim(teamID, challengeID); ret.OK {
		return model.SuccessRetVal()
	}
	challenge, ret := challengeRepo.GetByID(challengeID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Dockers": {Preloads: map[string]db.GetOptions{"ChallengeFlags": {}}}},
	})
	if !ret.OK {
		return ret
	}
	vOptions := db.CreateVictimOptions{
		UserID:             userID,
		TeamID:             sql.Null[uint]{V: teamID, Valid: teamID > 0},
		ContestID:          sql.Null[uint]{V: contestID, Valid: contestID > 0},
		ContestChallengeID: sql.Null[uint]{V: contestChallengeID, Valid: contestChallengeID > 0},
		ChallengeID:        challengeID,
		Start:              time.Now(),
		Duration:           time.Hour,
		NetworkPolicies:    challenge.NetworkPolicies,
	}
	var victim model.Victim
	if needVPC(challenge.Dockers) {
		pOptionsL := make(map[uint]db.CreatePodOptions)
		cOptionsL := make(map[uint]db.CreateContainerOptions)
		vpc := model.VPC{
			Name:    fmt.Sprintf("vpc-%d-%d-%s", contestChallengeID, userID, utils.RandStr(6)),
			Subnets: make([]*model.Subnet, 0),
		}
		subnets := make(map[string]*model.Subnet)
		// DNat 去重
		networkDockerExposeDNat := make([]string, 0)
		// SNat 去重
		networkExternalSNat := make([]string, 0)
		// port 去重
		dnatPort := make([]int32, 0)
		for _, docker := range challenge.Dockers {
			for _, network := range docker.Networks {
				subnet, ok := subnets[network.Name]
				if !ok {
					subnet = &model.Subnet{
						DefName:      network.Name,
						Name:         fmt.Sprintf("net-%d-%d-%s", contestChallengeID, userID, utils.RandStr(6)),
						CIDRBlock:    network.CIDR,
						Gateway:      network.Gateway,
						NetAttachDef: &model.NetAttachDef{Name: fmt.Sprintf("nad-%d-%d-%s", contestChallengeID, userID, utils.RandStr(6))},
					}
					vpc.Subnets = append(vpc.Subnets, subnet)
					subnets[network.Name] = subnet
				}
				if network.External || len(docker.Exposes) > 0 {
					snats := make([]*model.SNat, 0)
					dnats := make([]*model.DNat, 0)
					if network.External {
						if !slices.Contains(networkExternalSNat, network.Name) {
							snats = []*model.SNat{{Name: fmt.Sprintf("snat-%d-%d-%s", contestChallengeID, userID, utils.RandStr(6))}}
							networkExternalSNat = append(networkExternalSNat, network.Name)
						}
					}
					for _, expose := range docker.Exposes {
						key := fmt.Sprintf("%s-%d-%s", docker.Name, expose.Port, expose.Protocol)
						if !slices.Contains(networkDockerExposeDNat, key) {
							dnats = append(dnats, &model.DNat{
								Name: fmt.Sprintf("dnat-%d-%d-%s", contestChallengeID, userID, utils.RandStr(6)),
								ExternalPort: func() int32 {
									for {
										port, _ := rand.Int(rand.Reader, big.NewInt(65534))
										if !slices.Contains(dnatPort, int32(port.Int64())) {
											dnatPort = append(dnatPort, int32(port.Int64()))
											return int32(port.Int64())
										}
									}
								}(),
								InternalIP:   network.IP,
								InternalPort: expose.Port,
								Protocol:     expose.Protocol,
							})
							networkDockerExposeDNat = append(networkDockerExposeDNat, key)
						}
					}
					if len(snats) > 0 || len(dnats) > 0 {
						lanIP, err := utils.GetLastIP(subnet.CIDRBlock)
						if err != nil {
							return model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "IP", "Error": err.Error()}}
						}
						subnet.ExcludeIps = append(subnet.ExcludeIps, lanIP)
						if subnet.NatGateway == nil {
							subnet.NatGateway = &model.NatGateway{
								Name:  fmt.Sprintf("nat-%s", utils.RandStr(20)),
								LanIP: lanIP,
								EIP: &model.EIP{
									Name: fmt.Sprintf("eip-%s", utils.RandStr(20)),
								},
							}
						}
						subnet.NatGateway.EIP.DNats = append(subnet.NatGateway.EIP.DNats, dnats...)
						subnet.NatGateway.EIP.SNats = append(subnet.NatGateway.EIP.SNats, snats...)
					}
				}
			}
			pOptionsL[docker.ID] = db.CreatePodOptions{
				Name: fmt.Sprintf("pod-%d-%d-%s-%s", contestChallengeID, userID, func() string {
					name := strings.ToLower(docker.Name)
					if len(name) < 15 {
						return name
					}
					return name[:15]
				}(), utils.RandStr(6)),
				PodPorts: docker.Exposes,
				Networks: docker.Networks,
			}
			envFlagL := make(model.StringMap)
			volumeFlagL := make(model.StringMap)
			for _, challengeFlag := range docker.ChallengeFlags {
				value := challengeFlag.Value
				// teamID == 0 时为测试靶机
				if teamID > 0 {
					teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
						Conditions: map[string]any{"team_id": teamID, "challenge_flag_id": challengeFlag.ID},
					})
					if !ret.OK {
						return ret
					}
					value = teamFlag.Value
				}
				switch challengeFlag.InjectType {
				case model.EnvFlagInjectType:
					envFlagL[challengeFlag.Name] = value
				case model.VolumeFlagInjectType:
					volumeFlagL[challengeFlag.Path] = value
				default:
					return model.RetVal{Msg: i18n.Model.ChallengeFlag.InvalidType}
				}
			}
			cOptionsL[docker.ID] = db.CreateContainerOptions{
				Name:        docker.Name,
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
		victim, ret = victimRepo.Create(vOptions)
		if !ret.OK {
			return ret
		}
		for _, docker := range challenge.Dockers {
			pOptions, ok := pOptionsL[docker.ID]
			if !ok {
				return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Unknown docker.ID"}}
			}
			pOptions.VictimID = victim.ID
			pod, ret := podRepo.Create(pOptions)
			if !ret.OK {
				return ret
			}
			cOptions, ok := cOptionsL[docker.ID]
			if !ok {
				return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Unknown docker.ID"}}
			}
			cOptions.PodID = pod.ID
			container, ret := containerRepo.Create(cOptions)
			if !ret.OK {
				return ret
			}
			pod.Containers = append(pod.Containers, container)
			victim.Pods = append(victim.Pods, pod)
		}
	} else {
		victim, ret = victimRepo.Create(vOptions)
		if !ret.OK {
			return ret
		}
		pOptions := db.CreatePodOptions{
			VictimID: victim.ID,
			Name:     fmt.Sprintf("pod-%d-%d-%s", contestChallengeID, userID, utils.RandStr(6)),
			PodPorts: make(model.Exposes, 0),
		}
		cOptionsL := make([]db.CreateContainerOptions, 0)
		tmp := make([]string, 0)
		for _, docker := range challenge.Dockers {
			envFlagL := make(model.StringMap)
			volumeFlagL := make(model.StringMap)
			for _, challengeFlag := range docker.ChallengeFlags {
				value := challengeFlag.Value
				// team.ID == 0 时为测试靶机
				if teamID > 0 {
					teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
						Conditions: map[string]any{"team_id": teamID, "challenge_flag_id": challengeFlag.ID},
					})
					if !ret.OK {
						return ret
					}
					value = teamFlag.Value
				}
				switch challengeFlag.InjectType {
				case model.EnvFlagInjectType:
					envFlagL[challengeFlag.Name] = value
				case model.VolumeFlagInjectType:
					volumeFlagL[challengeFlag.Path] = value
				default:
					return model.RetVal{Msg: i18n.Model.ChallengeFlag.InvalidType}
				}
			}
			cOptionsL = append(cOptionsL, db.CreateContainerOptions{
				Name:        docker.Name,
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
		pod, ret := podRepo.Create(pOptions)
		if !ret.OK {
			return ret
		}
		for _, cOptions := range cOptionsL {
			cOptions.PodID = pod.ID
			container, ret := containerRepo.Create(cOptions)
			if !ret.OK {
				return ret
			}
			pod.Containers = append(pod.Containers, container)
		}
		victim.Pods = append(victim.Pods, pod)
	}
	if ret = victimRepo.Update(victim.ID, db.UpdateVictimOptions{Status: new(model.PendingVictimStatus)}); !ret.OK {
		return ret
	}
	go func() {
		ret = func() model.RetVal {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			ipExposesMap, ret := k8s.StartVictim(ctx, victim)
			if !ret.OK {
				return ret
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
			if config.Env.K8S.Frp.On {
				var frpc []string
				victim.ExposedEndpoints, frpc, ret = k8s.CreateFrpc(ctx, victim)
				if !ret.OK {
					return ret
				}
				for _, frpcPodName := range frpc {
					p, ret := podRepo.Create(db.CreatePodOptions{
						VictimID: victim.ID,
						Name:     frpcPodName,
					})
					if !ret.OK {
						return ret
					}
					victim.Pods = append(victim.Pods, p)
				}
			}
			if ret = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
				VPC:              &victim.VPC,
				Endpoints:        &victim.Endpoints,
				ExposedEndpoints: &victim.ExposedEndpoints,
				Start:            new(time.Now()),
				Status:           new(model.RunningVictimStatus),
			}); !ret.OK {
				return ret
			}
			return model.SuccessRetVal()
		}()
		if !ret.OK {
			victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
			if !ret.OK {
				return
			}
			StopVictim(db.DB, victim)
		}
	}()
	return model.SuccessRetVal()
}

func GetVictimStatus(tx *gorm.DB, teamID uint, challenge model.Challenge) gin.H {
	data := gin.H{
		"target":    make([]string, 0),
		"remaining": 0,
		"status":    "Down",
	}
	if challenge.Type != model.PodsChallengeType {
		data["status"] = "not_docker"
		return data
	}
	victim, ret := db.InitVictimRepo(tx).HasAliveVictim(teamID, challenge.ID)
	if !ret.OK {
		return data
	}
	targets := victim.RemoteAddr()
	data["target"] = targets
	data["status"] = victim.Status
	data["remaining"] = victim.Remaining().Seconds()
	return data
}

// StopVictim tx 无需开启事务
func StopVictim(tx *gorm.DB, victim model.Victim) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ret := k8s.StopVictim(ctx, victim)
	if !ret.OK {
		return ret
	}
	tx2 := tx.Begin()
	if ret = db.InitVictimRepo(tx2).Update(victim.ID, db.UpdateVictimOptions{
		Duration: new(time.Now().Sub(victim.Start)),
	}); !ret.OK {
		tx2.Rollback()
		return ret
	}
	LoadTraffic(tx2, victim)
	ret = db.InitVictimRepo(tx2).Delete(victim.ID)
	if !ret.OK {
		tx2.Rollback()
		return ret
	}
	tx2.Commit()
	return ret
}

func CountTeamVictims(tx *gorm.DB, team model.Team) (int64, model.RetVal) {
	return db.InitVictimRepo(tx).Count(db.CountOptions{Conditions: map[string]any{"team_id": team.ID}})
}

func GetVictims(tx *gorm.DB, contest model.Contest, form dto.GetVictimsForm) ([]model.Victim, int64, int64, model.RetVal) {
	var challengeID uint
	if form.ChallengeID != "" {
		challenge, ret := db.InitChallengeRepo(tx).GetByRandID(form.ChallengeID)
		if !ret.OK || challenge.Type != model.PodsChallengeType {
			return nil, 0, 0, ret
		}
		challengeID = challenge.ID
	}
	options := db.GetOptions{
		Conditions: make(map[string]any),
		Preloads: map[string]db.GetOptions{
			"Pods":             {},
			"User":             {},
			"Team":             {},
			"ContestChallenge": {},
		},
		Deleted: form.Deleted,
	}
	if contest.ID != 0 {
		options.Conditions["contest_id"] = contest.ID
	}
	if challengeID != 0 {
		options.Conditions["challenge_id"] = challengeID
	}
	if form.TeamID != 0 {
		options.Conditions["team_id"] = form.TeamID
	}
	if form.UserID != 0 {
		options.Conditions["user_id"] = form.UserID
	}
	victims, count, ret := db.InitVictimRepo(tx).List(form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, 0, ret
	}
	countOptions := db.CountOptions{Deleted: true}
	if contest.ID != 0 {
		countOptions.Conditions = map[string]any{"contest_id": contest.ID}
	}
	total, ret := db.InitVictimRepo(db.DB).Count(countOptions)
	if !ret.OK {
		total = count
	}
	return victims, count, total, ret
}

func StartVictims(tx *gorm.DB, contest model.Contest, form dto.StartVictimsForm) model.RetVal {
	if len(form.Challenges) == 0 || len(form.Teams) == 0 {
		return model.SuccessRetVal()
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.PodsChallengeType},
	})
	if !ret.OK {
		return ret
	}
	challengeIDL := make([]uint, 0)
	for _, challenge := range challenges {
		challengeIDL = append(challengeIDL, challenge.ID)
	}
	teams, _, ret := db.InitTeamRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "id": form.Teams},
	})
	if !ret.OK {
		return ret
	}
	if len(challengeIDL) == 0 || len(teams) == 0 {
		return model.SuccessRetVal()
	}
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "challenge_id": challengeIDL},
		Preloads:   map[string]db.GetOptions{"ContestFlags": {}},
	})
	if !ret.OK {
		return ret
	}
	if len(contestChallenges) == 0 {
		return model.SuccessRetVal()
	}
	for _, contestChallenge := range contestChallenges {
		for _, team := range teams {
			if CheckIfSolved(tx, team, contestChallenge.ContestFlags) {
				continue
			}
			if !CheckIfGenerated(tx, team, contestChallenge.ContestFlags) {
				if _, ret = CreateTeamFlag(tx, team, contest, contestChallenge); !ret.OK {
					continue
				}
			}
			StartVictim(tx, team.CaptainID, team.ID, contest.ID, contestChallenge.ID, contestChallenge.ChallengeID)
		}
	}
	return model.SuccessRetVal()
}

// StopVictims tx 无需开启事务
func StopVictims(tx *gorm.DB, form dto.StopVictimsForm) model.RetVal {
	if len(form.Victims) == 0 {
		return model.SuccessRetVal()
	}
	victims, _, ret := db.InitVictimRepo(tx).List(-1, -1, db.GetOptions{Conditions: map[string]any{"id": form.Victims}})
	if !ret.OK {
		return ret
	}
	for _, victim := range victims {
		if ret = StopVictim(tx, victim); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
