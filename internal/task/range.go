package task

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"golang.org/x/crypto/ssh"

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

// HandleRangeGatewayTask 创建玩家跳板并把队伍的登录私钥记录在实例上。
// 跳板是平台固定容器,只提供 SSH 登录,题目相关的 flag 与文件都在实例环境里。
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
		challenge, ret := db.InitChallengeRepo(db.TaskDB.WithContext(ctx)).GetByID(victim.ChallengeID)
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
		access, privateKey, err := rangeAccess(challenge.RangeConfig.AccessUser())
		if err != nil {
			return fmt.Errorf("generate range access key failed: %w", err)
		}
		known := make(map[string]struct{}, len(victim.Pods))
		for _, pod := range victim.Pods {
			known[pod.Name] = struct{}{}
		}
		victim, ret = k8s.StartRangeGateway(ctx, victim, gateway, access)
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
			AccessKey:        &privateKey,
		}); !ret.OK {
			return fmt.Errorf("update range gateway endpoints failed: %w", ret.AsError())
		}
		log.Logger.Infof("Range gateway ready: victim_id=%d team_id=%d endpoints=%d", victim.ID, victim.TeamID.V, len(victim.ExposedEndpoints))
		return nil
	})
}

// rangeAccess 生成队伍登录跳板的密钥对。
// 公钥进跳板容器,私钥随状态接口下发给队伍。
func rangeAccess(user string) (model.RangeAccess, string, error) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return model.RangeAccess{}, "", err
	}
	sshPublic, err := ssh.NewPublicKey(public)
	if err != nil {
		return model.RangeAccess{}, "", err
	}
	block, err := ssh.MarshalPrivateKey(private, "cbctf-range")
	if err != nil {
		return model.RangeAccess{}, "", err
	}
	return model.RangeAccess{
		User:      user,
		PublicKey: strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPublic))),
	}, string(pem.EncodeToMemory(block)), nil
}
