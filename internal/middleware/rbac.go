package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckPermission(ctx *gin.Context) {
	key := fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())
	permission, ok := model.RoutePermissions[key]
	if !ok {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	userID := GetSelf(ctx).ID
	pass, ret := redis.CheckUserRBAC(userID, permission)
	if !ret.OK {
		permissions, ret := db.InitPermissionRepo(db.DB).GetUserPermissions(userID)
		if !ret.OK {
			resp.AbortJSON(ctx, ret)
			return
		}
		if ret = redis.SetUserRBAC(userID, permissions); !ret.OK {
			resp.AbortJSON(ctx, ret)
			return
		}
		pass, ret = redis.CheckUserRBAC(userID, permission)
	}
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if !pass {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Next()
}
