package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"sync/atomic"
	"time"
)

const (
	FilePattern  = "f:%s"
	FilesPattern = "fs"
)

func GetFileCache(id string) (model.Avatar, bool) {
	if !config.Env.Redis.On {
		return model.Avatar{}, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf(FilePattern, id)).Result()
	if errors.Is(err, redis.Nil) {
		atomic.AddInt64(&CacheMiss, 1)
		return model.Avatar{}, false
	} else if err != nil {
		return model.Avatar{}, false
	}
	var file model.Avatar
	err = msgpack.Unmarshal([]byte(data), &file)
	if err != nil {
		return model.Avatar{}, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetFileCache: ", file.ID)
	return file, true
}

func GetFilesCache() ([]model.Avatar, bool) {
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
	var files []model.Avatar
	err = msgpack.Unmarshal([]byte(data), &files)
	if err != nil {
		return nil, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetFilesCache: ", len(files))
	return files, true
}

func SetFileCache(file model.Avatar) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(file)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, fmt.Sprintf(FilePattern, file.ID), data, 1*time.Hour).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetFileCache: ", file.ID)
	return nil
}

func SetFilesCache(files []model.Avatar) error {
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

func DelFileCache(id string) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	if err := RDB.Del(ctx, fmt.Sprintf(FilePattern, id)).Err(); err != nil {
		return err
	}
	log.Logger.Debug("DelFileCache: ", id)
	return nil
}

func DelFilesCache() error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(FilesPattern)
}
