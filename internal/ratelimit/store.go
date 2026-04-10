package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const redisRateLimitScript = `
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]

redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', now - window)

local count = redis.call('ZCARD', KEYS[1])
if count < limit then
	redis.call('ZADD', KEYS[1], now, member)
	count = count + 1
	redis.call('PEXPIRE', KEYS[1], window)
	return {1, count, limit - count, 0}
end

local oldest = redis.call('ZRANGE', KEYS[1], 0, 0, 'WITHSCORES')
local retry_after = 0
if oldest[2] ~= nil then
	retry_after = math.max(0, window - (now - tonumber(oldest[2])))
end
redis.call('PEXPIRE', KEYS[1], window)
return {0, count, 0, retry_after}
`

type RedisStore struct {
	client  func() *redis.Client
	keyFunc func(Rule, string) string
	timeout time.Duration
}

func NewRedisStoreFunc(client func() *redis.Client) *RedisStore {
	return &RedisStore{
		client:  client,
		keyFunc: defaultRedisKey,
		timeout: 3 * time.Second,
	}
}

func (s *RedisStore) Allow(ctx context.Context, rule Rule, subject string) (Decision, error) {
	if s.client == nil || s.client() == nil {
		return Decision{}, errors.New("rate limit redis client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if s.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}

	now := time.Now().UTC()
	result, err := s.client().Eval(
		ctx,
		redisRateLimitScript,
		[]string{s.keyFunc(rule, subject)},
		now.UnixMilli(),
		rule.Window.Milliseconds(),
		rule.Limit,
		uuid.NewString(),
	).Result()
	if err != nil {
		return Decision{}, err
	}

	values, ok := result.([]any)
	if !ok || len(values) != 4 {
		return Decision{}, fmt.Errorf("unexpected rate limit response: %T", result)
	}

	allowed, err := toInt64(values[0])
	if err != nil {
		return Decision{}, err
	}
	count, err := toInt64(values[1])
	if err != nil {
		return Decision{}, err
	}
	remaining, err := toInt64(values[2])
	if err != nil {
		return Decision{}, err
	}
	retryAfterMS, err := toInt64(values[3])
	if err != nil {
		return Decision{}, err
	}

	return Decision{
		Name:       rule.Name,
		Subject:    subject,
		Allowed:    allowed == 1,
		Limit:      rule.Limit,
		Count:      count,
		Remaining:  remaining,
		RetryAfter: time.Duration(retryAfterMS) * time.Millisecond,
	}, nil
}

func defaultRedisKey(rule Rule, subject string) string {
	return fmt.Sprintf("rl:%s:%s", rule.Name, subject)
}

func toInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("unexpected rate limit value type: %T", value)
	}
}
