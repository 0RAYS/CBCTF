package redis

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"time"
)

func GetFileCache(key string) (model.File, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.File{}, false
	} else if err != nil {
		return model.File{}, false
	}
	var file model.File
	err = msgpack.Unmarshal([]byte(data), &file)
	if err != nil {
		return model.File{}, false
	}
	log.Logger.Debug("GetFileCache: ", file.ID)
	return file, true
}

func GetFilesCache(key string) ([]model.File, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var files []model.File
	err = msgpack.Unmarshal([]byte(data), &files)
	if err != nil {
		return nil, false
	}
	log.Logger.Debug("GetFilesCache: ", len(files))
	return files, true
}

func SetFileCache(key string, file model.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
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

func SetFilesCache(key string, files []model.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	key := "file:" + id
	if err := RDB.Del(ctx, key).Err(); err != nil {
		return err
	}
	log.Logger.Debug("DelFileCache: ", id)
	return nil
}

func DelFilesCache() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, "files:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan file keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				return err
			}
			log.Logger.Debug("DelFilesCache: ", key)
		}
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelFilesCache")
	return nil
}
