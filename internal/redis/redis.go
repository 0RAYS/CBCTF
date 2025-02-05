package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"github.com/go-redis/redis/v8"
	"sync/atomic"
	"time"
)

// 目前整体来说，缓存只能说勉强能用，设计很有问题
// 数据模型之间相互关联，User Team Contest 嵌套比较多
// 遇到数据更新很多缓存没法及时同步，可能会导致一些功能问题
// 目前是缩短缓存时间来缓解这个问题，后续再说吧，想想更好的解决思路

var (
	RDB           *redis.Client
	CacheHit      int64
	CacheMiss     int64
	collectCTX    context.Context
	collectCancel context.CancelFunc
)

func Init() {
	if !config.Env.Redis.On {
		return
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:         config.Env.Redis.Addr,
		Password:     config.Env.Redis.Pwd,
		DB:           0,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Logger.Errorf("Failed to connect to Redis: %s", err)
		return
	}
	log.Logger.Debugf("Connected to Redis: %s", config.Env.Redis.Addr)

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
	if !config.Env.Redis.On {
		return 0, 0, 0
	}
	hit := atomic.LoadInt64(&CacheHit)
	miss := atomic.LoadInt64(&CacheMiss)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	count, err := RDB.DBSize(ctx).Result()
	if err != nil {
		log.Logger.Error("Failed to get cache total: ", err)
		return 0, hit, miss
	}
	return count, hit, miss
}
