package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Env.GetString("redis.addr"), // Redis 地址
		Password: config.Env.GetString("redis.pwd"),  // 如果没有密码，留空
		DB:       0,                                  // 默认 DB
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Logger.Fatal("Failed to connect to Redis")
	}
}
