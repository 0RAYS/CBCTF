package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

const (
	startRangeTaskType   = "tasks:range:start"
	stopRangeTaskType    = "tasks:range:stop"
	gatewayRangeTaskType = "tasks:range:gateway"
)

// EnqueueStartRangeTask 启动一个池化的靶场实例。
// 靶场实例的节点多、启动慢,单独排队,不和每队靶机抢同一个池子。
func EnqueueStartRangeTask(victim model.Victim) (*asynq.TaskInfo, error) {
	return enqueueMsgpackTask(startRangeTaskType, StartVictimPayload{victim}, asynq.MaxRetry(3), asynq.Timeout(15*time.Minute))
}

// EnqueueStopRangeTask 销毁一个靶场实例。
func EnqueueStopRangeTask(victim model.Victim) error {
	_, err := enqueueMsgpackTask(stopRangeTaskType, StopVictimPayload{victim}, asynq.MaxRetry(3), asynq.Timeout(5*time.Minute))
	return err
}

type RangeGatewayPayload struct {
	Victim model.Victim
}

// EnqueueRangeGatewayTask 为已租出的实例创建玩家网关。
func EnqueueRangeGatewayTask(victim model.Victim) (*asynq.TaskInfo, error) {
	return enqueueMsgpackTask(gatewayRangeTaskType, RangeGatewayPayload{victim}, asynq.MaxRetry(3), asynq.Timeout(2*time.Minute))
}

// HandleRangeGatewayTask 创建玩家网关并把 flag 下发到网关。
// flag 不进实例环境,环境内容因此与队伍无关,同批实例逐字节一致。
func HandleRangeGatewayTask(ctx context.Context, t *asynq.Task) error {
	payload, err := decodeMsgpackTask[RangeGatewayPayload](t)
	if err != nil {
		return err
	}
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "victim", payload.Victim.ID, func() error {
		repo := db.InitVictimRepo(db.TaskDB.WithContext(ctx))
		victim, ret := repo.GetByID(payload.Victim.ID, db.GetOptions{Preloads: map[string]db.GetOptions{"Pods": {}}})
		if !ret.OK {
			if ret.Msg == i18n.Model.NotFound {
				return nil
			}
			return fmt.Errorf("get range victim failed: %w", ret.AsError())
		}
		if victim.Status != model.RunningVictimStatus || !victim.TeamID.Valid {
			log.Logger.Debugf("Range gateway skipped: victim_id=%d status=%s team_id=%d", victim.ID, victim.Status, victim.TeamID.V)
			return nil
		}
		challenge, ret := db.InitChallengeRepo(db.TaskDB.WithContext(ctx)).GetByID(victim.ChallengeID, db.GetOptions{
			Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
		})
		if !ret.OK {
			return fmt.Errorf("get range challenge failed: %w", ret.AsError())
		}
		if challenge.RangeConfig == nil {
			return fmt.Errorf("challenge %d is not a range challenge", challenge.ID)
		}
		gateway, ok := model.GatewayOf(challenge.Template, challenge.RangeConfig.AccessNetwork)
		if !ok {
			return fmt.Errorf("challenge %d has no player gateway", challenge.ID)
		}
		teamFlags, ret := rangeTeamFlags(ctx, victim.TeamID.V, challenge.ChallengeFlags)
		if !ret.OK {
			return ret.AsError()
		}
		known := make(map[string]struct{}, len(victim.Pods))
		for _, pod := range victim.Pods {
			known[pod.Name] = struct{}{}
		}
		victim, ret = k8s.StartRangeGateway(ctx, victim, gateway, buildGatewayPayload(challenge, gateway, teamFlags, victim))
		if !ret.OK {
			return fmt.Errorf("create range gateway failed: %w", ret.AsError())
		}
		newPods := make([]model.Pod, 0)
		for _, pod := range victim.Pods {
			if _, ok := known[pod.Name]; ok {
				continue
			}
			newPods = append(newPods, model.Pod{VictimID: victim.ID, Name: pod.Name})
		}
		if len(newPods) > 0 {
			persisted, ret := db.InitPodRepo(db.TaskDB.WithContext(ctx)).CreateBatch(newPods)
			if !ret.OK {
				return fmt.Errorf("create gateway pod record failed: %w", ret.AsError())
			}
			victim.Pods = append(victim.Pods, persisted...)
		}
		if ret = repo.Update(victim.ID, db.UpdateVictimOptions{
			Resources:        &victim.Resources,
			Endpoints:        &victim.Endpoints,
			ExposedEndpoints: &victim.ExposedEndpoints,
		}); !ret.OK {
			return fmt.Errorf("update range gateway endpoints failed: %w", ret.AsError())
		}
		log.Logger.Infof("Range gateway ready: victim_id=%d team_id=%d endpoints=%d", victim.ID, victim.TeamID.V, len(victim.ExposedEndpoints))
		return nil
	})
}

// rangeTeamFlags 读取队伍在该题目上持有的 flag 值。
func rangeTeamFlags(ctx context.Context, teamID uint, flags []model.ChallengeFlag) (map[uint]string, model.RetVal) {
	ids := make([]uint, 0, len(flags))
	for _, flag := range flags {
		ids = append(ids, flag.ID)
	}
	values := make(map[uint]string)
	if len(ids) == 0 {
		return values, model.SuccessRetVal()
	}
	teamFlags, _, ret := db.InitTeamFlagRepo(db.TaskDB.WithContext(ctx)).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": teamID, "challenge_flag_id": ids},
	})
	if !ret.OK && ret.Msg != i18n.Model.NotFound {
		return nil, ret
	}
	for _, teamFlag := range teamFlags {
		values[teamFlag.ChallengeFlagID] = teamFlag.Value
	}
	return values, model.SuccessRetVal()
}

// buildGatewayPayload 生成租约载荷。
// flag 按出题人的绑定方式落到环境变量或文件,网关镜像按约定读取载荷目录。
func buildGatewayPayload(challenge model.Challenge, gateway model.RangeGateway, teamFlags map[uint]string, victim model.Victim) model.GatewayPayload {
	payload := model.GatewayPayload{
		Files: make(map[string]string),
		Env:   make(map[string]string),
	}
	paths := make(map[string]string)
	for _, flag := range challenge.ChallengeFlags {
		if !flag.Active {
			continue
		}
		value := flag.Value
		if teamValue, ok := teamFlags[flag.ID]; ok {
			value = teamValue
		} else {
			value = model.RenderFlagValue(flag.Value)
		}
		switch flag.Binding.Type {
		case model.EnvFlagBindingType:
			payload.Env[flag.Binding.Target] = value
		case model.FileFlagBindingType, model.CloudInitFileFlagBindingType:
			key := fmt.Sprintf("flag-%d", flag.ID)
			payload.Files[key] = value
			paths[key] = flag.Binding.Target
		}
	}
	user := ""
	keys := make([]string, 0)
	for _, userData := range gateway.UserData.Users {
		if user == "" {
			user = userData.Name
		}
		keys = append(keys, userData.SSHAuthorizedKeys...)
	}
	credentials, err := json.Marshal(map[string]any{
		"challenge":           challenge.Name,
		"team_id":             victim.TeamID.V,
		"lease_expire_at":     victim.LeaseExpireAt,
		"user":                user,
		"ssh_authorized_keys": keys,
		"paths":               paths,
	})
	if err != nil {
		log.Logger.Warningf("Failed to render gateway credentials: %v", err)
		return payload
	}
	payload.Credentials = string(credentials)
	return payload
}
