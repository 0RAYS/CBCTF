package ratelimit

import (
	"context"
	"time"
)

type Rule struct {
	Name   string
	Limit  int64
	Window time.Duration
}

type Subject struct {
	Key      string
	ClientIP string
}

type Decision struct {
	Name       string
	Subject    string
	Allowed    bool
	Bypassed   bool
	Limit      int64
	Count      int64
	Remaining  int64
	RetryAfter time.Duration
}

type Store interface {
	Allow(context.Context, Rule, string) (Decision, error)
}

type Limiter struct {
	store     Store
	allowlist Allowlist
}

func New(store Store, allowlist []string) *Limiter {
	return &Limiter{
		store:     store,
		allowlist: NewAllowlist(allowlist),
	}
}

func (l *Limiter) Allow(ctx context.Context, rule Rule, subject Subject) (Decision, error) {
	if rule.Limit <= 0 || rule.Window <= 0 {
		return Decision{
			Name:      rule.Name,
			Subject:   subject.Key,
			Allowed:   true,
			Bypassed:  true,
			Limit:     rule.Limit,
			Remaining: rule.Limit,
		}, nil
	}
	if l.allowlist.Contains(subject.ClientIP) {
		return Decision{
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
