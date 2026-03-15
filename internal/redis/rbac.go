package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

const userRBACKey = "users:rbac:%d"

func SetUserRBAC(userID uint, permissions []string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	data, _ := msgpack.Marshal(permissions)
	if err := RDB.Set(ctx, fmt.Sprintf(userRBACKey, userID), data, time.Hour).Err(); err != nil {
		log.Logger.Warningf("Failed to set user RBAC permissions: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": userRBACKey, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func CheckUserRBAC(userID uint, permission string) (bool, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	data, err := RDB.Get(ctx, fmt.Sprintf(userRBACKey, userID)).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get user RBAC permissions: %s", err)
		return false, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": userRBACKey, "Error": err.Error()}}
	}
	var permissions []string
	_ = msgpack.Unmarshal([]byte(data), &permissions)
	return slices.Contains(permissions, permission), model.SuccessRetVal()
}

func DeleteUserRBAC(userID uint) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Del(ctx, fmt.Sprintf(userRBACKey, userID)).Err(); err != nil {
		log.Logger.Warningf("Failed to delete user RBAC permissions: %s", err)
		return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": userRBACKey, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteRBAC() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var cursor uint64
	pattern := strings.TrimSuffix(userRBACKey, "%d") + "*"
	for {
		keys, nextCursor, err := RDB.Scan(ctx, cursor, pattern, 2000).Result()
		if err != nil {
			log.Logger.Warningf("Failed delete RBAC permissions: %s", err)
			return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": userRBACKey, "Error": err.Error()}}
		}

		if len(keys) > 0 {
			pipe := RDB.Pipeline()
			for _, k := range keys {
				pipe.Unlink(ctx, k) // 或 pipe.Del(ctx, k)
			}
			if _, err = pipe.Exec(ctx); err != nil {
				log.Logger.Warningf("Failed delete RBAC permissions: %s", err)
				return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": userRBACKey, "Error": err.Error()}}
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return model.SuccessRetVal()
}
