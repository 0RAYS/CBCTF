package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	captchaKey = "captcha:%s"
	captchaTTL = time.Minute
)

var consumeCaptchaScript = redis.NewScript(`
local answer = redis.call("GET", KEYS[1])
if not answer then
  return -1
end
redis.call("DEL", KEYS[1])
if string.lower(answer) == string.lower(ARGV[1]) then
  return 1
end
return 0
`)

func SetCaptchaAnswer(id, answer string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(captchaKey, id)
	if err := RDB.Set(ctx, key, strings.TrimSpace(answer), captchaTTL).Err(); err != nil {
		log.Logger.Warningf("Failed to set captcha answer: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func VerifyCaptcha(id, answer string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(captchaKey, id)
	result, err := consumeCaptchaScript.Run(ctx, RDB, []string{key}, strings.TrimSpace(answer)).Int()
	if err != nil {
		log.Logger.Warningf("Failed to verify captcha answer: %s", err)
		return model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	if result != 1 {
		return model.RetVal{Msg: i18n.Model.User.CaptchaWrong}
	}
	return model.SuccessRetVal()
}
