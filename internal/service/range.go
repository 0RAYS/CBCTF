package service

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
)

// rangeFlagTemplates 是实例里要清掉的 flag 模板。
var rangeFlagTemplates = []*regexp.Regexp{model.StaticFlagTmpl, model.DynamicFlagTmpl, model.UUIDFlagTmpl}

// buildRangeSpec 构造池化靶场实例的规格。
// 实例不携带队伍 flag,玩家网关也不在实例里,环境因此与队伍无关。
func buildRangeSpec(challenge model.Challenge, victim model.Victim) (model.VictimSpec, model.RetVal) {
	template := model.ChallengeTemplate{Pods: make([]model.ChallengePodTemplate, 0, len(challenge.Template.Pods))}
	for _, pod := range challenge.Template.Pods {
		if pod.Role == model.RangeGatewayRole {
			continue
		}
		template.Pods = append(template.Pods, pod)
	}
	source := victimSpecSource{Template: template, NetworkPolicies: challenge.NetworkPolicies}
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
	stripRangeFlags(&spec)
	return spec, model.SuccessRetVal()
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

// stripRangeFlags 清掉实例里的 flag 模板。
// 靶场环境不携带 flag,flag 由平台下发到玩家网关。
func stripRangeFlags(spec *model.VictimSpec) {
	for i := range spec.Pods {
		for j := range spec.Pods[i].Containers {
			container := &spec.Pods[i].Containers[j]
			for k := range container.VolumeMounts {
				volume := &container.VolumeMounts[k]
				for _, template := range rangeFlagTemplates {
					volume.Content = template.ReplaceAllString(volume.Content, "")
				}
			}
			for k := range container.UserData.WriteFiles {
				file := &container.UserData.WriteFiles[k]
				for _, template := range rangeFlagTemplates {
					file.Content = template.ReplaceAllString(file.Content, "")
				}
			}
		}
	}
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
	spec, ret := buildRangeSpec(challenge, instance)
	if !ret.OK {
		return model.Victim{}, ret
	}
	instance.Spec = spec
	instance.TemplateKey = model.TemplateKey(spec)
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
func StartRange(tx *gorm.DB, team model.Team, contest model.Contest, challenge model.Challenge) (model.Victim, model.RetVal) {
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
