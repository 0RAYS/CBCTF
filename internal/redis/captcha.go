package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

const (
	captchaKeyTmpl = "captcha:%s"
	captchaTTL     = time.Minute
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

func SetCaptchaAnswer(ctx context.Context, id, answer string) model.RetVal {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	key := fmt.Sprintf(captchaKeyTmpl, id)
	if err := RDB.Set(ctx, key, strings.TrimSpace(answer), captchaTTL).Err(); err != nil {
		log.Logger.Warningf("Failed to set captcha answer: %s", err)
		return model.FailRetVal(i18n.Redis.SetError, "Cache", key, err)
	}
	return model.SuccessRetVal()
}

func VerifyCaptcha(ctx context.Context, id, answer string) model.RetVal {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	key := fmt.Sprintf(captchaKeyTmpl, id)
	result, err := consumeCaptchaScript.Run(ctx, RDB, []string{key}, strings.TrimSpace(answer)).Int()
	if err != nil {
		log.Logger.Warningf("Failed to verify captcha answer: %s", err)
		return model.FailRetVal(i18n.Redis.GetError, "Cache", key, err)
	}
	if result != 1 {
		return model.FailRetVal(i18n.Model.User.CaptchaWrong, "Captcha", key, nil)
	}
	return model.SuccessRetVal()
}
