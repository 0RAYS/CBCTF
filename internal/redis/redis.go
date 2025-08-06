package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RDB       *redis.Client
	CacheHit  int64
	CacheMiss int64
)

func Init() {
	addr := fmt.Sprintf("%s:%d", config.Env.Redis.Host, config.Env.Redis.Port)
	log.Logger.Infof("Connecting to Redis: %s", addr)
	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Env.Redis.Pwd,
		DB:           0,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Logger.Warningf("Failed to connect to Redis: %s", err)
		return
	}
	log.Logger.Infof("Connected to Redis: %s", addr)

	go StartCollect()
}

func Status() (int64, int64, int64) {
	hit := atomic.LoadInt64(&CacheHit)
	miss := atomic.LoadInt64(&CacheMiss)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	count, err := RDB.DBSize(ctx).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get cache total: %s", err)
		return 0, hit, miss
	}
	return count, hit, miss
}
