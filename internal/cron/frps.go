package cron

import (
	"context"
	"slices"
	"strings"
	"time"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
)

// syncFrpsPortLocksTask 以数据库中的活跃靶机为准校准 FRPS 端口占用缓存。
// pending 靶机可能已锁定端口但尚未写回 ExposedEndpoints，跳过本轮避免误释放。
func syncFrpsPortLocksTask() model.RetVal {
	// 系统配置发生改变时可能会导致存在 key 残留, 但问题不大
	if !config.Env.K8S.Frp.On {
		return model.SuccessRetVal()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	revision, ret := redis.FrpsRevision(ctx)
	if !ret.OK {
		return ret
	}
	victimRepo := db.InitVictimRepo(db.CronDB.WithContext(ctx))
	pendingVictims, _, ret := victimRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"status": model.PendingVictimStatus},
	})
	if !ret.OK {
		return ret
	}
	if len(pendingVictims) > 0 {
		victim := pendingVictims[0]
		log.Logger.Debugf("Skip FRPS port lock reconciliation while victim is provisioning: victim_id=%d updated_at=%s", victim.ID, victim.UpdatedAt.Format(time.RFC3339))
		return model.SuccessRetVal()
	}

	expected := make(map[string]map[string][]int32)
	addExpectedKey := func(host, protocol string) {
		protocol = strings.ToLower(protocol)
		if expected[host] == nil {
			expected[host] = make(map[string][]int32)
		}
		if _, ok := expected[host][protocol]; !ok {
			expected[host][protocol] = make([]int32, 0)
		}
	}
	addExpectedPort := func(host, protocol string, port int32) {
		protocol = strings.ToLower(protocol)
		addExpectedKey(host, protocol)
		if slices.Contains(expected[host][protocol], port) {
			return
		}
		expected[host][protocol] = append(expected[host][protocol], port)
	}
	for _, frps := range config.Env.K8S.Frp.Frps {
		addExpectedKey(frps.Host, "tcp")
		addExpectedKey(frps.Host, "udp")
	}

	victims, _, ret := victimRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"status": []string{model.RunningVictimStatus, model.TerminatingVictimStatus}},
	})
	if !ret.OK {
		return ret
	}
	for _, victim := range victims {
		for _, endpoint := range victim.ExposedEndpoints {
			addExpectedPort(endpoint.IP, endpoint.Protocol, endpoint.Port)
		}
	}

	removedKeys, keptPorts, ret := redis.ReconcileFrpsPorts(expected, revision)
	if !ret.OK {
		return ret
	}
	log.Logger.Infof("FRPS port locks reconciled: active_victims=%d kept_ports=%d removed_keys=%d", len(victims), keptPorts, removedKeys)
	return model.SuccessRetVal()
}
