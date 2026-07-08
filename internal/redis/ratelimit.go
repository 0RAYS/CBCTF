package redis

import (
	"CBCTF/internal/config"
	"context"
	"errors"
	"fmt"
	"net/netip"
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

var RateLimiter *Limiter

func InitRateLimiter() {
	RateLimiter = NewLimiter(NewRateLimitRedisStore(func() *redis.Client {
		return RDB
	}), config.Env.Gin.RateLimit.Whitelist)
}

type RateLimitRule struct {
	Name   string
	Limit  int64
	Window time.Duration
}

type RateLimitSubject struct {
	Key      string
	ClientIP string
}

type RateLimitDecision struct {
	Name       string
	Subject    string
	Allowed    bool
	Bypassed   bool
	Limit      int64
	Count      int64
	Remaining  int64
	RetryAfter time.Duration
}

type RateLimitStore interface {
	Allow(context.Context, RateLimitRule, string) (RateLimitDecision, error)
}

type Limiter struct {
	store     RateLimitStore
	allowlist rateLimitAllowlist
}

func NewLimiter(store RateLimitStore, allowlist []string) *Limiter {
	return &Limiter{
		store:     store,
		allowlist: newRateLimitAllowlist(allowlist),
	}
}

func (l *Limiter) Allow(ctx context.Context, rule RateLimitRule, subject RateLimitSubject) (RateLimitDecision, error) {
	if rule.Limit <= 0 || rule.Window <= 0 {
		return RateLimitDecision{
			Name:      rule.Name,
			Subject:   subject.Key,
			Allowed:   true,
			Bypassed:  true,
			Limit:     rule.Limit,
			Remaining: rule.Limit,
		}, nil
	}
	if l.allowlist.Contains(subject.ClientIP) {
		return RateLimitDecision{
			Name:      rule.Name,
			Subject:   subject.Key,
			Allowed:   true,
			Bypassed:  true,
			Limit:     rule.Limit,
			Remaining: rule.Limit,
		}, nil
	}
	return l.store.Allow(ctx, rule, subject.Key)
}

type RateLimitRedisStore struct {
	client  func() *redis.Client
	keyFunc func(RateLimitRule, string) string
	timeout time.Duration
}

func NewRateLimitRedisStore(client func() *redis.Client) *RateLimitRedisStore {
	return &RateLimitRedisStore{
		client:  client,
		keyFunc: defaultRateLimitRedisKey,
		timeout: 3 * time.Second,
	}
}

func (s *RateLimitRedisStore) Allow(ctx context.Context, rule RateLimitRule, subject string) (RateLimitDecision, error) {
	if s.client == nil || s.client() == nil {
		return RateLimitDecision{}, errors.New("rate limit redis client is nil")
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
		return RateLimitDecision{}, err
	}

	values, ok := result.([]any)
	if !ok || len(values) != 4 {
		return RateLimitDecision{}, fmt.Errorf("unexpected rate limit response: %T", result)
	}

	allowed, err := rateLimitValueToInt64(values[0])
	if err != nil {
		return RateLimitDecision{}, err
	}
	count, err := rateLimitValueToInt64(values[1])
	if err != nil {
		return RateLimitDecision{}, err
	}
	remaining, err := rateLimitValueToInt64(values[2])
	if err != nil {
		return RateLimitDecision{}, err
	}
	retryAfterMS, err := rateLimitValueToInt64(values[3])
	if err != nil {
		return RateLimitDecision{}, err
	}

	return RateLimitDecision{
		Name:       rule.Name,
		Subject:    subject,
		Allowed:    allowed == 1,
		Limit:      rule.Limit,
		Count:      count,
		Remaining:  remaining,
		RetryAfter: time.Duration(retryAfterMS) * time.Millisecond,
	}, nil
}

type rateLimitAllowlist struct {
	addrs    map[netip.Addr]struct{}
	prefixes []netip.Prefix
}

func newRateLimitAllowlist(entries []string) rateLimitAllowlist {
	allowlist := rateLimitAllowlist{
		addrs: make(map[netip.Addr]struct{}, len(entries)),
	}
	for _, entry := range entries {
		if prefix, err := netip.ParsePrefix(entry); err == nil {
			allowlist.prefixes = append(allowlist.prefixes, prefix)
			continue
		}
		if addr, err := netip.ParseAddr(entry); err == nil {
			allowlist.addrs[addr.Unmap()] = struct{}{}
		}
	}
	return allowlist
}

func (a rateLimitAllowlist) Contains(rawIP string) bool {
	addr, err := netip.ParseAddr(rawIP)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	if _, ok := a.addrs[addr]; ok {
		return true
	}
	for _, prefix := range a.prefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

func defaultRateLimitRedisKey(rule RateLimitRule, subject string) string {
	return fmt.Sprintf("rl:%s:%s", rule.Name, subject)
}

func rateLimitValueToInt64(value any) (int64, error) {
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
