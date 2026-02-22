package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/websocket/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resource, operation := strings.Split(permission, ":")[0], strings.Split(permission, ":")[1]
		pass, ret := db.InitPermissionRepo(db.DB).CheckUserPermission(middleware.GetSelfID(ctx), resource, operation)
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
