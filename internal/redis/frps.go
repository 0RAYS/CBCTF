package redis

import (
	"context"
	"fmt"
	"time"
)

func LockFrpsPort(host string, port int32, protocol string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return RDB.SAdd(ctx, fmt.Sprintf("frps:%s:%s", host, protocol), port).Err()
}

func IsFrpsPortLocked(host string, port int32, protocol string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return RDB.SIsMember(ctx, fmt.Sprintf("frps:%s:%s", host, protocol), port).Result()
}

func UnlockFrpsPort(host string, port int32, protocol string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return RDB.SRem(ctx, fmt.Sprintf("frps:%s:%s", host, protocol), port).Err()
}
