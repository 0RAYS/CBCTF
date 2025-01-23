package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var RDB *redis.Client

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:         config.Env.Redis.Addr,
		Password:     config.Env.Redis.Pwd,
		DB:           0,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Logger.Error("Failed to connect to Redis")
		return
	}
	log.Logger.Debugf("Connected to Redis: %s", config.Env.Redis.Addr)
}

func Close() {
	if RDB != nil {
		_ = RDB.Close()
	}
	log.Logger.Info("Redis connection closed")
}
