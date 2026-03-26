package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
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

func shuffleTeams(teams []model.Team) model.RetVal {
	for i := len(teams) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			log.Logger.Errorf("Failed to shuffle teams: %s", err)
			return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
		}
		j := int(n.Int64())
		teams[i], teams[j] = teams[j], teams[i]
	}
	return model.SuccessRetVal()
}

func needVPC(pods []model.ChallengePodTemplate) bool {
	for _, pod := range pods {
		for _, network := range pod.Networks {
			if network.CIDR != "" {
				return true
			}
		}
	}
	return false
}

func buildVictimSpec(tx *gorm.DB, victim model.Victim, challenge model.Challenge) (model.VictimSpec, model.RetVal) {
	spec := model.VictimSpec{
		TemplateVersion: challenge.TemplateVersion,
		Pods:            make([]model.PodSpec, 0, len(challenge.Template.Pods)),
		NetworkPolicies: challenge.NetworkPolicies,
		FrpEnabled:      config.Env.K8S.Frp.On,
	}
	flagValues := make(map[uint]string)
	if victim.TeamID.Valid {
		challengeFlagIDs := make([]uint, 0, len(challenge.ChallengeFlags))
		for _, challengeFlag := range challenge.ChallengeFlags {
			challengeFlagIDs = append(challengeFlagIDs, challengeFlag.ID)
		}
		if len(challengeFlagIDs) > 0 {
			teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
				Conditions: map[string]any{
					"team_id":            victim.TeamID.V,
					"challenge_flag_id": challengeFlagIDs,
				},
			})
			if !ret.OK && ret.Msg != i18n.Model.NotFound {
				return model.VictimSpec{}, ret
			}
			for _, teamFlag := range teamFlags {
				flagValues[teamFlag.ChallengeFlagID] = teamFlag.Value
			}
		}
	}
	bindingKey := func(podKey, containerKey string) string {
		return podKey + "\x00" + containerKey
	}
	containerFlags := make(map[string][]model.ChallengeFlag, len(challenge.ChallengeFlags))
	for _, flag := range challenge.ChallengeFlags {
		key := bindingKey(flag.Binding.PodKey, flag.Binding.ContainerKey)
		containerFlags[key] = append(containerFlags[key], flag)
	}

	buildContainerSpec := func(
		podTemplate model.ChallengePodTemplate,
		containerTemplate model.ChallengeContainerTemplate,
	) (model.VictimContainerSpec, model.RetVal) {
		containerSpec := model.VictimContainerSpec{
			Key:         containerTemplate.Key,
			Name:        containerTemplate.Name,
			Image:       containerTemplate.Image,
			CPU:         containerTemplate.CPU,
			Memory:      containerTemplate.Memory,
			WorkingDir:  containerTemplate.WorkingDir,
			Command:     append(model.StringList(nil), containerTemplate.Command...),
			Environment: make(model.StringMap),
			Files:       make(model.StringMap),
			Exposes:     append(model.Exposes(nil), containerTemplate.Exposes...),
		}
		for key, value := range containerTemplate.Environment {
			containerSpec.Environment[key] = value
		}
		for _, flag := range containerFlags[bindingKey(podTemplate.Key, containerTemplate.Key)] {
			value := flag.Value
			if injected, ok := flagValues[flag.ID]; ok {
				value = injected
			}
			switch flag.Binding.Type {
			case model.EnvFlagBindingType:
				containerSpec.Environment[flag.Binding.Target] = value
			case model.FileFlagBindingType:
				containerSpec.Files[flag.Binding.Target] = value
			default:
				return model.VictimContainerSpec{}, model.RetVal{Msg: i18n.Model.ChallengeFlag.InvalidType}
			}
		}
		return containerSpec, model.SuccessRetVal()
	}

	networkPlans := make(map[string]*model.Subnet)
	dnatDedup := make(map[string]struct{})
	snatDedup := make(map[string]struct{})
	dnatPorts := make([]int32, 0)
	hasVPC := needVPC(challenge.Template.Pods)
	if hasVPC {
		spec.NetworkPlan = model.VPC{
			Name:    fmt.Sprintf("vpc-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandStr(6)),
			Subnets: make([]*model.Subnet, 0),
		}
	}
	if !hasVPC {
		containerCount := 0
		for _, podTemplate := range challenge.Template.Pods {
			containerCount += len(podTemplate.Containers)
		}
		mergedPod := model.PodSpec{
			Key:          "default",
			ServicePorts: make(model.Exposes, 0),
			Networks:     make(model.Networks, 0),
			Containers:   make([]model.VictimContainerSpec, 0, containerCount),
		}
		for _, podTemplate := range challenge.Template.Pods {
			for _, expose := range podTemplate.ServicePorts {
				if !slices.ContainsFunc(mergedPod.ServicePorts, func(p model.Expose) bool {
					return p.Port == expose.Port && strings.EqualFold(p.Protocol, expose.Protocol)
				}) {
					mergedPod.ServicePorts = append(mergedPod.ServicePorts, expose)
				}
			}
			for _, containerTemplate := range podTemplate.Containers {
				containerSpec, ret := buildContainerSpec(podTemplate, containerTemplate)
				if !ret.OK {
					return model.VictimSpec{}, ret
				}
				mergedPod.Containers = append(mergedPod.Containers, containerSpec)
			}
		}
		spec.Pods = append(spec.Pods, mergedPod)
		return spec, model.SuccessRetVal()
	}

	for _, podTemplate := range challenge.Template.Pods {
		podSpec := model.PodSpec{
			Key:          podTemplate.Key,
			ServicePorts: append(model.Exposes(nil), podTemplate.ServicePorts...),
			Networks:     append(model.Networks(nil), podTemplate.Networks...),
			Containers:   make([]model.VictimContainerSpec, 0, len(podTemplate.Containers)),
		}
		for _, containerTemplate := range podTemplate.Containers {
			containerSpec, ret := buildContainerSpec(podTemplate, containerTemplate)
			if !ret.OK {
				return model.VictimSpec{}, ret
			}
			podSpec.Containers = append(podSpec.Containers, containerSpec)
		}
		spec.Pods = append(spec.Pods, podSpec)

		for _, network := range podTemplate.Networks {
			if spec.NetworkPlan.Name == "" {
				continue
			}
			subnet, ok := networkPlans[network.Name]
			if !ok {
				subnet = &model.Subnet{
					DefName:      network.Name,
					Name:         fmt.Sprintf("net-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandStr(6)),
					CIDRBlock:    network.CIDR,
					Gateway:      network.Gateway,
					NetAttachDef: &model.NetAttachDef{Name: fmt.Sprintf("nad-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandStr(6))},
				}
				networkPlans[network.Name] = subnet
				spec.NetworkPlan.Subnets = append(spec.NetworkPlan.Subnets, subnet)
			}
			if network.External || len(podTemplate.ServicePorts) > 0 {
				snats := make([]*model.SNat, 0)
				dnats := make([]*model.DNat, 0)
				if network.External {
					if _, exists := snatDedup[network.Name]; !exists {
						snats = append(snats, &model.SNat{Name: fmt.Sprintf("snat-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandStr(6))})
						snatDedup[network.Name] = struct{}{}
					}
				}
				for _, expose := range podTemplate.ServicePorts {
					key := fmt.Sprintf("%s-%s-%d-%s", podTemplate.Key, network.Name, expose.Port, expose.Protocol)
					if _, exists := dnatDedup[key]; exists {
						continue
					}
					dnatDedup[key] = struct{}{}
					dnats = append(dnats, &model.DNat{
						Name: fmt.Sprintf("dnat-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandStr(6)),
						ExternalPort: func() int32 {
							for {
								port, _ := rand.Int(rand.Reader, big.NewInt(65534))
								portInt := int32(port.Int64())
								if !slices.Contains(dnatPorts, portInt) {
									dnatPorts = append(dnatPorts, portInt)
									return portInt
								}
							}
						}(),
						InternalIP:   network.IP,
						InternalPort: expose.Port,
						Protocol:     expose.Protocol,
					})
				}
				if len(snats) > 0 || len(dnats) > 0 {
					lanIP, err := utils.GetLastIP(subnet.CIDRBlock)
					if err != nil {
						return model.VictimSpec{}, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "IP", "Error": err.Error()}}
					}
					if !slices.Contains(subnet.ExcludeIps, lanIP) {
						subnet.ExcludeIps = append(subnet.ExcludeIps, lanIP)
					}
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
	}

	return spec, model.SuccessRetVal()
}

func buildPodRecords(victim model.Victim) []db.CreatePodOptions {
	options := make([]db.CreatePodOptions, 0, len(victim.Spec.Pods))
	for _, podSpec := range victim.Spec.Pods {
		options = append(options, db.CreatePodOptions{
			VictimID: victim.ID,
			Name: fmt.Sprintf("pod-%d-%d-%s-%s", victim.ContestChallengeID.V, victim.UserID, func() string {
				name := strings.ToLower(podSpec.Key)
				if len(name) < 15 {
					return name
				}
				return name[:15]
			}(), utils.RandStr(6)),
			Spec: podSpec,
		})
	}
	return options
}

func StartVictim(tx *gorm.DB, userID, teamID, contestID uint, contestChallengeID, challengeID uint, durationL ...time.Duration) model.RetVal {
	var (
		challengeRepo = db.InitChallengeRepo(tx)
		victimRepo    = db.InitVictimRepo(tx)
		podRepo       = db.InitPodRepo(tx)
	)
	if _, ret := victimRepo.HasAliveVictim(teamID, challengeID); ret.OK {
		return model.SuccessRetVal()
	}
	challenge, ret := challengeRepo.GetByID(challengeID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
	})
	if !ret.OK {
		return ret
	}
	duration := 2 * time.Hour
	if len(durationL) > 0 && durationL[0] > 0 {
		duration = durationL[0]
	}
	baseVictim := model.Victim{
		ChallengeID:        challengeID,
		ContestID:          sql.Null[uint]{V: contestID, Valid: contestID > 0},
		ContestChallengeID: sql.Null[uint]{V: contestChallengeID, Valid: contestChallengeID > 0},
		TeamID:             sql.Null[uint]{V: teamID, Valid: teamID > 0},
		UserID:             userID,
	}
	spec, ret := buildVictimSpec(tx, baseVictim, challenge)
	if !ret.OK {
		return ret
	}
	victim, ret := victimRepo.Create(db.CreateVictimOptions{
		UserID:             userID,
		TeamID:             sql.Null[uint]{V: teamID, Valid: teamID > 0},
		ContestID:          sql.Null[uint]{V: contestID, Valid: contestID > 0},
		ContestChallengeID: sql.Null[uint]{V: contestChallengeID, Valid: contestChallengeID > 0},
		ChallengeID:        challengeID,
		Start:              time.Now(),
		Duration:           duration,
		Spec:               spec,
	})
	if !ret.OK {
		return ret
	}
	for _, options := range buildPodRecords(victim) {
		pod, ret := podRepo.Create(options)
		if !ret.OK {
			return ret
		}
		victim.Pods = append(victim.Pods, pod)
	}
	_, err := task.EnqueueStartVictimTask(victim)
	if err != nil {
		log.Logger.Warningf("Failed to enqueue start victim task: %v", err)
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
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

func StopVictim(tx *gorm.DB, victim model.Victim) model.RetVal {
	switch victim.Status {
	case model.WaitingVictimStatus, model.PendingVictimStatus:
		return model.RetVal{Msg: i18n.Model.Victim.NotStoppable}
	}
	return ForceStopVictim(tx, victim)
}

func ForceStopVictim(tx *gorm.DB, victim model.Victim) model.RetVal {
	if victim.Status == model.TerminatingVictimStatus {
		return model.SuccessRetVal()
	}
	repo := db.InitVictimRepo(tx)
	rollbackStatus := victim.Status
	if ret := repo.Update(victim.ID, db.UpdateVictimOptions{
		Status: new(model.TerminatingVictimStatus),
	}); !ret.OK {
		return ret
	}
	victim.Status = model.TerminatingVictimStatus
	LoadTraffic(tx, victim)
	_, err := task.EnqueueStopVictimTask(victim)
	if err != nil {
		log.Logger.Warningf("Failed to enqueue stop victim task: %v", err)
		_ = repo.Update(victim.ID, db.UpdateVictimOptions{
			Status: &rollbackStatus,
		})
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
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
		Sort:    []string{"id DESC"},
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
	if len(form.Challenges) == 0 || form.TeamRatio <= 0 || form.TeamRatio >= 1 {
		return model.SuccessRetVal()
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.PodsChallengeType, "rand_id": form.Challenges},
	})
	if !ret.OK {
		return ret
	}
	challengeIDL := make([]uint, 0)
	for _, challenge := range challenges {
		challengeIDL = append(challengeIDL, challenge.ID)
	}
	teams, _, ret := db.InitTeamRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if !ret.OK {
		return ret
	}
	if len(challengeIDL) == 0 || len(teams) == 0 {
		return model.SuccessRetVal()
	}
	teamCount := int(float64(len(teams)) * form.TeamRatio)
	if teamCount <= 0 {
		teamCount = 1
	}
	if ret = shuffleTeams(teams); !ret.OK {
		return ret
	}
	teams = teams[:teamCount]
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
	duration := time.Duration(form.Duration) * time.Second
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
			StartVictim(tx, team.CaptainID, team.ID, contest.ID, contestChallenge.ID, contestChallenge.ChallengeID, duration)
		}
	}
	return model.SuccessRetVal()
}

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
			log.Logger.Warningf("Skip stopping victim %d: %s", victim.ID, ret.Msg)
		}
	}
	return model.SuccessRetVal()
}
