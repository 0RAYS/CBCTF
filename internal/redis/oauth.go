package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"time"
)

const oauthKey = "oauth:%s:%s"

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
