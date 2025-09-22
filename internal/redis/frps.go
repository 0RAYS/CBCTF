package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const frpsPortKey = "frps:%s:%s"

// 锁定端口范围中的随机可用端口
const lockFrpsPortScript = `
local key = KEYS[1]
local ports = cjson.decode(ARGV[1])
local protocol = ARGV[2]

local n = #ports
if n == 0 then
    return {0, 0}
end

-- 记录已尝试的索引，避免重复尝试
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
    return 1
end

return 0
`

func LockFrpsPort(host string, portRange []int32, protocol string) (int32, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	protocol = strings.ToLower(protocol)
	key := fmt.Sprintf(frpsPortKey, host, protocol)

	portsJSON, err := json.Marshal(portRange)
	if err != nil {
		return 0, false, fmt.Errorf("failed to marshal port range: %w", err)
	}
	result, err := RDB.Eval(ctx, lockFrpsPortScript, []string{key}, string(portsJSON), protocol).Result()
	if err != nil {
		return 0, false, err
	}
	resultSlice, ok := result.([]interface{})
	if !ok || len(resultSlice) != 2 {
		return 0, false, fmt.Errorf("unexpected result format from Lua script")
	}
	port, ok := resultSlice[0].(int64)
	if !ok {
		return 0, false, fmt.Errorf("invalid port in result")
	}
	success, ok := resultSlice[1].(int64)
	if !ok {
		return 0, false, fmt.Errorf("invalid success flag in result")
	}
	return int32(port), success == 1, nil
}

func UnlockFrpsPort(host string, port int32, protocol string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	protocol = strings.ToLower(protocol)
	key := fmt.Sprintf(frpsPortKey, host, protocol)

	_, err := RDB.Eval(ctx, unlockFrpsPortScript, []string{key}, port).Result()
	if err != nil {
		return err
	}
	return nil
}
