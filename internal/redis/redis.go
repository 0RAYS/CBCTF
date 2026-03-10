package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Init() {
	addr := fmt.Sprintf("%s:%d", config.Env.Redis.Host, config.Env.Redis.Port)
	log.Logger.Infof("Connecting to Redis: %s", addr)
	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Env.Redis.Pwd,
		DB:           0,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Logger.Warningf("Failed to connect to Redis: %s", err)
		return
	}
	log.Logger.Infof("Connected to Redis: %s", addr)
	// Mirror all logs to Redis list after a successful connection
	log.Logger.AddHook(NewLogHook(5000, log.Formatter{}))
}

func Stop() {
	if err := RDB.Close(); err != nil {
		log.Logger.Warningf("Failed to stop Redis: %s", err)
	}
}

func Count() int64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	count, err := RDB.DBSize(ctx).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get cache total: %s", err)
		return 0
	}
	return count
}
