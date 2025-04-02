package redis

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync/atomic"
	"time"
)

var (
	RDB           *redis.Client
	CacheHit      int64
	CacheMiss     int64
	collectCTX    context.Context
	collectCancel context.CancelFunc
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
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Logger.Errorf("Failed to connect to Redis: %s", err)
		return
	}
	log.Logger.Infof("Connected to Redis: %s", config.Env.Redis.Addr)

	collectCTX, collectCancel = context.WithCancel(context.Background())
	go StartCollect(collectCTX)
}

func Close() {
	if RDB != nil {
		_ = RDB.Close()
	}
	if collectCancel != nil {
		log.Logger.Info("Stop collecting Redis metrics")
		collectCancel()
	}
	log.Logger.Info("Redis connection closed")
}

func Status() (int64, int64, int64) {
	hit := atomic.LoadInt64(&CacheHit)
	miss := atomic.LoadInt64(&CacheMiss)
	ctx := context.Background()
	count, err := RDB.DBSize(ctx).Result()
	if err != nil {
		log.Logger.Error("Failed to get cache total: ", err)
		return 0, hit, miss
	}
	return count, hit, miss
}
