package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const TokenCookieName = "CBCTF_TOKEN"

// CheckAuth 是否登录, 用户是否被 ban
func CheckAuth(ctx *gin.Context) {
	var token string
	// 优先从 httpOnly Cookie 读取 token
	if cookie, err := ctx.Cookie(TokenCookieName); err == nil && cookie != "" {
		token = cookie
	} else {
		// 回退到 Authorization 请求头
		auth := strings.Fields(ctx.GetHeader("Authorization"))
		if len(auth) == 2 && auth[0] == "Bearer" {
			token = auth[1]
		}
	}
	if token == "" {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Unauthorized})
		return
	}
	claims, err := utils.ParseToken(token, config.Env.Gin.JWT.Secret)
	if err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Unauthorized})
		return
	}
	user, ret := db.InitUserRepo(db.DB).GetByID(claims.UserID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if user.Banned {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Set("Self", user)
	ctx.Next()
}

// GetSelf 获取当前登录 admin 或 user
func GetSelf(ctx *gin.Context) model.User {
	self, ok := ctx.Get("Self")
	if !ok || self == nil {
		return model.User{}
	}
	return self.(model.User)
}

// CheckCaptain 检查是否为队伍队长, 要求 uri 中必须包含 contestID, admin 路由不能使用
func CheckCaptain(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.CaptainID != GetSelf(ctx).ID {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Next()
}

// CheckVerified 检查邮箱是否已验证
func CheckVerified(ctx *gin.Context) {
	if !GetSelf(ctx).Verified {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Model.User.UnverifiedEmail})
		return
	}
	ctx.Next()
}

// CheckBanned 检查队伍是否被封禁
func CheckBanned(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.Banned {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Next()
}
