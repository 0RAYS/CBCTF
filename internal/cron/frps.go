package cron

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"slices"
	"strings"
	"time"
)

// syncFrpsPortLocksTask 以数据库中的活跃靶机为准校准 FRPS 端口占用缓存。
// 最近仍在 pending 的靶机可能刚锁定端口但尚未写回 ExposedEndpoints，此时跳过本轮避免误释放。
func syncFrpsPortLocksTask() model.RetVal {
	// 系统配置发生改变时可能会导致存在 key 残留, 但问题不大
	if !config.Env.K8S.Frp.On {
		return model.SuccessRetVal()
	}

	victimRepo := db.InitVictimRepo(db.CronDB)
	pendingVictims, _, ret := victimRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"status": model.PendingVictimStatus},
	})
	if !ret.OK {
		return ret
	}
	cutoff := time.Now().Add(-10 * time.Minute)
	for _, victim := range pendingVictims {
		if victim.UpdatedAt.After(cutoff) {
			log.Logger.Debugf("Skip FRPS port lock reconciliation while victim is provisioning: victim_id=%d updated_at=%s", victim.ID, victim.UpdatedAt.Format(time.RFC3339))
			return model.SuccessRetVal()
		}
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
			if endpoint.IP == "" || endpoint.Port <= 0 || endpoint.Port > 65535 {
				continue
			}
			addExpectedPort(endpoint.IP, endpoint.Protocol, endpoint.Port)
		}
	}

	removedKeys, keptPorts, ret := redis.ReconcileFrpsPorts(expected)
	if !ret.OK {
		return ret
	}
	log.Logger.Infof("FRPS port locks reconciled: active_victims=%d kept_ports=%d removed_keys=%d", len(victims), keptPorts, removedKeys)
	return model.SuccessRetVal()
}
