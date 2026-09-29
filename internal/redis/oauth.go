package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

const (
	oauthKeyTmpl     = "oauth:provider:%s:%s"
	oauthCodeKeyTmpl = "oauth:code:%s"
	oauthStateTTL    = 10 * time.Minute
	oauthCodeTTL     = 30 * time.Second
)

func SetOauthState(ctx context.Context, provider, state, verifier string) model.RetVal {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := RDB.Set(ctx, fmt.Sprintf(oauthKeyTmpl, provider, state), verifier, oauthStateTTL).Err(); err != nil {
		log.Logger.Warningf("Failed to set oauth state for provider %s: %s", provider, err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthKeyTmpl, provider, state), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func ConsumeOauthVerifier(ctx context.Context, provider, state string) (string, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	verifier, err := RDB.GetDel(ctx, fmt.Sprintf(oauthKeyTmpl, provider, state)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", model.RetVal{Msg: i18n.Redis.NotFound, Attr: map[string]any{"Key": fmt.Sprintf(oauthKeyTmpl, provider, state)}}
		}
		log.Logger.Warningf("Failed to get oauth state for provider %s: %s", provider, err)
		return "", model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthKeyTmpl, provider, state), "Error": err.Error()}}
	}
	return verifier, model.SuccessRetVal()
}

func SetOauthCode(ctx context.Context, code, token string) model.RetVal {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := RDB.Set(ctx, fmt.Sprintf(oauthCodeKeyTmpl, code), token, oauthCodeTTL).Err(); err != nil {
		log.Logger.Warningf("Failed to set oauth code: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthCodeKeyTmpl, code), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetAndDelOauthToken(ctx context.Context, code string) (string, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	token, err := RDB.GetDel(ctx, fmt.Sprintf(oauthCodeKeyTmpl, code)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", model.RetVal{Msg: i18n.Redis.NotFound, Attr: map[string]any{"Key": fmt.Sprintf(oauthCodeKeyTmpl, code)}}
		}
		log.Logger.Warningf("Failed to get oauth code: %s", err)
		return "", model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": fmt.Sprintf(oauthCodeKeyTmpl, code), "Error": err.Error()}}
	}
	return token, model.SuccessRetVal()
}
