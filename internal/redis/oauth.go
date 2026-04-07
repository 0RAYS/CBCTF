package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const oauthKey = "oauth:%s:%s"
const oauthCodeKey = "oauth:code:%s"

func SetOauthState(provider, state, verifier string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Set(ctx, fmt.Sprintf(oauthKey, provider, state), verifier, 10*time.Minute).Err(); err != nil {
		log.Logger.Warningf("Failed to set oauth state for provider %s: %s", provider, err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthKey, provider, state), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetOauthVerifier(provider, state string) (string, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	verifier, err := RDB.Get(ctx, fmt.Sprintf(oauthKey, provider, state)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", model.RetVal{Msg: i18n.Redis.NotFound, Attr: map[string]any{"Key": fmt.Sprintf(oauthKey, provider, state)}}
		}
		log.Logger.Warningf("Failed to get oauth state for provider %s: %s", provider, err)
		return "", model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthKey, provider, state), "Error": err.Error()}}
	}
	return verifier, model.SuccessRetVal()
}

func DelOauthState(provider string, state string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Del(ctx, fmt.Sprintf(oauthKey, provider, state)).Err(); err != nil {
		log.Logger.Warningf("Failed to delete oauth state for provider %s: %s", provider, err)
		return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": fmt.Sprintf(oauthKey, provider, state), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func SetOauthCode(code, token string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Set(ctx, fmt.Sprintf(oauthCodeKey, code), token, 30*time.Second).Err(); err != nil {
		log.Logger.Warningf("Failed to set oauth code: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthCodeKey, code), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetAndDelOauthToken(code string) (string, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	token, err := RDB.GetDel(ctx, fmt.Sprintf(oauthCodeKey, code)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", model.RetVal{Msg: i18n.Redis.NotFound, Attr: map[string]any{"Key": fmt.Sprintf(oauthCodeKey, code)}}
		}
		log.Logger.Warningf("Failed to get oauth code: %s", err)
		return "", model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthCodeKey, code), "Error": err.Error()}}
	}
	return token, model.SuccessRetVal()
}
