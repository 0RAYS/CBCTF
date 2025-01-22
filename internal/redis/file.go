package redis

import (
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"time"
)

func GetFileCache(ctx context.Context, key string) (model.File, bool) {
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
	return file, true
}

func GetFilesCache(ctx context.Context, key string) ([]model.File, bool) {
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
	return files, true
}

func SetFileCache(ctx context.Context, key string, file model.File) error {
	data, err := msgpack.Marshal(file)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func SetFilesCache(ctx context.Context, key string, files []model.File) error {
	data, err := msgpack.Marshal(files)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func DelFileCache(ctx context.Context, id string) error {
	return RDB.Del(ctx, fmt.Sprintf("file:%s", id)).Err()
}

func DelFilesCache(ctx context.Context) error {
	return RDB.Del(ctx, "file:list").Err()
}
