package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// setAuthCookie 写入 httpOnly 认证 cookie
// 当请求 Origin 在允许来源列表内（跨域前端）时, 设置 SameSite=None 使浏览器可携带 cookie
// 其余情况保持 SameSite=Lax, 避免无谓降低安全级别
func setAuthCookie(ctx *gin.Context, token string) {
	secure := strings.HasPrefix(config.Env.Host, "https://")
	sameSite := http.SameSiteLaxMode
	origin := ctx.GetHeader("Origin")
	if origin != "" {
		for _, allowed := range config.Env.Gin.Origins {
			if allowed == origin {
				sameSite = http.SameSiteNoneMode
				break
			}
		}
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     middleware.TokenCookieName,
		Value:    token,
		MaxAge:   int(time.Hour.Seconds()),
		Path:     "/",
		Secure:   secure || sameSite == http.SameSiteNoneMode,
		HttpOnly: true,
		SameSite: sameSite,
	})
}

func GetCaptcha(ctx *gin.Context) {
	id, image, answer, err := utils.NewCaptcha().Generate()
	if err != nil {
		log.Logger.Warningf("Failed to generate captcha: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	if ret := redis.SetCaptchaAnswer(id, answer); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"id": id, "image": image}))
}

func Register(ctx *gin.Context) {
	if !config.Env.Registration.Enabled {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.User.NotAllowedRegister})
		return
	}
	var form dto.RegisterForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if ret := redis.VerifyCaptcha(form.CaptchaID, form.Captcha); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.RegisterEventType)
	user, ret := service.RegisterUser(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if ret = service.SendEmail(user); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Name, config.Env.Gin.JWT.Secret)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set("Self", user)
	log.Logger.Infof("%s:%d register", user.Name, user.ID)
	setAuthCookie(ctx, token)
	prometheus.RecordUserRegister(oauth.LocalProvider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetUserResp(service.GetUserView(db.DB, user, false), false)))
}

func Login(ctx *gin.Context) {
	var form dto.LoginForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if form.CaptchaID == "" || form.Captcha == "" {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.User.CaptchaRequired})
		return
	}
	if ret := redis.VerifyCaptcha(form.CaptchaID, form.Captcha); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.LoginEventType)
	user, ret := service.VerifyUser(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Name, config.Env.Gin.JWT.Secret)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set("Self", user)
	log.Logger.Infof("%s:%d login", user.Name, user.ID)
	setAuthCookie(ctx, token)
	prometheus.RecordUserLogin(oauth.LocalProvider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetUserResp(service.GetUserView(db.DB, user, false), false)))
}

func Logout(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.LogoutEventType)
	secure := strings.HasPrefix(config.Env.Host, "https://")
	sameSite := http.SameSiteLaxMode
	origin := ctx.GetHeader("Origin")
	if origin != "" {
		for _, allowed := range config.Env.Gin.Origins {
			if allowed == origin {
				sameSite = http.SameSiteNoneMode
				break
			}
		}
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     middleware.TokenCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   secure || sameSite == http.SameSiteNoneMode,
		HttpOnly: true,
		SameSite: sameSite,
	})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

// ForgotPassword 接收邮箱地址，发送密码重置链接
func ForgotPassword(ctx *gin.Context) {
	var form dto.ForgotPasswordForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if ret := redis.VerifyCaptcha(form.CaptchaID, form.Captcha); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	// 无论邮箱是否存在均返回成功，防止用户枚举
	service.SendPasswordResetEmail(db.DB, form)
	resp.JSON(ctx, model.SuccessRetVal())
}

// ResetPassword 验证重置 token 并设置新密码，同时将邮箱置为已验证
func ResetPassword(ctx *gin.Context) {
	var form dto.ResetPasswordForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		prometheus.RecordPasswordReset(false)
		resp.JSON(ctx, ret)
		return
	}
	ret := service.ResetUserPassword(db.DB, form)
	if !ret.OK {
		prometheus.RecordPasswordReset(false)
		resp.JSON(ctx, ret)
		return
	}
	log.Logger.Info("Password reset via email token")
	prometheus.RecordPasswordReset(true)
	resp.JSON(ctx, model.SuccessRetVal())
}
