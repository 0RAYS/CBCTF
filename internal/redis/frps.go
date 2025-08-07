package redis

import (
	"context"
	"fmt"
	"time"
)

const frpsPortKey = "frps:%s:%s"

func LockFrpsPort(host string, port int32, protocol string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.SAdd(ctx, fmt.Sprintf(frpsPortKey, host, protocol), port).Err()
}

func IsFrpsPortLocked(host string, port int32, protocol string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.SIsMember(ctx, fmt.Sprintf(frpsPortKey, host, protocol), port).Result()
}

func UnlockFrpsPort(host string, port int32, protocol string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.SRem(ctx, fmt.Sprintf(frpsPortKey, host, protocol), port).Err()
}
