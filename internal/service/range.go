package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
)

// buildRangeSpec 构造池化靶场实例的规格。
// flag 渲染进实例环境,取值随实例而定,队伍领取时按这个取值给答案。
func buildRangeSpec(challenge model.Challenge, victim model.Victim, flags []victimFlagSpec) (model.VictimSpec, model.RetVal) {
	template := model.ChallengeTemplate{Pods: make([]model.ChallengePodTemplate, 0, len(challenge.Template.Pods))}
	for _, pod := range challenge.Template.Pods {
		if pod.Role == model.RangeGatewayRole {
			continue
		}
		template.Pods = append(template.Pods, pod)
	}
	source := victimSpecSource{
		Template:        template,
		NetworkPolicies: challenge.NetworkPolicies,
		Flags:           flags,
	}
	spec, ret := buildVictimSpec(nil, victim, source)
	if !ret.OK {
		return model.VictimSpec{}, ret
	}
	spec.FrpEnabled = false
	if challenge.RangeConfig != nil {
		spec.Routes = challenge.RangeConfig.Routes
	}
	if ret = addRangeAccessSubnet(&spec, challenge, victim); !ret.OK {
		return model.VictimSpec{}, ret
	}
	return spec, model.SuccessRetVal()
}

// rangeFlagSpecs 返回实例环境里要注入的 flag 及其取值。
// 取值在实例创建时渲染,队伍领取时按这个取值同步答案。
func rangeFlagSpecs(challenge model.Challenge, contest model.Contest) ([]victimFlagSpec, map[uint]string) {
	flags := make([]victimFlagSpec, 0, len(challenge.ChallengeFlags))
	values := make(map[uint]string)
	for _, flag := range challenge.ChallengeFlags {
		if !flag.Active || flag.IsFixed() {
			continue
		}
		value := model.RenderFlagValue(flag.Value)
		if contest.Prefix != "" {
			value = fmt.Sprintf("%s{%s}", contest.Prefix, value)
		}
		flags = append(flags, victimFlagSpec{
			Key:       flag.ID,
			Value:     flag.Value,
			Binding:   flag.Binding,
			TeamValue: value,
			HasTeam:   true,
		})
		values[flag.ID] = value
	}
	return flags, values
}

// rangeTemplateKey 计算靶场实例的拓扑指纹。
// flag 取值按模板算,实例里的实际取值不参与,同构实例的指纹因此一致。
func rangeTemplateKey(challenge model.Challenge, victim model.Victim, flags []victimFlagSpec) string {
	templateFlags := make([]victimFlagSpec, len(flags))
	for i, flag := range flags {
		templateFlags[i] = flag
		templateFlags[i].HasTeam = false
		templateFlags[i].TeamValue = ""
	}
	spec, ret := buildRangeSpec(challenge, victim, templateFlags)
	if !ret.OK {
		return ""
	}
	return model.TemplateKey(spec)
}

// addRangeAccessSubnet 把玩家接入网段写进实例的网段规划。
// 网关在租约时挂到这个网段上,网段因此由实例创建,随实例一起回收。
func addRangeAccessSubnet(spec *model.VictimSpec, challenge model.Challenge, victim model.Victim) model.RetVal {
	if challenge.RangeConfig == nil {
		return model.SuccessRetVal()
	}
	gateway, ok := model.GatewayOf(challenge.Template, challenge.RangeConfig.AccessNetwork)
	if !ok || gateway.Subnet.CIDR == "" || gateway.Subnet.Gateway == "" {
		return model.FailRetVal(i18n.Model.CreateError, "RangeAccessNetwork", challenge.RangeConfig.AccessNetwork, nil)
	}
	for _, subnet := range spec.NetworkPlan.Subnets {
		if subnet == nil || subnet.DefName != gateway.Subnet.Name {
			continue
		}
		subnet.ExcludeIps = append(subnet.ExcludeIps, gateway.Network.IP)
		return model.SuccessRetVal()
	}
	spec.NetworkPlan.Subnets = append(spec.NetworkPlan.Subnets, &model.Subnet{
		DefName:      gateway.Subnet.Name,
		Name:         fmt.Sprintf("net-%d-%d-%s", victim.ChallengeID, victim.UserID.V, utils.RandHexStr(6)),
		CIDRBlock:    gateway.Subnet.CIDR,
		Gateway:      gateway.Subnet.Gateway,
		ExcludeIps:   []string{gateway.Network.IP},
		NetAttachDef: &model.NetAttachDef{Name: fmt.Sprintf("nad-%d-%d-%s", victim.ChallengeID, victim.UserID.V, utils.RandHexStr(6))},
	})
	return model.SuccessRetVal()
}

// syncRangeTeamFlags 把队伍在该题的答案同步成实例里的取值。
// 实例里的 flag 是渲染过的,队伍拿到的答案必须与之一致。
func syncRangeTeamFlags(tx *gorm.DB, team model.Team, challenge model.Challenge, instance model.Victim) model.RetVal {
	if len(instance.FlagValues) == 0 {
		return model.SuccessRetVal()
	}
	repo := db.InitTeamFlagRepo(tx)
	for _, flag := range challenge.ChallengeFlags {
		if !flag.Active || flag.IsFixed() {
			continue
		}
		value, ok := instance.FlagValues[strconv.FormatUint(uint64(flag.ID), 10)]
		if !ok || value == "" {
			continue
		}
		teamFlag, ret := repo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "challenge_flag_id": flag.ID},
		})
		if !ret.OK {
			if ret.Msg != i18n.Model.NotFound {
				return ret
			}
			if _, ret = repo.CreateIfAbsent(model.TeamFlag{
				TeamID:          team.ID,
				ChallengeFlagID: flag.ID,
				Value:           value,
			}); !ret.OK {
				return ret
			}
			continue
		}
		if teamFlag.Solved || teamFlag.Value == value {
			continue
		}
		if ret = repo.Update(teamFlag.ID, db.UpdateTeamFlagRepo{Value: &value}); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}

// CreateRangeInstance 创建池中的一个靶场实例。
// 实例不带队伍,内容与出题人设计的拓扑一致。
func CreateRangeInstance(tx *gorm.DB, challenge model.Challenge) (model.Victim, model.RetVal) {
	instance := model.Victim{
		ChallengeID: challenge.ID,
		ContestID:   challenge.ContestID,
		Kind:        model.VictimKindRange,
		Status:      model.WaitingVictimStatus,
		Start:       time.Now(),
	}
	if challenge.RangeConfig != nil {
		instance.Duration = challenge.RangeConfig.LeaseDuration
	}
	contest, ret := rangeContest(tx, challenge)
	if !ret.OK {
		return model.Victim{}, ret
	}
	// flag 随题目走,池里的实例与现场启动的实例要带同一份声明。
	challengeFlags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		return model.Victim{}, ret
	}
	challenge.ChallengeFlags = challengeFlags
	flags, values := rangeFlagSpecs(challenge, contest)
	spec, ret := buildRangeSpec(challenge, instance, flags)
	if !ret.OK {
		return model.Victim{}, ret
	}
	instance.Spec = spec
	instance.FlagValues = flagValueMap(values)
	instance.TemplateKey = rangeTemplateKey(challenge, instance, flags)
	repo := db.InitVictimRepo(tx)
	created, ret := repo.Create(instance)
	if !ret.OK {
		return model.Victim{}, ret
	}
	pods, ret := db.InitPodRepo(tx).CreateBatch(buildPodRecords(created))
	if !ret.OK {
		return model.Victim{}, ret
	}
	created.Pods = pods
	return created, model.SuccessRetVal()
}

// rangeContest 返回题目所属比赛,题目不归属比赛时返回空比赛。
func rangeContest(tx *gorm.DB, challenge model.Challenge) (model.Contest, model.RetVal) {
	if !challenge.ContestID.Valid {
		return model.Contest{}, model.SuccessRetVal()
	}
	return db.InitContestRepo(tx).GetByID(challenge.ContestID.V)
}

// flagValueMap 把 flag 取值按 id 存进实例,领取时据此同步队伍答案。
func flagValueMap(values map[uint]string) map[string]string {
	result := make(map[string]string, len(values))
	for id, value := range values {
		result[strconv.FormatUint(uint64(id), 10)] = value
	}
	return result
}

// EnsureRangePool 补齐靶场实例池。
// 池中的实例内容完全一致,队伍领取后即销毁,租约结束后补建新的。
func EnsureRangePool(ctx context.Context, challenge model.Challenge) error {
	if challenge.RangeConfig == nil || challenge.RangeConfig.PoolSize <= 0 {
		return nil
	}
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "range-pool", challenge.ID, func() error {
		repo := db.InitVictimRepo(db.TaskDB.WithContext(ctx))
		instances, ret := repo.ListRanges(challenge.ID, []string{
			model.WaitingVictimStatus,
			model.PendingVictimStatus,
			model.RunningVictimStatus,
		})
		if !ret.OK {
			return ret.AsError()
		}
		for i := len(instances); i < challenge.RangeConfig.PoolSize; i++ {
			instance, ret := CreateRangeInstance(db.TaskDB.WithContext(ctx), challenge)
			if !ret.OK {
				return ret.AsError()
			}
			if _, err := task.EnqueueStartRangeTask(instance); err != nil {
				log.Logger.Warningf("Failed to enqueue range instance start: victim_id=%d error=%v", instance.ID, err)
				return err
			}
			log.Logger.Infof("Range instance created: victim_id=%d challenge_id=%d template_key=%s", instance.ID, challenge.ID, instance.TemplateKey)
		}
		return nil
	})
}

// StartRange 为队伍领取一个靶场实例。
// 池中有空闲实例就直接租出,否则现场启动一个,实例就绪后再创建玩家网关。
func StartRange(tx *gorm.DB, team model.Team, challenge model.Challenge) (model.Victim, model.RetVal) {
	if challenge.RangeConfig == nil {
		return model.Victim{}, model.FailRetVal(i18n.Model.Victim.NotStartable, "Range", fmt.Sprintf("challenge_id=%d", challenge.ID), nil)
	}
	repo := db.InitVictimRepo(tx)
	alive, ret := repo.Get(db.GetOptions{Conditions: map[string]any{
		"kind":         model.VictimKindRange,
		"challenge_id": challenge.ID,
		"team_id":      team.ID,
		"status": []string{
			model.WaitingVictimStatus,
			model.PendingVictimStatus,
			model.RunningVictimStatus,
		},
	}})
	if ret.OK {
		return alive, model.SuccessRetVal()
	}
	if ret.Msg != i18n.Model.NotFound {
		return model.Victim{}, ret
	}
	// 领取时要把队伍在该题的答案同步成实例里的取值。
	flags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		return model.Victim{}, ret
	}
	challenge.ChallengeFlags = flags
	if challenge.RangeConfig.MaxLeased > 0 {
		leased, ret := repo.Count(db.CountOptions{Conditions: map[string]any{
			"kind":         model.VictimKindRange,
			"challenge_id": challenge.ID,
			"status": []string{
				model.WaitingVictimStatus,
				model.PendingVictimStatus,
				model.RunningVictimStatus,
			},
		}})
		if !ret.OK {
			return model.Victim{}, ret
		}
		if leased >= int64(challenge.RangeConfig.MaxLeased) {
			return model.Victim{}, model.FailRetVal(i18n.Model.Victim.Limited, "Range", fmt.Sprintf("challenge_id=%d,limit=%d", challenge.ID, challenge.RangeConfig.MaxLeased), nil)
		}
	}
	instances, ret := repo.ListRanges(challenge.ID, []string{model.RunningVictimStatus})
	if !ret.OK {
		return model.Victim{}, ret
	}
	for _, instance := range instances {
		if instance.TeamID.Valid {
			continue
		}
		leaseExpire := time.Now().Add(challenge.RangeConfig.LeaseDuration)
		if ret = repo.LeaseRange(instance.ID, team.ID, leaseExpire); !ret.OK {
			continue
		}
		instance.TeamID = sql.Null[uint]{V: team.ID, Valid: true}
		instance.LeaseExpireAt = leaseExpire
		instance.Start = time.Now()
		instance.Duration = challenge.RangeConfig.LeaseDuration
		if ret = repo.Update(instance.ID, db.UpdateVictimOptions{
			Start:         &instance.Start,
			Duration:      &instance.Duration,
			TeamID:        &instance.TeamID,
			LeaseExpireAt: &leaseExpire,
		}); !ret.OK {
			return model.Victim{}, ret
		}
		log.Logger.Infof("Range instance leased: victim_id=%d team_id=%d challenge_id=%d", instance.ID, team.ID, challenge.ID)
		if ret = syncRangeTeamFlags(tx, team, challenge, instance); !ret.OK {
			return model.Victim{}, ret
		}
		if _, err := task.EnqueueRangeGatewayTask(instance); err != nil {
			log.Logger.Warningf("Failed to enqueue range gateway task: victim_id=%d error=%v", instance.ID, err)
			return model.Victim{}, model.FailRetVal(i18n.Task.EnqueueError, "Range", fmt.Sprintf("victim_id=%d", instance.ID), err)
		}
		return instance, model.SuccessRetVal()
	}
	// 池里没有空闲实例,现场启动一个并直接租给队伍。
	instance, ret := CreateRangeInstance(tx, challenge)
	if !ret.OK {
		return model.Victim{}, ret
	}
	leaseExpire := time.Now().Add(challenge.RangeConfig.LeaseDuration)
	instance.TeamID = sql.Null[uint]{V: team.ID, Valid: true}
	instance.LeaseExpireAt = leaseExpire
	instance.Start = time.Now()
	instance.Duration = challenge.RangeConfig.LeaseDuration
	if ret = repo.Update(instance.ID, db.UpdateVictimOptions{
		Start:         &instance.Start,
		Duration:      &instance.Duration,
		TeamID:        &instance.TeamID,
		LeaseExpireAt: &leaseExpire,
	}); !ret.OK {
		return model.Victim{}, ret
	}
	if ret = syncRangeTeamFlags(tx, team, challenge, instance); !ret.OK {
		return model.Victim{}, ret
	}
	if _, err := task.EnqueueStartRangeTask(instance); err != nil {
		log.Logger.Warningf("Failed to enqueue range instance start: victim_id=%d error=%v", instance.ID, err)
		return model.Victim{}, model.FailRetVal(i18n.Task.EnqueueError, "Range", fmt.Sprintf("victim_id=%d", instance.ID), err)
	}
	log.Logger.Infof("Range instance queued: victim_id=%d team_id=%d challenge_id=%d", instance.ID, team.ID, challenge.ID)
	return instance, model.SuccessRetVal()
}

// StopRange 结束队伍的靶场租约。
// 实例用完即销毁,不回收复用,渗透过程会改环境。
func StopRange(tx *gorm.DB, victim model.Victim) model.RetVal {
	if !victim.IsRange() {
		return model.FailRetVal(i18n.Model.Victim.NotStoppable, "Range", fmt.Sprintf("victim_id=%d", victim.ID), nil)
	}
	return ForceStopVictim(tx, victim)
}

// GetRangeVictim 返回队伍在靶场题上的实例。
func GetRangeVictim(tx *gorm.DB, teamID, challengeID uint) (model.Victim, model.RetVal) {
	return db.InitVictimRepo(tx).Get(db.GetOptions{Conditions: map[string]any{
		"kind":         model.VictimKindRange,
		"challenge_id": challengeID,
		"team_id":      teamID,
	}})
}

// BuildRangePoolView 汇总一个靶场题目的实例池状态。
func BuildRangePoolView(tx *gorm.DB, challenge model.Challenge) (view.RangePoolView, model.RetVal) {
	instances, ret := db.InitVictimRepo(tx).ListRanges(challenge.ID, nil)
	if !ret.OK {
		return view.RangePoolView{}, ret
	}
	pool := view.RangePoolView{
		ChallengeID:   challenge.ID,
		ChallengeName: challenge.Name,
		Instances:     make([]view.RangeInstanceView, 0, len(instances)),
	}
	if challenge.RangeConfig != nil {
		pool.PoolSize = challenge.RangeConfig.PoolSize
		pool.MaxLeased = challenge.RangeConfig.MaxLeased
	}
	for _, instance := range instances {
		if instance.TeamID.Valid {
			pool.Leased++
		}
		pool.Instances = append(pool.Instances, view.RangeInstanceView{
			ID:            instance.ID,
			Status:        instance.Status,
			TeamID:        instance.TeamID.V,
			TeamName:      instance.Team.Name,
			LeaseExpireAt: instance.LeaseExpireAt,
			TemplateKey:   instance.TemplateKey,
			Nodes:         len(instance.Spec.Pods),
			Start:         instance.Start,
		})
	}
	return pool, model.SuccessRetVal()
}

// BuildRangePoolViews 汇总比赛里所有靶场题目的实例池状态。
func BuildRangePoolViews(tx *gorm.DB, contestID uint) ([]view.RangePoolView, model.RetVal) {
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contestID},
	})
	if !ret.OK {
		return nil, ret
	}
	pools := make([]view.RangePoolView, 0)
	for _, challenge := range challenges {
		if challenge.RangeConfig == nil {
			continue
		}
		pool, ret := BuildRangePoolView(tx, challenge)
		if !ret.OK {
			return nil, ret
		}
		pools = append(pools, pool)
	}
	return pools, model.SuccessRetVal()
}

// EnsureContestRangePools 补齐比赛里所有靶场题目的实例池。
// 与定时任务走同一条路径,手动触发只是多一个口子。
func EnsureContestRangePools(ctx context.Context, contestID uint) model.RetVal {
	challenges, _, ret := db.InitChallengeRepo(db.TaskDB.WithContext(ctx)).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contestID},
	})
	if !ret.OK {
		return ret
	}
	for _, challenge := range challenges {
		if challenge.RangeConfig == nil {
			continue
		}
		if err := EnsureRangePool(ctx, challenge); err != nil {
			log.Logger.Warningf("Failed to ensure range pool: challenge_id=%d error=%v", challenge.ID, err)
			return model.FailRetVal(i18n.Common.UnknownError, "RangePool", fmt.Sprintf("challenge_id=%d", challenge.ID), err)
		}
	}
	return model.SuccessRetVal()
}
