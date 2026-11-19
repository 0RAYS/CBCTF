package service

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
)

func shuffleTeams(teams []model.Team) model.RetVal {
	for i := len(teams) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			log.Logger.Errorf("Failed to shuffle teams: %s", err)
			return model.FailRetVal(i18n.Common.UnknownError, "Team", "", err)
		}
		j := int(n.Int64())
		teams[i], teams[j] = teams[j], teams[i]
	}
	return model.SuccessRetVal()
}

func needVPC(pods []model.ChallengePodTemplate) bool {
	for _, pod := range pods {
		for _, network := range pod.Networks {
			if network.Definition.CIDR != "" {
				return true
			}
		}
	}
	return false
}

// victimFlagSpec 靶机规格中的一条 flag。
// Key 用于关联队伍手里的答案：比赛题目副本用 contest_flag_id，题库测试模式用 challenge_flag_id。
type victimFlagSpec struct {
	Key       uint
	Value     string
	Binding   model.FlagBinding
	TeamValue string
	HasTeam   bool
}

// victimSpecSource 靶机规格来源，比赛内一律使用比赛题目副本，题库题目只在管理员测试模式下使用。
type victimSpecSource struct {
	Template         model.ChallengeTemplate
	NetworkPolicies  model.NetworkPolicies
	Flags            []victimFlagSpec
	ChallengeFlagIDs []uint
	ContestFlagIDs   []uint
}

// buildVictimSpecSourceContestChallenge 依据比赛题目副本构造规格来源。
func buildVictimSpecSourceContestChallenge(tx *gorm.DB, victim model.Victim, contestChallenge model.ContestChallenge) (victimSpecSource, model.RetVal) {
	source := victimSpecSource{
		Template:        contestChallenge.Template,
		NetworkPolicies: contestChallenge.NetworkPolicies,
	}
	for _, flag := range contestChallenge.ContestFlags {
		if !flag.Active {
			continue
		}
		source.Flags = append(source.Flags, victimFlagSpec{
			Key:     flag.ID,
			Value:   flag.Value,
			Binding: flag.Binding,
		})
		source.ContestFlagIDs = append(source.ContestFlagIDs, flag.ID)
	}
	if !victim.TeamID.Valid || len(source.ContestFlagIDs) == 0 {
		return source, model.SuccessRetVal()
	}
	teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{
			"team_id":         victim.TeamID.V,
			"contest_flag_id": source.ContestFlagIDs,
		},
	})
	if !ret.OK && ret.Msg != i18n.Model.NotFound {
		return victimSpecSource{}, ret
	}
	values := make(map[uint]string, len(teamFlags))
	for _, teamFlag := range teamFlags {
		values[teamFlag.ContestFlagID] = teamFlag.Value
	}
	for i := range source.Flags {
		if value, ok := values[source.Flags[i].Key]; ok {
			source.Flags[i].TeamValue = value
			source.Flags[i].HasTeam = true
		}
	}
	return source, model.SuccessRetVal()
}

// buildVictimSpecSourceChallenge 依据题库题目构造规格来源，仅用于管理员测试模式。
func buildVictimSpecSourceChallenge(challenge model.Challenge) victimSpecSource {
	source := victimSpecSource{
		Template:        challenge.Template,
		NetworkPolicies: challenge.NetworkPolicies,
	}
	for _, flag := range challenge.ChallengeFlags {
		source.Flags = append(source.Flags, victimFlagSpec{
			Key:     flag.ID,
			Value:   flag.Value,
			Binding: flag.Binding,
		})
		source.ChallengeFlagIDs = append(source.ChallengeFlagIDs, flag.ID)
	}
	return source
}

func buildVictimSpec(tx *gorm.DB, victim model.Victim, source victimSpecSource) (model.VictimSpec, model.RetVal) {
	spec := model.VictimSpec{
		Pods:            make([]model.PodSpec, 0, len(source.Template.Pods)),
		NetworkPolicies: source.NetworkPolicies,
		FrpEnabled:      config.Env.K8S.Frp.On,
	}
	type containerKey struct {
		pod, container string
	}
	containerFlags := make(map[containerKey][]victimFlagSpec, len(source.Flags))
	for _, flag := range source.Flags {
		key := containerKey{flag.Binding.PodKey, flag.Binding.ContainerKey}
		containerFlags[key] = append(containerFlags[key], flag)
	}
	testMode := !victim.TeamID.Valid
	if !needVPC(source.Template.Pods) {
		containerCount := 0
		for _, podTemplate := range source.Template.Pods {
			containerCount += len(podTemplate.Containers)
		}
		mergedPod := model.PodSpec{
			Key:          "default",
			ServicePorts: make(model.Exposes, 0),
			Networks:     make(model.Networks, 0),
			Containers:   make([]model.VictimContainerSpec, 0, containerCount),
		}
		for _, podTemplate := range source.Template.Pods {
			for _, expose := range podTemplate.ServicePorts {
				if !slices.ContainsFunc(mergedPod.ServicePorts, func(p model.Expose) bool {
					return p.Port == expose.Port && strings.EqualFold(p.Protocol, expose.Protocol)
				}) {
					mergedPod.ServicePorts = append(mergedPod.ServicePorts, expose)
				}
			}
			for _, containerTemplate := range podTemplate.Containers {
				flags := containerFlags[containerKey{podTemplate.Key, containerTemplate.Key}]
				containerSpec, ret := buildVictimContainerSpec(containerTemplate, flags, testMode)
				if !ret.OK {
					return model.VictimSpec{}, ret
				}
				mergedPod.Containers = append(mergedPod.Containers, containerSpec)
			}
		}
		spec.Pods = append(spec.Pods, mergedPod)
		return spec, model.SuccessRetVal()
	}

	spec.NetworkPlan = model.VPC{
		Name:    fmt.Sprintf("vpc-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandHexStr(6)),
		Subnets: make([]*model.Subnet, 0),
	}
	networkPlans := make(map[string]*model.Subnet)
	for _, podTemplate := range source.Template.Pods {
		podSpec := model.PodSpec{
			Key:          podTemplate.Key,
			ServicePorts: append(model.Exposes(nil), podTemplate.ServicePorts...),
			Networks:     append(model.Networks(nil), podTemplate.Networks...),
			Containers:   make([]model.VictimContainerSpec, 0, len(podTemplate.Containers)),
		}
		for _, containerTemplate := range podTemplate.Containers {
			flags := containerFlags[containerKey{podTemplate.Key, containerTemplate.Key}]
			containerSpec, ret := buildVictimContainerSpec(containerTemplate, flags, testMode)
			if !ret.OK {
				return model.VictimSpec{}, ret
			}
			podSpec.Containers = append(podSpec.Containers, containerSpec)
		}
		spec.Pods = append(spec.Pods, podSpec)

		for _, network := range podTemplate.Networks {
			subnet, ok := networkPlans[network.Definition.Name]
			if !ok {
				subnet = &model.Subnet{
					DefName:      network.Definition.Name,
					Name:         fmt.Sprintf("net-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandHexStr(6)),
					CIDRBlock:    network.Definition.CIDR,
					Gateway:      network.Definition.Gateway,
					NetAttachDef: &model.NetAttachDef{Name: fmt.Sprintf("nad-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandHexStr(6))},
				}
				networkPlans[network.Definition.Name] = subnet
				spec.NetworkPlan.Subnets = append(spec.NetworkPlan.Subnets, subnet)
			}
		}
	}

	return spec, model.SuccessRetVal()
}

func buildVictimContainerSpec(template model.ChallengeContainerTemplate, flags []victimFlagSpec, testMode bool) (model.VictimContainerSpec, model.RetVal) {
	spec := model.VictimContainerSpec{
		Key:   template.Key,
		Name:  template.Name,
		Image: template.Image,
		Resources: model.ResourceSpec{
			CPUMillis:   int64(template.CPU * 1000),
			MemoryBytes: template.Memory,
		},
		WorkingDir:   template.WorkingDir,
		Command:      append([]string(nil), template.Command...),
		Environment:  make(map[string]string),
		KubeVirt:     template.KubeVirt,
		Bootloader:   template.Bootloader,
		SecureBoot:   template.SecureBoot,
		UserData:     template.UserData,
		VolumeMounts: append([]model.XVolume(nil), template.VolumeMounts...),
		Exposes:      append(model.Exposes(nil), template.Exposes...),
	}
	maps.Copy(spec.Environment, template.Environment)
	for _, flag := range flags {
		value := flag.Value
		if flag.HasTeam {
			value = flag.TeamValue
		} else if testMode {
			value = fmt.Sprintf("flag{%s}", renderChallengeFlagValue(flag.Value))
		}
		switch flag.Binding.Type {
		case model.EnvFlagBindingType:
			spec.Environment[flag.Binding.Target] = value
		case model.FileFlagBindingType:
			for i := range spec.VolumeMounts {
				volume := &spec.VolumeMounts[i]
				if volume.Path == flag.Binding.Target {
					volume.Content = strings.ReplaceAll(volume.Content, flag.Value, value)
				}
			}
		case model.CloudInitFileFlagBindingType:
			for i := range spec.UserData.WriteFiles {
				file := &spec.UserData.WriteFiles[i]
				if file.Path == flag.Binding.Target {
					file.Content = strings.ReplaceAll(file.Content, flag.Value, value)
				}
			}
		default:
			return model.VictimContainerSpec{}, model.FailRetVal(i18n.Model.ChallengeFlag.InvalidType, "ChallengeFlag", "type="+string(flag.Binding.Type), nil)
		}
	}
	return spec, model.SuccessRetVal()
}

func buildPodRecords(victim model.Victim) []model.Pod {
	pods := make([]model.Pod, 0, len(victim.Spec.Pods))
	for _, podSpec := range victim.Spec.Pods {
		name := strings.ToLower(podSpec.Key)
		name = name[:min(len(name), 15)]
		pods = append(pods, model.Pod{
			VictimID: victim.ID,
			Name:     fmt.Sprintf("pod-%d-%d-%s-%s", victim.ContestChallengeID.V, victim.UserID, name, utils.RandHexStr(6)),
			Spec:     podSpec,
		})
	}
	return pods
}

// StartVictim 启动靶机。
// contestChallengeID 大于 0 时使用比赛题目副本，为 0 时是管理员对题库题目的测试模式。
func StartVictim(tx *gorm.DB, userID, teamID, contestID uint, contestChallengeID, challengeID uint, durationL ...time.Duration) model.RetVal {
	var victim model.Victim
	victimRepo := db.InitVictimRepo(tx)
	var (
		source    victimSpecSource
		ret       model.RetVal
		hasVictim bool
	)
	if contestChallengeID > 0 {
		if _, ret = victimRepo.HasAliveVictimByContestChallenge(teamID, contestChallengeID); ret.OK {
			hasVictim = true
		} else if ret.Msg != i18n.Model.NotFound {
			return ret
		}
	} else {
		if _, ret = victimRepo.HasAliveVictim(teamID, challengeID); ret.OK {
			hasVictim = true
		} else if ret.Msg != i18n.Model.NotFound {
			return ret
		}
	}
	if hasVictim {
		log.Logger.Debugf("Start victim rejected: team_id=%d contest_challenge_id=%d challenge_id=%d already has alive victim", teamID, contestChallengeID, challengeID)
		return model.FailRetVal(i18n.Model.Victim.NotStartable, "Victim", fmt.Sprintf("team_id=%d,contest_challenge_id=%d,challenge_id=%d", teamID, contestChallengeID, challengeID), nil)
	}
	baseVictim := model.Victim{
		ChallengeID:        challengeID,
		ContestID:          sql.Null[uint]{V: contestID, Valid: contestID > 0},
		ContestChallengeID: sql.Null[uint]{V: contestChallengeID, Valid: contestChallengeID > 0},
		TeamID:             sql.Null[uint]{V: teamID, Valid: teamID > 0},
		UserID:             userID,
	}
	if contestChallengeID > 0 {
		contestChallenge, getRet := db.InitContestChallengeRepo(tx).GetByID(contestChallengeID, db.GetOptions{
			Preloads: map[string]db.GetOptions{"ContestFlags": {}},
		})
		if !getRet.OK {
			return getRet
		}
		source, ret = buildVictimSpecSourceContestChallenge(tx, baseVictim, contestChallenge)
		if !ret.OK {
			return ret
		}
	} else {
		challenge, getRet := db.InitChallengeRepo(tx).GetByID(challengeID, db.GetOptions{
			Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
		})
		if !getRet.OK {
			return getRet
		}
		source = buildVictimSpecSourceChallenge(challenge)
	}
	duration := 2 * time.Hour
	if len(durationL) > 0 && durationL[0] > 0 {
		duration = durationL[0]
	}
	spec, ret := buildVictimSpec(tx, baseVictim, source)
	if !ret.OK {
		return ret
	}
	ret = db.WithTransactionDB(tx, func(tx2 *db.Tx) model.RetVal {
		victimRepo := db.InitVictimRepo(tx2)
		podRepo := db.InitPodRepo(tx2)
		created, createRet := victimRepo.Create(model.Victim{
			UserID:             userID,
			TeamID:             sql.Null[uint]{V: teamID, Valid: teamID > 0},
			ContestID:          sql.Null[uint]{V: contestID, Valid: contestID > 0},
			ContestChallengeID: sql.Null[uint]{V: contestChallengeID, Valid: contestChallengeID > 0},
			ChallengeID:        challengeID,
			Start:              time.Now(),
			Duration:           duration,
			Spec:               spec,
			Status:             model.WaitingVictimStatus,
		})
		if !createRet.OK {
			return createRet
		}
		victim = created
		pods, createRet := podRepo.CreateBatch(buildPodRecords(victim))
		if !createRet.OK {
			return createRet
		}
		victim.Pods = append(victim.Pods, pods...)
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return ret
	}
	_, err := task.EnqueueStartVictimTask(victim)
	if err != nil {
		log.Logger.Warningf("Failed to enqueue start victim task: victim_id=%d user_id=%d team_id=%d challenge_id=%d error=%v", victim.ID, userID, teamID, challengeID, err)
		if cleanupRet := db.InitVictimRepo(tx).Delete(victim.ID); !cleanupRet.OK {
			log.Logger.Warningf("Failed to cleanup victim after enqueue failure: victim_id=%d reason=%s", victim.ID, cleanupRet.Msg)
		}
		return model.FailRetVal(i18n.Task.EnqueueError, "Victim", fmt.Sprintf("victim_id=%d", victim.ID), err)
	}
	log.Logger.Infof(
		"Start victim queued: victim_id=%d user_id=%d team_id=%d challenge_id=%d contest_id=%d duration=%s pods=%d",
		victim.ID, userID, teamID, challengeID, contestID, duration, len(victim.Pods),
	)
	return model.SuccessRetVal()
}

// GetVictimStatus 返回队伍在指定题目上的靶机状态。
// contestChallengeID 大于 0 时按比赛题目查找，否则按题库题目查找（测试模式）。
func GetVictimStatus(tx *gorm.DB, teamID, contestChallengeID, challengeID uint, challengeType model.ChallengeType) view.VictimStatusView {
	data := view.VictimStatusView{
		Targets:   make([]string, 0),
		Duration:  0,
		Remaining: 0,
		Status:    "Down",
	}
	if challengeType != model.PodsChallengeType {
		data.Status = "not_docker"
		return data
	}
	var (
		victim model.Victim
		ret    model.RetVal
	)
	if contestChallengeID > 0 {
		victim, ret = db.InitVictimRepo(tx).HasAliveVictimByContestChallenge(teamID, contestChallengeID)
	} else {
		victim, ret = db.InitVictimRepo(tx).HasAliveVictim(teamID, challengeID)
	}
	if !ret.OK {
		return data
	}
	data.Targets = victim.RemoteAddr()
	data.Duration = victim.Duration.Seconds()
	data.Status = victim.Status
	data.Remaining = victim.Remaining().Seconds()
	return data
}

// ExtendVictimDuration 延长队伍在比赛题目上的靶机时长。
func ExtendVictimDuration(tx *gorm.DB, team model.Team, contestChallengeID uint) (model.Victim, model.RetVal) {
	repo := db.InitVictimRepo(tx)
	victim, ret := repo.HasAliveVictimByContestChallenge(team.ID, contestChallengeID)
	if !ret.OK {
		return model.Victim{}, ret
	}
	if victim.Status != model.RunningVictimStatus {
		return model.Victim{}, model.FailRetVal(i18n.Model.Victim.NotExtendable, "Victim", fmt.Sprintf("victim_id=%d,status=%s", victim.ID, victim.Status), nil)
	}
	if !victim.Start.Add(victim.Duration).Before(time.Now().Add(20 * time.Minute)) {
		return model.Victim{}, model.FailRetVal(i18n.Model.Victim.HasMuchTime, "Victim", fmt.Sprintf("victim_id=%d,duration=%s", victim.ID, victim.Duration), nil)
	}
	duration := victim.Duration + time.Hour
	if ret = repo.Update(victim.ID, db.UpdateVictimOptions{
		Duration: &duration,
	}); !ret.OK {
		return model.Victim{}, ret
	}
	victim.Duration = duration
	log.Logger.Infof("Victim duration extended: victim_id=%d team_id=%d challenge_id=%d duration=%s", victim.ID, victim.TeamID.V, victim.ChallengeID, duration)
	return victim, model.SuccessRetVal()
}

func StopVictim(tx *gorm.DB, victim model.Victim) model.RetVal {
	if victim.Status != model.RunningVictimStatus {
		return model.FailRetVal(i18n.Model.Victim.NotStoppable, "Victim", fmt.Sprintf("victim_id=%d,status=%s", victim.ID, victim.Status), nil)
	}
	return ForceStopVictim(tx, victim)
}

// StopAliveVictim 停止队伍在比赛题目上的靶机。
func StopAliveVictim(tx *gorm.DB, team model.Team, contestChallengeID uint) model.RetVal {
	victim, ret := db.InitVictimRepo(tx).HasAliveVictimByContestChallenge(team.ID, contestChallengeID)
	if !ret.OK {
		return ret
	}
	return StopVictim(tx, victim)
}

func ForceStopVictim(tx *gorm.DB, victim model.Victim) model.RetVal {
	if victim.Status == model.TerminatingVictimStatus {
		log.Logger.Debugf("Stop victim skipped: victim_id=%d already terminating", victim.ID)
		return model.SuccessRetVal()
	}
	repo := db.InitVictimRepo(tx)
	originalStatus := victim.Status
	if ret := repo.UpdateIfStatus(victim.ID, originalStatus, db.UpdateVictimOptions{Status: new(model.TerminatingVictimStatus)}); !ret.OK {
		return ret
	}
	victim.Status = model.TerminatingVictimStatus
	err := task.EnqueueStopVictimTask(victim)
	if err != nil {
		log.Logger.Warningf(
			"Failed to enqueue stop victim task: victim_id=%d user_id=%d team_id=%d challenge_id=%d error=%v",
			victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID, err,
		)
		_ = repo.UpdateIfStatus(victim.ID, model.TerminatingVictimStatus, db.UpdateVictimOptions{Status: new(originalStatus)})
		return model.FailRetVal(i18n.Common.UnknownError, "Victim", fmt.Sprintf("victim_id=%d", victim.ID), err)
	}
	log.Logger.Infof("Stop victim queued: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
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
	options.Conditions["status"] = model.RunningVictimStatus
	running, ret := db.InitVictimRepo(tx).Count(db.CountOptions{Conditions: options.Conditions})
	if !ret.OK {
		return nil, 0, 0, ret
	}
	return victims, running, count, ret
}

func ListVictimHistories(tx *gorm.DB, team model.Team, form dto.ListModelsForm) ([]model.Victim, int64, model.RetVal) {
	return db.InitVictimRepo(tx).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
		Deleted:    true,
		Sort:       []string{"id DESC"},
	})
}

func StartVictims(tx *gorm.DB, contest model.Contest, form dto.StartVictimsForm) model.RetVal {
	if len(form.Challenges) == 0 || form.TeamRatio <= 0 || form.TeamRatio > 1 {
		return model.SuccessRetVal()
	}
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{
			"contest_id": contest.ID,
			"rand_id":    form.Challenges,
			"type":       model.PodsChallengeType,
		},
		Preloads: map[string]db.GetOptions{"ContestFlags": {}},
	})
	if !ret.OK {
		return ret
	}
	teams, _, ret := db.InitTeamRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if !ret.OK {
		return ret
	}
	if len(teams) == 0 {
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
	byRandID := make(map[string][]model.ContestChallenge, len(contestChallenges))
	for _, contestChallenge := range contestChallenges {
		byRandID[contestChallenge.RandID] = append(byRandID[contestChallenge.RandID], contestChallenge)
	}
	duration := time.Duration(form.Duration) * time.Second
	batch := model.NewBatch(len(form.Challenges) * len(teams))
	for _, randID := range form.Challenges {
		contestChallengeL := byRandID[randID]
		for _, team := range teams {
			if tx.Statement.Context.Err() != nil {
				return batch.Result(tx.Statement.Context)
			}
			if len(contestChallengeL) == 0 {
				batch.Fail(fmt.Sprintf("team:%d/challenge:%s", team.ID, randID), "lookup", model.FailRetVal(i18n.Model.NotFound, "ContestChallenge", fmt.Sprintf("rand_id=%s", randID), nil))
				continue
			}
			for _, contestChallenge := range contestChallengeL {
				id := fmt.Sprintf("team:%d/challenge:%s/instance:%d", team.ID, randID, contestChallenge.ID)
				solved, checkRet := CheckIfSolved(tx, team, contestChallenge.ContestFlags)
				if !checkRet.OK {
					batch.Fail(id, "read_solved", checkRet)
					return batch.Result(tx.Statement.Context)
				}
				if solved {
					batch.Skip(id, "already_solved")
					continue
				}
				generated, checkRet := CheckIfGenerated(tx, team, contestChallenge.ContestFlags)
				if !checkRet.OK {
					batch.Fail(id, "read_flags", checkRet)
					return batch.Result(tx.Statement.Context)
				}
				if !generated {
					if _, ret = CreateTeamFlag(tx, team, contest, contestChallenge); !ret.OK {
						batch.Fail(id, "create_flags", ret)
						if model.BatchDependencyFailed(ret) {
							return batch.Result(tx.Statement.Context)
						}
						continue
					}
				}
				if ret = StartVictim(tx, team.CaptainID, team.ID, contest.ID, contestChallenge.ID, contestChallenge.ChallengeID, duration); ret.OK {
					batch.Success(id, "queued")
				} else if ret.Msg == i18n.Model.Victim.NotStartable {
					batch.Skip(id, "already_active")
				} else {
					batch.Fail(id, "start", ret)
					if model.BatchDependencyFailed(ret) {
						return batch.Result(tx.Statement.Context)
					}
				}
			}
		}
	}
	return batch.Result(tx.Statement.Context)
}

func StopVictims(tx *gorm.DB, contestID uint, form dto.StopVictimsForm) model.RetVal {
	if len(form.Victims) == 0 {
		return model.SuccessRetVal()
	}
	options := db.GetOptions{Conditions: map[string]any{"id": form.Victims}, Deleted: true}
	if contestID > 0 {
		options.Conditions["contest_id"] = contestID
	}
	victims, _, ret := db.InitVictimRepo(tx).List(-1, -1, options)
	if !ret.OK {
		return ret
	}
	byID := make(map[uint]model.Victim, len(victims))
	for _, victim := range victims {
		byID[victim.ID] = victim
	}
	batch := model.NewBatch(len(form.Victims))
	for _, id := range form.Victims {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		key := fmt.Sprint(id)
		victim, exists := byID[id]
		if !exists {
			batch.Fail(key, "lookup", model.FailRetVal(i18n.Model.NotFound, "Victim", "id="+key, nil))
			continue
		}
		if victim.Status == model.TerminatingVictimStatus || victim.Status == model.StoppedVictimStatus {
			batch.Skip(key, "already_stopping")
			continue
		}
		ret = StopVictim(tx, victim)
		if ret.OK {
			batch.Success(key, "queued")
		} else {
			batch.Fail(key, "stop", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(tx.Statement.Context)
			}
		}
	}
	return batch.Result(tx.Statement.Context)
}
