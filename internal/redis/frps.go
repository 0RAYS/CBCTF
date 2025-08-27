package redis

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const frpsPortKey = "frps:%s:%s"

var portLock sync.Map

func LockFrpsPort(host string, port int32, protocol string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(frpsPortKey, host, protocol)
	mu, _ := portLock.LoadOrStore(key, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()
	locked, err := RDB.SIsMember(ctx, key, port).Result()
	if err != nil {
		return false, err
	}
	if locked {
		return false, nil
	}
	err = RDB.SAdd(ctx, key, port).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func UnlockFrpsPort(host string, port int32, protocol string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(frpsPortKey, host, protocol)
	mu, _ := portLock.LoadOrStore(key, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()
	locked, err := RDB.SIsMember(ctx, key, port).Result()
	if err != nil {
		return err
	}
	if locked {
		return RDB.SRem(ctx, key, port).Err()
	}
	return nil
}
