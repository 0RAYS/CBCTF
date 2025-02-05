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

func GetFileCache(key string) (model.Avatar, bool) {
	if !config.Env.Redis.On {
		return model.Avatar{}, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
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

func GetFilesCache(key string) ([]model.Avatar, bool) {
	if !config.Env.Redis.On {
		return nil, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
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

func SetFileCache(key string, file model.Avatar) error {
	if !config.Env.Redis.On {
		return errors.New("redis off")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := msgpack.Marshal(file)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetFileCache: ", file.ID)
	return nil
}

func SetFilesCache(key string, files []model.Avatar) error {
	if !config.Env.Redis.On {
		return errors.New("redis off")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := msgpack.Marshal(files)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetFilesCache: ", len(files))
	return nil
}

func DelFileCache(id string) error {
	if !config.Env.Redis.On {
		return errors.New("redis off")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	key := "file:" + id
	if err := RDB.Del(ctx, key).Err(); err != nil {
		return err
	}
	log.Logger.Debug("DelFileCache: ", id)
	return nil
}

func DelFilesCache() error {
	if !config.Env.Redis.On {
		return errors.New("redis off")
	}
	var cursor uint64
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
		keys, cursor, err := RDB.Scan(ctx, cursor, "files:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan file keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				cancel()
				return err
			}
			log.Logger.Debug("DelFilesCache: ", key)
		}
		cancel()
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelFilesCache")
	return nil
}
