package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

const frpsPortKeyTmpl = "frps:%s:%s"
const frpsRevisionKey = "frps:revision"

// 锁定端口范围中的随机可用端口
const lockFrpsPortScript = `
local key = KEYS[1]
local ports = cjson.decode(ARGV[1])
local protocol = ARGV[2]

local n = #ports
if n == 0 then
    return {0, 0}
end

-- 记录已尝试的索引, 避免重复尝试
local tried = {}
local triedCount = 0

while triedCount < n do
    local idx = math.random(1, n)
    if not tried[idx] then
        tried[idx] = true
        triedCount = triedCount + 1

        local selected_port = ports[idx]
        if redis.call('SISMEMBER', key, selected_port) == 0 then
            redis.call('SADD', key, selected_port)
            redis.call('INCR', KEYS[2])
            return {selected_port, 1}
        end
    end
end

return {0, 0}
`

// 解锁端口
const unlockFrpsPortScript = `
local key = KEYS[1]
local port = tonumber(ARGV[1])

local is_member = redis.call('SISMEMBER', key, port)
if is_member == 1 then
    redis.call('SREM', key, port)
    redis.call('INCR', KEYS[2])
    return 1
end

return 0
`

func LockFrpsPort(host string, portRange []int32, protocol string) (int32, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	protocol = strings.ToLower(protocol)
	key := fmt.Sprintf(frpsPortKeyTmpl, host, protocol)

	portsJSON, err := json.Marshal(portRange)
	if err != nil {
		log.Logger.Warningf("Failed to encode frps port range: %s", err)
		return 0, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	result, err := RDB.Eval(ctx, lockFrpsPortScript, []string{key, frpsRevisionKey}, string(portsJSON), protocol).Result()
	if err != nil {
		log.Logger.Warningf("Failed to eval lua script: %s", err)
		return 0, model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	resultSlice, ok := result.([]any)
	if !ok || len(resultSlice) != 2 {
		return 0, model.RetVal{Msg: i18n.Redis.InvalidScriptResult}
	}
	port, ok := resultSlice[0].(int64)
	if !ok {
		return 0, model.RetVal{Msg: i18n.Redis.InvalidScriptPort}
	}
	success, ok := resultSlice[1].(int64)
	if !ok {
		return 0, model.RetVal{Msg: i18n.Redis.InvalidScriptSuccessFlag}
	}
	if success != 1 {
		return int32(port), model.RetVal{Msg: i18n.Redis.NoAvailablePort}
	}
	return int32(port), model.SuccessRetVal()
}

func UnlockFrpsPort(host string, port int32, protocol string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	protocol = strings.ToLower(protocol)
	key := fmt.Sprintf(frpsPortKeyTmpl, host, protocol)

	_, err := RDB.Eval(ctx, unlockFrpsPortScript, []string{key, frpsRevisionKey}, port).Result()
	if err != nil {
		log.Logger.Warningf("Failed to eval lua script: %s", err)
		return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// ReconcileFrpsPorts 使用数据库中仍活跃靶机的暴露端口重建 Redis FRPS 端口锁集合。
// expected 的结构为 host -> protocol -> ports；未出现在 expected 中的 frps:*:* 键会被删除。
func FrpsRevision(ctx context.Context) (string, model.RetVal) {
	value, err := RDB.Get(ctx, frpsRevisionKey).Result()
	if errors.Is(err, goredis.Nil) {
		return "0", model.SuccessRetVal()
	}
	if err != nil {
		return "", model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": frpsRevisionKey, "Error": err.Error()}}
	}
	return value, model.SuccessRetVal()
}

const replaceFrpsSnapshotScript = `
if (redis.call('GET', KEYS[1]) or '0') ~= ARGV[1] then return 0 end
local entries = cjson.decode(ARGV[2])
for _, entry in ipairs(entries) do
  if entry[2] ~= '' and redis.call('EXISTS', entry[2]) == 0 then return -1 end
end
for _, entry in ipairs(entries) do
  if entry[2] == '' then redis.call('DEL', entry[1])
  else redis.call('RENAME', entry[2], entry[1]); redis.call('PERSIST', entry[1]) end
end
redis.call('INCR', KEYS[1])
return 1
`

func ReconcileFrpsPorts(expected map[string]map[string][]int32, revision string) (int64, int64, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	expectedKeys := make(map[string][]any)
	for host, protocols := range expected {
		for protocol, ports := range protocols {
			key := fmt.Sprintf(frpsPortKeyTmpl, host, protocol)
			members := make([]any, 0, len(ports))
			for _, port := range ports {
				members = append(members, strconv.FormatInt(int64(port), 10))
			}
			expectedKeys[key] = members
		}
	}

	var cursor uint64
	keysToReplace := make(map[string]bool)
	removedKeys := int64(0)
	for {
		keys, next, err := RDB.Scan(ctx, cursor, "frps:*:*", 100).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan frps port locks: %s", err)
			return 0, 0, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": "frps:*:*", "Error": err.Error()}}
		}
		for _, key := range keys {
			keysToReplace[key] = true
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}

	keptPorts := int64(0)
	for key := range expectedKeys {
		keysToReplace[key] = true
	}
	entries := make([][2]string, 0, len(keysToReplace))
	for key := range keysToReplace {
		members, exists := expectedKeys[key]
		if !exists {
			removedKeys++
		}
		staging := ""
		if len(members) > 0 {
			staging = "frps-staging:" + utils.UUID()
			// Failed staging never destroys the currently authoritative lock set.
			pipe := RDB.TxPipeline()
			pipe.SAdd(ctx, staging, members...)
			pipe.Expire(ctx, staging, 5*time.Minute)
			if _, err := pipe.Exec(ctx); err != nil {
				return 0, 0, model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
			}
		}
		entries = append(entries, [2]string{key, staging})
		keptPorts += int64(len(members))
	}
	data, err := json.Marshal(entries)
	if err != nil {
		return 0, 0, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	result, err := RDB.Eval(ctx, replaceFrpsSnapshotScript, []string{frpsRevisionKey}, revision, string(data)).Int()
	if err != nil || result < 0 {
		return 0, 0, model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": frpsRevisionKey, "Error": fmt.Sprintf("snapshot replacement failed: %v (result %d)", err, result)}}
	}
	if result == 0 {
		log.Logger.Debug("FRPS reconciliation deferred: allocations changed during database snapshot")
		return 0, 0, model.SuccessRetVal()
	}
	return removedKeys, keptPorts, model.SuccessRetVal()
}
