package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"sync/atomic"
	"time"
)

const (
	FilesPattern = "fs"
)

func GetFilesCache() ([]model.File, bool) {
	if !config.Env.Redis.On {
		return nil, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, FilesPattern).Result()
	if errors.Is(err, redis.Nil) {
		atomic.AddInt64(&CacheMiss, 1)
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var files []model.File
	err = msgpack.Unmarshal([]byte(data), &files)
	if err != nil {
		return nil, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetFilesCache: ", len(files))
	return files, true
}

func SetFilesCache(files []model.File) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(files)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, FilesPattern, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetFilesCache: ", len(files))
	return nil
}

func DelFilesCache() error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(FilesPattern)
}
