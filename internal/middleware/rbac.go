package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/websocket/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RBAC(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pass, ret := db.InitPermissionRepo(db.DB).CheckUserPermission(middleware.GetSelfID(ctx), permission)
		if !ret.OK {
			ctx.AbortWithStatusJSON(http.StatusOK, ret)
			return
		}
		if !pass {
			ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
			return
		}
		ctx.Next()
	}
}
